package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"timeProtocol/service"
)

type tcpServer struct {
	addr     string
	workers  int
	timeout  time.Duration
	wg       sync.WaitGroup
	listener net.Listener
	stopChan chan os.Signal
	connChan chan net.Conn
	svc      service.ITimeProtocolService
}

func newTCPTimeServer(
	addr string,
	workers int,
	timeout time.Duration,
	svc service.ITimeProtocolService,
) ITimeProtocolServer {
	return &tcpServer{
		addr:     addr,
		workers:  workers,
		timeout:  timeout,
		stopChan: make(chan os.Signal, 1),
		connChan: make(chan net.Conn),
		svc:      svc,
	}
}

func (s *tcpServer) Run() (err error) {
	log.Println("starting TCP TCPServer on", s.addr)

	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		return err
	}

	s.serve()
	return nil
}

func (s *tcpServer) Stop() {
	s.stopChan <- syscall.SIGINT
}

func (s *tcpServer) serve() {
	signal.Notify(
		s.stopChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)

	ctx, cancel := context.WithCancel(context.Background())
	go s.connListener(ctx)

	for i := 1; i <= s.workers; i++ {
		log.Println("starting connConsumer: ", i)
		go s.connConsumer(ctx)
	}

	log.Println("ready to accept connections")
	sig := <-s.stopChan
	log.Println("stopping, got:", sig)
	cancel()

	log.Println("closing listener")
	if err := s.listener.Close(); err != nil {
		log.Println(err)
	}

	log.Println("waiting workers for done")
	s.wg.Wait()
}

func (s *tcpServer) connListener(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Println(err)
				continue
			}
		}

		conn.SetWriteDeadline(time.Now().Add(s.timeout))
		s.connChan <- conn
	}
}

func (s *tcpServer) connConsumer(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case conn := <-s.connChan:
			log.Println("got connection from:", conn.RemoteAddr())
			go s.writeResponse(conn)
		case <-ctx.Done():
			return
		}
	}
}

func (s *tcpServer) writeResponse(conn net.Conn) {
	s.wg.Add(1)
	defer s.wg.Done()

	var err error
	if _, err = conn.Write(s.svc.GetBinarySeconds()); err != nil {
		log.Println(err)
	}

	if err = conn.Close(); err != nil {
		log.Println(err)
	}
}
