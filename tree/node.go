package tree

import (
	"blatt-3-salcon/messages"
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

type Remove struct {
	Key int
}

// Implementiere Receive
func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Add:
		fmt.Printf("\n# ADD: Got Request for %d -> %s\n\n", msg.Key, msg.Val)
		// Fall 1: Platz im Node und noch keine Teilbäume
		if state.LeftNode == nil && state.RightNode == nil && state.LeafSize > len(state.Data) {
			// Fall 1.1: Noch kein Data angelegt
			if state.Data == nil {
				state.Data = make(map[int]string)
			}
			state.Data[msg.Key] = msg.Val
			fmt.Printf("\n# ADD: Data successfully added to Node. PID: %s, {Key: %d, Val: %s} \n\n", context.Self().Address, msg.Key, msg.Val)
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
			fmt.Printf("\n# ADD: Maximum Key Val Left %d\n\n", state.MaxKeyVal)
			for _, key := range keys {
				// Rechts aufteilen
				if key < state.MaxKeyVal {
					fmt.Printf("\n# ADD: Set %d left\n\n", key)
					context.Send(state.LeftNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
				} else {
					fmt.Printf("\n# ADD: Set %d right\n\n", key)
					context.Send(state.RightNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
					// Links aufteilen
				}
			}
		}

	case *Find:
		fmt.Printf("\n# FIND: Got Request for %d\n\n", msg.Key)
		// Look, it the next nodes have the key, because Datalength is 0
		if state.LeftNode != nil && state.RightNode != nil && len(state.Data) == 0 {
			if msg.Key <= state.MaxKeyVal {
				// Search Left Node
				context.Send(state.LeftNode, msg)
			} else {
				// Search Right Node
				context.Send(state.RightNode, msg)
			}
		} else if state.LeftNode == nil && state.RightNode == nil {
			fmt.Printf("\n# FIND: Searching for Key %d\n", msg.Key)
			fmt.Printf("\n# FIND: Searching in map %s\n", state.Data)

			foundData := state.Data[msg.Key]
			fmt.Printf("\n# FIND: Data found %s\n", foundData)

			if foundData != "" {
				context.Send(msg.RequestFrom, &messages.Response{Key: int32(msg.Key), Value: foundData, Type: messages.FIND})
				fmt.Printf("# FIND: Key %d found\n", msg.Key)
			}
		} else {
			context.Send(msg.RequestFrom, &messages.Error{Message: "# FIND: Key not found"})
		}

	case *Remove:
		fmt.Printf("\n# REMOVE: Search for key %d to remove it\n\n", msg.Key)
		// Leaf
		if state.Data != nil {
			if _, ok := state.Data[msg.Key]; ok {
				delete(state.Data, msg.Key)
				fmt.Printf("\n # REMOVE: Key found in Tree -> Remove: %d\n\n", msg.Key)
			} else {
				fmt.Printf("\n # REMOVE: Could not find Key in Tree!: %d\n\n", msg.Key)
			}
		} else {
			// Inner Node
			if msg.Key <= state.MaxKeyVal {
				context.Send(state.LeftNode, &Remove{Key: msg.Key})
			} else {
				context.Send(state.RightNode, &Remove{Key: msg.Key})
			}
		}
	}

}
