package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	tunDevice = "/dev/net/tun"
	ifReqSize = unix.IFNAMSIZ + 64
)

func checkErr(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func main() {
	// Открываем TUN устройство
	nfd, err := unix.Open(tunDevice, os.O_RDWR, 0)
	checkErr(err)

	var ifr [ifReqSize]byte
	var flags uint16 = unix.IFF_TUN | unix.IFF_NO_PI
	name := []byte("wg1")
	copy(ifr[:unix.IFNAMSIZ], name)
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = flags

	// Печатаем имя интерфейса для отладки
	fmt.Println("Interface Name:", string(ifr[:unix.IFNAMSIZ]))

	// Выполняем IOCTL для настройки TUN устройства
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(nfd),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr[0])),
	)
	if errno != 0 {
		checkErr(fmt.Errorf("ioctl errno: %d", errno))
	}

	// Открываем файловый дескриптор для чтения
	fd := os.NewFile(uintptr(nfd), tunDevice)
	checkErr(err)

	for {
		buf := make([]byte, 1500)
		_, err := fd.Read(buf)
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			continue
		}

		fmt.Println("received packet")

		switch buf[0] & 0xF0 {
		case 0x40:
			fmt.Println("received ipv4")
			fmt.Printf("Length: %d\n", binary.BigEndian.Uint16(buf[2:4]))
			fmt.Printf("Protocol: %d (1=ICMP, 6=TCP, 17=UDP)\n", buf[9])
			fmt.Printf("Source IP: %s\n", net.IP(buf[12:16]))
			fmt.Printf("Destination IP: %s\n", net.IP(buf[16:20]))
		case 0x60:
			fmt.Println("received ipv6")
			fmt.Printf("Length: %d\n", binary.BigEndian.Uint16(buf[4:6]))
			fmt.Printf("Protocol: %d (1=ICMP, 6=TCP, 17=UDP)\n", buf[7])
			fmt.Printf("Source IP: %s\n", net.IP(buf[8:24]))
			fmt.Printf("Destination IP: %s\n", net.IP(buf[24:40]))
		}
	}
}

