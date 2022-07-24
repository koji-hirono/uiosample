package main

import (
	"context"
	"os"
	"sync"

	"uiosample/worker"
)

type Server struct {
	port1 *Port
	port2 *Port
}

func NewServer(port1, port2 *Port) *Server {
	s := new(Server)
	s.port1 = port1
	s.port2 = port2
	return s
}

func (s *Server) Serve(sig chan os.Signal) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	pkts := make([][]byte, 32, 32)
	worker.Worker(ctx, &wg, 2, func() {
		// N3
		n := s.port1.RxBurst(pkts)
		for i := 0; i < n; i++ {
			s.procN3(pkts[i])
		}
		// N6
		n = s.port2.RxBurst(pkts)
		for i := 0; i < n; i++ {
			s.procN6(pkts[i])
		}
	})
	defer wg.Wait()

	<-sig
	cancel()
}
