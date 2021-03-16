package endpoints

import (
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"go.api.backend/data/dto"
	"go.api.backend/data/mapper"
	"go.api.backend/service/utils"
	"time"

	"go.api.backend/data"
	"go.api.backend/repo/db"
	"go.api.backend/service"
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
		booksRouter.Delete("/{id:uint64}", h.delBookByID)
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
		(*h.response).RespErr(iris.StatusInternalServerError, data.ErrRepositoryOps, err.Error(), &ctx)
	} else {
		(*h.response).RespOKWithData(iris.StatusOK, books, &ctx)
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
func(h HBook) getBookById(ctx iris.Context) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	book, err := (*h.service).GetByID(&bookId)

	// Preparing the response
	if book.CreatedAt != *new(time.Time) && err == nil {											// 200 Founded
		(*h.response).RespOKWithData(iris.StatusOK, book, &ctx)
	} else if err != nil && err.Error()[4:11] == data.StrDB404 {									// 404 from repo
		(*h.response).RespErr(iris.StatusNotFound, data.ErrNotFound, data.ErrDetNotFound, &ctx)
	} else if err != nil {
		(*h.response).RespErr(iris.StatusInternalServerError, data.ErrGeneric, err.Error(), &ctx) // returning some other error may happen
	}

	// Regarding the "Nilnes" IDE warning, I think the book will not be null. Se the called service method.
}

// delBookByID deletes a Book by Id or 404 if doesn't exist
// @Summary Delete a Book
// @Description Deletes a Book by its Id
// @Tags Books
// @Accept  json
// @Produce  json
// @Param 	id	path	int true	"Book ID"	Format(uint32)
// @Success 204 "No Content"
// @Success 404 {object} dto.ApiError "err.not_found"
// @Failure 500 {object} dto.ApiError "err.repo_ops"
// @Router /books/{id} [delete]
func (h HBook) delBookByID(ctx iris.Context) {
	bookId := ctx.Params().GetUintDefault("id", 0)
	deleted, err := (*h.service).DelByID(&bookId)

	// Preparing the response
	if err == nil && deleted == 0 {
		(*h.response).RespErr(iris.StatusNotFound, data.ErrNotFound, data.ErrDetNotFound, &ctx) // 404 from repo
	} else if err == nil && deleted > 0 {
		(*h.response).RespDelete(&ctx) // 204 & empty data
	} else if err != nil {
		(*h.response).RespErr(iris.StatusInternalServerError, data.ErrRepositoryOps, err.Error(), &ctx) // returning some other error may happen
	}
}

// createBook create a new book
// @Summary Create a new book
// @Description Create a new book from the passed data
// @Tags Books
// @Accept	json
// @Produce json
// @Param	book	body	dto.BookCreateIn	true	"Book Data"
// @Success 201 {object} models.Book "OK"
// @Success 422 {object} dto.ApiError "err.duplicate_key || Invalid data"
// @Failure 500 {object} dto.ApiError "err.repo_ops || Internal error"
// @Router /books [post]
func (h HBook) createBook(ctx iris.Context) {
	var bDto dto.BookCreateIn

	// TIP: use ctx.ReadBody(&bDto) to bind any type of incoming data instead. E.g it comes in handy when the client request are using form-data
	if e := ctx.ReadJSON(&bDto); e != nil {
		(*h.response).RespErr(iris.StatusUnprocessableEntity, data.ErrVal, e.Error(), &ctx) // 422 ReadJSON do the validation here
		return
	}

	// TIP ❗ this is just a sample to show the mapper use. I think the ReadJSON method do this trick when unmarshal the data.
	// So, the mapper use maybe has more sense using it for the data responses for hiding fields to the client.
	// Another mapper utility is fot use in combination of specifically defined strut for validation purpose
	// Mapping
	book := mapper.ToBookCreateV(&bDto)

	err := (*h.service).Create(book)
	if err != nil {

		if err.Error() == data.ErrDuplicateKey { 															// 422 Unprocessable 'cause duplicate key
			(*h.response).RespErr(iris.StatusUnprocessableEntity, data.ErrDuplicateKey, data.ErrDetDuplicateKey, &ctx)
		} else {																							// 500
			(*h.response).RespErr(iris.StatusInternalServerError, data.ErrRepositoryOps, err.Error(), &ctx)
		}

	} else {		// All good
		(*h.response).RespOKWithData(iris.StatusCreated, book, &ctx)
	}
}

// updateBook update the book having the Id passed as path parameter, with the data passed in the request body
// @Summary Update the indicated book
// @Description Update the book having the specified Id with the data passed in the request body
// @Tags Books
// @Accept	json
// @Produce json
// @Param 	id		path	int					true	"Book ID"	Format(uint32)
// @Param	book	body	dto.BookUpdateIn	true	"Book Data"
// @Success 200 {object} models.Book "OK"
// @Success 404 {object} dto.ApiError "err.not_found"
// @Success 422 {object} dto.ApiError "err.duplicate_key || Invalid data"
// @Failure 500 {object} dto.ApiError "err.repo_ops || Internal error"
// @Router /books/{id} [put]
func (h HBook) updateBook(ctx iris.Context) {
	var bDto dto.BookUpdateIn

	// Getting the data
	bDto.Id = ctx.Params().GetUintDefault("id", 0)
	if e := ctx.ReadJSON(&bDto); e != nil {
		(*h.response).RespErr(iris.StatusUnprocessableEntity, data.ErrVal, e.Error(), &ctx) // 422 errors may happen in the marshaling process
		return
	}

	// Mapping
	book := mapper.ToBookUpdateV(&bDto)

	// Updating
	updated, err := (*h.service).UpdateBook(book)

	if err != nil && err.Error() == data.ErrNotFound {													// 404 Wrong ID
		(*h.response).RespErr(iris.StatusNotFound, data.ErrNotFound, data.ErrDetNotFound, &ctx)
	} else if err != nil && err.Error() == data.ErrDuplicateKey {										// Same unique field, name in this case
		(*h.response).RespErr(iris.StatusUnprocessableEntity, data.ErrDuplicateKey, data.ErrDetDuplicateKey, &ctx)
	} else if err != nil {																				// Something happen
		(*h.response).RespErr(iris.StatusInternalServerError, data.ErrRepositoryOps, err.Error(), &ctx)
	} else if updated > 0 {																				// All good, bDto was updated
		(*h.response).RespOKWithData(iris.StatusOK, book, &ctx)
	}
}

// endregion =============================================================================
