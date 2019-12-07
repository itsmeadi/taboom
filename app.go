package main

import (
	"context"
	"github.com/omise/go-tamboon/src/custom/constants"
	"github.com/omise/go-tamboon/src/repositries"
	"github.com/omise/go-tamboon/src/useCase"

	"log"
	"os"
	"runtime"
)

func main() {

	filePath := os.Args[1]

	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()
	fileRepo := repositries.InitFile(constants.FileReadBufferSize, constants.FileReadUserChannelSize)

	err := repositries.InitToken(constants.TimeOut, constants.RateLimitTokenAPI, constants.RateLimitTimeTokenAPI)
	if err != nil {
		log.Fatal("[Main]Cannot init Token Repo err=", err)
	}
	tokenRepo := repositries.GetTokenService()

	err = repositries.InitCharge(constants.TimeOut, constants.RateLimitChargeAPI, constants.RateLimitTimeChargeAPI)
	if err != nil {
		log.Fatal("[Main]Cannot init Charge Repo err=", err)
	}
	chargeRepo := repositries.GetChargeService()

	chargeUseCase := useCase.InitCharge(useCase.ChargeStruct{
		Token:  tokenRepo,
		Charge: chargeRepo,
	})

	userUseCase := useCase.InitUserUsecase(useCase.UserStruct{Repo: fileRepo})

	//StartToPushDataToUserChannel will read the filePath and start pushing the data to fileRepo.UserChannel channel
	err = userUseCase.StartToPushDataToUserChannel(filePath)
	if err != nil {
		log.Fatal("[Main]Cannot Read from userchannel err=", err)
	}
	<-fileRepo.UserChannel //Read out Header

	log.Println("performing donations...")
	result, err := chargeUseCase.ChargeAllUser(ctx, fileRepo.UserChannel)
	if err != nil {
		log.Fatal("[Main]Cannot ChargeAllUser err=", err)
	}

	log.Println("done.")
	log.Println(" total received:\t 		THB  ", result.TotalDonation)
	log.Println(" successfully donated:\t	THB  ", result.SuccessfulDonation)
	log.Println(" faulty donation:\t 		THB  ", result.FaultyDonation)
	log.Println(" average per person:\t 	THB  ", result.AvgDonation)
	log.Print(" top donors:  ")
	for _, user := range result.TopDonors {
		log.Println("\t\t", user.Name)
	}

}
