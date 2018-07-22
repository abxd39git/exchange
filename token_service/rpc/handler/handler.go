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
	/*
		t := new(model.ConfigQuenes).GetAllQuenes()
		rsp.Usdt = new(proto.SymbolsBaseData)
		rsp.Usdt.Data = make([]*proto.SymbolBaseData, 0)
		rsp.Btc = new(proto.SymbolsBaseData)
		rsp.Btc.Data = make([]*proto.SymbolBaseData, 0)
		rsp.Eth = new(proto.SymbolsBaseData)
		rsp.Eth.Data = make([]*proto.SymbolBaseData, 0)
		rsp.Sdc = new(proto.SymbolsBaseData)
		rsp.Sdc.Data = make([]*proto.SymbolBaseData, 0)
		for _, v := range t {
			if v.TokenId == 1 {
				rsp.Usdt.TokenId = int32(v.TokenId)
				cny := model.GetTokenCnyPrice(v.TokenId)

				rsp.Usdt.Data = append(rsp.Usdt.Data, &proto.SymbolBaseData{
					Symbol:       v.Name,
					Price:        convert.Int64ToStringBy8Bit(v.Price),
					CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, v.Price)),
					Scope:        v.Scope,
					TradeTokenId: int32(v.TokenTradeId),
				})
			} else if v.TokenId == 2 {
				rsp.Btc.TokenId = int32(v.TokenId)
				cny := model.GetTokenCnyPrice(v.TokenId)
				rsp.Btc.Data = append(rsp.Btc.Data, &proto.SymbolBaseData{
					Symbol:       v.Name,
					Price:        convert.Int64ToStringBy8Bit(v.Price),
					CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, v.Price)),
					Scope:        v.Scope,
					TradeTokenId: int32(v.TokenTradeId),
				})

			} else if v.TokenId == 3 {
				rsp.Eth.TokenId = int32(v.TokenId)
				cny := model.GetTokenCnyPrice(v.TokenId)
				rsp.Eth.Data = append(rsp.Eth.Data, &proto.SymbolBaseData{
					Symbol:       v.Name,
					Price:        convert.Int64ToStringBy8Bit(v.Price),
					CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, v.Price)),
					Scope:        v.Scope,
					TradeTokenId: int32(v.TokenTradeId),
				})
			} else if v.TokenId == 4 {
				rsp.Sdc.TokenId = int32(v.TokenId)
				cny := model.GetTokenCnyPrice(v.TokenId)
				rsp.Sdc.Data = append(rsp.Sdc.Data, &proto.SymbolBaseData{
					Symbol:       v.Name,
					Price:        convert.Int64ToStringBy8Bit(v.Price),
					CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(cny, v.Price)),
					Scope:        v.Scope,
					TradeTokenId: int32(v.TokenTradeId),
				})
			} else {
				continue
			}

		}
	*/
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
				SurplusNum: convert.Int64ToStringBy8Bit(convert.Int64DivInt64By8Bit(v.SurplusNum, v.OnPrice)),
				CnyPrice:   q.GetCnyPrice(v.OnPrice),
			}
			g.Price = convert.Int64ToStringBy8Bit(v.SurplusNum)
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
				CnyPrice:   q.GetCnyPrice(v.OnPrice),
			}

			g.Price = convert.Int64MulInt64By8BitString(v.OnPrice, v.SurplusNum)
			rsp.Sell = append(rsp.Sell, g)
		}
	}

	rsp.Err = ERRCODE_SUCCESS
	return nil
}

func (s *RPCServer) EntrustHistory(ctx context.Context, req *proto.EntrustHistoryRequest, rsp *proto.EntrustHistoryResponse) error {

	r := new(model.EntrustDetail).GetHistory(req.Uid, int(req.Limit), int(req.Page))
	var display string
	for _, v := range r {
		if v.Type == int(proto.ENTRUST_TYPE_MARKET_PRICE) {
			display = "市价"
		} else {
			display = convert.Int64ToStringBy8Bit(v.Mount)
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
		if v.Type == int(proto.ENTRUST_TYPE_MARKET_PRICE) {
			display = "市价"
		} else {
			display = convert.Int64ToStringBy8Bit(v.Mount)
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
	q, ok := model.GetQueneMgr().GetQueneByUKey(req.Symbol)
	if !ok {
		rsp.Err = ERR_TOKEN_QUENE_CONF
		rsp.Message = GetErrorMessage(rsp.Err)
		return nil
	}

	l := q.GetTradeList(5)
	for _, v := range l {
		rsp.Data = append(rsp.Data, &proto.TradeBaseData{
			CreateTime: time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05"),
			Price:      convert.Int64ToStringBy8Bit(v.TradePrice),
			Num:        convert.Int64ToStringBy8Bit(v.Num),
		})
	}
	return nil
}

func (s *RPCServer) TokenBalance(ctx context.Context, req *proto.TokenBalanceRequest, rsp *proto.TokenBalanceResponse) error {
	d := &model.UserToken{}
	err := d.GetUserToken(req.Uid, int(req.TokenId))
	if err != nil {
		rsp.Err = ERRCODE_UNKNOWN
		rsp.Message = err.Error()
		return nil
	}
	rsp.Balance = &proto.TokenBaseData{
		TokenId: int32(d.TokenId),
		Balance: convert.Int64ToStringBy8Bit(d.Balance),
	}

	return nil
}

func (s *RPCServer) Quotation(ctx context.Context, req *proto.QuotationRequest, rsp *proto.QuotationResponse) error {
	d := &model.ConfigQuenes{}
	g := d.GetQuenesByType(req.TokenId)
	for _, v := range g {
		rsp.Data = append(rsp.Data, &proto.QutationBaseData{
			Symbol: v.Name,
			/*
				Price:  convert.Int64ToStringBy8Bit(v.Price),

				Scope:  v.Scope,
				Low:    convert.Int64ToStringBy8Bit(v.Low),
				High:   convert.Int64ToStringBy8Bit(v.High),
				Amount: convert.Int64ToStringBy8Bit(v.Amount),
			*/
		})
	}
	return nil
}

func (s *RPCServer) GetConfigQuene(ctx context.Context, req *proto.NullRequest, rsp *proto.ConfigQueneResponse) error {
	t := new(model.ConfigQuenes).GetAllQuenes()

	for _, v := range t {
		rsp.Data = append(rsp.Data, &proto.ConfigQueneBaseData{
			TokenId:      int32(v.TokenId),
			TokenTradeId: int32(v.TokenTradeId),
			Name:         v.Name,
		})
	}

	m := model.GetCnyData()
	for _, v := range m {
		rsp.CnyData = append(rsp.CnyData, &proto.CnyPriceBaseData{
			TokenId:  int32(v.TokenId),
			CnyPrice: v.Price,
		})
	}

	return nil
}

func (s *RPCServer) DelEntrust(ctx context.Context, req *proto.DelEntrustRequest, rsp *proto.DelEntrustResponse) error {

	return nil
}
