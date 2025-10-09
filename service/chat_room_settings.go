package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"wechat-robot-client/interface/settings"
	"wechat-robot-client/model"
	"wechat-robot-client/pkg/robot"
	"wechat-robot-client/repository"
	"wechat-robot-client/utils"
	"wechat-robot-client/vars"
)

type ChatRoomSettingsService struct {
	ctx              context.Context
	Message          *model.Message
	gsRespo          *repository.GlobalSettings
	crsRespo         *repository.ChatRoomSettings
	globalSettings   *model.GlobalSettings
	chatRoomSettings *model.ChatRoomSettings
}

var _ settings.Settings = (*ChatRoomSettingsService)(nil)

func NewChatRoomSettingsService(ctx context.Context) *ChatRoomSettingsService {
	return &ChatRoomSettingsService{
		ctx:      ctx,
		gsRespo:  repository.NewGlobalSettingsRepo(ctx, vars.DB),
		crsRespo: repository.NewChatRoomSettingsRepo(ctx, vars.DB),
	}
}

func (s *ChatRoomSettingsService) GetChatRoomSettings(chatRoomID string) (*model.ChatRoomSettings, error) {
	return s.crsRespo.GetChatRoomSettings(chatRoomID)
}

func (s *ChatRoomSettingsService) InitByMessage(message *model.Message) error {
	s.Message = message
	globalSettings, err := s.gsRespo.GetGlobalSettings()
	if err != nil {
		log.Printf("获取全局设置失败: %v", err)
		return err
	}
	s.globalSettings = globalSettings

	chatRoomSettings, err := s.crsRespo.GetChatRoomSettings(message.FromWxID)
	if err != nil {
		log.Printf("获取群聊设置失败: %v", err)
		return err
	}

	// 如果群聊设置不存在，创建默认设置
	if chatRoomSettings == nil {
		log.Printf("群聊 %s 设置不存在，创建默认设置", message.FromWxID)
		aiEnabled := true
		aiTrigger := "AI"
		chatRoomSettings = &model.ChatRoomSettings{
			ChatRoomID:    message.FromWxID,
			ChatAIEnabled: &aiEnabled,
			ChatAITrigger: &aiTrigger,
		}
		// 保存默认设置到数据库
		err = s.crsRespo.Create(chatRoomSettings)
		if err != nil {
			log.Printf("保存默认群聊设置失败: %v", err)
			// 即使保存失败，也继续使用默认设置
		} else {
			log.Printf("成功创建群聊 %s 的默认设置", message.FromWxID)
		}
	}

	s.chatRoomSettings = chatRoomSettings
	return nil
}

func (s *ChatRoomSettingsService) GetAIConfig() settings.AIConfig {
	aiConfig := settings.AIConfig{}
	if s.globalSettings != nil {
		if s.globalSettings.ChatBaseURL != "" {
			aiConfig.BaseURL = s.globalSettings.ChatBaseURL
		}
		if s.globalSettings.ChatAPIKey != "" {
			aiConfig.APIKey = s.globalSettings.ChatAPIKey
		}
		if s.globalSettings.ChatModel != "" {
			aiConfig.Model = s.globalSettings.ChatModel
		}
		if s.globalSettings.WorkflowModel != "" {
			aiConfig.WorkflowModel = s.globalSettings.WorkflowModel
		}
		if s.globalSettings.ImageRecognitionModel != "" {
			aiConfig.ImageRecognitionModel = s.globalSettings.ImageRecognitionModel
		}
		if s.globalSettings.ChatPrompt != "" {
			aiConfig.Prompt = s.globalSettings.ChatPrompt
		}
		if s.globalSettings.MaxCompletionTokens != nil {
			aiConfig.MaxCompletionTokens = *s.globalSettings.MaxCompletionTokens
		}
		if s.globalSettings.ImageModel != "" {
			aiConfig.ImageModel = s.globalSettings.ImageModel
		}
		if s.globalSettings.ImageAISettings != nil {
			aiConfig.ImageAISettings = s.globalSettings.ImageAISettings
		}
		if s.globalSettings.TTSSettings != nil {
			aiConfig.TTSSettings = s.globalSettings.TTSSettings
		}
		if s.globalSettings.LTTSSettings != nil {
			aiConfig.LTTSSettings = s.globalSettings.LTTSSettings
		}
	}
	if s.chatRoomSettings != nil {
		if s.chatRoomSettings.ChatBaseURL != nil && *s.chatRoomSettings.ChatBaseURL != "" {
			aiConfig.BaseURL = *s.chatRoomSettings.ChatBaseURL
		}
		if s.chatRoomSettings.ChatAPIKey != nil && *s.chatRoomSettings.ChatAPIKey != "" {
			aiConfig.APIKey = *s.chatRoomSettings.ChatAPIKey
		}
		if s.chatRoomSettings.ChatModel != nil && *s.chatRoomSettings.ChatModel != "" {
			aiConfig.Model = *s.chatRoomSettings.ChatModel
		}
		if s.chatRoomSettings.WorkflowModel != nil && *s.chatRoomSettings.WorkflowModel != "" {
			aiConfig.WorkflowModel = *s.chatRoomSettings.WorkflowModel
		}
		if s.chatRoomSettings.ImageRecognitionModel != nil && *s.chatRoomSettings.ImageRecognitionModel != "" {
			aiConfig.ImageRecognitionModel = *s.chatRoomSettings.ImageRecognitionModel
		}
		if s.chatRoomSettings.ChatPrompt != nil && *s.chatRoomSettings.ChatPrompt != "" {
			aiConfig.Prompt = *s.chatRoomSettings.ChatPrompt
		}
		if s.chatRoomSettings.MaxCompletionTokens != nil {
			aiConfig.MaxCompletionTokens = *s.chatRoomSettings.MaxCompletionTokens
		}
		if s.chatRoomSettings.ImageModel != nil && *s.chatRoomSettings.ImageModel != "" {
			aiConfig.ImageModel = *s.chatRoomSettings.ImageModel
		}
		if s.chatRoomSettings.ImageAISettings != nil {
			aiConfig.ImageAISettings = s.chatRoomSettings.ImageAISettings
		}
		if s.chatRoomSettings.TTSSettings != nil {
			aiConfig.TTSSettings = s.chatRoomSettings.TTSSettings
		}
		if s.chatRoomSettings.LTTSSettings != nil {
			aiConfig.LTTSSettings = s.chatRoomSettings.LTTSSettings
		}
	}
	aiConfig.BaseURL = utils.NormalizeAIBaseURL(aiConfig.BaseURL)
	return aiConfig
}

