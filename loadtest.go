package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

//Loadtester is the type for handling actual load test.
type Loadtester struct {
	address string
	conns   []*net.TCPConn
}

func (obj *Loadtester) makeLoadRequest(conn *net.TCPConn, maxLines int, clientNo int, clientTotal int) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	line := r1.Intn(maxLines) + 1
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	cmdString := "GET " + strconv.Itoa(line) + "\n"

	start := time.Now()
	_, err1 := writer.Write([]byte(cmdString))
	obj.checkError(err1)
	writer.Flush()
	//print(cmdString + " is sent!")
	//result, err2 := reader.ReadByte()

	status, err2 := reader.ReadString('\n')
	obj.checkError(err2)
	if status == "OK" {
		_, err3 := reader.ReadString('\n')
		obj.checkError(err3)
	}
	elapsed := time.Since(start).Seconds()
	//log.Printf("The time taken: %s - %d of %d %s -> %d", elapsed, clientNo, clientTotal, cmdString, result)

	log.Printf(" %d %f ", clientTotal, elapsed*1000)

	//print(result)
}

func (obj *Loadtester) createConnections(numOfConn int) {
	obj.conns = make([]*net.TCPConn, numOfConn)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", obj.address)
	obj.checkError(err)
	for i := 0; i < numOfConn; i++ {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		obj.checkError(err)
		obj.conns[i] = conn
	}
}

func (obj *Loadtester) closeConnections() {
	for _, conn := range obj.conns {
		conn.Close()
	}
	obj.conns = obj.conns[:0]
}

// StartLoadTest the main method for starting a load test with numOfClients concurrently
func (obj *Loadtester) StartLoadTest(numOfClients int) {

	//obj.createConnections(numOfClients)

	var ch = make(chan bool, numOfClients)
	var wg sync.WaitGroup

	wg.Add(numOfClients)
	for i := 0; i < numOfClients; i++ {
		go func(i int) {
			start := <-ch
			if start {
				tcpAddr, err := net.ResolveTCPAddr("tcp4", obj.address)
				obj.checkError(err)
				conn, err := net.DialTCP("tcp", nil, tcpAddr)
				obj.checkError(err)
				for k := 0; k < 10; k++ {
					obj.makeLoadRequest(conn, 10000000, i, numOfClients)
				}
				conn.Close()
				wg.Done()
			}
		}(i)
	}

	for i := 0; i < numOfClients; i++ {
		ch <- true // start each thread
	}

	close(ch)
	wg.Wait()

	//obj.closeConnections()

	print("Load testing " + strconv.Itoa(numOfClients) + " is DONE!\n")
}

func (obj *Loadtester) checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
