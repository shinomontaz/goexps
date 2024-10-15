package main

import (
	"fmt"
	"io"
	"log"
	"syscall"

	"github.com/ncw/directio"
)

func main() {
	filepath := "./test.txt"
	in, err := directio.OpenFile(filepath, syscall.O_RDONLY, 0444)
	if err != nil {
		log.Fatal(err)
	}
	block := directio.AlignedBlock(directio.BlockSize)
	i := 0
	for {
		_, err = io.ReadFull(in, block)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Println("!!!", i, string(block))
		i++
	}
	fmt.Println("finished")

}
