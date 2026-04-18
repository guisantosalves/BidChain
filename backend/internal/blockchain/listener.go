package blockchain

import (
	"context"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AuctionCreatedEvent struct {
	AuctionAddress common.Address
	Seller         common.Address
	Descrition     string
}

type Listener struct {
	client          *ethclient.Client
	contractABI     abi.ABI
	contractAddress common.Address
	Events          chan AuctionCreatedEvent
}

func NewListener(rpcUrl string) (*Listener, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(factoryABI))
	if err != nil {
		return nil, err
	}

	return &Listener{
		client:          client,
		contractABI:     parsedABI,
		contractAddress: common.HexToAddress(factoryAddress),
		Events:          make(chan AuctionCreatedEvent, 100),
	}, nil
}

func (l *Listener) handleLog(vLog types.Log) {
	event := AuctionCreatedEvent{}

	err := l.contractABI.UnpackIntoInterface(&event, "AuctionCreated", vLog.Data)
	if err != nil {
		log.Printf("listener: failed to unpack log: %v", err)
		return
	}

	event.AuctionAddress = common.HexToAddress(vLog.Topics[1].Hex())
	event.Seller = common.HexToAddress(vLog.Topics[2].Hex())

	l.Events <- event
}

func (l *Listener) Start(ctx context.Context) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{l.contractAddress},
	}

	logs := make(chan types.Log, 100)

	sub, err := l.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Printf("listener: failed to subscribe: %v", err)
		return
	}
	defer l.client.Close()

	log.Println("blockchain listener started")

	for {
		select {
		case <-ctx.Done():
			log.Println("BlockChaing listener stopped")
			return
		case err := <-sub.Err():
			log.Printf("listener: subscription error: %v", err)
			return
		case vLog := <-logs:
			l.handleLog(vLog)
		}
	}
}
