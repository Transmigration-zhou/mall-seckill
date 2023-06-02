package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()
	app.HandleDir("/public", "./frontend/web/public")
	app.HandleDir("/html", "./frontend/web/htmlProductShow")

	app.Run(
		iris.Addr(":8083"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
