package main

import (
	"fmt"
  "log"
	"net"
  "net/http"
	"os"
	"io"
  "io/ioutil"

	"golang.org/x/net/proxy"
	"github.com/joho/godotenv"
)

// Global env initialization
func init() {
	err := godotenv.Load(".env")
	if err != nil {
     log.Fatalf("Error loading .env file: %v", err)
   }
} 


// Create a requst using proxy
func ProxyRequest(resource string) (string, error) {
  auth := proxy.Auth{
    User: os.Getenv("PROXY_USERNAME"),
    Password: os.Getenv("PROXY_PASSWORD"),
  }
  dialer, err := proxy.SOCKS5("tcp", os.Getenv("PROXY_ADDRESS"), &auth, nil)
  if err != nil {
    log.Fatal(err)
  }

  client := &http.Client{
    Transport: &http.Transport{
      Dial: dialer.Dial,
    },
  }

  r, err := client.Get(resource)
  if err != nil {
    log.Fatal(err)
  }
  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(string(body))
  return "you got access", nil
}

// Entry point for server side
func main() {
    port := ":"+os.Getenv("SERVER_PORT")
    listener, err := net.Listen("tcp", port)

    if err != nil {
        fmt.Println(err)
        return
    }

    defer listener.Close()
    fmt.Println("Server is listening...")
	
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }
		go handleConnection(conn)
    }
}

// Get message from client and call proxy function
func handleConnection(conn net.Conn){
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading from connection:", err)
		}
		return
	}
	clientMessage := string(buf[:n])
	fmt.Printf("Received message from client: %s\n", clientMessage)
	proxyResponse, err := ProxyRequest(clientMessage)

	if err != nil {
		conn.Write([]byte("Error accessing resource"))
		return
	}

	conn.Write([]byte(proxyResponse))
}
