package app

import (
	"classifier/internal/tickets"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB
}

func NewApp(db *gorm.DB) *App {
	return &App{db: db}
}

func (a *App) Run() {
	app := fiber.New()

	repo := tickets.NewTicketRepository(a.db)
	uc := tickets.NewTicketUsecase(repo)
	handler := tickets.NewTicketHandler(uc)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Start")
	})
	app.Post("/tickets", handler.CreateTicket)
	app.Get("/tickets", handler.GetTickets)

	app.Listen(":8088")

}
