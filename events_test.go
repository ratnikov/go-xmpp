package xmpp

import "testing"

type testMessageListener struct {
  message string
}

type testUnknownListener struct {
  unknown string
}

func (l *testMessageListener) onMessage(msg string) {
  l.message = msg
}

func (l* testUnknownListener) onUnknown(msg string) {
  l.unknown = msg
}

func TestFireOnMessage(t *testing.T) {
  list := listenerList{}

  l1 := testMessageListener{}
  l2 := testMessageListener{}

  list.subscribe(&l1)

  assertEqual(t, "", l1.message, "Should not fire any events")

  list.fireOnMessage("one")

  assertEqual(t, "one", l1.message, "Should forward the fired message")

  list.subscribe(&l2)

  list.fireOnMessage("two")

  assertEqual(t, "two", l1.message, "Should forward the second message as well")
  assertEqual(t, "two", l2.message, "Should fire for second listener as well")
}

func TestFireOnUnknown(t *testing.T) {

  list := listenerList{}

  l := testUnknownListener{}

  list.subscribe(&l)

  list.fireOnUnknown("one")

  assertEqual(t, "one", l.unknown, "Should have fired onUnknown")
}
