package plugins

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"wechat-robot-client/interface/plugin"
)

// KFCPlugin KFC疯狂星期四插件
type KFCPlugin struct{}

// NewKFCPlugin 创建KFC插件实例
func NewKFCPlugin() plugin.MessageHandler {
	return &KFCPlugin{}
}

// GetName 获取插件名称
func (p *KFCPlugin) GetName() string {
	return "KFC"
}

// GetLabels 获取插件标签
func (p *KFCPlugin) GetLabels() []string {
	return []string{"kfc", "fun", "thursday"}
}

// PreAction 前置处理
func (p *KFCPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *KFCPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *KFCPlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含KFC相关关键词
	if !p.containsKFCKeywords(messageContent) {
		return false
	}
	
	// 检查是否是星期四
	if !p.isThursday() {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, "今天不是星期四，等星期四再来吧！😄")
		return true
	}
	
	// 生成KFC疯狂星期四文案
	kfcText := p.generateKFCText()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, kfcText, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, kfcText)
	}
	
	return true
}

// containsKFCKeywords 检查是否包含KFC相关关键词
func (p *KFCPlugin) containsKFCKeywords(content string) bool {
	kfcKeywords := []string{
		"kfc", "肯德基", "疯狂星期四", "星期四", "周四",
		"炸鸡", "原味鸡", "汉堡", "薯条", "可乐",
		"疯狂", "v我50", "v我", "请我", "请客",
	}
	
	for _, keyword := range kfcKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	// 检查是否包含数字+元/块等金额相关词汇
	moneyPattern := regexp.MustCompile(`\d+[元块]`)
	if moneyPattern.MatchString(content) {
		return true
	}
	
	return false
}

// isThursday 检查今天是否是星期四
func (p *KFCPlugin) isThursday() bool {
	now := time.Now()
	return now.Weekday() == time.Thursday
}

// generateKFCText 生成KFC疯狂星期四文案
func (p *KFCPlugin) generateKFCText() string {
	rand.Seed(time.Now().UnixNano())
	
	// KFC疯狂星期四文案模板
	templates := []string{
		"今天是肯德基疯狂星期四！谁能v我50，我想吃原味鸡😋",
		"今天是星期四，肯德基疯狂星期四！谁请我吃炸鸡？🍗",
		"今天是肯德基疯狂星期四！v我50，我请你吃汉堡🍔",
		"今天是星期四，肯德基疯狂星期四！谁v我50，我想吃薯条🍟",
		"今天是肯德基疯狂星期四！谁能v我50，我想喝可乐🥤",
		"今天是星期四，肯德基疯狂星期四！谁请我吃原味鸡？🍗",
		"今天是肯德基疯狂星期四！v我50，我请你吃全家桶🍗",
		"今天是星期四，肯德基疯狂星期四！谁能v我50，我想吃鸡翅🍗",
		"今天是肯德基疯狂星期四！谁请我吃炸鸡？v我50😋",
		"今天是星期四，肯德基疯狂星期四！v我50，我请你吃汉堡🍔",
	}
	
	// 随机选择一个模板
	template := templates[rand.Intn(len(templates))]
	
	// 添加一些随机的表情符号
	emojis := []string{"😋", "🍗", "🍔", "🍟", "🥤", "😄", "🤤", "😍"}
	emoji := emojis[rand.Intn(len(emojis))]
	
	return fmt.Sprintf("%s %s", template, emoji)
}

// KFCStoryPlugin KFC故事插件
type KFCStoryPlugin struct{}

// NewKFCStoryPlugin 创建KFC故事插件实例
func NewKFCStoryPlugin() plugin.MessageHandler {
	return &KFCStoryPlugin{}
}

// GetName 获取插件名称
func (p *KFCStoryPlugin) GetName() string {
	return "KFCStory"
}

// GetLabels 获取插件标签
func (p *KFCStoryPlugin) GetLabels() []string {
	return []string{"kfc", "story", "fun"}
}

// PreAction 前置处理
func (p *KFCStoryPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *KFCStoryPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *KFCStoryPlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含故事相关关键词
	if !p.containsStoryKeywords(messageContent) {
		return false
	}
	
	// 检查是否是星期四
	if !p.isThursday() {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, "今天不是星期四，等星期四再来听故事吧！😄")
		return true
	}
	
	// 生成KFC故事
	story := p.generateKFCStory()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, story, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, story)
	}
	
	return true
}

