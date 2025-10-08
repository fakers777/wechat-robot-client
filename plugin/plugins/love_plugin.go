package plugins

import (
	"math/rand"
	"strings"
	"time"
	"wechat-robot-client/interface/plugin"
)

// LovePlugin 土味情话插件
type LovePlugin struct{}

// NewLovePlugin 创建土味情话插件实例
func NewLovePlugin() plugin.MessageHandler {
	return &LovePlugin{}
}

// GetName 获取插件名称
func (p *LovePlugin) GetName() string {
	return "Love"
}

// GetLabels 获取插件标签
func (p *LovePlugin) GetLabels() []string {
	return []string{"love", "romance", "fun"}
}

// PreAction 前置处理
func (p *LovePlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *LovePlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *LovePlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含情话相关关键词
	if !p.containsLoveKeywords(messageContent) {
		return false
	}
	
	// 生成土味情话
	loveText := p.generateLoveText()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, loveText, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, loveText)
	}
	
	return true
}

// containsLoveKeywords 检查是否包含情话相关关键词
func (p *LovePlugin) containsLoveKeywords(content string) bool {
	loveKeywords := []string{
		"情话", "土味", "表白", "喜欢", "爱",
		"甜言蜜语", "浪漫", "撩", "撩人",
		"想你", "爱你", "宝贝", "亲爱的",
		"情话", "土味情话", "说情话",
	}
	
	for _, keyword := range loveKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// generateLoveText 生成土味情话
func (p *LovePlugin) generateLoveText() string {
	rand.Seed(time.Now().UnixNano())
	
	// 土味情话模板
	loveTexts := []string{
		"你知道我的缺点是什么吗？是缺点你。💕",
		"你知道我想喝什么吗？我想呵护你。🥤",
		"你知道我想吃什么吗？我想痴痴地望着你。😍",
		"你知道我想看什么吗？我想看你的心。👀",
		"你知道我想听什么吗？我想听你的声音。👂",
		"你知道我想闻什么吗？我想闻你的味道。👃",
		"你知道我想摸什么吗？我想摸你的手。✋",
		"你知道我想抱什么吗？我想抱抱你。🤗",
		"你知道我想亲什么吗？我想亲亲你。😘",
		"你知道我想做什么吗？我想做你的男朋友。💑",
		"你知道我想成为什么吗？我想成为你的唯一。💎",
		"你知道我想拥有什么吗？我想拥有你的心。❤️",
		"你知道我想守护什么吗？我想守护你的笑容。😊",
		"你知道我想陪伴什么吗？我想陪伴你一生。👫",
		"你知道我想给你什么吗？我想给你我的全部。🎁",
		"你知道我想对你说什么吗？我想对你说我爱你。💕",
		"你知道我想和你做什么吗？我想和你一起变老。👴👵",
		"你知道我想去哪里吗？我想去你的心里。🏠",
		"你知道我想学什么吗？我想学如何爱你。📚",
		"你知道我想唱什么吗？我想唱情歌给你听。🎵",
		"你知道我想写什么吗？我想写情书给你。📝",
		"你知道我想画什么吗？我想画你的样子。🎨",
		"你知道我想拍什么吗？我想拍下你的美。📸",
		"你知道我想买什么吗？我想买下你的心。💳",
		"你知道我想送什么吗？我想送给你我的爱。💝",
		"你知道我想等什么吗？我想等你爱上我。⏰",
		"你知道我想追什么吗？我想追到你。🏃‍♂️",
		"你知道我想赢什么吗？我想赢得你的心。🏆",
		"你知道我想输什么吗？我想输给你。💔",
		"你知道我想赢什么吗？我想赢得你的爱。💖",
	}
	
	// 随机选择一个情话
	loveText := loveTexts[rand.Intn(len(loveTexts))]
	
	return loveText
}

// LoveStoryPlugin 爱情故事插件
type LoveStoryPlugin struct{}

// NewLoveStoryPlugin 创建爱情故事插件实例
func NewLoveStoryPlugin() plugin.MessageHandler {
	return &LoveStoryPlugin{}
}

// GetName 获取插件名称
func (p *LoveStoryPlugin) GetName() string {
	return "LoveStory"
}

// GetLabels 获取插件标签
func (p *LoveStoryPlugin) GetLabels() []string {
	return []string{"love", "story", "romance"}
}

// PreAction 前置处理
func (p *LoveStoryPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *LoveStoryPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *LoveStoryPlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含故事相关关键词
	if !p.containsStoryKeywords(messageContent) {
		return false
	}
	
	// 生成爱情故事
	story := p.generateLoveStory()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, story, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, story)
	}
	
	return true
}