func (s *ChatRoomSettingsService) IsAIChatEnabled() bool {
	if s.chatRoomSettings != nil && s.chatRoomSettings.ChatAIEnabled != nil {
		return *s.chatRoomSettings.ChatAIEnabled
	}
	if s.globalSettings != nil && s.globalSettings.ChatAIEnabled != nil {
		return *s.globalSettings.ChatAIEnabled
	}
	return false
}

func (s *ChatRoomSettingsService) IsAIDrawingEnabled() bool {
	if s.chatRoomSettings != nil && s.chatRoomSettings.ImageAIEnabled != nil {
		return *s.chatRoomSettings.ImageAIEnabled
	}
	if s.globalSettings != nil && s.globalSettings.ImageAIEnabled != nil {
		return *s.globalSettings.ImageAIEnabled
	}
	return false
}

func (s *ChatRoomSettingsService) IsTTSEnabled() bool {
	if s.chatRoomSettings != nil && s.chatRoomSettings.TTSEnabled != nil {
		return *s.chatRoomSettings.TTSEnabled
	}
	if s.globalSettings != nil && s.globalSettings.TTSEnabled != nil {
		return *s.globalSettings.TTSEnabled
	}
	return false
}

// 是否属于自动触发AI的指令
func (s *ChatRoomSettingsService) IsAutoAITrigger(message string) bool {
	matched, _ := NewAIWorkflowService(s.ctx, s).ChatIntentionSimple(message, nil)
	return matched
}

func (s *ChatRoomSettingsService) IsAITrigger() bool {
	messageContent := s.Message.Content
	if s.Message.AppMsgType == model.AppMsgTypequote {
		var xmlMessage robot.XmlMessage
		if err := vars.RobotRuntime.XmlDecoder(messageContent, &xmlMessage); err == nil {
			messageContent = xmlMessage.AppMsg.Title
		}
	}
	if s.Message.IsAtMe {
		// 是否是 @所有人
		atAllRegex := regexp.MustCompile(vars.AtAllRegexp)
		if atAllRegex.MatchString(messageContent) {
			// 如果是 @所有人，则不处理
			return false
		}
		return true
	}
	if s.chatRoomSettings == nil {
		if s.globalSettings == nil {
			return false
		}
		if s.globalSettings.ChatAIEnabled == nil || !*s.globalSettings.ChatAIEnabled {
			return false
		}
		if s.IsAutoAITrigger(messageContent) {
			return true
		}
		return *s.globalSettings.ChatAITrigger != "" && strings.HasPrefix(messageContent, *s.globalSettings.ChatAITrigger)
	}
	if s.chatRoomSettings.ChatAIEnabled == nil || !*s.chatRoomSettings.ChatAIEnabled {
		return false
	}
	if s.IsAutoAITrigger(messageContent) {
		return true
	}
	if s.chatRoomSettings.ChatAITrigger != nil && *s.chatRoomSettings.ChatAITrigger != "" {
		return *s.chatRoomSettings.ChatAITrigger != "" && strings.HasPrefix(messageContent, *s.chatRoomSettings.ChatAITrigger)
	}
	return s.globalSettings != nil && s.globalSettings.ChatAITrigger != nil && *s.globalSettings.ChatAITrigger != "" &&
		strings.HasPrefix(messageContent, *s.globalSettings.ChatAITrigger)
}

