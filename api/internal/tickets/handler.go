package tickets

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	usecase TicketUsecase
}

func NewTicketHandler(u TicketUsecase) *TicketHandler {
	return &TicketHandler{usecase: u}
}

func (h *TicketHandler) GetTickets(c *fiber.Ctx) error {
	tickets, err := h.usecase.GetTickets(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tickets)
}

func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	var body struct {
		Text string `json:"text"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	err := h.usecase.CreateTicket(context.Background(), body.Text)
	if err != nil {
		 log.Printf("Ошибка при создании тикета: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "created"})
}
