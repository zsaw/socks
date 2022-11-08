package socks

import (
	"errors"
	"net"
	"testing"
)

type handler struct{ Handler }

func (h *handler) Consult(w ConsultWriter, r ConsultRequest) (method, error) {
	methods := r.Methods()
	for i := 0; i < len(methods); i++ {
		switch methods[i] {
		case NoAuthenticationRequired:
			w.WriteMethod(NoAuthenticationRequired)
			return NoAuthenticationRequired, nil
		}
	}

	w.WriteMethod(NoAcceptableMethods)
	return NoAcceptableMethods, errors.New("no acceptable methods")
}

func (h *handler) Auth(AuthWriter, AuthRequest) error { return nil }

func (h *handler) Forward(w ForwardWriter, r ForwardRequest) (net.Conn, error) {
	conn, err := net.Dial("tcp", r.Addr())
	if err != nil {
		w.WriteRepAddr(NetworkUnreachable, r.Addr())
		return nil, err
	}
	w.WriteRepAddr(Succeeded, r.Addr())
	return conn, nil
}

func TestSocks(t *testing.T) {
	err := ListenAndServe(":1080", &handler{})
	if err != nil {
		t.Fatal(err.Error())
	}
}
