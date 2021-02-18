package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	comm1 := exec.Command("ps ax")
	stdout, err := comm1.StdoutPipe()
	stdin, err := comm1.StdinPipe()

	err = comm1.Start()
	if err != nil {
		log.Fatal(err)
	}

	if err := comm1.Wait(); err != nil {
		log.Fatal(err)
	}

	io.WriteString(stdin, "grep aaa")

	fmt.Printf("result: %s\n", stdout)
}
