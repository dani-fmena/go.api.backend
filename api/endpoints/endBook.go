package endpoints

import (
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"go.api.backend/data/models"
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

	booksRouter := app.Party("/books") 						// This is a go closure, but with a named function
	{
		// --- GROUP MIDDLEWARES ---
		// booksRouter.Use(iris.Compression)

		// --- DEPENDENCIES ---
		hero.Register(services.NewSvcBooks(&bookRepo))

		// TODO Luego que pinche cambia a ver si puedes registrar directo con el router
		// TODO investiga el tratamiento de error
		// TODO Loguea a un archivo, las operaciones ralizadas y pasa eso a los servicios, llamoago asyncor, kataras tiene un ejemplo
		// TODO mke a change log and a readme with the description of the shell project
		// TODO versionado de la APi
		// TODO document all the folders (readmes)

		booksRouter.Get("/", hero.Handler(getBooks))
		booksRouter.Get("/{id:uint64}", hero.Handler(getBookById))
		booksRouter.Delete("/{id:uint64}", hero.Handler(delBookByID))
		booksRouter.Post("/", hero.Handler(createBook))
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
// @Failure 500 {object} dto.ApiError "err.repo_ops"
// @Router /books [get]
func getBooks(ctx iris.Context, svc services.ServiceBook) {
	books, err := svc.GetAll()

	// Preparing the response
	if err != nil {
		services.ResponseErr(iris.StatusInternalServerError, data.ErrRepositoryOps, err.Error(), &ctx)
	} else {
		services.ResponseOKWithData(iris.StatusOK, books, &ctx)
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
// @Success 404 {object} dto.ApiError "err.not_found"
// @Failure 500 {object} dto.ApiError "Internal error"
// @Router /books/{id} [get]
func getBookById(ctx iris.Context, svc services.ServiceBook) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	book, err := svc.GetByID(&bookId)

	// Preparing the response
	if book.CreatedAt != *new(time.Time) && err == nil {											// 200 Founded
		services.ResponseOKWithData(iris.StatusOK, book, &ctx)
	} else if err != nil && err.Error()[4:11] == data.StrDB404 {									// 404 from repo
		services.ResponseErr(iris.StatusNotFound, data.ErrNotFound, "", &ctx)
	} else if err != nil {
		services.ResponseErr(iris.StatusInternalServerError, data.ErrGeneric, err.Error(), &ctx)	// returning some other error may happen
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
// @Success 404 {object} dto.ApiError "err.not_found"
// @Failure 500 {object} dto.ApiError "err.repo_ops"
// @Router /books/{id} [delete]
func delBookByID(ctx iris.Context, svc services.ServiceBook) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	deleted, err := svc.DelByID(&bookId)

	// Preparing the response
	if err == nil && deleted == 0 {
		services.ResponseErr(iris.StatusNotFound, data.ErrNotFound, "", &ctx)  				// 404 from repo
	} else if err == nil && deleted > 0 {
		services.ResponseDelete(&ctx) 																// 204 & empty data
	} else if err != nil {
		services.ResponseErr(iris.StatusInternalServerError, data.ErrRepositoryOps, err.Error(), &ctx)	// returning some other error may happen
	}
}

// createBook create a new book
// @Summary Create a new book
// @Description Create a new book from the passed data
// @Tags Books
// @Accept	json
// @Produce json
// @Param	book	body	models.Book	true "Book Data"
// @Success 201 {object} models.Book "OK"
// @Success 422 {object} dto.ApiError "err.duplicate_key || Invalid data"	// TODO learn to make validation of params and body, make the response 400 (bad request)
// @Failure 500 {object} dto.ApiError "err.repo_ops || Internal error"
// @Router /books [post]
func createBook(ctx iris.Context, svc services.ServiceBook) {
	var b models.Book

	// TIP: use ctx.ReadBody(&b) to bind any type of incoming data instead. E.g it comes in handy when the client request are using form-data
	if e1 := ctx.ReadJSON(&b); e1 != nil {
		services.ResponseErr(iris.StatusUnprocessableEntity, data.ErrGeneric, e1.Error(), &ctx) 		// 422 errors may happen in the marshaling process
	} else if e2 := svc.Create(&b); e2 != nil {

		if e2.Error() == data.ErrDuplicateKey {															// 422 Unprocessable 'cause duplicate key
			services.ResponseErr(iris.StatusUnprocessableEntity, data.ErrDuplicateKey, "", &ctx)
		} else {																						// 500
			services.ResponseErr(iris.StatusInternalServerError, data.ErrRepositoryOps, e2.Error(), &ctx)
		}

	} else {		// All good
		services.ResponseOKWithData(iris.StatusCreated, b, &ctx)
	}
}
// endregion =============================================================================
