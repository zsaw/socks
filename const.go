package socks

type method byte

// X’00’ NO AUTHENTICATION REQUIRED
// X’01’ GSSAPI
// X’02’ USERNAME/PASSWORD
// X’03’ to X’7F’ IANA ASSIGNED
// X’80’ to X’FE’ RESERVED FOR PRIVATE METHODS
// X’FF’ NO ACCEPTABLE METHODS
const (
	NoAuthenticationRequired method = iota
	Gssapi
	UsernamePassword
	NoAcceptableMethods = 0xff
)

type cmd byte

const (
	Connect cmd = iota + 1
	Bind
	Udp
)

type atyp byte

const (
	IPv4 atyp = iota + 1
	_
	Domain
	IPv6
)

// X’00’ succeeded
// X’01’ general SOCKS server failure
// X’02’ connection not allowed by ruleset
// X’03’ Network unreachable
// X’04’ Host unreachable
// X’05’ Connection refused
// X’06’ TTL expired
// X’07’ Command not supported
// X’08’ Address type not supported
// X’09’ to X’FF’ unassigned
type rep byte

const (
	Succeeded rep = iota
	GeneralSocks5ServerFailure
	ConnectionNotAllowedByRuleset
	NetworkUnreachable
	HostUnreachable
	ConnectionRefused
	TTLExpired
	CommandNotSupported
	AddressTypeNotSupported
)
