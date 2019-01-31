# Melodious

Tired of Discord's staff telling you to shut the fuck up about #DiscordUnbanQuackity?
Well, Project Melodious is **the** chat for you!

#### No bullshit, no biasing, fully opensource and can be hosted by anybody!

Written in Go, it is both fast and stable. Changes being made literally every 7-12 hours, it is growing fast and strong to be **a fully-featured Discord clone, without furries!**

Although it might seem like a shit idea since Discord is boasting a userbase of >150M people and counting, we are **free of furries**, which will probably attract more people.

## Now, seriously

Melodious is an attempt at writing a self-hosted Discord alternative initiated by recent Discord staff actions mentioned above.

Its server is being written in Go, while web client will be written in JavaScript.

### Building

```bash
$ git clone https://gitlab.com/melodious-chat/melodious.git
$ cd melodious
$ go build
```

### Configuring

Melodious server requires a running copy of PostgreSQL.

An example config file:

```json
{
    "db-addr": "postgres://handicraftsman:password@localhost/melodious?sslmode=disable",
    "http-addr": "0.0.0.0:8080",
    "delete-history-every": "PT1H",
    "store-history-for": "P1W"
}
```

`delete-history-every` and `store-history-for` are ISO 8601 duration string.

### Starting

```bash
/path/to/melodious -config /path/to/melodious.config.json
```

## Current To-Do list

* Implement a working prototype
  * Implement authentication
    * Get a working connection map (`*sync.Map`)
    * Implement `register` and `login` messages
  * Implement messaging
  * Allow uploading and receiving images
* Add a REST API for those who cannot use websockets
