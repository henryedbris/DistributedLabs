package stubs

type Response struct {
	Message [][]uint8
}

type Request struct {
	Message   [][]uint8
	ImgHeight int
	ImgWidth  int
	Turns     int
}
