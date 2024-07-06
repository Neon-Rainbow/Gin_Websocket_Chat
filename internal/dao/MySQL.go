package dao

import (
	"gorm.io/gorm"
	"websocket/model"
	"websocket/pkg/MySQL"
)

var database *gorm.DB

// GetHistoryMessage 获取历史消息
func GetHistoryMessage(SendID string, ReplyID string, HistoryCount int) ([]model.SQLMessage, error) {
	database = MySQL.MySQL
	var historyMessage []model.SQLMessage
	var result *gorm.DB
	if HistoryCount == -1 {
		result = database.
			Where("(send_id = ? AND receive_id = ?) OR (send_id = ? AND receive_id = ?)", SendID, ReplyID, ReplyID, SendID).
			Order("created_at DESC").Find(&historyMessage)
	} else {
		result = database.
			Where("(send_id = ? AND receive_id = ?) OR (send_id = ? AND receive_id = ?)", SendID, ReplyID, ReplyID, SendID).
			Limit(HistoryCount).
			Order("created_at DESC").Find(&historyMessage)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return historyMessage, nil
}

// SaveMessage 保存消息
func SaveMessage(sendID string, replyID string, s string, isRead bool) error {
	database = MySQL.MySQL
	msg := &model.SQLMessage{
		SendID:    sendID,
		ReceiveID: replyID,
		Content:   s,
		IsRead:    isRead,
	}
	result := database.Save(msg)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUnreadMessageCounts 获取某个人的未读消息数量以及每个发送者发送的未读消息数量
func GetUnreadMessageCounts(receiveID string) (int, []model.UnreadCount, error) {
	database = MySQL.MySQL
	// 获取总的未读消息数量
	var totalUnreadCount int64
	if err := database.Model(&model.SQLMessage{}).
		Where("receive_id = ? AND is_read = ?", receiveID, false).
		Count(&totalUnreadCount).Error; err != nil {
		return 0, nil, err
	}

	// 获取每个发送者的未读消息数量
	var unreadCounts []model.UnreadCount
	if err := database.Model(&model.SQLMessage{}).
		Select("send_id, COUNT(*) as unread_count").
		Where("receive_id = ? AND is_read = ?", receiveID, false).
		Group("send_id").
		Scan(&unreadCounts).Error; err != nil {
		return 0, nil, err
	}

	return int(totalUnreadCount), unreadCounts, nil
}

// GetUnreadMessages 获取未读消息
func GetUnreadMessages(sendID string, receiveID string) ([]model.SQLMessage, error) {
	database = MySQL.MySQL
	var messages []model.SQLMessage
	result := database.Where("send_id = ? AND receive_id = ? AND is_read = ?", sendID, receiveID, false).Find(&messages)
	return messages, result.Error
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(receiveID string) error {
	database = MySQL.MySQL
	result := database.Model(&model.SQLMessage{}).Where("receive_id = ? AND is_read = ?", receiveID, false).Update("is_read", true)
	return result.Error
}
