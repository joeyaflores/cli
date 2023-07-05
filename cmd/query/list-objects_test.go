package query

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_client "github.com/openfga/cli/mocks"
	"github.com/openfga/go-sdk/client"
)

var errMockListObjects = errors.New("mock error")

func TestListObjectsWithError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFgaClient := mock_client.NewMockSdkClient(mockCtrl)

	mockExecute := mock_client.NewMockSdkClientListObjectsRequestInterface(mockCtrl)

	var expectedResponse client.ClientListObjectsResponse

	mockExecute.EXPECT().Execute().Return(&expectedResponse, errMockListObjects)

	mockRequest := mock_client.NewMockSdkClientListObjectsRequestInterface(mockCtrl)
	options := client.ClientListObjectsOptions{}
	mockRequest.EXPECT().Options(options).Return(mockExecute)

	mockBody := mock_client.NewMockSdkClientListObjectsRequestInterface(mockCtrl)

	body := client.ClientListObjectsRequest{
		User:     "user:foo",
		Relation: "writer",
		Type:     "doc",
	}
	mockBody.EXPECT().Body(body).Return(mockRequest)

	mockFgaClient.EXPECT().ListObjects(context.Background()).Return(mockBody)

	_, err := listObjects(mockFgaClient, "user:foo", "writer", "doc")
	if err == nil {
		t.Error("Expect error but there is none")
	}
}

func TestListObjectsWithNoError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockFgaClient := mock_client.NewMockSdkClient(mockCtrl)

	mockExecute := mock_client.NewMockSdkClientListObjectsRequestInterface(mockCtrl)

	expectedResponse := client.ClientListObjectsResponse{
		Objects: &[]string{"doc:doc1", "doc:doc2"},
	}

	mockExecute.EXPECT().Execute().Return(&expectedResponse, nil)

	mockRequest := mock_client.NewMockSdkClientListObjectsRequestInterface(mockCtrl)
	options := client.ClientListObjectsOptions{}
	mockRequest.EXPECT().Options(options).Return(mockExecute)

	mockBody := mock_client.NewMockSdkClientListObjectsRequestInterface(mockCtrl)

	body := client.ClientListObjectsRequest{
		User:     "user:foo",
		Relation: "writer",
		Type:     "doc",
	}
	mockBody.EXPECT().Body(body).Return(mockRequest)

	mockFgaClient.EXPECT().ListObjects(context.Background()).Return(mockBody)

	output, err := listObjects(mockFgaClient, "user:foo", "writer", "doc")
	if err != nil {
		t.Error(err)
	}

	expectedOutput := "{\"objects\":[\"doc:doc1\",\"doc:doc2\"]}"
	if output != expectedOutput {
		t.Errorf("Expect %v but actual %v", expectedOutput, output)
	}
}