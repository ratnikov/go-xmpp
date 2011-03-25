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

func (buf *testIO) PushFixture(fixture string, data interface{}) {
  str := readFixture(fixture, data)

  var byt []byte

  byt = bytes.NewBufferString(str).Bytes()

  buf.in.Write(byt)
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

func readFixture(filename string, data interface{}) string {
  buf := bytes.NewBufferString("")

  template.MustParseFile(fmt.Sprintf("fixture/%s.xml.gt", filename), nil).Execute(buf, data)

  return buf.String()
}
