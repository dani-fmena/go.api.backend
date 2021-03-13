package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"go.api.backend/data/database"
	"go.api.backend/services"

	"github.com/iris-contrib/swagger/v12"              // swagger middleware for Iris
	"github.com/iris-contrib/swagger/v12/swaggerFiles" // swagger embed files

	_ "github.com/lib/pq"

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

	// TODO document all the folders

	// region ======== GLOBALS ===============================================================

	c := services.NewSvcConfig("D:\\Source\\Go\\src\\go.api.backend\\dev.yaml") // Configuration Services
	// endregion =============================================================================

	// region ======== MIDDLEWARES ===========================================================

	// Built-ins
	app.Use(logger.New())
	app.UseRouter(recover.New()) // Recovery middleware recovers from any panics and writes a 500 if there was one.

	// Customs
	// endregion =============================================================================

	// region ======== DATABASE BOOTSTRAPPING ================================================

	pgdb := database.Bootstrap(c) 							// Starting the database and creating the engine
	database.CreateSchema(pgdb, false)				// Table creation method
	// endregion =============================================================================

	// region ======== ENDPOINT REGISTRATIONS ================================================

	endpoints.BookRegister(app, pgdb)
	// endregion =============================================================================

	// region ======== SWAGGER REGISTRATION ==================================================

	// version https://gitee.com/luckymonkey006/swagger

	// sc == swagger config
	sc := &swagger.Config {
		DeepLinking: true,
		URL: "http://localhost:8080/swagger/apidoc.json", // The url pointing to API definition
	}

	// use swagger middleware to
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(sc, swaggerFiles.Handler))
	// endregion =============================================================================

	app.Run(iris.Addr(":8080"))
	//  app.Listen(":5000", iris.WithOptimizations)
}
