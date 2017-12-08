---
toc: false
comments: false
last_updated: 2017-10-09
---

# GitLab Runner

[![Build Status](https://gitlab.com/gitlab-org/gitlab-runner/badges/master/build.svg)](https://gitlab.com/gitlab-org/gitlab-runner)

GitLab Runner is the open source project that is used to run your jobs and
send the results back to GitLab. It is used in conjunction with [GitLab CI][ci],
the open-source continuous integration service included with GitLab that
coordinates the jobs.

GitLab Runner是一个开源项目，用于运行您的作业并将结果发回给GitLab。 它与GitLab CI（GitLab中包含的用于协调作业的开源持续集成服务）结合使用。


## Requirements--要求

GitLab Runner is written in [Go][golang] and can be run as a single binary, no
language specific requirements are needed.

Runner是用Go语言编写的，并可以以单个二进制文件来运行，因此无需安装特殊的语言包。

It is designed to run on the GNU/Linux, macOS, and Windows operating systems.
Other operating systems will probably work as long as you can compile a Go
binary on them.

它专门为Linux、macOS和window系统而设计，也可以运行在其他系统，只要你可以使用Go来编译出二进制包。

If you want to use Docker make sure that you have version `v1.5.0` at least
installed.
如果你想使用Docker，请确保安装`v1.5.0`及以上的版本。

## Features--Runner的功能

- Allows to run:支持以下功能
 - multiple jobs concurrently 同时运行多个作业
 - use multiple tokens with multiple server (even per-project) 不同服务器（甚至是项目）使用不同的token
 - limit number of concurrent jobs per-token 每个token能同时运行的作业有限
- Jobs can be run:作业可支持
 - locally 本地运行
 - using Docker containers 使用Docker容器运行
 - using Docker containers and executing job over SSH 使用Docker容器并通过SSH执行作业
 - using Docker containers with autoscaling on different clouds and virtualization hypervisors 在不同的云和虚拟化管理程序上使用Docker容器进行自动缩放
 - connecting to remote SSH server 可连接到远程SSH服务器
- Is written in Go and distributed as single binary without any other requirements 用Go语言编写，并以单个二进制包发行而无需其他要求
- Supports Bash, Windows Batch and Windows PowerShell 支持Bash、windows批处理和powerShell命令
- Works on GNU/Linux, OS X and Windows (pretty much anywhere you can run Docker) 支持在LInus、OS X和Windows上运行（任何可以运行Docker的地方）
- Allows to customize the job running environment 迟滞自定义作业的运行环境
- Automatic configuration reload without restart 自动重载配置而无需重启
- Easy to use setup with support for Docker, Docker-SSH, Parallels or SSH running environments 容易使用
- Enables caching of Docker containers 启动Docker容器缓存
- Easy installation as a service for GNU/Linux, OSX and Windows 可以安装为服务
- Embedded Prometheus metrics HTTP server 嵌入Prometheus指标HTTP服务器

## Compatibility chart--兼容性表

CAUTION: **Important:**
GitLab Runner >= 9.0 requires GitLab's API v4 endpoints, which were introduced
in GitLab CE/EE 9.0. Because of this change, GitLab Runner >= 9.0 requires
GitLab CE/EE >= 9.0 and will not work with older GitLab versions.
The old API used by GitLab Runner was deprecated in August 2017 and with this,
the v1.11.x version of GitLab Runner is deprecated as well.

In the following table you can see the compatibility chart between GitLab and
GitLab Runner.

|GitLab Runner / GitLab | 9.0.x (03.2017) | 9.1.x (04.2017) | 9.2.x (05.2017) | 9.3.x (06.2017) | 9.4.x (07.2017) | 9.5.x (08.2017) | 10.0.x (09.2017) |
|:---------------------:|:---------------:|:---------------:|:---------------:|:---------------:|:---------------:|:---------------:|:----------------:|
| v1.10.x  |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v1.11.x  |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v9.0.x   |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v9.1.x   |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v9.2.x   |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v9.3.x   |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v9.4.x   |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v9.5.x   |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |
| v10.0.x  |  鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�               | 鉁�                |

## [Install GitLab Runner](install/index.md)--安装GitLab-Runner

GitLab Runner can be installed and used on GNU/Linux, macOS, FreeBSD and Windows.
You can install it using Docker, download the binary manually or use the
repository for rpm/deb packages that GitLab offers. Below you can find
information on the different installation methods:

GitLab Runner可以在GNU / Linux，macOS，FreeBSD和Windows上安装和使用。 您可以使用Docker进行安装，手动下载二进制文件或使用GitLab提供的rpm/deb软件包的存储库。 以下您可以找到有关不同安装方法的信息：

- [Install using GitLab's repository for Debian/Ubuntu/CentOS/RedHat (preferred)](install/linux-repository.md)
- [Install on GNU/Linux manually (advanced)](install/linux-manually.md)
- [Install on macOS (preferred)](install/osx.md)
- [Install on Windows (preferred)](install/windows.md)
- [Install as a Docker Service](install/docker.md)
- [Install in Auto-scaling mode using Docker machine](install/autoscaling.md)
- [Install on FreeBSD](install/freebsd.md)
- [Install on Kubernetes](install/kubernetes.md)
- [Install the nightly binary manually (development)](install/bleeding-edge.md)

## [Register GitLab Runner](register/index.md)--注册Runner

Once GitLab Runner is installed, you need to register it with GitLab.
一旦安装了Runner后，你需要将其余GitLab注册。

Learn how to [register a GitLab Runner](register/index.md).
学习如何注册Runner。

## Using GitLab Runner--使用GitLab-Runner

- [See the commands documentation](commands/README.md)请查看commands目录下的readme文档。

## [Selecting the executor](executors/README.md)--如何选择executor

GitLab Runner implements a number of executors that can be used to run your
builds in different scenarios. If you are not sure what to select, read the
[I am not sure](executors/README.md#i-am-not-sure) section.
Visit the [compatibility chart](executors/README.md#compatibility-chart) to find
out what features each executor supports and what not.

GitLab Runner实现了许多可用于在不同场景下运行构建的执行程序（executors）。 如果您不确定要选择什么，请阅读我不确定的部分。 访问兼容性图表，找出每个执行程序支持哪些功能，哪些不支持。

To jump into the specific documentation of each executor, visit:
要进入每个executor的具体文档，请访问：

- [Shell](executors/shell.md)
- [Docker](executors/docker.md)
- [Docker Machine and Docker Machine SSH (auto-scaling)](install/autoscaling.md)
- [Parallels](executors/parallels.md)
- [VirtualBox](executors/virtualbox.md)
- [SSH](executors/ssh.md)
- [Kubernetes](executors/kubernetes.md)

## [Advanced Configuration](configuration/index.md)--高级配置

- [Advanced configuration options](configuration/advanced-configuration.md) Learn how to use the [TOML][] configuration file that GitLab Runner uses.--学习如何使用TOML配置文件
- [Use self-signed certificates](configuration/tls-self-signed.md) Configure certificates that are used to verify TLS peer when connecting to the GitLab server.配置用于验证TLS对（当连接到GitLab服务器）的证书，
- [Auto-scaling using Docker machine](configuration/autoscale.md) Execute jobs on machines that are created on demand using Docker machine.使用Docker机器进行自动缩放，在使用Docker机器按需创建的机器上执行作业。
- [Supported shells](shells/README.md) Learn what shell script generators are supported that allow to execute builds on different systems.
- [Security considerations](security/index.md) Be aware of potential security implications when running your jobs with GitLab Runner.安全考虑：使用GitLab Runner运行作业时，请注意潜在的安全隐患。
- [Runner monitoring](monitoring/README.md) Learn how to monitor Runner's behavior. Runner监控--学习如何监控Runner的行为
- [Cleanup the Docker images automatically]--自动清除Docker镜像(https://gitlab.com/gitlab-org/gitlab-runner-docker-cleanup) A simple Docker application that automatically garbage collects the GitLab Runner caches and images when running low on disk space.一个简单的Docker应用程序，在磁盘空间不足时自动收集GitLab Runner缓存和映像。

## Troubleshooting--故障排除

Read the [FAQ](faq/README.md) for troubleshooting common issues.

## Release process

The description of release process of the GitLab Runner project can be found in
the [release documentation](release_process/README.md).

## Contributing--欢迎贡献

Contributions are welcome, see [`CONTRIBUTING.md`][contribute] for more details.

## Development--GitLab Runner的发展史

See the [development documentation](development/README.md) to hack on GitLab
Runner.

## Changelog--变更史

Visit [Changelog] to view recent changes.

## License--证书

This code is distributed under the MIT license, see the [LICENSE][] file.

[ci]: https://about.gitlab.com/gitlab-ci
[Changelog]: https://gitlab.com/gitlab-org/gitlab-runner/blob/master/CHANGELOG.md
[contribute]: https://gitlab.com/gitlab-org/gitlab-runner/blob/master/CONTRIBUTING.md
[golang]: https://golang.org/
[LICENSE]: https://gitlab.com/gitlab-org/gitlab-runner/blob/master/LICENSE
[TOML]: https://github.com/toml-lang/toml
