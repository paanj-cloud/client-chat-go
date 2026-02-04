package chatclient

// Test sync: 2026-02-04
import (
	"github.com/paanj-cloud/paanj-go/client"
)

type ChatClient struct {
	client        *client.PaanjClient
	Conversations *ConversationsResource
	Users         *UsersResource
}

func NewChatClient(c *client.PaanjClient) *ChatClient {
	chatClient := &ChatClient{
		client: c,
	}

	chatClient.Conversations = NewConversationsResource(c)
	chatClient.Users = NewUsersResource(c)

	return chatClient
}
