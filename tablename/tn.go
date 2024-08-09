package main

import (
	"flag"
	"fmt"
	"hash/crc32"
)

var (
	id = flag.String("id", "", "id")
)

func main() {
	flag.Parse()
	if *id == "" {
		panic("id is empty")
	}
	fmt.Printf("table name - suffix : %03d\n", crc32.ChecksumIEEE([]byte(*id))%128)

}
