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

func (state *TreeServiceActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Request:
		switch msg.Type {
		case messages.CREATETREE:
			if msg.Id <= 0 {
				context.Respond(&messages.Error{"Leaf Size should be at least 1"})
			}
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

var bind = flag.String("bind", "localhost:8088", "Bind to address")

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	flag.Parse()
	remote.Start(*bind)
	remote.Register("TreeServiceActor", actor.PropsFromProducer(getTreeServiceActor))
}
