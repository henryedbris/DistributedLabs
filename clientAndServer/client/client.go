package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func read(conn *net.Conn) {
	reader := bufio.NewReader(*conn)
	msg, _ := reader.ReadString('\n')
	fmt.Println(msg)
}

func main() {
	stdin := bufio.NewReader(os.Stdin)
	conn, _ := net.Dial("tcp", "98.80.121.192:8030")
	for {
		fmt.Println("Enter text:")
		text, _ := stdin.ReadString('\n')
		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			return
		}
		read(&conn)
	}
}
