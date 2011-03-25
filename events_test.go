package xmpp

import "testing"

func TestOnMessage(t *testing.T) {
  list := listenerList{}

  var msg1, msg2 string

  list.onMessage(func(msg string) {
    msg1 = msg
  })

  assertEqual(t, "", msg1, "Should not fire any events")

  list.fireOnMessage("one")

  assertEqual(t, "one", msg1, "Should forward the fired message")

  list.onMessage(func(msg string) {
    msg2 = msg
  })

  list.fireOnMessage("two")

  assertEqual(t, "two", msg1, "Should forward the second message as well")
  assertEqual(t, "two", msg2, "Should fire for second listener as well")
}

func TestOnUnknown(t *testing.T) {

  list := listenerList{}

  msg1 := ""

  list.onUnknown(func(msg string) {
    msg1 = msg
  })

  list.fireOnUnknown("one")

  assertEqual(t, "one", msg1, "Should have fired onUnknown")
}
