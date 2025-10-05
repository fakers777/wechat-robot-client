package startup

import (
	"wechat-robot-client/plugin"
	"wechat-robot-client/plugin/plugins"
	"wechat-robot-client/vars"
)

func RegisterMessagePlugin() {
	vars.MessagePlugin = plugin.NewMessagePlugin()
	// 群聊聊天插件
	vars.MessagePlugin.Register(plugins.NewChatRoomAIChatSessionStartPlugin())
	vars.MessagePlugin.Register(plugins.NewChatRoomAIChatSessionEndPlugin())
	vars.MessagePlugin.Register(plugins.NewChatRoomAIChatPlugin())
	// 群聊绘画插件
	vars.MessagePlugin.Register(plugins.NewChatRoomAIDrawingSessionStartPlugin())
	vars.MessagePlugin.Register(plugins.NewChatRoomAIDrawingSessionEndPlugin())
	vars.MessagePlugin.Register(plugins.NewChatRoomAIDrawingPlugin())
	// 朋友聊天插件
	vars.MessagePlugin.Register(plugins.NewFriendAIChatPlugin())
	// 朋友绘画插件
	vars.MessagePlugin.Register(plugins.NewFriendAIDrawingSessionStartPlugin())
	vars.MessagePlugin.Register(plugins.NewFriendAIDrawingSessionEndPlugin())
	vars.MessagePlugin.Register(plugins.NewFriendAIDrawingPlugin())
	// 群聊拍一拍交互插件
	vars.MessagePlugin.Register(plugins.NewPatPlugin())
	// 图片自动上传插件
	vars.MessagePlugin.Register(plugins.NewImageAutoUploadPlugin())

	// === 新增功能插件 ===
	// Apilot多功能插件（星座运势、热榜、天气等）
	vars.MessagePlugin.Register(plugins.NewApilotPlugin())
	// 京东积存金价格查询插件
	vars.MessagePlugin.Register(plugins.NewJdjcjPlugin())
	// KFC相关插件
	vars.MessagePlugin.Register(plugins.NewKFCPlugin())
	vars.MessagePlugin.Register(plugins.NewKFCStoryPlugin())
	vars.MessagePlugin.Register(plugins.NewKFCMenuPlugin())
	vars.MessagePlugin.Register(plugins.NewKFCWenanPlugin())
	// 土味情话插件
	vars.MessagePlugin.Register(plugins.NewLovePlugin())
	vars.MessagePlugin.Register(plugins.NewLoveStoryPlugin())
	vars.MessagePlugin.Register(plugins.NewLoveAdvicePlugin())
	vars.MessagePlugin.Register(plugins.NewLoveTestPlugin())
	// 网易云音乐插件
	vars.MessagePlugin.Register(plugins.NewNeteasyPlugin())
	// 图像识别插件（iPad版本）
	vars.MessagePlugin.Register(plugins.NewImageRecognitionIPadPlugin())
}
