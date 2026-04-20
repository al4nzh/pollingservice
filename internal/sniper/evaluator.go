
package sniper

// ListingEvaluation holds the result of evaluating a listing
// Discount is relative to reference price and second best price
// IsDeal is true if discount is better than threshold/second best
// Reason gives explanation for decision
type ListingEvaluation struct {
	Listing   Listing
	DiscountFromRef      float64
	DiscountFromSecond   float64
	IsDeal   bool
	Reason   string
}

func (e ListingEvaluation) ListingID() string {
	return e.Listing.ID
}

func (e ListingEvaluation) AppID() int {
	return e.Listing.Item.AppID
}

func (e ListingEvaluation) MarketHashName() string {
	return e.Listing.Item.MarketHashName
}

func (e ListingEvaluation) DisplayName() string {
	if e.Listing.Item.ItemName != "" && e.Listing.Item.WearName != "" {
		return e.Listing.Item.ItemName + " (" + e.Listing.Item.WearName + ")"
	}
	if e.Listing.Item.ItemName != "" {
		return e.Listing.Item.ItemName
	}
	return e.Listing.Item.MarketHashName
}

func (e ListingEvaluation) DealPrice() float64 {
	return float64(e.Listing.Price) / 100.0
}

func (e ListingEvaluation) DealDiscountFromRef() float64 {
	return e.DiscountFromRef
}

func (e ListingEvaluation) DealDiscountFromSecond() float64 {
	return e.DiscountFromSecond
}

func (e ListingEvaluation) DealReason() string {
	return e.Reason
}

// EvaluateListings compares each listing to reference and second best price
// Returns evaluations for all listings
func EvaluateListings(listing Listing) []ListingEvaluation {
	if listing.ID == "" {
		return nil
	}

	// Get reference price from first listing (all should be for same item)
	Price := float64(listing.Price) / 100.0
	refPrice := float64(listing.Reference.Price) / 100.0
	secondBestPrice := listin
	secondrefPrice := 
}
