package main

import (
	"blatt-3-salcon/tree"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"sync"
	"time"
)

type TreeServiceActor struct {
}

func (state *TreeServiceActor) Receive(context actor.Context) {

}

func getTreeSericeActor() actor.Actor {
	fmt.Printf("# Tree-Service-Actor is ready\n")
	return &TreeServiceActor{}
}

func newToken() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func main() {
	var bind = flag.String("bind", "localhost:8088", "Bind to address")
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	flag.Parse()
	remote.Start(*bind)
	remote.Register("TreeServiceActor", actor.PropsFromProducer(getTreeSericeActor))
}
