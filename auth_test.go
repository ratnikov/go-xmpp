package xmpp;

import (
  "regexp"
  "testing"
)

func TestNewAuth(t *testing.T) {
  if auth := NewAuth("joe@test.com", "secret"); auth != nil {
    assertEqual(t, "joe", auth.user, "Should parse the user")
    assertEqual(t, "test.com", auth.domain, "Should parse the domain")
    assertEqual(t, "secret", auth.password, "Should assign password")
  } else {
    t.Errorf("Expected to parse %s as valid authentication info, but got nil.")
  }

  assertPanic(t, "Should not parse auth information with multiple '@'", func() {
    NewAuth("joe@johnson@test.com", "secret")
  })

  assertPanic(t, "Should complain about lack of @", func() {
    NewAuth("joe-someone.com", "secret")
  })
}

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
