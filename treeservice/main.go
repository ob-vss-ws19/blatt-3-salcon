package main

import (
	"blatt-3-salcon/messages"
	"blatt-3-salcon/tree"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"sync"
)

type TreeServiceActor struct {
}

var createdID = 1
var alltrees = make(map[int]*actor.PID)
var tokens = make(map[int]string)

// Kümmert sich darum, dass die Funktionalitäen
func (state *TreeServiceActor) Receive(context actor.Context) {
	fmt.Println(context.Message())
	switch msg := context.Message().(type) {
	case *messages.Request:
		fmt.Println("hi")
		switch msg.Type {

		case messages.CREATETREE:
			fmt.Println("createtree")

			if msg.LeafSize <= 0 {
				context.Respond(&messages.Error{"Leaf Size should be at least 1"})
			}
			fmt.Println(msg.Id)

			//Neue ID erhalten für Node
			newid := createdID
			createdID++

			props := actor.PropsFromProducer(func() actor.Actor {
				return &tree.Node{LeafSize: int(msg.LeafSize)}
			})
			newToken := newToken()

			pid := context.Spawn(props)
			alltrees[newid] = pid
			tokens[newid] = newToken

			context.Respond(&messages.Response{Key: int32(newid), Value: newToken, Type: messages.CREATETREE})

		case messages.FIND:

		case messages.ADD:

		case messages.REMOVE:

		case messages.DELETE:

		case messages.TRAVERSE:

		case messages.ALLTREES:

		}

	}
}

func getTreeServiceActor() actor.Actor {
	fmt.Printf("# Tree-Service-Actor is ready\n")
	return &TreeServiceActor{}
}

// Generates new Token, 4 Bytes long
func newToken() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

var bind = flag.String("bind", "localhost:8093", "Bind to address")

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	flag.Parse()
	remote.Start(*bind)
	remote.Register("treeService", actor.PropsFromProducer(getTreeServiceActor))
}
