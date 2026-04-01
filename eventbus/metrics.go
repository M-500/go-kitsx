package eventbus

import "sync/atomic"

type Metrics struct {
	Published     uint64
	PublishFailed uint64
	Delivered     uint64
	HandlerFailed uint64
	Retries       uint64
	Panics        uint64
	Dropped       uint64
}

type metricsRegistry struct {
	published     atomic.Uint64
	publishFailed atomic.Uint64
	delivered     atomic.Uint64
	handlerFailed atomic.Uint64
	retries       atomic.Uint64
	panics        atomic.Uint64
	dropped       atomic.Uint64
}

func (m *metricsRegistry) snapshot() Metrics {
	return Metrics{
		Published:     m.published.Load(),
		PublishFailed: m.publishFailed.Load(),
		Delivered:     m.delivered.Load(),
		HandlerFailed: m.handlerFailed.Load(),
		Retries:       m.retries.Load(),
		Panics:        m.panics.Load(),
		Dropped:       m.dropped.Load(),
	}
}
