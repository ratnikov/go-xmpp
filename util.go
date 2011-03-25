package xmpp


import (
  "bytes"
  "fmt"
  "io"
  "os"
)

func read(in io.Reader) (string, os.Error) {
  buf := make([]byte, 2048)

  num, err := in.Read(buf)

  if err != nil {
    fmt.Printf("Error occured (%s). Read %d bytes anyway: %s\n", err, num, buf)
  } else {
    fmt.Printf("Read %d bytes: %s\n", num, buf)
  }

  return bytes.NewBuffer(buf).String(), err
}

func mustRead(in io.Reader) (out string) {
  var err os.Error
  if out, err = read(in); err != nil {
    die("Failed to read (%s)", err)
  }

  return
}

func write(out io.Writer, format string, args ...interface{}) int {
  buf := bytes.NewBufferString(fmt.Sprintf(format + "\n", args...))

  n, err := out.Write(buf.Bytes())

  if err != nil {
    die("Failed to write to %s (%s)", out, err)
  }

  return n
}

func die(format string, args ...interface{}) {
  panic(os.NewError(fmt.Sprintf(format, args...)))
}

func log(format string, args ...interface{}) {
  fmt.Printf("LOG: %s\n", fmt.Sprintf(format, args... ))
}
