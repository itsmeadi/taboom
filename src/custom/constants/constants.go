package constants

import "time"

const (
	TimeOut = 3 * time.Second

	RateLimitTimeTokenAPI  = 200 * time.Millisecond
	RateLimitTokenAPI      = 1
	RateLimitTimeChargeAPI = 200 * time.Millisecond
	RateLimitChargeAPI     = 1

	BKCurrency = "THB"

	OmisePublicKey = "pkey_test_5hdm5ihi576sm7vel42"
	OmiseSecretKey = "skey_test_5hdm5ihif1lxufd2az0"

	RetryLimit    = 5
	RetryWaitTime = 1000 * time.Millisecond //Time wait before retrying user

	FileReadBufferSize      = 512
	FileReadUserChannelSize = 10
)
