# Test Cases

- 暂无

# Description
>     modify an existing validator account.

>     修改已存在的验证人账户。

# Example

> 下面实例中假设:
> - `Arya` 地址为: `address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`
> - `Sansa` 地址为: `address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

`Arya`可通过`modify-validator`添加/修改节点信息：
```bash
$ qoscli tx modify-validator --moniker "Arya's node" --owner Arya --tokens 1000 --logo "https://..." --website "https://..." --description "Long live Arya."
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"265"}
```

执行成为验证节点命令后将从`Arya`账户扣除1000QOS，绑定到验证节点中，验证节点参与投票、打块所获得的挖矿收益将直接增加到`Arya`账户。

# Usage
```
  qoscli tx modify-validator [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                | Required | Input Type | Default Input             | Input Range | Description                                            |
|:---------|:--------------------|:---------|:-----------|:--------------------------|:------------|:-------------------------------------------------------|
| `-h`     | `--help`            | ✖        | -          | -                         | -           | 帮助文档                                                   |
| -        | `--async`           | ✖        | -          | -                         | -           | 是否异步广播交易                                               |
| -        | `--chain-id`        | ✖        | string     | -                         | -           | Tendermint节点的链ID                                       |
| -        | `--indent`          | ✖        | -          | -                         | -           | 向JSON响应添加缩进                                            |
| -        | `--max-gas`         | ✖        | int        | `100000`                  | -           | 每个Tx设置的气体限制值                                           |
| `-n`     | `--node`            | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>`                |
| -        | `--nonce`           | ✖        | int        | -                         | -           | 要签署Tx的帐户nonce                                          |
| -        | `--nonce-node`      | ✖        | string     | -                         | -           | 用于其他链查询账户nonce的Tendermint RPC接口: `tcp://<host>:<port>` |
| -        | `--qcp`             | ✖        | -          | -                         | -           | 是否启用QCP模式(qcp mode), 发送QCP Tx                          |
| -        | `--qcp-blockheight` | ✖        | int        | -                         | -           | QCP模式Flag标志: 原始Tx块高度，块高度必须大于0                          |
| -        | `--qcp-extends`     | ✖        | string     | -                         | -           | QCP模式Flag标志: QCP Tx扩展信息                                |
| -        | `--qcp-from`        | ✖        | string     | -                         | -           | QCP模式Flag标志: QCP Tx源链ID                                |
| -        | `--qcp-seq`         | ✖        | int        | -                         | -           | QCP模式Flag标志: QCP顺序                                     |
| -        | `--qcp-signer`      | ✖        | string     | -                         | -           | QCP模式Flag标志: QCP Tx签名者key名称                            |
| -        | `--qcp-txindex`     | ✖        | int        | -                         | -           | QCP模式Flag标志: 原始Tx索引                                    |
| -        | `--trust-node`      | ✖        | -          | -                         | -           | 是否信任连接的完整节点（不验证其响应证据）                                  |
| -        | `--details`         | ✖        | string     | -                         | -           | (主要参数)验证人详细描述信息, `len(details) <= 1000`                |
| -        | `--logo`            | ✖        | string     | -                         | -           | (主要参数)logo链接， `len(logo) <= 255`                       |
| -        | `--moniker`         | ✖        | string     | -                         | -           | (主要参数)验证节点名字，`len(moniker) <= 300`                     |
| -        | `--owner`           | ✖        | string     | -                         | -           | (主要参数)`Owner`账户本地密钥库名字或账户地址                            |
| -        | `--website`         | ✖        | string     | -                         | -           | (主要参数)验证人网址， `len(website) <= 255`                     |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |