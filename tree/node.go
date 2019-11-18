package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
)

// Jeder Knoten ist ein eigener Actor
type Node struct {
	LeftNode  *actor.PID
	RightNode *actor.PID
	MaxKeyVal int
	LeafSize  int
	Data      map[int]string
}

type Add struct {
	Key int
	Val string
}

type Find struct {
	RequestFrom *actor.PID
	Key         int
}

// Implementiere Receive
func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Add:
		fmt.Printf("# ADD: Got Request for %d -> %s", msg.Key, msg.Val)
		// Fall 1: Platz im Node und noch keine Teilbäume
		if state.LeftNode == nil && state.RightNode == nil && state.LeafSize > len(state.Data) {
			// Fall 1.1: Noch kein Data angelegt
			if state.Data == nil {
				state.Data = make(map[int]string)
			}
			state.Data[msg.Key] = msg.Val
			fmt.Printf("# ADD: Data successfully added to Node. PID: %s, {Key: %d, Val: %s} \n", context.Self().Address, msg.Key, msg.Val)
			// Fall 2: Kein Platz im Node und noch keine Teilbäume
		} else if state.LeftNode == nil && state.RightNode == nil && state.LeafSize == len(state.Data) {
			// Erstelle zwei nodes (linke und Rechte Hälfte)
			props := actor.PropsFromProducer(func() actor.Actor {
				// Erstelle Node mit dem LeafSize vom Parent
				return &Node{LeafSize: state.LeafSize}
			})
			// Initialisiere Linke und Rechte Node
			state.LeftNode = context.Spawn(props)
			state.RightNode = context.Spawn(props)

			// Map in der Mitte auf die Leafs aufteilen
			state.Data[msg.Key] = msg.Val

			var keys []int
			for key := range state.Data {
				keys = append(keys, int(key))
			}
			sort.Ints(keys)

			state.MaxKeyVal = keys[(len(keys)/2)-1]
			fmt.Printf("# ADD: Maximum Key Val Left %d\n", state.MaxKeyVal)
			for _, key := range keys {
				// Rechts aufteilen
				if key > state.MaxKeyVal {
					fmt.Printf("# ADD: Set %d right\n", key)
					context.Send(state.RightNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
					// Links aufteilen
				} else {
					fmt.Printf("# ADD: Set %d left\n", key)
					context.Send(state.LeftNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
				}
			}
		}
	case *Find:
		fmt.Printf("# FIND: Got Request for %d -> %s", msg.Key)
		if state.Data == nil {

		}

	}

}
