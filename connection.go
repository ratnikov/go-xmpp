package xmpp;

import (
  "crypto/tls"
  "fmt"
  "os"
  "regexp"
  "net"
)

type Client struct {
  hostname string
  conn net.Conn
}

func (client *Client) startTls() {
  var plain_conn net.Conn
  var err os.Error

  if plain_conn, err = net.Dial("tcp", "", client.hostname + ":5222"); err != nil {
    die("Failed to establish plain connection: %s", err)
  }

  write(plain_conn, "<?xml version='1.0'?>")
  write(plain_conn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  log("Read: %s", read(plain_conn))

  // assuming need to start tls
  write(plain_conn, "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls' />")
  log("Read: %s", read(plain_conn))

  // assuming the server asked to proceed
  if client.conn, err = tls.Dial("tcp", "", client.hostname + ":https", nil); err != nil {
    die("Failed to establish tls connection (%s)", err)
  }
}


func NewClient(hostname, user, password string) (client *Client, failure os.Error) {
  failure = nil
  defer func() {
    if err := recover(); err != nil {
      client = nil

      var ok bool
      if failure, ok = err.(os.Error); !ok {
        failure = os.NewError(fmt.Sprintf("Weird error happened: %s\n", failure))
      }
    }
  }()

  client = &Client{hostname: hostname}

  client.startTls()

  client.write("<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  client.read()
  client.read()

  auth := NewAuth(user, password)

  client.write("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>%s</auth>", auth.Base64())

  client.read()

  client.write("<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  client.read()

  client.write("<iq type='set' id='xmpp-bot1029'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><resource>Home</resource></bind></iq>")

  client.read()
  client.read()
  client.read()

  client.Message("ratnikov@gmail.com", "Write me something and I will write back! (Please send 2 messages at first....)")

  client.read()

  for {
    raw := client.read()

    var recipient, received string

    if match := regexp.MustCompile("to=\"([^/\"]*)/").FindStringSubmatch(raw); len(match) > 1 {
      recipient = match[1]
    }

    if match := regexp.MustCompile("<body>(.*)</body>").FindStringSubmatch(raw); len(match) > 1 {
      received = match[1]
    }

    message := "You wrote: " + received

    fmt.Printf("Sending message %s to %s", message, recipient)

    client.Message(recipient, message)
  }

  return client, nil
}

func (client *Client) read() string {
  return read(client.conn)
}

func (client *Client) write(format string, args ...interface{}) int {
  return write(client.conn, format, args...)
}

func (client *Client) Message(recipient, msg string) {
  client.write("<message type='chat' id='xmpp-bot1029' to='%s'><body>%s</body></message>", recipient, msg)
}
