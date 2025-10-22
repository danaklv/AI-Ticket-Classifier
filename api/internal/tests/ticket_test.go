package tests_test

import (
	models "classifier/internal/models"
	"classifier/internal/tests"
	"classifier/internal/tickets"
	context "context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateR(t *testing.T) {
	cntrl := gomock.NewController(t)

	defer cntrl.Finish()

	mockRepo := tests.NewMockTicketRepository(cntrl)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	service := tickets.NewTicketUsecase(mockRepo)

	model := &models.Ticket{
		Text: "Hello World",
	}

	err := service.CreateTicket(context.Background(), model.Text)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
