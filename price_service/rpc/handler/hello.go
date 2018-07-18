package handler

import (
	"digicon/price_service/model"
	proto "digicon/proto/rpc"
	"golang.org/x/net/context"
	"log"
)

type RPCServer struct{}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Data = "Hello " + req.Cmd
	return nil
}

func (s *RPCServer) CurrentPrice(ctx context.Context, req *proto.CurrentPriceRequest, rsp *proto.CurrentPriceResponse) error {

	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		return nil
	}

	cny := model.GetTokenCnyPrice(q.TokenId)
	e := q.GetEntry()
	rsp.Data = model.Calculate(e.Price, e.Amount, cny, q.Symbol)
	return nil
}

func (s *RPCServer) LastPrice(ctx context.Context, req *proto.LastPriceRequest, rsp *proto.LastPriceResponse) error {
	p,ok:=model.GetPrice(req.Symbol)
	if ok {
		rsp.Data=&proto.PriceCache{
			Id:p.Id,
			Symbol:p.Symbol,
			Amount:p.Amount,
			Vol:p.Vol,
			CreatedTime:p.CreatedTime,
			Count:p.Count,
		}
		return nil
	}
	return nil
}