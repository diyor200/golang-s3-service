package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang-diplom-work/internal/service"
	"net/http"
)

type authRoutes struct {
	authService service.Auth
}

// func newAuthRoutes(g *echo.Group, authService service.Auth) {
func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}
	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

type singUpInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (r *authRoutes) signUp(c echo.Context) error {
	var input singUpInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		ID int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		ID: id,
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// @Summary Sign in
// @Description Sign in
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signInInput true "input"
// @Success 200 {object} v1.authRoutes.signIn.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/sign-in [post]
func (r *authRoutes) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	log.Println(input)
	if err := c.Validate(input); err != nil {
		log.Println(err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	log.Println("token=", token)

	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusBadRequest, "invalid username or password")
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, response{
		Token: token,
	})
}
