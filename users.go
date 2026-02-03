package chatclient

import (
	"github.com/paanj-cloud/paanj-go/client"
)

type UsersResource struct {
	client *client.PaanjClient
}

func NewUsersResource(c *client.PaanjClient) *UsersResource {
	return &UsersResource{
		client: c,
	}
}

func (r *UsersResource) GetBlocked() (map[string]interface{}, error) {
	return r.client.GetHttpClient().Request("GET", "/api/v1/users/blocked", nil, false)
}

func (r *UsersResource) User(userId string) *UserContext {
	return &UserContext{
		client:   r.client,
		userId:   userId,
		resource: r,
	}
}

// User Context
type UserContext struct {
	client   *client.PaanjClient
	userId   string
	resource *UsersResource
}

func (u *UserContext) Block() (map[string]interface{}, error) {
	// JS SDK: this.usersResource.block(this.userId)
	// We need generic block functionality in UsersResource or just call API here.
	// Let's call API directly here for simplicity or add to resource if shared.
	// The JS SDK has it in resource, let's follow suit for cleanliness though strict 1:1 isn't required if behavior matches.
	// Actually, let's implement the resource method too if it's cleaner, but for now direct API call is fine for the context.
	// POST /api/v1/users/block { blockedId: u.userId } ??
	// Checking admin chat... Admin uses POST /api/v1/admin/users/block.
	// Client chat uses... let's check JS SDK source.
	// It's likely POST /api/v1/users/block with { userId: targetId } or similar.
	// JS SDK `UsersResource.block` calls `this.httpClient.request('POST', '/api/v1/users/block', { userId })`
	return u.client.GetHttpClient().Request("POST", "/api/v1/users/block", map[string]interface{}{
		"userId": u.userId,
	}, false)
}

func (u *UserContext) Unblock() (map[string]interface{}, error) {
	return u.client.GetHttpClient().Request("POST", "/api/v1/users/unblock", map[string]interface{}{
		"userId": u.userId,
	}, false)
}

func (r *UsersResource) OnTokenRefresh(callback func(interface{})) {
	r.client.On("token.updated", callback)
}
