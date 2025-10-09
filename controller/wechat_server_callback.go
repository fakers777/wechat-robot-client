package controller

import (
	"log"
	"wechat-robot-client/dto"
	"wechat-robot-client/pkg/appx"
	"wechat-robot-client/pkg/robot"
	"wechat-robot-client/service"

	"github.com/gin-gonic/gin"
)

type WechatServerCallback struct {
}

func NewWechatServerCallbackController() *WechatServerCallback {
	return &WechatServerCallback{}
}

func (ct *WechatServerCallback) SyncMessageCallback(c *gin.Context) {
	wechatID := c.Param("wechatID")
	log.Printf("Received SyncMessageCallback for wechatID: %s", wechatID)
	var req robot.ClientResponse[robot.SyncMessage]
	resp := appx.NewResponse(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ToErrorResponse(err)
		return
	}

	// 打印收到的消息详情
	log.Printf("=== 收到消息回调详情 ===")
	log.Printf("微信ID: %s", wechatID)
	log.Printf("消息数量: %d", len(req.Data.AddMsgs))
	log.Printf("状态: %d, 时间: %d", req.Data.Status, req.Data.Time)

	// 打印每条消息的详细信息
	for i, msg := range req.Data.AddMsgs {
		log.Printf("--- 消息 %d ---", i+1)
		log.Printf("消息ID: %d", msg.NewMsgId)
		log.Printf("发送者: %s", *msg.FromUserName.String)
		log.Printf("接收者: %s", *msg.ToUserName.String)
		log.Printf("创建时间: %d", msg.CreateTime)
		//log.Printf("消息来源: %s", msg.MsgSource)
		log.Printf("状态: %d", msg.Status)

		// 使用优化的消息显示
		displayInfo := FormatMessageForDisplay(msg)
		log.Printf("消息类型: %s (%s)", displayInfo.Type, displayInfo.TypeDesc)
		log.Printf("消息内容: %s", displayInfo.Content)
		if displayInfo.IsTruncated {
			log.Printf("⚠️  消息内容已截断 (原始长度: %d 字符)", displayInfo.Length)
		}
		if msg.PushContent != "" && msg.PushContent != displayInfo.Content {
			log.Printf("推送内容: %s", truncateText(msg.PushContent, 100))
		}
	}
	log.Printf("=== 消息回调详情结束 ===")

	service.NewLoginService(c).SyncMessageCallback(wechatID, req.Data)

	resp.ToResponse(nil)
}

func (ct *WechatServerCallback) LogoutCallback(c *gin.Context) {
	wechatID := c.Param("wechatID")
	log.Printf("Received LogoutCallback for wechatID: %s", wechatID)
	var req dto.LogoutNotificationRequest
	resp := appx.NewResponse(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("LogoutCallback binding error: %v", err)
		resp.ToErrorResponse(err)
		return
	}
	err := service.NewLoginService(c).LogoutCallback(req)
	if err != nil {
		log.Printf("LogoutCallback failed: %v\n", err)
		resp.ToErrorResponse(err)
		return
	}
	resp.ToResponse(nil)
}
