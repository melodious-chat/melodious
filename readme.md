# Project Melodious

Tired of Discord's staff telling you to shut the fuck up about #DiscordUnbanQuackity?
Well, Project Melodious is **the** chat for you!

#### No bullshit, no biasing, fully opensource and can be hosted by anybody!

Written in Go, it is both fast and stable. Changes being made literally every 7-12 hours, it is growing fast and strong to be **a fully-featured Discord clone, without furries!**

Although it might seem like a shit idea since Discord is boasting a userbase of >150M people and counting, we are **free of furries**, which will probably attract more people.

### Current To-Do list

* Implement a working prototype
  * Implement authentication
    * Implement `func (db *Database) RegisterUser(name string, passhash string) error`
    * Implement `func (db *Database) CheckUser(name string, passhash string) (bool, error)`
    * Get a working connection map (`*sync.Map`)
    * Implement `register` and `login` messages
  * Implement messaging
  * Allow uploading and receiving images
* Add a REST API for those who cannot use websockets

```
<handicraftsman> raov: lmao
```