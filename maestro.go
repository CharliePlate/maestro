package maestro

import (
	"errors"
	"log/slog"
	"net"
	"time"
)

type Config struct {
	Logger slog.Logger
}

type ConnHandler interface {
	HandleConn(conn net.Conn) error
}

type Maestro struct {
	Config      Config
	Listener    net.Listener
	ConnHandler ConnHandler
	Queues      []Queue
}

func (m *Maestro) Listen() error {
	for {
		conn, err := m.Listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := m.ConnHandler.HandleConn(conn); err != nil {
				m.Config.Logger.Error("error handling conn", slog.String("error", err.Error()))
				conn.Close()
				return
			}
		}()
	}
}

type TCPConnHandler struct {
	KeepAlivePeriod int
}

func (t *TCPConnHandler) HandleConn(conn net.Conn) error {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return errors.New("expected tcp conn")
	}

	err := tcpConn.SetKeepAlive(true)
	if err != nil {
		return err
	}

	err = tcpConn.SetKeepAlivePeriod(time.Duration(t.KeepAlivePeriod))
	if err != nil {
		return err
	}

	return nil
}
