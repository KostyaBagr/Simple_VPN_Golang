package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"io"

	"github.com/joho/godotenv"
)

// Global env initialization
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}


// Create a request using proxy
func ProxyRequest(resource string) (string, error) {

    proxyStr := os.Getenv("PROXY_URL")
    proxyURL, err := url.Parse(proxyStr)
	
    if err != nil {
        fmt.Println("Error parsing proxy URL:", err)
        return "", err
    }

    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }

    fmt.Println("\nTrying to make request")
    resp, err := client.Get(resource)

    if err != nil {
        fmt.Println("Error making request:", err)
        return "", err
    }

    defer resp.Body.Close()

    if resp != nil {
        fmt.Println("Response Status:", resp.Status)
    }

	fmt.Print("\nGot access to: ", resource)
    return "You got access to resource", nil
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
