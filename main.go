package main

import (
	"encoding/hex"
	"errors"
	"log"
	"os"

	"github.com/WebauthnWorks/fdo-device-implementation/fdoshared"
	"github.com/fxamacker/cbor/v2"
	"github.com/urfave/cli/v2"
)

type HelloRV30 struct {
	_         struct{} `cbor:",toarray"`
	Guid      []byte
	EASigInfo fdoshared.SigInfo
}

func main() {
	cliapp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "genvoucher",
				Usage: "Generate new ownership voucher and DI. See _dis and _vouchers",
				Action: func(c *cli.Context) error {
					err := GenerateVoucher(fdoshared.StSECP256R1)

					if err != nil {
						return errors.New("Error generating voucher. " + err.Error())
					}
					return nil
				},
			},
			{
				Name:  "testdecode",
				Usage: "Test decoding",
				Action: func(c *cli.Context) error {
					sourceBytes, _ := hex.DecodeString("82508d62ddb18b404cf58cf5f22cef5c4576822678244920616d206120706f7461746f652120536d6172742c20496f542c20706f7461746f6521")

					var helloRv30 HelloRV30

					err := cbor.Unmarshal(sourceBytes, &helloRv30)

					if err != nil {
						log.Panic("Error decoding source: " + err.Error())
					}

					log.Println(helloRv30)
					return nil
				},
			},
			{
				Name:  "testto1",
				Usage: "",
				Action: func(c *cli.Context) error {
					vouchers, err := LoadLocalVouchers()
					if err != nil {
						log.Panic(err)
					}

					for _, voucher := range vouchers {
						to1requestor := NewTo1Requestor(RVEntry{
							RVURL:       "http://localhost:8083",
							AccessToken: "",
						}, voucher)

						helloRVAck31, err := to1requestor.HelloRV30()
						if err != nil {
							log.Panic(err)
						}

						// This below needs to be refactored ... ?
						// Need to include NonceTO1Proof in coseSignature --> generate with helloRVAck31.NonceTO1Proof
						var proveToRV32Payload fdoshared.EATPayloadBase = fdoshared.EATPayloadBase{
							EatNonce: helloRVAck31.NonceTO1Proof,
						}
						proveToRV32PayloadBytes, err := cbor.Marshal(proveToRV32Payload)

						privateKeyInst, err := fdoshared.ExtractPrivateKey(voucher.PrivateKeyX509)
						proveToRV32, err := fdoshared.GenerateCoseSignature(proveToRV32PayloadBytes, fdoshared.ProtectedHeader{}, fdoshared.UnprotectedHeader{}, privateKeyInst, sgType)
						if err != nil {
							return errors.New("ProveToRV32: Error generating ProveToRV32. " + err.Error())
						}
						// Above needs to be refactored...

						acceptOwner23, err := to1requestor.ProveToRV32(proveToRV32)
						if err != nil {
							log.Panic(err)
						}

						log.Println(acceptOwner23)
					}

					return nil
				},
			},
		},
	}

	err := cliapp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
