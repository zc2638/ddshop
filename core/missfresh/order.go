package missfresh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OrderResult struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		OrderNo string `json:"orderNo"`
		Id      int    `json:"id"`
	} `json:"data"`
}

func (s *Session) CreateOrder(ctx context.Context) error {
	u, err := url.Parse(HOST + ORDER_API)
	if err != nil {
		return fmt.Errorf("order url parse failed: %v", err)
	}

	urlPath := u.String()

	req := s.client.R()
	resp, err := s.execute(ctx, req, http.MethodPost, urlPath, maxRetryCount)
	if err != nil {
		return err
	}

	var cartResult CartResult
	if err := json.Unmarshal(resp.Body(), &cartResult); err != nil {
		return fmt.Errorf("parse response failed: %v", err)
	}

	return nil
}
