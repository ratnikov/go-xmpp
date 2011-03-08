package xmpp;

import (
  "testing"
)

func TestNewAuth(t *testing.T) {
  if auth, err := NewAuth("joe@test.com", "secret"); auth != nil {
    assertEqual(t, "joe", auth.user, "Should parse the user")
    assertEqual(t, "test.com", auth.domain, "Should parse the domain")
    assertEqual(t, "secret", auth.password, "Should assign password")
  } else {
    t.Errorf("Expected to parse %s as valid authentication info, but got nil. Accompanying error: %s", "joe@test.com", err)
  }

  if auth, err := NewAuth("joe@johnson@test.com", "secret"); auth != nil || err == nil {
    t.Error("Should not parse auth information with multiple '@'")
  }

  if auth, err := NewAuth("joe-someone.com", "secret"); auth != nil || err == nil {
    t.Error("Should complain about lack of @")
  }
}

func assertEqual(t *testing.T, expected, actual interface{}, message string) {
  if expected != actual {
    t.Errorf("Expected <%s> but got <%s>", expected, actual)
  }
}
