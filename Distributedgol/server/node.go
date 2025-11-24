package main

import (
	"flag"

	"net"
	"net/rpc"
	"os"
	"sync"
	"time"

	"Distributedgol/stubs"
)

type GameState struct {
	Lock   sync.Mutex
	World  [][]uint8
	Height int
	Width  int
	Turn   int
	Paused bool
	Quit   chan bool
}

func updateState(height int, width int, currentWorld [][]uint8, startY int, endY int) [][]uint8 {
	nextWorldHeight := endY - startY
	nextWorld := make([][]uint8, nextWorldHeight)
	for i := 0; i < nextWorldHeight; i++ {
		nextWorld[i] = make([]uint8, width)
	}

	for y := startY; y < endY; y++ {
		for x := 0; x < width; x++ {
			var sum int
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					ny := (y + dy + height) % height
					nx := (x + dx + width) % width
					if currentWorld[ny][nx] == 255 {
						sum++
					}
				}
			}

			if currentWorld[y][x] == 255 { // if cell is alive
				if sum < 2 { // flip cell if it has <2 alive neighbours
					nextWorld[y-startY][x] = 0

				} else if sum == 2 || sum == 3 { // cell stays alive if it has 2/3 neighbours
					nextWorld[y-startY][x] = 255
				} else { // flip cell if it has >3 neighbours
					nextWorld[y-startY][x] = 0

				}
			} else { // if cell is dead
				if sum == 3 { // flip cell if it has 3 neighbours
					nextWorld[y-startY][x] = 255

				} else {
					nextWorld[y-startY][x] = 0
				}

			}
		}
	}

	return nextWorld
}

func (g *GameState) HandleQuit(req stubs.QuitRequest, response *stubs.Response) error {
	go func() {
		time.Sleep(1 * time.Second)
		g.Quit <- true
	}()

	return nil
}

func getWorkerSlice(height, threads, id int) (int, int) {
	rows := height / threads
	extraRows := height % threads

	start := id * rows
	end := start + rows


	if id == threads-1 { //give last worker leftover rows
		end += extraRows
	}

	return start, end
}

func worker(height int, width int, currentWorld [][]uint8, workerOut chan<- [][]uint8, startY int, endY int) {
	out := updateState(height, width, currentWorld, startY, endY)
	workerOut <- out

}

func (g *GameState) HandleState(req stubs.Request, res *stubs.Response) error {
	threads := 4
	channels := make([]chan [][]uint8, threads)
	for i := 0; i < threads; i++ {
		channels[i] = make(chan [][]uint8)
	}
	height := req.EndY - req.StartY

	for i := 0; i < threads; i++ {
		startY, endY := getWorkerSlice(height, threads, i)
		go worker(height, req.ImgWidth, req.Message, channels[i], startY, endY)
	}
	//go updateState(req.ImgHeight, req.ImgWidth, req.Message, req.StartY, req.EndY, out)
	temp := make([][]uint8, 0, req.ImgHeight)
	for j := 0; j < threads; j++ {
		part := <-channels[j]
		temp = append(temp, part...)
	}

	res.Message = temp
	return nil
}

func main() {
	g := &GameState{}
	g.Quit = make(chan bool)
	rpc.Register(g)
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		panic(err)
	}
	go rpc.Accept(listener)
	<-g.Quit
	os.Exit(0)
}
