package xmpp

import (
  "bytes"
  "fmt"
  "os"
  "regexp"
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
  Type, To, From, Body, Botid string
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
    Type: "any-possible-type",
    Body: "Some stuff",
    From: "joe",
    To:   "sam" })

  var received *Message

  client.OnMessage(func(msg Message) {
    received = &msg
  })

  client.Loop()

  should(t, "invoked callback", func() bool {
    return received != nil
  })

  assertMatch(t, "any-possible-type", received.Type(), "Should capture its type")
}

func TestClientOnChatMessage(t *testing.T) {
  testio, client := setupClient()

  testio.PushFixture("message", fixtureMessage{
    Type: "chat",
    Body: "Hello world!",
    From: "joe@example.com",
    To: "sam@example.com" })

  var received *Message

  client.OnChatMessage(func(msg Message) {
    received = &msg
  })

  client.Loop()

  should(t, "have invoked onMessage callback", func() bool {
    return received != nil
  })

  assertMatch(t, "chat", received.Type(), "Should be a chat message")
  assertMatch(t, "Hello world!", received.Body(), "Should include the message")
  assertMatch(t, "sam@example.com", received.To(), "Should include To")
  assertMatch(t, "joe@example.com", received.From(), "Should include From")
}

func TestClientOnErrorMessage(t *testing.T) {
  testio, client := setupClient()

  testio.PushFixture("message", fixtureMessage{
    Type: "error",
    Body: "Something went wrong",
    From: "joe",
    To: "sam" })

  testio.PushFixture("message", fixtureMessage{
    Type: "chat",
    Body: "Some message" })

  var received *Message
  client.OnErrorMessage(func(msg Message) {
    received = &msg
  })

  client.Loop()

  should(t, "have invoked the callback", func() bool {
    return received != nil
  })

  assertMatch(t, "went.*wrong", received.Body(), "Should pick up the error only")
}

func TestClientOnUnknown(t *testing.T) {
  testio, client := setupClient()

  var received string
  client.OnUnknown(func(msg string) {
    received = msg
  })

  testio.PushString("<foobar>hello world!</foobar>")
  client.Loop()

  should(t, "invoke the unknown listener", func() bool {
    match, _ := regexp.MatchString("<foobar>hello world!</foobar>", received)

    return match
  })
}

func readFixture(filename string, data interface{}) string {
  buf := bytes.NewBufferString("")

  template.MustParseFile(fmt.Sprintf("fixture/%s.xml.gt", filename), nil).Execute(buf, data)

  return buf.String()
}
