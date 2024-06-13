package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"paepcke.de/dnsresolver"
)

var (
	resolverOnce = sync.OnceValue(func() *dnsresolver.Resolver {
		// ResolverAutoが`/etc/resolve.conf`をLstatして存在チェックする。
		// wslやdockerなど一部環境はこれがsymlinkなので、その環境ではうまく動かない
		// そのため手動であらかじめStatしておいて、存在していれば明示的にそこを読み込ませる。
		if _, err := os.Stat("/etc/resolv.conf"); err == nil {
			return dnsresolver.ResolverResolvConf()
		}
		return dnsresolver.ResolverAuto()
	})
)

func dialer(fallbackOnly bool) func(ctx context.Context, network, addr string) (net.Conn, error) {
	d := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var (
			host = addr
			port = ""
			err  error
		)
		if strings.Contains(addr, ":") {
			host, port, err = net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
		}

		raw, _ := strings.CutPrefix(host, "[")
		raw, _ = strings.CutSuffix(raw, "]")
		_, err = netip.ParseAddr(raw)
		if err == nil {
			// turns out it is already a raw IP addr.
			return d.DialContext(ctx, network, addr)
		}

		if fallbackOnly {
			_, err := net.LookupHost(host)
			if err == nil {
				return d.DialContext(ctx, network, addr)
			}
		}

		r := resolverOnce()
		resolved, err := r.LookupAddrs(host, []uint16{dns.TypeA, dns.TypeAAAA})
		if err != nil {
			return nil, err
		}
		if resolved[0].Is4() {
			addr = resolved[0].String()
		} else {
			addr = "[" + resolved[0].String() + "]"
		}
		if port != "" {
			addr += ":" + port
		}
		return d.DialContext(ctx, network, addr)
	}
}

func main() {
	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	t.DialContext = dialer(false)

	client := &http.Client{
		Transport: t,
	}
	var err error
	_, err = client.Get("https://www.google.com")
	fmt.Printf("err = %v\n", err)
	_, err = client.Get("https://142.250.206.228:443")
	fmt.Printf("err = %v\n", err)
	_, err = http.DefaultClient.Get("https://142.250.206.228:443")
	fmt.Printf("err = %v\n", err)
	_, err = client.Get("https://[2404:6800:400a:804::2004]:443")
	fmt.Printf("err = %v\n", err)
}
