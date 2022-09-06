package dbs

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/WebauthnWorks/fdo-fido-conformance-server/req_tests_deps"
	"github.com/WebauthnWorks/fdo-fido-conformance-server/testcom"
	"github.com/dgraph-io/badger/v3"
	"github.com/fxamacker/cbor/v2"
)

type RequestTestDB struct {
	db     *badger.DB
	prefix []byte
	ttl    int
}

func NewRequestTestDB(db *badger.DB) RequestTestDB {
	return RequestTestDB{
		db:     db,
		prefix: []byte("rvte-"),
		ttl:    60 * 60 * 24 * 183, //6months storage
	}
}

func (h *RequestTestDB) Save(rvte req_tests_deps.RequestTestInst) error {
	rvteBytes, err := cbor.Marshal(rvte)
	if err != nil {
		return errors.New("Failed to marshal rvte. The error is: " + err.Error())
	}

	rvteStorageId := append(h.prefix, rvte.Uuid...)

	dbtxn := h.db.NewTransaction(true)
	defer dbtxn.Discard()

	entry := badger.NewEntry(rvteStorageId, rvteBytes).WithTTL(time.Second * time.Duration(h.ttl)) // Session entry will only exist for 10 minutes
	err = dbtxn.SetEntry(entry)
	if err != nil {
		return errors.New("Failed creating rvte db entry instance. The error is: " + err.Error())
	}

	err = dbtxn.Commit()
	if err != nil {
		return errors.New("Failed saving rvte entry. The error is: " + err.Error())
	}

	return nil
}

func (h *RequestTestDB) Update(rvtId []byte, rvte req_tests_deps.RequestTestInst) error {
	rvteBytes, err := cbor.Marshal(rvte)
	if err != nil {
		return errors.New("Failed to marshal rvte. The error is: " + err.Error())
	}

	rvteStorageId := append(h.prefix, rvtId...)

	dbtxn := h.db.NewTransaction(true)
	defer dbtxn.Discard()

	entry := badger.NewEntry(rvteStorageId, rvteBytes).WithTTL(time.Second * time.Duration(h.ttl)) // Session entry will only exist for 10 minutes
	err = dbtxn.SetEntry(entry)
	if err != nil {
		return errors.New("Failed creating rvte db entry instance. The error is: " + err.Error())
	}

	err = dbtxn.Commit()
	if err != nil {
		return errors.New("Failed saving rvte entry. The error is: " + err.Error())
	}

	return nil
}

func (h *RequestTestDB) Get(rvtId []byte) (*req_tests_deps.RequestTestInst, error) {
	rvteStorageId := append(h.prefix, rvtId...)

	dbtxn := h.db.NewTransaction(true)
	defer dbtxn.Discard()

	item, err := dbtxn.Get(rvteStorageId)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, fmt.Errorf("The rvte entry with id %s does not exist", hex.EncodeToString(rvtId))
	} else if err != nil {
		return nil, errors.New("Failed locating rvte entry. The error is: " + err.Error())
	}

	itemBytes, err := item.ValueCopy(nil)
	if err != nil {
		return nil, errors.New("Failed reading rvte entry value. The error is: " + err.Error())
	}

	var rvteInst req_tests_deps.RequestTestInst
	err = cbor.Unmarshal(itemBytes, &rvteInst)
	if err != nil {
		return nil, errors.New("Failed cbor decoding rvte entry value. The error is: " + err.Error())
	}

	return &rvteInst, nil
}

func (h *RequestTestDB) GetMany(rvtids [][]byte) (*[]req_tests_deps.RequestTestInst, error) {
	var rvts []req_tests_deps.RequestTestInst

	for _, rvtid := range rvtids {
		rvt, err := h.Get(rvtid)
		if err != nil {
			return nil, fmt.Errorf("Error obtaining rvt for id %s. %s \n", hex.EncodeToString(rvtid), err.Error())
		}

		rvts = append(rvts, *rvt)
	}

	return &rvts, nil
}

func (h *RequestTestDB) StartNewRun(rvteid []byte) {
	rvte, err := h.Get(rvteid)
	if err != nil {
		log.Printf("%s test entry can not be found.", hex.EncodeToString(rvteid))
	}

	newRVTestRun := req_tests_deps.NewRVTestRun(rvte.Protocol)

	rvte.InProgress = true
	rvte.CurrentTestRun = newRVTestRun
	rvte.TestsHistory = append([]req_tests_deps.RequestTestRun{newRVTestRun}, rvte.TestsHistory...)

	err = h.Save(*rvte)
	if err != nil {
		log.Printf("%s error saving test entry.", hex.EncodeToString(rvteid))
	}
}

func (h *RequestTestDB) FinishRun(rvteid []byte) {
	rvte, err := h.Get(rvteid)
	if err != nil {
		log.Printf("%s test entry can not be found.", hex.EncodeToString(rvteid))
	}

	rvte.InProgress = false

	err = h.Save(*rvte)
	if err != nil {
		log.Printf("%s error saving test entry.", hex.EncodeToString(rvteid))
	}
}

func (h *RequestTestDB) ReportTest(rvteid []byte, testID testcom.FDOTestID, testResult testcom.FDOTestState) {
	rvte, err := h.Get(rvteid)
	if err != nil {
		log.Printf("%s test entry can not be found.", hex.EncodeToString(rvteid))
	}

	rvte.CurrentTestRun.Tests[testID] = testResult
	rvte.TestsHistory[0] = rvte.CurrentTestRun

	err = h.Save(*rvte)
	if err != nil {
		log.Printf("%s error saving test entry.", hex.EncodeToString(rvteid))
	}
}

func (h *RequestTestDB) RemoveTestRun(rvteid []byte, testRunId string) {
	rvte, err := h.Get(rvteid)
	if err != nil {
		log.Printf("%s test entry can not be found.", hex.EncodeToString(rvteid))
	}

	var updatedTestsHistory []req_tests_deps.RequestTestRun = []req_tests_deps.RequestTestRun{}
	for _, testRunEntry := range rvte.TestsHistory {
		if testRunEntry.Uuid != testRunId {
			updatedTestsHistory = append(updatedTestsHistory, testRunEntry)
		}
	}

	rvte.TestsHistory = updatedTestsHistory

	err = h.Save(*rvte)
	if err != nil {
		log.Printf("%s error saving test entry.", hex.EncodeToString(rvteid))
	}
}