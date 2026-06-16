# 订阅源出站模板（outbound-template）设计

日期：2026-06-16

## 背景与目标

当前所有出站节点都是解析原始订阅后直接生成的（`MergeToConfig` 末尾遍历用到的节点，逐个调用 `Convert2SingBoxOutbounds` 生成 sing-box 出站对象）。有时需要给某个订阅源生成的所有出站节点统一套用一组额外属性（拨号、多路复用、TLS、`domain_resolver`、`servername`、端口等任意 sing-box 出站合法字段）。

目标：允许在 config 模板中，针对 `providers.json` 的 `subscribes` 订阅数据源，定义一个出站模板。生成出站节点时先以该模板为基底，再用节点自身解析出的信息覆盖写入。

## 模板格式与放置位置

在 config 模板的 `outbounds` 数组中加入特殊条目，`type` 固定为字符串 `outbound-template`，`tag` 对应 `providers.json` 中订阅的 `tag` 值，其余字段即模板内容：

```json
{
  "type": "outbound-template",
  "tag": "AmyTelecom",
  "domain_resolver": "dns_amy",
  "multiplex": { "enabled": true, "protocol": "h2mux" },
  "tcp_fast_open": true
}
```

效果：所有 `FromSub == "AmyTelecom"` 的出站节点都会被套上这些字段。这些 `outbound-template` 条目本身不会出现在最终配置里——处理完即从 `outbounds` 中删除，与现有 `filter` 字段的处理方式一致。

## 合并语义（深度合并）

模板作为基底，节点自身解析出的字段覆盖在其上。规则：

- 节点字段优先级高于模板。同名标量字段 → 节点值覆盖模板值。
- 同名嵌套对象（如 `tls`）→ 递归深度合并：模板独有的子字段保留，节点独有的子字段保留，同名标量子字段以节点为准。
- 数组字段 → 整体覆盖（不做数组元素级合并，避免歧义）。

示例：

- 模板：`"tls": { "utls": { "enabled": true, "fingerprint": "chrome" } }`
- 节点生成：`"tls": { "enabled": true, "server_name": "a.com" }`
- 合并结果：

```json
"tls": { "enabled": true, "server_name": "a.com", "utls": { "enabled": true, "fingerprint": "chrome" } }
```

## 代码改动点

### `template/template.go`

1. 在 `MergeToConfig` 开头，先扫描 `outbounds`，提取所有 `type == "outbound-template"` 的条目，按 `tag` 建立映射 `map[string]map[string]any`（剥掉 `type` 和 `tag` 两个键），并将这些条目从 `outbounds` 中剔除，不进入后续占位符替换逻辑。
2. 在末尾生成出站节点的循环（当前 template.go:136-139）中，拿到 `Convert2SingBoxOutbounds(node)` 的结果 `o` 后，若 `templates[node.FromSub]` 存在，则以模板为基底、`o` 为覆盖做深度合并，再追加到 `outbounds`。

### `utils/merge.go`（新增）

新增深度合并工具函数：

```go
// DeepMerge 以 base 为基底，用 override 覆盖，返回合并后的新 map。
// override 字段优先；同名嵌套 map 递归合并；数组及标量整体覆盖。
func DeepMerge(base, override map[string]any) map[string]any
```

不修改入参，返回新 map。

### 无需改动

- 协议解析器（`converter/clash/*`、`converter/content/*`）——关联完全靠现有的 `ProxyNode.FromSub`。
- `providers.json` 与 `config` 包——模板放在 config 模板中，不涉及订阅配置结构。
- `converter/types/types.go`——无新增节点类型。

## 边界与错误处理

- 某订阅 tag 在模板里没有对应 `outbound-template` 条目 → 该订阅的节点保持原样，不受影响。
- 模板条目的 `tag` 在 `providers.json` 里找不到对应订阅 → 不影响任何节点，安静忽略（仅在生成节点时按 `FromSub` 匹配）。
- 模板字段与节点字段冲突（如模板写了 `server_port`，节点也有）→ 按合并规则节点值优先，不打日志，保持简洁。

## 测试

### `template/` 单元测试

构造含 `outbound-template` 的 config + 一组带 `FromSub` 的节点，断言：

1. 模板字段被正确套用到匹配订阅的出站节点上。
2. 嵌套 `tls` 块深度合并（模板子字段与节点子字段共存）。
3. 节点字段覆盖模板同名字段。
4. `outbound-template` 条目不出现在最终 `outbounds` 里。
5. 无匹配模板的订阅节点保持原样。

### `utils/` 表驱动测试

针对 `DeepMerge`：标量覆盖、嵌套 map 合并、数组整体覆盖、空 base、空 override 等用例。

## 数据流（改动后）

```
MergeToConfig(config, nodes)
  → 扫描 outbounds，提取 outbound-template 条目 → templates map，并从 outbounds 剔除
  → 占位符替换（{all} / {tag名} / filter）保持不变
  → 末尾生成出站节点循环：
      o = Convert2SingBoxOutbounds(node)
      if t, ok := templates[node.FromSub]; ok { o = DeepMerge(t, o) }
      outbounds = append(outbounds, o)
  → 返回最终配置
```
