package grpc

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeProcessPaymentUsecase struct {
	inputs []usecase.ProcessPaymentInput
	err    error
}

func (u *fakeProcessPaymentUsecase) Execute(_ context.Context, input usecase.ProcessPaymentInput) (usecase.ProcessPaymentOutput, error) {
	u.inputs = append(u.inputs, input)
	if u.err != nil {
		return usecase.ProcessPaymentOutput{}, u.err
	}
	return usecase.ProcessPaymentOutput{PaymentID: uuid.New()}, nil
}

func TestCreatePaymentIntent_UsesRequestedAmount(t *testing.T) {
	uc := &fakeProcessPaymentUsecase{}
	h := NewPaymentServiceHandler(uc, nil, nil, nil, nil)

	resp, err := h.CreatePaymentIntent(context.Background(), &pb.CreatePaymentIntentRequest{
		OrderId: uuid.New().String(),
		Amount:  "2500",
	})
	if err != nil {
		t.Fatalf("CreatePaymentIntent returned error: %v", err)
	}
	if resp.GetPaymentId() == "" {
		t.Error("expected payment id in response")
	}
	if len(uc.inputs) != 1 {
		t.Fatalf("expected 1 usecase call, got %d", len(uc.inputs))
	}
	if uc.inputs[0].Amount != 2500 {
		t.Errorf("usecase amount = %d, want 2500 (リクエストの金額を使うこと)", uc.inputs[0].Amount)
	}
}

func TestCreatePaymentIntent_RejectsInvalidAmount(t *testing.T) {
	for _, amount := range []string{"", "abc", "-100", "0"} {
		uc := &fakeProcessPaymentUsecase{}
		h := NewPaymentServiceHandler(uc, nil, nil, nil, nil)

		_, err := h.CreatePaymentIntent(context.Background(), &pb.CreatePaymentIntentRequest{
			OrderId: uuid.New().String(),
			Amount:  amount,
		})
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("amount %q: expected InvalidArgument, got %v", amount, err)
		}
		if len(uc.inputs) != 0 {
			t.Errorf("amount %q: usecase should not be called", amount)
		}
	}
}

func TestCreatePaymentIntent_RejectsInvalidOrderID(t *testing.T) {
	uc := &fakeProcessPaymentUsecase{}
	h := NewPaymentServiceHandler(uc, nil, nil, nil, nil)

	_, err := h.CreatePaymentIntent(context.Background(), &pb.CreatePaymentIntentRequest{
		OrderId: "not-a-uuid",
		Amount:  "1000",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", err)
	}
}