func (s *ChatRoomSettingsService) GetAITriggerWord() string {
	if s.chatRoomSettings != nil && s.chatRoomSettings.ChatAITrigger != nil && *s.chatRoomSettings.ChatAITrigger != "" {
		return *s.chatRoomSettings.ChatAITrigger
	}
	if s.globalSettings != nil && s.globalSettings.ChatAITrigger != nil && *s.globalSettings.ChatAITrigger != "" {
		return *s.globalSettings.ChatAITrigger
	}
	return ""
}

func (s *ChatRoomSettingsService) GetChatRoomWelcomeConfig(chatRoomID string) (*model.ChatRoomSettings, error) {
	globalSettings, err := s.gsRespo.GetGlobalSettings()
	if err != nil {
		return nil, err
	}
	if globalSettings == nil {
		return nil, fmt.Errorf("加载全局配置失败")
	}
	chatRoomSetting, err := s.crsRespo.GetChatRoomSettings(chatRoomID)
	if err != nil {
		return nil, err
	}
	if chatRoomSetting == nil {
		return &model.ChatRoomSettings{
			WelcomeEnabled:  globalSettings.WelcomeEnabled,
			WelcomeType:     globalSettings.WelcomeType,
			WelcomeText:     globalSettings.WelcomeText,
			WelcomeEmojiMD5: globalSettings.WelcomeEmojiMD5,
			WelcomeEmojiLen: globalSettings.WelcomeEmojiLen,
			WelcomeImageURL: globalSettings.WelcomeImageURL,
			WelcomeURL:      globalSettings.WelcomeURL,
		}, nil
	}
	return chatRoomSetting, nil
}

func (s *ChatRoomSettingsService) GetPatConfig() settings.PatConfig {
	if s.chatRoomSettings != nil {
		if s.chatRoomSettings.PatEnabled != nil {
			return settings.PatConfig{
				PatEnabled:     *s.chatRoomSettings.PatEnabled,
				PatType:        s.chatRoomSettings.PatType,
				PatText:        s.chatRoomSettings.PatText,
				PatVoiceTimbre: s.chatRoomSettings.PatVoiceTimbre,
			}
		}
	}
	if s.globalSettings != nil {
		if s.globalSettings.PatEnabled != nil {
			return settings.PatConfig{
				PatEnabled:     *s.globalSettings.PatEnabled,
				PatType:        s.globalSettings.PatType,
				PatText:        s.globalSettings.PatText,
				PatVoiceTimbre: s.globalSettings.PatVoiceTimbre,
			}
		}
	}
	return settings.PatConfig{}
}

func (s *ChatRoomSettingsService) GetLeaveChatRoomConfig(chatRoomID string) *model.ChatRoomSettings {
	globalSettings, err := s.gsRespo.GetGlobalSettings()
	if err != nil {
		return nil
	}
	chatRoomSettings, err := s.crsRespo.GetChatRoomSettings(chatRoomID)
	if err != nil {
		return nil
	}
	if chatRoomSettings != nil {
		return chatRoomSettings
	}
	if globalSettings != nil {
		return &model.ChatRoomSettings{
			LeaveChatRoomAlertEnabled: globalSettings.LeaveChatRoomAlertEnabled,
			LeaveChatRoomAlertText:    globalSettings.LeaveChatRoomAlertText,
		}
	}
	return nil
}

func (s *ChatRoomSettingsService) GetAllEnableChatRank() ([]*model.ChatRoomSettings, error) {
	if vars.RobotRuntime.Status == model.RobotStatusOffline {
		return []*model.ChatRoomSettings{}, nil
	}
	return s.crsRespo.GetAllEnableChatRank()
}

func (s *ChatRoomSettingsService) GetAllEnableAISummary() ([]*model.ChatRoomSettings, error) {
	if vars.RobotRuntime.Status == model.RobotStatusOffline {
		return []*model.ChatRoomSettings{}, nil
	}
	return s.crsRespo.GetAllEnableAISummary()
}

func (s *ChatRoomSettingsService) GetAllEnableGoodMorning() ([]*model.ChatRoomSettings, error) {
	if vars.RobotRuntime.Status == model.RobotStatusOffline {
		return []*model.ChatRoomSettings{}, nil
	}
	return s.crsRespo.GetAllEnableGoodMorning()
}

func (s *ChatRoomSettingsService) GetAllEnableNews() ([]*model.ChatRoomSettings, error) {
	if vars.RobotRuntime.Status == model.RobotStatusOffline {
		return []*model.ChatRoomSettings{}, nil
	}
	return s.crsRespo.GetAllEnableNews()
}

func (s *ChatRoomSettingsService) SaveChatRoomSettings(data *model.ChatRoomSettings) error {
	if data.ID == 0 {
		return s.crsRespo.Create(data)
	}
	return s.crsRespo.Update(data)
}
