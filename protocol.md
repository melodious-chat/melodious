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

### register (sent by client)

```json
{
    "type": "register",
    "name": "<string>",
    "pass": "<string>"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex
pass: Password. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

If user with username `name` already exists, server MUST send a `fatal` message.

Server MAY introduce additional protection like IP duplication protection.

If this is a first user ever registered, server MUST give that user admin permissions and send him a corresponding `note` message.

It is recommended that server stores hash sum of the hash sum of the password to prevent heavy damage on database leak.

After registering user MUST be treated as logged in.

### login (sent by client)

```json
{
    "type": "login",
    "name": "<string>",
    "pass": "<string>"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex
pass: Password. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

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

name: channel name; maximum 32 characters

Sent by client: Deletes a channel.

Sent by server: Notifies about a deleted channel. May mention non-existing channel.

### subscribe (sent by client)

```json
{
    "type": "subscribe",
    "name": "<string>",
    "subbed": <bool>
}
```

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

content: message contents; maximum 2048 characters

channel: channel name to send the message to

author: username of the user who sent this message

Sent by client: Posts a message in a specific channel (the "author" field does not need to be sent).

Sent by server: Notifies about a sent message in a specific channel.