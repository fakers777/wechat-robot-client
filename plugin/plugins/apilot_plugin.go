package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"wechat-robot-client/interface/plugin"
)

// ApilotPlugin Apilot多功能插件
type ApilotPlugin struct {
	config Config
}

// Config 插件配置
type Config struct {
	AlapiToken            string `json:"alapi_token"`
	MorningNewsTextEnabled bool   `json:"morning_news_text_enabled"`
	BaseURLVVHan          string `json:"base_url_vvhan"`
	BaseURLAlapi          string `json:"base_url_alapi"`
}

// ZodiacMapping 星座映射
var ZodiacMapping = map[string]string{
	"白羊座": "aries",
	"金牛座": "taurus",
	"双子座": "gemini",
	"巨蟹座": "cancer",
	"狮子座": "leo",
	"处女座": "virgo",
	"天秤座": "libra",
	"天蝎座": "scorpio",
	"射手座": "sagittarius",
	"摩羯座": "capricorn",
	"水瓶座": "aquarius",
	"双鱼座": "pisces",
}

// HotTrendTypes 热榜类型映射
var HotTrendTypes = map[string]string{
	"微博":     "wbHot",
	"虎扑":     "huPu",
	"知乎":     "zhihuHot",
	"知乎日报":  "zhihuDay",
	"哔哩哔哩":  "bili",
	"36氪":    "36Ke",
	"抖音":     "douyinHot",
	"IT":      "itNews",
	"虎嗅":     "huXiu",
	"产品经理":  "woShiPm",
	"头条":     "toutiao",
	"百度":     "baiduRD",
	"豆瓣":     "douban",
}

// NewApilotPlugin 创建Apilot插件实例
func NewApilotPlugin() plugin.MessageHandler {
	plugin := &ApilotPlugin{}
	plugin.loadConfig()
	return plugin
}

// GetName 获取插件名称
func (p *ApilotPlugin) GetName() string {
	return "Apilot"
}

// GetLabels 获取插件标签
func (p *ApilotPlugin) GetLabels() []string {
	return []string{"apilot", "news", "weather", "horoscope", "express", "hot", "fun"}
}

// PreAction 前置处理
func (p *ApilotPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *ApilotPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *ApilotPlugin) Run(ctx *plugin.MessageContext) bool {
	content := strings.TrimSpace(ctx.MessageContent)
	
	// 早报
	if content == "早报" {
		news := p.getMorningNews()
		replyType := "text"
		if p.isValidURL(news) {
			replyType = "image_url"
		}
		p.sendReply(ctx, replyType, news)
		return true
	}
	
	// 摸鱼
	if content == "摸鱼" {
		moyu := p.getMoyuCalendar()
		replyType := "text"
		if p.isValidURL(moyu) {
			replyType = "image_url"
		}
		p.sendReply(ctx, replyType, moyu)
		return true
	}
	
	// 摸鱼视频
	if content == "摸鱼视频" {
		moyu := p.getMoyuCalendarVideo()
		replyType := "text"
		if p.isValidURL(moyu) {
			replyType = "video_url"
		}
		p.sendReply(ctx, replyType, moyu)
		return true
	}
	
	// 八卦
	if content == "八卦" {
		bagua := p.getMxBagua()
		replyType := "text"
		if p.isValidURL(bagua) {
			replyType = "image_url"
		}
		p.sendReply(ctx, replyType, bagua)
		return true
	}
	
	// 举牌
	if strings.Contains(content, "举牌") {
		parts := strings.Split(content, " ")
		if len(parts) > 1 {
			jupai := p.getJupaiPic(parts[1])
			if jupai != "" {
				replyType := "text"
				if p.isValidURL(jupai) {
					replyType = "image_url"
				}
				p.sendReply(ctx, replyType, jupai)
			}
		}
		return true
	}
	
	// 快递查询
	if strings.HasPrefix(content, "快递") {
		trackingNumber := strings.TrimSpace(content[2:])
		trackingNumber = strings.ReplaceAll(trackingNumber, "：", ":")
		
		if p.config.AlapiToken == "" {
			p.sendReply(ctx, "text", "请先配置alapi的token")
			return true
		}
		
		// 检查顺丰快递格式
		if strings.HasPrefix(trackingNumber, "SF") && !strings.Contains(trackingNumber, ":") {
			p.sendReply(ctx, "text", "顺丰快递需要补充寄/收件人手机号后四位，格式：SF12345:0000")
			return true
		}
		
		result := p.queryExpressInfo(trackingNumber)
		p.sendReply(ctx, "text", result)
		return true
	}
	
	// 星座查询
	if zodiacEnglish, exists := ZodiacMapping[content]; exists {
		result := p.getHoroscope(zodiacEnglish)
		p.sendReply(ctx, "text", result)
		return true
	}
	
	// 热榜查询
	hotTrendMatch := regexp.MustCompile(`(.{1,6})热榜$`)
	if matches := hotTrendMatch.FindStringSubmatch(content); len(matches) > 1 {
		hotTrendsType := strings.TrimSpace(matches[1])
		result := p.getHotTrends(hotTrendsType)
		p.sendReply(ctx, "text", result)
		return true
	}
	
	// 天气查询
	weatherMatch := regexp.MustCompile(`^(?:(.{2,7}?)(?:市|县|区|镇)?|(\d{7,9}))(:?今天|明天|后天|7天|七天)?(?:的)?天气$`)
	if matches := weatherMatch.FindStringSubmatch(content); len(matches) > 1 {
		cityOrID := matches[1]
		if cityOrID == "" {
			cityOrID = matches[2]
		}
		date := matches[3]
		
		if p.config.AlapiToken == "" {
			p.sendReply(ctx, "text", "请先配置alapi的token")
			return true
		}
		
		result := p.getWeather(cityOrID, date, content)
		p.sendReply(ctx, "text", result)
		return true
	}
	
	return false
}

