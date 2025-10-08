package controller

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"wechat-robot-client/model"
	"wechat-robot-client/pkg/robot"
)

// MessageDisplayInfo 消息显示信息
type MessageDisplayInfo struct {
	Type        string `json:"type"`
	TypeDesc    string `json:"type_desc"`
	Content     string `json:"content"`
	IsTruncated bool   `json:"is_truncated"`
	Length      int    `json:"length"`
}

// AppMsgXML 用于解析APP消息的XML结构
type AppMsgXML struct {
	XMLName xml.Name `xml:"msg"`
	AppMsg  struct {
		Type        string `xml:"type,attr"`
		Title       string `xml:"title"`
		Des         string `xml:"des"`
		URL         string `xml:"url"`
		AppName     string `xml:"appname"`
		UserName    string `xml:"username"`
		DisplayName string `xml:"displayname"`
		Content     string `xml:"content"`
	} `xml:"appmsg"`
}

// SysMsgXML 用于解析系统消息的XML结构
type SysMsgXML struct {
	XMLName xml.Name `xml:"sysmsg"`
	Type    string   `xml:"type,attr"`
	Content string   `xml:",chardata"`
}

// FormatMessageForDisplay 格式化消息用于显示
func FormatMessageForDisplay(msg robot.Message) MessageDisplayInfo {
	msgType := model.MessageType(msg.MsgType)
	content := ""
	if msg.Content.String != nil {
		content = *msg.Content.String
	}

	displayInfo := MessageDisplayInfo{
		Type:    fmt.Sprintf("%d", msgType),
		Content: content,
		Length:  len(content),
	}

	// 根据消息类型获取描述和处理内容
	switch msgType {
	case model.MsgTypeText:
		displayInfo.TypeDesc = "文本消息"
		displayInfo.Content = truncateText(content, 200)

	case model.MsgTypeImage:
		displayInfo.TypeDesc = "图片消息"
		displayInfo.Content = "[图片消息]"

	case model.MsgTypeVoice:
		displayInfo.TypeDesc = "语音消息"
		displayInfo.Content = "[语音消息]"

	case model.MsgTypeVideo:
		displayInfo.TypeDesc = "视频消息"
		displayInfo.Content = "[视频消息]"

	case model.MsgTypeEmoticon:
		displayInfo.TypeDesc = "表情消息"
		displayInfo.Content = "[表情消息]"

	case model.MsgTypeLocation:
		displayInfo.TypeDesc = "地理位置消息"
		displayInfo.Content = "[地理位置消息]"

	case model.MsgTypeApp:
		displayInfo.TypeDesc = "APP消息"
		displayInfo.Content = formatAppMessage(content)

	case model.MsgTypeShareCard:
		displayInfo.TypeDesc = "名片消息"
		displayInfo.Content = "[名片消息]"

	case model.MsgTypeVerify:
		displayInfo.TypeDesc = "好友验证消息"
		displayInfo.Content = "[好友验证消息]"

	case model.MsgTypeSystem:
		displayInfo.TypeDesc = "系统消息"
		displayInfo.Content = formatSystemMessage(content)

	case model.MsgTypePrompt:
		displayInfo.TypeDesc = "系统提示消息"
		displayInfo.Content = truncateText(content, 100)

	case model.MsgTypeMicroVideo:
		displayInfo.TypeDesc = "小视频消息"
		displayInfo.Content = "[小视频消息]"

	case model.MsgTypeVoip:
		displayInfo.TypeDesc = "语音通话消息"
		displayInfo.Content = "[语音通话消息]"

	case model.MsgTypeVoipNotify:
		displayInfo.TypeDesc = "语音通话结束消息"
		displayInfo.Content = "[语音通话结束消息]"

	case model.MsgTypeVoipInvite:
		displayInfo.TypeDesc = "语音通话邀请消息"
		displayInfo.Content = "[语音通话邀请消息]"

	case model.MsgTypePossibleFriend:
		displayInfo.TypeDesc = "好友推荐消息"
		displayInfo.Content = "[好友推荐消息]"

	case model.MsgTypeInit:
		displayInfo.TypeDesc = "微信初始化消息"
		displayInfo.Content = "[微信初始化消息]"

	default:
		displayInfo.TypeDesc = "未知消息类型"
		displayInfo.Content = truncateText(content, 100)
	}

	// 检查是否被截断
	if len(content) > len(displayInfo.Content) {
		displayInfo.IsTruncated = true
	}

	return displayInfo
}

