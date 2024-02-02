package cmtest

import (
	"fmt"

	"github.com/eugenetriguba/checkmate"
)

// MockT is a mock structure for capturing test outputs.
type MockT struct {
	checkmate.TestingT

	FailNowCalled bool
	FailCalled    bool
	Logs          []string
}

func (m *MockT) FailNow() {
	m.FailNowCalled = true
}

func (m *MockT) Fail() {
	m.FailCalled = true
}

func (m *MockT) Log(args ...any) {
	m.Logs = append(m.Logs, fmt.Sprint(args...))
}

// MockHelperT has a MockT and is also a helperT.
type MockHelperT struct {
	MockT

	HelperCalled bool
}

func (m *MockHelperT) Helper() {
	m.HelperCalled = true
}
