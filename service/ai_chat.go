package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"wechat-robot-client/interface/ai"
	"wechat-robot-client/interface/settings"
	"wechat-robot-client/model"
	"wechat-robot-client/vars"

	"github.com/sashabaranov/go-openai"
)

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isGeminiAPI 检查是否是Gemini API
func isGeminiAPI(baseURL string) bool {
	return strings.Contains(strings.ToLower(baseURL), "gemini") ||
		strings.Contains(strings.ToLower(baseURL), "google")
}

type AIChatService struct {
	ctx    context.Context
	config settings.Settings
}

var _ ai.AIService = (*AIChatService)(nil)

func NewAIChatService(ctx context.Context, config settings.Settings) *AIChatService {
	return &AIChatService{
		ctx:    ctx,
		config: config,
	}
}

func (s *AIChatService) SetAISession(message *model.Message) error {
	return vars.RedisClient.Set(s.ctx, s.GetSessionID(message), true, defaultTTL).Err()
}

func (s *AIChatService) RenewAISession(message *model.Message) error {
	return vars.RedisClient.Expire(s.ctx, s.GetSessionID(message), defaultTTL).Err()
}

func (s *AIChatService) ExpireAISession(message *model.Message) error {
	return vars.RedisClient.Del(s.ctx, s.GetSessionID(message)).Err()
}

