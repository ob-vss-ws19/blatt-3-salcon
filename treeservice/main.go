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
var alltrees = make(map[int32]map[string]*actor.PID)

// Kümmert sich darum, dass die Funktionalitäen
func (state *TreeServiceActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Request:
		switch msg.Type {

		case messages.CREATETREE:
			fmt.Printf("# CREATETREE: Created Tree with ID %d\n", createdID)

			if msg.LeafSize <= 0 {
				context.Respond(&messages.Error{"Leaf Size should be at least 1"})
			}
			//Neue ID erhalten für Node
			newid := createdID
			createdID++

			props := actor.PropsFromProducer(func() actor.Actor {
				return &tree.Node{LeafSize: int(msg.LeafSize)}
			})
			newToken := getNewToken()

			pid := context.Spawn(props)
			alltrees[int32(newid)] = make(map[string]*actor.PID)
			alltrees[int32(newid)][newToken] = pid

			context.Respond(&messages.Response{Key: int32(newid), Value: newToken, Type: messages.CREATETREE})

		case messages.FIND:
			if pid := pidAccess(msg.Id, msg.Token); pid != nil {
				context.Send(pid, &tree.Find{Key: int(msg.Key), RequestFrom: context.Sender()})
				context.Respond(&messages.Response{Value: msg.Value, Type: messages.FIND})
			} else {
				accessDenied(context, context.Sender())
			}
		case messages.ADD:
			if pid := pidAccess(msg.Id, msg.Token); pid != nil {
				context.Send(pid, &tree.Add{Key: int(msg.Key), Val: msg.Value})
				context.Respond(&messages.Response{Type: messages.SUCCESS})
			} else {
				accessDenied(context, context.Sender())
			}

		case messages.REMOVE:
			if pid := pidAccess(msg.Id, msg.Token); pid != nil {
				context.Send(pid, &tree.Remove{Key: int(msg.Key)})
				context.Respond(&messages.Response{Type: messages.SUCCESS})
			} else {
				accessDenied(context, context.Sender())
			}

		case messages.TRAVERSE:
			if pid := pidAccess(msg.Id, msg.Token); pid != nil {
				context.Send(pid, &tree.Traverse{Caller: context.Sender(), Start: pid})
			} else {
				accessDenied(context, context.Sender())
			}
		case messages.DELETE:
			if pid := pidAccess(msg.Id, msg.Token); pid != nil {
				context.Send(pid, &tree.Delete{CurrentNode: pid})
				delete(alltrees, msg.Id)
				context.Respond(&messages.Response{Type: messages.SUCCESS})
			} else {
				accessDenied(context, context.Sender())
			}
		case messages.ALLTREES:
			tmp := ""
			for k := range alltrees {
				tmp += fmt.Sprintf("%d", k) + ", "
			}
			if len(tmp) > 2 {
				tmp = tmp[:len(tmp)-2]
			}

			context.Respond(&messages.Response{Type: messages.ALLTREES, Value: tmp})
		}

	}
}

func pidAccess(Id int32, token string) *actor.PID {
	pid := alltrees[Id][token]
	if pid == nil {
		return nil
	}
	return pid
}

func getTreeServiceActor() actor.Actor {
	fmt.Printf("# Tree-Service-Actor is ready\n")
	return &TreeServiceActor{}
}

// Generates new Token, 4 Bytes long
func getNewToken() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func accessDenied(context actor.Context, pid *actor.PID) {
	context.Send(pid, &messages.Error{Message: "Access Denied: Wrong token or id"})
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
