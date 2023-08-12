package main

func main() {
	ch := make(chan int64, 4)

	ch <- 1
	ch <- 2
	ch <- 3
	_ = <-ch

	ch <- 3
	ch <- 3

	close(ch)

	return
}