// containsStoryKeywords 检查是否包含故事相关关键词
func (p *LoveStoryPlugin) containsStoryKeywords(content string) bool {
	storyKeywords := []string{
		"故事", "讲个", "说个", "听故事", "讲故事",
		"爱情", "恋爱", "情侣", "夫妻", "恋人",
		"浪漫", "甜蜜", "幸福", "美好", "温馨",
	}
	
	for _, keyword := range storyKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// generateLoveStory 生成爱情故事
func (p *LoveStoryPlugin) generateLoveStory() string {
	rand.Seed(time.Now().UnixNano())
	
	// 爱情故事模板
	stories := []string{
		"从前有一个男孩，他每天都会在同一个咖啡厅里看到一个女孩。女孩总是坐在靠窗的位置，安静地看书。男孩鼓起勇气走过去，对女孩说：'你知道我想喝什么吗？我想呵护你。'女孩笑了，从此他们开始了美好的爱情故事。💕",
		"有一个女孩，她每天都会收到一束花，但不知道是谁送的。直到有一天，她发现送花的人是她暗恋已久的男孩。男孩对她说：'你知道我想送什么吗？我想送给你我的爱。'女孩感动得哭了，他们从此幸福地在一起。🌹",
		"有一个男孩，他为了追求心爱的女孩，每天都会在她家楼下等她。无论刮风下雨，他都会准时出现。女孩终于被他的坚持感动了，对他说：'你知道我想等什么吗？我想等你爱上我。'他们从此开始了甜蜜的恋爱。⏰",
		"有一个女孩，她总是觉得自己不够漂亮，不够优秀。直到有一天，一个男孩对她说：'你知道我想守护什么吗？我想守护你的笑容。'女孩终于明白，真正的爱情不是看外表，而是看内心。他们从此幸福地在一起。😊",
		"有一个男孩，他为了给心爱的女孩一个惊喜，学会了做她最爱吃的蛋糕。当女孩看到蛋糕上的字'你知道我想给你什么吗？我想给你我的全部'时，她感动得哭了。他们从此开始了甜蜜的生活。🎂",
		"有一个女孩，她总是觉得自己配不上优秀的男孩。直到有一天，男孩对她说：'你知道我想成为什么吗？我想成为你的唯一。'女孩终于明白，爱情没有配不配，只有爱不爱。他们从此幸福地在一起。💎",
		"有一个男孩，他为了追求心爱的女孩，每天都会写一首情诗送给她。女孩被他的才华和真心感动了，对他说：'你知道我想写什么吗？我想写情书给你。'他们从此开始了浪漫的爱情。📝",
		"有一个女孩，她总是觉得自己不够浪漫。直到有一天，一个男孩对她说：'你知道我想唱什么吗？我想唱情歌给你听。'女孩终于明白，浪漫不是形式，而是真心。他们从此幸福地在一起。🎵",
		"有一个男孩，他为了给心爱的女孩一个完美的求婚，准备了很久。当女孩看到他的真心时，对他说：'你知道我想做什么吗？我想做你的女朋友。'他们从此开始了美好的婚姻。💍",
		"有一个女孩，她总是觉得自己不够温柔。直到有一天，一个男孩对她说：'你知道我想抱什么吗？我想抱抱你。'女孩终于明白，温柔不是性格，而是爱情。他们从此幸福地在一起。🤗",
	}
	
	// 随机选择一个故事
	story := stories[rand.Intn(len(stories))]
	
	return story
}

// LoveAdvicePlugin 爱情建议插件
type LoveAdvicePlugin struct{}

// NewLoveAdvicePlugin 创建爱情建议插件实例
func NewLoveAdvicePlugin() plugin.MessageHandler {
	return &LoveAdvicePlugin{}
}

// GetName 获取插件名称
func (p *LoveAdvicePlugin) GetName() string {
	return "LoveAdvice"
}

// GetLabels 获取插件标签
func (p *LoveAdvicePlugin) GetLabels() []string {
	return []string{"love", "advice", "help"}
}

// PreAction 前置处理
func (p *LoveAdvicePlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *LoveAdvicePlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *LoveAdvicePlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含建议相关关键词
	if !p.containsAdviceKeywords(messageContent) {
		return false
	}
	
	// 生成爱情建议
	advice := p.generateLoveAdvice()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, advice, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, advice)
	}
	
	return true
}

