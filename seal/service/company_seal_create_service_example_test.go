package service_test

import (
	"fmt"

	"github.com/bytodream/official-seal-open/seal"
)

// ExampleCompanySealCreateService 演示服务构造与业务类型支持判断。
func ExampleCompanySealCreateService() {
	svc := seal.NewCompanySealCreateService()
	fmt.Println(svc.IsSupport(seal.SealTypeCompany, seal.BusinessTypeElectronic))
	// Output:
	// true
}
