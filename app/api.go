package app

import (
	"base/app/controller/facilities"
	questionset "base/app/controller/question-set"
	"base/database"

	"github.com/gofiber/fiber/v2"
)

func InitRoute(app *fiber.App) {
	DB := database.GetDB()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Question set
	qs := v1.Group("/qs")
	qs.Post("/", questionset.Create(DB))
	qs.Delete("/:id", questionset.Delete(DB))
	qs.Get("/:id", questionset.GetById(DB))
	qs.Get("/", questionset.GetAll(DB))
	qs.Put("/:id", questionset.UpdateById(DB))

	// facilities
	facilitiesRoutes := v1.Group("/facilities")
	facilitiesRoutes.Post("/", facilities.Create(DB))
	facilitiesRoutes.Delete("/:id", facilities.Delete(DB))
	facilitiesRoutes.Put("/:id", facilities.Update(DB))

	facilitiesRoutes.Get("/", facilities.GetAll(DB))
	facilitiesRoutes.Get("/:id", facilities.GetById(DB))
}
