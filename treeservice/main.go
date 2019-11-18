package main

import (
	"blatt-3-salcon/tree"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
)

func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.Node{}
	})

	pid := context.Spawn(props)
	context.Send(pid, &tree.Test1{Message: "hi", Name: "Name"})

	fmt.Println(pid)

}
