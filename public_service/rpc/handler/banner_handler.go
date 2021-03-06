package handler

import (
	proto "digicon/proto/rpc"
	//"digicon/public_service/log"
	log "github.com/sirupsen/logrus"
	"digicon/public_service/model"

	"golang.org/x/net/context"
)

func (s *RPCServer) GetBannerList(ctx context.Context, req *proto.BannerRequest, rsp *proto.BannerResponse) error {
	err := new(model.Banner).GetBannerList(req, rsp)
	//fmt.Println("1010110bannerlist")
	if err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
