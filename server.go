package main

import (
  "fmt"
  "net"
  "net/url"
)

func decode(encoded []byte) (decoded string) {
  decoded, _ = url.QueryUnescape(string(encoded))
  return
}

func main() {
  receiver, _ := net.Listen("tcp", ":9999")
  for {
    connection, _ := receiver.Accept()
    var message []byte
    fmt.Fscan(connection, &message)
    fmt.Println(decode(message))
  }
}
