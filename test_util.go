package xmpp

import (
  "testing"
  "regexp"
)

func assertEqual(t *testing.T, expected, actual interface{}, message string) {
  if expected != actual {
    t.Errorf("Expected <%s> but got <%s>", expected, actual)
  }
}

func assertMatch(t *testing.T, regex_str, str, message string) {
  regex := regexp.MustCompile(regex_str)

  if !regex.MatchString(str) {
    t.Errorf("%s: Expected <%s> to match <%s>", message, str, regex)
  } else {
    // all good
  }
}

func assertPanic(t *testing.T, message string, f func()) {
  defer func() {
    if err := recover(); err == nil {
      t.Errorf(message)
    }
  }()

  f()
}
