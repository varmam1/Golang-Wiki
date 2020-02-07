package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}
var (
	port = flag.String("port", "9000", "port used for ws connection")
)

// TODO: Add TLS on the client 
func connect() (*websocket.Conn, error){
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port),  "", mockedIP())
}

// Created just for the sake of the exercise. Returns a fake IP
func mockedIP() string{
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])   
}

func main(){
	flag.Parse()

	// This will allow you to connect

	ws, err := connect()
	if err != nil{
		log.Fatal(err)
	}
	defer ws.Close()

	// Allows you to recieve any messages
	var m Message
	go func(){
		for{
			err := websocket.JSON.Receive(ws, &m)
			if err != nil{
				fmt.Println("Error receiving message: ", err.Error())
				break		  
			}
			fmt.Println(m)
		}
	}()

	// Allows you to send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan(){
		text := scanner.Text()
		if text == ""{
			continue
		}
		m := Message{Text : text,}
		err = websocket.JSON.Send(ws, m)
		if err != nil{
			fmt.Println("Error sending message: ", err.Error())
			break
		}
	}
}