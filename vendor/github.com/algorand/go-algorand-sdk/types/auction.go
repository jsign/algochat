package types

// Bid represents a bid by a user as part of an auction.
type Bid struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// BidderKey identifies the bidder placing this bid.
	BidderKey Address `codec:"bidder"`

	// BidCurrency specifies how much external currency the bidder
	// is putting in with this bid.
	BidCurrency uint64 `codec:"cur"`

	// MaxPrice specifies the maximum price, in units of external
	// currency per Algo, that the bidder is willing to pay.
	// This must be at least as high as the current price of the
	// auction in the block in which this bid appears.
	MaxPrice uint64 `codec:"price"`

	// BidID identifies this bid.  The first bid by a bidder (identified
	// by BidderKey) with a particular BidID on the blockchain will be
	// considered, preventing replay of bids.  Specifying a different
	// BidID allows the bidder to place multiple bids in an auction.
	BidID uint64 `codec:"id"`

	// AuctionKey specifies the auction for this bid.
	AuctionKey Address `codec:"auc"`

	// AuctionID identifies the auction for which this bid is intended.
	AuctionID uint64 `codec:"aid"`
}
