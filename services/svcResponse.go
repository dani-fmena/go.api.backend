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
// - status [int] ~ Integer represent HTTP status for the response. iris.Status constants will be used
//
// - data [interface] ~ "Object" to be marshalled in to the context.
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseOKWithData(status int, data interface{}, ctx *iris.Context)  {
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
	if _, err := (*ctx).JSON(data);
		err != nil {																									// Logging *marshal* json if error occurs (come internally from iris)
		(*ctx).Application().Logger().Error(err.Error())
	}
	(*ctx).StatusCode(status)
}

// ResponseDelete create response 204. It's delete confirmation wit empty retrieving data.
// So, the client don't have to expect eny data and, we reduce some traffic.
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseDelete(ctx *iris.Context)  {
	(*ctx).StatusCode(iris.StatusNoContent)
}
// endregion =============================================================================

// region ======== ERROR RESPONSES =======================================================

// ResponseErr create and log an 'Error Response' to the stdout and set it up in the request context.
// Also set the response status = specific status code, so we can respond the request accordingly (application/problem+json), thanks to the iris.Context use.
// Ideally, this should be used for client error series (400s) or server error series (500)
//
// - status [int] ~ Integer represent HTTP status for the response. iris.Status constants will be used
//
// - title [string] ~ The title of the problem. We prefer use a code that can be use with an i18n system in the client.
//
// - detail [string] ~ Error detail
//
// - ctx [*iris.Context] ~ Iris Request context
func ResponseErr(status int, title string, detail string, ctx *iris.Context) {
	(*ctx).StopWithProblem(status, iris.NewProblem().Title(title).Detail(detail))

	// TODO if some debug var is set, don't send any detail ever in the response. Maybe you have to make the endpoint into structs with funcion to paas the config var easy

	// TODO log to a file
	// s.app.Logger().Warn(detail)

	return
}
// endregion =============================================================================
