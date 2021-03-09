package endpoints

import (
	"github.com/kataras/iris/v12"

	"go.api.backend/models"
)

// BookRegister register the Books endpoints
func BookRegister(app *iris.Application) {

	booksAPI := app.Party("/books") // This is a closure, but with a named function
	{
		// --- GROUP MIDDLEWARES ---
		booksAPI.Use(iris.Compression)

		booksAPI.Get("/", listBooks)
		booksAPI.Post("/", createBooks)
	}
}

//region ======== LOGIC =================================================================

// list endpoint for listing the  ...
// @Summary Books
// @Description Get all the books
// @Tag books
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func listBooks(ctx iris.Context) {
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
// @Summary Books
// @Description Get all the books
// @Tag books
// @Produce json
// @Success 200 {string} string	"ok"
// @Router /books [post]
func createBooks(ctx iris.Context) {
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