package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"html/template"
	"mall-seckill/datamodels"
	"mall-seckill/services"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

var (
	//生成的html保存目录
	htmlOutPath = "./frontend/web/htmlProductShow/"
	//静态文件模版目录
	templatePath = "./frontend/web/views/template/"
)

func (p *ProductController) GetGenerateHtml() {
	productId, err := strconv.ParseInt(p.Ctx.URLParam("productId"), 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//1.获取模版
	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2.获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	//3.获取模版渲染数据
	product, err := p.ProductService.GetProductById(productId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态文件
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	template.Execute(file, &product)
}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
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