// loadConfig 加载配置
func (p *ApilotPlugin) loadConfig() {
	configPath := "plugin/plugins/apilot_config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 使用默认配置
		p.config = Config{
			AlapiToken:            "",
			MorningNewsTextEnabled: false,
			BaseURLVVHan:          "https://api.vvhan.com/api/",
			BaseURLAlapi:          "https://v2.alapi.cn/api/",
		}
		return
	}
	
	json.Unmarshal(data, &p.config)
}

// sendReply 发送回复
func (p *ApilotPlugin) sendReply(ctx *plugin.MessageContext, replyType, content string) {
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content)
	}
}

// makeRequest 发送HTTP请求
func (p *ApilotPlugin) makeRequest(url, method string, headers map[string]string, data url.Values) (map[string]interface{}, error) {
	var req *http.Request
	var err error
	
	if method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		if len(data) > 0 {
			req, err = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		} else {
			req, err = http.NewRequest("POST", url, nil)
		}
	}
	
	if err != nil {
		return nil, err
	}
	
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	return result, err
}

// isValidURL 检查URL是否有效
func (p *ApilotPlugin) isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// getMorningNews 获取早报
func (p *ApilotPlugin) getMorningNews() string {
	if p.config.AlapiToken == "" {
		// 使用免费API
		apiURL := p.config.BaseURLAlapi + "zaobao"
		headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		data := make(url.Values)
		data.Set("format", "json")
		
		result, err := p.makeRequest(apiURL, "POST", headers, data)
		if err != nil {
			return "早报获取失败，请稍后再试"
		}
		
		if success, ok := result["success"].(bool); ok && success {
			if p.config.MorningNewsTextEnabled {
				if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
					var newsList []string
					for i, news := range data[:len(data)-1] {
						newsList = append(newsList, fmt.Sprintf("%d. %s", i+1, news))
					}
					formattedNews := fmt.Sprintf("☕ %s  今日早报\n\n%s", data[len(data)-1], strings.Join(newsList, "\n"))
					if imgURL, ok := result["imgUrl"].(string); ok {
						return fmt.Sprintf("%s\n\n图片url：%s", formattedNews, imgURL)
					}
					return formattedNews
				}
			} else {
				if imgURL, ok := result["imgUrl"].(string); ok {
					return imgURL
				}
			}
		}
		return "早报信息获取失败，可配置\"alapi token\"切换至 Alapi 服务，或者稍后再试"
	} else {
		// 使用Alapi
		apiURL := p.config.BaseURLAlapi + "zaobao"
		headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		data := make(url.Values)
		data.Set("token", p.config.AlapiToken)
		data.Set("format", "json")
		
		result, err := p.makeRequest(apiURL, "POST", headers, data)
		if err != nil {
			return "早报获取失败"
		}
		
		if code, ok := result["code"].(float64); ok && code == 200 {
			if data, ok := result["data"].(map[string]interface{}); ok {
				if imgURL, ok := data["image"].(string); ok {
					if p.config.MorningNewsTextEnabled {
						if news, ok := data["news"].([]interface{}); ok {
							if weiyu, ok := data["weiyu"].(string); ok {
								if date, ok := data["date"].(string); ok {
									var newsList []string
									for _, n := range news {
										newsList = append(newsList, fmt.Sprintf("• %s", n))
									}
									formattedNews := fmt.Sprintf("☕ %s  今日早报\n\n%s\n\n%s\n\n图片url：%s", date, strings.Join(newsList, "\n"), weiyu, imgURL)
									return formattedNews
								}
							}
						}
					}
					return imgURL
				}
			}
		}
		return "早报获取失败，请检查 token 是否有误"
	}
}