// containsAdviceKeywords 检查是否包含建议相关关键词
func (p *LoveAdvicePlugin) containsAdviceKeywords(content string) bool {
	adviceKeywords := []string{
		"建议", "帮助", "怎么办", "如何", "怎么",
		"表白", "追求", "恋爱", "分手", "复合",
		"爱情", "感情", "关系", "相处", "沟通",
		"问题", "困扰", "烦恼", "纠结", "迷茫",
	}
	
	for _, keyword := range adviceKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// generateLoveAdvice 生成爱情建议
func (p *LoveAdvicePlugin) generateLoveAdvice() string {
	rand.Seed(time.Now().UnixNano())
	
	// 爱情建议模板
	advices := []string{
		"💕 爱情建议 💕\n\n1. 真诚是爱情的基础，不要伪装自己\n2. 沟通是解决问题的关键，有话直说\n3. 尊重对方的想法和选择\n4. 给彼此一些空间和时间\n5. 学会包容和理解\n6. 保持新鲜感，偶尔制造惊喜\n7. 共同成长，一起进步\n8. 珍惜当下，不要轻易放弃\n\n记住：真正的爱情是相互的，不是单方面的付出。💖",
		"💕 表白建议 💕\n\n1. 选择合适的时间和地点\n2. 准备真诚的话语，不要过于华丽\n3. 了解对方的喜好和性格\n4. 不要给太大压力，给对方考虑的时间\n5. 如果被拒绝，要尊重对方的选择\n6. 保持友谊，不要因爱生恨\n7. 提升自己，让自己变得更好\n8. 相信缘分，不要强求\n\n记住：表白不是终点，而是新的开始。💖",
		"💕 恋爱建议 💕\n\n1. 保持独立，不要完全依赖对方\n2. 学会倾听，理解对方的感受\n3. 保持神秘感，不要完全透明\n4. 共同规划未来，有共同目标\n5. 学会妥协，但不要失去原则\n6. 保持浪漫，偶尔制造惊喜\n7. 信任对方，不要无端猜疑\n8. 珍惜感情，不要轻易说分手\n\n记住：恋爱是两个人的事，需要共同努力。💖",
		"💕 相处建议 💕\n\n1. 尊重对方的隐私和空间\n2. 学会换位思考，理解对方\n3. 保持幽默感，让生活有趣\n4. 学会道歉，承认错误\n5. 保持耐心，不要急躁\n6. 学会感恩，珍惜对方的好\n7. 保持健康的生活方式\n8. 共同面对困难，一起成长\n\n记住：相处是一门艺术，需要用心经营。💖",
		"💕 沟通建议 💕\n\n1. 选择合适的时间进行沟通\n2. 用'我'而不是'你'来表达感受\n3. 倾听对方的想法，不要打断\n4. 避免指责和批评，用建设性的语言\n5. 保持冷静，不要情绪化\n6. 寻求共同点，而不是分歧\n7. 学会妥协，找到平衡点\n8. 定期沟通，保持联系\n\n记住：沟通是爱情的桥梁，需要用心搭建。💖",
	}
	
	// 随机选择一个建议
	advice := advices[rand.Intn(len(advices))]
	
	return advice
}

// LoveTestPlugin 爱情测试插件
type LoveTestPlugin struct{}

// NewLoveTestPlugin 创建爱情测试插件实例
func NewLoveTestPlugin() plugin.MessageHandler {
	return &LoveTestPlugin{}
}

// GetName 获取插件名称
func (p *LoveTestPlugin) GetName() string {
	return "LoveTest"
}

// GetLabels 获取插件标签
func (p *LoveTestPlugin) GetLabels() []string {
	return []string{"love", "test", "fun"}
}

// PreAction 前置处理
func (p *LoveTestPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *LoveTestPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *LoveTestPlugin) Run(ctx *plugin.MessageContext) bool {
	messageContent := strings.ToLower(ctx.MessageContent)
	
	// 检查是否包含测试相关关键词
	if !p.containsTestKeywords(messageContent) {
		return false
	}
	
	// 生成爱情测试
	test := p.generateLoveTest()
	
	// 发送消息
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, test, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, test)
	}
	
	return true
}

