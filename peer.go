package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var i int
var n int
var v int
var curr int
var values []int
var sentToAll bool
var endListener bool

func getPort(offset int) string {
	port_num := 9000 + (offset * 4)
	return strconv.Itoa(port_num)
}

func getMin(arr []int) int {
	min := arr[0]
	for _, val := range arr {
		if val < min {
			min = val
		}
	}
	return min
}

func startListener() {
	port := getPort(i)
	// start listening
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("We are in err of the listener")
		log.Fatal(err)
	}
	defer l.Close()
	wg.Done()
	// defer fmt.Printf("when does this print\n")
	// fmt.Printf("Peer #%d - Listening on port: %s \n", i, port)

	for {
		// fmt.Printf("Waiting for connection \n")
		conn, err := l.Accept()
		// defer conn.Close()
		// defer fmt.Printf("DOes tis print in the for \n")
		if err != nil {
			// if endListener == true {
			// 	fmt.Println("Ending...")
			// 	return
			// } else {
			// 	log.Fatal(err)
			// }
		}
		// fmt.Printf("Port [%s] : Listener accepted a connection.\n", port)
		// fmt.Println("CURR: " + strconv.Itoa(curr))

		netData, err := bufio.NewReader(conn).ReadString('\n')
		receivedVal, _ := strconv.Atoi(netData)
		curr++

		go func(c net.Conn) {
			// echo all incoming data
			// fmt.Println("From process: " + netData)
			values = append(values, receivedVal)
			// fmt.Printf("Received value: %d\n", receivedVal)

			if curr == n {
				min := getMin(values)
				fmt.Printf("Consensus Value: %d\n", min)
				// l.Close()
			}
			io.Copy(c, c)
			// shut down the connection.
			c.Close()
		}(conn)

		// if curr == n && sentToAll == true {
		// 	fmt.Printf("We should close server here\n")
		// 	endListener = true
		// 	conn.Close()
		// 	l.Close()
		// } else {
		// 	fmt.Printf("Curr: %d\n", curr)
		// }
	}

}

// keeps on trying to get connection
func getConn(offset int) net.Conn {
	port := getPort(offset)
	for {
		d, err := net.Dial("tcp", "localhost:"+port)
		// no err, able to dial
		if err == nil {
			// fmt.Println("Connected success - port: " + port)
			return d
		}
		time.Sleep(1 * time.Second)
	}
}

// need to dial to all ports
func startDialer() {

	var dial net.Conn
	for j := 0; j < n; j++ {
		// don't need to send it itself
		if j == i {
			continue
		}
		dial = getConn(j)
		// fmt.Printf("Peer #%d : Sending value #%d to peer #%d\n", i, v, j)
		defer dial.Close()

		// send it!
		if _, err := dial.Write([]byte(strconv.Itoa(v))); err != nil {
			log.Fatal(err)
		}
	}
	sentToAll = true
	wg.Done()
	// fmt.Printf("End of dialer \n")
	// if curr == n && sentToAll == true {
	// 	fmt.Printf("We should end here - end of dilaer and serv \n")
	// }

	// dial := getConn(0)
	// defer dial.Close()

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print(">> Send Msg: ")
	// text, _ := reader.ReadString('\n')

	// if _, err := dial.Write([]byte(text)); err != nil {
	// 	log.Fatal(err)
	// }
}

// Usage:
//    go run main.go <i> <n>
//    i := the process number
//    n := the number of total processes
func main() {
	wg.Add(2)

	// Get command line args
	args := os.Args[1:]
	process_number, e1 := strconv.Atoi(args[0])
	total_processes, e2 := strconv.Atoi(args[1])
	val, e3 := strconv.Atoi(args[2])

	if e1 != nil || e2 != nil || e3 != nil {
		fmt.Println("Invalid command argument")
		os.Exit(-1)
	}

	i = process_number
	n = total_processes
	v = val
	values = append(values, v)
	sentToAll = false

	curr = 1

	// server
	go startListener()

	go startDialer()

	wg.Wait()
	fmt.Printf("The end\n")
}
