package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(WebResponseHelper(true, message, data, nil, nil))
}

func okList(c *fiber.Ctx, message string, data any, meta any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true, "message": message, "data": data, "meta": meta,
	})
}

func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(WebResponseHelper(true, message, data, nil, nil))
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(WebResponseHelper(false, message, nil, nil, nil))
}

func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(WebResponseHelper(false, "validasi gagal", nil, nil, errs))
}

// WebResponseHelper bikin bentuk map biar gak perlu import package model di sini.
func WebResponseHelper(success bool, message string, data any, meta any, errs any) fiber.Map {
	m := fiber.Map{"success": success, "message": message}
	if data != nil {
		m["data"] = data
	}
	if meta != nil {
		m["meta"] = meta
	}
	if errs != nil {
		m["errors"] = errs
	}
	return m
}

var allowedSort = map[string]bool{
	"id": true, "nim": true, "name": true, "grade": true, "created_at": true,
}

type ListQueryPlain struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

// reqCtx memberi batas waktu untuk setiap operasi basis data.
func reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func trimAll(values ...*string) {
	for _, v := range values {
		*v = strings.TrimSpace(*v)
	}
}