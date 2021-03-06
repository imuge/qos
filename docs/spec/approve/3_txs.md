# 交易

这里具体介绍预授权模块包含的所有交易类型，这些交易操作会直接影响[存储](2_state.md)。

## 创建预授权

[发送创建预授权交易](../../command/qoscli.md#创建预授权)创建预授权。

### 结构

```go
type TxCreateApprove struct {
	types.Approve
}
```
其中[Approve](2_state.md#预授权)为创建预授权，[增加预授权](#增加预授权)，[减少预授权](#减少预授权)，[使用预授权](#使用预授权)通用结构。

### 校验

创建预授权需要通过以下校验：
- `from`,`to`地址均不能为空，且不能相同
- `coins`中币种必须都存在，币值均为正，且不能重复
- `from`账户在QOS网络中必须存在
- 不存在已授权记录

通过校验并成功执行交易后，可通过[查询预授权](../../command/qoscli.md#查询预授权)搜索预授权信息。

### 签名

`from`账户

### 交易费

`0` 鼓励用户使用预授权功能

## 增加预授权

[发送增加预授权交易](../../command/qoscli.md#增加预授权)增加预授权。

### 结构

```go
type TxIncreaseApprove struct {
	types.Approve
}
```
其中[Approve](2_state.md#预授权)为[创建预授权](#创建预授权)，增加预授权，[减少预授权](#减少预授权)，[使用预授权](#使用预授权)通用结构。

### 校验

创建预授权需要通过以下校验：
- 存在授权记录

### 签名

`from`账户

### 交易费

`0` 鼓励用户使用预授权功能

## 减少预授权

[发送减少预授权交易](../../command/qoscli.md#减少预授权)减少预授权。

### 结构

```go
type TxDecreaseApprove struct {
	types.Approve
}
```
其中[Approve](2_state.md#预授权)为[创建预授权](#创建预授权)，[增加预授权](#增加预授权)，减少预授权，[使用预授权](#使用预授权)通用结构。

### 校验

创建预授权需要通过以下校验：
- 存在已授权记录
- `coins`中币种币值不能大于预授权记录值

### 签名

`from`账户

### 交易费

`0` 鼓励用户使用预授权功能

## 使用预授权

被授权账户提取授权币种币值到账户，相应扣除授权账户余额，减少更新预授权额度。

[发送使用预授权交易](../../command/qoscli.md#使用预授权)，使用预授权。

### 结构

```go
type TxUseApprove struct {
	types.Approve
}
```
其中[Approve](2_state.md#预授权)为[创建预授权](#创建预授权)，[增加预授权](#增加预授权)，[减少预授权](#减少预授权)，使用预授权通用结构。

### 校验

创建预授权需要通过以下校验：
- 存在已授权记录
- `coins`中币种币值不能大于预授权记录值

### 签名

`to`账户

### 交易费

`0` 鼓励用户使用预授权功能

> 当执行完减少预授权，预授权所有币种币值为零，预授权将被删除

## 取消预授权

取消已存在的预授权，将从数据库中删除授权信息。

[发送取消预授权交易](../../command/qoscli.md#取消预授权)，取消预授权

### 结构

```go
type TxCancelApprove struct {
    From btypes.AccAddress `json:"from"` // 授权账号
    To   btypes.AccAddress `json:"to"`   // 被授权账号
}
```

### 校验

创建预授权需要通过以下校验：
- 存在已授权记录

### 签名

`from`账户

### 交易费

`0` 鼓励用户使用预授权功能