# GitLab Runner

This is the repository of the official GitLab Runner written in Go.
这是官方GitLab Runner的存储库，用Go语言编写。
It runs tests and sends the results to GitLab.、
Runner运行测试并把测试结果返回给GitLab。
[GitLab CI](https://about.gitlab.com/gitlab-ci) is the open-source
continuous integration service included with GitLab that coordinates the testing.
GitLab CI是开源的持续集成服务。
The old name of this project was GitLab CI Multi Runner but please use "GitLab Runner" (without CI) from now on.
Runner项目的旧名为GitLab CI Multi Runner，但从现在开始请使用称呼此项目"GitLab Runner"。

![Build Status](https://gitlab.com/gitlab-org/gitlab-runner/badges/master/build.svg)

## Runner and GitLab CE/EE compatibility--Runner和GitLab的兼容性

For a list of compatible versions between GitLab and GitLab Runner, consult
the [compatibility chart](https://docs.gitlab.com/runner/#compatibility-chart).

## Release process

The description of release process of GitLab Runner project can be found in the [release documentation](docs/release_process/README.md).

## Contributing--向此项目贡献

Contributions are welcome, see [`CONTRIBUTING.md`](CONTRIBUTING.md) for more details.

### Closing issues and merge requests--issue和合并请求的关闭机制

GitLab is growing very fast and we have a limited resources to deal with reported issues
and merge requests opened by the community volunteers. We appreciate all the contributions
coming from our community. But to help all of us with issues and merge requests management
we need to create some closing policy.

If an issue or merge request has a ~"waiting for feedback" label and the response from the
reporter has not been received for 14 days, we can close it using the following response
template:

```
We haven't received an update for more than 14 days so we will assume that the
problem is fixed or is no longer valid. If you still experience the same problem
try upgrading to the latest version. If the issue persists, reopen this issue
or merge request with the relevant information.
```

## Documentation--文档

The documentation source files can be found under the [docs/](docs/) directory. You can
read the documentation online at https://docs.gitlab.com/runner/.

## Requirements--安装Runner的要求

[Read about the requirements of GitLab Runner.](https://docs.gitlab.com/runner/#requirements)

## Features--Runner的功能

[Read about the features of GitLab Runner.](https://docs.gitlab.com/runner/#features)

## Executors compatibility chart

[Read about what options each executor can offer.](https://docs.gitlab.com/runner/executors/#compatibility-chart)

## Install GitLab Runner--安装GitLab Runner

Visit the [installation documentation](https://docs.gitlab.com/runner/install/).

## Use GitLab Runner--使用Runner

See [https://docs.gitlab.com/runner/#using-gitlab-runner](https://docs.gitlab.com/runner/#using-gitlab-runner).

## Select executor--选择执行器

See [https://docs.gitlab.com/runner/executors/#selecting-the-executor](https://docs.gitlab.com/runner/executors/#selecting-the-executor).

## Troubleshooting--故障排除

Read the [FAQ](https://docs.gitlab.com/runner/faq/).

## Advanced Configuration--高级配置

See [https://docs.gitlab.com/runner/#advanced-configuration](https://docs.gitlab.com/runner/#advanced-configuration).

## Changelog--变更日志

Visit the [Changelog](CHANGELOG.md) to view recent changes.

## The future--未来

* Please see the [GitLab Direction page](https://about.gitlab.com/direction/).
* Feel free submit issues with feature proposals on the issue tracker.

## Author--作者

```
2014 - 2015   : [Kamil Trzciński](mailto:ayufan@ayufan.eu)
2015 - now    : GitLab Inc. team and contributors
```


## License--证书

This code is distributed under the MIT license, see the [LICENSE](LICENSE) file.
