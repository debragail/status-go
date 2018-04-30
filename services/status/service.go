package status

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/status-im/status-go/geth/account"
)

// Make sure that Service implements node.Service interface.
var _ node.Service = (*Service)(nil)

// Service represents out own implementation of status status operations.
type Service struct {
	am *account.Manager
	w  *whisper.Whisper
}

// New returns a new Service.
func New(w *whisper.Whisper) *Service {
	return &Service{w: w}
}

// Protocols returns a new protocols list. In this case, there are none.
func (s *Service) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

// APIs returns a list of new APIs.
func (s *Service) APIs() []rpc.API {

	return []rpc.API{
		{
			Namespace: "status",
			Version:   "1.0",
			Service:   NewAPI(s),
			Public:    false,
		},
	}
}

// SetAccountManager sets account manager for the API calls.
func (s *Service) SetAccountManager(a *account.Manager) {
	s.am = a
}

// Start is run when a service is started.
// It does nothing in this case but is required by `node.Service` interface.
func (s *Service) Start(server *p2p.Server) error {
	return nil
}

// Stop is run when a service is stopped.
// It does nothing in this case but is required by `node.Service` interface.
func (s *Service) Stop() error {
	return nil
}
