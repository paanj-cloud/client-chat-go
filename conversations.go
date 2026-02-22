package chatclient

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/paanj-cloud/paanj-go/client"
)

type ConversationsResource struct {
	runtime chatRuntime
}

func NewConversationsResource(c *client.PaanjClient) *ConversationsResource {
	return newConversationsResource(newPaanjRuntime(c))
}

func newConversationsResource(runtime chatRuntime) *ConversationsResource {
	return &ConversationsResource{runtime: runtime}
}

func (r *ConversationsResource) Create(data map[string]interface{}) (map[string]interface{}, error) {
	// Transform participants to members format like JS SDK
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[k] = v
	}

	// API expects 'members', not 'participants'
	if participants, ok := data["participants"]; ok {
		members := buildMembersPayload(participants)
		if len(members) > 0 {
			payload["members"] = members
			delete(payload, "participants")
		}
	}

	return r.runtime.Request("POST", "/api/v1/conversations", payload, false)
}

func (r *ConversationsResource) List(filters map[string]interface{}) (map[string]interface{}, error) {
	userID := r.runtime.UserID()
	if userID == "" {
		return nil, fmt.Errorf("user not authenticated")
	}

	path := fmt.Sprintf("/api/v1/users/%s/conversations", userID)
	query := buildPaginationQuery(filters)
	if query != "" {
		path = fmt.Sprintf("%s?%s", path, query)
	}

	return r.runtime.Request("GET", path, nil, false)
}

func (r *ConversationsResource) Get(conversationId string) (map[string]interface{}, error) {
	return r.runtime.Request("GET", fmt.Sprintf("/api/v1/conversations/%s", conversationId), nil, false)
}

func (r *ConversationsResource) OnMessage(callback func(interface{})) {
	r.runtime.On("message.create", callback)
}

func (r *ConversationsResource) OnUpdate(conversationId string, callback func(interface{})) {
	_ = r.runtime.Subscribe(map[string]interface{}{
		"type":     "subscribe",
		"resource": "conversation",
		"id":       conversationId,
		"events":   []string{"conversation.update"},
	})

	r.runtime.On(fmt.Sprintf("conversation:%s:conversation.update", conversationId), callback)
}

// Conversation Context helper (optional, simplified for Go)
type ConversationContext struct {
	ConversationId string
	resource       *ConversationsResource
}

func (r *ConversationsResource) Conversation(conversationId string) *ConversationContext {
	return &ConversationContext{
		ConversationId: conversationId,
		resource:       r,
	}
}

func (c *ConversationContext) Send(content string, metadata ...map[string]interface{}) (map[string]interface{}, error) {
	// JS SDK sends via WebSocket with format: { type: 'new', receiver: conversationId, message: content, hash: ... }
	hash := fmt.Sprintf("%d+%d", time.Now().UnixNano()%1000000, time.Now().UnixMilli())

	payload := map[string]interface{}{
		"type":     "new",
		"receiver": c.ConversationId,
		"message":  content,
		"hash":     hash,
	}
	if len(metadata) > 0 && metadata[0] != nil {
		payload["metadata"] = metadata[0]
	}

	err := c.resource.runtime.Send(payload)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sent": true}, nil
}

func (c *ConversationContext) AddParticipant(userId string, role string) (map[string]interface{}, error) {
	if role == "" {
		role = "member"
	}

	data := map[string]interface{}{
		"members": []map[string]interface{}{
			{
				"userId": normalizeUserID(userId),
				"role":   role,
			},
		},
	}
	return c.resource.runtime.Request("POST", fmt.Sprintf("/api/v1/conversations/%s/members", c.ConversationId), data, false)
}

func (c *ConversationContext) Leave() error {
	userId := c.resource.runtime.UserID()
	if userId == "" {
		return fmt.Errorf("user not authenticated")
	}

	data := map[string]interface{}{
		"userIds": []interface{}{normalizeUserID(userId)},
	}
	_, err := c.resource.runtime.Request("DELETE", fmt.Sprintf("/api/v1/conversations/%s/members", c.ConversationId), data, false)
	return err
}

func (c *ConversationContext) Get() (map[string]interface{}, error) {
	return c.resource.Get(c.ConversationId)
}

func (c *ConversationContext) OnUpdate(callback func(interface{})) {
	c.resource.OnUpdate(c.ConversationId, callback)
}

