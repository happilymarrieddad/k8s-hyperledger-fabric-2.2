package main

// ResourceType resource
type ResourceType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	PrivateName string `json:"privateName"`
}

// ResourceTypePrivateData
type ResourceTypePrivateData struct {
	ResourceTypeID string `json:"resource_type_id"`
	DisplayName    string `json:"display_name"`
	Active         bool   `json:"active"`
}

// ResourceTypeIndex
type ResourceTypeIndex struct {
	ResourceType
	PrivateData ResourceTypePrivateData `json:"private_data"`
}

// ResourceTypeTransactionItem transaction item
type ResourceTypeTransactionItem struct {
	TXID         string            `json:"tx_id"`
	ResourceType ResourceTypeIndex `json:"resource_type"`
	Timestamp    int64             `json:"timestamp"`
}
