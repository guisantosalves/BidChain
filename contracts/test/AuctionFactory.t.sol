// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/AuctionFactory.sol";
import "../src/Auction.sol";

contract AuctionFactoryTest is Test {
    AuctionFactory factory;
    address seller = makeAddr("seller");
    address buyer = makeAddr("buyer");

    event AuctionCreated(
        address indexed auctionAddress,
        address indexed seller,
        string description
    );

    function setUp() public {
        factory = new AuctionFactory();
        vm.deal(buyer, 10 ether);
    }

    // ─── createAuction retorna um endereço válido ─────────────────────────────

    function test_CreateAuctionReturnsValidAddress() public {
        vm.prank(seller);
        address auctionAddr = factory.createAuction(
            "Guitarra",
            "QmHash",
            1 days
        );

        // endereço válido = não é o endereço zero
        assertTrue(auctionAddr != address(0));
    }

    // ─── leilão criado tem o seller correto ───────────────────────────────────

    function test_CreatedAuctionHasCorrectSeller() public {
        vm.prank(seller);
        address auctionAddr = factory.createAuction(
            "Guitarra",
            "QmHash",
            1 days
        );

        Auction auction = Auction(auctionAddr);
        assertEq(auction.seller(), seller);
    }

    // ─── factory registra o leilão na lista ──────────────────────────────────

    function test_AuctionIsRegistered() public {
        vm.prank(seller);
        address auctionAddr = factory.createAuction(
            "Guitarra",
            "QmHash",
            1 days
        );

        assertEq(factory.getAuctionsCount(), 1);
        assertEq(factory.getAuctions()[0], auctionAddr);
    }

    // ─── múltiplos leilões são registrados corretamente ──────────────────────

    function test_MultipleAuctionsRegistered() public {
        vm.prank(seller);
        factory.createAuction("Guitarra", "QmHash1", 1 days);

        vm.prank(seller);
        factory.createAuction("Amplificador", "QmHash2", 2 days);

        assertEq(factory.getAuctionsCount(), 2);
    }

    // ─── evento AuctionCreated é emitido ─────────────────────────────────────

    function test_AuctionCreatedEventEmitted() public {
        vm.prank(seller);

        // não checamos o auctionAddress (topic1) porque não sabemos antes do deploy
        // checamos só o seller (topic2) e a description (data)
        vm.expectEmit(false, true, false, true);
        emit AuctionCreated(address(0), seller, "Guitarra");

        factory.createAuction("Guitarra", "QmHash", 1 days);
    }

    // ─── leilão criado pela factory aceita lances ─────────────────────────────

    function test_CreatedAuctionAcceptsBids() public {
        vm.prank(seller);
        address auctionAddr = factory.createAuction(
            "Guitarra",
            "QmHash",
            1 days
        );

        Auction auction = Auction(auctionAddr);

        vm.prank(buyer);
        auction.placeBid{value: 1 ether}();

        assertEq(auction.highestBidder(), buyer);
        assertEq(auction.highestBid(), 1 ether);
    }
}
