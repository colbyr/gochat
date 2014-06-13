package main

import (
  "bufio"
  "fmt"
  "net"
  "net/url"
  "os"
  "strings"
)

func decode(encoded []byte) (decoded string) {
  decoded, _ = url.QueryUnescape(string(encoded))
  return
}

func log(text string) {
  fmt.Printf("%s\n> ", text)
}

func watchConsole(input chan string) {
  console := bufio.NewReader(os.Stdin)
  for {
    text, _ := console.ReadString('\n')
    if (strings.TrimSpace(text) != "") {
      input <- strings.TrimSpace(text)
    } else {
      fmt.Print("> ")
    }
  }
}

func listenForMessages(address string, messages chan string) {
  for {
    receiver, _ := net.Listen("tcp", address)
    for {
      connection, err := receiver.Accept()
      if err != nil {
        fmt.Printf("Error: %s\n", err)
        continue
      }
      var command, name, message []byte
      fmt.Fscan(connection, &command, &name, &message)
      switch string(command) {
      case "message":
        messages <- fmt.Sprintf("(%s) %s", string(name), decode(message))
      }
    }
  }
}

func sendCommand(command, name, text string) {
  message := fmt.Sprintf("%s %s %s ", command, name, url.QueryEscape(text))
  connection, _ := net.Dial("tcp", ":9999")
  fmt.Fprint(connection, message)
  connection.Close()
  fmt.Print("> ")
}

func registerAddress(name string) (address string) {
  connection, _ := net.Dial("tcp", ":9999")
  address = connection.LocalAddr().String()
  message := fmt.Sprintf("register %s %s ", name, url.QueryEscape(address))
  fmt.Fprint(connection, message)
  connection.Close()
  fmt.Print("> ")
  return
}

func main() {
  input := make(chan string)
  messages := make(chan string)

  go watchConsole(input)

  fmt.Print("What's your name? ")
  name := <-input

  address := registerAddress(name)
  log("Connected to server!")
  go listenForMessages(address, messages)

  for {
    select {
    case text := <-input:
      sendCommand("send", name, text)
    case text := <-messages:
      log(text);
    }
  }
}
