package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"uk.ac.bris.cs/solutions/distributed2/secretstrings/stubs"
)

func makeCall(client *rpc.Client, message string) {
	request := stubs.Request{Message: message}
	response := new(stubs.Response)
	client.Call(stubs.ReverseHandler, request, response)
	fmt.Println("Responded: " + response.Message)
}

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	client, err := rpc.Dial("tcp", *server)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	file, err := os.Open("wordlist")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		fmt.Println("Called: " + t)
		makeCall(client, t)
	}

}
