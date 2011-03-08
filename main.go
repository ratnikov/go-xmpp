package main

import (
  "fmt"
  "xmpp"
)

func main() {
  client, err := xmpp.NewClient("ratnikov@gmail.com", "secret")

  if err != nil {
    fmt.Printf("Failed due to: %s\n", err)
  }

  fmt.Printf("Hello world!\n")

  fmt.Printf("Client: %x\n", client)
}
