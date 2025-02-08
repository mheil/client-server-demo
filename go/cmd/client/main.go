package main

import (
	"bufio"
	"client-server/pkg/errhandler"
	"client-server/pkg/msg"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var mp = msg.NewMessagePrinter()

func main() {
	host := flag.String("host", "localhost", "hostname of the server")
	port := flag.Int("port", 7666, "port number of the server")

	fmt.Println("client started")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	errhandler.PanicOnError(err)
	defer errhandler.CloseWithPanicOnError(conn)

	mp.PrintInOut(conn, "connected")

	outMsg := make(chan string)
	go readCommands(outMsg, conn)

	outQuit := make(chan bool)
	completed := make(chan bool)
	go handleOutgoingData(outMsg, outQuit, conn, completed)

	go readIncomingMessages(conn, outQuit, completed)

	<-completed
	<-completed

	mp.PrintInOut(conn, "finished")
}

func readIncomingMessages(conn net.Conn, quit chan bool, completed chan bool) {
	defer func() {
		completed <- true
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			mp.PrintInOut(conn, "closed")
			quit <- true
			return
		}

		errhandler.PanicOnError(err)

		message = strings.TrimSpace(message)
		mp.PrintIn(conn, message)
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

func handleOutgoingData(messages chan string, quit chan bool, conn net.Conn, completed chan bool) {
	defer func() {
		completed <- true
	}()

	for {
		select {
		case text := <-messages:
			mp.PrintOut(conn, text)
			_, err := fmt.Fprintln(conn, text)
			errhandler.PanicOnError(err)
		case <-quit:
			return
		}
	}
}
