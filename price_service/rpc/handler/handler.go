package handler

import (
	"digicon/common/convert"
	//"digicon/price_service/exchange"
	"digicon/price_service/model"
	proto "digicon/proto/rpc"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"strings"
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
	e := q.GetEntry()
	h,l:=q.GetDay1MaxPrice()
	rsp.Data = model.Calculate(e.Price, e.Amount, q.CnyPrice, q.Symbol,h,l)
	return nil
}

func (s *RPCServer) SelfSymbols(ctx context.Context, req *proto.SelfSymbolsRequest, rsp *proto.SelfSymbolsResponse) error {
	//t := new(model.QuenesConfig).GetQuenes(req.Uid)
	return nil
}

func (s *RPCServer) LastPrice(ctx context.Context, req *proto.LastPriceRequest, rsp *proto.LastPriceResponse) error {
	p, ok := model.GetPrice(req.Symbol)
	rsp.Ok = ok
	if ok {
		rsp.Data = &proto.PriceCache{
			Id:          p.Id,
			Symbol:      p.Symbol,
			Amount:      p.Amount,
			Vol:         p.Vol,
			CreatedTime: p.CreatedTime,
			Count:       p.Count,
		}

		return nil
	}
	return nil
}

func (s *RPCServer) SymbolTitle(ctx context.Context, req *proto.NullRequest, rsp *proto.SymbolTitleResponse) error {
	g := make([]*proto.TitleBaseData, 0)

	for _, v := range model.ConfigTitles {
		g = append(g, &proto.TitleBaseData{
			Mark:    v.Mark,
			TokenId: int32(v.TokenId),
		})
	}
	rsp.Data = g
	return nil
}

func (s *RPCServer) SymbolsById(ctx context.Context, req *proto.SymbolsByIdRequest, rsp *proto.SymbolsByIdResponse) error {
	g := model.GetConfigQuenesByType(req.TokenId)
	rsp.Data = make([]*proto.SymbolBaseData, 0)

	for _, v := range g {
		q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
		if !ok {
			return nil
		}

		p, ok := model.Get24HourPrice(v.Name)
		if !ok {
			return nil
		}

		price := q.GetEntry().Price
		rsp.Data = append(rsp.Data, &proto.SymbolBaseData{
			Symbol:       v.Name,
			Price:        convert.Int64ToStringBy8Bit(price),
			CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
			Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
			TradeTokenId: int32(v.TokenTradeId),
		})
	}

	return nil
}

