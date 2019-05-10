package wallet

import (
	"fmt"

	"github.com/gorilla/rpc/v2/json2"
)

// H is a helper map shortcut.
type H map[string]interface{}

// ErrorCode is a monero-wallet-rpc error code.
// Copied from https://github.com/monero-project/monero/blob/release-v0.11.0.0/src/wallet/wallet_rpc_server_error_codes.h
type ErrorCode int

const (
	// ErrUnknown - E_UNKNOWN_ERROR
	ErrUnknown ErrorCode = -1
	// ErrWrongAddress - E_WRONG_ADDRESS
	ErrWrongAddress ErrorCode = -2
	// ErrDaemonIsBusy - E_DAEMON_IS_BUSY
	ErrDaemonIsBusy ErrorCode = -3
	// ErrGenericTransferError - E_GENERIC_TRANSFER_ERROR
	ErrGenericTransferError ErrorCode = -4
	// ErrWrongPaymentID - E_WRONG_PAYMENT_ID
	ErrWrongPaymentID ErrorCode = -5
	// ErrTransferType - E_TRANSFER_TYPE
	ErrTransferType ErrorCode = -6
	// ErrDenied - E_DENIED
	ErrDenied ErrorCode = -7
	// ErrWrongTxID - E_WRONG_TXID
	ErrWrongTxID ErrorCode = -8
	// ErrWrongSignature - E_WRONG_SIGNATURE
	ErrWrongSignature ErrorCode = -9
	// ErrWrongKeyImage - E_WRONG_KEY_IMAGE
	ErrWrongKeyImage ErrorCode = -10
	// ErrWrongURI - E_WRONG_URI
	ErrWrongURI ErrorCode = -11
	// ErrWrongIndex - E_WRONG_INDEX
	ErrWrongIndex ErrorCode = -12
	// ErrNotOpen - E_NOT_OPEN
	ErrNotOpen ErrorCode = -13
)

// WalletError is the error structured returned by the monero-wallet-rpc
type WalletError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (we *WalletError) Error() string {
	return fmt.Sprintf("%v: %v", we.Code, we.Message)
}

// GetWalletError checks if an erro interface is a wallet-rpc error.
func GetWalletError(err error) (isWalletError bool, werr *WalletError) {
	if err == nil {
		return false, nil
	}
	gerr, ok := err.(*json2.Error)
	if !ok {
		return false, nil
	}
	werr = &WalletError{
		Code:    ErrorCode(gerr.Code),
		Message: gerr.Message,
	}
	isWalletError = true
	return
}

// Priority represents a transaction priority
type Priority uint

// Accepted Values are: 0-3 for: default, unimportant, normal, elevated, priority.
const (
	PriorityDefault     Priority = 0
	PriorityUnimportant Priority = 1
	PriorityNormal      Priority = 2
	PriorityElevated    Priority = 3
)

// GetTransferType is a string that contains the possible types:
// "all": all the transfers;
// "available": only transfers which are not yet spent;
// "unavailable": only transfers which are already spent.
type GetTransferType string

const (
	// TransferAll - all the transfers
	TransferAll GetTransferType = "all"
	// TransferAvailable - only transfers which are not yet spent
	TransferAvailable GetTransferType = "available"
	// TransferUnavailable - only transfers which are already spent
	TransferUnavailable GetTransferType = "unavailable"
)

// QueryKeyType is the parameter to send with client.QueryKey()
type QueryKeyType string

const (
	// QueryKeyMnemonic is the mnemonic seed
	QueryKeyMnemonic QueryKeyType = "mnemonic"
	// QueryKeyView is the private view key
	QueryKeyView QueryKeyType = "view_key"
	// QueryKeySpend is the private spend key
	QueryKeySpend QueryKeyType = "spend_key" //TODO: test
)