// formatAppMessage 格式化APP消息
func formatAppMessage(content string) string {
	// 尝试解析XML格式的APP消息
	if strings.Contains(content, "<msg") && strings.Contains(content, "</msg>") {
		var appMsg AppMsgXML
		if err := xml.Unmarshal([]byte(content), &appMsg); err == nil {
			appType := appMsg.AppMsg.Type
			title := appMsg.AppMsg.Title
			des := appMsg.AppMsg.Des

			switch appType {
			case "5": // 链接消息
				if title != "" {
					return fmt.Sprintf("[链接消息] %s", truncateText(title, 50))
				}
				return "[链接消息]"

			case "6": // 文件消息
				return "[文件消息]"

			case "33", "36": // 位置消息
				return "[位置消息]"

			case "57": // 引用消息
				if title != "" {
					return fmt.Sprintf("[引用消息] %s", truncateText(title, 50))
				}
				return "[引用消息]"

			case "74": // 附件上传中
				return "[附件上传中]"

			case "2000": // 转账消息
				if des != "" {
					return fmt.Sprintf("[转账消息] %s", truncateText(des, 30))
				}
				return "[转账消息]"

			case "2001": // 红包消息
				if des != "" {
					return fmt.Sprintf("[红包消息] %s", truncateText(des, 30))
				}
				return "[红包消息]"

			case "51": // 小程序消息
				if title != "" {
					return fmt.Sprintf("[小程序消息] %s", truncateText(title, 50))
				}
				return "[小程序消息]"

			default:
				if title != "" {
					return fmt.Sprintf("[APP消息-%s] %s", appType, truncateText(title, 50))
				}
				return fmt.Sprintf("[APP消息-%s]", appType)
			}
		}
	}

	// 如果不是XML格式，直接截断
	return truncateText(content, 100)
}

// formatSystemMessage 格式化系统消息
func formatSystemMessage(content string) string {
	// 尝试解析XML格式的系统消息
	if strings.Contains(content, "<sysmsg") && strings.Contains(content, "</sysmsg>") {
		var sysMsg SysMsgXML
		if err := xml.Unmarshal([]byte(content), &sysMsg); err == nil {
			sysType := sysMsg.Type

			switch sysType {
			case "revokemsg":
				return "[消息撤回]"
			case "pat":
				return "[拍一拍]"
			case "sysmsgtemplate":
				// 提取模板内容
				if strings.Contains(content, "template") {
					re := regexp.MustCompile(`<template><!\[CDATA\[(.*?)\]\]></template>`)
					matches := re.FindStringSubmatch(content)
					if len(matches) > 1 {
						return fmt.Sprintf("[系统消息] %s", truncateText(matches[1], 50))
					}
				}
				return "[系统消息]"
			case "mmchatroombarannouncememt":
				return "[群公告]"
			case "roomtoolstips":
				return "[群工具提示]"
			default:
				return fmt.Sprintf("[系统消息-%s]", sysType)
			}
		}
	}

	// 检查是否是邀请进群的消息
	if strings.Contains(content, "邀请你加入群聊") || strings.Contains(content, "invite you to join") {
		return "[邀请进群]"
	}

	// 检查是否是移除群聊的消息
	if strings.Contains(content, "移出了群聊") || strings.Contains(content, "removed from") {
		return "[移除群聊]"
	}

	// 检查是否是解散群聊的消息
	if strings.Contains(content, "已解散该群聊") || strings.Contains(content, "disbanded") {
		return "[解散群聊]"
	}

	// 检查是否是修改群名的消息
	if strings.Contains(content, "修改群名为") || strings.Contains(content, "changed group name") {
		return "[修改群名]"
	}

	// 检查是否是更换群主的消息
	if strings.Contains(content, "群主") || strings.Contains(content, "group owner") {
		return "[更换群主]"
	}

	// 如果不是特殊格式，直接截断
	return truncateText(content, 100)
}

// truncateText 截断文本
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// 如果文本包含换行符，优先在换行符处截断
	if strings.Contains(text, "\n") {
		lines := strings.Split(text, "\n")
		if len(lines[0]) <= maxLength {
			return lines[0] + "..."
		}
	}

	// 如果文本包含空格，尝试在单词边界截断
	if strings.Contains(text, " ") {
		words := strings.Split(text, " ")
		result := ""
		for _, word := range words {
			if len(result)+len(word)+1 <= maxLength {
				if result != "" {
					result += " "
				}
				result += word
			} else {
				break
			}
		}
		if result != "" {
			return result + "..."
		}
	}

	// 直接截断
	return text[:maxLength] + "..."
}
