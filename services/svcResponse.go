package services

import (
	"github.com/kataras/iris/v12"
)
// We could create a structure so we don't have to pass the Iris Context in all the functions, but
// that way we have to create a different instance of this service in each endpoint wasting memory,
// also the context may be difference on each request, so I think is better (but annoying) to pass
// a pointer to the the context in each case.


// region ======== RESPONSES =============================================================

// ResponseOKWithData create response 200 with json data in to the context.
//
// - data [interface] ~ "Objet" to be marshalled in to the context.
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseOKWithData(data interface{}, ctx *iris.Context)  {
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
	if _, err := (*ctx).JSON(data);
		err != nil {																														// Logging *marshal* json if error occurs (come internally from iris)
		(*ctx).Application().Logger().Error(err.Error())
	}

	// (*ctx).StatusCode(iris.StatusOK)
}

// ResponseDeleteOK create response 204. It's delete confirmation wit empty retrieving data.
// So, the client don't have to expect eny data and, we reduce some traffic.
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseDeleteOK(ctx *iris.Context)  {
	(*ctx).StatusCode(iris.StatusNoContent)
}
// endregion =============================================================================

// region ======== ERROR RESPONSES =======================================================

// ResponseErr500 create and log an 'Internal Server Error' to the stdout and set it up in the request context.
// Also set the response status = 500, so we can respond the request accordingly (application/problem+json), thanks to the iris.Context use.
//
// - title [string] ~ The title of the problem. We prefer use a code that can be use with an i18n system in the client.
//
// - detail [string] ~ Error detail
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseErr500(title string, detail string, ctx *iris.Context) {

	if _, err := (*ctx).Problem(iris.NewProblem().Title(title).Status(iris.StatusInternalServerError).Detail(detail));
		err != nil {																														// Logging *marshal* json if error occurs. If happens it comes internally from iris
		(*ctx).Application().Logger().Error(err.Error())
	}

	// TODO log to a file
	// s.app.Logger().Warn(detail)
}

// ResponseErr404 create and log an 'Not Found' error to the stdout and set it up in the request context.
// Also set the response status = 404, so we can respond the request accordingly (application/problem+json), thanks to the iris.Context use.
//
// - title [string] ~ The title of the problem. We prefer use a code that can be use with an i18n system in the client.
//
// - detail [string] ~ Error detail
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseErr404(title string, detail string, ctx *iris.Context) {

	if _, err := (*ctx).Problem(iris.NewProblem().Title(title).Status(iris.StatusNotFound).Detail(detail));
		err != nil {																														// Logging *marshal* json if error occurs (come internally from iris)
		(*ctx).Application().Logger().Error(err.Error())
	}

	// TODO log to a file
	// s.app.Logger().Warn(detail)
}
// endregion =============================================================================




