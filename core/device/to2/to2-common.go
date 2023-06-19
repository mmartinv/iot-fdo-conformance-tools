package to2

import (
	"github.com/fido-alliance/fdo-fido-conformance-server/core/device/common"
	fdoshared "github.com/fido-alliance/fdo-fido-conformance-server/core/shared"
	"github.com/fido-alliance/fdo-fido-conformance-server/core/shared/testcom"
)

var MaxDeviceMessageSize uint16 = 2048
var MaxOwnerServiceInfoSize uint16 = 2048

type To2Requestor struct {
	SrvEntry        common.SRVEntry
	Credential      fdoshared.WawDeviceCredential
	KexSuiteName    fdoshared.KexSuiteName
	CipherSuiteName fdoshared.CipherSuiteName

	AuthzHeader string
	SessionKey  fdoshared.SessionKeyInfo
	XAKex       []byte
	XBKEXParams fdoshared.KeXParams

	NonceTO2ProveOV60 fdoshared.FdoNonce
	NonceTO2ProveDv61 fdoshared.FdoNonce
	NonceTO2SetupDv64 fdoshared.FdoNonce

	ProveOVHdr61PubKey fdoshared.FdoPublicKey
	OvHmac             fdoshared.HashOrHmac

	Completed60 bool
	Completed62 bool
	Completed64 bool
}

func NewTo2Requestor(srvEntry common.SRVEntry, credential fdoshared.WawDeviceCredential, kexSuitName fdoshared.KexSuiteName, cipherSuitName fdoshared.CipherSuiteName) To2Requestor {
	return To2Requestor{
		SrvEntry:        srvEntry,
		Credential:      credential,
		KexSuiteName:    kexSuitName,
		CipherSuiteName: cipherSuitName,
	}
}

func (h *To2Requestor) confCheckResponse(bodyBytes []byte, fdoTestID testcom.FDOTestID, httpStatusCode int) testcom.FDOTestState {
	switch fdoTestID {
	case testcom.ExpectGroupTests(testcom.FIDO_TEST_LIST_DOT_60, fdoTestID):
		return testcom.ExpectAnyFdoError(bodyBytes, fdoTestID, fdoshared.MESSAGE_BODY_ERROR, httpStatusCode)

	case testcom.ExpectGroupTests(testcom.FIDO_TEST_LIST_DOT_62, fdoTestID):
		return testcom.ExpectAnyFdoError(bodyBytes, fdoTestID, fdoshared.MESSAGE_BODY_ERROR, httpStatusCode)

	case testcom.ExpectGroupTests(testcom.FIDO_TEST_LIST_DOT_64, fdoTestID):
		return testcom.ExpectAnyFdoError(bodyBytes, fdoTestID, fdoshared.MESSAGE_BODY_ERROR, httpStatusCode)

	case testcom.ExpectGroupTests(testcom.FIDO_TEST_LIST_DOT_66, fdoTestID):
		return testcom.ExpectAnyFdoError(bodyBytes, fdoTestID, fdoshared.MESSAGE_BODY_ERROR, httpStatusCode)

	case testcom.ExpectGroupTests(testcom.FIDO_TEST_LIST_DOT_68, fdoTestID):
		return testcom.ExpectAnyFdoError(bodyBytes, fdoTestID, fdoshared.MESSAGE_BODY_ERROR, httpStatusCode)

	case testcom.ExpectGroupTests(testcom.FIDO_TEST_LIST_DOT_70, fdoTestID):
		return testcom.ExpectAnyFdoError(bodyBytes, fdoTestID, fdoshared.MESSAGE_BODY_ERROR, httpStatusCode)

	}
	return testcom.FDOTestState{
		Passed: false,
		Error:  "Unsupported test " + string(fdoTestID),
	}
}