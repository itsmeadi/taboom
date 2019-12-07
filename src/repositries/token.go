package repositries

import (
	"context"
	"github.com/omise/go-tamboon/src/custom/constants"
	"github.com/omise/go-tamboon/src/custom/rate"
	"github.com/omise/go-tamboon/src/entities"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"time"
)

type TokenInterface interface {
	GetToken(ctx context.Context, user entities.UserInfo) (omise.Token, error)
}

type Token struct {
	Client      *omise.Client
	Timeout     time.Duration
	rateLimiter *rate.RateLimit
}

var _ TokenInterface = Token{}	//just to make sure interface implementation is right
var tokenService Token

//var rateLimiter rate.RateLimit

func GetTokenService() Token {
	return tokenService
}

func InitToken(timeOut time.Duration, rateNum int, rateLimitTime time.Duration) error {
	client, err := omise.NewClient(constants.OmisePublicKey, constants.OmiseSecretKey)
	if err != nil {
		return err
	}
	client.Timeout = timeOut
	rateLimiter := rate.InitRateLimiter(rateLimitTime, rateNum)

	tokenService = Token{
		Client:      client,
		Timeout:     timeOut,
		rateLimiter: &rateLimiter,
	}
	return nil
}

func (service Token) GetToken(ctx context.Context, user entities.UserInfo) (omise.Token, error) {

	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:            user.Name,
		Number:          string(user.CCNumber), //need to convert it to string, since thats what the API expects
		ExpirationMonth: time.Month(user.ExpMonth),
		ExpirationYear:  int(user.ExpYear),
		SecurityCode:    string(user.CVV), //need to convert it to string, since thats what the API expects
	}

	service.rateLimiter.Wait()
	err := service.Client.Do(token, createToken)

	return *token, err
}
