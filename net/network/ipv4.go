package network

import (
	"encoding/binary"
	"errors"
	"net"
	"sync/atomic"

	"github.com/solomonwzs/goxutil/net/xnetutil"
)

var _IPHdrID uint32 = 0

func NewIPHdrID() uint16 {
	id := atomic.AddUint32(&_IPHdrID, 1)
	return uint16(id & 0xffff)
}

type IPv4Header struct {
	Version    uint8
	IHL        uint8
	TOS        uint8
	Length     uint16
	Id         uint16
	Flags      uint16
	FragOffset uint16
	TTL        uint8
	Protocol   uint8
	Checksum   uint16
	SrcAddr    net.IP
	DstAddr    net.IP
	Options    []byte
}

func (h *IPv4Header) Marshal() (b []byte, err error) {
	if len(h.Options) == 0 {
		b = make([]byte, SIZEOF_IPV4_HEADER)
	} else {
		b = make([]byte, SIZEOF_IPV4_HEADER+4)
		copy(b[SIZEOF_IPV4_HEADER:], h.Options)
	}
	h.IHL = uint8(len(b))

	b[0] = h.Version<<4 | (h.IHL >> 2 & 0x0f)
	b[1] = h.TOS
	binary.BigEndian.PutUint16(b[2:], h.Length)
	binary.BigEndian.PutUint16(b[4:], h.Id)
	binary.BigEndian.PutUint16(b[6:],
		h.Flags<<13|(h.FragOffset&0x1fff))
	b[8] = h.TTL
	b[9] = h.Protocol

	if ip := h.SrcAddr.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}
	if ip := h.DstAddr.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	} else {
		return nil, errors.New("[ipv4] missing address")
	}

	h.Checksum = xnetutil.Checksum(b)
	binary.BigEndian.PutUint16(b[10:], h.Checksum)

	return
}

func IPv4HeaderUnmarshal(b []byte) (h *IPv4Header, err error) {
	h = &IPv4Header{}

	h.Version = b[0] >> 4
	h.IHL = (b[0] & 0x0f) << 2
	if len(b) < int(h.IHL) {
		return nil, errors.New("[ipv4] IHL error")
	}

	h.TOS = b[1]
	h.Length = binary.BigEndian.Uint16(b[2:])
	h.Id = binary.BigEndian.Uint16(b[4:])

	f := binary.BigEndian.Uint16(b[6:])
	h.Flags = f >> 13
	h.FragOffset = f & 0x1fff

	h.TTL = b[8]
	h.Protocol = b[9]
	h.Checksum = binary.BigEndian.Uint16(b[10:])

	h.SrcAddr = net.IP(b[12:16])
	h.DstAddr = net.IP(b[16:20])

	if h.IHL > SIZEOF_IPV4_HEADER {
		h.Options = make([]byte, h.IHL-SIZEOF_IPV4_HEADER)
		copy(h.Options, b[SIZEOF_IPV4_HEADER:])
	}

	return
}