// getMoyuCalendar 获取摸鱼日历
func (p *ApilotPlugin) getMoyuCalendar() string {
	apiURL := p.config.BaseURLVVHan + "moyu?type=json"
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	data := make(url.Values)
	data.Set("format", "json")
	
	result, err := p.makeRequest(apiURL, "POST", headers, data)
	if err != nil {
		return "摸鱼日历获取失败"
	}
	
	if success, ok := result["success"].(bool); ok && success {
		if url, ok := result["url"].(string); ok {
			return url
		}
	}
	
	// 备用API
	apiURL = "https://dayu.qqsuu.cn/moyuribao/apis.php?type=json"
	result, err = p.makeRequest(apiURL, "POST", headers, data)
	if err != nil {
		return "暂无可用\"摸鱼\"服务，认真上班"
	}
	
	if code, ok := result["code"].(float64); ok && code == 200 {
		if data, ok := result["data"].(string); ok {
			if p.isValidURL(data) {
				return data
			}
		}
	}
	
	return "周末无需摸鱼，愉快玩耍吧"
}

// getMoyuCalendarVideo 获取摸鱼视频
func (p *ApilotPlugin) getMoyuCalendarVideo() string {
	apiURL := "https://dayu.qqsuu.cn/moyuribaoshipin/apis.php?type=json"
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	data := make(url.Values)
	data.Set("format", "json")
	
	result, err := p.makeRequest(apiURL, "POST", headers, data)
	if err != nil {
		return "视频版没了，看看文字版吧"
	}
	
	if code, ok := result["code"].(float64); ok && code == 200 {
		if data, ok := result["data"].(string); ok {
			if p.isValidURL(data) {
				return data
			}
		}
	}
	
	return "视频版没了，看看文字版吧"
}

// getMxBagua 获取明星八卦
func (p *ApilotPlugin) getMxBagua() string {
	apiURL := "https://dayu.qqsuu.cn/mingxingbagua/apis.php?type=json"
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	data := make(url.Values)
	data.Set("format", "json")
	
	result, err := p.makeRequest(apiURL, "POST", headers, data)
	if err != nil {
		return "暂无明星八卦，吃瓜莫急"
	}
	
	if code, ok := result["code"].(float64); ok && code == 200 {
		if data, ok := result["data"].(string); ok {
			if p.isValidURL(data) {
				return data
			}
		}
	}
	
	return "周末不更新，请微博吃瓜"
}

