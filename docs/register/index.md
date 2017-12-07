# Registering Runners--注册Runner

Registering a Runner is the process that binds the Runner with a GitLab instance.

注册Runner的作用就是将Runner绑定到GitLab实例。

## Prerequisites--前提准备

Before registering a Runner, you need to first:
在注册Runner之前，你首先需要：

- [Install it](../install/index.md) on a server separate than where GitLab
  is installed on--安装GitLab-Runner，并且与GitLab安装位置隔离；
- [Obtain a token](https://docs.gitlab.com/ce/ci/runners/) for a shared or
  specific Runner via GitLab's interface--通过GitLab接口获取一个token

## GNU/Linux

To register a Runner under GNU/Linux:
在GNU/Linux中注册Runner的步骤：

1. Run the following command:
运行以下命令：

    ```sh
    sudo gitlab-runner register
    ```

1. Enter your GitLab instance URL:
输入GitLab实例的URL：

    ```
    Please enter the gitlab-ci coordinator URL (e.g. https://gitlab.com )
    https://gitlab.com
    ```

1. Enter the token you obtained to register the Runner:
输入token以注册Runner：

    ```
    Please enter the gitlab-ci token for this runner
    xxx
    ```

1. Enter a description for the Runner, you can change this later in GitLab's
   UI:
   输入一段Runner描述，后续可以在GitLab界面上修改：

    ```
    Please enter the gitlab-ci description for this runner
    [hostame] my-runner
    ```

1. Enter the [tags associated with the Runner][tags], you can change this later in GitLab's UI:
输入关联Runner的标签，后续可以在GitLab界面中修改：

    ```
    Please enter the gitlab-ci tags for this runner (comma separated):
    my-tag,another-tag
    ```

1. Choose whether the Runner should pick up jobs that do not [have tags][tags],
   you can change this later in GitLab's UI (defaults to false):
选择该Runner是否选择没有标签的作业，后续可以在GitLab界面中修改（默认为false）：

    ```
    Whether to run untagged jobs [true/false]:
    [false]: true
    ```

1. Choose whether to lock the Runner to the current project, you can change
   this later in GitLab's UI. Useful when the Runner is specific (defaults to
   true):
选择是否将改Runner锁定到当前项目，后续在GitLab界面中修改。当Runner是特定的时候很有用（默认为true）：
    ```
    Whether to lock Runner to current project [true/false]:
    [true]: true
    ```

1. Enter the [Runner executor](../executors/README.md):
输入Runner执行器：
    ```
    Please enter the executor: ssh, docker+machine, docker-ssh+machine, kubernetes, docker, parallels, virtualbox, docker-ssh, shell:
    docker
    ```

1. If you chose Docker as your executor, you'll be asked for the default
   image to be used for projects that do not define one in `.gitlab-ci.yml`:
如果你选择Docker作为你的执行器，你将被询问：当项目的`.gitlab-ci.yml`没有定义executor时，要使用的默认镜像：

    ```
    Please enter the Docker image (eg. ruby:2.1):
    alpine:latest
    ```

## macOS

To register a Runner under macOS:
在macOS系统注册Runner

1. Run the following command:
运行以下命令：

    ```sh
    gitlab-runner register
    ```

1. Enter your GitLab instance URL:
输入GitLab实例的URL：

    ```
    Please enter the gitlab-ci coordinator URL (e.g. https://gitlab.com )
    https://gitlab.com
    ```

1. Enter the token you obtained to register the Runner:
输入token以注册Runner：

    ```
    Please enter the gitlab-ci token for this runner
    xxx
    ```

1. Enter a description for the Runner, you can change this later in GitLab's
   UI:
输入Runner描述：

    ```
    Please enter the gitlab-ci description for this runner
    [hostame] my-runner
    ```

1. Enter the [tags associated with the Runner][tags], you can change this later in GitLab's UI:
输入标签以关联Runner，后续可修改

    ```
    Please enter the gitlab-ci tags for this runner (comma separated):
    my-tag,another-tag
    ```

1. Choose whether the Runner should pick up jobs that do not [have tags][tags],
   you can change this later in GitLab's UI (defaults to false):

    ```
    Whether to run untagged jobs [true/false]:
    [false]: true
    ```

1. Choose whether to lock the Runner to the current project, you can change
   this later in GitLab's UI. Useful when the Runner is specific (defaults to
   true):

    ```
    Whether to lock Runner to current project [true/false]:
    [true]: true
    ```

1. Enter the [Runner executor](../executors/README.md):

    ```
    Please enter the executor: ssh, docker+machine, docker-ssh+machine, kubernetes, docker, parallels, virtualbox, docker-ssh, shell:
    docker
    ```

