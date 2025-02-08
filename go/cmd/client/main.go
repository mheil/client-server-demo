package main

import (
	"bufio"
	"client-server/pkg/errhandler"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "hostname of the server")
	port := flag.Int("port", 7666, "port number of the server")

	fmt.Println("client started")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	errhandler.PanicOnError(err)
	defer errhandler.CloseWithPanicOnError(conn)

	fmt.Printf("%s <-> %s connected\n", conn.LocalAddr(), conn.RemoteAddr())

	outMsg := make(chan string)
	go readCommands(outMsg, conn)

	outQuit := make(chan bool)
	completed := make(chan bool)
	go handleOutgoingData(outMsg, outQuit, conn, completed)

	go readIncomingMessages(conn, outQuit, completed)

	<-completed
	<-completed

	fmt.Printf("%s <-> %s finished\n", conn.LocalAddr(), conn.RemoteAddr())
}

func readIncomingMessages(conn net.Conn, quit chan bool, completed chan bool) {
	defer func() {
		completed <- true
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("%s <-> %s closed\n", conn.LocalAddr(), conn.RemoteAddr())
			quit <- true
			return
		}

		errhandler.PanicOnError(err)

		message = strings.TrimSpace(message)
		fmt.Printf("%s <-- %s %s\n", conn.LocalAddr(), conn.RemoteAddr(), message)
	}
}

func readCommands(msg chan string, conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			errhandler.CloseWithPanicOnError(conn)
			return
		}
		errhandler.PanicOnError(err)
		msg <- strings.TrimSpace(text)
	}
}

func handleOutgoingData(msg chan string, quit chan bool, conn net.Conn, completed chan bool) {
	defer func() {
		completed <- true
	}()

	for {
		select {
		case text := <-msg:
			fmt.Printf("%s --> %s '%s'\n", conn.LocalAddr(), conn.RemoteAddr(), text)
			_, err := fmt.Fprintln(conn, text)
			errhandler.PanicOnError(err)
		case <-quit:
			return
		}
	}
}
