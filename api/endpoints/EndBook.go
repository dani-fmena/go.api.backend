package endpoints

import (
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"go.api.backend/repo/db"

	"go.api.backend/data"
	"go.api.backend/services"
)


// BookRegister register the Books endpoints as routes in the Iris app
//
// - app [*iris.Application] ~ Iris App instance
//
// - path [**pg.DB] ~ Postgres database instance
func BookRegister(app *iris.Application, dbCtx *pg.DB) {

	// --- VARS SETUP ---
	t := db.NewRepoDbBook(dbCtx)

	booksRouter := app.Party("/books") // This is a closure, but with a named function
	{
		// --- GROUP MIDDLEWARES ---
		booksRouter.Use(iris.Compression)

		// --- DEPENDENCIES ---

		hero.Register(services.NewSvcBooks(&t))

		// TODO Luego que pinche cambia a ver si puedes registrar directo con el router
		// TODO Loguea a un archivo y pasa eso a los servicios, llamoago asyncor, kataras tiene un ejemplo
		// TODO pass the a logger object, also implement a log to file

		booksRouter.Get("/", hero.Handler(test))
		// booksRouter.Get("/", hero.Handler(test))
		// booksRouter.Post("/", createBooks)
	}
}

// region ======== LOGIC =================================================================

func test(svc services.Service) {
	svc.GetAll()
}

// listBooks endpoint for listing the  ...
// @Summary Books
// @Description Get all the books
// @Tags books
// @Produce json
// @Success 200 {array} data.Temporal
// @Router /books [get]
func listBooks(ctx iris.Context) {
	// ctx.Application().Logger().Warnf("asdasd")

	books := []data.Temporal{
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

// createBooks endpoint ...
// @Summary Books
// @Description Create a book
// @Tags books
// @Produce json
// @Success 200 {string} string	"ok"
// @Router /books [post]
func createBooks(ctx iris.Context) {
	var b data.Temporal
	err := ctx.ReadJSON(&b)
	// TIP: use ctx.ReadBody(&b) to bind
	// any type of incoming data instead.
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().Title("Temporal creation failure").DetailErr(err))
		// TIP: use ctx.StopWithError(code, err) when only
		// plain text responses are expected on errors.
		return
	}

	println("Received Temporal: " + b.Title)

	ctx.StatusCode(iris.StatusCreated)
}

// endregion =============================================================================
