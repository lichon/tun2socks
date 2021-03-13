// Package proxy provides implementations of proxy protocols.
package proxy

import (
	"context"
	"net"
	"time"

	M "github.com/xjasonlyu/tun2socks/constant"
	"github.com/xjasonlyu/tun2socks/proxy/proto"

	"go.uber.org/atomic"
)

const (
	tcpConnectTimeout = 5 * time.Second
)

var (
	_defaultDialer atomic.Value
)

type Dialer interface {
	DialContext(context.Context, *M.Metadata) (net.Conn, error)
	DialUDP(*M.Metadata) (net.PacketConn, error)
}

type Proxy interface {
	Dialer
	Addr() string
	Proto() proto.Proto
}

func init() {
	_defaultDialer.Store(&Base{})
}

// SetDialer sets default Dialer.
func SetDialer(d Dialer) {
	_defaultDialer.Store(d)
}

// Dial uses default Dialer to dial TCP.
func Dial(metadata *M.Metadata) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpConnectTimeout)
	defer cancel()
	return _defaultDialer.Load().(Dialer).DialContext(ctx, metadata)
}

// DialContext uses default Dialer to dial TCP with context.
func DialContext(ctx context.Context, metadata *M.Metadata) (net.Conn, error) {
	return _defaultDialer.Load().(Dialer).DialContext(ctx, metadata)
}

// DialUDP uses default Dialer to dial UDP.
func DialUDP(metadata *M.Metadata) (net.PacketConn, error) {
	return _defaultDialer.Load().(Dialer).DialUDP(metadata)
}
