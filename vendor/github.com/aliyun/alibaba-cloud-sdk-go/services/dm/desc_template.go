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

// DescTemplate invokes the dm.DescTemplate API synchronously
// api document: https://help.aliyun.com/api/dm/desctemplate.html
func (client *Client) DescTemplate(request *DescTemplateRequest) (response *DescTemplateResponse, err error) {
	response = CreateDescTemplateResponse()
	err = client.DoAction(request, response)
	return
}

// DescTemplateWithChan invokes the dm.DescTemplate API asynchronously
// api document: https://help.aliyun.com/api/dm/desctemplate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescTemplateWithChan(request *DescTemplateRequest) (<-chan *DescTemplateResponse, <-chan error) {
	responseChan := make(chan *DescTemplateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescTemplate(request)
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

// DescTemplateWithCallback invokes the dm.DescTemplate API asynchronously
// api document: https://help.aliyun.com/api/dm/desctemplate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescTemplateWithCallback(request *DescTemplateRequest, callback func(response *DescTemplateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescTemplateResponse
		var err error
		defer close(result)
		response, err = client.DescTemplate(request)
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

// DescTemplateRequest is the request struct for api DescTemplate
type DescTemplateRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	TemplateId           requests.Integer `position:"Query" name:"TemplateId"`
	FromType             requests.Integer `position:"Query" name:"FromType"`
}

// DescTemplateResponse is the response struct for api DescTemplate
type DescTemplateResponse struct {
	*responses.BaseResponse
	RequestId        string `json:"RequestId" xml:"RequestId"`
	TemplateName     string `json:"TemplateName" xml:"TemplateName"`
	TemplateSubject  string `json:"TemplateSubject" xml:"TemplateSubject"`
	TemplateNickName string `json:"TemplateNickName" xml:"TemplateNickName"`
	TemplateStatus   string `json:"TemplateStatus" xml:"TemplateStatus"`
	TemplateType     string `json:"TemplateType" xml:"TemplateType"`
	CreateTime       string `json:"CreateTime" xml:"CreateTime"`
	TemplateText     string `json:"TemplateText" xml:"TemplateText"`
	SmsContent       string `json:"SmsContent" xml:"SmsContent"`
	SmsType          string `json:"SmsType" xml:"SmsType"`
	Remark           string `json:"Remark" xml:"Remark"`
}

// CreateDescTemplateRequest creates a request to invoke DescTemplate API
func CreateDescTemplateRequest() (request *DescTemplateRequest) {
	request = &DescTemplateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dm", "2015-11-23", "DescTemplate", "", "")
	return
}

// CreateDescTemplateResponse creates a response to parse from DescTemplate response
func CreateDescTemplateResponse() (response *DescTemplateResponse) {
	response = &DescTemplateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
