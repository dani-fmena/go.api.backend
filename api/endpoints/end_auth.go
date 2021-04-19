package endpoints

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.api.backend/lib"
	"go.api.backend/schema"
	"go.api.backend/schema/dto"
	"go.api.backend/schema/mapper"
	"go.api.backend/service/auth"
	"go.api.backend/service/utils"
)

type HAuth struct {
	response *utils.SvcResponse
	appConf *utils.SvcConfig
	providers map[string]bool
}

// NewAuthHandler create and register the authentication handlers for the App. For the moment, all the
// auth handlers emulates the Oauth2 "password" grant-type using the "client-credentials" flow.
//
// - app [*iris.Application] ~ Iris App instance
//
// - r [*utils.SvcResponse] ~ Response service instance
func NewAuthHandler (app *iris.Application, MdwAuthChecker *context.Handler, r *utils.SvcResponse, svcC *utils.SvcConfig) HAuth {

	// --- VARS SETUP ---
	h := HAuth{r, svcC, make(map[string]bool)}
	// filling providers
	h.providers["sisec"] = true
	// h.providers["another_provider"] = true
	// h.providers["another_provider"] = true

	svcA := auth.NewSvcAuthentication(h.providers, svcC) 					// creating authentication Service

	// registering unprotected router
	authRouter := app.Party("/auth")								// authorize
	{
		// --- GROUP / PARTY MIDDLEWARES ---

		// --- DEPENDENCIES ---
		hero.Register(depObtainUserCred)
		hero.Register(svcA)

		// --- REGISTERING ENDPOINTS ---
		// authRouter.Post("/<provider>")										// provider is the auth provider to be used.
		authRouter.Post("/{provider}", hero.Handler(h.authIntent)) 		// using a provider named 'sisec'.
	}

	// registering protected router
	protectedAuthRouter := app.Party("/auth")
	{
		// --- GROUP / PARTY MIDDLEWARES ---
		protectedAuthRouter.Use(*MdwAuthChecker)								// registering access token

		// --- DEPENDENCIES ---

		// --- REGISTERING ENDPOINTS ---
		protectedAuthRouter.Get("/protected", h.protectedSample)
		protectedAuthRouter.Get("/logout", h.logout)
	}

	return h
}

// region ======== ENDPOINT HANDLERS =====================================================

// protectedSample Sample protected endpoint
// @Summary Sample protected endpoint
// @Description This is a Bearer Token protected sample endpoint
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert access token" default(Bearer <Add access token here>)
// @Tags Auth
// @Produce  json
// @Success 200 {object} dto.AccessTokenData "OK"
// @Failure 401 {object} dto.ApiError "err.unauthorized"
// @Failure 500 {object} dto.ApiError "err.generic
// @Router /auth/protected [get]
func (h HAuth) protectedSample(ctx iris.Context) {

	// This 👇🏽 Open-Api declaration doesnt work with swaggo yet. Swaggo works only for Open-Api 2.0 https://github.com/swaggo/gin-swagger/issues/90 | https://github.com/swaggo/swag/issues/709 | https://stackoverflow.com/questions/32910065/how-can-i-represent-authorization-bearer-token-in-a-swagger-spec-swagger-j
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization

	claims := jwt.Get(ctx).(*dto.AccessTokenData)
	(*h.response).ResOKWithData(claims, &ctx)
}

// logout this endpoint invalidated a previously granted access token
// @Description This endpoint invalidated a previously granted access token
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert access token" default(Bearer <Add access token here>)
// @Tags Auth
// @Produce  json
// @Success 204 "OK"
// @Failure 401 {object} dto.ApiError "err.unauthorized"
// @Failure 500 {object} dto.ApiError "err.generic
// @Router /auth/logout [get]
func (h HAuth) logout(ctx iris.Context) {
	err := ctx.Logout()
	if err != nil {
		(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrGeneric, err.Error(), &ctx)
	} else {
		(*h.response).ResOK(&ctx)
	}
}

// authIntent Intent to grant authentication using the provider user's credentials and the specified  auth provider
// @Summary Auth the user credential through a provider
// @Description Intent to grant authentication using the provider user's credentials and the specified  auth provider
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param	provider	path	string			true	"Requested Book Id"
// @Param 	credential 	body 	dto.UserCredIn 	true	"User Login Credential"
// @Success 202 "Accepted"
// @Failure 401 {object} dto.ApiError "err.unauthorized"
// @Failure 400 {object} dto.ApiError "err.wrong_auth_provider"
// @Failure 504 {object} dto.ApiError "err.network"
// @Failure 500 {object} dto.ApiError "err.json_parse"
// @Router /auth/{provider} [post]
func (h HAuth) authIntent(ctx iris.Context, uCred *dto.UserCredIn, authService *auth.SvcAuthentication) {

	provider := ctx.Params().Get("provider")
	v, _ := h.providers[provider]									// v, ok := map[key]
	if !v {
		(*h.response).ResErr(iris.StatusBadRequest, schema.ErrWrongAuthProvider, schema.InvalidProvider, &ctx)
		return
	}

	// requesting authorization to SISEC with user credentials
	authGrantedData, e, eCode := authService.AuthProviders[provider].GrantIntent(uCred, &h.appConf)
	if e != nil && eCode == schema.ErrNetwork {
		(*h.response).ResErr(iris.StatusGatewayTimeout, eCode, e.Error(), &ctx)
	} else if e != nil && eCode == schema.ErrUnauthorized {
		(*h.response).ResErr(iris.StatusUnauthorized, eCode, e.Error(), &ctx)
	} else if e != nil {
		(*h.response).ResErr(iris.StatusInternalServerError, eCode, e.Error(), &ctx)
	}

	// if all good, we are going to create a token
	tokenData := mapper.ToAccessTokenDataV(&authGrantedData.Access_Token)
	accessToken, er := lib.MkAccessToken(tokenData, []byte(h.appConf.JWTSignKey), h.appConf.TkMaxAge)
	if er != nil {
		(*h.response).ResErr(iris.StatusInternalServerError, schema.ErrJwtGen, er.Error(), &ctx)
	}

	(*h.response).ResWithDataStatus(iris.StatusAccepted, string(accessToken), &ctx)
}
// endregion =============================================================================


// region ======== LOCAL DEPENDENCIES ====================================================

// depObtainUserCred is used as dependencies to obtain / create the user credential from request body (multipart/form-data).
// It return a dto.UserCredIn struct
func depObtainUserCred(ctx iris.Context) dto.UserCredIn {
	cred := dto.UserCredIn{}

	// Getting data
	cred.Username = ctx.PostValue("username")
	cred.Password = ctx.PostValue("password")
	cred.Domain = ctx.PostValue("domain")

	// TIP: We can do some validation here if we want
	return cred
}
// endregion =============================================================================