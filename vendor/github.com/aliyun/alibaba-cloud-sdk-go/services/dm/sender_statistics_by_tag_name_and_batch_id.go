package dm

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// SenderStatisticsByTagNameAndBatchID invokes the dm.SenderStatisticsByTagNameAndBatchID API synchronously
// api document: https://help.aliyun.com/api/dm/senderstatisticsbytagnameandbatchid.html
func (client *Client) SenderStatisticsByTagNameAndBatchID(request *SenderStatisticsByTagNameAndBatchIDRequest) (response *SenderStatisticsByTagNameAndBatchIDResponse, err error) {
	response = CreateSenderStatisticsByTagNameAndBatchIDResponse()
	err = client.DoAction(request, response)
	return
}

// SenderStatisticsByTagNameAndBatchIDWithChan invokes the dm.SenderStatisticsByTagNameAndBatchID API asynchronously
// api document: https://help.aliyun.com/api/dm/senderstatisticsbytagnameandbatchid.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SenderStatisticsByTagNameAndBatchIDWithChan(request *SenderStatisticsByTagNameAndBatchIDRequest) (<-chan *SenderStatisticsByTagNameAndBatchIDResponse, <-chan error) {
	responseChan := make(chan *SenderStatisticsByTagNameAndBatchIDResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SenderStatisticsByTagNameAndBatchID(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// SenderStatisticsByTagNameAndBatchIDWithCallback invokes the dm.SenderStatisticsByTagNameAndBatchID API asynchronously
// api document: https://help.aliyun.com/api/dm/senderstatisticsbytagnameandbatchid.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SenderStatisticsByTagNameAndBatchIDWithCallback(request *SenderStatisticsByTagNameAndBatchIDRequest, callback func(response *SenderStatisticsByTagNameAndBatchIDResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SenderStatisticsByTagNameAndBatchIDResponse
		var err error
		defer close(result)
		response, err = client.SenderStatisticsByTagNameAndBatchID(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// SenderStatisticsByTagNameAndBatchIDRequest is the request struct for api SenderStatisticsByTagNameAndBatchID
type SenderStatisticsByTagNameAndBatchIDRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	AccountName          string           `position:"Query" name:"AccountName"`
	StartTime            string           `position:"Query" name:"StartTime"`
	EndTime              string           `position:"Query" name:"EndTime"`
	TagName              string           `position:"Query" name:"TagName"`
}

// SenderStatisticsByTagNameAndBatchIDResponse is the response struct for api SenderStatisticsByTagNameAndBatchID
type SenderStatisticsByTagNameAndBatchIDResponse struct {
	*responses.BaseResponse
	RequestId  string                                    `json:"RequestId" xml:"RequestId"`
	TotalCount int                                       `json:"TotalCount" xml:"TotalCount"`
	Data       DataInSenderStatisticsByTagNameAndBatchID `json:"data" xml:"data"`
}

// CreateSenderStatisticsByTagNameAndBatchIDRequest creates a request to invoke SenderStatisticsByTagNameAndBatchID API
func CreateSenderStatisticsByTagNameAndBatchIDRequest() (request *SenderStatisticsByTagNameAndBatchIDRequest) {
	request = &SenderStatisticsByTagNameAndBatchIDRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dm", "2015-11-23", "SenderStatisticsByTagNameAndBatchID", "", "")
	return
}

// CreateSenderStatisticsByTagNameAndBatchIDResponse creates a response to parse from SenderStatisticsByTagNameAndBatchID response
func CreateSenderStatisticsByTagNameAndBatchIDResponse() (response *SenderStatisticsByTagNameAndBatchIDResponse) {
	response = &SenderStatisticsByTagNameAndBatchIDResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
