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

type Delete struct {
	CurrentNode *actor.PID
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
	RequestFrom *actor.PID
	Key         int
}

type KeyValuePair struct {
	Key   int
	Value string
}

type Traverse struct {
	Values         []KeyValuePair
	RemainingNodes []*actor.PID
	Caller         *actor.PID
	Start          *actor.PID
}

// Implementiere Receive
func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *Add:
		// If there is Room in Leaf
		if (len(state.Data) < state.LeafSize || state.Data[msg.Key] != "") && state.LeftNode == nil && state.RightNode == nil {
			// Add a map to leaf
			if state.Data == nil {
				state.Data = make(map[int]string)
			}
			state.Data[msg.Key] = msg.Val
			fmt.Printf("\n ### ADD: Added key: %d\n", msg.Key)
		} else if state.LeftNode != nil && state.RightNode != nil && len(state.Data) == 0 {
			// If left & right nodes are already defined
			if msg.Key <= state.MaxKeyVal {
				// add left
				context.Send(state.LeftNode, msg)
			} else {
				// add right
				context.Send(state.RightNode, msg)
			}
		} else if len(state.Data) == state.LeafSize && state.LeftNode == nil && state.RightNode == nil {
			// If left & right nodes are not defined
			fmt.Printf("\n ### ADD: Create new Leafs Left & Right\n")
			props := actor.PropsFromProducer(func() actor.Actor {
				return &Node{LeafSize: int(state.LeafSize)}
			})
			state.LeftNode = context.Spawn(props)

			state.RightNode = context.Spawn(props)

			state.Data[msg.Key] = msg.Val
			keys := sortKeys(state.Data)
			state.MaxKeyVal = keys[(len(keys)/2)-1]
			fmt.Printf("\n ### ADD: Maximum Key Value at left %d\n", state.MaxKeyVal)
			// Values will be sended to leafs
			for _, key := range keys {
				if key <= state.MaxKeyVal {
					// add half left
					fmt.Printf("\n ### ADD: Send Key %d to left leaf\n", msg.Key)
					context.Send(state.LeftNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
				} else {
					// add half right
					fmt.Printf("\n ### ADD: Send Key %d to right leaf\n", msg.Key)
					context.Send(state.RightNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
				}
			}
		}

	case *Find:
		fmt.Printf("\n### FIND: Got Request for %d\n", msg.Key)
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
			fmt.Printf("\n### FIND: Searching for Key %d\n", msg.Key)
			fmt.Print("\n### FIND: Searching in map: ")
			fmt.Println(state.Data)

			foundData := state.Data[msg.Key]
			fmt.Printf("\n### FIND: Data found %s\n", foundData)

			if foundData != "" {
				context.Send(msg.RequestFrom, &messages.Response{Key: int32(msg.Key), Value: foundData, Type: messages.FIND})
				fmt.Printf("### FIND: Key %d found\n", msg.Key)
			}
		} else {
			context.Send(msg.RequestFrom, &messages.Error{Message: "### FIND: Key not found"})
		}

	case *Remove:
		fmt.Printf("\n### REMOVE: Got Request for %d\n", msg.Key)
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
			fmt.Printf("\n### REMOVE: Searching for Key %d\n", msg.Key)
			fmt.Print("\n### REMOVE: Searching in map: \n")
			fmt.Println(state.Data)

			foundData := state.Data[msg.Key]
			fmt.Printf("\n### REMOVE: Data found %s\n", foundData)

			if foundData != "" {
				delete(state.Data, msg.Key)
				fmt.Printf("\n REMOVE: deleted key %d\n", msg.Key)
			}
		} else {
			context.Send(msg.RequestFrom, &messages.Error{Message: "### REMOVE: Key not found"})
		}

	case *Traverse:
		if msg.Start != nil {
			// set root node as start node for traverse
			msg.Values = make([]KeyValuePair, 0)
			msg.RemainingNodes = make([]*actor.PID, 0)
			tmp := msg.Start
			msg.Start = nil
			if state.LeftNode == nil && state.RightNode == nil {
				// if root is leaf create slices and set start to nil
				fmt.Printf("\n ### TRAVERSE: SEND TO START\n")
				context.Send(tmp, msg)
				return
			}
			// if root is node, create slices, set start to nil, add right node to remaining and forward
			msg.RemainingNodes = append(msg.RemainingNodes, state.RightNode)
			fmt.Printf("\n ### TRAVERSE: SEND TO LEFT NODE FROM START\n")
			context.Send(state.LeftNode, msg)
			return
		} else if state.LeftNode != nil && state.RightNode != nil {
			// node is not leaf
			// while remaining nodes add right node to remaining and send to left node
			fmt.Printf("\n ### TRAVERSE: ADD RIGHT NODE SEND TO LEFT\n")

			msg.RemainingNodes = append(msg.RemainingNodes, state.RightNode)
			context.Send(state.LeftNode, msg)
		} else if len(msg.RemainingNodes) != 0 && state.LeftNode == nil && state.RightNode == nil {
			// leaf with remaining nodes to traverse

			var keys []int
			for key := range state.Data {
				keys = append(keys, int(key))
			}
			sort.Ints(keys)

			for _, key := range keys {
				fmt.Printf("\n ### TRAVERSE: appending %d\n", key)
				msg.Values = append(msg.Values, KeyValuePair{key, state.Data[key]})
			}
			next := msg.RemainingNodes[len(msg.RemainingNodes)-1]
			msg.RemainingNodes = msg.RemainingNodes[:len(msg.RemainingNodes)-1]
			fmt.Println("\n ### TRAVERSE: send to next node from leaf")
			context.Send(next, msg)
		} else if len(msg.RemainingNodes) == 0 && state.LeftNode == nil && state.RightNode == nil {
			// leaf with no remaining nodes to traverse

			var keys []int
			for key := range state.Data {
				keys = append(keys, int(key))
			}
			sort.Ints(keys)

			for _, key := range keys {
				fmt.Printf("\n ### TRAVERSE: appending %d in last leaf\n", key)
				msg.Values = append(msg.Values, KeyValuePair{key, state.Data[key]})
			}

			response := make([]*messages.Response, 0)

			fmt.Println("\n ### TRAVERSE: send to caller")
			for _, pair := range msg.Values {
				response = append(response, &messages.Response{Value: pair.Value, Key: int32(pair.Key)})
			}
			context.Send(msg.Caller, &messages.Traverse{Values: response})
		} else {
			fmt.Printf("\n ### TRAVERSE: error in traverse\n")
			context.Send(msg.Caller, &messages.Error{"Error while traversing"})
		}

	case *Delete:
		fmt.Println("\n ### DELETE: DELETING NODE")
		if state.LeftNode != nil {
			context.Send(state.LeftNode, &Delete{CurrentNode: state.LeftNode})
		}

		if state.RightNode != nil {
			context.Send(state.RightNode, &Delete{CurrentNode: state.RightNode})
		}
		context.Stop(msg.CurrentNode)
		fmt.Println("still running?")
	}
}

func sortKeys(Values map[int]string) []int {
	var keys []int
	for k := range Values {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}