// getJupaiPic 获取举牌图片
func (p *ApilotPlugin) getJupaiPic(keyword string) string {
	if len(keyword) >= 20 {
		return ""
	}
	
	apiURL := "https://api.andeer.top/API/jupai.php?text=" + url.QueryEscape(keyword)
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	data := make(url.Values)
	data.Set("format", "json")
	
	_, err := p.makeRequest(apiURL, "GET", headers, data)
	if err != nil {
		return ""
	}
	
	return apiURL
}

// getHoroscope 获取星座运势
func (p *ApilotPlugin) getHoroscope(astroSign string) string {
	if p.config.AlapiToken == "" {
		// 使用免费API
		apiURL := p.config.BaseURLVVHan + "horoscope"
		headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		data := make(url.Values)
		data.Set("type", astroSign)
		data.Set("time", "today")
		
		result, err := p.makeRequest(apiURL, "GET", headers, data)
		if err != nil {
			return "星座信息获取失败，可配置\"alapi token\"切换至 Alapi 服务，或者稍后再试"
		}
		
		if success, ok := result["success"].(bool); ok && success {
			if data, ok := result["data"].(map[string]interface{}); ok {
				return p.formatHoroscopeData(data)
			}
		}
		return "星座信息获取失败，可配置\"alapi token\"切换至 Alapi 服务，或者稍后再试"
	} else {
		// 使用Alapi
		apiURL := p.config.BaseURLVVHan + "star"
		headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		data := make(url.Values)
		data.Set("token", p.config.AlapiToken)
		data.Set("star", astroSign)
		
		result, err := p.makeRequest(apiURL, "POST", headers, data)
		if err != nil {
			return "星座获取信息获取失败，请检查 token 是否有误"
		}
		
		if code, ok := result["code"].(float64); ok && code == 200 {
			if data, ok := result["data"].(map[string]interface{}); ok {
				if day, ok := data["day"].(map[string]interface{}); ok {
					return p.formatAlapiHoroscopeData(day)
				}
			}
		}
		return "星座获取信息获取失败，请检查 token 是否有误"
	}
}

// formatHoroscopeData 格式化星座数据
func (p *ApilotPlugin) formatHoroscopeData(data map[string]interface{}) string {
	var result strings.Builder
	
	if title, ok := data["title"].(string); ok {
		if time, ok := data["time"].(string); ok {
			result.WriteString(fmt.Sprintf("%s (%s):\n\n", title, time))
		}
	}
	
	// 每日建议
	if todo, ok := data["todo"].(map[string]interface{}); ok {
		result.WriteString("💡【每日建议】\n")
		if yi, ok := todo["yi"].(string); ok {
			result.WriteString(fmt.Sprintf("宜：%s\n", yi))
		}
		if ji, ok := todo["ji"].(string); ok {
			result.WriteString(fmt.Sprintf("忌：%s\n\n", ji))
		}
	}
	
	// 运势指数
	if index, ok := data["index"].(map[string]interface{}); ok {
		result.WriteString("📊【运势指数】\n")
		if all, ok := index["all"].(string); ok {
			result.WriteString(fmt.Sprintf("总运势：%s\n", all))
		}
		if love, ok := index["love"].(string); ok {
			result.WriteString(fmt.Sprintf("爱情：%s\n", love))
		}
		if work, ok := index["work"].(string); ok {
			result.WriteString(fmt.Sprintf("工作：%s\n", work))
		}
		if money, ok := index["money"].(string); ok {
			result.WriteString(fmt.Sprintf("财运：%s\n", money))
		}
		if health, ok := index["health"].(string); ok {
			result.WriteString(fmt.Sprintf("健康：%s\n\n", health))
		}
	}
	
	// 幸运提示
	if luckyNumber, ok := data["luckynumber"].(string); ok {
		result.WriteString("🍀【幸运提示】\n")
		result.WriteString(fmt.Sprintf("数字：%s\n", luckyNumber))
	}
	if luckyColor, ok := data["luckycolor"].(string); ok {
		result.WriteString(fmt.Sprintf("颜色：%s\n", luckyColor))
	}
	if luckyConstellation, ok := data["luckyconstellation"].(string); ok {
		result.WriteString(fmt.Sprintf("星座：%s\n\n", luckyConstellation))
	}
	
	// 简评
	if shortComment, ok := data["shortcomment"].(string); ok {
		result.WriteString(fmt.Sprintf("✍【简评】\n%s\n\n", shortComment))
	}
	
	// 详细运势
	if fortuneText, ok := data["fortunetext"].(map[string]interface{}); ok {
		result.WriteString("📜【详细运势】\n")
		if all, ok := fortuneText["all"].(string); ok {
			result.WriteString(fmt.Sprintf("总运：%s\n", all))
		}
		if love, ok := fortuneText["love"].(string); ok {
			result.WriteString(fmt.Sprintf("爱情：%s\n", love))
		}
		if work, ok := fortuneText["work"].(string); ok {
			result.WriteString(fmt.Sprintf("工作：%s\n", work))
		}
		if money, ok := fortuneText["money"].(string); ok {
			result.WriteString(fmt.Sprintf("财运：%s\n", money))
		}
		if health, ok := fortuneText["health"].(string); ok {
			result.WriteString(fmt.Sprintf("健康：%s\n", health))
		}
	}
	
	return result.String()
}

