package chatclient

import (
	"fmt"

	"github.com/paanj-cloud/client-go"
)

type UsersResource struct {
	runtime chatRuntime
}

func NewUsersResource(c *client.PaanjClient) *UsersResource {
	return newUsersResource(newPaanjRuntime(c))
}

func newUsersResource(runtime chatRuntime) *UsersResource {
	return &UsersResource{runtime: runtime}
}

func (r *UsersResource) GetBlocked() (map[string]interface{}, error) {
	return r.runtime.Request("GET", "/api/v1/users/me/blocked", nil, false)
}

func (r *UsersResource) Block(userId string) (map[string]interface{}, error) {
	return r.runtime.Request("POST", fmt.Sprintf("/api/v1/users/%s/block", userId), nil, false)
}

func (r *UsersResource) Unblock(userId string) (map[string]interface{}, error) {
	return r.runtime.Request("POST", fmt.Sprintf("/api/v1/users/%s/unblock", userId), nil, false)
}

func (r *UsersResource) User(userId string) *UserContext {
	return &UserContext{
		userId:   userId,
		resource: r,
	}
}

// User Context
type UserContext struct {
	userId   string
	resource *UsersResource
}

func (u *UserContext) Block() (map[string]interface{}, error) {
	return u.resource.Block(u.userId)
}

func (u *UserContext) Unblock() (map[string]interface{}, error) {
	return u.resource.Unblock(u.userId)
}

func (r *UsersResource) OnTokenRefresh(callback func(interface{})) {
	r.runtime.On("token.updated", callback)
}

func (u *UserContext) OnTokenRefresh(callback func(interface{})) {
	u.resource.OnTokenRefresh(callback)
}
