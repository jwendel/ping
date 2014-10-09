// Pinger
package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s {hostname} [hostnames...]\n", os.Args[0])
		os.Exit(1)
	}

	p := NewPinger()
	for i := 1; i < len(os.Args); i++ {
		host := os.Args[i]
		p.AddHost(host)
	}

	err := p.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case ping := <-p.Results:
			fmt.Println("avg ping:", ping.Avg)
		}
	}

}