// containsStoryKeywords 检查是否包含故事相关关键词
func (p *KFCStoryPlugin) containsStoryKeywords(content string) bool {
	storyKeywords := []string{
		"故事", "讲个", "说个", "听故事", "讲故事",
		"今天", "昨天", "前天", "发生", "事情",
		"我", "朋友", "同学", "同事", "家人",
	}
	
	for _, keyword := range storyKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// isThursday 检查今天是否是星期四
func (p *KFCStoryPlugin) isThursday() bool {
	now := time.Now()
	return now.Weekday() == time.Thursday
}

// generateKFCStory 生成KFC故事
func (p *KFCStoryPlugin) generateKFCStory() string {
	rand.Seed(time.Now().UnixNano())
	
	// KFC故事模板
	stories := []string{
		"今天我在路上看到一个老人摔倒了，我赶紧跑过去扶他起来。老人很感激，说要请我吃饭。我说不用了，老人坚持要请，最后他说：今天是肯德基疯狂星期四，谁能v我50？😅",
		"昨天我朋友跟我说他失恋了，很伤心。我安慰他说没关系，时间会治愈一切。他说：是的，时间会治愈一切，但今天是肯德基疯狂星期四，谁能v我50？😭",
		"今天上班的时候，老板突然把我叫到办公室，我以为要升职加薪了。结果老板说：今天是肯德基疯狂星期四，谁能v我50？😱",
		"今天我在医院排队挂号，前面有个小朋友一直哭。我问他怎么了，他说：今天是肯德基疯狂星期四，谁能v我50？😢",
		"今天我在银行取钱，ATM机突然坏了。银行工作人员说：今天是肯德基疯狂星期四，谁能v我50？🏦",
		"今天我在超市买菜，收银员突然说：今天是肯德基疯狂星期四，谁能v我50？🛒",
		"今天我在公交车上，司机突然停车说：今天是肯德基疯狂星期四，谁能v我50？🚌",
		"今天我在图书馆看书，管理员突然走过来小声说：今天是肯德基疯狂星期四，谁能v我50？📚",
		"今天我在健身房锻炼，教练突然停下来说：今天是肯德基疯狂星期四，谁能v我50？💪",
		"今天我在咖啡厅喝咖啡，服务员突然走过来神秘地说：今天是肯德基疯狂星期四，谁能v我50？☕",
	}
	
	// 随机选择一个故事
	story := stories[rand.Intn(len(stories))]
	
	// 添加一些随机的表情符号
	emojis := []string{"😅", "😭", "😱", "😢", "🏦", "🛒", "🚌", "📚", "💪", "☕"}
	emoji := emojis[rand.Intn(len(emojis))]
	
	return fmt.Sprintf("%s %s", story, emoji)
}

// KFCMenuPlugin KFC菜单插件
type KFCMenuPlugin struct{}

// NewKFCMenuPlugin 创建KFC菜单插件实例
func NewKFCMenuPlugin() plugin.MessageHandler {
	return &KFCMenuPlugin{}
}

// GetName 获取插件名称
func (p *KFCMenuPlugin) GetName() string {
	return "KFCMenu"
}

// GetLabels 获取插件标签
func (p *KFCMenuPlugin) GetLabels() []string {
	return []string{"kfc", "menu", "food"}
}

// PreAction 前置处理
func (p *KFCMenuPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *KFCMenuPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *KFCMenuPlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含菜单相关关键词
	if !p.containsMenuKeywords(messageContent) {
		return false
	}
	
	// 生成KFC菜单
	menu := p.generateKFCMenu()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, menu, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, menu)
	}
	
	return true
}

// containsMenuKeywords 检查是否包含菜单相关关键词
func (p *KFCMenuPlugin) containsMenuKeywords(content string) bool {
	menuKeywords := []string{
		"菜单", "有什么", "推荐", "好吃", "招牌",
		"kfc", "肯德基", "炸鸡", "汉堡", "薯条",
		"原味鸡", "鸡翅", "鸡腿", "全家桶",
	}
	
	for _, keyword := range menuKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// generateKFCMenu 生成KFC菜单
func (p *KFCMenuPlugin) generateKFCMenu() string {
	menu := `🍗 肯德基推荐菜单 🍗

🥤 饮品系列：
• 百事可乐 - 经典汽水
• 雪顶咖啡 - 香浓咖啡
• 柠檬红茶 - 清爽解腻

🍔 主食系列：
• 原味鸡 - 招牌经典
• 香辣鸡腿堡 - 香辣美味
• 新奥尔良烤鸡腿堡 - 嫩滑多汁

🍟 小食系列：
• 薯条 - 金黄酥脆
• 鸡米花 - 香脆可口
• 蛋挞 - 香甜嫩滑

🍗 全家桶：
• 原味鸡全家桶 - 全家共享
• 香辣鸡全家桶 - 香辣过瘾

今天是疯狂星期四，部分商品有优惠哦！😋`
	
	return menu
}