package chatclient

import "github.com/paanj-cloud/paanj-go/client"

// chatRuntime is a narrow runtime contract used by chat resources.
// It keeps resource logic decoupled from the concrete core client.
type chatRuntime interface {
	Request(method, path string, body interface{}, skipAuth bool) (map[string]interface{}, error)
	Send(data interface{}) error
	UserID() string
	On(event string, callback func(interface{}))
	Subscribe(subscription interface{}) error
}

type paanjRuntime struct {
	client *client.PaanjClient
}

func newPaanjRuntime(c *client.PaanjClient) chatRuntime {
	return &paanjRuntime{client: c}
}

func (r *paanjRuntime) Request(method, path string, body interface{}, skipAuth bool) (map[string]interface{}, error) {
	return r.client.GetHttpClient().Request(method, path, body, skipAuth)
}

func (r *paanjRuntime) Send(data interface{}) error {
	return r.client.GetWebSocket().Send(data)
}

func (r *paanjRuntime) UserID() string {
	return r.client.GetUserId()
}

func (r *paanjRuntime) On(event string, callback func(interface{})) {
	r.client.On(event, callback)
}

func (r *paanjRuntime) Subscribe(subscription interface{}) error {
	return r.client.Subscribe(subscription)
}
