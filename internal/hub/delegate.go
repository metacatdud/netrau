package hub

import "github.com/hashicorp/memberlist"

// Memberlist package delegation extension

type mlBroadcast struct {
	msgCh      chan []byte
	broadcasts *memberlist.TransmitLimitedQueue
}

func (m *mlBroadcast) Message() chan []byte {
	return m.msgCh
}

func (m *mlBroadcast) NodeMeta(limit int) []byte {
	return []byte("")
}

func (m *mlBroadcast) NotifyMsg(msg []byte) {
	m.msgCh <- msg
}

func (m *mlBroadcast) GetBroadcasts(overhead, limit int) [][]byte {
	return m.broadcasts.GetBroadcasts(overhead, limit)
}

func (m *mlBroadcast) LocalState(join bool) []byte {
	return []byte("")
}

func (m *mlBroadcast) MergeRemoteState(buf []byte, join bool) {
	// noop
}
