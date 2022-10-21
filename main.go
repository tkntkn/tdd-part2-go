package main

import (
	"fmt"
	"reflect"
)

/**
 * Common Method
 */
func assert(ok bool) {
	if !ok {
		panic("assertion failed")
	}
}

/**
 * TestCase
 */
type ITestCase interface {
	GetName() string
}

type TestCase struct {
	Name string
}

func (c *TestCase) GetName() string {
	return c.Name
}

func RunTestMethod(t ITestCase, r *TestResult) {
	defer func() {
		if err := recover(); err != nil {
			r.TestFailed()
		}
	}()

	method := reflect.ValueOf(t).MethodByName(t.GetName())
	method.Call([]reflect.Value{})
}

func RunTest(t ITestCase, result *TestResult) {
	{
		result.TestStarted()

		setUp := reflect.ValueOf(t).MethodByName("SetUp")
		if setUp.IsValid() {
			setUp.Call([]reflect.Value{})
		}

		RunTestMethod(t, result)

		tearDown := reflect.ValueOf(t).MethodByName("TearDown")
		if tearDown.IsValid() {
			tearDown.Call([]reflect.Value{})
		}
	}
}

/**
 * Test Result
 */
type TestResult struct {
	RunCount   int
	ErrorCount int
}

func (r *TestResult) TestStarted() {
	r.RunCount = r.RunCount + 1
}

func (r *TestResult) TestFailed() {
	r.ErrorCount = r.ErrorCount + 1
}

func (r *TestResult) Summary() string {
	return fmt.Sprintf("%d run, %d failed", r.RunCount, r.ErrorCount)
}

/**
 * Test Suite
 */
type TestSuite struct {
	Tests []ITestCase
}

func (s *TestSuite) Add(test ITestCase) {
	s.Tests = append(s.Tests, test)
}

func RunTests(s TestSuite, result *TestResult) {
	for _, test := range s.Tests {
		RunTest(test, result)
	}
}

/**
 * WasRun
 */
type WasRun struct {
	TestCase
	WasRun int
	Log    string
}

func (w *WasRun) SetUp() {
	w.WasRun = -1
	w.Log = "SetUp "
}

func (w *WasRun) TestMethod() {
	w.WasRun = 1
	w.Log = w.Log + "TestMethod "
}

func (w *WasRun) TestBrokenMethod() {
	panic("Broken Method")
}

func (w *WasRun) TearDown() {
	w.WasRun = 1
	w.Log = w.Log + "TearDown "
}

/**
 * Test for WasRun TestCase
 */
type TestCaseTest struct {
	TestCase
}

func (t *TestCaseTest) Test_TestCase_TemplateMethod() {
	test := &WasRun{TestCase: TestCase{"TestMethod"}}
	result := &TestResult{}
	RunTest(test, result)
	assert(test.Log == "SetUp TestMethod TearDown ")
}

func (t *TestCaseTest) Test_TestCase_Result() {
	test := &WasRun{TestCase: TestCase{"TestMethod"}}
	result := &TestResult{}
	RunTest(test, result)
	assert(result.Summary() == "1 run, 0 failed")
}

func (t *TestCaseTest) Test_TestCase_FailedResult() {
	test := &WasRun{TestCase: TestCase{"TestBrokenMethod"}}
	result := &TestResult{}
	RunTest(test, result)
	assert(result.Summary() == "1 run, 1 failed")
}

func (t *TestCaseTest) Test_TestResult_FailedResultFormatting() {
	result := &TestResult{}
	result.TestStarted()
	result.TestFailed()
	assert(result.Summary() == "1 run, 1 failed")
}

func (t *TestCaseTest) Test_TestSuite() {
	suite := TestSuite{}
	suite.Add(&WasRun{TestCase: TestCase{"TestMethod"}})
	suite.Add(&WasRun{TestCase: TestCase{"TestBrokenMethod"}})
	result := &TestResult{}
	RunTests(suite, result)
	assert(result.Summary() == "2 run, 1 failed")
}

/**
 * Main
 */
func main() {
	suite := TestSuite{}
	suite.Add(&TestCaseTest{TestCase{"Test_TestCase_TemplateMethod"}})
	suite.Add(&TestCaseTest{TestCase{"Test_TestCase_Result"}})
	suite.Add(&TestCaseTest{TestCase{"Test_TestCase_FailedResult"}})
	suite.Add(&TestCaseTest{TestCase{"Test_TestResult_FailedResultFormatting"}})
	suite.Add(&TestCaseTest{TestCase{"Test_TestSuite"}})

	result := &TestResult{}
	RunTests(suite, result)
	fmt.Println(result.Summary())
}
