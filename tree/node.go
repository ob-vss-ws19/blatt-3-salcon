package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
)

// Jeder Knoten ist ein eigener Actor
type Node struct {
	LeftNode  *actor.Actor
	RightNode *actor.Actor
	MaxKeyVal int
	LeafSize  int
	Data      map[int]string
}

type Add struct {
	Key int
	Val string
}

// Implementiere Receive
func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Add:
		fmt.Printf("# ADD: %d -> %s", msg.Key, msg.Val)
		// Fall 1: Platz im Node und noch keine Teilbäume
		if state.LeftNode == nil && state.RightNode == nil && state.LeafSize > len(state.Data) {
			// Fall 1.1: Noch kein Data angelegt
			if state.Data == nil {
				state.Data = make(map[int]string)
			}
			state.Data[msg.Key] = msg.Val
			fmt.Printf("# ADD: Data successfully added to Node. PID: %s, {Key: %d, Val: %s} \n", context.Self().Address, msg.Key, msg.Val)
		}
	}
}
