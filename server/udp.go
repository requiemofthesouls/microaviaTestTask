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

type udpServer struct {
	addr     string
	workers  int
	timeout  time.Duration
	wg       sync.WaitGroup
	listener *net.UDPConn
	stopChan chan os.Signal
	addrChan chan net.Addr
	svc      service.ITimeProtocolService
}

func newUDPTimeServer(addr string, workers int, timeout time.Duration, svc service.ITimeProtocolService) ITimeProtocolServer {
	return &udpServer{
		addr:     addr,
		workers:  workers,
		timeout:  timeout,
		stopChan: make(chan os.Signal, 1),
		addrChan: make(chan net.Addr),
		svc:      svc,
	}
}

func (s *udpServer) Run() (err error) {
	log.Println("starting UDP TCPServer on", s.addr)

	var udpAddr *net.UDPAddr
	if udpAddr, err = net.ResolveUDPAddr("udp4", s.addr); err != nil {
		return err
	}

	if s.listener, err = net.ListenUDP("udp", udpAddr); err != nil {
		return err
	}

	s.serve()
	return nil
}

func (s *udpServer) Stop() {
	s.stopChan <- syscall.SIGINT
}

func (s *udpServer) serve() {
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
	s.listener.Close()
	log.Println("waiting workers for done")
	s.wg.Wait()
}

func (s *udpServer) connListener(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		_, addr, err := s.listener.ReadFrom(nil)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Println("error to read datagramm: ", err)
				continue
			}
		}

		if addr != nil {
			log.Println("received packet from: ", addr.String())
			s.addrChan <- addr
		}
	}
}

func (s *udpServer) connConsumer(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case addr := <-s.addrChan:
			log.Println("got connection from:", addr.String())
			go s.writeResponse(addr)
		case <-ctx.Done():
			return
		}
	}
}

func (s *udpServer) writeResponse(addr net.Addr) {
	s.wg.Add(1)
	defer s.wg.Done()

	s.listener.SetWriteDeadline(time.Now().Add(s.timeout))
	var err error
	if _, err = s.listener.WriteTo(s.svc.GetBinarySeconds(), addr); err != nil {
		log.Println("fail to write response:", err)
	}
}
