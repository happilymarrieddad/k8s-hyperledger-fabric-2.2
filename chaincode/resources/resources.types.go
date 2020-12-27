package main

// Resource resource
type Resource struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ResourceTypeID string `json:"resource_type_id"`
	Active         bool   `json:"active"`
	Owner          string `json:"owner"`
}

// ResourceTransactionItem the transaction item
type ResourceTransactionItem struct {
	TXID      string   `json:"tx_id"`
	Resource  Resource `json:"resource"`
	Timestamp int64    `json:"timestamp"`
}
