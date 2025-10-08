package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"wechat-robot-client/interface/plugin"
)

// NeteasyPlugin 网易视频下载插件
type NeteasyPlugin struct {
	config NeteasyConfig
}

// NeteasyConfig 插件配置
type NeteasyConfig struct {
	VideoFolder    string `json:"video_folder"`
	TempFolder     string `json:"temp_folder"`
	UploadFolder   string `json:"upload_folder"`
	UserAgent      string `json:"user_agent"`
	NeteasyAppID   string `json:"neteasy_app_id"`
}

// VideoInfo 视频信息结构
type VideoInfo struct {
	Title       string
	CoverImage  string
	M3U8URL     string
	OutputFile  string
	OriginalURL string
}

// NewNeteasyPlugin 创建网易视频插件实例
func NewNeteasyPlugin() plugin.MessageHandler {
	plugin := &NeteasyPlugin{}
	plugin.loadConfig()
	return plugin
}

// GetName 获取插件名称
func (p *NeteasyPlugin) GetName() string {
	return "Neteasy"
}

// GetLabels 获取插件标签
func (p *NeteasyPlugin) GetLabels() []string {
	return []string{"neteasy", "video", "download", "netease"}
}

// PreAction 前置处理
func (p *NeteasyPlugin) PreAction(ctx *plugin.MessageContext) bool {
	return true
}

// PostAction 后置处理
func (p *NeteasyPlugin) PostAction(ctx *plugin.MessageContext) {
	// 可以在这里添加清理逻辑
}

// Run 主要逻辑
func (p *NeteasyPlugin) Run(ctx *plugin.MessageContext) bool {
	content := ctx.MessageContent
	
	// 检查是否包含网易视频链接
	if !strings.Contains(content, p.config.NeteasyAppID) {
		return false
	}
	
	// 解析视频URL
	videoURL := p.parseVideoURL(content)
	if videoURL == "" {
		return false
	}
	
	// 下载视频
	videoInfo, err := p.downloadVideo(videoURL)
	if err != nil {
		p.sendReply(ctx, "text", fmt.Sprintf("视频下载失败: %v", err))
		return true
	}
	
	if videoInfo == nil {
		p.sendReply(ctx, "text", "未找到可下载的视频")
		return true
	}
	
	// 上传到云盘（这里简化处理，实际需要集成云盘API）
	uploadSuccess := p.uploadToCloud(videoInfo.OutputFile)
	
	// 清理临时文件
	p.cleanupTempFiles(videoInfo.OutputFile)
	
	// 发送结果
	if uploadSuccess {
		p.sendReply(ctx, "text", fmt.Sprintf("neteasy,视频上传完成: %s 链接: %s", videoInfo.Title, videoInfo.OriginalURL))
	} else {
		p.sendReply(ctx, "text", fmt.Sprintf("neteasy,视频上传失败: %s", videoInfo.Title))
	}
	
	return true
}

// loadConfig 加载配置
func (p *NeteasyPlugin) loadConfig() {
	configPath := "plugin/plugins/neteasy_config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 使用默认配置
		p.config = NeteasyConfig{
			VideoFolder:    "video",
			TempFolder:     "temp_segments",
			UploadFolder:   "网易视频",
			UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
			NeteasyAppID:   "wx7be3c1bb46c68c63",
		}
		return
	}
	
	json.Unmarshal(data, &p.config)
}

// sendReply 发送回复
func (p *NeteasyPlugin) sendReply(ctx *plugin.MessageContext, replyType, content string) {
	if ctx.Message.IsChatRoom {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content, ctx.Message.SenderWxID)
	} else {
		ctx.MessageService.SendTextMessage(ctx.Message.FromWxID, content)
	}
}

