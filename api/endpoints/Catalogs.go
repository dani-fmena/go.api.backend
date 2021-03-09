package endpoints

import (
	"github.com/kataras/iris/v12"

	"go.api.backend/models"
)

// BookRegister register the Books endpoints
func CatalogRegister(app *iris.Application) {

	catalogAPI := app.Party("/catalogs") // This is a closure, but with a named function
	{
		// --- GROUP MIDDLEWARES ---
		catalogAPI.Use(iris.Compression)

		catalogAPI.Get("/", listCatalogs)
		catalogAPI.Post("/", createCatalogs)
	}
}

//region ======== LOGIC =================================================================

// list endpoint for listing the  ...
// @Summary Catalog
// @Description Get all the books
// @Tag catalogs
// Group Catalog
// @Produce json
// @Success 200 {object} models.Book "ok"
// @Router /catalogs [get]
func listCatalogs(ctx iris.Context) {
	books := []models.Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}

	ctx.JSON(books)
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

// create endpoint ...
// @Summary Catalog
// @Description Get all the books
// @Tag catalogs
// Group Catalog
// @Produce json
// @Success 200 {string} string	"ok"
// @Router /catalogs [post]
func createCatalogs(ctx iris.Context) {
	var b models.Book
	err := ctx.ReadJSON(&b)
	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Book creation failure").DetailErr(err))
		// TIP: use ctx.StopWithError(code, err) when only
		// plain text responses are expected on errors.
		return
	}

	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}
//endregion =============================================================================
