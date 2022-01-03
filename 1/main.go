package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			return
		} else {
			text := strings.Split(text, " ")
			if text[0] == "/listen" {
				var conn net.Conn
				switch text[1] {
				case "udp":
					laddr, _ := net.ResolveUDPAddr("udp", ":4000")
					conn, _ = net.ListenUDP("udp", laddr)
				case "tcp":
					listener, _ := net.Listen("tcp", ":4000")
					conn, _ = listener.Accept()
				case "ip":
					laddr, err := net.ResolveIPAddr("ip:254", "127.0.0.1")
					if err != nil {
						panic(err)
					}
					conn, err = net.ListenIP("ip:254", laddr)
					if err != nil {
						panic(err)
					}
				}
				func() {
					buf := make([]byte, 1024)
					for {
						n, err := conn.Read(buf)
						if err == nil {
							var data string
							if text[1] == "ip" {
								data = string(buf[16:n])
							} else {
								data = string(buf[:n])
							}
							fmt.Println(data)
						}
					}
				}()
			} else if text[0] == "/dial" {
				switch text[1] {
				case "udp":
					conn, _ := net.Dial(text[1], "127.0.0.1:4000")
					for i := 'a'; i < 'z'; i += 2 {
						time.Sleep(time.Microsecond * 100)
						data := string(i) + string(i+1)
						conn.Write([]byte(data))
					}
				case "tcp":
					conn, _ := net.Dial(text[1], "127.0.0.1:4000")
					for j := 0; j < 2; j++ {
						for i := 'a'; i < 'z'; i += 2 {
							time.Sleep(time.Microsecond * 100)
							data := string(i) + string(i+1)
							conn.Write([]byte(data))
						}
					}
				case "ip":
					conn, err := net.Dial("ip:254", "127.0.0.1")
					if err != nil {
						panic(err)
					}
					for j := 0; j < 3; j++ {
						for i := 'A'; i < 'Z'; i += 2 {
							time.Sleep(time.Microsecond * 100)
							data := string(i) + string(i+1)
							conn.Write([]byte(data))
						}
					}
				}
			}
		}
	}
}
