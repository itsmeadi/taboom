package useCase

import (
	"context"
	"errors"
	"github.com/omise/go-tamboon/src/custom/constants"
	"github.com/omise/go-tamboon/src/custom/customError"
	"github.com/omise/go-tamboon/src/entities"
	"github.com/omise/go-tamboon/src/repositries"
	"github.com/omise/omise-go"
	"sync"
	"time"
)

type ChargeStruct struct {
	File   repositries.FileSysInterface
	Token  repositries.TokenInterface
	Charge repositries.ChargeInterface
}

func InitCharge(charge ChargeStruct) ChargeStruct {
	return charge
}

func (chargeInst *ChargeStruct) ChargeOnUserInChannel(ctx context.Context, tran *entities.Transaction, tranCh chan entities.Transaction, wg *sync.WaitGroup) {

	var amount int64
	var err error
	chErr := make(chan error)

	go func() {
		amount, err = chargeInst.ChargeOnUser(ctx, tran.User)
		chErr <- err
	}()

	select {
	case <-ctx.Done():
		tran.Err = customError.ErrorTimeOut

	case err := <-chErr:
		tran.Err = err
		tran.Amount = amount
	}

	if err != nil && IsTransactionRepeatable(*tran) {
		chargeInst.RetryChargeOnUserInChannel(ctx, *tran, tranCh, wg)
	} else {
		tranCh <- *tran
		wg.Done()
	}

}

func (chargeInst *ChargeStruct) RetryChargeOnUserInChannel(ctx context.Context, tran entities.Transaction, amountCh chan entities.Transaction, wg *sync.WaitGroup) {

	time.Sleep(constants.RetryWaitTime)
	tran.Retry = tran.Retry + 1
	chargeInst.ChargeOnUserInChannel(ctx, &tran, amountCh, wg)
}

func (chargeInst *ChargeStruct) ChargeOnUser(ctx context.Context, user entities.UserInfo) (int64, error) {

	token, err := chargeInst.Token.GetToken(ctx, user)

	if err != nil {
		return 0, err
	}
	charge, err := chargeInst.Charge.Charge(ctx, user, token)
	if err != nil {
		return 0, err
	}

	if charge.Status == omise.ChargeSuccessful {
		return charge.Amount, nil
	} else {
		return 0, errors.New("something unexpected happened") //TODO
	}
}