func (s *RPCServer) Symbols(ctx context.Context, req *proto.NullRequest, rsp *proto.SymbolsResponse) error {

	rsp.Usdt = new(proto.SymbolsBaseData)
	rsp.Usdt.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Btc = new(proto.SymbolsBaseData)
	rsp.Btc.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Eth = new(proto.SymbolsBaseData)
	rsp.Eth.Data = make([]*proto.SymbolBaseData, 0)
	rsp.Sdc = new(proto.SymbolsBaseData)
	rsp.Sdc.Data = make([]*proto.SymbolBaseData, 0)
	for _, v := range model.ConfigQueneData {
		if v.TokenId == 1 {
			rsp.Usdt.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}

			price := q.GetEntry().Price

			rsp.Usdt.Data = append(rsp.Usdt.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 2 {
			rsp.Btc.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Btc.Data = append(rsp.Btc.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 3 {
			rsp.Eth.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Eth.Data = append(rsp.Eth.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		} else if v.TokenId == 4 {
			rsp.Sdc.TokenId = int32(v.TokenId)

			q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
			if !ok {
				return nil
			}

			p, ok := model.Get24HourPrice(v.Name)
			if !ok {
				return nil
			}
			price := q.GetEntry().Price
			rsp.Sdc.Data = append(rsp.Sdc.Data, &proto.SymbolBaseData{
				Symbol:       v.Name,
				Price:        convert.Int64ToStringBy8Bit(price),
				CnyPrice:     convert.Int64ToStringBy8Bit(convert.Int64MulInt64By8Bit(q.CnyPrice, price)),
				Scope:        convert.Int64DivInt64StringPercent(price-p.Price, p.Price),
				TradeTokenId: int32(v.TokenTradeId),
			})
		}
	}
	return nil
}

func (s *RPCServer) Quotation(ctx context.Context, req *proto.QuotationRequest, rsp *proto.QuotationResponse) error {
	g := model.GetConfigQuenesByType(req.TokenId)

	for _, v := range g {

		q, ok := model.GetQueneMgr().GetQueneByUKey(v.Name)
		if !ok {
			return nil
		}

		price := q.GetEntry().Price
		h,l:=q.GetDay1MaxPrice()
		r := model.Calculate(price, q.GetEntry().Amount, q.CnyPrice, q.Symbol,h,l)

		rsp.Data = append(rsp.Data, &proto.QutationBaseData{
			Symbol: v.Name,
			Price:  r.Price,
			Scope:  r.Scope,
			Low:    r.Low,
			High:   r.High,
			Amount: r.Amount,
		})
	}
	return nil
}

/*
	获取一个币对的价格比
*/
func (s *RPCServer) GetSymbolsRate(ctx context.Context, req *proto.GetSymbolsRateRequest, rsp *proto.GetSymbolsRateResponse) error {
	type BaseData struct {
		Symbol   string `json:"symbol"`    //  btc/usdt
		Price    string `json:"price"`     //  1btc = xxx usdt
		CnyPrice string `json:"cny_price"` //  cny
	}
	data := map[string]*proto.RateBaseData{}
	for _, symbol := range req.Symbols {
		var ok bool
		data[symbol], ok = getSymbolRate(symbol)
		if !ok {
			tmpdata, ok := getOtherSymbolRage(symbol)
			if !ok {
				continue
			}
			data[symbol] = tmpdata
		}

	}
	rsp.Data = data
	return nil
}

func getSymbolRate(symbol string) (data *proto.RateBaseData, ok bool) {
	q, ok := model.GetQueneMgr().GetQueneByUKey(symbol)
	if !ok {
		//fmt.Println(symbol, ok)
		//return getOtherSymbolRage(symbol)
		return
	} else {
		e := q.GetEntry()
		data = &proto.RateBaseData{
			Symbol:   q.Symbol,
			Price:    convert.Int64ToStringBy8Bit(e.Price),
			CnyPrice: convert.Int64ToStringBy8Bit(q.CnyPrice),
		}
	}
	return
}

/*
   如果没有找到币对
*/
func getOtherSymbolRage(symbol string) (data *proto.RateBaseData, ok bool) {
	//fmt.Println("other symbol rage ...", symbol)
	tmpSym := strings.Split(symbol, "/")
	if len(tmpSym) < 2 {
		return
	}
	tmpSymToUSDT := fmt.Sprintf("%s/USDT", tmpSym[0])
	toUSDTQ, ok := model.GetQueneMgr().GetQueneByUKey(tmpSymToUSDT)
	if !ok {
		//fmt.Println(tmpSymToUSDT, " not ok!!!")
		return
	}
	usdtToTmpSym := fmt.Sprintf("%s/USDT", tmpSym[1])
	usdtToQ, ok := model.GetQueneMgr().GetQueneByUKey(usdtToTmpSym)
	if !ok {
		//fmt.Println(usdtToTmpSym, " not ok!!!!")
		return
	}
	BTCPrice := toUSDTQ.GetEntry()
	OtherPrice := usdtToQ.GetEntry()
	price := convert.Int64DivInt64By8Bit(BTCPrice.Price, OtherPrice.Price)
	cnyPrice := convert.Int64MulInt64By8Bit(OtherPrice.Price, toUSDTQ.CnyPrice)
	data = &proto.RateBaseData{
		Symbol:   symbol,
		Price:    convert.Int64ToStringBy8Bit(price),
		CnyPrice: convert.Int64ToStringBy8Bit(cnyPrice),
	}
	return
}

func (s *RPCServer) Volume(ctx context.Context, req *proto.VolumeRequest, rsp *proto.VolumeResponse) error {
	data := model.GetVolumeTotal()
	if data != nil {
		rsp.DayVolume = data.DayVolume / 100000000
		rsp.WeekVolume = data.WeekVolume / 100000000
		rsp.MonthVolume = data.MonthVolume / 100000000
	}
	return nil
}
