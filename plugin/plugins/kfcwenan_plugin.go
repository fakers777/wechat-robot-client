package plugins

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"wechat-robot-client/interface/plugin"
)

// KFCWenanPlugin KFC文案插件
type KFCWenanPlugin struct {
	config KFCWenanConfig
}

// KFCWenanConfig 插件配置
type KFCWenanConfig struct {
	BaseURL string `json:"base_url"`
}

// APIResponse API响应结构
type APIResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// NewKFCWenanPlugin 创建KFC文案插件实例
func NewKFCWenanPlugin() plugin.MessageHandler {
	plugin := &KFCWenanPlugin{}
	plugin.loadConfig()
	return plugin
}

// GetName 获取插件名称
func (p *KFCWenanPlugin) GetName() string {
	return "KFCWenan"
}

// GetLabels 获取插件标签
func (p *KFCWenanPlugin) GetLabels() []string {
	return []string{"text", "kfc", "wenan", "fun", "thursday"}
}

// PreAction 前置处理
func (p *KFCWenanPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *KFCWenanPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *KFCWenanPlugin) Run(ctx *plugin.MessageContext) bool {
	content := strings.ToLower(strings.TrimSpace(ctx.MessageContent))

	// 检查是否是KFC相关命令
	if content == "kfc" || content == "疯狂星期四" {
		result := p.getKFCWenan()
		if result != "" {
			p.sendReply(ctx, "text", result)
		} else {
			p.sendReply(ctx, "text", "获取失败,等待修复⌛️")
		}
		return true
	}

	// 检查是否是舔狗相关命令
	if content == "舔狗" {
		result := p.getDogWenan()
		if result != "" {
			p.sendReply(ctx, "text", result)
		} else {
			p.sendReply(ctx, "text", "获取失败,等待修复⌛️")
		}
		return true
	}

	return false
}

// loadConfig 加载配置
func (p *KFCWenanPlugin) loadConfig() {
	configPath := "plugin/plugins/kfcwenan_config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 使用默认配置
		p.config = KFCWenanConfig{
			BaseURL: "https://api.pearktrue.cn/api/",
		}
		return
	}

	json.Unmarshal(data, &p.config)
}

// sendReply 发送回复
func (p *KFCWenanPlugin) sendReply(ctx *plugin.MessageContext, replyType, content string) {
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content)
	}
}

// getKFCWenan 获取KFC文案
func (p *KFCWenanPlugin) getKFCWenan() string {
	apiURL := p.config.BaseURL + "kfc"
	return p.makeAPIRequest(apiURL)
}

// getDogWenan 获取舔狗文案
func (p *KFCWenanPlugin) getDogWenan() string {
	apiURL := p.config.BaseURL + "dog"
	return p.makeAPIRequest(apiURL)
}

// makeAPIRequest 发送API请求
func (p *KFCWenanPlugin) makeAPIRequest(apiURL string) string {
	// 构建请求URL
	u, err := url.Parse(apiURL)
	if err != nil {
		return ""
	}

	// 添加参数
	params := url.Values{}
	params.Set("type", "json")
	u.RawQuery = params.Encode()

	// 创建请求
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return ""
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != 200 {
		return ""
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// 解析响应
	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return ""
	}

	// 检查响应状态
	if apiResp.Code == 200 && apiResp.Text != "" {
		return apiResp.Text
	}

	return ""
}
