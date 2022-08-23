package testcom

import (
	"errors"
	"log"

	fdodeviceimplementation "github.com/WebauthnWorks/fdo-device-implementation"
	fdoshared "github.com/WebauthnWorks/fdo-shared"
)

type FDOTestState struct {
	_      struct{} `cbor:",toarray"`
	Passed bool     `json:"passed"`
	Error  string   `json:"error"`
}

func GenerateTestVoucherSet() ([]fdodeviceimplementation.VDANDV, error) {
	var resultVDANDV []fdodeviceimplementation.VDANDV
	for _, sgType := range fdoshared.DeviceSgTypeList {
		if sgType == fdoshared.StEPID10 || sgType == fdoshared.StEPID11 {
			log.Println("Generating test vouchers. EPID is not currently supported!")
			continue
		}

		vdanv, err := fdodeviceimplementation.NewVirtualDeviceAndVoucher(sgType)
		if err != nil {
			return resultVDANDV, errors.New("Error generating test VDANDV. " + err.Error())
		}

		resultVDANDV = append(resultVDANDV, *vdanv)
	}

	return resultVDANDV, nil
}
