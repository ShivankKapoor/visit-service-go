package service

import (
	"os"
	"strconv"
)

func IsProd() bool {
	isProdStr := os.Getenv("PROD")
	isProd, err := strconv.ParseBool(isProdStr)
	if err != nil {
		isProd = true
	}
	return isProd
}
