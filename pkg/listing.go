package universalis

import (
	"context"
	"fmt"
)

type ListingService service

type Listing struct {
	LastReviewTime int           `json:"lastReviewTime"`
	PricePerUnit   int           `json:"pricePerUnit"`
	Quantity       int           `json:"quantity"`
	StainID        int           `json:"stainID"`
	CreatorName    string        `json:"creatorName"`
	CreatorID      string        `json:"creatorID"`
	HQ             bool          `json:"hq"`
	IsCrafted      bool          `json:"isCrafted"`
	ListingID      interface{}   `json:"listingID"`
	Materia        []interface{} `json:"materia"`
	OnMannequin    bool          `json:"onMannequin"`
	RetainerCity   int           `json:"retainerCity"`
	RetainerID     string        `json:"retainerID"`
	RetainerName   string        `json:"retainerName"`
	SellerID       string        `json:"sellerID"`
	Total          int           `json:"total"`
}

type ListingResult struct {
	ItemID         int        `json:"itemID"`
	WorldID        int        `json:"worldID"`
	LastUploadTime int64      `json:"lastUploadTime"`
	Listings       *[]Listing `json:"listings"`
	RecentHistory  []struct {
		Hq           bool   `json:"hq"`
		PricePerUnit int    `json:"pricePerUnit"`
		Quantity     int    `json:"quantity"`
		Timestamp    int    `json:"timestamp"`
		BuyerName    string `json:"buyerName"`
		Total        int    `json:"total"`
	} `json:"recentHistory"`
	CurrentAveragePrice   float64 `json:"currentAveragePrice"`
	CurrentAveragePriceNQ float64 `json:"currentAveragePriceNQ"`
	CurrentAveragePriceHQ float64 `json:"currentAveragePriceHQ"`
	RegularSaleVelocity   float64 `json:"regularSaleVelocity"`
	NqSaleVelocity        float64 `json:"nqSaleVelocity"`
	HqSaleVelocity        float64 `json:"hqSaleVelocity"`
	AveragePrice          float64 `json:"averagePrice"`
	AveragePriceNQ        float64 `json:"averagePriceNQ"`
	AveragePriceHQ        float64 `json:"averagePriceHQ"`
	MinPrice              int     `json:"minPrice"`
	MinPriceNQ            int     `json:"minPriceNQ"`
	MinPriceHQ            int     `json:"minPriceHQ"`
	MaxPrice              int     `json:"maxPrice"`
	MaxPriceNQ            int     `json:"maxPriceNQ"`
	MaxPriceHQ            int     `json:"maxPriceHQ"`
	StackSizeHistogram    struct {
		Num1 int `json:"1"`
	} `json:"stackSizeHistogram"`
	StackSizeHistogramNQ struct {
	} `json:"stackSizeHistogramNQ"`
	StackSizeHistogramHQ struct {
		Num1 int `json:"1"`
	} `json:"stackSizeHistogramHQ"`
	WorldName string `json:"worldName"`
}

type ListingOptions struct {
	HQOnly bool
}

func (s *ListingService) Listings(ctx context.Context, world string, query string) (*ListingResult, *Response, error) {
	result := new(ListingResult)
	resp, err := s.getListings(ctx, world, query, false, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

func (s *ListingService) ListingsWithOptions(ctx context.Context, world string, query string, options ListingOptions) (*ListingResult, *Response, error) {
	result := new(ListingResult)
	resp, err := s.getListings(ctx, world, query, options.HQOnly, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

func (s *ListingService) getListings(ctx context.Context, world string, query string, hq bool, result interface{}) (*Response, error) {

	u := fmt.Sprintf("%s/%s", world, query)
	if hq {
		u = fmt.Sprintf("%s%s", u, "?hq=1")
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, result)
}
