package model

import "gorm.io/gorm"

type SQLMessage struct {
	gorm.Model
	SendID    string `json:"send_id"`
	ReceiveID string `json:"receive_id"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
}

type UnreadCount struct {
	SendID      string `json:"send_id"`
	UnreadCount int    `json:"unread_count"`
}

type UnreadMessage struct {
	TotalCounts  int           `json:"total_counts"`
	UnreadCounts []UnreadCount `json:"unread_counts"`
}
