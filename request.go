package socks

import (
	"encoding/binary"
	"fmt"
	"net/netip"
)

type ConsultRequest []byte

func (r ConsultRequest) Methods() []method {
	methods := make([]method, r[1])
	for i := 0; i < len(methods); i++ {
		methods[i] = method(r[2+i])
	}
	return methods
}

type AuthRequest []byte

func (r AuthRequest) UName() []byte { return r[2 : 2+r[1]] }

func (r AuthRequest) Passwd() []byte { return r[3+r[1]:] }

type ForwardRequest []byte

func (r ForwardRequest) Cmd() cmd { return cmd(r[1]) }

func (r ForwardRequest) Addr() string {
	switch r[3] {
	case byte(IPv4):
		addr := netip.AddrFrom4([4]byte{r[4], r[5], r[6], r[7]})
		addrport := netip.AddrPortFrom(addr, binary.BigEndian.Uint16(r[8:]))
		return addrport.String()
	case byte(IPv6):
		addr := netip.AddrFrom16([16]byte{r[4], r[5], r[6], r[7], r[8], r[9], r[10], r[11], r[12], r[13], r[14], r[15], r[16], r[17], r[18], r[19]})
		addrport := netip.AddrPortFrom(addr, binary.BigEndian.Uint16(r[20:]))
		return addrport.String()
	case byte(Domain):
		return fmt.Sprintf("%s:%d", r[5:5+r[4]], binary.BigEndian.Uint16(r[5+r[4]:]))
	default:
		return ""
	}
}
