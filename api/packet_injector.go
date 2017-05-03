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

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/abbot/go-http-auth"
	shttp "github.com/skydive-project/skydive/http"
	"github.com/skydive-project/skydive/packet_injector"
	"github.com/skydive-project/skydive/packet_injector/common"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/graph"
	"github.com/skydive-project/skydive/validator"
)

type PacketInjectorAPI struct {
	PIClient *packet_injector.PacketInjectorClient
	Graph    *graph.Graph
}

type PacketParamsReq struct {
	Src     string
	Dst     string
	SrcIP   string
	DstIP   string
	SrcMAC  string
	DstMAC  string
	Type    string
	Payload string
	Count   int
}

func (pi *PacketInjectorAPI) requestToParams(ppr *PacketParamsReq) (string, *packet_injector_common.PacketParams, error) {
	pi.Graph.RLock()
	defer pi.Graph.RUnlock()

	srcNode := pi.getNode(ppr.Src)
	dstNode := pi.getNode(ppr.Dst)

	if srcNode == nil {
		return "", nil, errors.New("Not able to find source node")
	}

	if ppr.DstIP == "" && dstNode == nil {
		return "", nil, errors.New("Not able to find a dest node and dest IP also empty")
	}

	if ppr.DstMAC == "" && dstNode == nil {
		return "", nil, errors.New("Not able to find a dest node and dest MAC also empty")
	}

	pp := packet_injector_common.PacketParams{
		SrcNodeID: srcNode.ID,
		SrcIP:     ppr.SrcIP,
		SrcMAC:    ppr.SrcMAC,
		DstNodeID: dstNode.ID,
		DstIP:     ppr.DstIP,
		DstMAC:    ppr.DstMAC,
		Type:      ppr.Type,
		Payload:   ppr.Payload,
		Count:     ppr.Count,
	}
	if errs := validator.Validate(pp); errs != nil {
		return "", nil, errors.New("All the parms not set properly.")
	}

	return srcNode.Host(), &pp, nil
}

func (pi *PacketInjectorAPI) injectPacket(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	decoder := json.NewDecoder(r.Body)
	var ppr PacketParamsReq
	if err := decoder.Decode(&ppr); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	host, pp, err := pi.requestToParams(&ppr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := pi.PIClient.InjectPacket(host, pp); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pp); err != nil {
		panic(err)
	}
}

func (pi *PacketInjectorAPI) getNode(gremlinQuery string) *graph.Node {
	res, err := topology.ExecuteGremlinQuery(pi.Graph, gremlinQuery)
	if err != nil {
		return nil
	}

	for _, value := range res.Values() {
		switch value.(type) {
		case *graph.Node:
			return value.(*graph.Node)
		default:
			return nil
		}
	}
	return nil
}

func (pi *PacketInjectorAPI) registerEndpoints(r *shttp.Server) {
	routes := []shttp.Route{
		{
			Name:        "InjectPacket",
			Method:      "POST",
			Path:        "/api/injectpacket",
			HandlerFunc: pi.injectPacket,
		},
	}

	r.RegisterRoutes(routes)
}

func RegisterPacketInjectorAPI(pic *packet_injector.PacketInjectorClient, g *graph.Graph, r *shttp.Server) {
	pia := &PacketInjectorAPI{
		PIClient: pic,
		Graph:    g,
	}

	pia.registerEndpoints(r)
}
