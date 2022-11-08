package socks

import (
	"encoding/binary"
	"errors"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

type ConsultWriter interface {
	Write([]byte) (int, error)
	WriteMethod(method)
}

type AuthWriter interface {
	Write([]byte) (int, error)
	WriteStatus(bool)
}

type ForwardWriter interface {
	Write([]byte) (int, error)
	WriteRepAddr(rep, string)
}

type consultWriter struct{ conn net.Conn }

func (w *consultWriter) Write(b []byte) (int, error) { return w.conn.Write(b) }

func (w *consultWriter) WriteMethod(m method) { w.conn.Write([]byte{5, byte(m)}) }

type authWriter struct{ conn net.Conn }

func (w *authWriter) Write(b []byte) (int, error) { return w.conn.Write(b) }

func (w *authWriter) WriteStatus(s bool) {
	if s {
		w.conn.Write([]byte{5, 0})
	} else {
		w.conn.Write([]byte{5, 1})
	}
}

type forwardWriter struct{ conn net.Conn }

func (w *forwardWriter) Write(b []byte) (int, error) { return w.conn.Write(b) }

func (w *forwardWriter) WriteRepAddr(r rep, addr string) {
	var byts []byte

	addrport, err := netip.ParseAddrPort(addr)
	if err != nil {
		arr := strings.Split(addr, ":")
		if len(arr) != 2 {
			panic(errors.New("domain format error"))
		}

		byts = append(byts, byte(len(arr[0])))
		byts = append(byts, []byte(arr[0])...)

		uintp, err := strconv.ParseUint(arr[1], 10, 16)
		if err != nil {
			panic(err.Error())
		}

		port := make([]byte, 2)
		binary.BigEndian.PutUint16(port, uint16(uintp))
		byts = append(byts, port...)

		w.conn.Write(byts)
		return
	}

	if addrport.Addr().Is4() {
		byts = []byte{5, byte(r), 0, byte(IPv4)}

		as4 := addrport.Addr().As4()
		for i := 0; i < len(as4); i++ {
			byts = append(byts, as4[i])
		}
	} else {
		byts = []byte{5, byte(r), 0, byte(IPv6)}

		as16 := addrport.Addr().As16()
		for i := 0; i < len(as16); i++ {
			byts = append(byts, as16[i])
		}

	}

	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, addrport.Port())
	byts = append(byts, port...)

	w.conn.Write(byts)
}
