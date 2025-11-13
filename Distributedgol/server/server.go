package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"

	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
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

func updateState(height int, width int, currentWorld [][]uint8, out chan [][]uint8, g *GameState) {

	nextWorld := make([][]uint8, height)
	for i := 0; i < height; i++ {
		nextWorld[i] = make([]uint8, width)
	}

	for y := 0; y < height; y++ {
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

			if currentWorld[y][x] == 255 {
				if sum < 2 {
					nextWorld[y][x] = 0
				} else if sum == 2 || sum == 3 {
					nextWorld[y][x] = 255
				} else {
					nextWorld[y][x] = 0
				}
			} else {
				if sum == 3 {
					nextWorld[y][x] = 255
				} else {
					nextWorld[y][x] = 0
				}

			}
		}
	}
	out <- nextWorld
}

func calculateAliveCells(imgHeight int, imgWidth int, world [][]byte) []util.Cell {
	var aliveCells []util.Cell
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			if world[y][x] == 255 {
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}
	return aliveCells
}

func makeWorld(imgHeight int, imgWidth int, world [][]byte) [][]uint8 {
	currentWorld := make([][]uint8, imgHeight)
	for i := range currentWorld {
		currentWorld[i] = make([]uint8, imgWidth)
	}
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			currentWorld[y][x] = world[y][x]
		}
	}
	return currentWorld
}
func switchPause(g *GameState, res *stubs.Response) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	g.Paused = !g.Paused
	res.Paused = g.Paused
	res.Turn = g.Turn
	return nil
}

func saveWorld(g *GameState, res *stubs.Response) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	res.Message = g.World
	res.Turn = g.Turn
	return nil
}
func (g *GameState) HandleAlive(req stubs.CellRequest, res *stubs.Response) error {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	aliveCells := len(calculateAliveCells(g.Height, g.Width, g.World))
	res.Cells = aliveCells
	res.Turn = g.Turn
	return nil
}

func (g *GameState) HandleKeyPress(req stubs.KeyRequest, res *stubs.Response) error {
	keyPress := req.Key
	switch {
	case keyPress == 'p': // pause gol
		return switchPause(g, res)
	case keyPress == 's': // send pgm image of current gol
		return saveWorld(g, res)
	case keyPress == 'q': // send pgm image of final turn , then quit gol
		saveWorld(g, res)
	case keyPress == 'k':
		saveWorld(g, res)
	default:
	}
	return nil
}

func (g *GameState) HandleQuit(req stubs.QuitRequest, response *stubs.Response) error {
	if req.Flag {
		saveWorld(g, response)

		go func() {
			time.Sleep(1 * time.Second)
			g.Quit <- true
		}()

	}
	return nil
}

func (g *GameState) HandleState(req stubs.Request, res *stubs.Response) error {
	currentWorld := makeWorld(req.ImgHeight, req.ImgWidth, req.Message)
	g.Lock.Lock()
	g.Height = req.ImgHeight
	g.Width = req.ImgWidth
	g.World = currentWorld
	g.Paused = false
	g.Lock.Unlock()
	channel := make(chan [][]uint8)
	turn := 0
	for i := 0; i < req.Turns; {
		if !g.Paused {
			nextWorld := makeWorld(req.ImgHeight, req.ImgWidth, currentWorld)
			go updateState(req.ImgHeight, req.ImgWidth, currentWorld, channel, g)
			nextWorld = <-channel
			currentWorld, nextWorld = nextWorld, currentWorld
			turn += 1
			g.Lock.Lock()
			g.Turn = turn
			fmt.Println(turn)
			g.World = currentWorld
			g.Lock.Unlock()
			i++
		} else { // wait until world is unpaused to continue game
			for g.Paused {
				time.Sleep(50 * time.Millisecond)
			}
		}

	}
	res.Message = g.World
	res.Turn = turn
	return nil
}

func main() {
	pAddr := flag.String("port", "8020", "Port to listen on")
	flag.Parse()

	g := &GameState{Quit: make(chan bool)}
	rpc.Register(g)

	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go rpc.Accept(listener)
	<-g.Quit
	fmt.Println("Shutting down...")
	os.Exit(0)
}
