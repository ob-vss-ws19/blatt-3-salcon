package main

import (
	"blatt-3-salcon/tree"
	"github.com/AsynkronIT/protoactor-go/actor"
	"time"
)

func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.Node{}
	})
	pid := context.Spawn(props)
	context.Send(pid, &tree.Find{Key: 2, RequestFrom: pid})

	time.Sleep(2 * time.Second)
}
