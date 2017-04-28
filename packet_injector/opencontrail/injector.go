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

package opencontrail_injector

import (
	"errors"
	"fmt"
	"net"

	"github.com/nlewo/contrail-introspect-cli/collection"
	"github.com/nlewo/contrail-introspect-cli/descriptions"

	"github.com/skydive-project/skydive/config"
	"github.com/skydive-project/skydive/packet_injector/common"
	"github.com/skydive-project/skydive/topology/graph"
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
		nodeIP, _ := srcNode.GetFieldString("Neutron/IPs")
		srcIP = packet_injector_common.GetIP(nodeIP)
	} else {
		srcIP = packet_injector_common.GetIP(pp.SrcIP)
	}

	if srcIP == nil {
		return errors.New("Source Node doesn't have proper IP")
	}

	if pp.DstIP == "" {
		nodeIP, _ := dstNode.GetFieldString("Neutron/IPs")
		dstIP = packet_injector_common.GetIP(nodeIP)
	} else {
		dstIP = packet_injector_common.GetIP(pp.DstIP)
	}

	if dstIP == nil {
		return errors.New("Destination doesn't have proper IP")
	}

	if pp.SrcMAC == "" {
		nodeMAC, _ := srcNode.GetFieldString("ExtID/attached-mac")
		srcMAC = packet_injector_common.GetMAC(nodeMAC)
	} else {
		srcMAC = packet_injector_common.GetMAC(pp.SrcMAC)
	}

	if srcMAC == nil {
		return errors.New("Source Node doesn't have proper MAC")
	}

	if pp.DstMAC == "" {
		nodeMAC, _ := dstNode.GetFieldString("ExtID/attached-mac")
		dstMAC = packet_injector_common.GetMAC(nodeMAC)
	} else {
		dstMAC = packet_injector_common.GetMAC(pp.DstMAC)
	}

	if dstMAC == nil {
		return errors.New("Destination doesn't have proper MAC")
	}

	if pp.SrcPort == 0 {
		return errors.New("Source port missing")
	}

	if pp.DstPort == 0 {
		return errors.New("Destination port missing")
	}

	vrfName := srcNode.Metadata()["Contrail/VRF"].(string)
	host := config.GetConfig().GetString("opencontrail.host")

	switch pp.Type {
	case "tcp":
		params := []string{
			host,
			vrfName,
			srcIP.String(),
			fmt.Sprint(pp.SrcPort),
			dstIP.String(),
			fmt.Sprint(pp.DstPort),
		}
		col, e := collection.LoadCollection(descriptions.AgentPing(), params)
		if e != nil {
			return fmt.Errorf("Failed to contact vrouter-agent")
		}
		col.Short()
		return nil
	default:
		return fmt.Errorf("Unsupported traffic type '%s'", pp.Type)
	}

	return nil
}
