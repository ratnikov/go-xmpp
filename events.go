package xmpp

type listenerList struct {
  onMessageListeners []func(message string)
  onUnknownListeners []func(unknown string)
}

func (list *listenerList) onMessage(callback func(string)) {
  list.onMessageListeners = append(list.onMessageListeners, callback)
}

func (list *listenerList) onUnknown(callback func(string)) {
  list.onUnknownListeners = append(list.onUnknownListeners, callback)
}

func (list *listenerList) fireOnMessage(msg string) {
  for i := range list.onMessageListeners {
    list.onMessageListeners[i](msg)
  }
}

func (list *listenerList) fireOnUnknown(msg string) {
  for i := range list.onUnknownListeners {
    list.onUnknownListeners[i](msg)
  }
}
