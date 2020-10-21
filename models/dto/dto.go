package dto

// ReturnSlice ...
type ReturnSlice struct {
	Date []string
	Izcr []float64
	Izr  []float64
	Ld   []float64
	Plt  []float64
	Slc  []float64
	Yen  []float64
	YMax int
	YMin int
}

type Message struct {
	Message string `json:"message"`
}

type Activity struct {
	Title string
	URL   string
	Time  string
}