func (c *ConversationContext) OnMessage(callback func(interface{})) {
	c.resource.runtime.On(fmt.Sprintf("conversation:%s:message.create", c.ConversationId), callback)
}

// Participants Helper
type ParticipantsHelper struct {
	context *ConversationContext
}

func (c *ConversationContext) Participants() *ParticipantsHelper {
	return &ParticipantsHelper{context: c}
}

func (p *ParticipantsHelper) List() (interface{}, error) {
	// JS SDK gets conversation and returns members field
	conv, err := p.context.resource.Get(p.context.ConversationId)
	if err != nil {
		return nil, err
	}
	// Return members field if it exists
	if members, ok := conv["members"]; ok {
		return members, nil
	}
	if members, ok := conv["participants"]; ok {
		return members, nil
	}
	return []interface{}{}, nil
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
	path := fmt.Sprintf("/api/v1/conversations/%s/messages", m.context.ConversationId)
	query := buildMessageQuery(filters)
	if query != "" {
		path = fmt.Sprintf("%s?%s", path, query)
	}

	response, err := m.context.resource.runtime.Request("GET", path, nil, false)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return []interface{}{}, nil
	}
	if messages, ok := response["messages"]; ok {
		return messages, nil
	}

	return response, nil
}

func buildMembersPayload(participants interface{}) []map[string]interface{} {
	switch partArray := participants.(type) {
	case []map[string]string:
		members := make([]map[string]interface{}, 0, len(partArray))
		for _, p := range partArray {
			role := p["role"]
			if role == "" {
				role = "member"
			}
			members = append(members, map[string]interface{}{
				"userId": normalizeUserID(p["userId"]),
				"role":   role,
			})
		}
		return members
	case []map[string]interface{}:
		members := make([]map[string]interface{}, 0, len(partArray))
		for _, p := range partArray {
			userID, ok := p["userId"]
			if !ok {
				continue
			}
			role := "member"
			if rawRole, ok := p["role"].(string); ok && rawRole != "" {
				role = rawRole
			}

			members = append(members, map[string]interface{}{
				"userId": normalizeUserID(userID),
				"role":   role,
			})
		}
		return members
	case []interface{}:
		members := make([]map[string]interface{}, 0, len(partArray))
		for _, raw := range partArray {
			participant, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}

			userID, ok := participant["userId"]
			if !ok {
				continue
			}
			role := "member"
			if rawRole, ok := participant["role"].(string); ok && rawRole != "" {
				role = rawRole
			}

			members = append(members, map[string]interface{}{
				"userId": normalizeUserID(userID),
				"role":   role,
			})
		}
		return members
	default:
		return nil
	}
}

func normalizeUserID(userID interface{}) interface{} {
	switch v := userID.(type) {
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
		return v
	case float64:
		return int(v)
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		return userID
	}
}

func buildPaginationQuery(filters map[string]interface{}) string {
	if len(filters) == 0 {
		return ""
	}

	values := url.Values{}
	if rawLimit, ok := filters["limit"]; ok {
		appendIntQuery(values, "limit", rawLimit)
	}
	if rawOffset, ok := filters["offset"]; ok {
		appendIntQuery(values, "offset", rawOffset)
	}

	return values.Encode()
}

func buildMessageQuery(filters map[string]interface{}) string {
	if len(filters) == 0 {
		return ""
	}

	values := url.Values{}
	if rawLimit, ok := filters["limit"]; ok {
		appendIntQuery(values, "limit", rawLimit)
	}
	if rawOffset, ok := filters["offset"]; ok {
		appendIntQuery(values, "offset", rawOffset)
	}
	if rawBefore, ok := filters["before"]; ok {
		appendStringQuery(values, "before", rawBefore)
	}
	if rawAfter, ok := filters["after"]; ok {
		appendStringQuery(values, "after", rawAfter)
	}

	return values.Encode()
}

func appendIntQuery(values url.Values, key string, value interface{}) {
	switch v := value.(type) {
	case int:
		values.Set(key, strconv.Itoa(v))
	case int32:
		values.Set(key, strconv.Itoa(int(v)))
	case int64:
		values.Set(key, strconv.FormatInt(v, 10))
	case float64:
		values.Set(key, strconv.Itoa(int(v)))
	case string:
		if _, err := strconv.Atoi(v); err == nil {
			values.Set(key, v)
		}
	}
}

func appendStringQuery(values url.Values, key string, value interface{}) {
	if str, ok := value.(string); ok && str != "" {
		values.Set(key, str)
	}
}
