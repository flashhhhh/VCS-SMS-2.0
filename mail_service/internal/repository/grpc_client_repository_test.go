package repository_test

import (
	"errors"
	"testing"

	"context"
	"mail_service/internal/repository"
	"mail_service/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMailGRPCClient implements grpcclient.MailGRPCClient for testing
type mockMailGRPCClient struct {
	mock.Mock
}

func (m *mockMailGRPCClient) GetServersInformation(ctx context.Context, req *proto.TimeRequest) (*proto.ServersInformationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*proto.ServersInformationResponse), args.Error(1)
}

func TestMailGRPCClientRepository_GetServersInformation_Success(t *testing.T) {
	mockClient := new(mockMailGRPCClient)
	repo := repository.NewMailGRPCClientRepository(mockClient)

	startTime := "2024-01-01T00:00:00Z"
	endTime := "2024-01-02T00:00:00Z"
	expectedResp := &proto.ServersInformationResponse{
		NumServers:      10,
		NumOnServers:    7,
		NumOffServers:   3,
		MeanUpTimeRatio: 0.85,
	}

	mockClient.
		On("GetServersInformation", mock.Anything, &proto.TimeRequest{
			StartTime: startTime,
			EndTime:   endTime,
		}).
		Return(expectedResp, nil).
		Once()

	numServers, numOn, numOff, meanRatio, err := repo.GetServersInformation(startTime, endTime)
	assert.NoError(t, err)
	assert.Equal(t, 10, numServers)
	assert.Equal(t, 7, numOn)
	assert.Equal(t, 3, numOff)
	assert.Equal(t, 0.85, meanRatio)
	mockClient.AssertExpectations(t)
}

func TestMailGRPCClientRepository_GetServersInformation_Error(t *testing.T) {
	mockClient := new(mockMailGRPCClient)
	repo := repository.NewMailGRPCClientRepository(mockClient)

	startTime := "2024-01-01T00:00:00Z"
	endTime := "2024-01-02T00:00:00Z"
	expectedErr := errors.New("grpc error")

	mockClient.
		On("GetServersInformation", mock.Anything, &proto.TimeRequest{
			StartTime: startTime,
			EndTime:   endTime,
		}).
		Return(nil, expectedErr).
		Once()

	numServers, numOn, numOff, meanRatio, err := repo.GetServersInformation(startTime, endTime)
	assert.Error(t, err)
	assert.Equal(t, 0, numServers)
	assert.Equal(t, 0, numOn)
	assert.Equal(t, 0, numOff)
	assert.Equal(t, 0.0, meanRatio)
	assert.Equal(t, expectedErr, err)
	mockClient.AssertExpectations(t)
}