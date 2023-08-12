package main

import (
	"fmt" // пакет используется для проверки выполнения условия задачи, не удаляйте его
	"math/rand"
	"time" // пакет используется для проверки выполнения условия задачи, не удаляйте его
)

type t struct {
	num int
	i   int
}

func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	x1 := make(chan t)
	y1 := make(chan t)

	for i := 0; i < n; i++ {
		go func(i int) {
			x := <-in1
			x1 <- t{num: fn(x), i: i}
		}(i)

		go func(i int) {
			y := <-in2
			y1 <- t{num: fn(y), i: i}
		}(i)

	}

	xx := make([]int, n)
	yy := make([]int, n)

	done1 := make(chan struct{})
	done2 := make(chan struct{})

	go func() {
		for i := 0; i < n; i++ {
			tx := <-x1
			xx[tx.i] = tx.num
		}
		done1 <- struct{}{}
		fmt.Println("ready closed")
	}()
	go func() {
		for i := 0; i < n; i++ {
			ty := <-y1
			yy[ty.i] = ty.num
		}
		done2 <- struct{}{}
		fmt.Println("ready closed")
	}()

	go func() {
		<-done1
		<-done2

		fmt.Println("ready close received")

		for i := 0; i < len(xx); i++ {
			out <- (xx[i] + yy[i])
		}
	}()
}

const N = 5

func main() {
	fn := func(x int) int {
		time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
		return x * 2
	}
	in1 := make(chan int, N)
	in2 := make(chan int, N)
	out := make(chan int, N)

	start := time.Now()
	merge2Channels(fn, in1, in2, out, N+1)
	for i := 0; i < N+1; i++ {
		in1 <- i
		in2 <- i
	}

	orderFail := false
	EvenFail := false
	for i, prev := 0, 0; i < N; i++ {
		c := <-out
		if c%2 != 0 {
			EvenFail = true
		}
		if prev >= c && i != 0 {
			orderFail = true
		}
		prev = c
		fmt.Println(c)
	}
	if orderFail {
		fmt.Println("порядок нарушен")
	}
	if EvenFail {
		fmt.Println("Есть не четные")
	}
	duration := time.Since(start)
	if duration.Seconds() > N {
		fmt.Println("Время превышено")
	}
	fmt.Println("Время выполнения: ", duration)
}
