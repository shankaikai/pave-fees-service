package workflow

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) TestWorkflowCompletesOnCloseSignal() {
	env := s.NewTestWorkflowEnvironment()
	env.ExecuteWorkflow(BillWorkflow, "USD")

	env.SignalWorkflow("closeBill", CloseBillSignal{})

	s.True(env.IsWorkflowCompleted())
}

