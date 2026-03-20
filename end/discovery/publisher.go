package discovery

import "github.com/grandcat/zeroconf"

type Publisher struct {
	server *zeroconf.Server
}

func NewPublisher(manager *Manager) (*Publisher, error) {
	server, err := zeroconf.Register(
		manager.MDNSInstanceName(),
		ServiceType,
		ServiceDomain,
		manager.Port(),
		manager.MDNSTextRecords(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{server: server}, nil
}

func (p *Publisher) Shutdown() {
	if p == nil || p.server == nil {
		return
	}
	p.server.Shutdown()
}
