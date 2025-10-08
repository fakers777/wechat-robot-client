package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"wechat-robot-client/interface/plugin"
)

// JdjcjPlugin 京东积存金价格查询插件
type JdjcjPlugin struct {
	config JdjcjConfig
}

// JdjcjConfig 插件配置
type JdjcjConfig struct {
	VoiceReply   bool   `json:"voice_reply"`
	APIBaseURL   string `json:"api_base_url"`
}

// PriceResponse 价格响应结构
type PriceResponse struct {
	ResultCode int `json:"resultCode"`
	ResultData struct {
		Datas struct {
			Price float64 `json:"price"`
			Time  int64   `json:"time"`
		} `json:"datas"`
	} `json:"resultData"`
}

// NewJdjcjPlugin 创建京东积存金插件实例
func NewJdjcjPlugin() plugin.MessageHandler {
	plugin := &JdjcjPlugin{}
	plugin.loadConfig()
	return plugin
}

// GetName 获取插件名称
func (p *JdjcjPlugin) GetName() string {
	return "Jdjcj"
}

// GetLabels 获取插件标签
func (p *JdjcjPlugin) GetLabels() []string {
	return []string{"jdjcj", "gold", "price", "finance"}
}

// PreAction 前置处理
func (p *JdjcjPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *JdjcjPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *JdjcjPlugin) Run(ctx *plugin.MessageContext) bool {
	content := strings.ToLower(strings.TrimSpace(ctx.MessageContent))
	
	// 检查是否是积存金相关命令
	if !p.containsJdjcjKeywords(content) {
		return false
	}
	
	// 处理语音开关命令
	if content == "积存金语音开" || content == "积存金语音打开" {
		p.config.VoiceReply = true
		p.saveConfig()
		p.sendReply(ctx, "text", "已开启积存金语音回复功能")
		return true
	}
	
	if content == "积存金语音关" || content == "积存金语音关闭" {
		p.config.VoiceReply = false
		p.saveConfig()
		p.sendReply(ctx, "text", "已关闭积存金语音回复功能")
		return true
	}
	
	// 处理查询命令
	if content == "jcj" || content == "积存金" || content == "激存金" {
		price, timestamp, err := p.getJdjcjPrice()
		if err != nil {
			p.sendReply(ctx, "text", "获取失败,等待修复⌛️")
			return true
		}
		
		if price != 0 {
			priceTime := time.Unix(timestamp/1000, 0)
			priceText := fmt.Sprintf("💰 京东积存金当前价格: %.2f 元/克\n⏰ 更新时间: %s", 
				price, priceTime.Format("2006-01-02 15:04:05"))
			
			p.sendReply(ctx, "text", priceText)
			
			// 根据配置决定是否使用语音回复
			if p.config.VoiceReply {
				// 这里可以添加语音回复逻辑
				// 暂时用文字表示
				p.sendReply(ctx, "text", "🔊 语音回复: 京东积存金当前价格"+fmt.Sprintf("%.2f元每克", price))
			}
		} else {
			p.sendReply(ctx, "text", "获取失败,等待修复⌛️")
		}
		
		return true
	}
	
	return false
}

// containsJdjcjKeywords 检查是否包含积存金相关关键词
func (p *JdjcjPlugin) containsJdjcjKeywords(content string) bool {
	keywords := []string{
		"jcj", "积存金", "激存金", "京东", "黄金", "金价",
		"语音开", "语音关", "语音打开", "语音关闭",
	}
	
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// loadConfig 加载配置
func (p *JdjcjPlugin) loadConfig() {
	configPath := "plugin/plugins/jdjcj_config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 使用默认配置
		p.config = JdjcjConfig{
			VoiceReply: false,
			APIBaseURL: "https://api.jdjygold.com/gw/generic/hj/h5/m/",
		}
		return
	}
	
	json.Unmarshal(data, &p.config)
}

// saveConfig 保存配置
func (p *JdjcjPlugin) saveConfig() {
	configPath := "plugin/plugins/jdjcj_config.json"
	data, err := json.MarshalIndent(p.config, "", "    ")
	if err != nil {
		return
	}
	
	os.WriteFile(configPath, data, 0644)
}

// sendReply 发送回复
func (p *JdjcjPlugin) sendReply(ctx *plugin.MessageContext, replyType, content string) {
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content)
	}
}

// getJdjcjPrice 获取京东积存金价格
func (p *JdjcjPlugin) getJdjcjPrice() (float64, int64, error) {
	url := p.config.APIBaseURL + "latestPrice"
	
	// 创建请求
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return 0, 0, err
	}
	
	// 设置请求头
	req.Header.Set("Host", "api.jdjygold.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; ELS-AN00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.18 Mobile Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Origin", "https://m.jdjygold.com")
	req.Header.Set("Referer", "https://m.jdjygold.com/finance-gold/newgold/index/?jrcontainer=h5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	
	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	
	// 解析响应
	var priceResp PriceResponse
	err = json.Unmarshal(body, &priceResp)
	if err != nil {
		return 0, 0, err
	}
	
	if priceResp.ResultCode == 0 {
		return priceResp.ResultData.Datas.Price, priceResp.ResultData.Datas.Time, nil
	}
	
	return 0, 0, fmt.Errorf("API返回错误代码: %d", priceResp.ResultCode)
}