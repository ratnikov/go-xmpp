package xmpp;

import (
  "bufio"
  "bytes"
  "crypto/tls"
  "io"
  "fmt"
  "os"
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

  hostname := "talk.google.com:https"
  tlsconn, err := tls.Dial("tcp", "", hostname, nil)

  if err != nil {
    panic(os.NewError(fmt.Sprintf("Failed to connect to %s (%s)", hostname, err)))
  }

  client = &Client{
    tls: tlsconn,
    tls_read: bufio.NewReader(tlsconn),
    parser: xml.NewParser(tlsconn),
  }

  client.authenticate(user, password)

  return
}

func (client *Client) writef(format string, args ...interface{}) {
  buf := bytes.NewBufferString(fmt.Sprintf(format + "\n", args...))

  client.log("Sending %s", buf.String())

  if num, err := client.tls.Write(buf.Bytes()); err != nil {
    panic(os.NewError(fmt.Sprintf("Failed to write to server (%s)", err)))
  } else {
    client.log("Wrote %d bytes to server", num)
  }
}

func (client *Client) authenticate(login, password string) {
  auth := NewAuth(login, password)

  client.writef("<?xml version='1.0'?>")
  client.writef("<stream:stream to='%s' xmlns:stream='http://etherx.jabber.org/streams' xmlns='jabber:client' version='1.0' />", auth.domain)

  features := client.readUntilEOF()

  client.log("Gotten following features: %s", features)

  client.writef("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN' xmlns:ga='http://www.google.com/talk/protocol/auth' ga:client-uses-full-bind-result='true'>AGp1bGlldAByMG0zMG15cjBtMzA=</auth>")

  //client.writef("<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>%s</auth>\n", auth.Base64())

  auth_response := client.readUntilEOF()

  client.log("Auth response: %s\n", auth_response)
}

func (client *Client) readUntilEOF() string {
  var buf bytes.Buffer

  L: for {
    b, err := client.tls_read.ReadByte()

    switch (err) {
      case nil: buf.WriteByte(b)
      case os.EOF:
        if buf.Len() > 0 {
          break L
        }
      default: panic(os.NewError(fmt.Sprintf("Failed to read from server: %s\n", err)))
    }
  }

  return buf.String()
}

func (client *Client) log(format string, args ...interface{}) {
  fmt.Printf("LOG: %s\n", fmt.Sprintf(format, args... ))
}
