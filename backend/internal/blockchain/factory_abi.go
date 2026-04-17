package blockchain

// it will be used to decode the bytes coming from blockchain
const factoryABI = `[
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

const factoryAddress = "0x47b30924Bf389f2489B66FB1B65D6281a7f534BC"
