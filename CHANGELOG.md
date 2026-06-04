### 修复

- **AnyTLS**：补齐 `tls.server_name` 默认值处理，支持 `alpn` 字符串/数组兼容与 `detour` 输出
- **AnyTLS**：规范化 `client-fingerprint`，仅输出 sing-box 支持的 uTLS 指纹值
- **AnyTLS**：移除不安全的类型断言，避免配置字段类型不符时解析直接 panic
