package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"go.api.backend/schema/database"
	"go.api.backend/service/utils"

	"github.com/iris-contrib/swagger/v12"              // swagger middleware for Iris
	"github.com/iris-contrib/swagger/v12/swaggerFiles" // swagger embed files

	"github.com/go-playground/validator/v10"

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
	// region ======== GLOBALS ===============================================================
	v := validator.New()	// Validator instance
	// TIP validation reference https://github.com/kataras/iris/wiki/Model-validation | https://github.com/go-playground/validator | https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835

	app := iris.New()		// App instance
	app.Validator = v		// Register validation on the iris app

	// Services
	svcC := utils.NewSvcConfig("D:\\Source\\Go\\src\\go.api.backend\\conf.dev.yaml") 			// Creating Configuration Service
	svcR := utils.NewSvcResponse(svcC)                                               				// Creating Response Service
	// endregion =============================================================================

	// region ======== MIDDLEWARES ===========================================================

	// Built-ins
	app.Use(logger.New())
	app.UseRouter(recover.New()) // Recovery middleware recovers from any panics and writes a 500 if there was one.

	// Customs
	// endregion =============================================================================

	// region ======== DATABASE BOOTSTRAPPING ================================================

	pgdb := database.Bootstrap(svcC) // Starting the database and creating the engine
	// database.CreateSchema(pgdb, false) 				// Table creation method
	// database.MkMigrations(svcC)						// Making migrations
	// endregion =============================================================================

	// region ======== ENDPOINT REGISTRATIONS ================================================

	endpoints.NewBookHandler(app, pgdb, svcR)
	endpoints.NewAuthHandler(app, svcR, svcC)
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
	//  app.Listen(":5000", iris.WithOptimizations) see https://github.com/kataras/iris/issues/1739, check if it related to the context.go 2307 line
}
