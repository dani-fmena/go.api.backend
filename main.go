package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"

	"github.com/iris-contrib/swagger/v12"              // swagger middleware for Iris
	"github.com/iris-contrib/swagger/v12/swaggerFiles" // swagger embed files

	"go.api.backend/api/endpoints"
	_ "go.api.backend/docs"
)

// @title Shell Project
// @version 0.0
// @description Api description shell project

// @contact.name Name Test
// @contact.url http://contact.sample/text
// @contact.email sample@mail.io

// @host localhost:8080
// @BasePath /
func main() {
	app := iris.New()

	//region ======== GLOBAL MIDDLEWARES ====================================================

	// Built-in
	app.UseRouter(recover.New()) //Recovery middleware recovers from any panics and writes a 500 if there was one.

	// Customs

	//endregion =============================================================================

	//region ======== ENDPOINT REGISTRATIONS ================================================

	endpoints.BookRegister(app)
	endpoints.CatalogRegister(app)

	//endregion =============================================================================

	//region ======== SWAGGER REGISTRATIONS =================================================

	// version https://gitee.com/luckymonkey006/swagger
	config := &swagger.Config {
		URL: "http://localhost:8080/swagger/doc.json", //The url pointing to API definition
	}

	// use swagger middleware to
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))
	//endregion =============================================================================

	app.Run(iris.Addr(":8080"))
}
