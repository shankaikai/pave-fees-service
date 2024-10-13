package fees

import (
	"context"
	"errors"
	"testing"

	workflow "encore.app/fees/workflow"
	"encore.dev/beta/errs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/mocks"
)

type UnitTestSuite struct {
	suite.Suite
}

type MockEncodedValue struct {
  mock.Mock
}

func (m *MockEncodedValue) HasValue() bool {
  args := m.Called()
  return args.Get(0).(bool)
}

var mockBill = workflow.Bill{
	Currency: "USD",
	LineItems: []workflow.LineItem{},
	TotalAmount: 1.0,
}

func (m *MockEncodedValue) Get(valuePtr interface{}) error {
	*valuePtr.(*workflow.Bill) = mockBill
	return nil
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_CreateBill_Success() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}
	mockWorkflowRun :=  mocks.NewWorkflowRun(s.T())
	mockWorkflowRun.On("GetID").Return("123")
	mockClient.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)

	ctx := context.Background()
	req := &CreateBillRequest{
		Currency: "USD",
	}

	resp, err := service.CreateBill(ctx, req)
	s.NoError(err)
	s.Equal("123", resp.Id)
}

func (s *UnitTestSuite) Test_CreateBill_WorkflowCreationFail() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	mockClient.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	ctx := context.Background()
	req := &CreateBillRequest{
		Currency: "USD",
	}

	resp, err := service.CreateBill(ctx, req)
	s.Error(err)
	s.Nil(resp)
}

func (s *UnitTestSuite) Test_CreateBill_InvalidCurrency() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()
	req := &CreateBillRequest{
		Currency: "EUR",
	}

	resp, err := service.CreateBill(ctx, req)
	s.Error(err)
	s.EqualError(err, "invalid_argument: unsupported currency, only USD or GEL")
	s.Nil(resp)
}

func (s *UnitTestSuite) Test_CloseBill_Success() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockClient.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockEncodedValue := &MockEncodedValue{}
	mockEncodedValue.On("Get").Return(nil)
	mockClient.On("QueryWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockEncodedValue, nil)

	req := &CloseBillRequest{
		Id: "1234",
	}

	resp, err := service.CloseBill(ctx, req)
	s.NoError(err)
	s.Equal("1234", resp.Id)
}

func (s *UnitTestSuite) Test_CloseBill_SignalFail() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockClient.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

	req := &CloseBillRequest{
		Id: "1234",
	}

	resp, err := service.CloseBill(ctx, req)
	s.Error(err)
	s.Nil(resp)
}

func (s *UnitTestSuite) Test_GetBill_Success() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockEncodedValue := &MockEncodedValue{}
	mockEncodedValue.On("Get").Return(nil)
	mockClient.On("QueryWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockEncodedValue, nil)

	bill, err := service.GetBill(ctx, "1234")
	s.NoError(err)
	s.Equal(mockBill.Currency, bill.Currency)
	s.Equal(mockBill.TotalAmount, bill.TotalAmount)
	s.Empty(bill.LineItems)
}

func (s *UnitTestSuite) Test_GetBill_Fail() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockClient.On("QueryWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	bill, err := service.GetBill(ctx, "1234")
	s.Error(err)
	s.Nil(bill)
}

func (s *UnitTestSuite) Test_AddLineItem_Success() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockClient.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockEncodedValue := &MockEncodedValue{}
	mockEncodedValue.On("Get").Return(nil)
	mockClient.On("QueryWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockEncodedValue, nil)

	req := &AddLineItemRequest{
		BillId: "1234",
		Description: "item1",
		Amount: 10.0,
	}

	resp, err := service.AddLineItem(ctx, req)
	s.NoError(err)
	s.Equal(1.0, resp.CurrentTotal)
	s.Equal(0, resp.NumberOfItems)
}

func (s *UnitTestSuite) Test_AddLineItem_SignalFail() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockClient.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

	req := &AddLineItemRequest{
		BillId: "1234",
		Description: "item1",
		Amount: 10.0,
	}

	resp, err := service.AddLineItem(ctx, req)
	s.Error(err)
	s.Equal(err.Error(), "internal: unable to add line item to bill")
	s.Nil(resp)
}

func (s *UnitTestSuite) Test_AddLineItem_QueryFail() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	mockClient.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockClient.On("QueryWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	req := &AddLineItemRequest{
		BillId: "1234",
		Description: "item1",
		Amount: 10.0,
	}

	resp, err := service.AddLineItem(ctx, req)
	s.Error(err)
	s.Equal(err.Error(), "internal: unable to get bill")
	s.Nil(resp)
}

func (s *UnitTestSuite) Test_AddLineItem_InvalidAmount() {
	mockClient := mocks.NewClient(s.T())
	service := &Service{
		client: mockClient,
		worker: nil,
		eb:     *errs.B(),
	}

	ctx := context.Background()

	req := &AddLineItemRequest{
		BillId: "1234",
		Description: "item1",
		Amount: -10.0,
	}

	resp, err := service.AddLineItem(ctx, req)
	s.Error(err)
	s.Equal(err.Error(), "invalid_argument: amount must be greater than 0")
	s.Nil(resp)
}