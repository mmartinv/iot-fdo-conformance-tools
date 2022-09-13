package fdodeviceimplementation

import (
	"errors"

	"github.com/WebauthnWorks/fdo-fido-conformance-server/testcom"
	fdoshared "github.com/WebauthnWorks/fdo-shared"
	"github.com/fxamacker/cbor/v2"
)

func (h *To2Requestor) DeviceServiceInfoReady66(fdoTestID testcom.FDOTestID) (*fdoshared.OwnerServiceInfoReady67, *testcom.FDOTestState, error) {
	var testState testcom.FDOTestState

	deviceSrvInfoReady := fdoshared.DeviceServiceInfoReady66{
		ReplacementHMac:       &h.OvHmac,
		MaxOwnerServiceInfoSz: &MaxOwnerServiceInfoSize,
	}
	deviceSrvInfoReadyBytes, _ := cbor.Marshal(deviceSrvInfoReady)

	deviceSrvInfoReadyBytesEnc, err := fdoshared.AddEncryptionWrapping(deviceSrvInfoReadyBytes, h.SessionKey, h.CipherSuiteName)
	if err != nil {
		return nil, nil, errors.New("DeviceServiceInfoReady66: Error encrypting... " + err.Error())
	}

	rawResultBytes, authzHeader, httpStatusCode, err := SendCborPost(h.SrvEntry, fdoshared.TO2_66_DEVICE_SERVICE_INFO_READY, deviceSrvInfoReadyBytesEnc, &h.AuthzHeader)
	if fdoTestID != testcom.NULL_TEST {
		testState = h.confCheckResponse(rawResultBytes, fdoTestID, httpStatusCode)
	}

	if err != nil {
		return nil, nil, err
	}

	h.AuthzHeader = authzHeader

	bodyBytes, err := fdoshared.RemoveEncryptionWrapping(rawResultBytes, h.SessionKey, h.CipherSuiteName)
	if err != nil {
		return nil, nil, errors.New("DeviceServiceInfoReady66: Error decrypting... " + err.Error())
	}

	var ownerServiceInfoReady67 fdoshared.OwnerServiceInfoReady67
	err = cbor.Unmarshal(bodyBytes, &ownerServiceInfoReady67)
	if err != nil {
		return nil, nil, errors.New("DeviceServiceInfoReady66: Error decoding OwnerServiceInfoReady67... " + err.Error())
	}

	return &ownerServiceInfoReady67, &testState, nil
}
