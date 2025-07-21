package templates

const Controller = `package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type {{.Name}} struct {
	Placeholder string
}

func {{.Name}}Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func {{.Name}}Create(c *fiber.Ctx) error {
	p := new({{.Name}})

	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func {{.Name}}Show(c *fiber.Ctx) error {
	param := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func {{.Name}}Update(c *fiber.Ctx) error {
	param := c.Params("id")
	p := new({{.Name}})

	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func {{.Name}}Delete(c *fiber.Ctx) error {
	param := c.Params("id")
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
`