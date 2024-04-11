package v1

import (
	"github.com/gin-gonic/gin"
	"go-gin/pkg/app"
	"go-gin/pkg/e"
	"go-gin/pkg/logging"
	"go-gin/service/userService"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6,max=30"`
}

type AddUserInput struct {
	LoginInput
	RePassword string `form:"rePassword" binding:"required,min=6,max=30"`
}

// AddUser godoc
// @Summary      添加用户
// @Description  添加用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body AddUserInput true "user obj"
// @Success      200
// @Router       /add_user [post]
func AddUser(ctx *gin.Context) {

	var userForm AddUserInput
	valid, validErr := app.BindAndValid(ctx, &userForm)

	logging.Info("%v", userForm)
	if !valid {
		logging.Error("app.BindAndValid errs: %v", validErr)
		app.Failed(ctx, e.ERROR_PARAMS, nil)
		return
	}

	if userForm.Password != userForm.RePassword {
		app.Failed(ctx, e.ERROR_ADD_USER_PASS_FAIL, nil)
		return
	}

	user := userService.User{
		Email:    userForm.Email,
		Password: userForm.Password,
	}

	if err := user.AddUser(); err != nil {
		app.Failed(ctx, e.ERROR_ADD_USER_FAIL, nil)
		return
	}

	app.Success(ctx, nil)
}

type GetUserInput struct {
	Email string `form:"email" binding:"required,email"`
}

// GetUser godoc
// @Summary      获取用户信息
// @Description  用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Security 	 Bearer 123456
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success      200
// @Router       /get_user [get]
func GetUser(ctx *gin.Context) {
	a, _ := ctx.Get("email")
	logging.Info("%v", a)

	app.Success(ctx, map[string]interface{}{"email": a})
}

// Login godoc
// @Summary      用户登录
// @Description  用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        loginInput body LoginInput true "GetUserInput"
// @Success      200
// @Router       /login [post]
func Login(ctx *gin.Context) {
	var loginInput LoginInput
	valid, validErr := app.BindAndValid(ctx, &loginInput)

	if !valid {
		logging.Error("app.BindAndValid errs: %v", validErr)
		app.Failed(ctx, e.ERROR_PARAMS, nil)
		return
	}

	user := userService.User{
		Email:    loginInput.Email,
		Password: loginInput.Password,
	}
	token, err := user.Login()
	if err != nil {
		app.Failed(ctx, e.ERROR_LOGIN_USER_FAIL, nil)
		return
	}

	if token.Token != "" {
		app.Success(ctx, token)
	}
}
