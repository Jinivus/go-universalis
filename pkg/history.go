package universalis

import (
	"context"
	"fmt"
)

type HistoryService service

type HistoryResult struct {
	ItemID         int   `json:"itemID"`
	WorldID        int   `json:"worldID"`
	LastUploadTime int64 `json:"lastUploadTime"`
	Entries        []struct {
		Hq           bool `json:"hq"`
		PricePerUnit int  `json:"pricePerUnit"`
		Quantity     int  `json:"quantity"`
		Timestamp    int  `json:"timestamp"`
	} `json:"entries"`
	StackSizeHistogram struct {
		Num1 int `json:"1"`
	} `json:"stackSizeHistogram"`
	StackSizeHistogramNQ struct {
		Num1 int `json:"1"`
	} `json:"stackSizeHistogramNQ"`
	StackSizeHistogramHQ struct {
		Num1 int `json:"1"`
	} `json:"stackSizeHistogramHQ"`
	RegularSaleVelocity float64 `json:"regularSaleVelocity"`
	NqSaleVelocity      float64 `json:"nqSaleVelocity"`
	HqSaleVelocity      float64 `json:"hqSaleVelocity"`
	WorldName           string  `json:"worldName"`
}

func (s *HistoryService) History(ctx context.Context, world string, query string) (*HistoryResult, *Response, error) {
	result := new(HistoryResult)
	resp, err := s.getHistory(ctx, world, query, false, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

func (s *HistoryService) HistoryWithOptions(ctx context.Context, world string, query string, options ListingOptions) (*HistoryResult, *Response, error) {
	result := new(HistoryResult)
	resp, err := s.getHistory(ctx, world, query, options.HQOnly, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

func (s *HistoryService) getHistory(ctx context.Context, world string, query string, hq bool, result interface{}) (*Response, error) {

	u := fmt.Sprintf("history/%s/%s", world, query)
	if hq {
		u = fmt.Sprintf("%s%s", u, "?hq=1")
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, result)
}
