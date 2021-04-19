package utils

import (
	"github.com/kataras/iris/v12"
)


type SvcResponse struct {
	appConf *SvcConfig
}

// NewSvcResponse create a response service instance. Depends on the app configuration instance
//
// - appConf [*SvcConfig] ~ App conf instance pointer
func NewSvcResponse(appConf *SvcConfig) *SvcResponse {
	return &SvcResponse{appConf: appConf}
}

// region ======== OK RESPONSES ==========================================================

// ResWithDataStatus create response with specified status and specified data converted to json
// in to the context.
//
// - status [int] ~ Integer represent HTTP status for the response. iris.Status constants will be used
//
// - data [interface] ~ "Object" to be marshalled in to the context.
//
// - ctx [*iris.Context] ~ Iris Request context
func (s SvcResponse) ResWithDataStatus(status int, data interface{}, ctx *iris.Context)  {
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
	if _, err := (*ctx).JSON(data); err != nil {																									// Logging *marshal* json if error occurs (come internally from iris)
		(*ctx).Application().Logger().Error(err.Error())
	}
	(*ctx).StatusCode(status)
}

// ResWithDataStatus create response 200 with specified data converted to json in to the context.
//
// - data [interface] ~ "Object" to be marshalled in to the context.
//
// - ctx [*iris.Context] ~ Iris Request context
func (s SvcResponse) ResOKWithData(data interface{}, ctx *iris.Context) {
	if _, err := (*ctx).JSON(data); err != nil {																									// Logging *marshal* json if error occurs (come internally from iris)
		(*ctx).Application().Logger().Error(err.Error())
	}
	(*ctx).StatusCode(iris.StatusOK)
}

// ResOK create a response OK but with an empty content (204)
//
// - ctx [*iris.Context] ~ Iris Request context
func (s SvcResponse) ResOK(ctx *iris.Context)  {
	(*ctx).StatusCode(iris.StatusNoContent)
}

// ResDelete create response 204. It's delete confirmation wit empty retrieving data.
// So, the client don't have to expect eny data and, we reduce some traffic.
//
// - ctx [*iris.Context] ~ Iris Request context
func (s SvcResponse) ResDelete(ctx *iris.Context)  {
	(*ctx).StatusCode(iris.StatusNoContent)
}
// endregion =============================================================================

// region ======== ERROR RESPONSES =======================================================

// ResErr create and log an 'Error Response' to the stdout and setup the request context properly.
// Also set the response status = specific status code, so we can respond the request accordingly (application/problem+json).
// Ideally, this should be used for client error series (400s) or server error series (500)
//
// - status [int] ~ Integer represent HTTP status for the response. iris.Status constants will be used
//
// - title [string] ~ The title of the problem. We prefer use a code that can be use with an i18n system in the client.
//
// - detail [string] ~ Error detail
//
// - ctx [*iris.Context] ~ Iris Request context
func (s SvcResponse) ResErr(status int, title string, detail string, ctx *iris.Context) {
	d := detail

	// If the environment debug config isn't true then retrieve no details
	if s.appConf.Debug != true {d = ""}

	(*ctx).StopWithProblem(status, iris.NewProblem().Title(title).Detail(d))

	// TODO log to a file
	// s.app.Logger().Warn(detail)

	return
}
// endregion =============================================================================