// formatAlapiHoroscopeData 格式化Alapi星座数据
func (p *ApilotPlugin) formatAlapiHoroscopeData(data map[string]interface{}) string {
	var result strings.Builder
	
	if date, ok := data["date"].(string); ok {
		result.WriteString(fmt.Sprintf("📅 日期：%s\n\n", date))
	}
	
	// 每日建议
	if yi, ok := data["yi"].(string); ok {
		result.WriteString("💡【每日建议】\n")
		result.WriteString(fmt.Sprintf("宜：%s\n", yi))
	}
	if ji, ok := data["ji"].(string); ok {
		result.WriteString(fmt.Sprintf("忌：%s\n\n", ji))
	}
	
	// 运势指数
	result.WriteString("📊【运势指数】\n")
	if all, ok := data["all"].(string); ok {
		result.WriteString(fmt.Sprintf("总运势：%s\n", all))
	}
	if love, ok := data["love"].(string); ok {
		result.WriteString(fmt.Sprintf("爱情：%s\n", love))
	}
	if work, ok := data["work"].(string); ok {
		result.WriteString(fmt.Sprintf("工作：%s\n", work))
	}
	if money, ok := data["money"].(string); ok {
		result.WriteString(fmt.Sprintf("财运：%s\n", money))
	}
	if health, ok := data["health"].(string); ok {
		result.WriteString(fmt.Sprintf("健康：%s\n\n", health))
	}
	
	// 提醒
	if notice, ok := data["notice"].(string); ok {
		result.WriteString(fmt.Sprintf("🔔【提醒】：%s\n\n", notice))
	}
	
	// 幸运提示
	result.WriteString("🍀【幸运提示】\n")
	if luckyNumber, ok := data["lucky_number"].(string); ok {
		result.WriteString(fmt.Sprintf("数字：%s\n", luckyNumber))
	}
	if luckyColor, ok := data["lucky_color"].(string); ok {
		result.WriteString(fmt.Sprintf("颜色：%s\n", luckyColor))
	}
	if luckyStar, ok := data["lucky_star"].(string); ok {
		result.WriteString(fmt.Sprintf("星座：%s\n\n", luckyStar))
	}
	
	// 简评
	result.WriteString("✍【简评】\n")
	if allText, ok := data["all_text"].(string); ok {
		result.WriteString(fmt.Sprintf("总运：%s\n", allText))
	}
	if loveText, ok := data["love_text"].(string); ok {
		result.WriteString(fmt.Sprintf("爱情：%s\n", loveText))
	}
	if workText, ok := data["work_text"].(string); ok {
		result.WriteString(fmt.Sprintf("工作：%s\n", workText))
	}
	if moneyText, ok := data["money_text"].(string); ok {
		result.WriteString(fmt.Sprintf("财运：%s\n", moneyText))
	}
	if healthText, ok := data["health_text"].(string); ok {
		result.WriteString(fmt.Sprintf("健康：%s\n", healthText))
	}
	
	return result.String()
}

