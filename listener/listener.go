package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	cfg "earhart.com/config"
	//parseFunc "earhart.com/parseFunc"
	parser "earhart.com/parser"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"

	eventClient "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
)

var (
	parsedBlock parser.Block

	seekType   = seek.FromBlock
	startBlock = 340
)

func main() {
	cfg.LoadConfig()
	cfg.InitializeSDK()
	cfg.InitializeUserIdentity()

	session := cfg.Sdk.Context(fabsdk.WithIdentity(cfg.User))

	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(session, cfg.ChannelId)
	}

	args := os.Args[1:]
	SetArgs(args)

	ListenToBlockEvents(channelProvider, seek.Type(seekType), uint64(startBlock))
}

func SetArgs(args []string) {
	if len(args) == 0 || len(args) > 2 {
		seekType = seek.Newest
		return
	}
	switch args[0] {
	case "oldest":
		seekType = seek.Oldest
	case "newest":
		seekType = seek.Newest
	case "from":
		if args[1] == "" {
			panic("not enough arguments. 'from' should be followed by a number indicating the starting block\n")
		}
		seekType = seek.FromBlock

		sb, err := strconv.Atoi(args[1])
		if err != nil {
			panic(fmt.Errorf("error in arg to int conversion: %v", err))
		}
		startBlock = sb
	default:
		seekType = seek.Newest
	}
}

func ListenToBlockEvents(channelProvider context.ChannelProvider, seekType seek.Type, startBlock uint64) {
	client, err := eventClient.New(
		channelProvider,
		eventClient.WithBlockEvents(),
		eventClient.WithSeekType(seekType),
		eventClient.WithBlockNum(startBlock),
	)

	if err != nil {
		panic(fmt.Errorf("failed to create event client: %v", err))
	}

	eventRegister, blockEvents, err := client.RegisterBlockEvent()
	defer client.Unregister(eventRegister)

	fmt.Println("... start listening to events ...")

	cwd, _ := os.Getwd()

	// Since blockEvents is a channel of [][]bytes
	// this for receives values indefinitely or until the channel is closed (by the sender)
	//
	for event := range blockEvents {
		//parseFunc.ParseBlock(event.Block)
		parsedBlock.Init(event.Block)

		//isParse := parsedBlock.BlockData.Envelopes[0].IsTransaction

		//fmt.Println(isParse)
		//
		fileName := "blockEvent" + strconv.Itoa(int(parsedBlock.BlockProto.Header.Number)) + ".json"

		bInfo, _ := json.Marshal(parsedBlock)
		err := ioutil.WriteFile(cwd+"/block-event-parses/"+fileName, bInfo, 0644)
		if err != nil {
			panic(err)
		}

		bInfoIn, _ := json.MarshalIndent(parsedBlock, "", " ")
		err = ioutil.WriteFile(cwd+"/block-event-parses/in"+fileName, bInfoIn, 0644)
		if err != nil {
			panic(err)
		}
	}

}
