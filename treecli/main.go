package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/remotewatch/messages"
)

type TreeCliActor struct {}

func (state *TreeCliActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.HelloRequest:
		fmt.Println(msg);

		//TODO
	}
}

//Global Variables
var (
	id        *int
	token     int
	pid       *actor.PID
	remotePid *actor.PID
)

// Command Line Interface
func main() {
	fmt.Println("Hello Tree-CLI!")
	//Using Flags for Input
	//Define Flags
	flag.Usage = func() {
		fmt.Println("This is not helpful")
	}

	flag.Parse() //after defining all flags, call is necessary




	fmt.Println("Bye Tree-CLI!")
}
