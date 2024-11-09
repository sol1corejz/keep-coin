package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sol1corejz/keep-coin/internal/auth"
	"github.com/sol1corejz/keep-coin/internal/logger"
	"github.com/sol1corejz/keep-coin/internal/models"
	"github.com/sol1corejz/keep-coin/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func RegisterHandler(c *fiber.Ctx) error {
	var request models.AuthRequest
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case <-ctx.Done():
		logger.Log.Warn("Context canceled or timeout exceeded")
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "Request timed out",
		})
	default:
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		userID := uuid.New()
		token, err := auth.GenerateToken(userID)
		if err != nil {
			logger.Log.Error("Error generating token: ", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Log.Error("Error hashing password: ", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		userData := models.User{
			Uuid:      userID,
			FirstName: "",
			LastName:  "",
			Email:     request.Email,
			Password:  string(hashedPassword),
		}

		err = storage.RegisterUser(&userData)
		if err != nil {
			logger.Log.Error("Error registering user: ", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(auth.TokenExp),
			HTTPOnly: true,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User registered successfully",
		})
	}
}

func LoginHandler(c *fiber.Ctx) error {
	var request models.AuthRequest
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	select {
	case <-ctx.Done():
		logger.Log.Warn("Context canceled or timeout exceeded")
		return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
			"error": "Request timed out",
		})
	default:
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		existingUser, err := storage.GetUserByLogin(request.Email)
		if err != nil {
			logger.Log.Error("Error while querying user: ", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		if existingUser.Uuid.String() == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong email or password",
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(request.Password))
		if err != nil {
			logger.Log.Error("Error while comparing hash: ", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong login or password",
			})
		}

		token, err := auth.GenerateToken(existingUser.Uuid)
		if err != nil {
			logger.Log.Error("Error generating token: ", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(auth.TokenExp),
			HTTPOnly: true,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User authorized successfully",
		})
	}
}
