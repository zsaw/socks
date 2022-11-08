package socks

import (
	"bufio"
	"errors"
)

func ReadConsultRequest(b *bufio.Reader) (ConsultRequest, error) {
	req := make(ConsultRequest, 2)
	_, err := b.Read(req)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, req[1])
	_, err = b.Read(buf)
	if err != nil {
		return nil, err
	}

	req = append(req, buf...)
	return req, nil
}

func ReadAuthRequest(b *bufio.Reader) (AuthRequest, error) {
	req := make(AuthRequest, 1)
	_, err := b.Read(req)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 2; i++ {
		len, err := b.ReadByte()
		if err != nil {
			return nil, err
		}
		req = append(req, len)

		buf := make([]byte, len)
		_, err = b.Read(buf)
		if err != nil {
			return nil, err
		}
		req = append(req, buf...)
	}

	return req, nil
}

func ReadForwardRequest(b *bufio.Reader) (ForwardRequest, error) {
	req := make(ForwardRequest, 4)
	_, err := b.Read(req)
	if err != nil {
		return nil, err
	}

	switch req[3] {
	case byte(IPv4):
		var buf [4]byte
		_, err := b.Read(buf[:])
		if err != nil {
			return nil, err
		}
		req = append(req, buf[:]...)
	case byte(Domain):
		buf := make([]byte, 1)
		_, err := b.Read(buf)
		if err != nil {
			return nil, err
		}
		req = append(req, buf...)

		buf = make([]byte, buf[0])
		_, err = b.Read(buf)
		if err != nil {
			return nil, err
		}
		req = append(req, buf...)
	case byte(IPv6):
		var buf [16]byte
		_, err := b.Read(buf[:])
		if err != nil {
			return nil, err
		}
		req = append(req, buf[:]...)
	default:
		return nil, errors.New("unexpected atyp")
	}
	buf := make([]byte, 2)
	_, err = b.Read(buf)
	if err != nil {
		return nil, err
	}
	req = append(req, buf...)
	return req, nil
}
