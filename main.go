package main

import (
	"go.api.backend/endpoints"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := iris.New()

	//region ======== GLOBAL MIDDLEWARES ====================================================

	app.UseRouter(recover.New()) //Recovery middleware recovers from any panics and writes a 500 if there was one.

	//endregion =============================================================================

	//region ======== ENDPOINT REGISTRATIONS ================================================

	endpoints.BookRegister(app)

	//endregion =============================================================================

	app.Listen(":8080")
}
