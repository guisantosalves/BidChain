// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Auction.sol";

contract AuctionTest is Test {
    Auction auction;
    address seller = makeAddr("seller");
    address alice = makeAddr("alice");
    address bob = makeAddr("bob");

    event NewBid(address indexed bidder, uint amount);

    // setup roda antes de cada teste
    function setUp() public {
        // vm.prank -> prox chamada tem msg.sender = seller
        vm.prank(seller);
        auction = new Auction(seller, "Guitarra fender", "QmFakeHash", 1 days);

        // eth de mentira
        vm.deal(alice, 10 ether);
        vm.deal(bob, 10 ether);
    }

    // lance inicial
    function test_FirstBidUpdatesState() public {
        vm.prank(alice);
        auction.placeBid{value: 1 ether}();

        assertEq(auction.highestBidder(), alice);
        assertEq(auction.highestBid(), 1 ether);
    }

    // lance superado
    function test_HigherBidFundsPrevious() public {
        vm.prank(alice);
        auction.placeBid{value: 1 ether}();

        uint aliceBalanceBefore = alice.balance;

        vm.prank(bob);
        auction.placeBid{value: 2 ether}();

        assertEq(auction.highestBidder(), bob);
        assertEq(auction.highestBid(), 2 ether);

        // make sure she has receveid it back
        assertEq(alice.balance, aliceBalanceBefore + 1 ether);
    }

    // lance baixo
    function test_LowBidReverts() public {
        vm.prank(alice);
        auction.placeBid{value: 1 ether}();

        vm.expectRevert("Bid too low");

        vm.prank(bob);
        auction.placeBid{value: 0.5 ether}();
    }

    // lance apos deadline
    function test_BidAfterDeadLineReverts() public {
        vm.warp(block.timestamp + 2 days);

        vm.expectRevert("Auction Ended");

        vm.prank(alice);
        auction.placeBid{value: 1 ether}();
    }

    // finalizar antes do prazo
    function test_FinalizeBeforeDeadlineReverts() public {
        vm.prank(alice);
        auction.placeBid{value: 1 ether}();

        vm.expectRevert("Auction still active");
        auction.finalize();
    }

    // finalizar corretamente
    function test_FinalizeTransfersFundsToSeller() public {
        vm.prank(alice);
        auction.placeBid{value: 3 ether}();

        uint sellerBalanceBefore = seller.balance; // sem os 3

        vm.warp(block.timestamp + 2 days);
        auction.finalize();

        assertEq(auction.finalized(), true);
        // seller recebeu os 3 ether
        assertEq(seller.balance, sellerBalanceBefore + 3 ether); // deve conter os 3 apos o finalize
    }

    // finalizar duas vezes
    function test_DoubleFinalize_Reverts() public {
        vm.prank(alice);
        auction.placeBid{value: 1 ether}();

        vm.warp(block.timestamp + 2 days);
        auction.finalize();

        vm.expectRevert("Already finalized");
        auction.finalize();
    }

    // evento emitido
    function test_NewBidEventEmitted() public {
        // expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData)
        // true, false, false, true → verifica indexed bidder + amount
        vm.expectEmit(true, false, false, true);
        emit NewBid(alice, 1 ether);

        vm.prank(alice);
        auction.placeBid{value: 1 ether}();
    }
}
