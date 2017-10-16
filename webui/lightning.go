package webui

import (
	"fmt"
	"strconv"
	"strings"

	lr "github.com/cdecker/kugelblitz/lightningrpc"
)

// This file takes care of handing incoming JSON-RPC requests over to
// the matching JSON-RPC call in lightning. These are just the shims
// that are exposed over the jsonrpc2

type Lightning struct {
	lrpc *lr.LightningRpc
}

func NewLightning(lrpc *lr.LightningRpc) Lightning {
	return Lightning{
		lrpc: lrpc,
	}
}

func (l *Lightning) IsAlive(_ *lr.Empty) bool {
	_, err := l.lrpc.GetInfo()
	return err == nil
}

func (l *Lightning) Close(req *lr.PeerReference, res *lr.Empty) error {
	return l.lrpc.Close(req.PeerId)
}

func (l *Lightning) GetInfo(_ *lr.Empty, res *lr.GetInfoResponse) error {
	info, err := l.lrpc.GetInfo()
	*res = info
	return err
}

func (l *Lightning) GetPeers(_ *lr.Empty, res *lr.GetPeersResponse) error {
	peers, err := l.lrpc.GetPeers()
	*res = peers
	return err
}

func (l *Lightning) GetRoute(req *lr.GetRouteRequest, res *lr.Route) error {
	route, err := l.lrpc.GetRoute(req.Destination, req.Amount, req.RiskFactor)
	*res = route
	return err
}

func (l *Lightning) NewAddress(_ *lr.Empty, res *lr.NewAddressResponse) error {
	addr, err := l.lrpc.NewAddress()
	*res = addr
	return err
}

func (l *Lightning) SendPayment(req *lr.SendPaymentRequest, res *lr.SendPaymentResponse) error {
	response, err := l.lrpc.SendPayment(req.Route, req.PaymentHash)
	*res = response
	return err
}

func (l *Lightning) AddFunds(rawtx string) error {
	return l.lrpc.AddFunds(rawtx)
}
func (l *Lightning) FundChannel(nodeid string, capacity uint64) error {
	return l.lrpc.FundChannel(nodeid, capacity)
}

func (l *Lightning) Connect(req *lr.ConnectRequest, _ *lr.Empty) error {
	return l.lrpc.Connect(req.Host, req.Port, req.NodeId)
}

type PaymentRequestInfoRequest struct {
	Destination string `json:"destination"`
}

type PaymentRequestInfoResponse struct {
	Hops        []lr.RouteHop `json:"route"`
	PaymentHash string        `json:"paymenthash"`
	Amount      uint64        `json:"amount"`
}

func (l *Lightning) GetPaymentRequestInfo(req *PaymentRequestInfoRequest, res *PaymentRequestInfoResponse) error {
	var route lr.Route
	fmt.Println(req.Destination)
	parts := strings.Split(req.Destination, ":")
	amount, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return err
	}
	routeReq := &lr.GetRouteRequest{
		Destination: parts[0],
		Amount:      amount,
		RiskFactor:  1,
	}
	fmt.Println(parts[0])
	res.Amount = amount
	res.PaymentHash = parts[1]
	err = l.GetRoute(routeReq, &route)
	res.Hops = route.Hops
	return nil
}
