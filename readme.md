# PROJECT FROZEN; Contributors needed

I have too much going on, so the project is currently in the frozen state. I have little to no time on developing Melodious because I'm not going to sacrifice my time for this project due to having a full-time job.  
I (or we) may get back on developing Melodious by January 1st 2020; but for now, if you want to contribute, please let me know. 
*- Polyzium, December 22nd 2019*  

UPDATE: I'm no longer commited to its development due to Fosscord's existence. I recommend you use that instead when it's ready.  
*Updated on July 21st 2022*

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
    "delete-history-every": "1h",
    "store-history-for": "P1W"
}
```

`store-history-for` is a ISO 8601 duration string.

`delete-history-every` is a duration string as described at https://golang.org/pkg/time/#ParseDuration

### Starting

```bash
/path/to/melodious -config /path/to/melodious.config.json
```

## Current To-Do list

* Implement misc features
  * VoIP
  * editable messages
  * custom user status
* Add a REST API for those who cannot use websockets
