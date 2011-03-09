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

  tlsconn, err := tls.Dial("tcp", "", "talk.google.com:https", nil)

  if err != nil {
    panic(err)
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
  fmt.Fprintf(client.tls, format, args...)
}

func (client *Client) authenticate(login, password string) {
  auth := NewAuth(login, password)

  client.writef("<?xml version='1.0'?><stream:stream to='%s' xmlns:stream='http://etherx.jabber.org/streams' xmlns='jabber:client' version='1.0' />", auth.domain)

  features := client.readUntilEOF()

  client.log("Gotten folloiwing features: %s", features)
}

func (client *Client) readUntilEOF() string {
  var buf bytes.Buffer

  L: for {
    b, err := client.tls_read.ReadByte()

    switch (err) {
      case nil: buf.WriteByte(b)
      case os.EOF: break L
      default: panic(os.NewError(fmt.Sprintf("Failed to read from server: %s\n", err)))
    }
  }

  return buf.String()
}

func (client *Client) log(format string, args ...interface{}) {
  fmt.Printf("LOG: %s\n", fmt.Sprintf(format, args... ))
}