1. If you chose Docker as your executor, you'll be asked for the default
   image to be used for projects that do not define one in `.gitlab-ci.yml`:

    ```
    Please enter the Docker image (eg. ruby:2.1):
    alpine:latest
    ```
    **Note** _[be sure Docker.app is installed on your mac](https://docs.docker.com/docker-for-mac/install/)_

## Windows--在windows注册Runner

To register a Runner under Windows:

1. Run the following command:

    ```sh
    ./gitlab-runner.exe register
    ```

1. Enter your GitLab instance URL:

    ```
    Please enter the gitlab-ci coordinator URL (e.g. https://gitlab.com )
    https://gitlab.com
    ```

1. Enter the token you obtained to register the Runner:

    ```
    Please enter the gitlab-ci token for this runner
    xxx
    ```

1. Enter a description for the Runner, you can change this later in GitLab's
   UI:

    ```
    Please enter the gitlab-ci description for this runner
    [hostame] my-runner
    ```

1. Enter the [tags associated with the Runner][tags], you can change this later in GitLab's UI:

    ```
    Please enter the gitlab-ci tags for this runner (comma separated):
    my-tag,another-tag
    ```

1. Choose whether the Runner should pick up jobs that do not [have tags][tags],
   you can change this later in GitLab's UI (defaults to false):

    ```
    Whether to run untagged jobs [true/false]:
    [false]: true
    ```

1. Choose whether to lock the Runner to the current project, you can change
   this later in GitLab's UI. Useful when the Runner is specific (defaults to
   true):

    ```
    Whether to lock Runner to current project [true/false]:
    [true]: true
    ```

1. Enter the [Runner executor](../executors/README.md):

    ```
    Please enter the executor: ssh, docker+machine, docker-ssh+machine, kubernetes, docker, parallels, virtualbox, docker-ssh, shell:
    docker
    ```

1. If you chose Docker as your executor, you'll be asked for the default
   image to be used for projects that do not define one in `.gitlab-ci.yml`:

    ```
    Please enter the Docker image (eg. ruby:2.1):
    alpine:latest
    ```

If you'd like to register multiple Runners on the same machine with different
configurations repeat the `./gitlab-runner.exe register` command.

## FreeBSD

To register a Runner under FreeBSD:

1. Run the following command:

    ```sh
    sudo -u gitlab-runner -H /usr/local/bin/gitlab-runner register
    ```

1. Enter your GitLab instance URL:

    ```
    Please enter the gitlab-ci coordinator URL (e.g. https://gitlab.com )
    https://gitlab.com
    ```

1. Enter the token you obtained to register the Runner:

    ```
    Please enter the gitlab-ci token for this runner
    xxx
    ```

1. Enter a description for the Runner, you can change this later in GitLab's
   UI:

    ```
    Please enter the gitlab-ci description for this runner
    [hostame] my-runner
    ```

1. Enter the [tags associated with the Runner][tags], you can change this later in GitLab's UI:

    ```
    Please enter the gitlab-ci tags for this runner (comma separated):
    my-tag,another-tag
    ```

1. Choose whether the Runner should pick up jobs that do not [have tags][tags],
   you can change this later in GitLab's UI (defaults to false):

    ```
    Whether to run untagged jobs [true/false]:
    [false]: true
    ```

1. Choose whether to lock the Runner to the current project, you can change
   this later in GitLab's UI. Useful when the Runner is specific (defaults to
   false):

    ```
    Whether to lock Runner to current project [true/false]:
    [true]: true
    ```

1. Enter the [Runner executor](../executors/README.md):

    ```
    Please enter the executor: ssh, docker+machine, docker-ssh+machine, kubernetes, docker, parallels, virtualbox, docker-ssh, shell:
    docker
    ```

1. If you chose Docker as your executor, you'll be asked for the default
   image to be used for projects that do not define one in `.gitlab-ci.yml`:

    ```
    Please enter the Docker image (eg. ruby:2.1):
    alpine:latest
    ```

## Docker--使用Docker容器注册Runner

To register a Runner using a Docker container:
要在Docker容器中注册Runner：

1. Run the following command:
运行以下命令：

    ```sh
    docker exec -it gitlab-runner gitlab-runner register
    ```

1. Enter your GitLab instance URL:
输入GitLab实例的URL
    ```
    Please enter the gitlab-ci coordinator URL (e.g. https://gitlab.com )
    https://gitlab.com
    ```

1. Enter the token you obtained to register the Runner:
输入token
    ```
    Please enter the gitlab-ci token for this runner
    xxx
    ```

1. Enter a description for the Runner, you can change this later in GitLab's
   UI:

    ```
    Please enter the gitlab-ci description for this runner
    [hostame] my-runner
    ```

1. Enter the [tags associated with the Runner][tags], you can change this later in GitLab's UI:

    ```
    Please enter the gitlab-ci tags for this runner (comma separated):
    my-tag,another-tag
    ```

1. Choose whether the Runner should pick up jobs that do not [have tags][tags],
   you can change this later in GitLab's UI (defaults to false):

    ```
    Whether to run untagged jobs [true/false]:
    [false]: true
    ```

1. Choose whether to lock the Runner to the current project, you can change
   this later in GitLab's UI. Useful when the Runner is specific (defaults to
   true):

    ```
    Whether to lock Runner to current project [true/false]:
    [true]: true
    ```

1. Enter the [Runner executor](../executors/README.md):

    ```
    Please enter the executor: ssh, docker+machine, docker-ssh+machine, kubernetes, docker, parallels, virtualbox, docker-ssh, shell:
    docker
    ```

1. If you chose Docker as your executor, you'll be asked for the default
   image to be used for projects that do not define one in `.gitlab-ci.yml`:

    ```
    Please enter the Docker image (eg. ruby:2.1):
    alpine:latest
    ```

[tags]: https://docs.gitlab.com/ce/ci/runners/#using-tags
