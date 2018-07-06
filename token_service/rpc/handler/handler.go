package handler

import (
	"digicon/common/convert"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"digicon/token_service/model"
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
	"log"
	"time"
)

type RPCServer struct {
}

func (s *RPCServer) AdminCmd(ctx context.Context, req *proto.AdminRequest, rsp *proto.AdminResponse) error {
	log.Print("Received Say.Hello request")
	rsp.Data = "Hello " + req.Cmd
	return nil
}

func (s *RPCServer) EntrustOrder(ctx context.Context, req *proto.EntrustOrderRequest, rsp *proto.CommonErrResponse) error {
	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		rsp.Err = ERR_TOKEN_QUENE_CONF
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	ret, err := q.EntrustReq(req)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(rsp.Err)
	return nil
}

func (s *RPCServer) Symbols(ctx context.Context, req *proto.NullRequest, rsp *proto.SymbolsResponse) error {
	t := new(model.QuenesConfig).GetAllQuenes()

	for _, v := range t {
		if v.TokenId == 1 {

			rsp.Usdt.TokenId = int32(v.TokenId)
			rsp.Usdt.Data = append(rsp.Usdt.Data, &proto.SymbolBaseData{
				Symbol:   v.Name,
				Price:    convert.Int64ToStringBy8Bit(v.Price),
				CnyPrice: convert.Int64ToStringBy8Bit(7 * v.Price),
				Scope:    v.Scope,
			})
		} else if v.TokenId == 2 {
			//rsp.Btc.Data=make([]*proto.SymbolBaseData,0)

			rsp.Btc.TokenId = int32(v.TokenId)
			rsp.Btc.Data = append(rsp.Btc.Data, &proto.SymbolBaseData{
				Symbol:   v.Name,
				Price:    convert.Int64ToStringBy8Bit(v.Price),
				CnyPrice: convert.Int64ToStringBy8Bit(7 * v.Price),
				Scope:    v.Scope,
			})

		} else if v.TokenId == 3 {
			//rsp.Eth.Data=make([]*proto.SymbolBaseData,0)

			rsp.Eth.TokenId = int32(v.TokenId)
			rsp.Eth.Data = append(rsp.Eth.Data, &proto.SymbolBaseData{
				Symbol:   v.Name,
				Price:    convert.Int64ToStringBy8Bit(v.Price),
				CnyPrice: convert.Int64ToStringBy8Bit(7 * v.Price),
				Scope:    v.Scope,
			})
		} else {
			log.Fatalf("errf type %d", v.TokenId)
		}

	}

	return nil
}
func (s *RPCServer) SelfSymbols(ctx context.Context, req *proto.SelfSymbolsRequest, rsp *proto.SelfSymbolsResponse) error {
	//t := new(model.QuenesConfig).GetQuenes(req.Uid)
	return nil
}

func (s *RPCServer) AddTokenNum(ctx context.Context, req *proto.AddTokenNumRequest, rsp *proto.CommonErrResponse) error {
	ret, err := model.AddTokenSess(req)
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
	}
	rsp.Err = ret
	rsp.Message = GetErrorMessage(ret)
	return nil
}

func (s *RPCServer) HistoryKline(ctx context.Context, req *proto.HistoryKlineRequest, rsp *proto.HistoryKlineResponse) error {
	return nil
}

func (s *RPCServer) EntrustQuene(ctx context.Context, req *proto.EntrustQueneRequest, rsp *proto.EntrustQueneResponse) error {
	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		rsp.Err = ERR_TOKEN_QUENE_CONF
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}
	others, err := q.PopFirstEntrust(proto.ENTRUST_OPT_BUY, 2, 5)
	if err == redis.Nil {

	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	} else {
		for _, v := range others {
			g := &proto.EntrustBaseData{
				OnPrice:    convert.Int64ToStringBy8Bit(v.OnPrice),
				SurplusNum: convert.Int64ToStringBy8Bit(v.SurplusNum),
			}
			g.Price = convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(v.OnPrice, v.SurplusNum))
			rsp.Buy = append(rsp.Buy, g)
		}
	}

	others, err = q.PopFirstEntrust(proto.ENTRUST_OPT_SELL, 2, 5)
	if err == redis.Nil {

	} else if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	} else {
		for _, v := range others {
			g := &proto.EntrustBaseData{
				OnPrice:    convert.Int64ToStringBy8Bit(v.OnPrice),
				SurplusNum: convert.Int64ToStringBy8Bit(v.SurplusNum),
			}
			g.Price = convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(v.OnPrice, v.SurplusNum))
			rsp.Sell = append(rsp.Sell, g)
		}
	}

	rsp.Err = ERRCODE_SUCCESS
	return nil
}

