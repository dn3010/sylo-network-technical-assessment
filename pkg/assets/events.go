package assets

import (
	"context"
	"errors"
	"strings"

	token "github.com/dn3010/sylo-network-technical-assessment/contract_bindings"
	"github.com/dn3010/sylo-network-technical-assessment/geth"
	"github.com/dn3010/sylo-network-technical-assessment/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func queryTransferEventsLogs(ctx context.Context, userAccount string, contracts []common.Address) (*Assets, error) {
	query := ethereum.FilterQuery{
		Addresses: contracts,
	}

	client := geth.From(ctx)
	if client == nil {
		return nil, errors.New("error no geth client was found")
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		return nil, err
	}

	logTransfers := make(map[string][]*LogTransfer, 0)
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {
		var transferEvent LogTransfer

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			//should the loop continue if decoding one of the events fails?
			if err != nil {
				return nil, err
			}

			//the third indexed value in the topic indicates to what account the
			//transfer was made, if it's not matches with the user account
			//passed-in via url, then skip it
			transferTo := common.HexToAddress(vLog.Topics[2].Hex())
			if transferTo.Hex() != userAccount {
				logger.From(ctx).Info().Msgf("Transfer event do not belong to account %s, skipping it", userAccount)
				continue
			}

			tokenAddress := common.HexToAddress(vLog.Topics[3].Hex())
			if tokenAddress.String() != "" {
				transferEvent.TokenId = tokenAddress.Hash().Big()
			}

			if _, ok := logTransfers[vLog.Address.Hex()]; ok {
				logTransfers[vLog.Address.Hex()] = append(logTransfers[vLog.Address.Hex()], &transferEvent)
			} else {
				logTransfers[vLog.Address.Hex()] = []*LogTransfer{&transferEvent}
			}
		}
	}
	return transformToAssets(logTransfers), nil
}

func transformToAssets(logTransfers map[string][]*LogTransfer) *Assets {
	tokens := make([]*Token, 0)
	for contract, transfers := range logTransfers {
		token := &Token{
			Token: contract,
		}

		tokenIds := make([]int64, 0)
		for _, transfer := range transfers {
			tokenIds = append(tokenIds, transfer.TokenId.Int64())
		}

		token.Ids = tokenIds
		tokens = append(tokens, token)
	}
	return &Assets{Tokens: tokens}
}
