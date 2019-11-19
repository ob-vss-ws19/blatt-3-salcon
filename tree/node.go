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
		if (len(state.Data) < state.LeafSize || state.Data[msg.Key] != "") && state.LeftNode == nil && state.RightNode == nil {
			// add key to leaf
			if state.Data == nil {
				state.Data = make(map[int]string)
			}
			state.Data[msg.Key] = msg.Val
			fmt.Printf("added key: %d\n", msg.Key)

		} else if len(state.Data) == 0 && state.LeftNode != nil && state.RightNode != nil {
			// not a leaf
			if msg.Key <= state.MaxKeyVal {
				// add left
				context.Send(state.LeftNode, msg)
			} else {
				// add right
				context.Send(state.RightNode, msg)
			}
		} else if len(state.Data) == state.LeafSize && state.LeftNode == nil && state.RightNode == nil {
			// leaf full create new leafs
			fmt.Println("created new leafs")
			props := actor.PropsFromProducer(func() actor.Actor {
				return &Node{LeafSize: int(state.LeafSize)}
			})
			state.LeftNode = context.Spawn(props)

			state.RightNode = context.Spawn(props)

			// send values to leafs
			state.Data[msg.Key] = msg.Val
			keys := sortKeys(state.Data)
			state.MaxKeyVal = keys[(len(keys)/2)-1]
			fmt.Printf("left max key %d\n", state.MaxKeyVal)
			for _, key := range keys {
				if key <= state.MaxKeyVal {
					// add half left
					fmt.Printf("send %d left\n", key)
					context.Send(state.LeftNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
				} else {
					// add half right
					fmt.Printf("send %d right\n", key)
					context.Send(state.RightNode, &Add{Key: key, Val: state.Data[key]})
					delete(state.Data, key)
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
			fmt.Print("\n# FIND: Searching in map: ")
			fmt.Println(state.Data)

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
			fmt.Print("\n# FIND: Searching in map: ")
			fmt.Println(state.Data)

			foundData := state.Data[msg.Key]
			fmt.Printf("\n# FIND: Data found %s\n", foundData)

			if foundData != "" {
				delete(state.Data, msg.Key)
				fmt.Printf("deleted key %d\n", msg.Key)
			}
		} else {
			context.Send(msg.RequestFrom, &messages.Error{Message: "# FIND: Key not found"})
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
				fmt.Println("send to start")
				context.Send(tmp, msg)
				return
			}

			// if root is node create slices, set start to nil, add right node to remaining and forward
			msg.RemainingNodes = append(msg.RemainingNodes, state.RightNode)
			fmt.Println("send to left node from start")
			context.Send(state.LeftNode, msg)
			return
		} else if state.LeftNode != nil && state.RightNode != nil {
			// node is not leaf
			// while remaining nodes add right node to remaining and send to left node
			fmt.Printf("add right node send to left\n")
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
				fmt.Printf("appending %d\n", key)
				msg.Values = append(msg.Values, KeyValuePair{key, state.Data[key]})
			}
			next := msg.RemainingNodes[len(msg.RemainingNodes)-1]
			msg.RemainingNodes = msg.RemainingNodes[:len(msg.RemainingNodes)-1]
			fmt.Println("send to next node from leaf")
			context.Send(next, msg)
		} else if len(msg.RemainingNodes) == 0 && state.LeftNode == nil && state.RightNode == nil {
			// leaf with no remaining nodes to traverse

			var keys []int
			for key := range state.Data {
				keys = append(keys, int(key))
			}
			sort.Ints(keys)

			for _, key := range keys {
				fmt.Printf("appending %d in last leaf\n", key)
				msg.Values = append(msg.Values, KeyValuePair{key, state.Data[key]})
			}

			respon := make([]*messages.Response, 0)

			fmt.Println("send to caller")
			for _, pair := range msg.Values {
				respon = append(respon, &messages.Response{Value: pair.Value, Key: int32(pair.Key)})
			}
			context.Send(msg.Caller, &messages.Traverse{Values: respon})
		} else {
			fmt.Printf("error in traverse\n")
			context.Send(msg.Caller, &messages.Error{"Error while traversing"})
		}
	case *Delete:
		fmt.Println("stopping node")
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
