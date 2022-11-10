package testexec

import (
	fdodeviceimplementation "github.com/WebauthnWorks/fdo-device-implementation"
	"github.com/WebauthnWorks/fdo-device-implementation/to2"
	fdoshared "github.com/WebauthnWorks/fdo-shared"
	"github.com/WebauthnWorks/fdo-shared/testcom"
	"github.com/WebauthnWorks/fdo-shared/testcom/dbs"
	reqtestsdeps "github.com/WebauthnWorks/fdo-shared/testcom/request"
)

func executeTo2_60(reqte reqtestsdeps.RequestTestInst, reqtDB *dbs.RequestTestDB) {
	for _, fdoTestId := range testcom.FIDO_TEST_LIST_DOT_60 {
		testCred, err := reqte.TestVouchers.GetVoucher(testcom.NULL_TEST)
		if err != nil {
			errTestState := testcom.NewFailTestState(fdoTestId, "Error getting voucher for TO2 60. "+err.Error())

			reqtDB.ReportTest(reqte.Uuid, testcom.NULL_TEST, errTestState)
			return
		}

		// Generating TO0 handler
		to2requestor := to2.NewTo2Requestor(fdodeviceimplementation.SRVEntry{
			SrvURL: reqte.URL,
		}, testCred.WawDeviceCredential, fdoshared.KEX_ECDH256, fdoshared.CIPHER_A128GCM) // TODO

		switch fdoTestId {
		case testcom.FIDO_DOT_60_POSITIVE:
			var errTestState testcom.FDOTestState
			_, _, err := to2requestor.HelloDevice60(fdoTestId)
			if err != nil {
				errTestState := testcom.NewFailTestState(fdoTestId, err.Error())

				reqtDB.ReportTest(reqte.Uuid, fdoTestId, errTestState)
				return
			} else {
				errTestState = testcom.NewSuccessTestState(fdoTestId)
				reqtDB.ReportTest(reqte.Uuid, fdoTestId, errTestState)
			}

		default:
			_, rvtTestState, err := to2requestor.HelloDevice60(fdoTestId)
			if rvtTestState == nil && err != nil {
				errTestState := testcom.NewFailTestState(fdoTestId, err.Error())
				rvtTestState = &errTestState
			}

			reqtDB.ReportTest(reqte.Uuid, fdoTestId, *rvtTestState)
		}
	}
}