package v1

import "testing"

// 生成コードのディスクリプタが壊れていないか(init が panic しないか)の煙テスト
func TestDescriptorLoads(t *testing.T) {
	if OrderService_ServiceDesc.ServiceName == "" {
		t.Fatal("service descriptor is empty")
	}
}