// containsTestKeywords 检查是否包含测试相关关键词
func (p *LoveTestPlugin) containsTestKeywords(content string) bool {
	testKeywords := []string{
		"测试", "测验", "测", "算", "算算",
		"爱情", "恋爱", "感情", "缘分", "匹配",
		"星座", "血型", "性格", "配对", "合适",
	}
	
	for _, keyword := range testKeywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	
	return false
}

// generateLoveTest 生成爱情测试
func (p *LoveTestPlugin) generateLoveTest() string {
	rand.Seed(time.Now().UnixNano())
	
	// 爱情测试模板
	tests := []string{
		"💕 爱情测试 💕\n\n问题：你最喜欢什么颜色？\n\nA. 红色 - 热情如火，爱情指数：95%\nB. 蓝色 - 深沉内敛，爱情指数：85%\nC. 绿色 - 自然清新，爱情指数：80%\nD. 粉色 - 温柔浪漫，爱情指数：90%\n\n选择你的答案，看看你的爱情指数吧！💖",
		"💕 爱情测试 💕\n\n问题：你最喜欢什么季节？\n\nA. 春天 - 充满希望，爱情指数：88%\nB. 夏天 - 热情奔放，爱情指数：92%\nC. 秋天 - 成熟稳重，爱情指数：85%\nD. 冬天 - 冷静理智，爱情指数：78%\n\n选择你的答案，看看你的爱情指数吧！💖",
		"💕 爱情测试 💕\n\n问题：你最喜欢什么动物？\n\nA. 猫 - 独立优雅，爱情指数：82%\nB. 狗 - 忠诚热情，爱情指数：95%\nC. 兔子 - 温柔可爱，爱情指数：88%\nD. 鸟 - 自由浪漫，爱情指数：85%\n\n选择你的答案，看看你的爱情指数吧！💖",
		"💕 爱情测试 💕\n\n问题：你最喜欢什么花？\n\nA. 玫瑰 - 浪漫热情，爱情指数：95%\nB. 百合 - 纯洁美好，爱情指数：88%\nC. 向日葵 - 阳光积极，爱情指数：90%\nD. 薰衣草 - 神秘浪漫，爱情指数：85%\n\n选择你的答案，看看你的爱情指数吧！💖",
		"💕 爱情测试 💕\n\n问题：你最喜欢什么食物？\n\nA. 甜食 - 甜蜜浪漫，爱情指数：90%\nB. 辣食 - 热情如火，爱情指数：92%\nC. 清淡 - 温和体贴，爱情指数：85%\nD. 酸食 - 个性独特，爱情指数：80%\n\n选择你的答案，看看你的爱情指数吧！💖",
	}
	
	// 随机选择一个测试
	test := tests[rand.Intn(len(tests))]
	
	return test
}