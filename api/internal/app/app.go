package app

import (
	"classifier/internal/kafka"
	"classifier/internal/outbox"
	"classifier/internal/tickets"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	db       *gorm.DB
	producer *kafka.Producer
}

func NewApp(db *gorm.DB, producer *kafka.Producer) *App {
	return &App{db: db, producer: producer}
}

func (a *App) Run() {
	app := fiber.New()

	outboxRepo := outbox.NewOutboxRepository(a.db)
	ticketRepo := tickets.NewTicketRepository(a.db)
	usecace := tickets.NewTicketUsecase(ticketRepo, outboxRepo)
	handler := tickets.NewTicketHandler(usecace)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Start")
	})
	app.Post("/tickets", handler.CreateTicket)
	app.Get("/tickets", handler.GetTickets)

	outboxWorker := outbox.NewOutboxWorker(outboxRepo, a.producer)

	go func() {
		ctx := context.Background()
		log.Println("Outbox worker started...")
		outboxWorker.Run(ctx)
	}()

	log.Fatal(app.Listen(":8080"))

}
