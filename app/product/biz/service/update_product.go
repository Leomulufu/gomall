package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
)

type UpdateProductService struct {
	ctx context.Context
}

func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	if err := mysql.DB.WithContext(s.ctx).Model(&model.Product{}).Where("id = ?", req.Id).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"picture":     req.Picture,
		"price":       req.Price,
	}).Error; err != nil {
		return nil, err
	}
	return &product.UpdateProductResp{Success: true}, nil
}