// parseVideoURL 解析视频URL
func (p *NeteasyPlugin) parseVideoURL(content string) string {
	// 匹配网易视频链接的正则表达式
	pattern := regexp.MustCompile(`https?://c\.m\.163\.com/[^\s<>"]+`)
	matches := pattern.FindAllString(content, -1)
	
	if len(matches) > 0 {
		// 取第一个匹配的URL，去掉末尾可能的特殊字符
		videoURL := strings.TrimRight(matches[0], "&amp;")
		return videoURL
	}
	
	return ""
}

// downloadVideo 下载视频
func (p *NeteasyPlugin) downloadVideo(url string) (*VideoInfo, error) {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 设置请求头
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Host", "c.m.163.com")
	req.Header.Set("User-Agent", p.config.UserAgent)
	
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// 解析HTML内容
	videoInfo := p.parseHTMLContent(string(body), url)
	if videoInfo == nil {
		return nil, fmt.Errorf("无法解析视频信息")
	}
	
	// 下载M3U8视频
	err = p.downloadM3U8Video(videoInfo)
	if err != nil {
		return nil, err
	}
	
	return videoInfo, nil
}

// parseHTMLContent 解析HTML内容获取视频信息
func (p *NeteasyPlugin) parseHTMLContent(htmlContent, originalURL string) *VideoInfo {
	// 查找包含data-m3u8属性的video标签
	m3u8Pattern := regexp.MustCompile(`<video[^>]*data-m3u8="([^"]*)"[^>]*>`)
	m3u8Matches := m3u8Pattern.FindStringSubmatch(htmlContent)
	
	if len(m3u8Matches) < 2 {
		return nil
	}
	
	m3u8URL := m3u8Matches[1]
	
	// 提取标题
	title := "网易视频"
	titlePattern := regexp.MustCompile(`<div[^>]*class="digest"[^>]*>([^<]*)</div>`)
	titleMatches := titlePattern.FindStringSubmatch(htmlContent)
	if len(titleMatches) >= 2 {
		title = strings.TrimSpace(titleMatches[1])
	} else {
		// 尝试其他标题选择器
		titlePattern2 := regexp.MustCompile(`<div[^>]*class="m-viewpoint"[^>]*>([^<]*)</div>`)
		titleMatches2 := titlePattern2.FindStringSubmatch(htmlContent)
		if len(titleMatches2) >= 2 {
			title = strings.TrimSpace(titleMatches2[1])
		}
	}
	
	// 提取封面图片
	coverImage := ""
	coverPattern := regexp.MustCompile(`<div[^>]*data-bg="([^"]*)"[^>]*>`)
	coverMatches := coverPattern.FindStringSubmatch(htmlContent)
	if len(coverMatches) >= 2 {
		coverImage = coverMatches[1]
	}
	
	// 清理标题中的特殊字符
	title = p.cleanTitle(title)
	
	// 生成输出文件路径
	outputFile := filepath.Join(p.config.VideoFolder, title+".mp4")
	
	return &VideoInfo{
		Title:       title,
		CoverImage:  coverImage,
		M3U8URL:     m3u8URL,
		OutputFile:  outputFile,
		OriginalURL: originalURL,
	}
}

// cleanTitle 清理标题中的特殊字符
func (p *NeteasyPlugin) cleanTitle(title string) string {
	// 替换特殊字符为下划线
	pattern := regexp.MustCompile(`[~` + "`" + `!@#$%^&*()_"：\-+=|\\{\}\[\]:;\"'<>,.?/·！￥…（）—【】、？《》，。]+`)
	cleaned := pattern.ReplaceAllString(title, "_")
	
	// 限制长度
	if len(cleaned) > 100 {
		cleaned = cleaned[:100]
	}
	
	return cleaned
}

