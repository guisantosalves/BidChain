// this file will call the contract on chain
package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Caller struct {
	client          *ethclient.Client
	contractABI     abi.ABI
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
}

func NewCaller(rpcURL, privateKeyHex, contractAddr string) (*Caller, error) {
	// abre uma conexão com a Sepolia — é como conectar no banco, mas na blockchain
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(factoryABI))
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return nil, err
	}

	return &Caller{
		client:          client,
		contractABI:     parsedABI,
		contractAddress: common.HexToAddress(contractAddr),
		privateKey:      privateKey,
	}, nil
}

func (c *Caller) CreateAuction(ctx context.Context, description, ipfsHash string, durationSeconds uint64) error {
	// get sepolia chain id
	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return err
	}

	// login into blockchain
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		return err
	}

	// instance of the contract
	contract := bind.NewBoundContract(c.contractAddress, c.contractABI, c.client, c.client, c.client)

	// do the transaction
	_, err = contract.Transact(auth, "createAuction", description, ipfsHash, new(big.Int).SetUint64(durationSeconds))

	return err
}
