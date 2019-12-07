package repositries

import (
	"context"
	"github.com/omise/go-tamboon/src/custom/constants"
	"github.com/omise/go-tamboon/src/custom/rate"
	"github.com/omise/go-tamboon/src/entities"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"sync"
	"time"
)

type ChargeInterface interface {
	Charge(ctx context.Context, user entities.UserInfo, token omise.Token) (omise.Charge, error)
}

type Charge struct {
	Client      *omise.Client
	Timeout     time.Duration
	rateLimiter *rate.RateLimit
}

var _ ChargeInterface = Charge{}//making sure it implements correctly
var chargeService Charge

func GetChargeService() Charge {
	return chargeService
}

func InitCharge(timeout time.Duration, rateNum int, rateLimitTime time.Duration) error {
	client, err := omise.NewClient(constants.OmisePublicKey, constants.OmiseSecretKey)
	if err != nil {
		return err
	}
	client.Timeout = timeout

	rateLimiter := rate.InitRateLimiter(rateLimitTime, rateNum)

	chargeService = Charge{
		Client:      client,
		Timeout:     timeout,
		rateLimiter: &rateLimiter,
	}
	return nil
}

var lock1 sync.Mutex

func (service Charge) Charge(ctx context.Context, user entities.UserInfo, token omise.Token) (omise.Charge, error) {

	// Creates a charge from the token
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   user.Amount, // ฿ 1,000.00
		Currency: "thb",
		Card:     token.ID,
	}

	service.rateLimiter.Wait()

	err := service.Client.Do(charge, createCharge)

	return *charge, err
}
