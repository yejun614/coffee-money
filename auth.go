package main

import (
	"os"
	"fmt"
	"errors"
	"coffee-money/models"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
)

// SignCheck godoc
// @Tags Auth
// @Summery 로그인 확인
// @Description 로그인 확인
// @Router /auth [get]
func SignCheck(c *fiber.Ctx) error {
	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: get username
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// OK
	return c.SendString(fmt.Sprintf("%v", username))
}

// SignIn godoc
// @Tags Auth
// @Summery 로그인
// @Description 로그인
// @Accept json
// @Param request body UserDTO true "body"
// @Router /auth [post]
func SignIn(c *fiber.Ctx) error {
	// payload: body parser
	payload := new(UserDTO)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// payload: validation
	if err := Validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// user
	user := models.User{
		Username: payload.Username,
		Password: SHA512(payload.Password),  // encryption by SHA512
	}

	// user: find user from DB
	if result := DB.Where(&user).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusUnauthorized)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// session
	sess, err := Store.Get(c)
	if err != nil {
		return err
	}

	// session: set user information
	sess.Set("auth", "SERVER")
	sess.Set("username", user.Username)
	sess.Set("ID", user.ID)

	// session: save
	if err := sess.Save(); err != nil {
		return err
	}

	// OK
	return nil
}

// SignWithGithub godoc
// @Tags Auth
// @Summery Github 로그인
// @Description Github 로그인
// @Router /auth/github [get]
func SignWithGithub(c *fiber.Ctx) error {
	// Github ClientID and Secrets from Environment Variables
	clientID := os.Getenv("SERVER_CLIENT_ID")

	// OK
	return c.SendString(fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s", clientID))
}

// CallbackSignWithGithub godoc
// @Tags Auth
// @Summery Github 로그인 Callback
// @Description Github 로그인 Callback
// @Param request body CallbackGithubDTO true "body"
// @Router /auth/github/callback [get]
func CallbackSignWithGithub(c *fiber.Ctx) error {
	// Github ClientID and Secrets from Environment Variables
	clientID := os.Getenv("SERVER_CLIENT_ID")
	clientSecrets := os.Getenv("SERVER_CLIENT_SECRETS")

	// payload: query parser
	payload := new(CallbackGithubDTO)
	if err := c.QueryParser(payload); err != nil {
		return err
	}

	// payload: validation
	if err := Validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// get github access token
	githubToken, err := GetGithubAccessToken(clientID, clientSecrets, payload.Code)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// get github user information
	githubUser, err := GetGithubAuthUser(githubToken.AccessToken)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// session
	sess, err := Store.Get(c)
	if err != nil {
		return err
	}

	// session: set user information
	sess.Set("auth", "GITHUB")
	sess.Set("username", githubUser.Name)
	sess.Set("ID", nil)

	// session: save
	if err := sess.Save(); err != nil {
		return err
	}

	// OK
	return nil
}

// SignOut godoc
// @Tags Auth
// @Summery 로그아웃
// @Description 로그아웃
// @Router /auth [delete]
func SignOut(c *fiber.Ctx) error {
	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: destroy
	if err := sess.Destroy(); err != nil {
		return err
	}

	// OK
	return nil
}
