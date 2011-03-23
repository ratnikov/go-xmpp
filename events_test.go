package xmpp

import "testing"

type testListener struct {
  message string
}

func (l *testListener) onMessage(msg string) {
  l.message = msg
}

func TestFireOnMessage(t *testing.T) {
  list := listenerList{}

  l1 := testListener{}
  l2 := testListener{}

  list.subscribe(&l1)

  assertEqual(t, "", l1.message, "Should not fire any events")

  list.fireOnMessage("one")

  assertEqual(t, "one", l1.message, "Should forward the fired message")

  list.subscribe(&l2)

  list.fireOnMessage("two")

  assertEqual(t, "two", l1.message, "Should forward the second message as well")
  assertEqual(t, "two", l2.message, "Should fire for second listener as well")
}
