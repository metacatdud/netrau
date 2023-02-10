package hub

import (
	"context"
	"fmt"
	"github.com/hashicorp/memberlist"
	"gitlab.com/macroscope-lab/atomika"
	"gitlab.com/macroscope-lab/atomika/log"
	"gitlab.com/macroscope-lab/atomika/types"
	"time"
)

type Hub interface {
	atomika.Service
}

type hub struct {
	options Options
	log     log.Ctx

	memberListCfg     *memberlist.Config
	list              *memberlist.Memberlist
	broadcast         *memberlist.TransmitLimitedQueue
	broadcastDelegate *mlBroadcast
}

func (h hub) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	// Create Memeberlist here
	memberList, err := memberlist.Create(h.memberListCfg)
	if err != nil {
		return err
	}

	local := memberList.LocalNode()
	joinList := []string{
		fmt.Sprintf("%s:%d", local.Addr.To4().String(), local.Port),
	}

	//TODO: add proper ip:port validation
	if h.options.Join != "" {
		joinList = append(joinList, h.options.Join)
	}

	if _, err = memberList.Join(joinList); err != nil {
		return err
	}

	h.list = memberList

	// DEV Purposes only
	tick := time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-tick.C:
				msg := &Message{
					Payload: fmt.Sprintf("NodeID: %s | M:%s", h.memberListCfg.Name, types.GenUUIDv4().String()),
				}

				// This broadcast should be used in business logic
				h.broadcast.QueueBroadcast(msg)
				h.log.Debug("send", "data", msg)

			case data := <-h.broadcastDelegate.Message():
				// If the message cannot be parsed move to the next one
				msg, ok := ParseMessage(data)
				if ok != true {
					continue
				}

				h.log.Debug("received", "data", msg)
			}
		}
	}()

	h.log.Info("HUB started", "port", h.list.LocalNode().Port)

	for {
		select {
		case <-ctx.Done():
			shutdown(context.Background(), memberList)
			return nil
		case err = <-errChan:
			log.Error(fmt.Sprintf("errChan %s", err))
			return nil
		}
	}
}

func New(opts ...Option) (Hub, error) {
	return newMemberService(setOptions(opts...))
}

func newMemberService(opts Options) (Hub, error) {
	lg := log.WithContext("service", "Hub")

	// Config memberlist pkg
	memberListCfg := memberlist.DefaultWANConfig()
	memberListCfg.Name = types.GenUUIDv4().String()

	if opts.BindAddr != "" {
		memberListCfg.BindAddr = opts.BindAddr
	}

	memberListCfg.BindPort = opts.BindPort
	memberListCfg.AdvertisePort = memberListCfg.BindPort

	// Attach custom event delegation handler
	mde := &mlDelegateEvents{}
	memberListCfg.Events = mde

	// Attach custom broadcast handler
	msgCh := make(chan []byte)
	broadcast := &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return mde.Count()
		},
		RetransmitMult: opts.ResendLimit,
	}

	broadcastDelegate := &mlBroadcast{
		msgCh:      msgCh,
		broadcasts: broadcast,
	}

	memberListCfg.Delegate = broadcastDelegate

	h := &hub{
		options:           opts,
		log:               lg,
		memberListCfg:     memberListCfg,
		broadcast:         broadcast,
		broadcastDelegate: broadcastDelegate,
	}

	return h, nil
}

func shutdown(ctx context.Context, memberList *memberlist.Memberlist) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := memberList.Leave(5 * time.Second); err != nil {
		log.Error(fmt.Sprintf("node did not left correctly: %s", err.Error()), "node", memberList.LocalNode().String())
	}

	time.Sleep(2 * time.Second)
	log.Info("node left")
}
