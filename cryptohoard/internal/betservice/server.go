package betservice

import (
	fmt "fmt"
	math "math"
	"sync"
	"time"

	"github.com/cryptohoard/experimental/cryptohoard/internal/helper"
	"github.com/cryptohoard/experimental/cryptohoard/internal/pricetracker"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/golang/protobuf/ptypes/timestamp"
	context "golang.org/x/net/context"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

type BetSvc struct {
	logger          log.Logger
	pt              *pricetracker.PriceTracker
	Bets            sync.Map
	productToStrMap map[Product]string
}

func NewBetSvc(logger log.Logger, pt *pricetracker.PriceTracker) *BetSvc {
	return &BetSvc{
		logger: logger,
		pt:     pt,
		productToStrMap: map[Product]string{
			Product_BTCUSD: pricetracker.BTCUSD,
			Product_BCHUSD: pricetracker.BCHUSD,
			Product_ETHUSD: pricetracker.ETHUSD,
			Product_LTCUSD: pricetracker.LTCUSD,
		},
	}
}

func (betsvc *BetSvc) startbetworker(bet *Bet) {
	method := "BetSvc.startbetworker"

	pr := betsvc.pt.Register(betsvc.productToStrMap[bet.Product])
	p := pr.Read()
	level.Debug(betsvc.logger).Log(
		"method",
		method,
		"msg",
		fmt.Sprintf("price: %+v", p),
	)
	bet.State = State_PLACED
	bet.Price = p.Price
	bet.CryptoCurrency = (bet.InitialAmount / p.Price)
	bet.CurrentAmount = toFixed((bet.CryptoCurrency * p.Price), 2)
	bet.ProfitPercent = toFixed((bet.CurrentAmount / bet.InitialAmount), 2)
	if bet.CurrentAmount < bet.InitialAmount {
		bet.ProfitPercent = (bet.ProfitPercent * -1)
	}
	betsvc.Bets.Store(bet.BetId, bet)
	level.Debug(betsvc.logger).Log(
		"method",
		method,
		"msg",
		fmt.Sprintf("bet: %s", bet),
	)
	betsvc.betworker(bet, pr)
}

func (betsvc *BetSvc) betworker(bet *Bet, pr *pricetracker.PriceReceiver) {
	method := "BetSvc.betworker"

	for {
		p := pr.Read()
		level.Debug(betsvc.logger).Log(
			"method",
			method,
			"msg",
			fmt.Sprintf("price: %+v", p),
		)
		bet.Price = p.Price
		bet.CurrentAmount = toFixed((bet.CryptoCurrency * p.Price), 2)
		bet.ProfitPercent = toFixed((bet.CurrentAmount / bet.InitialAmount), 2)
		if bet.CurrentAmount < bet.InitialAmount {
			bet.ProfitPercent = (bet.ProfitPercent * -1)
		}
		betsvc.Bets.Store(bet.BetId, bet)
		level.Debug(betsvc.logger).Log(
			"method",
			method,
			"msg",
			fmt.Sprintf("bet: %s", bet),
		)
	}
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
		send := true

		if filter.Product != Product_PRODUCT_NONE &&
			filter.Product != bet.Product {
			send = false
		}

		if filter.State != State_STATE_NONE &&
			filter.State != bet.State {
			send = false
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

func (betsvc *BetSvc) validateRequest(betReq *PlaceBetRequest) error {
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

func (betsvc *BetSvc) PlaceBet(
	ctx context.Context,
	betReq *PlaceBetRequest,
) (*BetServiceResponse, error) {
	method := "BetSvc.PlaceBet"

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
		InitialAmount: toFixed(betReq.Amount, 2),
		ProfitPercent: toFixed(betReq.ProfitPercent, 2),
		LossPercent:   toFixed(betReq.LossPercent, 2),
		State:         State_PROCESSING,
		CreationTime:  &timestamp.Timestamp{Seconds: int64(time.Now().Unix())},
	}

	if bet.Exchange == "" {
		bet.Exchange = "GDAX"
	}

	if bet.ProfitPercent == float64(0) {
		bet.ProfitPercent = toFixed(float64(25), 2)
	}

	if bet.LossPercent == float64(0) {
		bet.LossPercent = toFixed(float64(10), 2)
	}

	level.Debug(betsvc.logger).Log(
		"method",
		method, "msg",
		fmt.Sprintf("%+v", bet),
	)

	betsvc.Bets.Store(bet.BetId, bet)

	go betsvc.startbetworker(bet)

	return &BetServiceResponse{Success: true, Bet: bet}, nil
}

func (betsvc *BetSvc) CashoutBet(
	ctx context.Context,
	cashoutReq *CashoutBetRequest,
) (*BetServiceResponse, error) {
	method := "BetSvc.CashoutBet"

	if cashoutReq.BetId == "" {
		err := fmt.Errorf("invalide bet ID")
		level.Error(betsvc.logger).Log("method", method, "msg", err)
		return nil, err
	}

	_bet, ok := betsvc.Bets.Load(cashoutReq.BetId)
	if !ok {
		err := fmt.Errorf("bet ID %s not present", cashoutReq.BetId)
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

	bet.State = State_CASHOUT
	bet.ClosedTime = &timestamp.Timestamp{Seconds: int64(time.Now().Unix())}

	return &BetServiceResponse{Success: true, Bet: bet}, nil
}

func (betsvc *BetSvc) PayoutBet(
	ctx context.Context,
	payoutReq *PayoutBetRequest,
) (*BetServiceResponse, error) {
	method := "BetSvc.PayoutBet"

	if payoutReq.BetId == "" {
		err := fmt.Errorf("invalide bet ID")
		level.Error(betsvc.logger).Log("method", method, "msg", err)
		return nil, err
	}

	_bet, ok := betsvc.Bets.Load(payoutReq.BetId)
	if !ok {
		err := fmt.Errorf("bet ID %s not present", payoutReq.BetId)
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

	bet.State = State_PAYOUT
	bet.ClosedTime = &timestamp.Timestamp{Seconds: int64(time.Now().Unix())}

	return &BetServiceResponse{Success: true, Bet: bet}, nil
}