// getHotTrends 获取热榜
func (p *ApilotPlugin) getHotTrends(hotTrendsType string) string {
	hotTrendsTypeEn, exists := HotTrendTypes[hotTrendsType]
	if !exists {
		var supportedTypes []string
		for k := range HotTrendTypes {
			supportedTypes = append(supportedTypes, k)
		}
		return fmt.Sprintf("👉 已支持的类型有：\n\n    %s\n\n📝 请按照以下格式发送：\n    类型+热榜  例如：微博热榜", strings.Join(supportedTypes, "/"))
	}
	
	apiURL := p.config.BaseURLVVHan + "hotlist/" + hotTrendsTypeEn
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	
	result, err := p.makeRequest(apiURL, "GET", headers, nil)
	if err != nil {
		return "热榜获取失败，请稍后再试"
	}
	
	if success, ok := result["success"].(bool); ok && success {
		var output []string
		if updateTime, ok := result["update_time"].(string); ok {
			output = append(output, fmt.Sprintf("更新时间：%s\n", updateTime))
		}
		
		if topics, ok := result["data"].([]interface{}); ok {
			for i, topic := range topics {
				if i >= 15 {
					break
				}
				if topicMap, ok := topic.(map[string]interface{}); ok {
					title := ""
					hot := "无热度参数, 0"
					url := ""
					
					if t, ok := topicMap["title"].(string); ok {
						title = t
					}
					if h, ok := topicMap["hot"].(string); ok {
						hot = h
					}
					if u, ok := topicMap["url"].(string); ok {
						url = u
					}
					
					output = append(output, fmt.Sprintf("%d. %s (%s 浏览)\nURL: %s\n", i+1, title, hot, url))
				}
			}
		}
		
		return strings.Join(output, "\n")
	}
	
	return "热榜获取失败，请稍后再试"
}

// queryExpressInfo 查询快递信息
func (p *ApilotPlugin) queryExpressInfo(trackingNumber string) string {
	apiURL := p.config.BaseURLVVHan + "kd"
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	data := make(url.Values)
	data.Set("token", p.config.AlapiToken)
	data.Set("number", trackingNumber)
	data.Set("com", "")
	data.Set("order", "asc")
	
	result, err := p.makeRequest(apiURL, "POST", headers, data)
	if err != nil {
		return "快递查询失败"
	}
	
	if code, ok := result["code"].(float64); ok && code == 200 {
		if data, ok := result["data"].(map[string]interface{}); ok {
			var formattedResult []string
			
			if nu, ok := data["nu"].(string); ok {
				formattedResult = append(formattedResult, fmt.Sprintf("快递编号：%s", nu))
			}
			if com, ok := data["com"].(string); ok {
				formattedResult = append(formattedResult, fmt.Sprintf("快递公司：%s", com))
			}
			if statusDesc, ok := data["status_desc"].(string); ok {
				formattedResult = append(formattedResult, fmt.Sprintf("状态：%s", statusDesc))
			}
			formattedResult = append(formattedResult, "状态信息：")
			
			if info, ok := data["info"].([]interface{}); ok {
				for _, item := range info {
					if itemMap, ok := item.(map[string]interface{}); ok {
						time := ""
						statusDesc := ""
						content := ""
						
						if t, ok := itemMap["time"].(string); ok {
							time = t[5:11] // 提取时间部分
						}
						if s, ok := itemMap["status_desc"].(string); ok {
							statusDesc = s
						}
						if c, ok := itemMap["content"].(string); ok {
							content = c
						}
						
						formattedResult = append(formattedResult, fmt.Sprintf("%s - %s\n    %s", time, statusDesc, content))
					}
				}
			}
			
			return strings.Join(formattedResult, "\n")
		}
	}
	
	if msg, ok := result["msg"].(string); ok {
		return fmt.Sprintf("查询失败，%s", msg)
	}
	
	return "查询失败：api响应为空"
}

