package main

import (
  "fmt"
  "net"
  "net/url"
  "time"
)

var clients = make(map[string]string)

func decode(encoded []byte) (decoded string) {
  decoded, _ = url.QueryUnescape(string(encoded))
  return
}

func distribute(sender, message string) {
  for to, _ := range clients {
    if to == sender {
      continue
    }
    go forwardMessage("message", to, sender, message)
  }
}

func forwardMessage(command, to, sender, text string) {
  message := fmt.Sprintf("%s %s %s ", command, sender, url.QueryEscape(text))
  connection, err := net.Dial("tcp", clients[to])
  if err != nil {
    delete(clients, to)
    distribute(to, "DISCONNECTED")
    return
  }
  fmt.Fprint(connection, message)
  connection.Close()
}

func ping(name string) {
  connection, err := net.Dial("tcp", clients[name])
  if err != nil {
    go distribute(name, "DISCONNECTED")
    delete(clients, name)
    return
  }
  connection.Close()
}

func startPinging() {
  for {
    time.Sleep(1000 * time.Millisecond)
    for name, _ := range clients {
      go forwardMessage("ping", name, "SERVER", "")
    }
  }
}

func handleCommand(command, name, message string) {
  switch command {
  case "register":
    clients[name] = message
    fmt.Printf("%s registered at %s\n", name, message)
    distribute(name, "CONNECTED")
  case "send":
    fmt.Printf("%s: %s\n", name, message)
    distribute(name, message)
  }
}

func main() {
  receiver, _ := net.Listen("tcp", ":9999")
  go startPinging()
  for {
    connection, err := receiver.Accept()
    if err != nil {
      fmt.Printf("Error: %s\n", err)
      continue
    }
    var command, name, message []byte
    fmt.Fscan(connection, &command, &name, &message)
    handleCommand(string(command), string(name), decode(message))
  }
}
