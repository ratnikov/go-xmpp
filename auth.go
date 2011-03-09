package xmpp;

import (
  "fmt"
  "os"
  "regexp"
)

const (
  auth_regex = "^([^@]+)@([^@]+)$"
)

type Auth struct {
  user, domain, password string
}

func NewAuth(login, password string) (*Auth) {
  chunks := regexp.MustCompile(auth_regex).FindStringSubmatch(login)

  if len(chunks) == 0 {
    panic(os.NewError(fmt.Sprintf("Authentication identifier has to match %s", auth_regex)))
  }

  return &Auth{user: chunks[1], domain: chunks[2], password: password }
}
