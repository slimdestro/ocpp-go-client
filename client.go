// Client package for Open Charge Point Protocol (OCPP)
// Author: Mukul kumar(https://github.com/slimdestro/)

package ocpp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Client struct {
	endpoint string
	mutex    sync.Mutex
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func (c *Client) sendRequest(action string, request interface{}, response interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	requestXML, err := xml.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", c.endpoint, action)
	resp, err := http.Post(url, "application/xml", strings.NewReader(xml.Header+string(requestXML)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("non-200 status code received")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(respBody, response)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) BootNotification(chargeBoxIdentity string) error {
	request := BootNotificationRequest{
		ChargeBoxIdentity: chargeBoxIdentity,
	}
	var response BootNotificationResponse
	return c.sendRequest("BootNotification", request, &response)
}

func (c *Client) Heartbeat() error {
	request := HeartbeatRequest{}
	var response HeartbeatResponse
	return c.sendRequest("Heartbeat", request, &response)
}

func (c *Client) Authorize(idTag string) (AuthorizeResponse, error) {
	request := AuthorizeRequest{
		IdTag: idTag,
	}
	var response AuthorizeResponse
	err := c.sendRequest("Authorize", request, &response)
	return response, err
}

func (c *Client) StartTransaction(connectorId int, idTag string) (StartTransactionResponse, error) {
	request := StartTransactionRequest{
		ConnectorId: connectorId,
		IdTag:       idTag,
	}
	var response StartTransactionResponse
	err := c.sendRequest("StartTransaction", request, &response)
	return response, err
}

func (c *Client) StopTransaction(transactionId int) (StopTransactionResponse, error) {
	request := StopTransactionRequest{
		TransactionId: transactionId,
	}
	var response StopTransactionResponse
	err := c.sendRequest("StopTransaction", request, &response)
	return response, err
}

func (c *Client) MeterValues(values []MeterValue) error {
	request := MeterValuesRequest{
		Values: values,
	}
	var response MeterValuesResponse
	return c.sendRequest("MeterValues", request, &response)
}

func (c *Client) StatusNotification(status Status) error {
	request := StatusNotificationRequest{
		Status: status,
	}
	var response StatusNotificationResponse
	return c.sendRequest("StatusNotification", request, &response)
}

func (c *Client) DataTransfer(vendorId string, messageData string) (DataTransferResponse, error) {
	request := DataTransferRequest{
		VendorId:    vendorId,
		MessageData: messageData,
	}
	var response DataTransferResponse
	err := c.sendRequest("DataTransfer", request, &response)
	return response, err
}
