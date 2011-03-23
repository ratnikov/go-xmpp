package xmpp

type Listener interface {
  onMessage(message string)
}

type listenerList struct {
  listeners []Listener
}

func (list *listenerList) subscribe(listener Listener) {
  list.listeners = append(list.listeners, listener)
}

func (list *listenerList) fireOnMessage(msg string) {
  list.each(func(listener Listener) {
    listener.onMessage(msg)
  })
}

func (list *listenerList) each(f func(listener Listener)) {
  for l := range list.listeners {
    f(list.listeners[l])
  }
}
