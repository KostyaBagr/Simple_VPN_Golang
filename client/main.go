package main
import (
    "fmt"
    "os"
    "net"
    "io"
    "log"

    "github.com/joho/godotenv"

)
func main() {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
        return 
    }
    serverAddress := net.JoinHostPort(os.Getenv("SERVER_IP"), os.Getenv("SERVER_PORT"))
    conn, err := net.Dial("tcp", serverAddress) 

    if err != nil { 
        fmt.Println(err) 
        return
    } 
	fmt.Print("Message sent to server. Waiting for response\n")
	_, err = conn.Write([]byte("https://www.linkedin.com/feed/"))
    defer conn.Close() 
  
    io.Copy(os.Stdout, conn) 
    fmt.Println("\nDone")
}


