package v1

import "testing"

// 生成コードの descriptor が壊れていないことの恒久チェック
// (過去に pb.go の rawDesc 破損で init panic が起きたための再発防止)
func TestProtoInit(t *testing.T) {
	if ShipmentStatus_SHIPMENT_STATUS_DELIVERED != 5 {
		t.Error("unexpected enum value")
	}
	req := &CreateShipmentRequest{OrderId: "x"}
	if req.GetOrderId() != "x" {
		t.Error("getter failed")
	}
}
