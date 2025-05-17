package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	sso "github.com/sol1corejz/keep-coin/internal/clients/sso/grpc"
	"github.com/sol1corejz/keep-coin/internal/domain/models"
	"time"
)

type Handler struct {
	SSOClient *sso.Client
}

func New(ssoClient *sso.Client) *Handler {
	return &Handler{
		SSOClient: ssoClient,
	}
}

func (h *Handler) RegisterGRPC(c *fiber.Ctx) error {
	log.Info("Register GRPC")
	var payloadData models.AuthRequest

	if err := c.BodyParser(&payloadData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), time.Second*5)
	defer cancel()

	resp, err := h.SSOClient.Register(ctx, payloadData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": resp,
	})
}

func (h *Handler) LoginGRPC(c *fiber.Ctx) error {
	log.Info("Login GRPC")
	var payloadData models.AuthRequest

	if err := c.BodyParser(&payloadData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.SSOClient.Login(ctx, payloadData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": resp,
	})
}
