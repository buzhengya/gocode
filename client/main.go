package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	BufSize = 1024
)

func main() {
	for i := 0; i < 1024; i++{
		go Conn(i)
		time.Sleep(time.Duration(3) * time.Second)
		if i == 1023 {
			i = 0
		}
	}
	time.Sleep(time.Duration(1000) * time.Second)
}

func Conn(i int) {
	defer fmt.Println("exit conn. ", i)
	conn, err := net.Dial("tcp", "10.224.14.205:8089")
	if err != nil {
		fmt.Println("dial error : ", err)
		return
	}

	chanRead := make(chan string, 16)
	defer conn.Close()
	go Send(conn, chanRead)

	sliBuf := make([]byte, BufSize)
	for {
		nLen, err := conn.Read(sliBuf)
		if err != nil {
			fmt.Println("read error : ", err)
			return
		}

		strContent := <-chanRead
		if strContent != string(sliBuf[:nLen]) {
			//fmt.Println("send : ", strContent)
			//fmt.Println("recv", string(sliBuf[:nLen]))
			panic("send and recv content not match.")
		}
		//fmt.Println("finish deal : ", strContent)
		if rand.Intn(BufSize) == 0 {
			return
		}
	}
}


func Send(conn net.Conn, chanRead chan string) {
	for {
	strContent := GetStr(int32(rand.Intn(BufSize)+1))
	//fmt.Println("Send : ", strContent)
	conn.Write([]byte(strContent))
	chanRead <- strContent
	time.Sleep(time.Duration(1) * time.Second)
	select {
	case <-chanRead :
		var tmp []byte
		if _, err := conn.Read(tmp); err != nil {
			return
		}
		panic("msg not return from server.")
	default:
		}
	}
}

func GetStr(nLen int32) string {
	sliContent := make([]byte, nLen)
	for i := 0; i < int(nLen); i++ {
		sliContent[i] = 'a' + byte(rand.Intn(26))
	}
	return string(sliContent)
}