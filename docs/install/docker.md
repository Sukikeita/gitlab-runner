# Run GitLab Runner in a container--������������GitLab Runner

This is how you can run GitLab Runner inside a Docker container.
���̳��ǹ��������Docker����������GitLab Runner��

## Docker image installation and configuration--��װDocker���������

1. Install Docker first:
	����Ҫ��װDocker

    ```bash
    curl -sSL https://get.docker.com/ | sh
    ```

1. You need to mount a config volume into the `gitlab-runner` container to
   be used for configs and other resources:
	����Ҫ����config���ݾ���ص�`gitlab-runner`������config����������Դ��
    ```bash
    docker run -d --name gitlab-runner --restart always \
      -v /srv/gitlab-runner/config:/etc/gitlab-runner \
      -v /var/run/docker.sock:/var/run/docker.sock \
      gitlab/gitlab-runner:latest
    ```

    *On OSX, substitute the path "/Users/Shared" for "/srv".*
    	��OSX��/srv�滻Ϊ"/Users/Shared"

    Or, you can use a config container to mount your custom data volume:
	���������������Ϊconfig�����Թ����Զ�������ݾ�
	
    ```bash
    docker run -d --name gitlab-runner-config \
        -v /etc/gitlab-runner \
        busybox:latest \
        /bin/true

    docker run -d --name gitlab-runner --restart always \
        --volumes-from gitlab-runner-config \
        gitlab/gitlab-runner:latest
    ```

    If you plan on using Docker as the method of spawning Runners, you will need to
    mount your docker socket like this:
	�����ƻ�ʹ��Docker��Ϊ����Runner�ķ�ʽ���㽫��Ҫ������ķ�ʽ����docker socket��
    ```bash
    docker run -d --name gitlab-runner --restart always \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -v /srv/gitlab-runner/config:/etc/gitlab-runner \
      gitlab/gitlab-runner:latest
    ```

1. [Register the Runner](../register/index.md)--ע��Runner

Make sure that you read the [FAQ](../faq/README.md) section which describes
some of the most common problems with GitLab Runner.

����ȷ���Ѿ��Ķ���FAQ���֣��ò���ʱʹ��GitLab Runner�ĳ������⼯�ϡ�

## Update--����gitlab-runner

Pull the latest version:
��ȡgitlab-runner:

```bash
docker pull gitlab/gitlab-runner:latest
```

ֹͣ���Ƴ��Ѵ��ڵ�������
Stop and remove the existing container:

```bash
docker stop gitlab-runner && docker rm gitlab-runner
```

������������������
Start the container as you did originally:

```bash
docker run -d --name gitlab-runner --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /srv/gitlab-runner/config:/etc/gitlab-runner \
  gitlab/gitlab-runner:latest
```

>**Note**:
you need to use the same method for mounting you data volume as you
did originally (`-v /srv/gitlab-runner/config:/etc/gitlab-runner` or
`--volumes-from gitlab-runner`).
ע�⣺����Ҫʹ��ԭ���������ݾ�ķ������������ݾ�(`-v /srv/gitlab-runner/config:/etc/gitlab-runner` ��
`--volumes-from gitlab-runner`)��

## Installing trusted SSL server certificates--��װ���ŵ�ssl����֤��

If your GitLab CI server is using self-signed SSL certificates then you should
make sure the GitLab CI server certificate is trusted by the gitlab-runner
container for them to be able to talk to each other.

������GitLab CI����������ʹ����ǩ��SSL֤�飬��ô��Ӧ�ñ���GitLab CI������֤�鱻Ҳ��gitlab-runner���Σ��Ա�����֮����Ի��ཻ����

The `gitlab/gitlab-runner` image is configured to look for the trusted SSL
certificates at `/etc/gitlab-runner/certs/ca.crt`, this can however be changed using the
`-e "CA_CERTIFICATES_PATH=/DIR/CERT"` configuration option.

`gitlab/gitlab-runner`����Ҫ��������SSL֤���λ��Ϊ��`/etc/gitlab-runner/certs/ca.crt`��Ȼ�������ʹ��`-e "CA_CERTIFICATES_PATH=/DIR/CERT"`����ѡ�����޸ġ�

Copy the `ca.crt` file into the `certs` directory on the data volume (or container).
The `ca.crt` file should contain the root certificates of all the servers you
want gitlab-runner to trust. The gitlab-runner container will
import the `ca.crt` file on startup so if your container is already running you
may need to restart it for the changes to take effect.

����`ca.crt`�ļ������ݾ��������`certs`Ŀ¼��`ca.crt`�ļ�Ҫ������������gitlab-runner�ķ������ĸ�֤�顣��ôgitlab-runner������������ʱ����`ca.crt`�ļ��������������Ѿ��������У������Ҫ����ϵ���Ǹ�����Ч��

## Alpine Linux--

You can also use alternative [Alpine Linux](https://www.alpinelinux.org/) based image with much smaller footprint:
�㻹����ѡ��ʹ�û���Alpine Linux�ľ������Ƚ�С��
```
gitlab/gitlab-runner    latest              3e8077e209f5        13 hours ago        304.3 MB
gitlab/gitlab-runner    alpine              7c431ac8f30f        13 hours ago        25.98 MB
```

**Alpine Linux image is designed to use only Docker as the method of spawning runners.**
Alpine Linux ����ר��Ϊʹ��docker����runner��Ƶġ�

The original `gitlab/gitlab-runner:latest` is based on Ubuntu 14.04 LTS.
`gitlab/gitlab-runner:latest`����İ汾�ǻ���Ubuntu 14.04 LTS��

## SELinux--

Some distributions (CentOS, RedHat, Fedora) use SELinux by default to enhance the security of the underlying system.

һЩLinux���а棬��(CentOS, RedHat, Fedora)Ĭ��ʹ��SELinux�Լ�ǿϵͳ��ȫ��

The special care must be taken when dealing with such configuration.
������һ������ʱ��Ҫע���������

1. If you want to use Docker executor to run builds in containers you need to access the `/var/run/docker.sock`.
�����ϣ����Docker�����������й���������Ҫ����`/var/run/docker.sock`��
However, if you have a SELinux in enforcing mode, you will see the `Permission denied` when accessing the `/var/run/docker.sock`.
Ȼ�����������ʹ�ü�ǿģʽ��SELinux�����ڷ���`/var/run/docker.sock`ʱ������`Permission denied ���Ʒ���`��
Install the `selinux-dockersock` and to resolve the issue: https://github.com/dpw/selinux-dockersock.
��װ`selinux-dockersock`�����������: https://github.com/dpw/selinux-dockersock

1. Make sure that persistent directory is created on host: `mkdir -p /srv/gitlab-runner/config`.
	ȷ���������ϴ�������Ŀ¼��`mkdir -p /srv/gitlab-runner/config`��

1. Run docker with `:Z` on volumes:
����Runner docker����ļ���Ŀ¼�м���`:Z`

```bash
docker run -d --name gitlab-runner --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /srv/gitlab-runner/config:/etc/gitlab-runner:Z \
  gitlab/gitlab-runner:latest
```

More information about the cause and resolution can be found here:
�������ʺͽ���������������������ҵ���
http://www.projectatomic.io/blog/2015/06/using-volumes-with-docker-can-cause-problems-with-selinux/
