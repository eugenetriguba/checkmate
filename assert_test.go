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

// assertFn defines a common interface to be able to
// wrap all the assertion functions in the library to
// a common interface to make testing their common behavior
// with testing tables easier.
type assertFn func(t TestingT, args []any)

func wrapCheck(t TestingT, args []any) {
	if len(args) > 1 {
		Check(t, args[0].(bool), args[1:]...)
	} else {
		Check(t, args[0].(bool))
	}
}

func wrapAssert(t TestingT, args []any) {
	if len(args) > 1 {
		Assert(t, args[0].(bool), args[1:]...)
	} else {
		Assert(t, args[0].(bool))
	}
}

func wrapAssertEqual(t TestingT, args []any) {
	if len(args) > 2 {
		AssertEqual(t, args[0], args[1], args[2:]...)

	} else {
		AssertEqual(t, args[0], args[1])
	}
}

func wrapAssertDeepEqual(t TestingT, args []any) {
	if len(args) > 2 {
		AssertDeepEqual(t, args[0], args[1], args[2:]...)

	} else {
		AssertDeepEqual(t, args[0], args[1])
	}
}

func wrapAssertNotEqual(t TestingT, args []any) {
	if len(args) > 2 {
		AssertNotEqual(t, args[0], args[1], args[2:]...)

	} else {
		AssertNotEqual(t, args[0], args[1])
	}
}

func wrapAssertLenEqualInt(t TestingT, args []any) {
	slice, ok := args[0].([]int)
	if !ok {
		t.Log("Invalid argument type for AssertLenEqualInt")
	}

	expectedLen, ok := args[1].(int)
	if !ok {
		t.Log("Invalid length argument for AssertLenEqualInt")
	}

	if len(args) > 2 {
		AssertLenEqual(t, slice, expectedLen, args[2:]...)

	} else {
		AssertLenEqual(t, slice, expectedLen)
	}
}

var passingTestFns = []struct {
	name        string
	assertionFn assertFn
	args        []any
}{
	{"Check", wrapCheck, []any{true}},
	{"Assert", wrapAssert, []any{true}},
	{"AssertEqual", wrapAssertEqual, []any{5, 5}},
	{"AssertDeepEqual", wrapAssertDeepEqual, []any{5, 5}},
	{"AssertNotEqual", wrapAssertNotEqual, []any{5, 10}},
	{"AssertLenEqual", wrapAssertLenEqualInt, []any{[]int{1, 2}, 2}},
}

var failingTestFns = []struct {
	name        string
	assertionFn assertFn
	args        []any
}{
	{"Check", wrapCheck, []any{false}},
	{"Assert", wrapAssert, []any{false}},
	{"AssertEqual", wrapAssertEqual, []any{5, 10}},
	{"AssertDeepEqual", wrapAssertDeepEqual, []any{5, 10}},
	{"AssertNotEqual", wrapAssertNotEqual, []any{5, 5}},
	{"AssertLenEqual", wrapAssertLenEqualInt, []any{[]int{1, 2}, 0}},
}

func TestOptionalMessageAndArgs(t *testing.T) {
	testCases := []struct {
		name         string
		args         []any
		expectedLogs []string
	}{
		{"Plain message", []any{"my message"}, []string{"my message"}},
		{"Message with format placeholders", []any{"my message: %d", 5}, []string{"my message: 5"}},
		{
			"Invalid message type",
			[]any{5},
			[]string{
				"checkmate: assertion called with a non-string message, using default message",
				"check failed",
			},
		},
	}

	for _, testFn := range failingTestFns {
		for _, testCase := range testCases {
			testName := fmt.Sprintf("%s: %s", testFn.name, testCase.name)
			t.Run(testName, func(t *testing.T) {
				mockT := &MockT{}

				testFn.assertionFn(mockT, append(testFn.args, testCase.args...))

				if len(mockT.Logs) != len(testCase.expectedLogs) {
					t.Errorf(
						"%s: expected %d (%v) log messages but got %d (%v)",
						testName, len(mockT.Logs), mockT.Logs, len(testCase.expectedLogs),
						testCase.expectedLogs,
					)
				}
				if len(mockT.Logs) == len(testCase.expectedLogs) {
					for i, expectedLogMsg := range testCase.expectedLogs {
						if mockT.Logs[i] != expectedLogMsg {
							t.Errorf(
								"%s: expected log message '%s', got '%s'",
								testName, expectedLogMsg, mockT.Logs[i],
							)
						}
					}
				}
			})
		}
	}
}

func TestAssertFnsCallHelper(t *testing.T) {
	for _, testFn := range failingTestFns {
		t.Run(testFn.name, func(t *testing.T) {
			mockHelperT := &MockHelperT{}

			testFn.assertionFn(mockHelperT, testFn.args)

			if !mockHelperT.HelperCalled {
				t.Errorf("%s: expected HelperT to be called", testFn.name)
			}
		})
	}
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

	t.Run(
		"should default to standard message if incorrect type for message passed",
		func(t *testing.T) {
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
					expectedFailureLog,
					mockT.Logs[0],
				)
			}
		},
	)

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
	testCases := []struct {
		name        string
		actual      any
		expected    any
		shouldFail  bool
		logMessages []string
	}{
		{"EqualIntegers", 5, 5, false, []string{}},
		{"UnequalIntegers", 5, 10, true, []string{"expected 10 to equal 5"}},
		{"EqualFloats", 5.123, 5.123, false, []string{}},
		{"UnequalFloats", 5.123, 5.1234, true, []string{"expected 5.1234 to equal 5.123"}},
		{"EqualStrings", "test", "test", false, []string{}},
		{"UnequalStrings", "test", "fail", true, []string{"expected fail to equal test"}},
		{"EqualBooleans", true, true, false, []string{}},
		{"UnequalBooleans", false, true, true, []string{"expected true to equal false"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &MockT{}
			AssertEqual(mockT, tc.actual, tc.expected)

			if mockT.FailNowCalled != tc.shouldFail {
				t.Errorf(
					"%s: FailNowCalled = %v, want %v",
					tc.name, mockT.FailNowCalled, tc.shouldFail,
				)
			}
			if len(mockT.Logs) != len(tc.logMessages) {
				t.Errorf(
					"%s: expected %d (%v) log messages but got %d (%v)",
					tc.name, len(mockT.Logs), mockT.Logs, len(tc.logMessages), tc.logMessages,
				)
			}
			if len(mockT.Logs) == len(tc.logMessages) {
				for i, expectedLogMsg := range tc.logMessages {
					if mockT.Logs[i] != expectedLogMsg {
						t.Errorf(
							"%s: expected log message '%s', got '%s'",
							tc.name, expectedLogMsg, mockT.Logs[0],
						)
					}
				}
			}
		})
	}
}

func TestAssertEqualCallsHelper(t *testing.T) {
	mockHelperT := &MockHelperT{}

	AssertEqual(mockHelperT, true, true)

	if !mockHelperT.HelperCalled {
		t.Fatal("HelperT not called")
	}
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
