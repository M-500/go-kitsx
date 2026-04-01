# eventbus

`eventbus` 是一个面向生产环境设计的统一事件总线包，目标不是只做一个“进程内回调集合”，而是提供一套统一的事件模型、订阅模型和传输抽象，让业务代码可以在本地模式、channel 模式和 Kafka 模式之间平滑切换。

当前目录：`pkg/eventbus`

## 设计目标

- 同一套 `Message` / `Handler` / `SubscribeOptions` API 适配多种传输。
- 本地可直接用，不依赖外部基础设施。
- 面向生产环境补齐基本能力：优雅关闭、可配置并发、缓冲队列、超时、重试、panic 恢复、指标统计。
- Kafka 不和某个 SDK 强绑定，避免核心包被 `sarama`、`segmentio-kafka-go` 之类客户端锁死。
- 业务 handler 只关心事件本身，不关心底层传输细节。

## 模块结构

- `Bus`：统一入口，负责注册 transport、发布消息、订阅消息、汇总指标、管理生命周期。
- `Message`：统一事件信封。包含 `ID`、`Topic`、`Headers`、`Metadata`、`Payload` 等字段。
- `Driver`：传输驱动接口。当前实现：
  - `local`
  - `channel`
  - `kafka`
- `Middleware`：围绕 handler 的横切扩展点。
- `RetryPolicy`：失败重试策略。

## 三种模式的定位

### 1. local

适合单进程内业务解耦。

特点：

- 发布后将消息分发到匹配订阅者的独立队列。
- 每个订阅者可配置并发数和缓冲区。
- 不依赖额外基础设施。

适用场景：

- 站内消息联动
- 缓存刷新
- 审计记录
- SSE/WebSocket 内部分发

### 2. channel

适合“进程内异步队列化”的模式，比 local 多一层中心输入通道。

特点：

- 发布方先写入统一 channel。
- 由 dispatcher 再 fan-out 给匹配订阅者。
- 发布和执行更解耦，适合削峰和突发流量。

适用场景：

- 高频但轻量的异步处理
- 想降低发布路径阻塞感
- 不需要跨进程，但需要更明显的队列语义

### 3. kafka

适合跨服务、跨进程事件流。

特点：

- `eventbus` 提供统一 Kafka 适配层接口：
  - `KafkaProducer`
  - `KafkaConsumer`
- 你可以用自己的 `sarama`、`segmentio-kafka-go`、云厂商 SDK 来实现这两个接口。
- 核心包不会和某一个 SDK 强耦合。

适用场景：

- 服务间事件驱动
- 削峰、解耦、可回放
- 大规模异步任务

## 为什么 Kafka 不直接内置某个客户端

生产环境里 Kafka 客户端选型通常受以下因素影响：

- 团队现有标准库
- SASL / TLS / 云厂商认证方式
- 事务、幂等、批量发送策略
- 监控埋点和 tracing

如果在核心包里直接绑定某个 SDK，后续迁移成本会很高。所以这里采用“核心统一抽象 + 外部适配”的方式，更适合长期维护。

## 关键设计思路

### 统一消息信封

所有模式都使用同一份 `Message`：

```go
type Message struct {
    ID        string
    Topic     string
    Key       string
    Source    string
    Timestamp time.Time
    Headers   map[string]string
    Metadata  map[string]string
    Payload   json.RawMessage
}
```

这样有几个好处：

- 本地模式和 Kafka 模式结构对齐。
- 更容易做审计、日志、幂等、追踪。
- 业务层不用为不同传输写不同 handler。

### 统一订阅模型

`SubscribeOptions` 中统一了以下维度：

- `Transport`
- `ConsumerGroup`
- `Concurrency`
- `Buffer`
- `Timeout`
- `RetryPolicy`
- `MatchMode`

这让“同一个业务订阅逻辑”可以切换不同传输，而不是重写一套消费代码。

### 可靠性增强

当前内建能力：

- `panic` 恢复
- 超时控制
- 重试控制
- 指标统计
- 优雅关闭

### 优雅关闭

`Bus.Close(ctx)` 会：

- 拒绝新发布和新订阅
- 取消订阅上下文
- 关闭 transport
- 等待 worker goroutine 退出

这比简单地清空 map 更接近生产环境需要的生命周期管理。

## 使用方式

### 1. 本地模式

