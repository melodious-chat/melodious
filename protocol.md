# Melodious Chat Protocol (websocket)

Communication is established over a WebSocket _connection_ with JSON _messages_ used to exchange data.

## Messages

Each message is a JSON document which necessarily has `type` field which describes message type

Each message MAY have an `_id` field. All responses to such message MUST contain the same `_id`.
If length of the id is more than 64 characters, then only first 64 characters are used. 

### quit (sent by server and client)

```json
{
    "type": "quit",
    "message": "<string>",
}
```

Sent by both server and client to indicate that one of them is no longer interested in the connection.

After sending/receiving server and client MUST close the connection.

### fatal (sent by server)

```json
{
    "type": "fatal",
    "message": "<string>",
}
```

Send by server to indicate that a fatal error has occured and the connection must be closed.

Server MUST close the connection after sending this message.

### note (sent by server)

```json
{
    "type": "note",
    "message": "<string>"
}
```

### ok, fail (sent by server)

```json
{
    "type": "ok",
    "message": "<string>"
}
```

```json
{
    "type": "fail",
    "message": "<string>"
}
```

These messages are used to notify user about results of operations started by the user.

### register 

```json
{
    "type": "register",
    "name": "<string>",
    "pass": "<string>"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

pass: Password. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

Sent by client: Registers the client on the server.

Sent by server: Indicates a user register event.

If user with username `name` already exists, server MUST send a `fatal` message.

Server MAY introduce additional protection like IP duplication protection.

If this is a first user ever registered, server MUST give that user admin permissions and send him a corresponding `note` message.

It is recommended that server stores hash sum of the hash sum of the password to prevent heavy damage on database leak.

After registering user MUST be treated as logged in.

### login

```json
{
    "type": "login",
    "name": "<string>",
    "pass": "<string>"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

pass: Password. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

Sent by client: Logs the client in.

Sent by server: Indicates a login (online) event.

If user with username `name` does not exist or SHA256 hash/checksum `hash` is invalid, server MUST send a `fatal` message.

Server MAY introduce additional protection like banning users from connecting.

If user has administrator privileges, server MUST send him a corresponding `note` message.

It is recommended that server stores hash sum of the hash sum of the password to prevent heavy damage on database leak.

### new-channel

```json 
{
    "type": "new-channel",
    "name": "<string>",
    "topic": "<string>"
}
```

User needs perms.new-channel(user=$user,channel=$channel) flag or owner status to do that.

name: channel name; maximum 32 characters

Sent by client: Creates a new channel. If such channel already exists or user is not an owner, server MUST return a `fail` message.

Sent by server: Notifies about a new channel.

### delete-channel

```json 
{
    "type": "delete-channel",
    "name": "<string>"
}
```

User needs perms.delete-channel flag or owner status to do that.

name: channel name; maximum 32 characters

Sent by client: Deletes a channel.

Sent by server: Notifies about a deleted channel. May mention non-existing channel.

### channel-topic (sent by client)

```json
{
    "type": "channel-topic",
    "name": "<string>",
    "topic": "<string>"
}
```

User needs perms.channel-topic flag or owner status to do that.

name: channel name; maximum 32 characters

topic: channel topic message; maximum 1024 characters

Changes a channel's topic.

### subscribe (sent by client)

```json
{
    "type": "subscribe",
    "name": "<string>",
    "subbed": <bool>
}
```

User needs perms.subscribe flag or owner status to do that.

name: channel name; maximum 32 characters

subbed: true or false maps to subscribed or unsubscribed respectively

(Un)Subscribes to a channel. "post-message" (below) messages will be sent to the client by the server from other clients accordingly.

### post-message

```json
{
    "type": "post-message",
    "content": "<string>",
    "channel": "<string>",
    "author": "<string>"
}
```

User needs perms.post-message flag or owner status to do that.

content: message contents; maximum 2048 characters

channel: channel name to send the message to

author: username of the user who sent this message

Sent by client: Posts a message in a specific channel (the "author" field does not need to be sent).

Sent by server: Notifies about a sent message in a specific channel.

### get-messages (sent by client)

```json
{
    "type": "get-messages",
    "channel-id": <int>,
    "message-id": <int>,
    "amount": <int>
}
```

User needs perms.get-messages flag or owner status to do that.

### get-messages-result (sent by server)

```json
{
    "type": "get-messages-result",
    "messages": []
}
```

TODO

### list-channels

```json
{
    "type": "list-channels",
    "channels": {
        "<string>": <int>,
        ...
    }
}
```

User needs perms.list-channels flag or owner status to do that.

channels: a map in which channel name is a key and ID is a value

Sent by client: Tells the server to fetch all channels that exist (the "channels" field does not need to be sent).

Sent by server: Returns the client a channelname:id map

### list-users

```json
{
    "type": "list-users",
    "users": [
        {
            "user": {
                "id": <int>,
                "username": "<string>",
                "owner": <bool>
            },
            "online": <bool>
        },
        ...
    ]
}
```

User needs perms.list-users flag or owner status to do that.

users: an array:

User: an user.

Online: whether or not the user is connected to the server.

Sent by client: Tells the server to fetch all users that are registered (the "users" field does not need to be sent).

Sent by server: Returns the client an array of users with their status of connection.

### user-quit (sent by server)

```json
{
    "type": "user-quit",
    "username": "<string>"
}
```

username: user's name

Indicates a user disconnect (offline) event.