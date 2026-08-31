# 自动化构建与发布

## 持续集成

`.github/workflows/ci.yml` 会在每次 push 和 Pull Request 时执行：

1. 下载 Go 模块依赖；
2. 运行 `go test ./...`；
3. 生成 Revel 服务端入口并编译；
4. 构建 Docker 镜像。

数据库集成测试仍需显式提供临时 MongoDB/PostgreSQL，普通 CI 中未配置数据库的集成测试会按测试代码约定自动跳过。

## 正式发布

正式版本标签必须采用 `vMAJOR.MINOR.PATCH` 格式，并与 `app/version/version.go` 中的 `Current` 完全一致。例如当前版本的发布标签为：

```bash
git tag -a v1.0.0 -m "Pearlnote 1.0.0"
git push origin v1.0.0
```

标签推送后，`.github/workflows/release.yml` 会：

1. 校验标签格式和代码版本；
2. 再次执行全部 Go 测试；
3. 构建 Linux amd64、Linux arm64、Windows amd64、macOS amd64 和 macOS arm64 发布包；
4. 生成 `checksums.txt`；
5. 根据提交记录生成 Release Notes 并创建 GitHub Release。

每个发布包都包含服务端、MongoDB/PostgreSQL 迁移工具、配置、页面资源、数据库 Schema、初始 MongoDB 数据和文档。

## 发布新版本

发布前应同时完成以下事项：

1. 修改 `app/version/version.go` 中的版本；
2. 在 `app/db/migrations.go` 添加该版本需要的幂等数据库迁移；
3. 更新文档并执行 `go test ./...`；
4. 合并到正式发布分支；
5. 创建并推送对应标签。

如果标签与应用版本不一致，Release 工作流会直接失败，不会产生错误版本的附件。

## 权限

CI 工作流只有仓库内容只读权限。Release 工作流只在正式标签触发，并使用 GitHub 自动提供的 `GITHUB_TOKEN` 创建 Release，无需配置个人访问令牌。
