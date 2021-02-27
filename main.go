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
		// fmt.Printf("Port [%s] : Listener accepted a connection.\n", port)
		// fmt.Println("CURR: " + strconv.Itoa(curr))

		netData, err := bufio.NewReader(conn).ReadString('\n')
		receivedVal, _ := strconv.Atoi(netData)

		wg.Add(1)

		go func(c net.Conn) {
			// echo all incoming data
			// fmt.Println("From process: " + netData)
			values = append(values, receivedVal)
			curr++
			if curr == n {
				min := getMin(values)
				fmt.Printf("Consensus Value: %d\n", min)
			}
			io.Copy(c, c)
			// shut down the connection.
			c.Close()
		}(conn)

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
		defer dial.Close()

		// send it!
		if _, err := dial.Write([]byte(strconv.Itoa(v))); err != nil {
			log.Fatal(err)
		}
	}

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

	curr = 1

	// server
	go startListener()

	go startDialer()

	wg.Wait()
}
