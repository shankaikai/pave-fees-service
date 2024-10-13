package workflow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_BillCreation() {
	bill := Bill{
		LineItems: make([]LineItem, 0),
		Currency:  "USD",
		TotalAmount: 0.0,
	}
	
	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow("getBill")
		s.NoError(err)
		err = res.Get(&bill)
		s.NoError(err)
		s.Equal(len(bill.LineItems), 0)
		s.Equal("USD", bill.Currency)
	}, time.Millisecond)

	s.env.ExecuteWorkflow(BillWorkflow, bill)
}

func (s *UnitTestSuite) Test_BillCreationWithLineItem() {
	bill := Bill{
		LineItems: make([]LineItem, 0),
		Currency:  "USD",
		TotalAmount: 0.0,
	}
	
	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow("getBill")
		s.NoError(err)
		err = res.Get(&bill)
		s.NoError(err)
		s.Equal(len(bill.LineItems), 0)
		s.Equal("USD", bill.Currency)
		s.Equal(0.0, bill.TotalAmount)

		s.env.SignalWorkflow("addLineItem", AddLineItemSignal{
			Description: "item1",
			Amount:      10.0,
		})
		s.env.SignalWorkflow("addLineItem", AddLineItemSignal{
			Description: "item2",
			Amount:      11.0,
		})
	}, time.Millisecond)

	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow("getBill")
		s.NoError(err)
		err = res.Get(&bill)
		s.NoError(err)
		s.Equal(len(bill.LineItems), 2)
		s.Equal("USD", bill.Currency)
		s.Equal(21.0, bill.TotalAmount)
		s.Equal("item1", bill.LineItems[0].Description)
		s.Equal(10.0, bill.LineItems[0].Amount)
		s.Equal("item2", bill.LineItems[1].Description)
		s.Equal(11.0, bill.LineItems[1].Amount)
	}, time.Millisecond * 2)

	s.env.ExecuteWorkflow(BillWorkflow, bill)
}

func (s *UnitTestSuite) Test_BillClose() {
	bill := Bill{
		LineItems: make([]LineItem, 0),
		Currency:  "USD",
		TotalAmount: 0.0,
	}
	
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("closeBill", CloseBillSignal{})
	}, time.Millisecond)

	s.env.RegisterDelayedCallback(func() {
		s.True(s.env.IsWorkflowCompleted())
	}, time.Millisecond * 3)

	s.env.ExecuteWorkflow(BillWorkflow, bill)
}