package main

import (
  "fmt"
  "xmpp"
)

type Bot struct {
}

func (b *Bot) onMessage(message string) {
  log("Bot received message: %s", message)
}

func (b *Bot) onUnkown(msg string) {
  log("Bot received unsupported message: %s", msg)
}

func main() {
  bot := Bot{}

  client, err := xmpp.NewClient("talk.google.com", "xmpp.chatterbox@gmail.com", "XXX")

  if err != nil {
    fmt.Printf("Failed due to: %s\n", err)
  } else {
    client.Subscribe(&bot)

    client.Message("ratnikov@gmail.com", "Hello world!")

    client.Loop()
  }
}

func log(format string, args ...interface{}) {
  fmt.Printf("MAIN LOG: " + format, args...)
}
