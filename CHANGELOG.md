### 新增

- **订阅源出站模板（outbound-template）**：在 config 模板的 `outbounds` 数组中可声明一种特殊条目 `{"type": "outbound-template", "tag": "<订阅 tag>", ...}`，所有 `FromSub` 等于该 tag 的出站节点都会以该条目为基底套用其余字段（如 `domain_resolver`、`multiplex`、`tls.utls`、`tcp_fast_open` 等任意 sing-box 出站合法字段），实现按订阅源批量注入公共属性。`outbound-template` 条目处理完后会从最终 `outbounds` 中剔除，不污染输出。
- **`utils.DeepMerge` 深度合并工具**：以 base 为基底、override 覆盖返回新 map；同名嵌套对象递归合并，标量与数组整体覆盖；不修改入参。供出站模板特性及后续合并场景复用。

### 合并语义说明

- 节点解析出的字段优先级高于模板：同名标量以节点为准，节点端口、地址等真实参数不会被模板误覆盖。
- 同名嵌套对象（如 `tls`）执行深度合并：模板独有子字段保留、节点独有子字段保留、同名标量子字段以节点为准。
- 数组字段整体覆盖，避免数组元素级合并的语义歧义。
- 模板 `tag` 在订阅中没有对应节点 → 安静忽略；订阅没有匹配模板 → 节点保持原样。