// getWeather 获取天气信息
func (p *ApilotPlugin) getWeather(cityOrID, date, content string) string {
	apiURL := p.config.BaseURLAlapi + "tianqi"
	isFuture := strings.Contains(date, "明天") || strings.Contains(date, "后天") || strings.Contains(date, "七天") || strings.Contains(date, "7天")
	if isFuture {
		apiURL = p.config.BaseURLVVHan + "tianqi/seven"
	}
	
	var data url.Values
	if regexp.MustCompile(`^\d+$`).MatchString(cityOrID) {
		// 使用城市ID
		data = make(url.Values)
		data.Set("city_id", cityOrID)
		data.Set("token", p.config.AlapiToken)
	} else {
		// 使用城市名称
		data = make(url.Values)
		data.Set("city", cityOrID)
		data.Set("token", p.config.AlapiToken)
	}
	
	result, err := p.makeRequest(apiURL, "GET", nil, data)
	if err != nil {
		return "获取天气信息失败"
	}
	
	if code, ok := result["code"].(float64); ok && code == 200 {
		if data, ok := result["data"].(map[string]interface{}); ok {
			if isFuture {
				// 对于未来天气，data应该是数组
				if futureData, ok := result["data"].([]interface{}); ok {
					return p.formatFutureWeather(futureData, date)
				}
			} else {
				return p.formatCurrentWeather(data, cityOrID, content)
			}
		}
	}
	
	return "获取失败，请查看服务器log"
}

// formatCurrentWeather 格式化当前天气
func (p *ApilotPlugin) formatCurrentWeather(data map[string]interface{}, cityOrID, content string) string {
	var result strings.Builder
	
	if city, ok := data["city"].(string); ok {
		if province, ok := data["province"].(string); ok {
			result.WriteString(fmt.Sprintf("🏙️ 城市: %s (%s)\n", city, province))
		}
		
		// 检查城市是否匹配
		if !regexp.MustCompile(`^\d+$`).MatchString(cityOrID) && !strings.Contains(content, city) {
			return "输入不规范，请输<国内城市+(今天|明天|后天|七天|7天)+天气>，比如 '广州天气'"
		}
	}
	
	if updateTime, ok := data["update_time"].(string); ok {
		if t, err := time.Parse("2006-01-02 15:04:05", updateTime); err == nil {
			result.WriteString(fmt.Sprintf("🕒 更新: %s\n", t.Format("01-02 15:04")))
		}
	}
	
	if weather, ok := data["weather"].(string); ok {
		result.WriteString(fmt.Sprintf("🌦️ 天气: %s\n", weather))
	}
	
	if minTemp, ok := data["min_temp"].(string); ok {
		if temp, ok := data["temp"].(string); ok {
			if maxTemp, ok := data["max_temp"].(string); ok {
				result.WriteString(fmt.Sprintf("🌡️ 温度: ↓%s℃| 现%s℃| ↑%s℃\n", minTemp, temp, maxTemp))
			}
		}
	}
	
	if wind, ok := data["wind"].(string); ok {
		result.WriteString(fmt.Sprintf("🌬️ 风向: %s\n", wind))
	}
	
	if humidity, ok := data["humidity"].(string); ok {
		result.WriteString(fmt.Sprintf("💦 湿度: %s\n", humidity))
	}
	
	if sunrise, ok := data["sunrise"].(string); ok {
		if sunset, ok := data["sunset"].(string); ok {
			result.WriteString(fmt.Sprintf("🌅 日出/日落: %s / %s\n", sunrise, sunset))
		}
	}
	
	// 运动指数
	if index, ok := data["index"].([]interface{}); ok && len(index) > 4 {
		if chuangyi, ok := index[4].(map[string]interface{}); ok {
			if level, ok := chuangyi["level"].(string); ok {
				if content, ok := chuangyi["content"].(string); ok {
					result.WriteString(fmt.Sprintf("👚 运动指数: %s - %s\n", level, content))
				}
			}
		}
	}
	
	// 未来10小时天气
	if hour, ok := data["hour"].([]interface{}); ok {
		var futureWeather []string
		now := time.Now()
		tenHoursLater := now.Add(10 * time.Hour)
		
		for _, hourData := range hour {
			if hourMap, ok := hourData.(map[string]interface{}); ok {
				if timeStr, ok := hourMap["time"].(string); ok {
					if forecastTime, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
						if now.Before(forecastTime) && forecastTime.Before(tenHoursLater) {
							if wea, ok := hourMap["wea"].(string); ok {
								if temp, ok := hourMap["temp"].(string); ok {
									futureWeather = append(futureWeather, fmt.Sprintf("     %02d:00 - %s - %s°C", forecastTime.Hour(), wea, temp))
								}
							}
						}
					}
				}
			}
		}
		
		if len(futureWeather) > 0 {
			result.WriteString("⏳ 未来10小时的天气预报:\n")
			result.WriteString(strings.Join(futureWeather, "\n"))
		}
	}
	
	// 预警信息
	if alarm, ok := data["alarm"].([]interface{}); ok && len(alarm) > 0 {
		result.WriteString("\n⚠️ 预警信息:\n")
		for _, alarmItem := range alarm {
			if alarmMap, ok := alarmItem.(map[string]interface{}); ok {
				if title, ok := alarmMap["title"].(string); ok {
					result.WriteString(fmt.Sprintf("🔴 标题: %s\n", title))
				}
				if level, ok := alarmMap["level"].(string); ok {
					result.WriteString(fmt.Sprintf("🟠 等级: %s\n", level))
				}
				if alarmType, ok := alarmMap["type"].(string); ok {
					result.WriteString(fmt.Sprintf("🟡 类型: %s\n", alarmType))
				}
				if tips, ok := alarmMap["tips"].(string); ok {
					result.WriteString(fmt.Sprintf("🟢 提示: \n%s\n", tips))
				}
				if content, ok := alarmMap["content"].(string); ok {
					result.WriteString(fmt.Sprintf("🔵 内容: \n%s\n\n", content))
				}
			}
		}
	}
	
	return result.String()
}

