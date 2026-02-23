package chatclient

// Test sync: 2026-02-04
import (
	"github.com/paanj-cloud/client-go"
)

type ChatClient struct {
	Conversations *ConversationsResource
	Users         *UsersResource
}

func NewChatClient(c *client.PaanjClient) *ChatClient {
	runtime := newPaanjRuntime(c)

	chatClient := &ChatClient{}
	chatClient.Conversations = newConversationsResource(runtime)
	chatClient.Users = newUsersResource(runtime)

	return chatClient
}

func (c *ChatClient) Conversation(conversationId string) *ConversationContext {
	return c.Conversations.Conversation(conversationId)
}

func (c *ChatClient) User(userId string) *UserContext {
	return c.Users.User(userId)
}
