package packet_injector_common

import (
	"net"
	"strings"

	"github.com/skydive-project/skydive/topology/graph"
)

type PacketParams struct {
	SrcNodeID graph.Identifier `valid:"nonzero"`
	SrcIP     string           `valid:"nonzero"`
	SrcMAC    string           `valid:"nonzero"`
	SrcPort   int32            `valid:"nonzero"`
	DstNodeID graph.Identifier `valid:"nonzero"`
	DstIP     string           `valid:"nonzero"`
	DstMAC    string           `valid:"nonzero"`
	DstPort   int32            `valid:"nonzero"`
	Type      string           `valid:"nonzero"`
	Payload   string
	Count     int `valid:"nonzero"`
}

func GetIP(cidr string) net.IP {
	if len(cidr) <= 0 {
		return nil
	}
	ips := strings.Split(cidr, ",")
	//TODO(masco): currently taking first IP, need to implement to select a proper IP
	ip, _, err := net.ParseCIDR(ips[0])
	if err != nil {
		return nil
	}
	return ip
}

func GetMAC(mac string) net.HardwareAddr {
	valid, err := net.ParseMAC(mac)
	if err != nil {
		return nil
	}
	return valid
}
