package main

import (
	"fmt"
	"errors"
	"coffee-money/models"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
)

// CreateLedger godoc
// @Tags ledger
// @Summery 가계부 데이터 추가
// @Description 가계부 데이터 추가
// @Accept json
// @Produce json
// @Param request body CreateLedgerDTO true "body"
// @Router /ledger [post]
func CreateLedger(c *fiber.Ctx) error {
	// payload: body parser
	payload := new(CreateLedgerDTO)
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
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// ledger
	ledger := models.Ledger{
		StoreName: payload.StoreName,
		Balance: payload.Balance,
		Description: payload.Description,
		Username: fmt.Sprintf("%v", username),
		IsDisabled: false,
	}

	// ledger: create
	if result := DB.Create(&ledger); result.Error != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// OK
	return c.JSON(ledger)
}

// UpdateLedger godoc
// @Tags ledger
// @Summery 가계부 데이터 수정
// @Description 가계부 데이터 수정
// @Accept json
// @Produce json
// @Param request body UpdateLedgerDTO true "body"
// @Router /ledger [put]
func UpdateLedger(c *fiber.Ctx) error {
	// payload: body parser
	payload := new(UpdateLedgerDTO)
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
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// ledger
	ledger := models.Ledger{}

	// ledger: find item from DB
	if result := DB.First(&ledger, payload.ID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusUnauthorized)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// ledger: update values
	ledger.ID = payload.ID
	ledger.StoreName = payload.StoreName
	ledger.Balance = payload.Balance
	ledger.Description = payload.Description
	ledger.IsDisabled = payload.IsDisabled

	// ledger: save
	if result := DB.Save(&ledger); result.Error != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// OK
	return c.JSON(ledger)
}

// GetAllLedger godoc
// @Tags ledger
// @Summery 가계부 전체 데이터 조회
// @Description 가계부 전체 데이터 조회
// @Accept json
// @Produce json
// @Router /ledger [get]
func GetAllLedger(c *fiber.Ctx) error {
	// payload: query parser
	payload := new(SearchLedgerDTO)
	if err := c.QueryParser(payload); err != nil {
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
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// ledgers
	var ledgers []models.Ledger

	// ledgers: find items from DB
	if result := DB.Where(models.Ledger{
		IsDisabled: false,
	}).Order("created_at").Find(&ledgers); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// OK
	return c.JSON(ledgers)
}

// GetLedger godoc
// @Tags ledger
// @Summery 가계부 데이터 조회
// @Description 가계부 데이터 조회
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Router /ledger/item/{id} [get]
func GetLedger(c *fiber.Ctx) error {
	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: get data
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// params: get ID
	id := c.Params("id")

	// ledger
	ledger := models.Ledger{}

	// ledger: find item from DB
	if result := DB.First(&ledger, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// OK
	return c.JSON(ledger)
}

// SearchLedger godoc
// @Tags ledger
// @Summery 가계부 데이터 검색
// @Description 가계부 데이터 검색
// @Accept json
// @Produce json
// @Param store_name query string false "Store Name"
// @Param balance_begin query string false "Balance Begin"
// @Param balance_end query string false "Balance End"
// @Param description query string false "Description"
// @Param is_disable query string false "IsDisabled"
// @Param useranme query string false "Username"
// @Param created_at_begin query string false "Created At Begin"
// @Param created_at_end query string false "Created At End"
// @Param updated_at_begin query string false "Updated At Begin"
// @Param updated_at_end query string false "Updated At End"
// @Router /ledger/search [get]
func SearchLedger(c *fiber.Ctx) error {
	// payload: query parser
	payload := new(SearchLedgerDTO)
	if err := c.QueryParser(payload); err != nil {
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
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// make where query
	where := "(store_name LIKE ?)"
	where += " OR (balance >= ? AND balance <= ?)"
	where += " OR (description LIKE ?)"
	where += " OR (is_disable = ?)"
	where += " OR (username = ?)"
	where += " OR (created_at >= ? AND created_at <= ?)"
	where += " OR (updated_at >= ? AND updated_at <= ?)"

	// ledgers
	var ledgers []models.Ledger

	// ledgers: find items from DB
	if result := DB.Where(
		where,
		payload.StoreName,
		payload.BalanceBegin,
		payload.BalanceEnd,
		payload.Description,
		payload.IsDisabled,
		payload.Username,
		payload.CreatedAtBegin,
		payload.CreatedAtEnd,
		payload.UpdatedAtBegin,
		payload.UpdatedAtEnd,
	).Find(&ledgers); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// OK
	return c.JSON(ledgers)
}

// FilterStoreLedger godoc
// @Tags ledger
// @Summery 가계부 데이터 조회 (가계 기준 검색)
// @Description 가계부 데이터 조회 (가계 기준 검색)
// @Accept json
// @Produce json
// @Param store path string true "Store Name"
// @Router /ledger/filter/store/{store} [get]
func FilterStoreLedger(c *fiber.Ctx) error {
	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: get data
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// params: get ID
	storeName := c.Params("store")

	// ledger
	var ledgers []models.Ledger

	// ledger: find item from DB
	if result := DB.Where(models.Ledger{
		StoreName: storeName,
		IsDisabled: false,
	}).Order("created_at").Find(&ledgers); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// OK
	return c.JSON(ledgers)
}

// FilterUserLedger godoc
// @Tags ledger
// @Summery 가계부 데이터 조회 (사용자 기준 검색)
// @Description 가계부 데이터 조회 (사용자 기준 검색)
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Router /ledger/filter/user/{username} [get]
func FilterUserLedger(c *fiber.Ctx) error {
	// session
	sess, err := Store.Get(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// session: get data
	username := sess.Get("username")

	// authentication
	if username == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// params: get ID
	paramUsername := c.Params("username")

	// ledger
	var ledgers []models.Ledger

	// ledger: find item from DB
	if result := DB.Where(models.Ledger{
		Username: paramUsername,
		IsDisabled: false,
	}).Order("created_at").Find(&ledgers); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// OK
	return c.JSON(ledgers)
}