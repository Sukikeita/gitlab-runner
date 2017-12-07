package kubernetes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	clientcmd "k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	clientcmdapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"

	"gitlab.com/gitlab-org/gitlab-runner/common"
)

type kubeConfigProvider func() (*restclient.Config, error)

var (
	// inClusterConfig parses kubernets configuration reading in cluster values
	inClusterConfig kubeConfigProvider = restclient.InClusterConfig
	// defaultKubectlConfig parses kubectl configuration ad loads the default cluster
	defaultKubectlConfig kubeConfigProvider = loadDefaultKubectlConfig
)

func init() {
	clientcmd.DefaultCluster = clientcmdapi.Cluster{}
}

func getKubeClientConfig(config *common.KubernetesConfig, overwrites *overwrites) (kubeConfig *restclient.Config, err error) {
	if len(config.Host) > 0 {
		kubeConfig, err = getOutClusterClientConfig(config)
	} else {
		kubeConfig, err = guessClientConfig()
	}
	if err != nil {
		return nil, err
	}

	//apply overwrites
	kubeConfig.BearerToken = string(overwrites.bearerToken)

	return kubeConfig, nil
}

func getOutClusterClientConfig(config *common.KubernetesConfig) (*restclient.Config, error) {
	kubeConfig := &restclient.Config{
		Host:        config.Host,
		BearerToken: config.BearerToken,
		TLSClientConfig: restclient.TLSClientConfig{
			CAFile: config.CAFile,
		},
	}

	// certificate based auth
	if len(config.CertFile) > 0 {
		if len(config.KeyFile) == 0 || len(config.CAFile) == 0 {
			return nil, fmt.Errorf("ca file, cert file and key file must be specified when using file based auth")
		}

		kubeConfig.TLSClientConfig.CertFile = config.CertFile
		kubeConfig.TLSClientConfig.KeyFile = config.KeyFile
	}

	return kubeConfig, nil
}

func guessClientConfig() (*restclient.Config, error) {
	// Try in cluster config first
	if inClusterCfg, err := inClusterConfig(); err == nil {
		return inClusterCfg, nil
	}

	// in cluster config failed. Reading default kubectl config
	return defaultKubectlConfig()
}

func loadDefaultKubectlConfig() (*restclient.Config, error) {
	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, err
	}

	return clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func getKubeClient(config *common.KubernetesConfig, overwrites *overwrites) (*client.Client, error) {
	restConfig, err := getKubeClientConfig(config, overwrites)
	if err != nil {
		return nil, err
	}

	return client.New(restConfig)
}

func closeKubeClient(client *client.Client) bool {
	if client == nil || client.Client == nil || client.Client.Transport == nil {
		return false
	}
	if transport, _ := client.Client.Transport.(*http.Transport); transport != nil {
		transport.CloseIdleConnections()
		return true
	}
	return false
}

func isRunning(pod *api.Pod) (bool, error) {
	switch pod.Status.Phase {
	case api.PodRunning:
		return true, nil
	case api.PodSucceeded:
		return false, fmt.Errorf("pod already succeeded before it begins running")
	case api.PodFailed:
		return false, fmt.Errorf("pod status is failed")
	default:
		return false, nil
	}
}

type podPhaseResponse struct {
	done  bool
	phase api.PodPhase
	err   error
}

func getPodPhase(c *client.Client, pod *api.Pod, out io.Writer) podPhaseResponse {
	pod, err := c.Pods(pod.Namespace).Get(pod.Name)
	if err != nil {
		return podPhaseResponse{true, api.PodUnknown, err}
	}

	ready, err := isRunning(pod)

	if err != nil {
		return podPhaseResponse{true, pod.Status.Phase, err}
	}

	if ready {
		return podPhaseResponse{true, pod.Status.Phase, nil}
	}

	// check status of containers
	for _, container := range pod.Status.ContainerStatuses {
		if container.Ready {
			continue
		}
		if container.State.Waiting == nil {
			continue
		}

		switch container.State.Waiting.Reason {
		case "ErrImagePull", "ImagePullBackOff":
			err = fmt.Errorf("image pull failed: %s", container.State.Waiting.Message)
			err = &common.BuildError{Inner: err}
			return podPhaseResponse{true, api.PodUnknown, err}
		}
	}

	fmt.Fprintf(out, "Waiting for pod %s/%s to be running, status is %s\n", pod.Namespace, pod.Name, pod.Status.Phase)
	return podPhaseResponse{false, pod.Status.Phase, nil}

}

func triggerPodPhaseCheck(c *client.Client, pod *api.Pod, out io.Writer) <-chan podPhaseResponse {
	errc := make(chan podPhaseResponse)
	go func() {
		defer close(errc)
		errc <- getPodPhase(c, pod, out)
	}()
	return errc
}

// waitForPodRunning will use client c to detect when pod reaches the PodRunning
// state. It returns the final PodPhase once either PodRunning, PodSucceeded or
// PodFailed has been reached. In the case of PodRunning, it will also wait until
// all containers within the pod are also Ready.
// It returns error if the call to retrieve pod details fails or the timeout is
// reached.
// The timeout and polling values are configurable through KubernetesConfig
// parameters.
func waitForPodRunning(ctx context.Context, c *client.Client, pod *api.Pod, out io.Writer, config *common.KubernetesConfig) (api.PodPhase, error) {
	pollInterval := config.GetPollInterval()
	pollAttempts := config.GetPollAttempts()
	for i := 0; i <= pollAttempts; i++ {
		select {
		case r := <-triggerPodPhaseCheck(c, pod, out):
			if !r.done {
				time.Sleep(time.Duration(pollInterval) * time.Second)
				continue
			}
			return r.phase, r.err
		case <-ctx.Done():
			return api.PodUnknown, ctx.Err()
		}
	}
	return api.PodUnknown, errors.New("timedout waiting for pod to start")
}

// limits takes a string representing CPU & memory limits,
// and returns a ResourceList with appropriately scaled Quantity
// values for Kubernetes. This allows users to write "500m" for CPU,
// and "50Mi" for memory (etc.)
func limits(cpu, memory string) (api.ResourceList, error) {
	var rCPU, rMem resource.Quantity
	var err error

	parse := func(s string) (resource.Quantity, error) {
		var q resource.Quantity
		if len(s) == 0 {
			return q, nil
		}
		if q, err = resource.ParseQuantity(s); err != nil {
			return q, fmt.Errorf("error parsing resource limit: %s", err.Error())
		}
		return q, nil
	}

	if rCPU, err = parse(cpu); err != nil {
		return api.ResourceList{}, nil
	}

	if rMem, err = parse(memory); err != nil {
		return api.ResourceList{}, nil
	}

	l := make(api.ResourceList)

	q := resource.Quantity{}
	if rCPU != q {
		l[api.ResourceCPU] = rCPU
	}
	if rMem != q {
		l[api.ResourceMemory] = rMem
	}

	return l, nil
}

// buildVariables converts a common.BuildVariables into a list of
// kubernetes EnvVar objects
func buildVariables(bv common.JobVariables) []api.EnvVar {
	e := make([]api.EnvVar, len(bv))
	for i, b := range bv {
		e[i] = api.EnvVar{
			Name:  b.Key,
			Value: b.Value,
		}
	}
	return e
}
