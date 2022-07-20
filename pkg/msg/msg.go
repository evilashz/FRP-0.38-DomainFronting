// Copyright 2016 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msg

import "net"

const (
	TypeLogin                 = 'o'
	TypeLoginResp             = '1'
	TypeNewProxy              = 'p'
	TypeNewProxyResp          = '2'
	TypeCloseProxy            = 'c'
	TypeNewWorkConn           = 'w'
	TypeReqWorkConn           = 'r'
	TypeStartWorkConn         = 's'
	TypeNewVisitorConn        = 'v'
	TypeNewVisitorConnResp    = '3'
	TypePing                  = 'h'
	TypePong                  = '4'
	TypeUDPPacket             = 'u'
	TypeNatHoleVisitor        = 'i'
	TypeNatHoleClient         = 'n'
	TypeNatHoleResp           = 'm'
	TypeNatHoleClientDetectOK = 'd'
	TypeNatHoleSid            = '5'
)

var (
	msgTypeMap = map[byte]interface{}{
		TypeLogin:                 Login{},
		TypeLoginResp:             LoginResp{},
		TypeNewProxy:              NewProxy{},
		TypeNewProxyResp:          NewProxyResp{},
		TypeCloseProxy:            CloseProxy{},
		TypeNewWorkConn:           NewWorkConn{},
		TypeReqWorkConn:           ReqWorkConn{},
		TypeStartWorkConn:         StartWorkConn{},
		TypeNewVisitorConn:        NewVisitorConn{},
		TypeNewVisitorConnResp:    NewVisitorConnResp{},
		TypePing:                  Ping{},
		TypePong:                  Pong{},
		TypeUDPPacket:             UDPPacket{},
		TypeNatHoleVisitor:        NatHoleVisitor{},
		TypeNatHoleClient:         NatHoleClient{},
		TypeNatHoleResp:           NatHoleResp{},
		TypeNatHoleClientDetectOK: NatHoleClientDetectOK{},
		TypeNatHoleSid:            NatHoleSid{},
	}
)

// When frpc start, client send this message to login to server.
type Login struct {
	Version      string            `json:"NA5f6Elv"`
	Hostname     string            `json:"I9dZ4SKU"`
	Os           string            `json:"FbPeDZbt"`
	Arch         string            `json:"WGrELROn"`
	User         string            `json:"sbJMzEPx"`
	PrivilegeKey string            `json:"4nglzYMT"`
	Timestamp    int64             `json:"uF5dvtAV"`
	RunID        string            `json:"ywgHozPK"`
	Metas        map[string]string `json:"EgSd3K5D"`

	// Some global configures.
	PoolCount int `json:"BdCC9mzc"`
}

type LoginResp struct {
	Version       string `json:"1AU2EZzb"`
	RunID         string `json:"E6ToTLdJ"`
	ServerUDPPort int    `json:"KwgXadlW"`
	Error         string `json:"DmW7EnHs"`
}

// When frpc login success, send this message to frps for running a new proxy.
type NewProxy struct {
	ProxyName      string            `json:"kjPD8ymd"`
	ProxyType      string            `json:"oOK3woP9"`
	UseEncryption  bool              `json:"lio3pBG6"`
	UseCompression bool              `json:"dU1SruHr"`
	Group          string            `json:"Jum4BsFF"`
	GroupKey       string            `json:"mAs2FSQa"`
	Metas          map[string]string `json:"rBUyX5GO"`

	// tcp and udp only
	RemotePort int `json:"eAx7SmjL"`

	// http and https only
	CustomDomains     []string          `json:"jVKtJp3i"`
	SubDomain         string            `json:"zRK1kS59"`
	Locations         []string          `json:"ONc37dJz"`
	HTTPUser          string            `json:"08luWoTX"`
	HTTPPwd           string            `json:"4tvPnutE"`
	HostHeaderRewrite string            `json:"ZXer20Ho"`
	Headers           map[string]string `json:"mvIc7DzD"`

	// stcp
	Sk string `json:"yi1w4HCX"`

	// tcpmux
	Multiplexer string `json:"NWnkQx2Z"`
}

type NewProxyResp struct {
	ProxyName  string `json:"H8l1hsgF"`
	RemoteAddr string `json:"EB1l7q1s"`
	Error      string `json:"o6Em0DKt"`
}

type CloseProxy struct {
	ProxyName string `json:"c3Ndlv18"`
}

type NewWorkConn struct {
	RunID        string `json:"2QMWUhdS"`
	PrivilegeKey string `json:"RBmJbSOI"`
	Timestamp    int64  `json:"VoZHt5JO"`
}

type ReqWorkConn struct {
}

type StartWorkConn struct {
	ProxyName string `json:"qr35aqvs"`
	SrcAddr   string `json:"gbfsyOG2"`
	DstAddr   string `json:"uYP9tTJL"`
	SrcPort   uint16 `json:"rGVoEDZk"`
	DstPort   uint16 `json:"KHmOsYyP"`
	Error     string `json:"GfwqiRkd"`
}

type NewVisitorConn struct {
	ProxyName      string `json:"meBXugzZ"`
	SignKey        string `json:"DJzEYDlA"`
	Timestamp      int64  `json:"tGfl6qZM"`
	UseEncryption  bool   `json:"QluggTRS"`
	UseCompression bool   `json:"6CfFSPOG"`
}

type NewVisitorConnResp struct {
	ProxyName string `json:"2t1l2mbi"`
	Error     string `json:"iYmwDrAB"`
}

type Ping struct {
	PrivilegeKey string `json:"Vvuwhtz0"`
	Timestamp    int64  `json:"pIt2fdc1"`
}

type Pong struct {
	Error string `json:"7qurzJeB"`
}

type UDPPacket struct {
	Content    string       `json:"c"`
	LocalAddr  *net.UDPAddr `json:"l"`
	RemoteAddr *net.UDPAddr `json:"r"`
}

type NatHoleVisitor struct {
	ProxyName string `json:"uG57MyRa"`
	SignKey   string `json:"y7CowXmC"`
	Timestamp int64  `json:"QqmyBkl6"`
}

type NatHoleClient struct {
	ProxyName string `json:"ttDDrVYD"`
	Sid       string `json:"SG7CLdqR"`
}

type NatHoleResp struct {
	Sid         string `json:"buGmOcO7"`
	VisitorAddr string `json:"zFhGu2Ry"`
	ClientAddr  string `json:"5HWh3BXL"`
	Error       string `json:"msbjcwD1"`
}

type NatHoleClientDetectOK struct {
}

type NatHoleSid struct {
	Sid string `json:"8aAQW4Bj"`
}
