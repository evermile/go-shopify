package goshopify

import (
	"fmt"
)

// FulfillmentOrderService is an interface for interfacing with the fulfillment order endpoints
// of the Shopify API.
// https://shopify.dev/api/admin-rest/2022-10/resources/fulfillmentorder
type FulfillmentOrderService interface {
	GetFulfillmentOrders(orderID int64, options any) ([]FulfillmentOrder, error)
	MoveFulfillmentOrder(fulfillmentOrderID int64, locationID int64) (MoveFulfillmentOrderResource, error)
}

// FulfillmentOrderServiceOp handles communication with the fulfillment order
// related methods of the Shopify API.
type FulfillmentOrderServiceOp struct {
	client     *Client
	resource   string
	resourceID int64
}

// FulfillmentOrder represents a Shopify fulfillment order.
type FulfillmentOrder struct {
	ID                       int64      `json:"id,omitempty"`
	ShopID                   int64      `json:"shop_id"`
	OrderID                  int64      `json:"order_id,omitempty"`
	AssignedLocationID       int64      `json:"assigned_location_id"`
	LineItems                []LineItem `json:"line_items"`
	Status                   string     `json:"status,omitempty"`
	FulfillmentServiceHandle string     `json:"fulfillment_service_handle,omitempty"`
	RequestStatus            string     `json:"request_status,omitempty"`
}

// FulfillmentOrderResource represents the result from the orders/x/fulfillment_orders.json endpoint
type FulfillmentOrderResource struct {
	FulfillmentOrder []FulfillmentOrder `json:"fulfillment_orders"`
}

// MoveLocationRequest represents the request for the move inventory call
type MoveLocationRequest struct {
	FulfillmentOrderNewLocation FulfillmentOrderNewLocation `json:"fulfillment_order"`
}

type FulfillmentOrderNewLocation struct {
	NewLocationID int64 `json:"new_location_id"`
}

type MoveFulfillmentOrderResource struct {
	OriginalFulfillmentOrder  FulfillmentOrder `json:"original_fulfillment_order"`
	MovedFulfillmentOrder     FulfillmentOrder `json:"moved_fulfillment_order"`
	RemainingFulfillmentOrder FulfillmentOrder `json:"remaining_fulfillment_order"`
}

// GetFulfillmentOrders retrieves fulfillment orders for an order given its order id
func (s *FulfillmentOrderServiceOp) GetFulfillmentOrders(orderID int64, options interface{}) ([]FulfillmentOrder, error) {
	prefix := OrderPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d/fulfillment_orders.json", prefix, orderID)
	resource := new(FulfillmentOrderResource)
	err := s.client.Get(path, resource, options)
	return resource.FulfillmentOrder, err
}

// MoveFulfillmentOrder moves the fulfillment location to the new inventory location. Fails if items are not available
func (s *FulfillmentOrderServiceOp) MoveFulfillmentOrder(fulfillmentOrderID, locationID int64) (MoveFulfillmentOrderResource, error) {
	prefix := FulfillmentOrderPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d/move.json", prefix, fulfillmentOrderID)
	resource := new(MoveFulfillmentOrderResource)

	m := MoveLocationRequest{
		FulfillmentOrderNewLocation: FulfillmentOrderNewLocation{
			NewLocationID: locationID,
		},
	}
	err := s.client.Post(path, m, resource)
	return *resource, err
}
