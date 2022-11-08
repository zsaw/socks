package socks5

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestSocks(t *testing.T) {
	go func() {
		time.Sleep(time.Second)

		conn, err := net.Dial("tcp", "127.0.0.1:1080")
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer conn.Close()

		conn, err = Client(conn, "93.184.216.34:80")
		if err != nil {
			log.Println(err.Error())
			return
		}

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
		err = req.Write(conn)
		if err != nil {
			log.Println(err.Error())
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			log.Println(err.Error())
			return
		}

		resp.Write(os.Stdout)
	}()

	err := ListenAndServe(":1080")
	if err != nil {
		t.Fatal(err.Error())
		return
	}
}
