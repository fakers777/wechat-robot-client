package plugins

import (
	"encoding/base64"
	"fmt"
	"strings"
	"wechat-robot-client/interface/plugin"
	"wechat-robot-client/model"
	"wechat-robot-client/service"

	"github.com/sashabaranov/go-openai"
)

// ImageRecognitionIPadPlugin iPad图片识别插件
type ImageRecognitionIPadPlugin struct{}

// NewImageRecognitionIPadPlugin 创建iPad图片识别插件实例
func NewImageRecognitionIPadPlugin() plugin.MessageHandler {
	return &ImageRecognitionIPadPlugin{}
}

// GetName 获取插件名称
func (p *ImageRecognitionIPadPlugin) GetName() string {
	return "ImageRecognitionIPad"
}

// GetLabels 获取插件标签
func (p *ImageRecognitionIPadPlugin) GetLabels() []string {
	return []string{"image", "recognition", "ai", "ipad"}
}

// PreAction 前置处理
func (p *ImageRecognitionIPadPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *ImageRecognitionIPadPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *ImageRecognitionIPadPlugin) Run(ctx *plugin.MessageContext) bool {
	// 检查消息内容是否包含"识别"关键词
	if !strings.Contains(ctx.MessageContent, "识别") {
		return false
	}

	// 检查是否有引用的消息
	if ctx.ReferMessage == nil {
		p.sendReply(ctx, "请先发送一张图片，然后引用该图片消息并输入包含'识别'的文字。")
		return true
	}

	// 检查引用的消息是否是图片
	if !p.isImageMessage(ctx.ReferMessage) {
		p.sendReply(ctx, "请引用一张图片消息进行识别。")
		return true
	}

	// 获取图片数据
	imageDataURL, err := p.getImageData(ctx)
	if err != nil {
		p.sendReply(ctx, fmt.Sprintf("获取图片失败: %v", err))
		return true
	}

	// 使用AI模型识别图片
	result, err := p.recognizeImage(ctx, imageDataURL)
	if err != nil {
		p.sendReply(ctx, fmt.Sprintf("图片识别失败: %v", err))
		return true
	}

	// 发送识别结果
	p.sendReply(ctx, result)
	return true
}

// isImageMessage 检查消息是否是图片消息
func (p *ImageRecognitionIPadPlugin) isImageMessage(message *model.Message) bool {
	// 检查消息类型是否为图片
	if message.Type == model.MsgTypeImage {
		return true
	}
	
	// 或者检查是否有图片相关的字段
	if message.AttachmentUrl != "" {
		return true
	}
	
	return false
}

// getImageData 获取图片数据
func (p *ImageRecognitionIPadPlugin) getImageData(ctx *plugin.MessageContext) (string, error) {
	// 如果引用消息已经有附件URL，直接使用
	if ctx.ReferMessage.AttachmentUrl != "" {
		return ctx.ReferMessage.AttachmentUrl, nil
	}

	// 下载引用的图片
	attachDownloadService := service.NewAttachDownloadService(ctx.Context)
	imageBytes, contentType, _, err := attachDownloadService.DownloadImage(ctx.ReferMessage.ID)
	if err != nil {
		return "", err
	}

	// 转换为base64格式
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Image)

	return dataURL, nil
}

// recognizeImage 使用AI模型识别图片
func (p *ImageRecognitionIPadPlugin) recognizeImage(ctx *plugin.MessageContext, imageDataURL string) (string, error) {
	// 构建AI消息
	aiContext := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleUser,
			MultiContent: []openai.ChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL: imageDataURL,
					},
				},
				{
					Type: openai.ChatMessagePartTypeText,
					Text: "请详细描述这张图片的内容，包括主要物体、场景、颜色、文字等所有可见的元素。",
				},
			},
		},
	}

	// 使用AI聊天服务进行识别
	aiChatService := service.NewAIChatService(ctx.Context, ctx.Settings)
	aiReply, err := aiChatService.Chat(aiContext)
	if err != nil {
		return "", err
	}

	// 提取AI回复的文本内容
	var aiReplyText string
	if aiReply.Content != "" {
		aiReplyText = aiReply.Content
	} else if len(aiReply.MultiContent) > 0 {
		aiReplyText = aiReply.MultiContent[0].Text
	}

	if aiReplyText == "" {
		aiReplyText = "AI识别结果为空，请检查图片或重试。"
	}

	return aiReplyText, nil
}

// sendReply 发送回复
func (p *ImageRecognitionIPadPlugin) sendReply(ctx *plugin.MessageContext, content string) {
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content)
	}
}