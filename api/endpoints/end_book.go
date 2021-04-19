package endpoints

import (
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"go.api.backend/repo/db"
	"go.api.backend/schema"
	"go.api.backend/schema/dto"
	"go.api.backend/schema/mapper"
	"go.api.backend/service"
	"go.api.backend/service/utils"
	"time"
)

type HBook struct {
	response *utils.SvcResponse
	service *service.SvcBook
}

// NewBookHandler create and register the Books handler and endpoints respectively. The registration create iris app routes
// with their dedicated handlers. If you look at the params you will notice the dependencies used for instantiation of other
// dependencies and for passing it to the handlers. Another way to do this is using iris DI system. This way we don't
// have to create a struct for handler, we can just register the dependencies and the handlers.
//
// - app [*iris.Application] ~ Iris App instance
//
// - path [*pg.DB] ~ Postgres database instance
//
// - r [*utils.SvcResponse] ~ Response service instance
func NewBookHandler(app *iris.Application, dbCtx *pg.DB, r *utils.SvcResponse) HBook {

	// --- VARS SETUP ---
	// TIP As an alternative, we may not use a pointer and leave the cleaning job to the GO garbage collector
	bookRepo := db.NewRepoDbBook(dbCtx)									// Instantiating repo
	bookService := service.NewSvcBooks(&bookRepo)						// Instantiating service

	h := HBook{r, &bookService}

	// --- REGISTERING ENDPOINTS ---
	booksRouter := app.Party("/books") 						// This is a go closure, but with a named function
	{
		// --- GROUP / PARTY MIDDLEWARES ---
		// booksRouter.Use(iris.Compression)

		// --- DEPENDENCIES ---
		// hero.Register(service.NewSvcBooks(&bookRepo))

		booksRouter.Get("/", h.getBooks)
		booksRouter.Get("/{id:uint64}", h.getBookById)
		booksRouter.Post("/", h.createBook)
		booksRouter.Put("/{id:uint64}", h.updateBook)				// PUT vs PATCH https://stackoverflow.com/a/34400076/4196056
		booksRouter.Delete("/{id:uint64}", h.delBookById)
		// booksRouter.Get("/", hero.Handler(getBooks))					// sample with dependency injection
		// booksRouter.Post("/", createBooks)							// when no dependencies injection (but context) is needed
	}

	return h
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
func (h HBook) getBooks(ctx iris.Context) {
	books, err := (*h.service).GetAll()

	// Preparing the response
	if err != nil {
		(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrRepositoryOps, err.Error(), &ctx)
	} else {
		(*h.response).ResWithDataStatus(iris.StatusOK, books, &ctx)
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
// @Failure 404 {object} dto.ApiError "err.not_found"
// @Failure 500 {object} dto.ApiError "Internal error"
// @Router /books/{id} [get]
func(h HBook) getBookById(ctx iris.Context) {

	bookId := ctx.Params().GetUintDefault("id", 1)
	book, err := (*h.service).GetByID(&bookId)

	// Preparing the response
	if book.CreatedAt != *new(time.Time) && err == nil {											// 200 Founded
		(*h.response).ResWithDataStatus(iris.StatusOK, book, &ctx)
	} else if err != nil && err.Error()[4:11] == schema.StrDB404 { // 404 from repo
		(*h.response).ResErr(iris.StatusNotFound, schema.ErrNotFound, schema.ErrDetNotFound, &ctx)
	} else if err != nil {
		(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrGeneric, err.Error(), &ctx) // returning some other error may happen
	}

	// Regarding the "Nilnes" IDE warning, I think the book will not be null. Se the called service method.
}

// delBookById deletes a Book by Id or 404 if doesn't exist
// @Summary Delete a Book
// @Description Deletes a Book by its Id
// @Tags Books
// @Accept  json
// @Produce  json
// @Param 	id	path	int true	"Book ID"	Format(uint32)
// @Success 204 "No Content"
// @Failure 404 {object} dto.ApiError "err.not_found"
// @Failure 500 {object} dto.ApiError "err.repo_ops"
// @Router /books/{id} [delete]
func (h HBook) delBookById(ctx iris.Context) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	deleted, err := (*h.service).DelByID(&bookId)

	// Preparing the response
	if err == nil && deleted == 0 {
		(*h.response).ResErr(iris.StatusNotFound, schema.ErrNotFound, schema.ErrDetNotFound, &ctx) // 404 from repo
	} else if err == nil && deleted > 0 {
		(*h.response).ResDelete(&ctx) // 204 & empty schema
	} else if err != nil {
		(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrRepositoryOps, err.Error(), &ctx) // returning some other error may happen
	}
}

// createBook create a new book
// @Summary Create a new book
// @Description Create a new book from the passed schema
// @Tags Books
// @Accept	json
// @Produce json
// @Param	book	body	dto.BookCreateIn	true	"Book Data"
// @Success 201 {object} models.Book "OK"
// @Failure 422 {object} dto.ApiError "err.duplicate_key || Invalid schema"
// @Failure 500 {object} dto.ApiError "err.repo_ops || Internal error"
// @Router /books [post]
func (h HBook) createBook(ctx iris.Context) {
	var bDto dto.BookCreateIn

	// TIP: use ctx.ReadBody(&bDto) to bind any type of incoming schema instead. E.g it comes in handy when the client request are using form-schema
	if e := ctx.ReadJSON(&bDto); e != nil {
		(*h.response).ResErr(iris.StatusUnprocessableEntity, schema.ErrVal, e.Error(), &ctx) // 422 ReadJSON do the validation here
		return
	}

	// TIP ❗ this is just a sample to show the mapper use.
	// Mapping
	book := mapper.ToBookCreateV(&bDto)

	err := (*h.service).Create(book)
	if err != nil {

		if err.Error() == schema.ErrDuplicateKey { // 422 Unprocessable 'cause duplicate key
			(*h.response).ResErr(iris.StatusUnprocessableEntity, schema.ErrDuplicateKey, schema.ErrDetDuplicateKey, &ctx)
		} else {																							// 500
			(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrRepositoryOps, err.Error(), &ctx)
		}

	} else {		// All good
		(*h.response).ResWithDataStatus(iris.StatusCreated, book, &ctx)
	}
}

// updateBook update the book having the Id passed as path parameter, with the schema passed in the request body
// @Summary Update the indicated book
// @Description Update the book having the specified Id with the schema passed in the request body
// @Tags Books
// @Accept	json
// @Produce json
// @Param 	id		path	int					true	"Book ID"	Format(uint32)
// @Param	book	body	dto.BookUpdateIn	true	"Book Data"
// @Success 200 {object} models.Book "OK"
// @Failure 404 {object} dto.ApiError "err.not_found"
// @Failure 422 {object} dto.ApiError "err.duplicate_key || Invalid schema"
// @Failure 500 {object} dto.ApiError "err.repo_ops || Internal error"
// @Router /books/{id} [put]
func (h HBook) updateBook(ctx iris.Context) {
	var bDto dto.BookUpdateIn

	// Getting the schema
	bDto.Id = ctx.Params().GetUintDefault("id", 0)
	if e := ctx.ReadJSON(&bDto); e != nil {
		(*h.response).ResErr(iris.StatusUnprocessableEntity, schema.ErrVal, e.Error(), &ctx) // 422 errors may happen in the marshaling process
		return
	}

	// Mapping
	book := mapper.ToBookUpdateV(&bDto)

	// Updating
	updated, err := (*h.service).UpdateBook(book)

	if err != nil && err.Error() == schema.ErrNotFound { // 404 Wrong ID
		(*h.response).ResErr(iris.StatusNotFound, schema.ErrNotFound, schema.ErrDetNotFound, &ctx)
	} else if err != nil && err.Error() == schema.ErrDuplicateKey { // Same unique field, name in this case
		(*h.response).ResErr(iris.StatusUnprocessableEntity, schema.ErrDuplicateKey, schema.ErrDetDuplicateKey, &ctx)
	} else if err != nil {																				// Something happen
		(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrRepositoryOps, err.Error(), &ctx)
	} else if updated > 0 {																				// All good, bDto was updated
		(*h.response).ResWithDataStatus(iris.StatusOK, book, &ctx)
	}
}

// endregion =============================================================================

// region ======== LOCAL DEPENDENCIES ====================================================
// endregion =============================================================================
