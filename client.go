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
  listeners listenerList
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
  client.authenticate(NewAuth(user, password))

  return client, nil
}

func (client *Client) OnMessage(callback func(string)) {
  client.listeners.onMessage(callback);
}

func (client *Client) Loop() {
  for {
    msg := client.read()

    switch {
      case regexp.MustCompile("^<message").MatchString(msg): client.listeners.fireOnMessage(msg)
      default: client.listeners.fireOnUnknown(msg)
    }
  }
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

func (client *Client) authenticate(auth *Auth) {
  client.write("<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  // get stream response with id back
  client.read()

  // get auth mechanisms...
  client.read()

  // assuming we can do plain authentication
  client.write("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>%s</auth>", auth.Base64())

  // get "success" response
  client.read()

  // re-start the stream
  client.write("<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  client.read() // get stream acknowledgement
  client.read() // get session information

  // identify as xmpp-bot1029
  client.write("<iq type='set' id='xmpp-bot1029'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><resource>Home</resource></bind></iq>")
  client.read() // get return as to what we're bound to... or something...

  // anyhow, assuming authentication is complete
}

func (client *Client) read() string {
  msg := read(client.conn)
  log(">> " + msg)
  return msg
}

func (client *Client) write(format string, args ...interface{}) int {
  log("<< " + format, args...)
  return write(client.conn, format, args...)
}

func (client *Client) Message(recipient, msg string) {
  client.write("<message type='chat' id='xmpp-bot1029' to='%s'><body>%s</body></message>", recipient, msg)
}