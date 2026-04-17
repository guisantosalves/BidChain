// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./Auction.sol";

contract AuctionFactory {
    // guarda o endereço de cada leilao criado pois o auction tmb é um contrato
    address[] public auctions;

    event AuctionCreated(
        address indexed auctionAddress,
        address indexed seller,
        string description
    );

    // storage  │ blockchain permanente       │ variáveis de estado 
    // memory   │ memória temporária          │ parâmetros e variáveis locais
    // calldata │ dados da chamada, read-only │ parâmetros de funções external
    function createAuction(
        string memory _description,
        string memory _ipfsHash,
        uint _durationSeconds
    ) external returns (address) {
        // aqui no create action o msg.sender é o address do vendedor que criou o leilao
        Auction auction = new Auction(
            msg.sender,
            _description,
            _ipfsHash,
            _durationSeconds
        );

        auctions.push(address(auction));

        emit AuctionCreated(address(auction), msg.sender, _description);

        return address(auction);
    }

    function getAuctions() external view returns (address[] memory) {
        return auctions;
    }

    function getAuctionsCount() external view returns (uint) {
        return auctions.length;
    }
}
