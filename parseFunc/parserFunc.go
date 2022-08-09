package parseFunc

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// []Block.BlockData.Data, @>
// Envelope >
// 		Payload >
// 		 	Header >
// 		 	Transaction>
// 		 		[]TransactionAction, @ >
// 		 			Header
// 		 			ChaincodeActionPayload >
// 		 				ChaincodeProposalPayload >
// 		 				ChaincodeEndorsedAction >
// 		 				 	ProposalResponsePayload >
// 		 				 	 	ProposalHash
// 		 				 	 	ChainCodeAction >
// 		 				 	 	 	TxReadWriteSet >
// 		 				 	 	 	 	DataModel
// 		 				 	 	 	 	[]NsReadWriteSet, @ >
// 		 				 	 	 	 	 	Namespace
// 		 				 	 	 	 	 	KVRWset >
// 		 				 	 	 	 	 	 	[]KVRead, @ >
// 		 				 	 	 	 	 	 	 	Key
// 		 				 	 	 	 	 	 	 	Version
// 		 				 	 	 	 	 	 	[]RangeQueryInfo
// 		 				 	 	 	 	 	 	[]KVWrite, @ >
// 		 				 	 	 	 	 	 	 	Key
// 		 				 	 	 	 	 	 	 	IsDelete
// 		 				 	 	 	 	 	 	 	Value
// 		 				 	 	 	 	 	 	[]KVMetadataWrite
// 		 				 	 	 	 	 	[]CollectionHashReadWriteSet
// 		 				 	 	 	ChaincodeEvent
// 		 				 	 	 	Response
// 		 				 	 	 	ChaincodeId
// 		 				 	[]Endorsement
// 		Signature >

// 1 Block -> Multiple Transactions [Multiple Stakeholder Proposal]
// 1 Transaction -> Multiple Actions
// 1 Action -> Multiple NamespaceReadWriteSet, Multiple Events
// 1 NsRWSet -> 1 Asset Interaction

var (
	Block *common.Block
)

func ParseBlock(block *common.Block) {

	Block = block

	fmt.Printf("=== Block #%v ===\n", block.GetHeader().Number)

	txVCodes := GetTxValidationCodes(Block)
	for _, txVCode := range txVCodes {
		fmt.Printf("Validation Code: %v\n", txVCode)
	}

	payloads := GetPayloads(Block)
	for i, payload := range payloads {
		fmt.Printf("Transaction #%v\n", i)
		//fmt.Printf("Payload: \n%v\n", payload)

		// Get Channel Header
		channelHeader := &common.ChannelHeader{}
		err := proto.Unmarshal(payload.GetHeader().GetChannelHeader(), channelHeader)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal channel header: %v", err))
		}
		fmt.Printf("Channel Header: %v\n", channelHeader)

		if channelHeader.GetType() == int32(common.HeaderType_ENDORSER_TRANSACTION) {
			fmt.Println("is endorsed")
		}

		if channelHeader.GetType() == 1 {
			fmt.Println("is Configuration Block")
			return
		}

		// Get RWSet
		tx := &peer.Transaction{}
		err = proto.Unmarshal(payload.GetData(), tx)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal transaction: %v", err))
		}

		txActions := tx.GetActions()
		for _, txAction := range txActions {
			//fmt.Printf("Tx Actions:\n%v\n", txAction)

			ccActionPayload := &peer.ChaincodeActionPayload{}
			err := proto.Unmarshal(txAction.Payload, ccActionPayload)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal ccActionPayload: %v", err))
			}

			// --- Check for invocation spec ---
			// this is no longer necessary, config block filter handled in header
			ccProposalPayload := &peer.ChaincodeProposalPayload{}
			err = proto.Unmarshal(ccActionPayload.ChaincodeProposalPayload, ccProposalPayload)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal ccProposalPayload: %v", err))
			}

			ccInvocationSpec := &peer.ChaincodeInvocationSpec{}
			err = proto.Unmarshal(ccProposalPayload.Input, ccInvocationSpec)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal ccInvocationSpec: %v", err))
			}

			if ccInvocationSpec.ChaincodeSpec == nil {
				fmt.Println("is Genesis Block?")
				return
			}

			name := ccInvocationSpec.ChaincodeSpec.ChaincodeId.Name
			fmt.Printf("Chaincode Name: %v\n", name)

			// --- Check for invoc spec done ---

			proposalResponsePayload := &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(ccActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal proposalResponsePayload: %v", err))
			}

			ccAction := &peer.ChaincodeAction{}
			err = proto.Unmarshal(proposalResponsePayload.Extension, ccAction)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal chaincodeAction: %v", err))
			}

			//fmt.Printf("Chaincode Action:\n %v\n", ccAction)

			txRWSet := &rwset.TxReadWriteSet{}
			err = proto.Unmarshal(ccAction.Results, txRWSet)
			if err != nil {
				panic(fmt.Errorf("failed to unmarshal txRWSet: %v", err))
			}

			//fmt.Printf("TxRWSet:\n%v\n", txRWSet)

			for _, nsRWSet := range txRWSet.NsRwset {
				fmt.Printf("Namespace: %v\n", nsRWSet.Namespace)

				rwSet := &kvrwset.KVRWSet{}
				err := proto.Unmarshal(nsRWSet.Rwset, rwSet)
				if err != nil {
					panic(fmt.Errorf("failed to unmarshal rwset: %v", err))
				}

				//fmt.Printf("rwSet: %v\n", rwSet.String())

				fmt.Println("Reads")
				for _, read := range rwSet.Reads {
					fmt.Printf("key: %v\n", read.Key)
				}

				fmt.Println("Writes")
				for _, write := range rwSet.Writes {
					fmt.Printf("key: %v, value: %v\n", write.Key, string(write.Value))
				}
			}

		}
	}

}

func GetPayloads(block *common.Block) (payloads []*common.Payload) {
	dataArray := block.GetData().GetData()
	for _, data := range dataArray {
		envelope := &common.Envelope{}
		err := proto.Unmarshal(data, envelope)
		if err != nil {
			panic(fmt.Errorf("failed to create envelope: %v", err))
		}
		//fmt.Printf("Envelope: \n%v\n", envelope)

		payload := &common.Payload{}
		err = proto.Unmarshal(envelope.GetPayload(), payload)
		if err != nil {
			panic(fmt.Errorf("failed to create payload: %v", err))
		}
		payloads = append(payloads, payload)
	}
	return payloads
}

// Metadata format (?)
// [i][v]
// i -> Block metadata index, listed as enum:
// 		{SIGNATURES=0, LAST_CONFIG, TRANSACTIONS_FILTER, ORDERER, COMMIT HASH}
// v -> value of int that also refers to enum of TxValidationCode
// 		{VALID=0, NIL ENVELOPE, BAD PAYLOAD, ...}
func GetTxValidationCodes(block *common.Block) (txVCodes []*peer.TxValidationCode) {
	metadataArray := block.GetMetadata().GetMetadata()[common.BlockMetadataIndex_TRANSACTIONS_FILTER]
	for _, metadata := range metadataArray {
		metadataInt := int32(metadata)
		txVCodes = append(txVCodes, (*peer.TxValidationCode)(&metadataInt))
	}
	return txVCodes
}
