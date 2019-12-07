package entities

type UserInfo struct {
	Name     string
	Amount   int64
	CCNumber []byte
	CVV      []byte//TODO
	ExpMonth int64
	ExpYear  int64
	//Retries  int

	//IsLastUser bool
}

type Transaction struct {
	User   UserInfo
	Amount int64
	Retry  int
	Err    error
}

type Keys struct {
	OmisePublicKey string
	OmiseSecretKey string
}

type Result struct {
	TotalDonation      int64
	SuccessfulDonation int64
	FaultyDonation     int64
	AvgDonation        int64
	TopDonors          []UserInfo
	Currency           string
}
