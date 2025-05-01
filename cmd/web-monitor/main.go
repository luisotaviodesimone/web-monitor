package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func ping(host string, count int) {
	url := fmt.Sprintf("%s", host)
	pinger, err := probing.NewPinger(url)
	if err != nil {
		panic(err)
	}
	pinger.Count = count

	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	slog.Info("Packets:",
		slog.Int("sent", stats.PacketsSent),
		slog.Int("received", stats.PacketsRecv),
		slog.Float64("packet loss", stats.PacketLoss),
	)

	slog.Info("Approximate round trip times in milli-seconds",
		slog.Float64("minimum", stats.MinRtt.Seconds()*1000),
		slog.Float64("average", stats.AvgRtt.Seconds()*1000),
		slog.Float64("maximum", stats.MaxRtt.Seconds()*1000),
		slog.Float64("standardDeviation", stats.StdDevRtt.Seconds()*1000),
	)
}

func httpPing(host string, times int) {
	url := fmt.Sprintf("http://%s", host)
	counterChan := make(chan int)
	counter := 0
	var totalTime time.Duration

	headers := make(http.Header)
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	httpCaller := probing.NewHttpCaller(url,
		probing.WithHTTPCallerCallFrequency(time.Second),
		probing.WithHTTPCallerHeaders(headers),
		probing.WithHTTPCallerOnResp(func(suite *probing.TraceSuite, info *probing.HTTPCallInfo) {
			requestTime := suite.GetGeneralEnd().Sub(suite.GetGeneralStart())
			fmt.Printf("got resp, status code: %d, latency: %s\n",
				info.StatusCode,
				requestTime,
			)
			counter++
			totalTime += requestTime
			counterChan <- counter
		}),
	)

	go httpCaller.Run()

	for count := range counterChan {
		if count >= times {
			httpCaller.Stop()
			// avgTime := time.Duration(totalTime / count)
			slog.Info("Average time in miliseconds", slog.Float64("average", float64(totalTime.Milliseconds())/float64(count)))
			break
		}
	}
}
func main() {
	// hosts := []string{"google.com", "rnp.br", "youtube.com"}
	// for _, host := range hosts {
	// 	ping(host, 11)
	// 	httpPing(host, 11)
	// }
	httpPing("rnp.br", 11)
}
