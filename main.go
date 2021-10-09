package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/benbjohnson/clock"
	"timeProtocol/client"
	"timeProtocol/server"
	"timeProtocol/service"
)

func main() {
	var (
		tcp       = flag.Bool("tcp", false, "use tcp (default)")
		udp       = flag.Bool("udp", false, "use udp")
		workers   = flag.Int("workers", 2, "workers")
		timeout   = flag.Int("timeout", 10, "write timeout (in seconds)")
		port      = flag.Int("port", 8080, "port")
		hostname  = flag.String("h", "", "hostname")
		useClient = flag.Bool("c", false, "use client")
	)
	flag.Parse()

	var (
		addr         = fmt.Sprintf("%s:%d", *hostname, *port)
		writeTimeout = time.Duration(*timeout) * time.Second
		s            server.ITimeProtocolServer
		err          error
		proto        server.Protocol
	)

	if *tcp {
		proto = server.TCP
	} else if *udp {
		proto = server.UDP
	} else {
		proto = server.TCP
	}

	if *useClient {
		var c client.ITimeProtocolClient
		if c, err = client.NewTimeProtocolClient(proto, addr, writeTimeout); err != nil {
			log.Fatal(err)
		}
		var seconds uint32
		if seconds, err = c.Get(); err != nil {
			log.Fatal(err)
		}

		log.Println("current time is:", seconds)
		log.Println("converted time:", service.TPStartTime.Add(time.Second*time.Duration(seconds)))
		return
	}

	var svc = service.NewTimeProtocolService(clock.New())

	if s, err = server.NewTimeServer(proto, addr, *workers, writeTimeout, svc); err != nil {
		log.Fatal(err)
	}

	if err = s.Run(); err != nil {
		log.Fatal(err)
	}

}
