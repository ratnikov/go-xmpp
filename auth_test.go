package xmpp;

import (
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
