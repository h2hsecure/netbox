package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nadoo/ipset"
	"github.com/otterize/go-procnet/procnet"
)

func main() {
	// must call Init first
	if err := ipset.Init(); err != nil {
		log.Printf("error in ipset Init: %s", err)
		return
	}

	// default is ipv4 without timeout
	ipset.Destroy("myset")
	ipset.Create("myset")

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				socks, err := procnet.TCPSocks()
				if err != nil {
					panic(err)
				}
				ipmap := make(map[string]int)
				for _, sock := range socks {
					// fmt.Printf("local ip: %s local port: %d remote IP: %s remote port: %d state: %s",
					// 	sock.LocalAddr.IP, sock.LocalAddr.Port, sock.RemoteAddr.IP, sock.RemoteAddr.Port, sock.State)
					if _, has := ipmap[sock.RemoteAddr.IP.String()]; has {
						ipmap[sock.RemoteAddr.IP.String()]++
					} else {
						ipmap[sock.RemoteAddr.IP.String()] = 1
					}

				}

				for k, v := range ipmap {
					fmt.Printf("ip: %s, count: %d\n", k, v)
					ipset.Add("myset", k)
				}
				time.Sleep(5 * time.Second)
			}
		}
	}(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done
	cancel()

	// //ipset.Flush("myset")

	// // ipv6 and timeout example
	// // ipset create myset6 hash:net family inet6 timeout 60
	// ipset.Create("myset6", ipset.OptIPv6(), ipset.OptTimeout(60))
	// ipset.Flush("myset6")

	// ipset.Add("myset6", "2022::1", ipset.OptTimeout(10))
	// ipset.Add("myset6", "2022::1/32")
}
