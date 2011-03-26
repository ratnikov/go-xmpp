package xmpp

type listenerList struct {
  anyCallbacks []func(message string)
  messageCallbacks []func(message Message)
  unknownCallbacks []func(unknown string)
}

func (list *listenerList) onAny(callback func(string)) {
  list.anyCallbacks = append(list.anyCallbacks, callback)
}

func (list *listenerList) onMessage(callback func(Message)) {
  list.messageCallbacks = append(list.messageCallbacks, callback)
}

func (list *listenerList) onUnknown(callback func(string)) {
  list.unknownCallbacks = append(list.unknownCallbacks, callback)
}

func (list *listenerList) fireOnMessage(msg *Message) {
  list.fireOnAny(msg.Raw())

  for i := range list.messageCallbacks {
    list.messageCallbacks[i](*msg)
  }
}

func (list *listenerList) fireOnUnknown(msg string) {
  list.fireOnAny(msg)

  for i := range list.unknownCallbacks {
    list.unknownCallbacks[i](msg)
  }
}

func (list *listenerList) fireOnAny(msg string) {
  for i := range list.anyCallbacks {
    list.anyCallbacks[i](msg)
  }
}
