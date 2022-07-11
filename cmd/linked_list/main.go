package main

import (
	"fmt"

	"github.com/loopholelabs/common/pkg/pool"
)

func main() {
	newLinkedList := pool.NewLinkedList[string]()

	newLinkedList.Insert("Hey!")
	newLinkedList.Insert("You")
	newLinkedList.Insert("Are")
	newLinkedList.Insert("Cool")
	newLinkedList.Insert(":)")

	fmt.Println(newLinkedList.ToArray(), newLinkedList.Len())
}
