package main

import (
	"ps-sdk/sdk/sungrow"
)

func main() {
	sdk, err := sungrow.NewSungrowSDK(
		sungrow.Credentials{
			AppID:        "EFCACCA7D56F186180F0AB1B708662D9",
			AppSecret:    "0dq39edbjub2rhq7hff9r7qqpeguxk5p",
			UserAccount:  "js",
			UserPassword: "fxdlyw12.",
		},
		sungrow.WithDevMode(true),
	)
	if err != nil {
		panic(err)
	}

	sdk.Login()
}
