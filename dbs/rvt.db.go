package dbs

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/WebauthnWorks/fdo-fido-conformance-server/rvtests"
	"github.com/dgraph-io/badger/v3"
	"github.com/fxamacker/cbor/v2"
)

type RendezvousServerTestDB struct {
	db *badger.DB
}

var RVT_TTLS int = 60 * 60 * 24 * 183 //6months storage

var rvtdbpref []byte = []byte("rvte-")

func NewRendezvousServerTestDB(db *badger.DB) RendezvousServerTestDB {
	return RendezvousServerTestDB{
		db: db,
	}
}

func (h *RendezvousServerTestDB) Save(rvte rvtests.RendezvousServerTestDBEntry) error {
	rvteBytes, err := cbor.Marshal(rvte)
	if err != nil {
		return errors.New("Failed to marshal rvte. The error is: " + err.Error())
	}

	rvteStorageId := append(rvtdbpref, rvte.ID...)

	dbtxn := h.db.NewTransaction(true)
	defer dbtxn.Discard()

	entry := badger.NewEntry(rvteStorageId, rvteBytes).WithTTL(time.Second * time.Duration(RVT_TTLS)) // Session entry will only exist for 10 minutes
	err = dbtxn.SetEntry(entry)
	if err != nil {
		return errors.New("Failed creating rvte db entry instance. The error is: " + err.Error())
	}

	dbtxn.Commit()
	if err != nil {
		return errors.New("Failed saving rvte entry. The error is: " + err.Error())
	}

	return nil
}

func (h *RendezvousServerTestDB) Update(rvtId []byte, rvte rvtests.RendezvousServerTestDBEntry) error {
	rvteBytes, err := cbor.Marshal(rvte)
	if err != nil {
		return errors.New("Failed to marshal rvte. The error is: " + err.Error())
	}

	rvteStorageId := append(rvtdbpref, rvtId...)

	dbtxn := h.db.NewTransaction(true)
	defer dbtxn.Discard()

	entry := badger.NewEntry(rvteStorageId, rvteBytes).WithTTL(time.Second * time.Duration(RVT_TTLS)) // Session entry will only exist for 10 minutes
	err = dbtxn.SetEntry(entry)
	if err != nil {
		return errors.New("Failed creating rvte db entry instance. The error is: " + err.Error())
	}

	dbtxn.Commit()
	if err != nil {
		return errors.New("Failed saving rvte entry. The error is: " + err.Error())
	}

	return nil
}

func (h *RendezvousServerTestDB) Get(rvtId []byte) (*rvtests.RendezvousServerTestDBEntry, error) {
	rvteStorageId := append(rvtdbpref, rvtId...)

	dbtxn := h.db.NewTransaction(true)
	defer dbtxn.Discard()

	item, err := dbtxn.Get(rvteStorageId)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, fmt.Errorf("The rvte entry with id %s does not exist", hex.EncodeToString(rvteStorageId))
	} else if err != nil {
		return nil, errors.New("Failed locating rvte entry. The error is: " + err.Error())
	}

	itemBytes, err := item.ValueCopy(nil)
	if err != nil {
		return nil, errors.New("Failed reading rvte entry value. The error is: " + err.Error())
	}

	var rvteInst rvtests.RendezvousServerTestDBEntry
	err = cbor.Unmarshal(itemBytes, &rvteInst)
	if err != nil {
		return nil, errors.New("Failed cbor decoding rvte entry value. The error is: " + err.Error())
	}

	return &rvteInst, nil
}
