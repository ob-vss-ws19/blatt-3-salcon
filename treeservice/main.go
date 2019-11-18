package main

import (
	"blatt-3-salcon/tree"
	"github.com/AsynkronIT/protoactor-go/actor"
)

func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.Node{}
	})
	pid := context.Spawn(props)
	context.Send(pid, &tree.Add{Key: 1, Val: "Sali"})
	context.Send(pid, &tree.Add{Key: 2, Val: "Salihh"})

}
