package main

import (
	"fmt"
	"testing"
)

func Test_Reader(t *testing.T) {
	rdr := dummyReader{Size: 500}

	buf := make([]byte, 10)

	i, err := rdr.Read(buf)
	i, err = rdr.Read(buf)

	fmt.Println(i)
	fmt.Println(err)
}
