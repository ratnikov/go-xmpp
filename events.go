package xmpp

type Listener interface {
}

type MessageListener struct {
  callback func(message Message)
}

type UnknownListener struct {
  callback func(raw string)
}

type AnyListener struct {
  callback func(raw string)
}

type listenerList struct {
  listeners []Listener
}

func (list *listenerList) Subscribe(listener Listener) {
  list.listeners = append(list.listeners, listener)
}

func (list *listenerList) onAny(callback func(string)) {
  list.Subscribe(AnyListener{ callback : callback })
}

func (list *listenerList) onMessage(callback func(Message)) {
  list.Subscribe(MessageListener{ callback: callback })
}

func (list *listenerList) onUnknown(callback func(string)) {
  list.Subscribe(UnknownListener{ callback : callback })
}

func (list *listenerList) fireOnMessage(msg *Message) {
  list.eachListener(func(l Listener) {
    fireIfAny(l, msg.Raw())

    if msg_l, ok := l.(MessageListener); ok {
      msg_l.callback(*msg)
    }
  })
}

func (list *listenerList) fireOnUnknown(msg string) {
  list.eachListener(func(l Listener) {
    fireIfAny(l, msg)

    if unknown_l, ok := l.(UnknownListener); ok {
      unknown_l.callback(msg)
    }
  })
}

func fireIfAny(l Listener, msg string) {
  if any_l, ok := l.(AnyListener); ok {
    any_l.callback(msg)
  } else {
    // not any listener, so nothing to do
  }
}

func (list *listenerList) eachListener(callback func(Listener)) {
  for i := range list.listeners {
    callback(list.listeners[i])
  }
}
