package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
  "strings"
  "net/url"
)

func watchConsole(input chan string) {
  console := bufio.NewReader(os.Stdin)
  for {
    text, _ := console.ReadString('\n')
    input <- strings.TrimSpace(text)
    fmt.Print("> ")
  }
}

func main() {
  input := make(chan string)
  go watchConsole(input)
  fmt.Print("What's your name? ")
  name := <-input
  for {
    receiver, _ := net.Dial("tcp", ":9999")
    message := fmt.Sprintf("%s: %s", name, <-input)
    fmt.Fprintf(receiver, "%s\n", url.QueryEscape(message))
  }
}
