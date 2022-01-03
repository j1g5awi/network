package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const TOKEN = "BLADE"

type Client struct {
	laddr         *net.UDPAddr
	raddr         *net.UDPAddr
	conn          *net.UDPConn
	maxDataLength int
}

func (c *Client) Listen() {
	for {
		buf := make([]byte, 1024)
		n, addr, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		fmt.Println(buf[:n])

		msg := &Message{}
		msg.UnMarsal(buf[:n])
		if int(msg.length) != n {
			c.Send(addr, NewMessage(ERROR, "Wrong length!"))
			continue
		}
		c.Receive(addr, msg)
	}
}

func (c *Client) Send(raddr *net.UDPAddr, msg *Message) {
	if raddr == nil {
		raddr = c.raddr
	}
	if msg.type_ == DISCONNECT {
		if msg.data == "T" {
			c.conn.WriteToUDP(msg.Marsal(), raddr)
			c.raddr = nil
			fmt.Println("Disconnected")
		}
	}
	c.conn.WriteToUDP(msg.Marsal(), raddr)
}

func (c *Client) Receive(raddr *net.UDPAddr, msg *Message) {
	switch msg.type_ {
	case ERROR:
		fmt.Println("Error: " + msg.data)
	case CONNECT:
		if msg.data != "" {
			if msg.data[0] == byte(AUTH_REQUIRE) {
				fmt.Println("Require auth!")
			} else {
				fmt.Println("Don't require auth!")
				c.raddr = raddr
			}
			switch msg.data[1] {
			case byte(MAX_DATA_LENGTH_10):
				c.maxDataLength = 10
			case byte(MAX_DATA_LENGTH_20):
				c.maxDataLength = 20
			case byte(MAX_DATA_LENGTH_50):
				c.maxDataLength = 50
			}
		} else {
			fmt.Println("Require auth!")
		}
	case AUTH:
		if msg.data != "" {
			if msg.data != TOKEN {
				fmt.Println("Auth fail")
				c.Send(raddr, NewMessage(ERROR, "Wrong token!"))
			} else {
				fmt.Println("Auth success")
				c.raddr = raddr
				c.Send(nil, NewMessage(AUTH, ""))
			}
		} else {
			fmt.Println("Auth success")
			c.raddr = raddr
		}
	case DATA:
		fmt.Println("Receive data: " + msg.data)
	case DISCONNECT:
		if msg.data != "" {
			if msg.data == "T" {
				c.raddr = nil
				fmt.Println("Disconnected")
			}
		}
	}
}

func main() {
	client := &Client{}
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if !strings.HasPrefix(text, "/") {
			continue
		} else {
			text := strings.Split(text, " ")
			switch text[0] {
			case "/listen":
				laddr, err := net.ResolveUDPAddr("udp", text[1])
				if err != nil {
					panic(err)
				}
				client.laddr = laddr

				conn, err := net.ListenUDP("udp", laddr)
				if err != nil {
					panic(err)
				}
				client.conn = conn

				go client.Listen()
			case "/connect":
				if len(text) < 3 {
					text = append(text, "")
				}

				raddr, err := net.ResolveUDPAddr("udp", text[1])
				if err != nil {
					panic(err)
				}

				client.Send(raddr, NewMessage(CONNECT, text[2]))
			case "/auth":
				if len(text) < 3 {
					text = append(text, "")
				}

				raddr, err := net.ResolveUDPAddr("udp", text[1])
				if err != nil {
					panic(err)
				}

				client.Send(raddr, NewMessage(AUTH, text[2]))
			case "/data":
				if len(text) < 2 {
					text = append(text, "")
				}
				client.Send(nil, NewMessage(DATA, text[1]))
			case "/disconnect":
				if len(text) < 2 {
					text = append(text, "")
				}
				client.Send(nil, NewMessage(DISCONNECT, text[1]))
			}
		}
	}
}
