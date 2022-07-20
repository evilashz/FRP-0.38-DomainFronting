package mux

import (
	"bufio"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func runHttpSvr(ln net.Listener) *httptest.Server {
	svr := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("http service"))
	}))
	svr.Listener = ln
	svr.Start()
	return svr
}

func runHttpsSvr(ln net.Listener) *httptest.Server {
	svr := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("https service"))
	}))
	svr.Listener = ln
	svr.StartTLS()
	return svr
}

func runEchoSvr(ln net.Listener) {
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			rd := bufio.NewReader(conn)
			data, err := rd.ReadString('\n')
			if err != nil {
				return
			}
			conn.Write([]byte(data))
			conn.Close()
		}
	}()
}

func runTcpSvr(ln net.Listener, respContent string) {
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			rd := bufio.NewReader(conn)
			_, err = rd.ReadString('\n')
			if err != nil {
				return
			}
			conn.Write([]byte(respContent))
			conn.Close()
		}
	}()
}

func TestMux(t *testing.T) {
	assert := assert.New(t)

	ln, err := net.Listen("tcp", "127.0.0.1:")
	assert.NoError(err)

	mux := NewMux(ln)
	httpLn := mux.ListenHttp(0)
	httpsLn := mux.ListenHttps(0)
	defaultLn := mux.DefaultListener()
	go mux.Serve()
	time.Sleep(100 * time.Millisecond)

	httpSvr := runHttpSvr(httpLn)
	defer httpSvr.Close()
	httpsSvr := runHttpsSvr(httpsLn)
	defer httpsSvr.Close()
	runEchoSvr(defaultLn)
	defer ln.Close()

	// test http service
	resp, err := http.Get(httpSvr.URL)
	assert.NoError(err)
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal("http service", string(data))

	// test https service
	client := httpsSvr.Client()
	resp, err = client.Get(httpsSvr.URL)
	assert.NoError(err)
	data, err = ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal("https service", string(data))

	// test echo service
	conn, err := net.Dial("tcp", ln.Addr().String())
	assert.NoError(err)
	_, err = conn.Write([]byte("test echo\n"))
	assert.NoError(err)
	data = make([]byte, 1024)
	n, err := conn.Read(data)
	assert.NoError(err)
	assert.Equal("test echo\n", string(data[:n]))
}

func TestMuxPriority(t *testing.T) {
	assert := assert.New(t)

	ln, err := net.Listen("tcp", "127.0.0.1:")
	assert.NoError(err)

	mux := NewMux(ln)
	ln1 := mux.Listen(0, 2, func(data []byte) bool {
		if data[0] == '1' {
			return true
		} else {
			return false
		}
	})
	ln2 := mux.Listen(1, 2, func(data []byte) bool {
		if data[0] == '1' {
			return true
		} else {
			return false
		}
	})
	runTcpSvr(ln1, "aaa")
	runTcpSvr(ln2, "bbb")
	go mux.Serve()
	time.Sleep(100 * time.Millisecond)

	// priority 0, '1' -> 'aaa'
	// priority 1, '1' -> 'bbb'
	conn, err := net.Dial("tcp", ln.Addr().String())
	assert.NoError(err)
	_, err = conn.Write([]byte("111\n"))
	assert.NoError(err)
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	assert.NoError(err)
	assert.Equal("aaa", string(data[:n]))

	// No match
	ln1.Close()
	ln2.Close()
	conn, err = net.Dial("tcp", ln.Addr().String())
	assert.NoError(err)
	_, err = conn.Write([]byte("111\n"))
	assert.NoError(err)
	data = make([]byte, 1024)
	_, err = conn.Read(data)
	assert.Error(err)

	// priority 0, '1' -> 'bbb'
	// priority 1, '1' -> 'aaa'
	ln1 = mux.Listen(0, 2, func(data []byte) bool {
		if data[0] == '1' {
			return true
		} else {
			return false
		}
	})
	ln2 = mux.Listen(1, 2, func(data []byte) bool {
		if data[0] == '2' {
			return true
		} else {
			return false
		}
	})
	runTcpSvr(ln2, "aaa")
	runTcpSvr(ln1, "bbb")

	conn, err = net.Dial("tcp", ln.Addr().String())
	assert.NoError(err)
	_, err = conn.Write([]byte("111\n"))
	assert.NoError(err)
	data = make([]byte, 1024)
	n, err = conn.Read(data)
	assert.NoError(err)
	assert.Equal("bbb", string(data[:n]))
}
