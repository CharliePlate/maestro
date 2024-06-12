package maestro

import (
	"context"
	"log/slog"
	"net"
	"os"
	"sync"
)

type Server struct {
	Logger   slog.Logger
	Listener net.Listener
	Opts     ServerOpts
}

type ServerOpts struct {
	Addr string
	Port int
}

func NewServer(l net.Listener, opts ServerOpts) *Server {
	s := &Server{
		Opts:     opts,
		Logger:   *slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
		Listener: l,
	}

	return s
}

func (s *Server) Start(ctx context.Context) {
	s.Logger.Info("starting server", slog.String("addr", s.Opts.Addr), slog.Int("port", s.Opts.Port))
	s.Logger.Info("server started", slog.String("addr", s.Opts.Addr), slog.Int("port", s.Opts.Port))

	go func() {
		<-ctx.Done()
		s.Logger.Info("stopping server")
		s.Listener.Close()
	}()
}

type Peer struct{}

type PeerMap struct {
	peers map[string]*Peer
	mutex sync.RWMutex
}

func NewPeerMap() *PeerMap {
	return &PeerMap{
		peers: make(map[string]*Peer),
		mutex: sync.RWMutex{},
	}
}

func (pm *PeerMap) AddPeer(connID string, peer *Peer) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.peers[connID] = peer
}

func (pm *PeerMap) RemovePeer(connID string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	delete(pm.peers, connID)
}

func (pm *PeerMap) GetPeer(connID string) *Peer {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.peers[connID]
}
