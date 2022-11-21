package listener

import (
	"time"

	fdoshared "github.com/WebauthnWorks/fdo-shared"
	"github.com/WebauthnWorks/fdo-shared/testcom"
	"github.com/google/uuid"
)

type ListenerTestRun struct {
	_         struct{}                `cbor:",toarray"`
	Uuid      string                  `json:"uuid"`
	Timestamp int64                   `json:"timestamp"`
	TestRuns  []testcom.FDOTestState  `json:"tests"`
	Protocol  fdoshared.FdoToProtocol `json:"protocol"`
	Completed bool                    `json:"completed"`
}

func NewListenerTestRun(protocol fdoshared.FdoToProtocol) ListenerTestRun {
	newUuid, _ := uuid.NewRandom()
	uuidStr, _ := newUuid.MarshalText()
	newRVTestRun := ListenerTestRun{
		Uuid:      string(uuidStr),
		Timestamp: time.Now().Unix(),
		TestRuns:  []testcom.FDOTestState{},
		Protocol:  protocol,
	}

	return newRVTestRun
}

func (h *ListenerTestRun) Complete() {
	h.Completed = true
}
