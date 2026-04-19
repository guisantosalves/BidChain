package blockchain

// it will be used to decode the bytes coming from blockchain
const factoryABI = `[
	{
      "type": "function",
      "name": "createAuction",
      "inputs": [
        {"name": "_description", "type": "string"},
        {"name": "_ipfsHash", "type": "string"},
        {"name": "_durationSeconds", "type": "uint256"}
      ],
      "outputs": [{"name": "", "type": "address"}],
      "stateMutability": "nonpayable"
    },
	{
		"type": "event",
		"name": "AuctionCreated",
		"inputs": [
			{"name": "auctionAddress", "type": "address", "indexed": true},
			{"name": "seller", "type": "address", "indexed": true},
        	{"name": "description", "type": "string", "indexed": false}
		]
	}
]`

// endereço do contrato na chain
const FactoryAddress = "0x47b30924Bf389f2489B66FB1B65D6281a7f534BC"
