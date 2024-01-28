package checkmate

import (
	"fmt"
	"strings"
	"testing"
)

// MockT is a mock structure for capturing test outputs.
type MockT struct {
	TestingT

	FailNowCalled bool
	Logs          []string
}

func (m *MockT) FailNow() {
	m.FailNowCalled = true
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

func TestCheck(t *testing.T) {
	t.Run("should not fail or emit any logs on true condition", func(t *testing.T) {
		mockT := &MockT{}

		passed := Check(mockT, true)

		if !passed {
			t.Fatal("Check should have returned true on true condition")
		}
		if mockT.FailNowCalled {
			t.Fatal("Check should not be failing any tests")
		}
		if len(mockT.Logs) > 0 {
			t.Fatal("Check should not emit any logs on success")
		}
	})

	t.Run("should not fail but it should emit a log on false condition", func(t *testing.T) {
		mockT := &MockT{}

		passed := Check(mockT, false)

		if passed {
			t.Fatal("Check should have returned false on false condition")
		}
		if mockT.FailNowCalled {
			t.Fatal("Check should not be failing any tests")
		}
		if len(mockT.Logs) == 0 {
			t.Fatal("Check should be emiting a logs on failure")
		}
	})

	t.Run("should allow optional message with no format placeholders", func(t *testing.T) {
		mockT := &MockT{}
		failureLog := "my failure message"

		Check(mockT, false, failureLog)

		if len(mockT.Logs) != 1 {
			t.Fatal("Check should emit one failure log")
		}
		if mockT.Logs[0] != failureLog {
			t.Fatalf(
				"Check should emit custom message %s on failure, got %s",
				failureLog, mockT.Logs[0],
			)
		}
	})

	t.Run("should allow optional message with format placeholders filled in", func(t *testing.T) {
		mockT := &MockT{}

		Check(mockT, false, "my failure message: %d", 5)

		if len(mockT.Logs) != 1 {
			t.Fatal("Check should emit one failure log")
		}
		expectedFailureLog := "my failure message: 5"
		if mockT.Logs[0] != expectedFailureLog {
			t.Fatalf(
				"Check should emit custom message %s on failure, got %s",
				expectedFailureLog, mockT.Logs[0],
			)
		}
	})

	t.Run("should default to standard message if incorrect type for message passed", func(t *testing.T) {
		mockT := &MockT{}

		Check(mockT, false, 5)

		if len(mockT.Logs) != 2 {
			t.Fatal("Check should emit one failure log and one warning log")
		}
		expectedWarningLog := "checkmate: assertion called with a non-string message, using default message"
		if mockT.Logs[0] != expectedWarningLog {
			t.Fatalf(
				"Check should emit warning message %s on wrong type for message, got %s",
				expectedWarningLog, mockT.Logs[0],
			)
		}
		expectedFailureLog := "check failed"
		if mockT.Logs[1] != expectedFailureLog {
			t.Fatalf(
				"Check should emit default failure message %s on wrong type for message, got %s",
				expectedFailureLog, mockT.Logs[0],
			)
		}
	})

	t.Run("should called Helper() when passed a helperT", func(t *testing.T) {
		mockHelperT := &MockHelperT{}

		Check(mockHelperT, true)

		if !mockHelperT.HelperCalled {
			t.Fatal("HelperT not called")
		}
	})
}

func TestAssert(t *testing.T) {
	t.Run("should not fail test or emit log on true assertion", func(t *testing.T) {
		mockT := &MockT{}

		Assert(mockT, true)

		if mockT.FailNowCalled {
			t.Fatal("Assert failed when it should have passed")
		}
		if len(mockT.Logs) > 0 {
			t.Fatal("Check should not emit any logs on success")
		}
	})

	t.Run("should call fail now and emit a log on false assertion", func(t *testing.T) {
		mockT := &MockT{}

		Assert(mockT, false)

		if !mockT.FailNowCalled {
			t.Fatal("Assert passed when it should have failed")
		}
		if len(mockT.Logs) == 0 {
			t.Fatal("Assert should have emitted a failure log")
		}
	})

	t.Run("should allow optional message with no format placeholders", func(t *testing.T) {
		mockT := &MockT{}
		failureLog := "my failure message"

		Assert(mockT, false, failureLog)

		if !mockT.FailNowCalled {
			t.Fatal("Assert passed when it should have failed")
		}
		if len(mockT.Logs) != 1 {
			t.Fatalf("Assert should emit one failure log, got %d logs", len(mockT.Logs))
		}
		if mockT.Logs[0] != failureLog {
			t.Fatalf(
				"Assert should emit custom message %s on failure, got %s",
				failureLog, mockT.Logs[0],
			)
		}
	})

	t.Run("should allow optional message with format placeholders filled in", func(t *testing.T) {
		mockT := &MockT{}

		Assert(mockT, false, "my failure message: %d", 5)

		if !mockT.FailNowCalled {
			t.Fatal("Assert passed when it should have failed")
		}
		if len(mockT.Logs) != 1 {
			t.Fatalf("Assert should emit one failure log, got %d logs", len(mockT.Logs))
		}
		expectedFailureLog := "my failure message: 5"
		if mockT.Logs[0] != expectedFailureLog {
			t.Fatalf(
				"Assert should emit custom message %s on failure, got %s",
				expectedFailureLog, mockT.Logs[0],
			)
		}
	})

	t.Run("should called Helper() when passed a helperT", func(t *testing.T) {
		mockHelperT := &MockHelperT{}

		Assert(mockHelperT, true)

		if !mockHelperT.HelperCalled {
			t.Fatal("HelperT not called")
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

	t.Run("should called Helper() when passed a helperT", func(t *testing.T) {
		mockHelperT := &MockHelperT{}

		AssertEqual(mockHelperT, true, true)

		if !mockHelperT.HelperCalled {
			t.Fatal("HelperT not called")
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

	t.Run("should called Helper() when passed a helperT", func(t *testing.T) {
		mockHelperT := &MockHelperT{}

		AssertNotEqual(mockHelperT, true, false)

		if !mockHelperT.HelperCalled {
			t.Fatal("HelperT not called")
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

	t.Run("should called Helper() when passed a helperT", func(t *testing.T) {
		mockHelperT := &MockHelperT{}

		AssertLenEqual(mockHelperT, []string{""}, 1)

		if !mockHelperT.HelperCalled {
			t.Fatal("HelperT not called")
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

	t.Run("should called Helper() when passed a helperT", func(t *testing.T) {
		mockHelperT := &MockHelperT{}

		AssertDeepEqual(mockHelperT, []string{""}, []string{""})

		if !mockHelperT.HelperCalled {
			t.Fatal("HelperT not called")
		}
	})
}

func containsDiffMessage(log string) bool {
	return strings.Contains(log, "(-expected +actual)")
}
