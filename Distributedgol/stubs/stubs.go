package stubs

type Response struct {
	Message [][]uint8
	Cells   int
	Turn    int
}

type Request struct {
	Message   [][]uint8
	ImgHeight int
	ImgWidth  int
	Turns     int
}

type CellRequest struct {
	Flag bool
}
