package socks

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type Handler interface {
	Consult(ConsultWriter, ConsultRequest) (method, error)
	Auth(AuthWriter, AuthRequest) error
	Forward(ForwardWriter, ForwardRequest) (net.Conn, error)
}

func Serve(l net.Listener, handler Handler, consult ConsultRequest, auth AuthRequest, forward ForwardRequest) error {
	srv := &Server{Addr: l.Addr().String(), Handler: handler}
	return srv.Serve(l, consult, auth, forward)
}

func ListenAndServe(addr string, handler Handler) error {
	srv := &Server{Addr: addr, Handler: handler}
	return srv.ListenAndServe()
}

type Server struct {
	Addr          string
	Handler       Handler
	consultWriter ConsultWriter
	authWriter    AuthWriter
	forwardWriter ForwardWriter
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":1080"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln, nil, nil, nil)
}

func (srv *Server) Serve(l net.Listener, consult ConsultRequest, auth AuthRequest, forward ForwardRequest) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		if consult == nil {
			srv.consultWriter = &consultWriter{conn: conn}
		}
		if auth == nil {
			srv.authWriter = &authWriter{conn: conn}
		}
		if forward == nil {
			srv.forwardWriter = &forwardWriter{conn: conn}
		}

		go func(conn net.Conn, h Handler) {
			defer conn.Close()

			creq, err := ReadConsultRequest(bufio.NewReader(conn))
			if err != nil {
				log.Println(err.Error())
				return
			}

			method, err := h.Consult(srv.consultWriter, creq)
			if err != nil {
				log.Println(err.Error())
				return
			}

			switch method {
			case NoAuthenticationRequired:
			case UsernamePassword:
				areq, err := ReadAuthRequest(bufio.NewReader(conn))
				if err != nil {
					log.Println(err.Error())
					return
				}

				err = h.Auth(srv.authWriter, areq)
				if err != nil {
					log.Println(err.Error())
					return
				}
			default:
				log.Println(errors.New("unsupported method"))
				return
			}

			freq, err := ReadForwardRequest(bufio.NewReader(conn))
			if err != nil {
				log.Println(err.Error())
				return
			}

			dconn, err := h.Forward(srv.forwardWriter, freq)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer dconn.Close()

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				_, err := io.Copy(conn, dconn)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}()
			go func() {
				defer wg.Done()
				_, err := io.Copy(dconn, conn)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}()
			wg.Wait()
		}(conn, srv.Handler)
	}
}
