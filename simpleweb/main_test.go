package main

import (
	"fmt"
	"testing"
)

func Test_Reader(t *testing.T) {
	rdr := dummyReader{Size: 500}

	buf := make([]byte, 1000)

	i, err := rdr.Read(buf)

	fmt.Println(i)
	fmt.Println(err)
}
