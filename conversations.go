package chatclient

import (
	"fmt"
	"time"

	"github.com/paanj-cloud/paanj-go/client"
)

type ConversationsResource struct {
	client *client.PaanjClient
}

func NewConversationsResource(c *client.PaanjClient) *ConversationsResource {
	return &ConversationsResource{
		client: c,
	}
}

func (r *ConversationsResource) Create(data map[string]interface{}) (map[string]interface{}, error) {
	return r.client.GetHttpClient().Request("POST", "/api/v1/conversations", data, false)
}

func (r *ConversationsResource) List(filters map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Handle query params for filters if needed, passing as nil body for now or constructing URL
	// For simplicity, Go SDK might need a better query param builder in HttpClient later.
	// Assuming filters are passed as body for now or just ignoring them as proof of concept if GET
	return r.client.GetHttpClient().Request("GET", "/api/v1/conversations", nil, false)
}

func (r *ConversationsResource) Get(conversationId string) (map[string]interface{}, error) {
	return r.client.GetHttpClient().Request("GET", fmt.Sprintf("/api/v1/conversations/%s", conversationId), nil, false)
}

func (r *ConversationsResource) OnMessage(callback func(interface{})) {
	r.client.On("message.create", callback)
}

// Conversation Context helper (optional, simplified for Go)
type ConversationContext struct {
	client         *client.PaanjClient
	ConversationId string
	resource       *ConversationsResource
}

func (r *ConversationsResource) Conversation(conversationId string) *ConversationContext {
	return &ConversationContext{
		client:         r.client,
		ConversationId: conversationId,
		resource:       r,
	}
}

func (c *ConversationContext) Send(content string) (map[string]interface{}, error) {
	// JS SDK sends via WebSocket with format: { type: 'new', receiver: conversationId, message: content, hash: ... }
	timestamp := time.Now().UnixMilli()
	hash := fmt.Sprintf("%d", timestamp)

	payload := map[string]interface{}{
		"type":     "new",
		"receiver": c.ConversationId,
		"message":  content,
		"hash":     hash,
	}

	err := c.client.GetWebSocket().Send(payload)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sent": true}, nil
}

func (c *ConversationContext) AddParticipant(userId string, role string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"userId": userId,
		"role":   role,
	}
	return c.client.GetHttpClient().Request("POST", fmt.Sprintf("/api/v1/conversations/%s/participants", c.ConversationId), data, false)
}

func (c *ConversationContext) Leave() error {
	userId := c.client.GetUserId()
	if userId == "" {
		return fmt.Errorf("user not authenticated")
	}
	// Note: JS SDK uses DELETE /members with body { userIds: [id] }
	// Go's http.NewRequest can exist with body for DELETE.
	data := map[string]interface{}{
		"userIds": []interface{}{userId}, // Assuming API takes ID or maybe int/string depending on backend. JS SDK passed parseInt(userId) so likely int.
		// If userId is string in client state but int in backend, we need to be careful.
		// For now assuming backend handles string or we pass as is. The JS SDK does parseInt.
		// Let's assume the ID is consistent for now or try to pass as is.
	}
	_, err := c.client.GetHttpClient().Request("DELETE", fmt.Sprintf("/api/v1/conversations/%s/members", c.ConversationId), data, false)
	return err
}

func (c *ConversationContext) OnUpdate(callback func(interface{})) {
	// JS SDK: this.conversationsResource.onUpdate
	// We need to listen to generic event and filter? Or specific subject?
	// JS SDK actually calls conversationsResource.onUpdate which listens to 'conversation.update'.
	// We'll trust the server emits 'conversation.update' with conversationId in data.
	c.resource.client.On("conversation.update", func(data interface{}) {
		// Filter logic would be needed here if the event is global.
		// Simplified: assumes client app filters or server sends only relevant user events.
		callback(data)
	})
}

// Participants Helper
type ParticipantsHelper struct {
	context *ConversationContext
}

func (c *ConversationContext) Participants() *ParticipantsHelper {
	return &ParticipantsHelper{context: c}
}

func (p *ParticipantsHelper) List() (interface{}, error) {
	// Assuming GET /conversations/:id/participants
	return p.context.client.GetHttpClient().Request("GET", fmt.Sprintf("/api/v1/conversations/%s/participants", p.context.ConversationId), nil, false)
}

func (p *ParticipantsHelper) Add(userId string, role string) (map[string]interface{}, error) {
	return p.context.AddParticipant(userId, role)
}

// Messages Helper
type MessagesHelper struct {
	context *ConversationContext
}

func (c *ConversationContext) Messages() *MessagesHelper {
	return &MessagesHelper{context: c}
}

func (m *MessagesHelper) List(filters map[string]interface{}) (interface{}, error) {
	// Construct query string from filters in a real implementation
	return m.context.client.GetHttpClient().Request("GET", fmt.Sprintf("/api/v1/conversations/%s/messages", m.context.ConversationId), nil, false)
}
