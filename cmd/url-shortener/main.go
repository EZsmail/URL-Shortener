package main

type window struct {
	place string
}

type car struct {
	window,
	wheelNum int
	psOne, psTwo string
}

var (
	name = "Anatoly"
	age  = 10
)

func main() {
	// TODO: init config: cleanenv

	// TODO: init logger: log

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server:
}
