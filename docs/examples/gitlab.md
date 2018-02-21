# How to configure GitLab Runner for GitLab CE integration tests如何配置GitLab Runner使用GitLab CE集成测试

We will register the Runner using a confined Docker executor.我们将注册该Runner使用指定的Docker执行器。

The registration token can be found at `https://gitlab.com/project_namespace/project_name/runners`.
You can export it as a variable and run the command below as is:
注册令牌可在`https://gitlab.com/project_namespace/project_name/runners`找到。你将其导出做作为一个变量并运行以下命令：

```bash
gitlab-runner register \
--non-interactive \
--url "https://gitlab.com" \
--registration-token "$REGISTRATION_TOKEN" \
--description "gitlab-ce-ruby-2.1" \
--executor "docker" \
--docker-image ruby:2.1 --docker-mysql latest \
--docker-postgres latest --docker-redis latest
```

You now have a GitLab CE integration testing instance with bundle caching.
Push some commits to test it.
现在你拥有一个GitLab CE集成测试实例，可以推送一些提交来测试它：

For [advanced configuration](../configuration/advanced-configuration.md), look into
`/etc/gitlab-runner/config.toml` and tune it.
高级配置内容可在/etc/gitlab-runner/config.toml调整。
