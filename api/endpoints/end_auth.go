package endpoints

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"go.api.backend/schema/dto"
	"go.api.backend/service/auth"
	"go.api.backend/service/utils"
)

type HAuth struct {
	response *utils.SvcResponse
	// appConf *utils.SvcConfig
}

// NewAuthHandler create and register the authentication handlers for the App. For the moment, all the
// auth handlers emulates the Oauth2 "password" grant-type using the "client-credentials" flow.
//
// - app [*iris.Application] ~ Iris App instance
//
// - r [*utils.SvcResponse] ~ Response service instance
func NewAuthHandler (app *iris.Application, r *utils.SvcResponse, c *utils.SvcConfig) HAuth {

	// --- VARS SETUP ---
	h := HAuth{r}
	svcA := auth.NewSvcAuthentication([]string{"sisec"}, c)						// creating authentication Service

	authRouter := app.Party("/auth")									// authorize
	{
		// --- GROUP / PARTY MIDDLEWARES ---

		// --- DEPENDENCIES ---
		hero.Register(DepObtainUserCred)
		hero.Register(svcA)

		// --- REGISTERING ENDPOINTS ---
		// authRouter.Post("/<provider>")										// provider is the auth provider to be used.
		authRouter.Post("/sisec", hero.Handler(h.authIntentSISEC)) 		// using a provider named 'sisec'.
	}

	return h
}

// region ======== ENDPOINT HANDLERS =====================================================

// authIntentSISEC Try to make the authentication of the user credentials through the SISEC auth provider service
// @Summary Auth the user credential through SISEC
// @Description Try to make the authentication of the user credentials through the SISEC auth provider service
// @Tags Auth
// @Accept  multipart/form-data
// @Produce  json
// @Param credential body dto.UserCredIn true "User Login Credential"
// @Success 202 "Accepted"
// @Failure 401 {object} dto.ApiError "err.unauthorized"
// @Failure 500 {object} dto.ApiError "err.network || err.json_parse"
// @Router /auth/sisec [post]
func (h HAuth) authIntentSISEC(ctx iris.Context, uCred *dto.UserCredIn, authService *auth.SvcAuthentication) {

	grantData, err, eCode := authService.AuthProviders["sisec"].GrantIntent(uCred)

	if err != nil {
		(*h.response).RespErr(iris.StatusInternalServerError, err.Error(), eCode, &ctx)
	}

	ctx.JSON(grantData)
}
// endregion =============================================================================


// region ======== LOCAL DEPENDENCIES ====================================================

// DepObtainUserCred is used as dependencies to obtain / create the user credential from request body (multipart/form-data).
// It return a dto.UserCredIn struct
func DepObtainUserCred(ctx iris.Context) dto.UserCredIn {
	cred := dto.UserCredIn{}

	// Getting data
	cred.Username = ctx.PostValue("username")
	cred.Password = ctx.PostValue("password")
	cred.Scope = ctx.PostValue("scope")

	// TIP: We can do some validation here if we want
	return cred
}
// endregion =============================================================================