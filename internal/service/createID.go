package service

import "fmt"

// GenerateId 根据发送方和接收方生成唯一ID
func GenerateId(sendID string, replyID string) (generateID string) {
	generateID = fmt.Sprintf("%s->%s", sendID, replyID)
	return
}
