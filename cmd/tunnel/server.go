package main

import (
	"context"
	"net"
	"os"
	"sync"

	"uiosample/worker"
)

type Server struct {
	GTPAddr net.IP
	GTPPort uint16
	port1   *Port
	port2   *Port
	pdrtbl  *PDRTable
	fartbl  *FARTable
	qertbl  *QERTable
}

func NewServer(port1, port2 *Port) *Server {
	s := new(Server)
	s.GTPAddr = net.IPv4(30, 30, 0, 1).To4()
	s.GTPPort = 2152
	s.port1 = port1
	s.port2 = port2
	s.pdrtbl = NewPDRTable()
	s.fartbl = NewFARTable()
	s.qertbl = NewQERTable()
	// uplink(N3)
	s.qertbl.Put(1, 1, &QER{
		SEID: 1,
		ID:   1,
		QFI:  9,
	})
	s.fartbl.Put(1, 1, &FAR{
		SEID:   1,
		ID:     1,
		Action: ApplyActionFORW,
	})
	s.pdrtbl.Put(1, 1, &PDR{
		SEID: 1,
		ID:   1,
		PDI: &PDI{
			UEAddr: net.IPv4(60, 60, 0, 2).To4(),
			FTEID: &FTEID{
				TEID:     78,
				GTPuAddr: net.IPv4(30, 30, 0, 1).To4(),
			},
		},
		FARID: 1,
		QERID: 1,
	})
	// downlink(N6)
	s.qertbl.Put(1, 2, &QER{
		SEID: 1,
		ID:   1,
		QFI:  9,
	})
	s.fartbl.Put(1, 2, &FAR{
		SEID:   1,
		ID:     2,
		Action: ApplyActionFORW,
		Param: &ForwardParam{
			Creation: &HeaderCreation{
				TEID:     87,
				PeerAddr: net.IPv4(30, 30, 0, 2).To4(),
				Port:     2152,
			},
		},
	})
	s.pdrtbl.Put(1, 2, &PDR{
		SEID: 1,
		ID:   2,
		PDI: &PDI{
			UEAddr: net.IPv4(60, 60, 0, 2).To4(),
		},
		FARID: 2,
		QERID: 2,
	})
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