/*
func (s *RPCServer) EntrustList(ctx context.Context, req *proto.EntrustQueneRequest, rsp *proto.EntrustQueneResponse) error {
	a := make([][]interface{}, 0)
	godump.Dump(a)
	return nil
}
*/
func (s *RPCServer) EntrustHistory(ctx context.Context, req *proto.EntrustHistoryRequest, rsp *proto.EntrustHistoryResponse) error {

	r := new(model.EntrustDetail).GetHistory(req.Uid, int(req.Limit), int(req.Page))
	var display string
	for _, v := range r {
		if v.Type==int(proto.ENTRUST_TYPE_MARKET_PRICE) {
			display="市价"
		}else {
			display=convert.Int64ToStringBy8Bit(v.Mount)
		}
		rsp.Data = append(rsp.Data, &proto.EntrustHistoryBaseData{
			EntrustId:  v.EntrustId,
			Symbol:     v.Symbol,
			Opt:        proto.ENTRUST_OPT(v.Opt),
			Type:       proto.ENTRUST_TYPE(v.Type),
			AllNum:     convert.Int64ToStringBy8Bit(v.AllNum),
			OnPrice:    convert.Int64ToStringBy8Bit(v.OnPrice),
			TradeNum:   convert.Int64ToStringBy8Bit(v.AllNum - v.SurplusNum),
			Mount:      display,
			CreateTime: time.Unix(v.CreatedTime, 0).Format("2006-01-02 15:04:05"),
			States:     int32(v.States),
		})
	}
	return nil
}

func (s *RPCServer) EntrustList(ctx context.Context, req *proto.EntrustHistoryRequest, rsp *proto.EntrustListResponse) error {
	r := new(model.EntrustDetail).GetList(req.Uid, int(req.Limit), int(req.Page))
	var display string
	for _, v := range r {
		if v.Type==int(proto.ENTRUST_TYPE_MARKET_PRICE) {
			display="市价"
		}else {
			display=convert.Int64ToStringBy8Bit(v.Mount)
		}
		rsp.Data = append(rsp.Data, &proto.EntrustListBaseData{
			EntrustId:  v.EntrustId,
			Symbol:     v.Symbol,
			Opt:        proto.ENTRUST_OPT(v.Opt),
			Type:       proto.ENTRUST_TYPE(v.Type),
			AllNum:     convert.Int64ToStringBy8Bit(v.AllNum),
			OnPrice:    convert.Int64ToStringBy8Bit(v.OnPrice),
			Mount:      display,
			TradeNum:   convert.Int64ToStringBy8Bit(v.AllNum - v.SurplusNum),
			CreateTime: time.Unix(v.CreatedTime, 0).Format("2006-01-02 15:04:05"),
			States:     int32(v.States),
		})
	}
	return nil
}

func (s *RPCServer) Trade(ctx context.Context, req *proto.TradeRequest, rsp *proto.TradeRespone) error {
	q,ok:=model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		rsp.Err = ERR_TOKEN_QUENE_CONF
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	l:=q.GetTradeList(5)
	for _,v:=range l  {
		rsp.Data=append(rsp.Data,&proto.TradeBaseData{
			CreateTime:time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05"),
			Price:convert.Int64ToStringBy8Bit(v.TradePrice),
			Num:convert.Int64ToStringBy8Bit(v.Num),
		})
	}
	return nil
}