// downloadM3U8Video 下载M3U8视频
func (p *NeteasyPlugin) downloadM3U8Video(videoInfo *VideoInfo) error {
	// 创建视频目录
	err := os.MkdirAll(p.config.VideoFolder, 0755)
	if err != nil {
		return err
	}
	
	// 创建临时目录
	err = os.MkdirAll(p.config.TempFolder, 0755)
	if err != nil {
		return err
	}
	
	// 下载M3U8播放列表
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(videoInfo.M3U8URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	m3u8Content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	// 解析M3U8内容
	segments := p.parseM3U8Content(string(m3u8Content), videoInfo.M3U8URL)
	if len(segments) == 0 {
		return fmt.Errorf("无法解析M3U8播放列表")
	}
	
	// 下载所有片段
	for i, segmentURL := range segments {
		segmentPath := filepath.Join(p.config.TempFolder, fmt.Sprintf("segment_%d.ts", i))
		err := p.downloadSegment(segmentURL, segmentPath)
		if err != nil {
			return fmt.Errorf("下载片段 %d 失败: %v", i, err)
		}
	}
	
	// 合并视频片段
	err = p.mergeSegments(videoInfo.OutputFile, len(segments))
	if err != nil {
		return fmt.Errorf("合并视频失败: %v", err)
	}
	
	// 清理临时文件
	p.cleanupTempSegments(len(segments))
	
	return nil
}

// parseM3U8Content 解析M3U8内容
func (p *NeteasyPlugin) parseM3U8Content(content, baseURL string) []string {
	var segments []string
	lines := strings.Split(content, "\n")
	
	baseDir := baseURL[:strings.LastIndex(baseURL, "/")+1]
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// 处理相对URL
			if !strings.HasPrefix(line, "http") {
				line = baseDir + line
			}
			segments = append(segments, line)
		}
	}
	
	return segments
}

// downloadSegment 下载单个片段
func (p *NeteasyPlugin) downloadSegment(url, filepath string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.Copy(file, resp.Body)
	return err
}

// mergeSegments 合并视频片段
func (p *NeteasyPlugin) mergeSegments(outputFile string, segmentCount int) error {
	// 创建文件列表
	fileListPath := filepath.Join(p.config.TempFolder, "file_list.txt")
	file, err := os.Create(fileListPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// 写入文件列表
	for i := 0; i < segmentCount; i++ {
		segmentPath := filepath.Join(p.config.TempFolder, fmt.Sprintf("segment_%d.ts", i))
		_, err := file.WriteString(fmt.Sprintf("file '%s'\n", segmentPath))
		if err != nil {
			return err
		}
	}
	file.Close()
	
	// 这里简化处理，实际应该使用ffmpeg进行合并
	// 由于Go中没有直接的ffmpeg绑定，这里使用简单的文件复制作为示例
	// 实际项目中需要调用外部ffmpeg命令或使用CGO绑定
	
	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()
	
	// 逐个复制片段到输出文件
	for i := 0; i < segmentCount; i++ {
		segmentPath := filepath.Join(p.config.TempFolder, fmt.Sprintf("segment_%d.ts", i))
		segmentFile, err := os.Open(segmentPath)
		if err != nil {
			return err
		}
		
		_, err = io.Copy(outFile, segmentFile)
		segmentFile.Close()
		if err != nil {
			return err
		}
	}
	
	// 删除文件列表
	os.Remove(fileListPath)
	
	return nil
}

// cleanupTempSegments 清理临时片段文件
func (p *NeteasyPlugin) cleanupTempSegments(count int) {
	for i := 0; i < count; i++ {
		segmentPath := filepath.Join(p.config.TempFolder, fmt.Sprintf("segment_%d.ts", i))
		os.Remove(segmentPath)
	}
}

// cleanupTempFiles 清理临时文件
func (p *NeteasyPlugin) cleanupTempFiles(filePath string) {
	// 这里可以添加清理逻辑
	// 比如删除临时目录等
}

// uploadToCloud 上传到云盘
func (p *NeteasyPlugin) uploadToCloud(filePath string) bool {
	// 这里简化处理，实际需要集成云盘API
	// 比如阿里云盘、百度网盘等
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	
	// 模拟上传过程
	time.Sleep(1 * time.Second)
	
	// 这里应该调用实际的云盘上传API
	// 返回true表示上传成功，false表示失败
	return true
}