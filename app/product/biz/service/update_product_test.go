package service

import (
	"context"
	"testing"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProduct_Run(t *testing.T) {
	ctx := context.Background()
	s := NewUpdateProductService(ctx)

	// 初始化请求
	req := &product.UpdateProductReq{
		Id:          1, // 假设商品 ID 为 1
		Name:        "Updated Product",
		Description: "This product has been updated",
		Picture:     "http://example.com/updated.jpg",
		Price:       150.0,
	}

	// 调用服务
	resp, err := s.Run(req)

	// 断言
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, resp, "unexpected nil response")
	assert.True(t, resp.Success, "update should be successful")
}
