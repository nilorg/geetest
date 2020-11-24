package main

import (
	"fmt"

	"github.com/nilorg/geetest/gt3"
)

func main() {
	id := "c9c4facd1a6feeb80802222cbb74ca8e"
	key := "f7475f921a41f7ba79ae15e41658627c"
	client := gt3.NewClient(id, key)
	var (
		registerResponse *gt3.RegisterResponse
		err              error
	)
	registerResponse, err = client.Register("md5")
	if err != nil {
		fmt.Printf("RegisterResponse Err: %v\n", err)
		return
	}
	fmt.Printf("registerResponse: %+v\n", registerResponse)
	var validateResponse *gt3.ValidateResponse
	validateResponse, err = client.Validate(registerResponse.Challenge)
	if err != nil {
		fmt.Printf("ValidateResponse Err: %v\n", err)
		return
	}
	fmt.Printf("validateResponse: %+v\n", validateResponse)
}
