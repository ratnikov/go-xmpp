package xmpp

import "testing"

func TestBadXmppMessage(t *testing.T) {
  should(t, "create a nil message", func() bool {
    return NewMessage("<message>no good info</message>") == nil
  })
}

func TestXmppMessage(t *testing.T) {
  msg := NewMessage("<message from=\"me\" to=\"you\"><body>Hello world</body></message>")

  should(t, "parse 'from'", func() bool {
    return msg.From() == "me"
  })

  should(t, "parse 'to'", func() bool {
    return msg.To() == "you"
  })

  should(t, "parse body", func() bool { return msg.Body() == "Hello world" })
}
