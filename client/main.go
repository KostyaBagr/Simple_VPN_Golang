package main
import (
    "fmt"
    "os"
    "net"
    "io"
	"log"

    "github.com/songgao/water"
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

func createTun(ip string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		return nil, err
	}
	log.Printf("Interface Name: %s\n", iface.Name())
	out, err := cmd.RunCommand(fmt.Sprintf("sudo ip addr add %s/24 dev %s", ip, iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}

	out, err = cmd.RunCommand(fmt.Sprintf("sudo ip link set dev %s up", iface.Name()))
	if err != nil {
		fmt.Println(out)
		return nil, err
	}
	return iface, nil
}