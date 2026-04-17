package blockchain

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Listener struct {
	client *ethclient.Client
}

func NewListener(rpcUrl string) (*Listener, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	return &Listener{client: client}, nil
}

func (l *Listener) Start(ctx context.Context) {
	log.Println("Blockchain listener started")

	<-ctx.Done() // return a channel and it will be closed when the context is cancealed

	log.Println("blockchain listener stopped")
	l.client.Close()
}
