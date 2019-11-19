package main

import (
	"blatt-3-salcon/messages"
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"strconv"
	"sync"
)

//Global Variables
var (
	id          *int
	token       = flag.String("token", "", "tree token")
	pid         *actor.PID
	remotePid   *actor.PID
	wg          sync.WaitGroup //A WaitGroup waits for a collection of goroutines to finish.
	flagBind    = flag.String("bind", "localhost:8090", "Bind to Address")
	flagRemote  = flag.String("remote", "localhost:8091", "remote host:port")
	forceDelete = flag.Bool("no-preserve-tree", false, "force deletion of tree")
	rootContext *actor.RootContext
)

type TreeCliActor struct{}

func (state *TreeCliActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		switch msg.Type {
		case messages.CREATETREE:
			fmt.Printf("Id: %d\n", msg.Key)
			fmt.Printf("Token: %s\n", msg.Value)
			wg.Done()
		case messages.FIND:
			fmt.Printf("Value: %s\n", msg.Key)
			wg.Done()
		case messages.SUCCESS:
			wg.Done()
		case messages.ALLTREES:
			fmt.Println("ID's for Trees: " + msg.Value)
			wg.Done()
		}
	case *messages.Traverse:
		for i, pair := range msg.Values {
			fmt.Printf("{%d,%s}", pair.Key, pair.Value)
			if i < len(msg.Values)-1 {
				fmt.Printf(",")
			}
		}
		fmt.Printf("\n")
		wg.Done()
	case *messages.Error:
		fmt.Println(msg.Message + "\n")
		wg.Done()
	}
}

// Command Line Interface
func main() {
	fmt.Println("Hello Tree-CLI!")
	remote.Start(*flagBind)
	//siehe folie
	rootContext = actor.EmptyRootContext //initliaze empty root context
	props := actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1) // wait one goroutine
		return &TreeCliActor{}
	})
	pid = rootContext.Spawn(props) //starts actor after being created
	pidResp, err := remote.SpawnNamed(*flagRemote, "remote", "treeService", 0);
	if err != nil {
		//handle error
		panic(err)
	} else {
		remotePid = pidResp.Pid
		//handle commands
		switch flag.Args()[0] {
		case "newtree":
			newTree()
		case "insert":
			insert()
		case "search":
			search()
		case "remove":
			remove()
		case "delete":
			deleteTree()
		case "traverse":
			traverse()
		case "trees":
			trees()
		default:
			fmt.Println("Error....")
			return
		}
		wg.Wait()
	}
}

//function definitions for every command
func search() {
	if len(flag.Args()) != 1 || isNotValid(id, token) {
		handleError()
		return
	}
	tmp, _ := strconv.Atoi(flag.Args()[1])
	rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.FIND, Key: int32(tmp), Token: *token, Id: int32(*id)}, pid)
}

func insert() {
	if len(flag.Args()) != 2 || isNotValid(id, token) {
		handleError()
		return
	} else {
		tmp, _ := strconv.Atoi(flag.Args()[1])
		rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.ADD, Key: int32(tmp), Value: flag.Args()[2], Token: *token, Id: int32(*id)}, pid)
	}
}

func remove() {
	if len(flag.Args()) != 3 || isNotValid(id, token) {
		handleError()
		return
	}
	tmp, _ := strconv.Atoi(flag.Args()[1])
	rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.REMOVE, Key: int32(tmp), Token: *token, Id: int32(*id)}, pid)
}

func traverse() {
	if len(flag.Args()) > 1 || isNotValid(id, token) {
		handleError()
		return
	}
	rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.TRAVERSE, Token: *token, Id: int32(*id)}, pid)
}

func trees() {
	if len(flag.Args()) != 4 {
		handleError()
		return
	}
	rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.ALLTREES}, pid)
}

func newTree() {
	if len(flag.Args()) != 5 {
		handleError()
		return
	}
	tmp, _ := strconv.Atoi(flag.Args()[1])
	rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.CREATETREE, LeafSize: int32(tmp)}, pid)
}

func deleteTree() {
	if len(flag.Args()) != 6 || isNotValid(id, token) {
		handleError()
		return
	}
	tmp, _ := strconv.Atoi(flag.Args()[1])
	rootContext.RequestWithCustomSender(remotePid, &messages.Request{Type: messages.REMOVE, Key: int32(tmp), Token: *token, Id: int32(*id)}, pid)
}

func handleError() {
	fmt.Println("Syntax Error...")
	wg.Done()
}

//bool function for checking if id and token are given
func isNotValid(id *int, token *string) bool {
	if *id == -1 || *token == "" {
		return true
	}
	return false
}
