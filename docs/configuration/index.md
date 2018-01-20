---
comments: false
---

# Configuring GitLab Runner

Below you can find some specific documentation on configuring GitLab Runner, the
shells supported, the security implications using the various executors, as
well as information how to set up Prometheus metrics:

您可以在下面找到关于配置GitLab Runner的一些特定文档，支持的shell，使用各种执行程序的安全含义，以及如何设置Prometheus指标的信息：


- [Advanced configuration options](advanced-configuration.md) Learn how to use the [TOML][] configuration file that GitLab Runner uses.
学习如何使用GitLab Runner使用的TOML配置文件。

- [Use self-signed certificates](tls-self-signed.md) Configure certificates that are used to verify TLS peer when connecting to the GitLab server.
连接到GitLab服务器时，配置用于验证TLS对等体的证书。

- [Auto-scaling using Docker machine](autoscale.md) Execute jobs on machines that are created on demand using Docker machine.
在使用Docker机器按需创建的机器上执行作业。

- [Supported shells](../shells/README.md) Learn what shell script generators are supported that allow to execute builds on different systems.
了解哪些shell脚本生成器支持允许在不同系统上执行构建。

- [Security considerations](../security/index.md) Be aware of potential security implications when running your jobs with GitLab Runner.
使用GitLab Runner运行作业时，请注意潜在的安全隐患。

- [Prometheus monitoring](../monitoring/README.md) Learn how to use the Prometheus metrics HTTP server.
了解如何使用Prometheus指标HTTP服务器。

- [Cleanup the Docker images automatically](https://gitlab.com/gitlab-org/gitlab-runner-docker-cleanup) A simple Docker application that automatically garbage collects the GitLab Runner caches and images when running low on disk space.
一个简单的Docker应用程序，在磁盘空间不足时自动收集GitLab Runner缓存和映像。


[TOML]: https://github.com/toml-lang/toml
