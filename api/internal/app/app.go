package app

import (
	"classifier/internal/kafka"
	"classifier/internal/outbox"
	"classifier/internal/tickets"
	"classifier/internal/worker"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	db       *gorm.DB
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewApp(db *gorm.DB, producer *kafka.Producer, consumer *kafka.Consumer) *App {
	return &App{db: db, producer: producer, consumer: consumer}
}

func (a *App) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signCh
		log.Println("Caught shutdown signal, stopping gracefully...")
		cancel()
	}()

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
	classifierUsecase := tickets.NewTicketClassifierUsecase(ticketRepo)
	consumerWorker := worker.NewTicketConsumerWorker(a.consumer, classifierUsecase)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		consumerWorker.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		outboxWorker.Run(ctx)
	}()

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Println("Fiber stopped:", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down application...")

	_ = app.Shutdown()
	wg.Wait() 
	log.Println("Application stopped gracefully")

}
