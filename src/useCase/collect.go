package useCase

import (
	"context"
	"github.com/omise/go-tamboon/src/custom/constants"
	"github.com/omise/go-tamboon/src/custom/customError"
	"github.com/omise/go-tamboon/src/entities"
	"sync"
)

func (chargeInst *ChargeStruct) ChargeAllUser(ctx context.Context, userCh chan entities.UserInfo) (entities.Result, error) {

	transactionCh := make(chan entities.Transaction, 30)

	var wg sync.WaitGroup

	go chargeInst.Producer(ctx, userCh, transactionCh, &wg)
	return chargeInst.CollectResults(ctx, userCh, transactionCh)

}

func (chargeInst *ChargeStruct) Producer(ctx context.Context, userCh chan entities.UserInfo, transactionCh chan entities.Transaction, wg *sync.WaitGroup) {

	for user := range userCh {

		wg.Add(1)
		go chargeInst.ChargeOnUserInChannel(ctx, &entities.Transaction{User: user}, transactionCh, wg)
	}
	wg.Wait()
	close(transactionCh)
}

//result
func (chargeInst *ChargeStruct) CollectResults(ctx context.Context, userCh chan entities.UserInfo, transactionCh chan entities.Transaction) (entities.Result, error) {
	var ttl int64
	var ttlDonation, successDonation int64
	topUsers := make([]entities.UserInfo, 0)
	for transaction := range transactionCh {

		ttlDonation = ttlDonation + transaction.User.Amount
		ttl = ttl + 1

		if transaction.Err != nil {
			//log.Printf("Failed Transaction after %+v attempts=%+v Error=%+v", transaction.Retry+1, transaction, transaction.Err)
			continue //Transaction failed after multiple attempts, can push user to new file to retry later
		}

		successDonation = successDonation + transaction.Amount

		Top3Users(&topUsers, transaction)
		//log.Println(transaction)
	}

	var result entities.Result

	if ttl == 0 {
		return result, customError.ErrorZeroUserFound
	}
	result.TotalDonation = ttlDonation
	result.AvgDonation = successDonation / ttl
	result.FaultyDonation = ttlDonation - successDonation
	result.SuccessfulDonation = successDonation
	result.TopDonors = topUsers
	result.Currency = constants.BKCurrency

	return result, nil
}

func IsTransactionRepeatable(tran entities.Transaction) bool {
	if customError.IsInvalidCard(tran.Err) {
		return false
	}
	if tran.Retry <= constants.RetryLimit {
		return true
	}
	return false
}

func Top3Users(userArr *[]entities.UserInfo, trans entities.Transaction) {
	if len(*userArr) < 3 {
		*userArr = append(*userArr, trans.User)
		return
	}
	for i := 0; i < 3; i++ {
		if trans.Amount > (*userArr)[i].Amount {
			(*userArr)[i] = trans.User
			return
		}
	}
}