func (s *AIChatService) ExpireAllAISessionByChatRoomID(chatRoomID string) error {
	sessionID := fmt.Sprintf("ai_chat_session_%s:", chatRoomID)
	keys, err := vars.RedisClient.Keys(s.ctx, sessionID+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return vars.RedisClient.Del(s.ctx, keys...).Err()
}

func (s *AIChatService) IsInAISession(message *model.Message) (bool, error) {
	cnt, err := vars.RedisClient.Exists(s.ctx, s.GetSessionID(message)).Result()
	return cnt == 1, err
}

func (s *AIChatService) GetSessionID(message *model.Message) string {
	return fmt.Sprintf("ai_chat_session_%s:%s", message.FromWxID, message.SenderWxID)
}

func (s *AIChatService) IsAISessionStart(message *model.Message) bool {
	if message.Content == "#进入AI会话" {
		err := s.SetAISession(message)
		return err == nil
	}
	return false
}

func (s *AIChatService) GetAISessionStartTips() string {
	return "AI会话已开始，请输入您的问题。10分钟不说话会话将自动结束，您也可以输入 #退出AI会话 来结束会话。"
}

func (s *AIChatService) IsAISessionEnd(message *model.Message) bool {
	if message.Content == "#退出AI会话" {
		err := s.ExpireAISession(message)
		return err == nil
	}
	return false
}

func (s *AIChatService) GetAISessionEndTips() string {
	return "AI会话已结束，您可以输入 #进入AI会话 来重新开始。"
}

func (s *AIChatService) Chat(aiMessages []openai.ChatCompletionMessage) (openai.ChatCompletionMessage, error) {
	aiConfig := s.config.GetAIConfig()

	// 验证AI配置是否完整
	if aiConfig.APIKey == "" {
		return openai.ChatCompletionMessage{}, fmt.Errorf("AI API Key 未配置，请联系管理员")
	}
	if aiConfig.BaseURL == "" {
		return openai.ChatCompletionMessage{}, fmt.Errorf("AI Base URL 未配置，请联系管理员")
	}
	if aiConfig.Model == "" {
		return openai.ChatCompletionMessage{}, fmt.Errorf("AI Model 未配置，请联系管理员")
	}

	log.Printf("AI配置验证通过 - BaseURL: %s, Model: %s", aiConfig.BaseURL, aiConfig.Model)

	if aiConfig.Prompt != "" {
		systemMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: aiConfig.Prompt,
		}
		if aiConfig.MaxCompletionTokens > 0 {
			systemMessage.Content += fmt.Sprintf("\n\n请注意，每次回答不能超过%d个汉字。", aiConfig.MaxCompletionTokens)
		}
		aiMessages = append([]openai.ChatCompletionMessage{systemMessage}, aiMessages...)
	}
	openaiConfig := openai.DefaultConfig(aiConfig.APIKey)

	// 检查是否是Gemini API，如果是则需要特殊处理
	if isGeminiAPI(aiConfig.BaseURL) {
		log.Printf("检测到Gemini API，使用原始配置地址")
		// 对于Gemini代理服务，直接使用原始配置的URL，不添加任何端点
		openaiConfig.BaseURL = aiConfig.BaseURL
	} else {
		openaiConfig.BaseURL = aiConfig.BaseURL
	}

	client := openai.NewClientWithConfig(openaiConfig)
	req := openai.ChatCompletionRequest{
		Model:    aiConfig.Model,
		Messages: aiMessages,
		Stream:   false,
	}

	// 记录详细的请求信息
	log.Printf("=== AIConfig 完整配置 ===")
	log.Printf("BaseURL: %s", aiConfig.BaseURL)
	log.Printf("APIKey: %s...", aiConfig.APIKey[:min(8, len(aiConfig.APIKey))])
	log.Printf("Model: %s", aiConfig.Model)
	log.Printf("WorkflowModel: %s", aiConfig.WorkflowModel)
	log.Printf("ImageRecognitionModel: %s", aiConfig.ImageRecognitionModel)
	log.Printf("Prompt: %s", aiConfig.Prompt)
	log.Printf("MaxCompletionTokens: %d", aiConfig.MaxCompletionTokens)
	log.Printf("ImageModel: %s", aiConfig.ImageModel)
	log.Printf("ImageAISettings: %s", string(aiConfig.ImageAISettings))
	log.Printf("TTSSettings: %s", string(aiConfig.TTSSettings))
	log.Printf("LTTSSettings: %s", string(aiConfig.LTTSSettings))
	log.Printf("=== AIConfig 配置结束 ===")

	log.Printf("AI请求详情:")
	log.Printf("  原始BaseURL: %s", aiConfig.BaseURL)
	log.Printf("  调整后BaseURL: %s", openaiConfig.BaseURL)
	log.Printf("  Model: %s", aiConfig.Model)
	log.Printf("  APIKey: %s...", aiConfig.APIKey[:min(8, len(aiConfig.APIKey))])
	log.Printf("  消息数量: %d", len(aiMessages))
	for i, msg := range aiMessages {
		log.Printf("  消息%d: Role=%s, Content长度=%d", i+1, msg.Role, len(msg.Content))
	}
	// 判断一下aiMessages是否包含图片，如果包含，则使用多模态模型
	for _, msg := range aiMessages {
		if len(msg.MultiContent) > 0 {
			for _, part := range msg.MultiContent {
				if part.Type == openai.ChatMessagePartTypeImageURL {
					req.Model = aiConfig.ImageRecognitionModel
					break
				}
			}
		}
	}
	if aiConfig.MaxCompletionTokens > 0 {
		req.MaxCompletionTokens = aiConfig.MaxCompletionTokens
	}
	log.Printf("开始调用AI服务 - 消息数量: %d", len(aiMessages))
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		log.Printf("AI服务调用失败: %v", err)

		// 如果是Gemini API且返回404，记录详细错误信息
		if isGeminiAPI(aiConfig.BaseURL) && strings.Contains(err.Error(), "404") {
			log.Printf("Gemini API 404错误，请检查代理服务配置")
			log.Printf("当前配置: BaseURL=%s, Model=%s", aiConfig.BaseURL, aiConfig.Model)
		}

		return openai.ChatCompletionMessage{}, fmt.Errorf("AI服务调用失败: %v", err)
	}
	if len(resp.Choices) == 0 {
		log.Printf("AI返回了空内容")
		return openai.ChatCompletionMessage{}, fmt.Errorf("AI返回了空内容，请联系管理员")
	}
	log.Printf("AI服务调用成功，返回内容长度: %d", len(resp.Choices[0].Message.Content))
	return resp.Choices[0].Message, nil
}
