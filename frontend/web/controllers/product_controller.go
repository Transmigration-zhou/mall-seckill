package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"mall-seckill/datamodels"
	"mall-seckill/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	productId, err := strconv.ParseInt(p.Ctx.URLParam("productId"), 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductById(productId)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	userId, err := strconv.ParseInt(p.Ctx.GetCookie("uid"), 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	productId, err := strconv.ParseInt(p.Ctx.URLParam("productId"), 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductById(productId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderId int64
	showMessage := "抢购失败"
	if product.ProductNum > 0 {
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserId:      userId,
			ProductId:   productId,
			OrderStatus: datamodels.OrderSuccess,
		}
		orderId, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功"
		}
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderId":     orderId,
			"showMessage": showMessage,
		},
	}
}
