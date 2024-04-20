// Package ocpp implements an Open Charge Point Protocol (OCPP OCPP 1.6) client for interacting with charge points.
//
// This package provides functionality to communicate with OCPP-compliant charge points,
// including methods for BootNotification, Heartbeat, Authorization, Transaction management,
// and other OCPP-defined actions.
// Author: Mukul kumar(https://github.com/slimdestro/)

package ocpp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	// DefaultHTTPClient can be overridden for custom HTTP client configurations.
	DefaultHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// Client represents an OCPP client that communicates with a charge point.
type Client struct {
	endpoint string
	client   *http.Client
	logger   *zap.Logger
}

// NewClient creates a new OCPP client with the specified endpoint URL and optional logger.
func NewClient(endpoint string, logger *zap.Logger) *Client {
	if logger == nil {
		logger = zap.NewNop() // Default to no-op logger if not provided
	}
	return &Client{
		endpoint: endpoint,
		client:   DefaultHTTPClient,
		logger:   logger,
	}
}

// SetHTTPClient allows setting a custom HTTP client for the client.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

// sendRequest sends a request to the OCPP server and parses the response.
func (c *Client) sendRequest(action string, request interface{}, response interface{}) error {
	requestXML, err := xml.Marshal(request)
	if err != nil {
		c.logger.Error("failed to marshal request XML", zap.Error(err))
		return err
	}

	url := fmt.Sprintf("%s/%s", c.endpoint, action)
	resp, err := c.client.Post(url, "application/xml", bytes.NewBuffer(requestXML))
	if err != nil {
		c.logger.Error("failed to send HTTP request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("non-200 status code received: %s", resp.Status)
		c.logger.Error(errMsg)
		return errors.New(errMsg)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("failed to read response body", zap.Error(err))
		return err
	}

	if err := xml.Unmarshal(respBody, response); err != nil {
		c.logger.Error("failed to unmarshal response XML", zap.Error(err))
		return err
	}

	return nil
}

// BootNotification sends a BootNotification request to the charge point.
func (c *Client) BootNotification(chargeBoxIdentity string) (*BootNotificationResponse, error) {
	request := BootNotificationRequest{
		ChargeBoxIdentity: chargeBoxIdentity,
	}
	response := &BootNotificationResponse{}
	if err := c.sendRequest("BootNotification", request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// Heartbeat sends a Heartbeat request to the charge point.
func (c *Client) Heartbeat() (*HeartbeatResponse, error) {
	request := HeartbeatRequest{}
	response := &HeartbeatResponse{}
	if err := c.sendRequest("Heartbeat", request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// Authorize sends an Authorize request to the charge point.
func (c *Client) Authorize(idTag string) (*AuthorizeResponse, error) {
	request := AuthorizeRequest{
		IdTag: idTag,
	}
	response := &AuthorizeResponse{}
	if err := c.sendRequest("Authorize", request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// StartTransaction sends a StartTransaction request to the charge point.
func (c *Client) StartTransaction(connectorId int, idTag string) (*StartTransactionResponse, error) {
	request := StartTransactionRequest{
		ConnectorId: connectorId,
		IdTag:       idTag,
	}
	response := &StartTransactionResponse{}
	if err := c.sendRequest("StartTransaction", request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// StopTransaction sends a StopTransaction request to the charge point.
func (c *Client) StopTransaction(transactionId int) (*StopTransactionResponse, error) {
	request := StopTransactionRequest{
		TransactionId: transactionId,
	}
	response := &StopTransactionResponse{}
	if err := c.sendRequest("StopTransaction", request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// MeterValues sends a MeterValues request to the charge point.
func (c *Client) MeterValues(values []MeterValue) error {
	request := MeterValuesRequest{
		Values: values,
	}
	response := &MeterValuesResponse{}
	if err := c.sendRequest("MeterValues", request, response); err != nil {
		return err
	}
	return nil
}

// StatusNotification sends a StatusNotification request to the charge point.
func (c *Client) StatusNotification(status Status) error {
	request := StatusNotificationRequest{
		Status: status,
	}
	response := &StatusNotificationResponse{}
	if err := c.sendRequest("StatusNotification", request, response); err != nil {
		return err
	}
	return nil
}

// DataTransfer sends a DataTransfer request to the charge point.
func (c *Client) DataTransfer(vendorId string, messageData string) (*DataTransferResponse, error) {
	request := DataTransferRequest{
		VendorId:    vendorId,
		MessageData: messageData,
	}
	response := &DataTransferResponse{}
	if err := c.sendRequest("DataTransfer", request, response); err != nil {
		return nil, err
	}
	return response, nil
}

// Define request and response structs for OCPP actions here.

// MeterValue represents a meter value in OCPP.
type MeterValue struct {
	// Define fields for MeterValue here
}

// BootNotificationRequest represents a BootNotification request in OCPP.
type BootNotificationRequest struct {
	ChargeBoxIdentity string `xml:"chargeBoxIdentity"`
}

// BootNotificationResponse represents a BootNotification response in OCPP.
type BootNotificationResponse struct {
	Status      RegistrationStatus `xml:"status"`
	CurrentTime string             `xml:"currentTime,omitempty"`
	Interval    int                `xml:"interval,omitempty"`
	Heartbeat   int                `xml:"heartbeat,omitempty"`
}

// HeartbeatRequest represents a Heartbeat request in OCPP.
type HeartbeatRequest struct{}

// HeartbeatResponse represents a Heartbeat response in OCPP.
type HeartbeatResponse struct {
	CurrentTime string `xml:"currentTime,omitempty"`
}

// AuthorizeRequest represents an Authorize request in OCPP.
type AuthorizeRequest struct {
	IdTag string `xml:"idTag"`
}

// AuthorizeResponse represents an Authorize response in OCPP.
type AuthorizeResponse struct {
	IdTagInfo IdTagInfo `xml:"idTagInfo"`
}

// IdTagInfo represents information about an ID tag in OCPP.
type IdTagInfo struct {
	Status      AuthorizationStatus `xml:"status"`
	ExpiryDate  *string             `xml:"expiryDate,omitempty"`
	ParentIdTag *string             `xml:"parentIdTag,omitempty"`
}

// StartTransactionRequest represents a StartTransaction request in OCPP.
type StartTransactionRequest struct {
	ConnectorId int    `xml:"connectorId"`
	IdTag       string `xml:"idTag"`
}

// StartTransactionResponse represents a StartTransaction response in OCPP.
type StartTransactionResponse struct {
	TransactionId int `xml:"transactionId"`
}

// StopTransactionRequest represents a StopTransaction request in OCPP.
type StopTransactionRequest struct {
	TransactionId int `xml:"transactionId"`
}

// StopTransactionResponse represents a StopTransaction response in OCPP.
type StopTransactionResponse struct {
	Status TransactionEventStatus `xml:"status"`
}

type MeterValuesResponse struct {
	Status RegistrationStatus `xml:"status"`
}

// StatusNotificationRequest represents a StatusNotification request in OCPP.
type StatusNotificationRequest struct {
	Status Status `xml:"status"`
}

// ErrorCode represents an error code in OCPP.
type ErrorCode string

const (
	// Common OCPP error codes (example)
	ErrorCodeGenericError       ErrorCode = "GenericError"
	ErrorCodeProtocolError      ErrorCode = "ProtocolError"
	ErrorCodeInternalError      ErrorCode = "InternalError"
	ErrorCodeNotImplemented     ErrorCode = "NotImplemented"
	ErrorCodeUnknownMessageType ErrorCode = "UnknownMessageType"
	// Add more error codes as needed
)

// Status represents the status information in OCPP.
type Status struct {
	ErrorCode    ErrorCode `xml:"errorCode,omitempty"`
	StatusDetail string    `xml:"statusDetail,omitempty"`
}

type StatusNotificationResponse struct {
	Status Status `xml:"status"`
}

// DataTransferRequest represents a DataTransfer request in OCPP.
type DataTransferRequest struct {
	VendorId    string `xml:"vendorId"`
	MessageData string `xml:"messageData"`
}

// DataTransferResponse represents a DataTransfer response in OCPP.
type DataTransferResponse struct {
	Status DataTransferStatus `xml:"status"`
	Data   string             `xml:"data,omitempty"`
}

// DataTransferStatus represents the status of a DataTransfer operation.
type DataTransferStatus string

const (
	DataTransferStatusAccepted DataTransferStatus = "Accepted"
	DataTransferStatusRejected DataTransferStatus = "Rejected"
)

// RegistrationStatus represents the registration status of a charge point.
type RegistrationStatus string

const (
	RegistrationStatusAccepted       RegistrationStatus = "Accepted"
	RegistrationStatusPending        RegistrationStatus = "Pending"
	RegistrationStatusRejected       RegistrationStatus = "Rejected"
	RegistrationStatusScheduled      RegistrationStatus = "Scheduled"
	RegistrationStatusUnscheduled    RegistrationStatus = "Unscheduled"
	RegistrationStatusRecurring      RegistrationStatus = "Recurring"
	RegistrationStatusCancelled      RegistrationStatus = "Cancelled"
	RegistrationStatusInstallation   RegistrationStatus = "Installation"
	RegistrationStatusRegistration   RegistrationStatus = "Registration"
	RegistrationStatusDeregistration RegistrationStatus = "Deregistration"
)

// AuthorizationStatus represents the authorization status of an ID tag.
type AuthorizationStatus string

const (
	AuthorizationStatusAccepted   AuthorizationStatus = "Accepted"
	AuthorizationStatusBlocked    AuthorizationStatus = "Blocked"
	AuthorizationStatusExpired    AuthorizationStatus = "Expired"
	AuthorizationStatusInvalid    AuthorizationStatus = "Invalid"
	AuthorizationStatusConcurrent AuthorizationStatus = "ConcurrentTx"
)

// TransactionEventStatus represents the status of a transaction event.
type TransactionEventStatus string

const (
	TransactionEventStatusAccepted TransactionEventStatus = "Accepted"
	TransactionEventStatusRejected TransactionEventStatus = "Rejected"
)
