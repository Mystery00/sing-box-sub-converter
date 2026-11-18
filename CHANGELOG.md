### 新增

- 订阅配置支持 `providers.yaml` 和 `providers.yml`，可以用 `#` 为配置项添加注释。字段名与 `providers.json` 相同。
- 启动时按 `providers.json`、`providers.yaml`、`providers.yml` 的顺序选取文件。已有 JSON 配置可以直接使用。迁移到 YAML 时，请将旧 JSON 改名备份并重启服务。此后编辑选中的文件即可热重载。

### 修复

- 启动时配置读取或解析失败会返回错误。运行期间热重载失败会保留上一次成功加载的配置。
