## 0.0.6

### 新增

- **订阅用户信息回传**：当仅请求单个订阅时，通过 HTTP 响应头 `Subscription-Userinfo` 返回 `upload`、`download`、`total`、`expire` 字段，便于客户端展示流量用量和到期时间（适用于 `/api/generate`、`/api/quickstart`、Vercel 接口）。
- **可视化订阅信息节点**：启用 `show_sub_in_nodes: true` 且订阅端返回流量信息时，自动在节点列表末尾插入一条汇总节点，显示"剩余流量 / 剩余天数"；永不过期的订阅显示"永久"。

### 问题修复

- **Vercel 接口编译错误**：`api/vercel.go` 调用 `ProcessSubscribes` 时未正确接收返回的订阅头信息，导致无法编译。
- **Hysteria2 解析 panic**：`Obfs` 字段未初始化，向 nil map 写入时触发 panic；同时修正 `fingerprint` 字段被误写入 `Obfs` 而非 `Tls.utls` 的问题，现可正确生成 uTLS 指纹配置。
- **VLESS 节点默认名称错误**：无名称的 VLESS 节点默认 tag 后缀被错误设为 `_hysteria2`，已修正为 `_vless`。
- **SS 纯文本格式 plugin 解析 panic**：当 plugin 参数不含分隔符 `;` 时，`strings.Index` 返回 `-1` 导致切片越界；现在无分隔符时将整个字符串作为 plugin 名称。
- **多 filter 只保留最后一条结果**：模板 outbound 配置多个 `filter` 时，每条 filter 均对原始节点列表重新过滤，导致前面的结果被覆盖；已改为链式应用，各 filter 依次作用于上一步的输出。
- **表情旗帜匹配逻辑无效**：`renameNodeTagWithEmoji` 中两处特殊关键词判断逻辑完全无效（一处缺少 return/continue、一处对 emoji 字符做英文字符串检查永远为 false），已删除；同时将原来仅针对 `🇺🇲` 的旗帜替换特殊处理，泛化为对任意已有旗帜前缀的通用剥离逻辑。
- **`calculateRemainDays` 时间戳为 0 时显示异常**：`expire` 为 `0` 时应视为永不过期，之前会计算出负数天数；已修正为返回"永久"。
- **`GenName` 传入负数时 panic**：`make([]byte, n)` 在 `n < 0` 时触发 runtime panic；现在 `n <= 0` 时自动回退为默认长度 8。

### 性能优化

- **正则表达式预编译**：旗帜 Emoji 匹配正则由每次调用时动态编译改为包级别变量，避免重复编译开销。

### 依赖更新

- Go 工具链升级至 `1.26.1`
- `gin` 升级至 `v1.12.0`
- `bytedance/sonic` 升级至 `v1.15.0`
- `quic-go` 升级至 `v0.59.0`
- 其余间接依赖同步升级

### 其他

- 移除从未被调用的死代码文件 `converter/singbox.go`
- 移除 GitHub Actions 中无用的架构后缀标签清理步骤
- Vercel 接口响应 `Content-Type` 补充 `charset=utf-8`
