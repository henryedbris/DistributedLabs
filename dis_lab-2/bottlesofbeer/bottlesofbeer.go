package main
import (
	"flag"
	"net/rpc"
	"fmt"
	"time"
	"net")

var initialised = false
var nextround *rpc.Client
var nextAddr string

func Beers(i int) {
	time.Sleep(1*time.Second) //Delay just so the round is visible
	if i > 1 {
		fmt.Printf("%v bottles of beer on the wall, %v bottles of beer\nTake one down, pass it around...", i, i)
	} else if i == 1 {
		fmt.Printf("One bottle of beer on the wall, one bottle of beer\nTake it down, pass it around...")
	} else {
		fmt.Printf("No more bottles of beer on the wall, no more bottles of beer\nGo to the store and buy some more!\n")
	}
}

func PassItAround(bottles int){
	reqs := Token{Bottles: bottles}
	resp := new(Token)
	if initialised == false {
		nextround, _ = rpc.Dial("tcp", nextAddr)
		initialised = true
	}
	nextround.Go("BottlesOfBeer.Round", reqs, resp, nil)
}

type BottlesOfBeer struct {}
type Token struct {
	Bottles int
}

func (b *BottlesOfBeer) Round(intoken Token, outtoken *Token) (err error) {
	bottles := intoken.Bottles
	Beers(bottles)
	if bottles > 0 {
		PassItAround(bottles-1)
	}
	return
}

func main(){
	thisPort := flag.String("this", "8030", "Port to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	bottles := flag.Int("n",0, "Bottles of Beer (launches round)")
	flag.Parse()
	rpc.Register(&BottlesOfBeer{})
	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer listener.Close()
	if *bottles > 0 {
		Beers(*bottles)
		go PassItAround(*bottles-1)
	}
	rpc.Accept(listener)
}