```go
bus := eventbus.NewBus()
defer bus.Close(context.Background())

_, err := bus.Subscribe("user.created", func(ctx context.Context, msg *eventbus.Message) error {
    var payload struct {
        Name string `json:"name"`
    }
    return msg.DecodeJSON(&payload)
}, eventbus.SubscribeOptions{
    Transport:   eventbus.TransportLocal,
    Concurrency: 4,
    Buffer:      128,
})

err = bus.PublishJSON(context.Background(), "user.created", map[string]any{
    "name": "alice",
}, eventbus.PublishOptions{
    Transport: eventbus.TransportLocal,
})
```

### 2. Channel 模式

```go
bus := eventbus.NewBus()

_, err := bus.Subscribe("order.*", handleOrder, eventbus.SubscribeOptions{
    Transport: eventbus.TransportChannel,
    MatchMode: eventbus.MatchPattern,
    Buffer:    256,
})

err = bus.PublishJSON(ctx, "order.created", payload, eventbus.PublishOptions{
    Transport: eventbus.TransportChannel,
})
```

### 3. Kafka 模式

先实现两个接口：

```go
type KafkaProducer interface {
    Publish(ctx context.Context, record KafkaRecord) error
}

type KafkaConsumer interface {
    Subscribe(ctx context.Context, req KafkaConsumeRequest, handler func(context.Context, KafkaRecord) error) error
}
```

然后注册：

```go
producer := NewYourKafkaProducer(...)
consumer := NewYourKafkaConsumer(...)

bus := eventbus.NewBus()
_ = bus.RegisterTransport(eventbus.TransportKafka, eventbus.NewKafkaTransport(producer, consumer))

_, err := bus.Subscribe("audit.*", handleAudit, eventbus.SubscribeOptions{
    Transport:     eventbus.TransportKafka,
    ConsumerGroup: "audit-service",
    MatchMode:     eventbus.MatchPattern,
})

err = bus.PublishJSON(ctx, "audit.created", payload, eventbus.PublishOptions{
    Transport: eventbus.TransportKafka,
})
```

## MatchMode

- `MatchExact`：精确匹配
- `MatchPrefix`：前缀匹配
- `MatchPattern`：基于 `path.Match` 的通配符匹配，适合 `order.*` 这类主题

建议：

- 高吞吐核心路径优先使用 `MatchExact`
- 只有确实需要模糊订阅时再使用 `MatchPattern`

## 指标

当前总线暴露聚合指标：

- `Published`
- `PublishFailed`
- `Delivered`
- `HandlerFailed`
- `Retries`
- `Panics`
- `Dropped`

用法：

```go
metrics := bus.Metrics()
```

## 生产环境建议

### 建议 1：主题要稳定命名

建议统一采用点分主题，例如：

- `user.created`
- `user.updated`
- `order.created`
- `audit.login`

避免随意拼接和多套命名风格并存。

### 建议 2：事件要尽量做成不可变事实

推荐发布“事实型事件”，比如：

- 用户已创建
- 订单已支付
- 权限已变更

而不是发布模糊命令，比如“处理一下用户”。

### 建议 3：Kafka 消费 handler 做幂等

即使底层客户端支持提交 offset，业务 handler 仍要设计为可重复执行：

- 基于事件 ID 去重
- 基于业务主键幂等更新
- 对外部副作用加防重保护

### 建议 4：重试要区分临时错误和永久错误

当前框架层提供通用重试能力，但生产上最好在业务里区分：

- 临时网络抖动：可重试
- 参数非法、数据缺失：不应重试

### 建议 5：高并发场景关注缓冲和背压

核心参数：

- `Concurrency`
- `Buffer`
- ChannelTransport 的总输入缓冲
- Kafka 客户端自己的批量发送与拉取参数

这些参数需要按业务流量压测后再定。

## 当前实现边界

这一版已经具备“生产可落地”的基础骨架，但还保留了一些扩展空间：

- 还没有内置 DLQ
- 还没有内置 tracing instrumentation
- 还没有内置 schema registry
- Kafka offset 提交策略由外部 consumer 适配器控制
- 没有把指标直接对接 Prometheus

如果你后续要继续增强，建议优先级如下：

1. Prometheus metrics
2. OpenTelemetry tracing
3. DLQ / poison message 策略
4. 延迟重试队列
5. 事件 schema 校验

## 测试

当前包内已经包含：

- local 模式示例
- channel 模式示例
- local 重试与指标测试
- Kafka mock 适配测试

运行：

```bash
go test ./pkg/eventbus
```
