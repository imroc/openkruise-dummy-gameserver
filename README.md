# openkruise-dummy-gameserver

使用 OpenKruiseGame 部署的房间类型游戏服 demo。

## 游戏服特点

暴露两个接口：
* `:8099/metrics`: 暴露 Prometheus 指标，包含当前 Pod 中游戏房间数量、已分配房间的数量、当前是否空闲的状态（结合 KEDA 实现基于房间占用率的弹性伸缩）。
* `:8080/api/idle`: 探测当前 Pod 是否空闲，如果全部房间空闲，则返回空闲状态，否则返回繁忙状态（结合 OpenKruiseGame 实现缩容时优先删除空闲 Pod，避免占用中的游戏房间被中断）。
