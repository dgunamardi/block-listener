package config

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
)

func InitializeSDK() {
	configProvider := config.FromFile(Vars.ChannelConfigPath)
	var opts []fabsdk.Option

	opts, err := getOptstoInitalizeSDK(Vars.ChannelConfigPath)
	if err != nil {
		panic(fmt.Errorf("Failed to create new SDK: %s\n", err))
	}

	Sdk, err = fabsdk.New(configProvider, opts...)
	if err != nil {
		panic(fmt.Errorf("Failed to create new SDK: %s\n", err))
	}
	fmt.Println("fabric SDK initialized")
}

func getOptstoInitalizeSDK(configPath string) ([]fabsdk.Option, error) {
	var opts []fabsdk.Option

	vc := viper.New()
	vc.SetConfigFile(configPath)
	err := vc.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to create new SDK: %s\n", err))
	}

	org := vc.GetString("client.originalOrganization")
	if org == "" {
		org = vc.GetString("client.organization")
	}

	opts = append(opts, fabsdk.WithOrgid(org))
	opts = append(opts, fabsdk.WithUserName(Vars.UserName))
	return opts, nil

}
