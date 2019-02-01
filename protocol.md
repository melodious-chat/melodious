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
    "pass": "<sha256>"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex
hash: SHA256 hash sum of the password. MUST match `[a-f0-9]{64}` regex

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
    "pass": "<sha256>"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex
hash: SHA256 hash sum of the password. MUST match `[a-f0-9]{32-64}` regex

If user with username `name` does not exist or SHA256 hash/checksum `hash` is invalid, server MUST send a `fatal` message.

Server MAY introduce additional protection like banning users from connecting.

If user has administrator privileges, server MUST send him a corresponding `note` message.

It is recommended that server stores hash sum of the hash sum of the password to prevent heavy damage on database leak.

### new-channel (sent by client)

```json 
{
    "type": "new-channel",
    "name": "<string>"
}
```

name: channel name; maximum 32 characters

Creates a new channel. If such channel already exists or user is not an owner, server MUST return a `fail` message.