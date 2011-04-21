package main

import (
  "fmt"
  "xmpp"
)

func main() {
  client, err := xmpp.NewClient("talk.google.com", "xmpp.chatterbox@gmail.com", "XXX")

  if err != nil {
    fmt.Printf("Failed due to: %s\n", err)
  } else {
    client.OnAny(func(msg string) {
      log(msg)
    })

    client.OnMessage(func(msg xmpp.Message) {
      log("Got a message from %s to %s: %s\n", msg.From(), msg.To(), msg.Body())
    })

    client.SendChat("ratnikov@gmail.com", "Hello world!")

    client.Loop()
  }
}

func log(format string, args ...interface{}) {
  fmt.Printf("MAIN LOG: " + format + "\n", args...)
}