// formatFutureWeather 格式化未来天气
func (p *ApilotPlugin) formatFutureWeather(data []interface{}, date string) string {
	var result strings.Builder
	
	for num, d := range data {
		if dayMap, ok := d.(map[string]interface{}); ok {
			if num == 0 {
				if city, ok := dayMap["city"].(string); ok {
					if province, ok := dayMap["province"].(string); ok {
						result.WriteString(fmt.Sprintf("🏙️ 城市: %s (%s)\n", city, province))
					}
				}
			}
			
			// 根据日期过滤
			if date == "明天" && num != 1 {
				continue
			}
			if date == "后天" && num != 2 {
				continue
			}
			
			if dayDate, ok := dayMap["date"].(string); ok {
				result.WriteString(fmt.Sprintf("🕒 日期: %s\n", dayDate))
			}
			
			if weaDay, ok := dayMap["wea_day"].(string); ok {
				if weaNight, ok := dayMap["wea_night"].(string); ok {
					result.WriteString(fmt.Sprintf("🌦️ 天气: 🌞%s| 🌛%s\n", weaDay, weaNight))
				}
			}
			
			if tempDay, ok := dayMap["temp_day"].(string); ok {
				if tempNight, ok := dayMap["temp_night"].(string); ok {
					result.WriteString(fmt.Sprintf("🌡️ 温度: 🌞%s℃| 🌛%s℃\n", tempDay, tempNight))
				}
			}
			
			if sunrise, ok := dayMap["sunrise"].(string); ok {
				if sunset, ok := dayMap["sunset"].(string); ok {
					result.WriteString(fmt.Sprintf("🌅 日出/日落: %s / %s\n", sunrise, sunset))
				}
			}
			
			if index, ok := dayMap["index"].([]interface{}); ok {
				for _, idx := range index {
					if idxMap, ok := idx.(map[string]interface{}); ok {
						if name, ok := idxMap["name"].(string); ok {
							if level, ok := idxMap["level"].(string); ok {
								result.WriteString(fmt.Sprintf("%s: %s\n", name, level))
							}
						}
					}
				}
			}
			
			result.WriteString("\n")
		}
	}
	
	return result.String()
}