package assert

import (
	"fmt"
	"strings"
	"testing"
)

// MockT is a mock testing.T structure for capturing test outputs.
type MockT struct {
	TestingT

	HelperCalled  bool
	FailNowCalled bool
	Logs          []string
}

func (m *MockT) Helper() {
	m.HelperCalled = true
}

func (m *MockT) FailNow() {
	m.FailNowCalled = true
}

func (m *MockT) Log(args ...any) {
	m.Logs = append(m.Logs, fmt.Sprint(args...))
}

func TestAssert(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		mockT := &MockT{}
		Assert(mockT, true)
		if mockT.FailNowCalled {
			t.Fatal("Assert failed when it should have passed")
		}
	})

	t.Run("False", func(t *testing.T) {
		mockT := &MockT{}
		Assert(mockT, false)
		if !mockT.FailNowCalled {
			t.Fatal("Assert passed when it should have failed")
		}
	})
}

func TestAssertEqual(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		mockT := &MockT{}
		AssertEqual(mockT, 5, 5)
		if mockT.FailNowCalled {
			t.Fatal("AssertEqual failed when it should have passed")
		}
	})

	t.Run("Not equal", func(t *testing.T) {
		mockT := &MockT{}
		AssertEqual(mockT, 5, 10)
		if !mockT.FailNowCalled {
			t.Fatal("AssertEqual passed when it should have failed")
		}
	})
}

func TestAssertNotEqual(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		mockT := &MockT{}
		AssertNotEqual(mockT, 5, 5)
		if !mockT.FailNowCalled {
			t.Fatal("AssertNotEqual failed when it should have passed")
		}
	})

	t.Run("Not equal", func(t *testing.T) {
		mockT := &MockT{}
		AssertNotEqual(mockT, 5, 10)
		if mockT.FailNowCalled {
			t.Fatal("AssertNotEqual passed when it should have failed")
		}
	})
}

func TestAssertLenEqual(t *testing.T) {
	t.Run("Equal lengths", func(t *testing.T) {
		mockT := &MockT{}
		AssertLenEqual(mockT, []int{1, 2, 3}, 3)
		if mockT.FailNowCalled {
			t.Error("AssertLenEqual failed when it should have passed")
		}
	})

	t.Run("Lengths aren't equal", func(t *testing.T) {
		mockT := &MockT{}
		AssertLenEqual(mockT, []int{1, 2, 3}, 2)
		if !mockT.FailNowCalled {
			t.Error("AssertLenEqual passed when it should have failed")
		}
	})
}

func TestAssertDeepEqual(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	t.Run("Equal structs", func(t *testing.T) {
		mockT := &MockT{}
		AssertDeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Alice", 30})
		if mockT.FailNowCalled {
			t.Error("AssertDeepEqual failed when it should have passed")
		}
	})

	t.Run("Not equal structs", func(t *testing.T) {
		mockT := &MockT{}
		AssertDeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Bob", 30})
		if !mockT.FailNowCalled {
			t.Error("AssertDeepEqual passed when it should have failed")
		}
		if len(mockT.Logs) == 0 || !containsDiffMessage(mockT.Logs[0]) {
			t.Error("AssertDeepEqual did not log the correct diff message")
		}
	})
}

func containsDiffMessage(log string) bool {
	return strings.Contains(log, "(-expected +actual)")
}
