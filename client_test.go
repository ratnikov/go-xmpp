package xmpp

import (
  "bytes"
  "fmt"
  "os"
  "template"
  "testing"
)

type testIO struct {
  out bytes.Buffer
  in bytes.Buffer
}

func (buf *testIO) Read(p []byte) (int, os.Error) {
  return buf.in.Read(p)
}

func (buf *testIO) Write(p []byte) (int, os.Error) {
  return buf.out.Write(p)
}

// Resets both in and out buffers
func (buf *testIO) Reset() {
  buf.in.Reset()
  buf.out.Reset()
}

func (buf *testIO) PushString(str string) {
  buf.in.Write(bytes.NewBufferString(str).Bytes())
}

func (buf *testIO) PushFixture(fixture string, data interface{}) {
  buf.PushString(readFixture(fixture, data))
}

func (buf *testIO) PopString() string {
  var byte_buf []byte

  buf.out.Read(byte_buf)

  return bytes.NewBuffer(byte_buf).String()
}

func setupClient() (io *testIO, client *Client) {
  io = new(testIO)
  client = &Client{conn: io}

  return
}

type fixtureMessage struct {
  To, From, Body, Botid string
}

func TestLoop(t *testing.T) {
  testio, client := setupClient()

  testio.Reset()

  // should return from looping, if testio returns EOF
  client.Loop()

  // if we got here, then loop must have existed, and we're good
}

func TestClientOnAny(t *testing.T) {
  testio, client := setupClient()

  var received string

  client.OnAny(func(msg string) {
    received = msg
  })

  testio.PushString("<hello world! />")
  client.Loop()

  assertMatch(t, "hello world!", received, "Should return the unknown message")

  testio.PushString("<message><body>hello world!</body></message>")
  client.Loop()

  assertMatch(t, "message.*body.*hello world!.*/body.*/message", received, "Should return the message as well")
}

func TestClientOnMessage(t *testing.T) {
  testio, client := setupClient()

  testio.PushFixture("message", fixtureMessage{
    Body: "Hello world!",
    From: "joe@example.com",
    To: "sam@example.com" })

  var received *Message
  client.OnMessage(func(msg Message) {
    received = &msg
  })

  client.Loop()

  assertMatch(t, "Hello world!", received.Body(), "Should include the message")
  assertMatch(t, "sam@example.com", received.To(), "Should include To")
  assertMatch(t, "joe@example.com", received.From(), "Should include From")
}

func readFixture(filename string, data interface{}) string {
  buf := bytes.NewBufferString("")

  template.MustParseFile(fmt.Sprintf("fixture/%s.xml.gt", filename), nil).Execute(buf, data)

  return buf.String()
}
