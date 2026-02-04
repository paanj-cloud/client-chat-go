# Paanj Chat Client SDK for Go

Official Go Chat Client SDK for Paanj - Build real-time chat applications with ease.

[![Go Reference](https://pkg.go.dev/badge/github.com/paanj-cloud/client-chat-go.svg)](https://pkg.go.dev/github.com/paanj-cloud/client-chat-go)

## Installation

```bash
go get github.com/paanj-cloud/client-go@latest
go get github.com/paanj-cloud/client-chat-go@latest
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    client "github.com/paanj-cloud/client-go"
    chatclient "github.com/paanj-cloud/client-chat-go"
)

func main() {
    // Initialize base client
    paanjClient := client.NewClient(client.ClientOptions{
        ApiKey: "pk_live_your_api_key",
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })

    // Authenticate
    auth, _ := paanjClient.AuthenticateAnonymous(map[string]interface{}{
        "name": "Alice",
    }, nil)

    // Connect to WebSocket
    paanjClient.Connect()
    defer paanjClient.Disconnect()

    // Initialize chat client
    chat := chatclient.NewChatClient(paanjClient)

    // Create a conversation
    conv, err := chat.Conversations.Create(map[string]interface{}{
        "name": "My First Chat",
        "participants": []map[string]string{
            {"userId": "user2_id", "role": "member"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Conversation created: %s\\n", conv["id"])
}
```

## Features

- ✅ Create and manage conversations
- ✅ Send and receive messages in real-time
- ✅ Manage participants
- ✅ Real-time message events
- ✅ User management
- ✅ Message history

## Complete Examples

### Example 1: Basic Chat Application

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    client "github.com/paanj-cloud/client-go"
    chatclient "github.com/paanj-cloud/client-chat-go"
)

func main() {
    // Initialize client
    paanjClient := client.NewClient(client.ClientOptions{
        ApiKey: "pk_live_your_api_key",
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })

    // Authenticate
    auth, err := paanjClient.AuthenticateAnonymous(map[string]interface{}{
        "name": "Alice",
    }, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✅ Authenticated as: %s\\n", auth.UserId)

    // Connect to WebSocket
    if err := paanjClient.Connect(); err != nil {
        log.Fatal(err)
    }
    defer paanjClient.Disconnect()

    // Initialize chat
    chat := chatclient.NewChatClient(paanjClient)

    // Listen for incoming messages
    chat.Conversations.OnMessage(func(data interface{}) {
        msg := data.(map[string]interface{})
        fmt.Printf("📩 New message: %s\\n", msg["content"])
    })

    // Create a conversation
    conv, err := chat.Conversations.Create(map[string]interface{}{
        "name": "Team Chat",
        "participants": []map[string]string{
            {"userId": "user2_id", "role": "member"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    conversationId := conv["id"].(string)
    fmt.Printf("✅ Conversation created: %s\\n", conversationId)

    // Send a message
    _, err = chat.Conversations.Conversation(conversationId).Send("Hello, team!")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("✅ Message sent!")

    // Keep alive to receive messages
    time.Sleep(5 * time.Second)
}
```

### Example 2: Two-User Chat

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    client "github.com/paanj-cloud/client-go"
    chatclient "github.com/paanj-cloud/client-chat-go"
)

func main() {
    // User 1
    client1 := client.NewClient(client.ClientOptions{
        ApiKey: "pk_live_your_api_key",
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })
    auth1, _ := client1.AuthenticateAnonymous(map[string]interface{}{
        "name": "Alice",
    }, nil)
    client1.Connect()
    defer client1.Disconnect()
    chat1 := chatclient.NewChatClient(client1)

    // User 2
    client2 := client.NewClient(client.ClientOptions{
        ApiKey: "pk_live_your_api_key",
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })
    auth2, _ := client2.AuthenticateAnonymous(map[string]interface{}{
        "name": "Bob",
    }, nil)
    client2.Connect()
    defer client2.Disconnect()
    chat2 := chatclient.NewChatClient(client2)

    // Create conversation with both users
    conv, _ := chat1.Conversations.Create(map[string]interface{}{
        "name": "Alice & Bob",
        "participants": []map[string]string{
            {"userId": auth2.UserId, "role": "member"},
        },
    })
    conversationId := conv["id"].(string)

    // Wait for membership to propagate
    time.Sleep(2 * time.Second)

    // User 2 listens for messages
    chat2.Conversations.OnMessage(func(data interface{}) {
        msg := data.(map[string]interface{})
        fmt.Printf("Bob received: %s\\n", msg["content"])
    })

    // User 1 sends message
    chat1.Conversations.Conversation(conversationId).Send("Hi Bob!")

    // Keep alive
    time.Sleep(3 * time.Second)
}
```

### Example 3: Group Chat with Participants Management

```go
package main

import (
    "fmt"
    "log"
    
    client "github.com/paanj-cloud/client-go"
    chatclient "github.com/paanj-cloud/client-chat-go"
)

func main() {
    // Initialize
    paanjClient := client.NewClient(client.ClientOptions{
        ApiKey: "pk_live_your_api_key",
        ApiUrl: "https://api1.paanj.com",
        WsUrl:  "wss://ws1.paanj.com",
    })
    auth, _ := paanjClient.AuthenticateAnonymous(map[string]interface{}{
        "name": "Admin",
    }, nil)
    paanjClient.Connect()
    defer paanjClient.Disconnect()

    chat := chatclient.NewChatClient(paanjClient)

    // Create group conversation
    conv, _ := chat.Conversations.Create(map[string]interface{}{
        "name": "Project Team",
        "participants": []map[string]string{
            {"userId": "user2_id", "role": "member"},
            {"userId": "user3_id", "role": "member"},
        },
    })
    conversationId := conv["id"].(string)

    // Get conversation context
    conversation := chat.Conversations.Conversation(conversationId)

    // List participants
    participants, _ := conversation.Participants().List()
    fmt.Printf("Participants: %+v\\n", participants)

    // Add a new participant
    conversation.Participants().Add("user4_id", "member")
    fmt.Println("✅ New participant added")

    // Send message to group
    conversation.Send("Welcome to the team!")

    // Get message history
    messages, _ := conversation.Messages().List(nil)
    fmt.Printf("Message history: %+v\\n", messages)
}
```

## API Reference

### Conversations

#### Create Conversation

```go
conv, err := chat.Conversations.Create(map[string]interface{}{
    "name": "Conversation Name",
    "participants": []map[string]string{
        {"userId": "user_id", "role": "member"},
    },
    "metadata": map[string]interface{}{
        "topic": "general",
    },
})
```

#### List Conversations

```go
conversations, err := chat.Conversations.List(nil)
```

#### Get Conversation

```go
conv, err := chat.Conversations.Get(conversationId)
```

#### Listen for Messages

```go
chat.Conversations.OnMessage(func(data interface{}) {
    msg := data.(map[string]interface{})
    fmt.Printf("Message: %s from %s\\n", msg["content"], msg["senderId"])
})
```

### Conversation Context

#### Send Message

```go
conversation := chat.Conversations.Conversation(conversationId)
result, err := conversation.Send("Hello, world!")
```

#### Add Participant

```go
result, err := conversation.AddParticipant(userId, "member")
```

#### Leave Conversation

```go
err := conversation.Leave()
```

#### Listen for Updates

```go
conversation.OnUpdate(func(data interface{}) {
    fmt.Printf("Conversation updated: %+v\\n", data)
})
```

### Participants

#### List Participants

```go
participants, err := conversation.Participants().List()
```

#### Add Participant

```go
result, err := conversation.Participants().Add(userId, "admin")
```

### Messages

#### Get Message History

```go
messages, err := conversation.Messages().List(map[string]interface{}{
    "limit":  50,
    "offset": 0,
})
```

### Users

#### Block User

```go
result, err := chat.Users.Block(blockedUserId)
```

#### Unblock User

```go
result, err := chat.Users.Unblock(blockedUserId)
```

#### List Blocked Users

```go
blockedUsers, err := chat.Users.ListBlocked()
```

## Event Types

| Event | Description | Data |
|-------|-------------|------|
| `message.create` | New message received | `{content, senderId, conversationId, timestamp, ...}` |
| `conversation.update` | Conversation updated | `{conversationId, ...}` |

## Error Handling

```go
conv, err := chat.Conversations.Create(data)
if err != nil {
    log.Printf("Failed to create conversation: %v", err)
    return
}

_, err = conversation.Send("Hello")
if err != nil {
    log.Printf("Failed to send message: %v", err)
    return
}
```

## Best Practices

1. **Wait for membership propagation** after creating conversations (2 seconds recommended)
2. **Handle errors gracefully** - network issues can occur
3. **Use context for conversation operations** - cleaner API
4. **Listen for events before sending** - ensures you don't miss messages
5. **Clean up connections** - always defer `Disconnect()`

## License

MIT License - see LICENSE file for details.

## Support

- Documentation: https://docs.paanj.com
- Issues: https://github.com/paanj-cloud/client-chat-go/issues
- Email: support@paanj.com
