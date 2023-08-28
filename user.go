package main

import (
	"fmt"
	"coffee-money/models"
	"github.com/gofiber/fiber/v2"
)

// CreateUser godoc
// @Tags user
// @Summery 새로운 사용자 생성
// @Description 새로운 사용자 생성
// @Accept json
// @Param request body UserDTO true "body"
// @Router /user [post]
func CreateUser(c *fiber.Ctx) error {
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
	}

	// user: find user from DB
	if result := DB.Where(&user).First(&user); result.Error == nil {
		return c.Status(fiber.StatusBadRequest).SendString("duplicate username")
	}

	// user: set password
	user.Password = SHA512(payload.Password)  // encryption by SHA512

	// user: create new user
	if result := DB.Create(&user); result.Error != nil {
		return result.Error
	}

	// OK
	return nil
}

// ChangePasswordUser godoc
// @Tags user
// @Summery 사용자 패스워드 변경
// @Description 사용자 패스워드 변경
// @Accept json
// @Param request body UserPasswordDTO true "body"
// @Router /user [patch]
func ChangePasswordUser(c *fiber.Ctx) error {
	// payload: body parser
	payload := new(UserPasswordDTO)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	// payload: validation
	if err := Validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: get data
	auth := sess.Get("auth")
	username := sess.Get("username")

	// authentication
	if auth == nil || fmt.Sprintf("%v", auth) != "SERVER" || username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// user
	user := models.User{
		Username: fmt.Sprintf("%v", username),
	}

	// user: find user from DB
	if result := DB.Where(&user).First(&user); result.Error != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// user: change password
	user.Password = SHA512(payload.Password)  // encryption by SHA512

	// user: save to DB
	if result := DB.Save(&user); result.Error != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// OK
	return nil
}

// DeleteUser godoc
// @Tags user
// @Summery 사용자 삭제
// @Description 사용자 삭제
// @Accept json
// @Router /user [delete]
func DeleteUser(c *fiber.Ctx) error {
	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: get data
	auth := sess.Get("auth")
	username := sess.Get("username")

	// authentication
	if auth == nil || fmt.Sprintf("%v", auth) != "SERVER" || username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// user
	user := models.User{
		Username: fmt.Sprintf("%v", username),
	}

	// user: delete a user from DB
	if result := DB.Where(&user).Delete(&user); result.Error != nil {
		return result.Error
	}

	// session: destroy
	if err := sess.Destroy(); err != nil {
		return err
	}

	// OK
	return nil
}
