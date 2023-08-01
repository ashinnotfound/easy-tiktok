package logic

import (
	"context"
	pb "easy-tiktok/apps/social/proto"
)

// MessageServiceImpl //
// MessageService接口的实现类
// Author lql
type MessageServiceImpl struct {
	pb.MessageServer
}

// Chat //
// 获取对话消息
func (impl *MessageServiceImpl) Chat(ctx context.Context, request *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {

	return nil, nil
}

// Action //
// 消息操作：发送消息
func (impl *MessageServiceImpl) Action(ctx context.Context, request *pb.DouyinMessageActionRequest) (*pb.DouyinMessageActionResponse, error) {

	return nil, nil
}
