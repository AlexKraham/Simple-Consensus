package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func startListener(port string) {
	// start listening
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	fmt.Printf("Listening on port: %s \n", port)

	for {
		conn, err := l.Accept()
		// defer conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Port [%s] : Listener accepted a connection.\n", port)

		netData, err := bufio.NewReader(conn).ReadString('\n')

		go func(c net.Conn) {
			// echo all incoming data
			fmt.Println("From process: " + netData)
			io.Copy(c, c)
			// shut down the connection.
			c.Close()
		}(conn)
	}
}

// keeps on trying to get connection
func getConn(port string) net.Conn {
	for {
		d, err := net.Dial("tcp", "localhost:"+port)
		// no err, able to dial
		if err == nil {
			fmt.Println("Connected success - port: " + port)
			return d
		}
		fmt.Println("Failed to get connection on port: " + port)
		time.Sleep(1 * time.Second)
	}
}

func startDialer(port string) {
	dial := getConn(port)
	defer dial.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> Send Msg: ")
	text, _ := reader.ReadString('\n')

	if _, err := dial.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	wg.Add(2)

	// Get command line args
	args := os.Args[1:]

	// server
	go startListener(args[0])

	go startDialer(args[1])

	wg.Wait()
}
