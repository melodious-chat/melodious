# Melodious Chat Protocol (websocket)

Communication is established over a WebSocket _connection_ with JSON _messages_ used to exchange data.

## Messages

Each message is a JSON document which necessarily has `type` field which describes message type

### fatal (sent by server)

```json
{
    "type": "fatal",
    "message": "string",
}
```

Send by server to indicate that a fatal error has occured and the connection must be closed.

Server MUST close the connection after sending this message.

### register (sent by client)

```json
{
    "type": "register",
    "name": "string",
    "pass": "string"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex
pass: Password. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

If user with username `name` already exists, server MUST send a `fatal` message.

Server MAY introduce additional protection like IP duplication protection.

### login (sent by client)

```json
{
    "type": "login",
    "name": "string",
    "pass": "string"
}
```

name: Username. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex
pass: Password. MUST match `[a-zA-Z0-9\-_\.]{3,32}` regex

If user with username `name` does not exist or password `pass` is invalid, server MUST send a `fatal` message.

Server MAY introduce additional protection like banning users from connecting.

### ping (sent by server)

```json
{
    "type": "ping"
}
```

Server MUST send this every N seconds. If user does not reply with a `pong` message in N seconds, server MUST send a `fatal` message.

### pong (sent by client)

```json
{
    "type": "pong"
}
```

Client MUST send this message in response to a `ping` message. Clients which fail to do so will be disconnected when next `ping` message should be sent.