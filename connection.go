package xmpp;

import (
  "bytes"
  "crypto/tls"
  "io"
  "fmt"
  "os"
  "net"
  "xml"
)

type Client struct {
  tls *tls.Conn
  tls_read io.ByteReader
  parser *xml.Parser
  err os.Error
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

  c := &Client{}


  hostname := "talk.google.com"

  if conn, err := net.Dial("tcp", "", "talk.google.com:5222"); err != nil {
    fmt.Printf("crap!!! %s\n", err)
  } else {
    c.write(conn, "<?xml version='1.0'?>")
    c.write(conn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

    fmt.Printf("Read: %s\n", client.read(conn))

    // assuming need to start tls
    c.write(conn, "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls' />")

    fmt.Printf("Read: %s\n", client.read(conn))
  }

  tlsconn, err := tls.Dial("tcp", "", hostname + ":https", nil)

  if err != nil {
    panic(os.NewError(fmt.Sprintf("Failed to connect to %s (%s)", hostname, err)))
  }

  tlsconn.Handshake()

  client.write(tlsconn, "<stream:stream to='gmail.com' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>")

  client.read(tlsconn)
  client.read(tlsconn)

  auth := NewAuth(user, password)

  client.write(tlsconn, "<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>%s</auth>", auth.Base64())

  client.read(tlsconn)

  return client, nil
}

func (client *Client) read(conn io.Reader) (out string) {
  buf := make([]byte, 2048)

  num, err := conn.Read(buf)

  if err != nil {
    panic("Failed to read....")
  } else {
    fmt.Printf("Read %d bytes: %s\n", num, buf)
  }

  return
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
