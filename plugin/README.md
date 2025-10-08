# 微信机器人插件系统说明

## 概述

微信机器人插件系统是一个基于消息处理的插件架构，允许开发者创建各种功能插件来扩展机器人的能力。每个插件都实现了 `MessageHandler` 接口，可以处理不同类型的消息和触发条件。

## 插件架构

### 核心接口

```go
type MessageHandler interface {
    GetName() string                    // 插件名称
    GetLabels() []string                // 插件标签
    PreAction(ctx *MessageContext) bool // 前置处理
    PostAction(ctx *MessageContext)     // 后置处理
    Run(ctx *MessageContext) bool       // 主要逻辑
}
```

### 消息上下文

```go
type MessageContext struct {
    Context        context.Context      // 上下文
    Settings       settings.Settings    // 设置配置
    Message        *model.Message       // 消息对象
    MessageContent string              // 消息内容
    Pat            bool                // 是否拍一拍
    ReferMessage   *model.Message      // 引用消息
    MessageService MessageServiceIface // 消息服务接口
}
```

## 现有插件功能

### 1. AI 聊天插件 (`ai_chat.go`)
- **功能**: 基于AI触发词的智能聊天
- **标签**: `["internal", "chat"]`
- **触发条件**: 包含AI触发词的消息
- **支持**: 群聊和私聊，自动去除触发词

### 2. 好友AI聊天插件 (`friend_chat.go`)
- **功能**: 私聊AI聊天，不响应自己发送的消息
- **标签**: `["chat"]`
- **触发条件**: 私聊消息且AI聊天功能开启
- **特点**: 自动更新AI会话上下文

### 3. AI绘图插件 (`ai_drawing.go`)
- **功能**: 基于文本描述生成AI图片
- **标签**: `["internal", "drawing"]`
- **支持模型**: 豆包、即梦、智谱、混元等
- **特点**: 支持多种AI绘图服务

### 4. 拍一拍插件 (`pat.go`)
- **功能**: 响应拍一拍动作
- **标签**: `["pat"]`
- **支持类型**: 文本回复、语音回复
- **特点**: 可配置回复内容和语音音色

### 5. 抖音视频解析插件 (`douyin_video_parse.go`)
- **功能**: 解析抖音分享链接，获取视频信息
- **标签**: `["douyin"]`
- **特点**: 自动识别抖音链接，返回视频详情和下载链接

### 6. 自动加群插件 (`auto_join_group.go`)
- **功能**: 自动邀请用户加入指定群聊
- **标签**: `["auto"]`
- **触发条件**: "申请进群 群名" 格式的消息
- **特点**: 简化群聊邀请流程

### 7. 群聊AI聊天插件 (`chatroom_chat.go`)
- **功能**: 群聊中的AI聊天功能
- **标签**: `["chatroom", "chat"]`
- **特点**: 专门处理群聊场景的AI交互

### 8. AI语音合成插件 (`ai_tts.go`)
- **功能**: 将文本转换为语音消息
- **标签**: `["internal", "tts"]`
- **支持**: 多种语音合成服务

### 9. AI图片识别插件 (`ai_image_recognizer.go`)
- **功能**: 识别图片内容并生成描述
- **标签**: `["internal", "image"]`
- **特点**: 支持多种图片识别模型

### 10. AI图片编辑插件 (`ai_image_edit.go`)
- **功能**: 基于文本指令编辑图片
- **标签**: `["internal", "image"]`
- **特点**: 支持图片风格转换、内容修改等

## 插件使用方式

### 1. 注册插件

```go
pluginManager := plugin.NewMessagePlugin()
pluginManager.Register(plugins.NewAIChatPlugin())
pluginManager.Register(plugins.NewPatPlugin())
```

### 2. 插件执行流程

1. **PreAction**: 前置检查，返回false可跳过插件
2. **Run**: 主要逻辑处理，返回true表示消息已处理
3. **PostAction**: 后置处理，用于清理资源等

### 3. 消息服务接口

```go
type MessageServiceIface interface {
    SendTextMessage(toWxID, content string, at ...string) error
    MsgUploadImg(toWxID string, image io.Reader) (*model.Message, error)
    MsgSendVoice(toWxID string, voice io.Reader, voiceExt string) error
    MsgSendVideo(toWxID string, video io.Reader, videoExt string) error
    SendMusicMessage(toWxID string, songTitle string) error
    ShareLink(toWxID string, shareLinkInfo robot.ShareLinkMessage) error
    // ... 更多方法
}
```

## 开发新插件

### 1. 实现接口

```go
type MyPlugin struct{}

func (p *MyPlugin) GetName() string {
    return "MyPlugin"
}

func (p *MyPlugin) GetLabels() []string {
    return []string{"custom"}
}

func (p *MyPlugin) PreAction(ctx *plugin.MessageContext) bool {
    // 前置检查逻辑
    return true
}

func (p *MyPlugin) PostAction(ctx *plugin.MessageContext) {
    // 后置处理逻辑
}

func (p *MyPlugin) Run(ctx *plugin.MessageContext) bool {
    // 主要业务逻辑
    return true
}
```

### 2. 注册插件

```go
pluginManager.Register(NewMyPlugin())
```

## 插件开发最佳实践

### 1. 错误处理
- 使用适当的错误处理机制
- 向用户发送友好的错误消息
- 记录详细的错误日志

### 2. 性能优化
- 避免阻塞操作
- 合理使用缓存
- 优化网络请求

### 3. 用户体验
- 提供清晰的使用说明
- 支持多种触发方式
- 响应及时且准确

### 4. 配置管理
- 支持动态配置
- 提供默认设置
- 允许用户自定义

## 插件标签说明

- `internal`: 内部插件，通常由系统自动触发
- `chat`: 聊天相关插件
- `drawing`: 绘图相关插件
- `pat`: 拍一拍相关插件
- `douyin`: 抖音相关插件
- `auto`: 自动化功能插件
- `chatroom`: 群聊专用插件
- `tts`: 语音合成插件
- `image`: 图片处理插件

## 扩展功能

### 1. 自定义触发条件
可以通过修改 `PreAction` 方法来实现自定义的触发条件，比如：
- 关键词匹配
- 正则表达式
- 消息类型检查
- 用户权限验证

### 2. 多步骤处理
可以在 `Run` 方法中实现复杂的多步骤处理逻辑，比如：
- 数据收集
- 外部API调用
- 数据处理
- 结果返回

### 3. 状态管理
可以通过上下文或外部存储来管理插件状态，实现：
- 会话记忆
- 用户偏好
- 历史记录
- 缓存数据

## 注意事项

1. **消息处理**: 确保插件正确处理不同类型的消息
2. **资源管理**: 及时释放资源，避免内存泄漏
3. **并发安全**: 考虑多用户并发访问的情况
4. **配置验证**: 验证必要的配置参数
5. **错误恢复**: 实现适当的错误恢复机制

## 示例插件

参考 `plugins/` 目录下的现有插件实现，了解不同场景下的插件开发模式。