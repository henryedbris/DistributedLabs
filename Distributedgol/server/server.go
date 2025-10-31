package main

import (
	"Distributedgol/stubs"
	"flag"
	"net"
	"net/rpc"
)

type GameServer struct{}

func updateState(height int, width int, currentWorld [][]uint8, nextWorld [][]uint8) {
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

}

func (g *GameServer) HandleState(req stubs.Request, res *stubs.Response) error {
	currentWorld := make([][]uint8, req.ImgWidth)
	nextWorld := make([][]uint8, req.ImgWidth)
	for i := range currentWorld {
		currentWorld[i] = make([]uint8, req.ImgHeight)
		nextWorld[i] = make([]uint8, req.ImgWidth)
	}
	for y := 0; y < req.ImgHeight; y++ {
		for x := 0; x < req.ImgWidth; x++ {
			currentWorld[y][x] = req.Message[y][x]
		}
	}

	for i := 0; i < req.Turns; i++ {
		updateState(req.ImgHeight, req.ImgWidth, currentWorld, nextWorld)
		currentWorld, nextWorld = nextWorld, currentWorld
	}
	res.Message = currentWorld
	return nil
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rpc.Register(&GameServer{})
	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	rpc.Accept(listener)
}
