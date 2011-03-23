package xmpp;

import (
  "bytes"
  "crypto/tls"
  "io"
  "fmt"
  "os"
  "regexp"
  "net"
  "xml"
)

type Client struct {
  tls *tls.Conn
  tls_read io.ByteReader
  parser *xml.Parser
  err os.Error
}

func startTls(host string) {
  if conn, err := net.Dial("tcp", "", "talk.google.com:5222"); err != nil {
    die("Failed to establish plain connection: %s", err)
  } else {
    write(conn, "<?xml version='1.0'?>")
    write(conn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

    fmt.Printf("Read: %s\n", read(conn))

    // assuming need to start tls
    write(conn, "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls' />")

    fmt.Printf("Read: %s\n", read(conn))
  }
}


func NewClient(user, password string) (client *Client, failure os.Error) {
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

  client = &Client{}


  
  hostname := "talk.google.com"

  if conn, err := net.Dial("tcp", "", "talk.google.com:5222"); err != nil {
    die("Failed to establish plain connection: %s", err)
  } else {
    client.write(conn, "<?xml version='1.0'?>")
    client.write(conn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

    fmt.Printf("Read: %s\n", client.read(conn))

    // assuming need to start tls
    client.write(conn, "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls' />")

    fmt.Printf("Read: %s\n", client.read(conn))
  }

  tlsconn, err := tls.Dial("tcp", "", hostname + ":https", nil)

  if err != nil {
    panic(os.NewError(fmt.Sprintf("Failed to connect to %s (%s)", hostname, err)))
  }

  client.tls = tlsconn

  tlsconn.Handshake()

  client.write(tlsconn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  client.read(tlsconn)
  client.read(tlsconn)

  auth := NewAuth(user, password)

  client.write(tlsconn, "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>%s</auth>", auth.Base64())

  client.read(tlsconn)

  client.write(tlsconn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  client.read(tlsconn)

  client.write(tlsconn, "<iq type='set' id='xmpp-bot1029'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><resource>Home</resource></bind></iq>")

  client.read(tlsconn)
  client.read(tlsconn)
  client.read(tlsconn)

  client.Message("ratnikov@gmail.com", "Write me something and I will write back! (Please send 2 messages at first....)")

  client.read(tlsconn)

  for {
    raw := client.read(tlsconn)

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

func (client *Client) Message(recipient, msg string) {
  client.write(client.tls, fmt.Sprintf("<message type='chat' id='xmpp-bot1029' to='%s'><body>%s</body></message>", recipient, msg))
}

func (client *Client) read(conn io.Reader) (out string) {
  buf := make([]byte, 2048)

  num, err := conn.Read(buf)

  if err != nil {
    fmt.Printf("Error occured (%s). Read %d bytes anyway: %s\n", err, num, buf)
  } else {
    fmt.Printf("Read %d bytes: %s\n", num, buf)
  }

  return bytes.NewBuffer(buf).String()
}


func (client *Client) write(conn io.Writer, format string, args ...interface{}) {
  buf := bytes.NewBufferString(fmt.Sprintf(format + "\n", args...))

  client.log("Sending %s\n", buf.String())

  if num, err := conn.Write(buf.Bytes()); err != nil {
    panic(os.NewError(fmt.Sprintf("Failed to write... (%s)", err)))
  } else {
    client.log("Wrote %d bytes to connection", num)
  }
}

func (client *Client) log(format string, args ...interface{}) {
  fmt.Printf("LOG: %s\n", fmt.Sprintf(format, args... ))
}
