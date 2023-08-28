package main

import (
	"os"
	"log"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"coffee-money/models"
	_ "coffee-money/docs"
	"github.com/sethvargo/go-password/password"
	"github.com/gofiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/go-playground/validator/v10"
)

var (
	err error
	DB *gorm.DB
	Validate = validator.New()
	Store = session.New()
)

// hello godoc
// @Summery hello, world!
// @Description hello, world!
// @Accept */*
// @Router /hello [get]
func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// @title Coffee Money API
// @version 0.0.1
// @description 카페에 맡겨둔 돈 관리해 주는 전자장부
func main() {
	// Database
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Database: auto migrate
	DB.AutoMigrate(
		&models.User{},
		&models.Ledger{},
	)

	// Generate ADMIN password
	adminPW, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ADMIN PW: %s\n", adminPW)

	// Create Fiber app
	app := fiber.New()

	// Middleware
	logFile, err := os.OpenFile("./server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	app.Use(logger.New(logger.Config{
	    Output: logFile,
	}))

	// Routing
	app.Get("/hello", Hello)

	// Routing: Admin
	admin := app.Group("/admin")
	admin.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": adminPW,
		},
	}))
	admin.Get("/swagger/*", swagger.HandlerDefault)
	admin.Get("/metrics", monitor.New())

	// Routing: User
	user := app.Group("/user")
	user.Post("", CreateUser)
	user.Patch("", ChangePasswordUser)
	user.Delete("", DeleteUser)

	// Routing: Authentication
	auth := app.Group("/auth")
	auth.Get("", SignCheck)
	auth.Post("", SignIn)
	auth.Get("/github", SignWithGithub)
	auth.Get("/github/callback", CallbackSignWithGithub)
	auth.Delete("", SignOut)

	// Routing: Ledger
	ledger := app.Group("/ledger")
	ledger.Post("", CreateLedger)
	ledger.Put("", UpdateLedger)
	ledger.Get("", GetAllLedger)
	ledger.Get("/item/:id", GetLedger)
	ledger.Get("/search", SearchLedger)
	ledger.Get("/filter/store/:store", FilterStoreLedger)
	ledger.Get("/filter/user/:username", FilterUserLedger)

	// Start the server
	log.Fatal(app.Listen("localhost:3000"))
}
