package hub

import (
	"fmt"
	"sync"

	"github.com/hashicorp/memberlist"
)

// mlDelegateEvents will help keep track of nodes
// TODO: Implement this
type mlDelegateEvents struct {
	memberCount int
	mux         sync.Mutex
}

func (mde *mlDelegateEvents) Count() int {
	mde.mux.Lock()
	defer mde.mux.Unlock()

	return mde.memberCount
}

func (mde *mlDelegateEvents) add() {
	mde.mux.Lock()
	defer mde.mux.Unlock()

	mde.memberCount += 1
}

func (mde *mlDelegateEvents) remove() {
	mde.mux.Lock()
	defer mde.mux.Unlock()

	mde.memberCount -= 1
}

func (mde *mlDelegateEvents) NotifyJoin(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)

	mde.add()

	fmt.Println("Node has joined: "+node.String(), "ADDR:", hostPort)
}

func (mde *mlDelegateEvents) NotifyLeave(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)

	mde.remove()

	fmt.Println("Node has left: "+node.String(), "ADDR:", hostPort)
}

func (mde *mlDelegateEvents) NotifyUpdate(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	fmt.Println("Node was updated: "+node.String(), "ADDR:", hostPort)
}
