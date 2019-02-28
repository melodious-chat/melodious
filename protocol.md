# Melodious Chat Protocol (websocket)

**Please note that this project is WIP and the protocol is subject to change.**

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

Sent by server: Indicates a user register event (no "pass" field sent).

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

Sent by server: Indicates a login (online) event (no "pass" field sent).

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

Client:
```json
{
    "type": "post-message",
    "content": "<string>",
    "channel": "<string>"
}
```
Server:
```json
{
    "type": "post-message",
    "message": {
        "content": "<string>",
        "pings": ["<string>", ...],
        "id": <int>,
        "timestamp": "string",
        "author": "<string>",
        "author_id": <int>
    },
    "channel": "<string>"
}
```

User needs perms.post-message flag or owner status to do that.

content: message contents; maximum 2048 characters

channel: channel name to send the message to or the channel it was received from

author: username of the user who sent this message

pings: usernames of people that were mentioned in the message

id: message ID

timestamp: ISO 8601 timestamp

author: username of the user who sent the message

author_id: user's ID who sent the message

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
    "messages": [{
        "content": "<string>",
        "pings": ["<string>", ...],
        "id": <int>,
        "timestamp": "<string>",
        "author": "<string>",
        "author_id": <int>
    }, ...]
}
```

Sent by client: requests messages from the server.

Sent by server: returns a list of messages.

### list-channels

```json
{
    "type": "list-channels",
    "channels": [{
        "id": <int>,
        "name": "<string>",
        "topic": "<string>"
    }, ...]
}
```

User needs perms.list-channels flag or owner status to do that.

channels: an array of channel objects

Sent by client: Tells the server to fetch all channels that exist (the "channels" field does not need to be sent).

Sent by server: Returns the client an array of channels

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

### kick (sent by client)

```json
{
    "type": "kick",
    "id": <int>, 
    "username": "<string>",
    "ban": <bool>
}
```
id: user ID

username: user's name

ban: whether or not to set the user's banned flag to true

Kicks and optionally bans a user. You MUSTN'T have both id and username fields.

The banned user will be logged off with a user-quit event.

### new-group (sent by client)

```json
{
    "type": "new-group",
    "name": "<string>"
}
```

name: group name

Creates a new group with a specified name. 

### delete-group (sent by client)

```json
{
    "type": "new-group",
    "id": <int>
}
```

id: group ID

Deletes a group with a specified ID.

### set-flag (sent by client)

```json
{
    "type": "set-flag",
    "group": "<string>",
    "name": "<string>",
    "flag": {
        "<any-json>": "<here>",
        ...
    }
}
```

group: group name

name: flag name

flag: any JSON for additional data

Sets/creates a flag for the specified group with optional additional data. Required for setting up permissions for a group.

### delete-flag (sent by client)

```json
{
    "type": "set-flag",
    "group": "<string>",
    "name": "<string>"
}
```

group: group name

name: flag name

Removes/deletes a flag from the specified group.

### typing

```json
{
    "type": "typing",
    "channel": "<string>",
    "username": "<string>"
}
```

channel: channel name

username: user's name

Sent by client: sends a typing status to a channel (the "username" field does not need to be sent).

Sent by server: indicates a typing event from someone.

### new-group-holder (sent by client)

```json
{
    "type": "new-group-holder",
    "group": "<string>",
    "user": "<string>",
    "channel": "<string>"
}
```

group: group name

user: user's name

channel: channel name

Makes a link (called a group holder) of a group and a user and/or a channel.

The logic is following:

TODO EXPLAIN HOW THIS WORKS

### delete-group-holder (sent by client)

```json
{
    "type": "delete-group-holder",
    "id": <int>
}
```

id: group holder id

Removes/deletes an already made link (group holder).

### get-group-holders

```json
{
    "type": "get-group-holders",
    "group-holders": [{
        "id": <int>,
        "group": "<string>",
        "user": "<string>",
        "channel": "<string>"
    }, ...]
}
```

Sent by client: requests a list of group holders from the server (the "group-holders" field does not need to be sent).

Sent by server: returns a list of group holders.

### ping (sent by server)

```json
{
    "type": "ping",
    "message": {
        "content": "<string>",
        "pings": ["<string>", ...],
        "id": <int>,
        "timestamp": "string",
        "author": "<string>",
        "author_id": <int>
    },
    "channel": "<string>"
}
```

message: a message object

channel: name of the channel it's coming from

Pings/mentions a mentioned user when the author sends a post-message event if message content contains a mention in the format of <@USERID>, regardless of the pinged/metioned user's subscription status.

### delete-message (sent by client)

```json
{
    "type": "delete-message",
    "id": <int>
}
```

id: message id

Deletes a message with a specified id permanently.