package tree

import (
	cr "crypto/rand"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ws19/blatt-3-salcon/messages"
	mr "math/rand"
	"sync"
	"testing"
	"time"
)

type TestActor struct {
	t       *testing.T
	wg      *sync.WaitGroup
	indices []int
}

const LEAFS = 5
const NUMBEROFVALUES = 10000

var values = make([]KeyValuePair, 0)

func (state *TestActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Traverse:
		i := 0
		for _, k := range state.indices {
			if msg.Values[i].Value != values[k].Value || int(msg.Values[i].Key) != values[k].Key {
				fmt.Printf("should be: %d %s but is %d %s\n", msg.Values[i].Key, msg.Values[i].Value, values[k].Key, values[k].Value)
				state.t.Error()
			}
			i++
		}
		state.wg.Done()
	case *messages.Response:
		switch msg.Type {
		case messages.FIND:
			if int(msg.Key) != values[state.indices[0]].Key || msg.Value != values[state.indices[0]].Value {
				fmt.Printf("should be: %d %s but is %d %s\n", msg.Key, msg.Value, values[state.indices[0]].Key, values[state.indices[0]].Value)
				state.t.Error()
			}
			state.wg.Done()
		}

	}
}

func createRandValues() {

	values = make([]KeyValuePair, 0)

	pairs := make(map[int]string)

	for i := 0; i < NUMBEROFVALUES; i++ {
		pairs[int(mr.Int31n(NUMBEROFVALUES*10))] = newToken()
	}

	for _, v := range sortKeys(pairs) {
		values = append(values, KeyValuePair{Key: v, Value: pairs[v]})
	}

}

func TestAdd(t *testing.T) {
	createRandValues()
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: LEAFS}
	})
	tree := context.Spawn(props)
	var wg sync.WaitGroup

	indices := make(map[int]string)
	for range values {
		tmp := mr.Intn(len(values))
		indices[tmp] = ""
		context.Send(tree, &Add{Key: values[tmp].Key, Val: values[tmp].Value})
	}
	time.Sleep(1 * time.Second)

	props = actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &TestActor{t, &wg, sortKeys(indices)}
	})
	testAct := context.Spawn(props)
	context.Send(tree, &Traverse{Caller: testAct, Start: tree})
	time.Sleep(1 * time.Second)
	wg.Wait()
}

func TestFind(t *testing.T) {
	createRandValues()

	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: LEAFS}
	})
	tree := context.Spawn(props)
	var wg sync.WaitGroup

	indices := make(map[int]string)
	for range values {
		tmp := mr.Intn(len(values))
		indices[tmp] = ""
		context.Send(tree, &Add{Key: values[tmp].Key, Val: values[tmp].Value})
	}
	time.Sleep(2 * time.Second)

	for k := range indices {
		props = actor.PropsFromProducer(func() actor.Actor {
			wg.Add(1)
			return &TestActor{t, &wg, []int{k}}
		})
		testAct := context.Spawn(props)
		context.Send(tree, &Find{RequestFrom: testAct, Key: values[k].Key})
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
}

func newToken() string {
	b := make([]byte, 4)
	_, _ = cr.Read(b)
	return fmt.Sprintf("%x", b)
}

func TestDelete(t *testing.T) {
	createRandValues()

	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: LEAFS}
	})
	tree := context.Spawn(props)
	var wg sync.WaitGroup

	indices := make(map[int]string)
	for range values {
		tmp := mr.Intn(len(values))
		indices[tmp] = ""
		context.Send(tree, &Add{Key: values[tmp].Key, Val: values[tmp].Value})
	}
	time.Sleep(2 * time.Second)

	for k := range indices {
		if mr.Int31n(100) < 50 {
			delete(indices, k)
			context.Send(tree, &Remove{Key: values[k].Key})
		}
	}

	time.Sleep(2 * time.Second)

	props = actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &TestActor{t, &wg, sortKeys(indices)}
	})
	testAct := context.Spawn(props)
	context.Send(tree, &Traverse{Caller: testAct, Start: tree})
	wg.Wait()
	time.Sleep(1 * time.Second)
}
