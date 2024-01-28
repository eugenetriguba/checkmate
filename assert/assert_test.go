package assert

import (
	"errors"
	"fmt"
	"os"
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
		Equal(t, args[0], args[1], args[2:]...)

	} else {
		Equal(t, args[0], args[1])
	}
}

func wrapAssertDeepEqual(t TestingT, args []any) {
	if len(args) > 2 {
		DeepEqual(t, args[0], args[1], args[2:]...)

	} else {
		DeepEqual(t, args[0], args[1])
	}
}

func wrapAssertNotEqual(t TestingT, args []any) {
	if len(args) > 2 {
		NotEqual(t, args[0], args[1], args[2:]...)

	} else {
		NotEqual(t, args[0], args[1])
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
		LenEqual(t, slice, expectedLen, args[2:]...)

	} else {
		LenEqual(t, slice, expectedLen)
	}
}

func wrapAssertErrorIs(t TestingT, args []any) {
	if len(args) > 2 {
		ErrorIs(t, args[0].(error), args[1].(error), args[2:]...)

	} else {
		ErrorIs(t, args[0].(error), args[1].(error))
	}
}

func wrapAssertErrorContains(t TestingT, args []any) {
	if len(args) > 2 {
		ErrorContains(t, args[0].(error), args[1].(string), args[2:]...)

	} else {
		ErrorContains(t, args[0].(error), args[1].(string))
	}
}

func wrapAssertTrue(t TestingT, args []any) {
	if len(args) > 1 {
		True(t, args[0].(bool), args[1:]...)

	} else {
		True(t, args[0].(bool))
	}
}

func wrapAssertFalse(t TestingT, args []any) {
	if len(args) > 1 {
		False(t, args[0].(bool), args[1:]...)

	} else {
		False(t, args[0].(bool))
	}
}

func wrapAssertNil(t TestingT, args []any) {
	if len(args) > 1 {
		Nil(t, args[0], args[1:]...)

	} else {
		Nil(t, args[0])
	}
}

func wrapAssertNotNil(t TestingT, args []any) {
	if len(args) > 1 {
		NotNil(t, args[0], args[1:]...)

	} else {
		NotNil(t, args[0])
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
	{"AssertErrorIs", wrapAssertErrorIs, []any{os.ErrExist, os.ErrExist}},
	{"AssertErrorContains", wrapAssertErrorContains, []any{errors.New("error 1"), "error 1"}},
	{"AssertTrue", wrapAssertTrue, []any{true}},
	{"AssertFalse", wrapAssertFalse, []any{false}},
	{"AssertNil", wrapAssertNil, []any{nil}},
	{"AssertNotNil", wrapAssertNotNil, []any{1}},
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
	{"AssertErrorIs", wrapAssertErrorIs, []any{errors.New("error 1"), errors.New("error 2")}},
	{"AssertErrorContains", wrapAssertErrorContains, []any{errors.New("error 1"), "error 2"}},
	{"AssertTrue", wrapAssertTrue, []any{false}},
	{"AssertFalse", wrapAssertFalse, []any{true}},
	{"AssertNil", wrapAssertNil, []any{1}},
	{"AssertNotNil", wrapAssertNotNil, []any{nil}},
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

func TestAssertFnsShouldNotFailOrEmitLogsOnSuccess(t *testing.T) {
	for _, testFn := range passingTestFns {
		t.Run(testFn.name, func(t *testing.T) {
			mockT := &MockT{}
			testFn.assertionFn(mockT, testFn.args)

			if mockT.FailNowCalled {
				t.Fatalf(
					"%s: fail now called, expected the test to not be failed on true condition",
					testFn.name,
				)
			}
			if len(mockT.Logs) > 0 {
				t.Fatalf(
					"%s: %v logs emitted, expected the test to not emit any logs on true condition",
					testFn.name, mockT.Logs,
				)
			}
		})
	}
}

func TestCheckReturnValues(t *testing.T) {
	t.Run("should return true on true condition", func(t *testing.T) {
		mockT := &MockT{}

		passed := Check(mockT, true)

		if !passed {
			t.Fatal("Check should have returned true on true condition")
		}
	})

	t.Run("should return false on false condition", func(t *testing.T) {
		mockT := &MockT{}

		passed := Check(mockT, false)

		if passed {
			t.Fatal("Check should have returned false on false condition")
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
			Equal(mockT, tc.actual, tc.expected)

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

func TestAssertDeepEqual(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	t.Run("Equal structs", func(t *testing.T) {
		mockT := &MockT{}
		DeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Alice", 30})
		if mockT.FailNowCalled {
			t.Error("AssertDeepEqual failed when it should have passed")
		}
	})

	t.Run("Not equal structs", func(t *testing.T) {
		mockT := &MockT{}
		DeepEqual(mockT, TestStruct{"Alice", 30}, TestStruct{"Bob", 30})
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

func TestAssertErrorIsChecksErrTree(t *testing.T) {
	mockT := &MockT{}
	actual := fmt.Errorf("wrapped: %w", os.ErrInvalid)
	expected := os.ErrInvalid

	ErrorIs(mockT, actual, expected)

	if mockT.FailNowCalled {
		t.Fatal("AssertErrorIs failed when it should have passed")
	}
}
