# Run GitLab Runner in a container--在容器中运行GitLab Runner

This is how you can run GitLab Runner inside a Docker container.
本教程是关于如何在Docker容器中运行GitLab Runner。

## Docker image installation and configuration--安装Docker镜像和配置

1. Install Docker first:
	首先要安装Docker

    ```bash
    curl -sSL https://get.docker.com/ | sh
    ```

1. You need to mount a config volume into the `gitlab-runner` container to
   be used for configs and other resources:
	你需要配置config数据卷挂载到`gitlab-runner`容器的config还有其他资源：
    ```bash
    docker run -d --name gitlab-runner --restart always \
      -v /srv/gitlab-runner/config:/etc/gitlab-runner \
      -v /var/run/docker.sock:/var/run/docker.sock \
      gitlab/gitlab-runner:latest
    ```

    *On OSX, substitute the path "/Users/Shared" for "/srv".*
    	在OSX，/srv替换为"/Users/Shared"

    Or, you can use a config container to mount your custom data volume:
	或者你可以配置以为config容器以挂载自定义的数据卷：
	
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
	如果你计划使用Docker作为运行Runner的方式，你将需要用下面的方式挂载docker socket：
    ```bash
    docker run -d --name gitlab-runner --restart always \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -v /srv/gitlab-runner/config:/etc/gitlab-runner \
      gitlab/gitlab-runner:latest
    ```

1. [Register the Runner](../register/index.md)--注册Runner

Make sure that you read the [FAQ](../faq/README.md) section which describes
some of the most common problems with GitLab Runner.

请您确保已经阅读了FAQ部分，该部分时使用GitLab Runner的常见问题集合。

## Update--更新gitlab-runner

Pull the latest version:
拉取gitlab-runner:

```bash
docker pull gitlab/gitlab-runner:latest
```

停止和移除已存在的容器：
Stop and remove the existing container:

```bash
docker stop gitlab-runner && docker rm gitlab-runner
```

按正常启动该容器：
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
注意：你需要使用原来挂载数据卷的方法来挂载数据卷(`-v /srv/gitlab-runner/config:/etc/gitlab-runner` 或
`--volumes-from gitlab-runner`)。

## Installing trusted SSL server certificates--安装授信的ssl服务证书

If your GitLab CI server is using self-signed SSL certificates then you should
make sure the GitLab CI server certificate is trusted by the gitlab-runner
container for them to be able to talk to each other.

如果你的GitLab CI服务器正在使用自签名SSL证书，那么你应该宝藏GitLab CI服务器证书被也被gitlab-runner信任，以便他们之间可以互相交流。

The `gitlab/gitlab-runner` image is configured to look for the trusted SSL
certificates at `/etc/gitlab-runner/certs/ca.crt`, this can however be changed using the
`-e "CA_CERTIFICATES_PATH=/DIR/CERT"` configuration option.

`gitlab/gitlab-runner`镜像要配置搜索SSL证书的位置为：`/etc/gitlab-runner/certs/ca.crt`，然而这可以使用`-e "CA_CERTIFICATES_PATH=/DIR/CERT"`配置选项来修改。

Copy the `ca.crt` file into the `certs` directory on the data volume (or container).
The `ca.crt` file should contain the root certificates of all the servers you
want gitlab-runner to trust. The gitlab-runner container will
import the `ca.crt` file on startup so if your container is already running you
may need to restart it for the changes to take effect.

拷贝`ca.crt`文件到数据卷或容器的`certs`目录，`ca.crt`文件要包括所有授信gitlab-runner的服务器的根证书。那么gitlab-runner容器将在启动时导入`ca.crt`文件，如果你的容器已经正在运行，你可能要重启系以是更改生效。

## Alpine Linux--

You can also use alternative [Alpine Linux](https://www.alpinelinux.org/) based image with much smaller footprint:
你还可是选择使用基于Alpine Linux的镜像，它比较小：
```
gitlab/gitlab-runner    latest              3e8077e209f5        13 hours ago        304.3 MB
gitlab/gitlab-runner    alpine              7c431ac8f30f        13 hours ago        25.98 MB
```

**Alpine Linux image is designed to use only Docker as the method of spawning runners.**
Alpine Linux 镜像专门为使用docker运行runner设计的。

The original `gitlab/gitlab-runner:latest` is based on Ubuntu 14.04 LTS.
`gitlab/gitlab-runner:latest`最初的版本是基于Ubuntu 14.04 LTS。

## SELinux--

Some distributions (CentOS, RedHat, Fedora) use SELinux by default to enhance the security of the underlying system.

一些Linux发行版，如(CentOS, RedHat, Fedora)默认使用SELinux以加强系统安全。

The special care must be taken when dealing with such configuration.
当处理一下配置时，要注意以下事项：

1. If you want to use Docker executor to run builds in containers you need to access the `/var/run/docker.sock`.
如果你希望用Docker运行器来运行构建，你需要访问`/var/run/docker.sock`。
However, if you have a SELinux in enforcing mode, you will see the `Permission denied` when accessing the `/var/run/docker.sock`.
然而，如果你是使用加强模式的SELinux，你在访问`/var/run/docker.sock`时将看到`Permission denied 限制访问`。
Install the `selinux-dockersock` and to resolve the issue: https://github.com/dpw/selinux-dockersock.
安装`selinux-dockersock`并解决该问题: https://github.com/dpw/selinux-dockersock

1. Make sure that persistent directory is created on host: `mkdir -p /srv/gitlab-runner/config`.
	确保在主机上创建永久目录：`mkdir -p /srv/gitlab-runner/config`。

1. Run docker with `:Z` on volumes:
运行Runner docker命令的加载目录中加上`:Z`

```bash
docker run -d --name gitlab-runner --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /srv/gitlab-runner/config:/etc/gitlab-runner:Z \
  gitlab/gitlab-runner:latest
```

More information about the cause and resolution can be found here:
更多疑问和解决方案可在以下链接中找到：
http://www.projectatomic.io/blog/2015/06/using-volumes-with-docker-can-cause-problems-with-selinux/
