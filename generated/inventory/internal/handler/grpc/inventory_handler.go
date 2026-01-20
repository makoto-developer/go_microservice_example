package grpc

// InventoryServiceHandler is a mock implementation for Inventory Service
// Note: Full proto implementation is temporarily disabled due to namespace conflict
type InventoryServiceHandler struct {
}

func NewInventoryServiceHandler() *InventoryServiceHandler {
	return &InventoryServiceHandler{}
}
