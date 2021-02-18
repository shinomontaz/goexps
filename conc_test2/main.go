package main

import "fmt"
import "time"
import "runtime"

func main() {
	var x int
	threads := runtime.GOMAXPROCS(0) - 1
	for i := 0; i < threads; i++ {
		go func() {
			for {
				x++
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println("x =", x)
}
