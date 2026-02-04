package chatclient

// Test sync: 2026-02-04

type ChatClient struct {
	client        *client_go.PaanjClient
	Conversations *ConversationsResource
	Users         *UsersResource
}

func NewChatClient(c *client_go.PaanjClient) *ChatClient {
	chatClient := &ChatClient{
		client: c,
	}

	chatClient.Conversations = NewConversationsResource(c)
	chatClient.Users = NewUsersResource(c)

	return chatClient
}
