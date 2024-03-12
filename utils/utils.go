package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var env_keys = []string{
	"ENV",
	"HTTP_PORT",
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASS",
	"DB_NAME",
}

type Response struct {
	Result  bool        `json:"result" example:"true"`
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty" `
}

func CheckEnvReady() {

	env := os.Getenv("ENV")

	if env == "LOCAL" || env == "" {
		err := godotenv.Load("local.env")
		if err != nil {
			fmt.Println("Not found environment file")
			panic("Not found environment file")
		}
	}

	var missEnv = []string{}

	for _, e := range env_keys {
		value, ok := os.LookupEnv(e)
		if !ok {
			missEnv = append(missEnv, e)
		} else {
			if value == "" {
				missEnv = append(missEnv, e)
			}
		}
	}

	if len(missEnv) != 0 {
		fmt.Print("Environment missing keys: ")
		fmt.Print(strings.Join(missEnv, ", "))
		fmt.Println("")
		panic("Env missing")
	}

	fmt.Println("Environment setup success!")
}

var CommonRespCode = map[string]Response{
	"OK": {
		Code:    2000,
		Message: "OK",
	},
	"VALIDATE_ERROR": {
		Code:    4001,
		Message: "Field is required",
	},
	"UNAUTHORIZED": {
		Code:    4002,
		Message: "Unauthorized, token is required",
	},
	"TOKEN_EXPIRED": {
		Code:    4003,
		Message: "Token has expired",
	},
	"INVALID_TOKEN": {
		Code:    4004,
		Message: "Invalid token",
	},
	"INTERNAL_SERVER_ERROR": {
		Code:    9000,
		Message: "Internal Server Error",
	},
}

func ContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func GenerateRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result += string(charset[randomIndex.Int64()])
	}
	return result, nil
}
