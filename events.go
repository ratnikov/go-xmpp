package xmpp

type Listener interface {
}

type MessageListener interface {
  onMessage(message string)
}

type UnknownListener interface {
  onUnknown(msg string)
}

type AllListener interface {
  MessageListener
  UnknownListener
}

type listenerList struct {
  listeners []Listener
}

func (list *listenerList) subscribe(listener Listener) {
  list.listeners = append(list.listeners, listener)
}

func (list *listenerList) fireOnMessage(msg string) {
  list.each(func(raw_l Listener) {
    if l, ok := raw_l.(MessageListener); ok {
      l.onMessage(msg)
    } else {
      // not our listener, oh well...
    }
  })
}

func (list *listenerList) fireOnUnknown(msg string) {
  list.each(func(raw_l Listener) {
    if l, ok := raw_l.(UnknownListener); ok { 
      l.onUnknown(msg)
    } else {
      // not our listener, oh well...
    }
  })
}

func (list *listenerList) each(f func(listener Listener)) {
  for l := range list.listeners {
    f(list.listeners[l])
  }
}
