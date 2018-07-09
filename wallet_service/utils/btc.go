package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"io/ioutil"
	"net/http"
)

//var btcClient *rpcclient.Client
type BtcClient struct {
}

// btc 客户端链接
func (p *BtcClient) NewClient(host, user, pass string) (btcclient *rpcclient.Client, err error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	btcclient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		Log.Errorln(err.Error())
	}
	return
}

/*
 curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "sendtoaddress", "params": ["1M72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd", 0.1, "donation", "seans outpost"] }' -H 'content-type: text/plain;' http://user:pass@127.0.0.1:8332/
*/
func BtcSendToAddressFunc(url string, address string, mount string) (string, error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "sendtoaddress"
	//param := []string{}
	params := make([]interface{}, 0, 2)
	params = append(params, address)
	params = append(params, mount)
	data["params"] = params

	rsp, err := BtcRpcPost(url, data)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(rsp, &ret)
	if err != nil {
		fmt.Println("unmarshal error: ", err)
	}
	result, ok := ret["result"]
	if !ok {
		return "", err
	}
	txHash, ok := result.(string)
	if !ok {
		msg := "sendtoaddress error!"
		err = errors.New(msg)
		Log.Errorln(msg)
		return "", err
	}
	return txHash, nil
}

/*
  curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "walletpassphrase", "params": ["my pass phrase", 60] }' -H 'content-type: text/plain;' http://user@pass127.0.0.1:18332/
  解锁钱包
*/
func BtcWalletPhrase(url string, pass string, keepTime int64) error {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "walletpassphrase"
	params := make([]interface{}, 0, 2)
	params = append(params, pass)
	params = append(params, keepTime)
	data["params"] = params
	_, err := BtcRpcPost(url, data)
	if err != nil {
		fmt.Println(err.Error())
		Log.Errorf(err.Error())
		return err
	}
	return nil
}

/*
  curl  --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "listunspent", "params": [6, 9999999, [] , true, { "minimumAmount": 0.005 } ] }' -H 'content-type: text/plain;' http://user:pass@127.0.0.1:8332/
*/
func BtcListunspent(url string, start int, end int) (string, error) {
	data := make(map[string]interface{})
	data["jsonrpc"] = "1.0"
	data["id"] = 1
	data["method"] = "listunspent"

	if start == 0 {
		start = 0
	}
	if end == 0 {
		end = 5
	}
	param := []int{}
	param = append(param, start)
	param = append(param, end)
	data["params"] = param
	rsp, err := BtcRpcPost(url, data)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(rsp, &ret)
	result, ok := ret["result"]
	if !ok {
		fmt.Println("result:", result)
		return "", err
	}
	//fmt.Println("result:", result)
	return "", nil
}

/*

 */
func BtcRpcPost(url string, send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println("rpc post:", err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	//fmt.Println("resp:", resp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//byte数组直接转成string，优化内存
	return respBytes, nil
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println(*str)
}