// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Auction {
    // === Variáveis de estado ===
    // Cada uma dessas fica gravada permanentemente na blockchain.
    // Isso custa gas pra escrever, mas é de graça pra ler.

    address public seller; // quem criou o leilao
    address public highestBidder;
    uint public highestBid;
    uint public deadline;
    bool public finalized;

    string public itemDescription;
    string public ipfsHash;

    // === Events ===
    // Logs stored on the blockchain.
    // Almost no gas cost and it's how the backend
    // knows what happened in the contract.
    event NewBid(address indexed bidder, uint amount);
    event AuctionFinalized(address indexed winner, uint amount);

    // === Constructor ===
    constructor(
        address _seller,
        string memory _description,
        string memory _ipfsHash,
        uint _durationSeconds
    ) {
        seller = _seller; // quem chamou o createAuction na factory
        itemDescription = _description; //
        ipfsHash = _ipfsHash;
        deadline = block.timestamp + _durationSeconds;
    }

    // === Place bid ===
    // payable = this function accepts ETH along with the call.
    // msg.value is how much ETH was sent
    function placeBid() external payable {
        require(block.timestamp < deadline, "Auction Ended");
        require(msg.value > highestBid, "Bid too low");

        // refund the previous highestBidder -> overwrite with the new one
        if (highestBidder != address(0)) {
            // != null
            (bool success, ) = payable(highestBidder).call{value: highestBid}(
                ""
            );
            require(success, "Refund failed");
        }

        highestBidder = msg.sender;
        highestBid = msg.value;

        emit NewBid(msg.sender, msg.value);
    }

    function finalize() external {
        require(block.timestamp >= deadline, "Auction still active");
        require(!finalized, "Already finalized");

        finalized = true;
        // it will send the money to the seller, because the current highestBidder won the item
        (bool success, ) = payable(seller).call{value: highestBid}("");
        require(success, "Transfer Failed");

        emit AuctionFinalized(highestBidder, highestBid);
    }
}
