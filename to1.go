package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/WebauthnWorks/fdo-device-implementation/fdoshared"
	"github.com/fxamacker/cbor/v2"
)

type RVEntry struct {
	RVURL       string
	AccessToken string
}

type To1Requestor struct {
	rvEntry     RVEntry
	credential  fdoshared.WawDeviceCredential
	authzHeader string
}

type VoucherDBEntry struct {
	_              struct{} `cbor:",toarray"`
	Voucher        fdoshared.OwnershipVoucher
	PrivateKeyX509 []byte
}

func NewTo1Requestor(rvEntry RVEntry, credential fdoshared.WawDeviceCredential) To1Requestor {
	return To1Requestor{
		rvEntry:    rvEntry,
		credential: credential,
	}
}

func SendCborPost(rvEntry RVEntry, cmd fdoshared.FdoCmd, payload []byte, authzHeader *string) ([]byte, string, error) {
	url := rvEntry.RVURL + fdoshared.FDO_101_URL_BASE + cmd.ToString()

	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, "", errors.New("Error creating new request. " + err.Error())
	}

	if authzHeader != nil {
		req.Header.Set("Authorization", *authzHeader)
	}

	req.Header.Set("Content-Type", "application/cbor")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("Error sending post request to %s url. %s", url, err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("Error reading body bytes for %s url. %s", url, err.Error())
	}

	return bodyBytes, resp.Header.Get("Authorization"), nil
}

func (h *To1Requestor) HelloRV30() (fdoshared.HelloRVAck31, error) {
	var helloRVAck31 fdoshared.HelloRVAck31

	helloRV30Bytes, err := cbor.Marshal(fdoshared.HelloRV30{
		Guid:      h.credential.DCGuid,
		EASigInfo: h.credential.DCSigInfo,
	})

	if err != nil {
		return helloRVAck31, errors.New("HelloRV30: Error marshaling HelloRV30. " + err.Error())
	}

	resultBytes, authzHeader, err := SendCborPost(h.rvEntry, fdoshared.TO1_HELLO_RV_30, helloRV30Bytes, &h.rvEntry.AccessToken)
	if err != nil {
		return helloRVAck31, errors.New("Hello30: " + err.Error())
	}

	h.authzHeader = authzHeader

	err = cbor.Unmarshal(resultBytes, &helloRVAck31)
	if err != nil {
		return helloRVAck31, errors.New("HelloRV30: Failed to unmarshal HelloRVAck31. " + err.Error())
	}

	return helloRVAck31, nil
}

func (h *To1Requestor) ProveToRV32(helloRVAck31 fdoshared.HelloRVAck31) (*fdoshared.RVRedirect33, error) {

	var proveToRV32Payload fdoshared.EATPayloadBase = fdoshared.EATPayloadBase{
		EatNonce: helloRVAck31.NonceTO1Proof,
	}

	proveToRV32PayloadBytes, err := cbor.Marshal(proveToRV32Payload)
	if err != nil {
		return nil, errors.New("ProveToRV32: Error generating ProveToRV32. " + err.Error())
	}

	privateKeyInst, err := fdoshared.ExtractPrivateKey(h.credential.DCPrivateKeyDer)
	if err != nil {
		return nil, errors.New("ProveToRV32: Error extracting private key from voucher. " + err.Error())
	}

	sgType := helloRVAck31.EBSigInfo.SgType

	proveToRV32, err := fdoshared.GenerateCoseSignature(proveToRV32PayloadBytes, fdoshared.ProtectedHeader{}, fdoshared.UnprotectedHeader{}, privateKeyInst, sgType)
	if err != nil {
		return nil, errors.New("ProveToRV32: Error generating ProveToRV32. " + err.Error())
	}

	proveToRV32Bytes, err := cbor.Marshal(proveToRV32)
	if err != nil {
		return nil, errors.New("ProveToRV32: Error marshaling proveToRV32. " + err.Error())
	}

	resultBytes, authzHeader, err := SendCborPost(h.rvEntry, fdoshared.TO1_PROVE_TO_RV_32, proveToRV32Bytes, &h.authzHeader)
	if err != nil {
		return nil, errors.New("Hello30: " + err.Error())
	}

	h.authzHeader = authzHeader

	var rvRedirect33 fdoshared.RVRedirect33

	err = cbor.Unmarshal(resultBytes, &rvRedirect33)

	if err != nil {
		return nil, errors.New("RVRedirect33: Failed to unmarshal RVRedirect33. " + err.Error())
	}

	return &rvRedirect33, nil
}
