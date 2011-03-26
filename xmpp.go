package xmpp

import (
  "os"
  "regexp"
)

type Message struct {
  raw, to, from, body string
}

func NewMessage(raw string) *Message {
  var err os.Error
  msg := new(Message)

  msg.raw = raw

  if msg.to, err = parse(raw, "to=\"([^\"]+)\""); err != nil {
    return nil // failed to parse to
  }

  if msg.from, err = parse(raw, "from=\"([^\"]+)\""); err != nil {
    return nil // failed to parse from
  }

  if msg.body, err = parse(raw, "<body>(.*)</body>"); err != nil {
    return nil // failed to parse body
  }

  return msg
}

func (m *Message) Raw() string {
  return m.raw
}

func (m *Message) From() (from string) {
  return m.from
}

func (m *Message) To() string {
  return m.to
}

func (m *Message) Body() string {
  return m.body
}

func parse(raw, regex string) (out string, err os.Error) {
  if match := regexp.MustCompile(regex).FindStringSubmatch(raw); len(match) > 0 {
    out = match[1]
  } else {
    err = os.NewError("Failed to parse message")
  }

  return
}
