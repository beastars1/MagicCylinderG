package MagicCylinderG

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestSocks(t *testing.T) {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9999", nil, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	if resp, err := httpClient.Get("https://www.qq.com"); err != nil {
		t.Fatal(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		t.Logf("body:%s\n", body)
	}
}

func BenchmarkSocks(b *testing.B) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9999", nil, proxy.Direct)
	if err != nil {
		log.Fatal("can't connect to the proxy:", err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateSocksRequest(httpClient, "https://www.qq.com")
	}
}

func CreateSocksRequest(httpClient *http.Client, url string) {
	if _, err := httpClient.Get(url); err != nil {
		log.Fatal(err)
	}
}
