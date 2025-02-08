package main

import (
	"bufio"
	"client-server/pkg/errhandler"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func main() {
	fmt.Println("server started")

	port := flag.Int("port", 7666, "port number to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	errhandler.PanicOnError(err)

	for {
		fmt.Printf("%s listening\n", ln.Addr())
		conn, err := ln.Accept()
		errhandler.PanicOnErrorWithMessage("listening for server failed", err)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("%s <-> %s closing failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
		}
	}()

	fmt.Printf("%s <-> %s opened\n", conn.LocalAddr(), conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Printf("%s <-> %s closed\n", conn.LocalAddr(), conn.RemoteAddr())
			break
		}

		if err != nil {
			fmt.Printf("%s <-> %s failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
			break
		}

		message = strings.TrimSpace(message)
		fmt.Printf("%s <-- %s '%s'\n", conn.LocalAddr(), conn.RemoteAddr(), message)
		if message == "exit" {
			fmt.Printf("%s --> %s bye\n", conn.LocalAddr(), conn.RemoteAddr())
			_, err = fmt.Fprintln(conn, "bye")
			if err != nil {
				fmt.Printf("%s <-> %s sending bye failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
			}
			fmt.Printf("%s <-> %s closing\n", conn.LocalAddr(), conn.RemoteAddr())
			errhandler.CloseWithPanicOnError(conn)
			break
		}
		if message == "time" {
			timeMsg := time.Now().Format(time.RFC3339)
			fmt.Printf("%s --> %s %s\n", conn.LocalAddr(), conn.RemoteAddr(), timeMsg)
			_, err = fmt.Fprintln(conn, timeMsg)
			if err != nil {
				fmt.Printf("%s <-> %s sending bye failed: %s\n", conn.LocalAddr(), conn.RemoteAddr(), err)
				fmt.Printf("%s <-> %s closing\n", conn.LocalAddr(), conn.RemoteAddr())
				errhandler.CloseWithPanicOnError(conn)
				break
			}
		}
	}
	fmt.Printf("%s <-> %s finished\n", conn.LocalAddr(), conn.RemoteAddr())
}
