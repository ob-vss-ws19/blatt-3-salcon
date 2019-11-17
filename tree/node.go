package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
)

// Jeder Knoten ist ein eigener Actor
type Node struct {
	Left    *actor.Actor
	Right   *actor.Actor
	MaxKeys int
}

type Test1 struct {
	Message string
	Name    string
}

type HelloActor struct{}

// Recieve will als Parameter Context von einem Actor
func (state *HelloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Test1:
		fmt.Println(msg.Message + "" + msg.Name)
	}
}
