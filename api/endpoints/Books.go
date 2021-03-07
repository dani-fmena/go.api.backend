package endpoints

import (
	"github.com/kataras/iris/v12"

	"go.api.backend/models"
)

// BookRegister register the Books endpoints
func BookRegister(app *iris.Application) {

	booksAPI := app.Party("/books") // This is a closure, but with a named function
	{
		// Middlewares
		booksAPI.Use(iris.Compression)

		// EndPoints
		booksAPI.Get("/", list)
		booksAPI.Post("/", create)
	}
}

//region ======== LOGIC =================================================================
// list endpoint for listing the ...
func list(ctx iris.Context) {
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
func create(ctx iris.Context) {
	var b models.Book
	err := ctx.ReadJSON(&b)
	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Book creation failure").DetailErr(err))
		// TIP: use ctx.StopWithError(code, err) when only
		// plain text responses are expected on errors.
		return
	}

	println("Received Book: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}
//endregion =============================================================================
