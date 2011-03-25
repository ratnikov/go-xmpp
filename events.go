package xmpp

type listenerList struct {
  messageCallbacks []func(message string)
  unknownCallbacks []func(unknown string)
}

func (list *listenerList) onMessage(callback func(string)) {
  list.messageCallbacks = append(list.messageCallbacks, callback)
}

func (list *listenerList) onUnknown(callback func(string)) {
  list.unknownCallbacks = append(list.unknownCallbacks, callback)
}

func (list *listenerList) fireOnMessage(msg string) {
  for i := range list.messageCallbacks {
    list.messageCallbacks[i](msg)
  }
}

func (list *listenerList) fireOnUnknown(msg string) {
  for i := range list.unknownCallbacks {
    list.unknownCallbacks[i](msg)
  }
}
