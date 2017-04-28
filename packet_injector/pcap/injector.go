/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package pcap_injector

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"github.com/skydive-project/skydive/packet_injector/common"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/graph"
)

var (
	options = gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
)

func InjectPacket(pp *packet_injector_common.PacketParams, g *graph.Graph) error {
	var srcIP net.IP
	var dstIP net.IP
	var srcMAC net.HardwareAddr
	var dstMAC net.HardwareAddr

	g.RLock()
	defer g.RUnlock()
	srcNode := g.GetNode(pp.SrcNodeID)
	dstNode := g.GetNode(pp.DstNodeID)
	if srcNode == nil {
		return errors.New("Unable to find source node")
	}

	if pp.SrcIP == "" {
		nodeIP, _ := srcNode.GetFieldString("IPv4")
		srcIP = packet_injector_common.GetIP(nodeIP)
	} else {
		srcIP = packet_injector_common.GetIP(pp.SrcIP)
	}

	if srcIP == nil {
		return errors.New("Source Node doesn't have proper IP")
	}

	if pp.DstIP == "" {
		nodeIP, _ := dstNode.GetFieldString("IPv4")
		dstIP = packet_injector_common.GetIP(nodeIP)
	} else {
		dstIP = packet_injector_common.GetIP(pp.DstIP)
	}

	if dstIP == nil {
		return errors.New("Destination doesn't have proper IP")
	}

	if pp.SrcMAC == "" {
		nodeMAC, _ := srcNode.GetFieldString("MAC")
		srcMAC = packet_injector_common.GetMAC(nodeMAC)
	} else {
		srcMAC = packet_injector_common.GetMAC(pp.SrcMAC)
	}

	if srcMAC == nil {
		return errors.New("Source Node doesn't have proper MAC")
	}

	if pp.DstMAC == "" {
		nodeMAC, _ := dstNode.GetFieldString("MAC")
		dstMAC = packet_injector_common.GetMAC(nodeMAC)
	} else {
		dstMAC = packet_injector_common.GetMAC(pp.DstMAC)
	}

	if dstMAC == nil {
		return errors.New("Destination doesn't have proper MAC")
	}

	//create packet
	buffer := gopacket.NewSerializeBuffer()
	ipLayer := &layers.IPv4{Version: 4, SrcIP: srcIP, DstIP: dstIP}
	ethLayer := &layers.Ethernet{EthernetType: layers.EthernetTypeIPv4, SrcMAC: srcMAC, DstMAC: dstMAC}

	switch pp.Type {
	case "icmp":
		ipLayer.Protocol = layers.IPProtocolICMPv4
		gopacket.SerializeLayers(buffer, options,
			ethLayer,
			ipLayer,
			&layers.ICMPv4{
				TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
			},
			gopacket.Payload([]byte(pp.Payload)),
		)
	default:
		return fmt.Errorf("Unsupported traffic type '%s'", pp.Type)
	}

	ifName, _ := srcNode.GetFieldString("Name")

	nscontext, err := topology.NewNetNSContextByNode(g, srcNode)
	defer nscontext.Close()

	if err != nil {
		return err
	}

	handle, err := pcap.OpenLive(ifName, 1024, false, 2000)
	if err != nil {
		return fmt.Errorf("Unable to open the source node: %s", err.Error())
	}
	defer handle.Close()

	packet := buffer.Bytes()
	for i := 0; i < pp.Count; i++ {
		if err := handle.WritePacketData(packet); err != nil {
			return fmt.Errorf("Write error: %s", err.Error())
		}
	}

	return nil
}
