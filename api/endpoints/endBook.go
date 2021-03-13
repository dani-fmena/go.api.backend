package endpoints

import (
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"time"

	"go.api.backend/data"
	"go.api.backend/repo/db"
	"go.api.backend/services"
)

// BookRegister register the Books endpoints as routes in the Iris app
//
// - app [*iris.Application] ~ Iris App instance
//
// - path [**pg.DB] ~ Postgres database instance
func BookRegister(app *iris.Application, dbCtx *pg.DB) {

	// --- VARS SETUP ---
	bookRepo := db.NewRepoDbBook(dbCtx)

	booksRouter := app.Party("/books") 			// This is a go closure, but with a named function
	{
		// --- GROUP MIDDLEWARES ---
		// booksRouter.Use(iris.Compression)

		// --- DEPENDENCIES ---
		hero.Register(services.NewSvcBooks(&bookRepo))

		// TODO learn to make validation of params and body, make the response 400 (bad request)
		// TODO Luego que pinche cambia a ver si puedes registrar directo con el router
		// TODO investiga el tratamiento de error
		// TODO validación de entidades
		// TODO Loguea a un archivo, las operaciones ralizadas y pasa eso a los servicios, llamoago asyncor, kataras tiene un ejemplo
		// TODO pass the a logger object, also implement a log to file
		// TODO mke a change log and a readme with the description of the shell project

		booksRouter.Get("/", hero.Handler(getBooks))
		booksRouter.Get("/{id:uint64}", hero.Handler(getBookById))
		booksRouter.Delete("/{id:uint64}", hero.Handler(delBookByID))
		// booksRouter.Post("/", createBooks)							// when no dependencies (but context) is needed
	}
}
// region ======== ENDPOINT HANDLERS =====================================================

// getBooks list all the books in the repository
// @Summary Get Books
// @Description Get the books in the repository
// @Tags Books
// @Produce json
// @Success 200 {array} models.Book "List of Books"
// @Failure 500 {object} dto.ApiError "Internal error, struct same as Iris.Problem"
// @Router /books [get]
func getBooks(ctx iris.Context, svc services.ServiceBook) {
	books, err := svc.GetAll()

	// Preparing the response
	if err != nil {
		services.ResponseErr500(data.ErrRepositoryOps, err.Error(), &ctx)
	} else {
		services.ResponseOKWithData(books, &ctx)
	}
}

// getBookById Get a book by Id or 404 if doesn't exist
// @Summary Get book by Id
// @Description Get a book through its Id
// @Tags Books
// @Accept  json
// @Produce json
// @Param	id	path	int	true	"Requested Book Id"	Format(uint32)
// @Success 200 {object} models.Book "OK"
// @Success 404 {object} dto.ApiError "Book not found"
// @Failure 500 {object} dto.ApiError "Internal error, same struct as Iris.Problem"
// @Router /books/{id} [get]
func getBookById(ctx iris.Context, svc services.ServiceBook) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	book, err := svc.GetByID(&bookId)

	// Preparing the response
	if book.CreatedAt != *new(time.Time) && err == nil {							// 200 Founded
		services.ResponseOKWithData(book, &ctx)
	} else if err != nil && err.Error()[4:11] == data.StrDB404 {					// 404 from repo
		services.ResponseErr404(data.ErrNotFound, "", &ctx)
	} else if err != nil {
		services.ResponseErr500(data.ErrGeneric, err.Error(), &ctx)					// returning some other error may happen
	}

	// Regarding the "Nilnes" IDE warning, I think the book will not be null. Se the called service method.
}

// delBookByID deletes a Book by Id or 404 if doesn't exist
// @Summary Delete a Book
// @Description Deletes a Book by its Id
// @Tags Books
// @Accept  json
// @Produce  json
// @Param 	id	path	int true	"Account ID"	Format(uint32)
// @Success 204 "No Content"
// @Success 404 {object} dto.ApiError "Book not found"
// @Failure 500 {object} dto.ApiError "Internal error, same struct as Iris.Problem"
// @Router /books/{id} [delete]
func delBookByID(ctx iris.Context, svc services.ServiceBook) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	deleted, err := svc.DelByID(&bookId)

	// Preparing the response
	if err == nil && deleted == 0 {
		services.ResponseErr404(data.ErrNotFound, "", &ctx)					// 404 from repo
	} else if err == nil && deleted > 0 {
		services.ResponseDeleteOK(&ctx)												// 204 & empty data
	} else if err != nil {
		services.ResponseErr500(data.ErrGeneric, err.Error(), &ctx)					// returning some other error may happen
	}
}

// crate endpoint ...
// @Summary Books
// @Description Create a book
// @Tags books
// @Produce json
// @Success 200 {string} string	"ok"
// @Router /books [post]
func crate(ctx iris.Context) {
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
