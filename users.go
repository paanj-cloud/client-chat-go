package chatclient

import (
	"fmt"

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
	// JS SDK: POST /api/v1/users/:userId/block
	return u.client.GetHttpClient().Request("POST", fmt.Sprintf("/api/v1/users/%s/block", u.userId), nil, false)
}

func (u *UserContext) Unblock() (map[string]interface{}, error) {
	// JS SDK: POST /api/v1/users/:userId/unblock
	return u.client.GetHttpClient().Request("POST", fmt.Sprintf("/api/v1/users/%s/unblock", u.userId), nil, false)
}

func (r *UsersResource) OnTokenRefresh(callback func(interface{})) {
	r.client.On("token.updated", callback)
}
