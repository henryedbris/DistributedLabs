package stubs

var StateHandler = "GameState.HandleState"
var CellHandler = "GameState.HandleAlive"
var KeyHandler = "GameState.HandleKeyPress"
var QuitHandler = "GameState.HandleQuit"

var BStateHandler = "Broker.HandleState"
var BKeyHandler = "Broker.HandleKeyPress"
var BCellHandler = "Broker.HandleAlive"

type Response struct {
	Message [][]uint8
	Cells   int
	Turn    int
	Paused  bool
}

type Request struct {
	Message   [][]uint8
	ImgHeight int
	ImgWidth  int
	Turns     int
	StartY    int
	EndY      int
}

type CellRequest struct {
	Flag bool
}

type KeyRequest struct {
	Key rune
}

type QuitRequest struct {
	Flag bool
}
