package betservice

import (
	fmt "fmt"
	"sync"
	"time"

	"github.com/cryptohoard/experimental/cryptohoard/internal/helper"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/protobuf/ptypes/timestamp"
	context "golang.org/x/net/context"
)

type BetSvc struct {
	logger log.Logger
	Bets   sync.Map
}

func NewBetSvc(logger log.Logger) *BetSvc {
	return &BetSvc{logger: logger}
}

func (betsvc *BetSvc) ListBets(
	filter *BetFilter,
	stream BetService_ListBetsServer,
) error {
	method := "BetSvc.ListBets"

	f := func(_k, _v interface{}) bool {
		if _, ok := _v.(*Bet); !ok {
			return false
		}

		bet := _v.(*Bet)
		send := false

		if filter.Product != Product_PRODUCT_NONE &&
			filter.Product == bet.Product {
			send = true
		}

		if filter.State != State_STATE_NONE &&
			filter.State == bet.State {
			send = true
		}

		if send {
			if err := stream.Send(bet); err != nil {
				level.Error(betsvc.logger).Log(
					"method",
					method,
					"msg",
					fmt.Sprintf("failed to send bets, error: %v", err),
				)

				return false
			}
		}

		return true
	}

	betsvc.Bets.Range(f)

	return nil
}

func (betsvc *BetSvc) validateRequest(betReq *BetServiceRequest) error {
	if betReq.CustomerId == "" {
		return fmt.Errorf("invalid customer ID")
	}

	if betReq.Product == Product_PRODUCT_NONE {
		return fmt.Errorf("invalid product")
	}

	if betReq.Amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	return nil
}

func (betsvc *BetSvc) CreateBet(
	ctx context.Context,
	betReq *BetServiceRequest,
) (*BetServiceResponse, error) {
	method := "BetSvc.CreateBet"

	if err := betsvc.validateRequest(betReq); err != nil {
		level.Error(betsvc.logger).Log("method", method, "msg", err)
		return nil, err
	}

	bet := &Bet{
		BetId: helper.CreateID(
			fmt.Sprintf(
				"%s%d%d",
				betReq.CustomerId,
				betReq.Product,
				betReq.Amount,
			),
		),
		CustomerId:    betReq.CustomerId,
		Exchange:      betReq.Exchange,
		Product:       betReq.Product,
		InitialAmount: betReq.Amount,
		ProfitPercent: betReq.ProfitPercent,
		LossPercent:   betReq.LossPercent,
		State:         State_PROCESSING,
		CreationTime:  &timestamp.Timestamp{Seconds: int64(time.Now().Unix())},
	}

	if bet.Exchange == "" {
		bet.Exchange = "GDAX"
	}

	if bet.ProfitPercent == 0 {
		bet.ProfitPercent = 25
	}

	if bet.LossPercent == 0 {
		bet.LossPercent = 10
	}

	level.Debug(betsvc.logger).Log(
		"method",
		method, "msg",
		fmt.Sprintf("%+v", bet),
	)

	betsvc.Bets.Store(bet.BetId, bet)

	return &BetServiceResponse{Success: true, Bet: bet}, nil
}

func (betsvc *BetSvc) CancelBet(
	ctx context.Context,
	cancelReq *BetCancelRequest,
) (*BetServiceResponse, error) {
	method := "BetSvc.CancelBet"

	if cancelReq.BetId == "" {
		err := fmt.Errorf("invalide bet ID")
		level.Error(betsvc.logger).Log("method", method, "msg", err)
		return nil, err
	}

	_bet, ok := betsvc.Bets.Load(cancelReq.BetId)
	if !ok {
		err := fmt.Errorf("bet ID %s not present", cancelReq.BetId)
		level.Error(betsvc.logger).Log("method", method, "msg", err)
		return nil, err
	}

	_, ok = _bet.(*Bet)
	if !ok {
		err := fmt.Errorf("fatal error not a Bet object")
		level.Error(betsvc.logger).Log("method", method, "msg", err)
		return nil, err
	}

	bet := _bet.(*Bet)
	level.Debug(betsvc.logger).Log(
		"method",
		method,
		"msg",
		fmt.Sprintf("%+v", bet),
	)

	return &BetServiceResponse{Success: true, Bet: bet}, nil
}
