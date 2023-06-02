package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"mall-seckill/common"
	"mall-seckill/frontend/middleware"
	"mall-seckill/frontend/web/controllers"
	"mall-seckill/repositories"
	"mall-seckill/services"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	template := iris.HTML("./frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	app.HandleDir("/public", "./frontend/web/public")
	app.HandleDir("/html", "./frontend/web/htmlProductShow")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userRepository := repositories.NewUserManager("user", db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	productParty := app.Party("/product")
	productParty.Use(middleware.AuthConProduct)
	product := mvc.New(productParty)
	product.Register(ctx, productService, orderService)
	product.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr(":8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
