package wallet

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// Client is a monero-wallet-rpc client.
type Client interface {
	// Return the wallet's balance.
	GetBalance(*RequestGetBalance) (*ResponseGetBalance, error)
	// Return the wallet's addresses for an account. Optionally filter for specific set of subaddresses.
	GetAddress(*RequestGetAddress) (*ResponseGetAddress, error)
	// Get account and address indexes from a specific (sub)address
	GetAddressIndex(*RequestGetAddressIndex) (*ResponseGetAddressIndex, error)
	// Create a new address for an account. Optionally, label the new address.
	CreateAddress(*RequestCreateAddress) (*ResponseCreateAddress, error)
	// Label an address.
	LabelAddress(*RequestLabelAddress) error
	// Validate an address.
	ValidateAddress(*RequestValidateAddress) (*ResponseValidateAddress, error)
	// Get all accounts for a wallet. Optionally filter accounts by tag.
	GetAccounts(*RequestGetAccounts) (*ResponseGetAccounts, error)
	// Create a new account with an optional label.
	CreateAccount(*RequestCreateAccount) (*ResponseCreateAccount, error)
	// Label an account.
	LabelAccount(*RequestLabelAccount) error
	// Get a list of user-defined account tags.
	GetAccountTags() (*ResponseGetAccountTags, error)
	// Apply a filtering tag to a list of accounts.
	TagAccounts(*RequestTagAccounts) error
	// Remove filtering tag from a list of accounts.
	UntagAccounts(*RequestUntagAccounts) error
	// Set description for an account tag.
	SetAccountTagDescription(*RequestSetAccountTagDescription) error
	// Returns the wallet's current block height.
	GetHeight() (*ResponseGetHeight, error)
	// Send monero to a number of recipients.
	Transfer(*RequestTransfer) (*ResponseTransfer, error)
	// Same as transfer, but can split into more than one tx if necessary.
	TransferSplit(*RequestTransferSplit) (*ResponseTransferSplit, error)
	// Sign a transaction created on a read-only wallet (in cold-signing process)
	SignTransfer(*RequestSignTransfer) (*ResponseSignTransfer, error)
	// Submit a previously signed transaction on a read-only wallet (in cold-signing process).
	SubmitTransfer(*RequestSubmitTransfer) (*ResponseSubmitTransfer, error)
	// Send all dust outputs back to the wallet's, to make them easier to spend (and mix).
	SweepDust(*RequestSweepDust) (*ResponseSweepDust, error)
	// Send all unlocked balance to an address.
	SweepAll(*RequestSweepAll) (*ResponseSweepAll, error)
	// Send all of a specific unlocked output to an address.
	SweepSingle(*RequestSweepSingle) (*ResponseSweepSingle, error)
	// Relay a transaction previously created with "do_not_relay":true.
	RelayTx(*RequestRelayTx) (*ResponseRelayTx, error)
	// Save the wallet file.
	Store() error
	// Get a list of incoming payments using a given payment id.
	GetPayments(*RequestGetPayments) (*ResponseGetPayments, error)
	// Get a list of incoming payments using a given payment id, or a list of payments ids, from a given height.
	// This method is the preferred method over get_payments because it has the same functionality but is more extendable.
	// Either is fine for looking up transactions by a single payment ID.
	GetBulkPayments(*RequestGetBulkPayments) (*ResponseGetBulkPayments, error)
	// Return a list of incoming transfers to the wallet.
	IncomingTransfers(*RequestIncomingTransfers) (*ResponseIncomingTransfers, error)
	// Return the spend or view private key.
	QueryKey(*RequestQueryKey) (*ResponseQueryKey, error)
	// Make an integrated address from the wallet address and a payment id.
	MakeIntegratedAddress(*RequestMakeIntegratedAddress) (*ResponseMakeIntegratedAddress, error)
	// Retrieve the standard address and payment id corresponding to an integrated address.
	SplitIntegratedAddress(*RequestSplitIntegratedAddress) (*ResponseSplitIntegratedAddress, error)
	// Stops the wallet, storing the current state.
	StopWallet() error
	// Rescan the blockchain from scratch, losing any information which can not be recovered from the blockchain itself.
	// This includes destination addresses, tx secret keys, tx notes, etc.
	RescanBlockchain() error
	// Set arbitrary string notes for transactions.
	SetTxNotes(*RequestSetTxNotes) error
	// Get string notes for transactions.
	GetTxNotes(*RequestGetTxNotes) (*ResponseGetTxNotes, error)
	// Set arbitrary attribute.
	SetAttribute(*RequestSetAttribute) error
	// Get attribute value by name.
	GetAttribute(*RequestGetAttribute) (*ResponseGetAttribute, error)
	// Get transaction secret key from transaction id.
	GetTxKey(*RequestGetTxKey) (*ResponseGetTxKey, error)
	// Check a transaction in the blockchain with its secret key.
	CheckTxKey(*RequestCheckTxKey) (*ResponseCheckTxKey, error)
	// Get transaction signature to prove it.
	GetTxProof(*RequestGetTxProof) (*ResponseGetTxProof, error)
	// Prove a transaction by checking its signature.
	CheckTxProof(*RequestCheckTxProof) (*ResponseCheckTxProof, error)
	// Generate a signature to prove a spend. Unlike proving a transaction, it does not requires the destination public address.
	GetSpendProof(*RequestGetSpendProof) (*ResponseGetSpendProof, error)
	// Prove a spend using a signature. Unlike proving a transaction, it does not requires the destination public address.
	CheckSpendProof(*RequestCheckSpendProof) (*ResponseCheckSpendProof, error)
	// Generate a signature to prove of an available amount in a wallet.
	GetReserveProof(*RequestGetReserveProof) (*ResponseGetReserveProof, error)
	// Proves a wallet has a disposable reserve using a signature.
	CheckReserveProof(*RequestCheckReserveProof) (*ResponseCheckReserveProof, error)
	// Returns a list of transfers.
	GetTransfers(*RequestGetTransfers) (*ResponseGetTransfers, error)
	// Show information about a transfer to/from this address.
	GetTransferByTxID(*RequestGetTransferByTxID) (*ResponseGetTransferByTxID, error)
	// Sign a string.
	Sign(*RequestSign) (*ResponseSign, error)
	// Verify a signature on a string.
	Verify(*RequestVerify) (*ResponseVerify, error)
	// Export all outputs in hex format.
	ExportOutputs() (*ResponseExportOutputs, error)
	// Import outputs in hex format.
	ImportOutputs(*RequestImportOutputs) (*ResponseImportOutputs, error)
	// Export a signed set of key images.
	ExportKeyImages() (*ResponseExportKeyImages, error)
	// Import signed key images list and verify their spent status.
	ImportKeyImages(*RequestImportKeyImages) (*ResponseImportKeyImages, error)
	// Create a payment URI using the official URI spec.
	MakeURI(*RequestMakeURI) (*ResponseMakeURI, error)
	// Parse a payment URI to get payment information.
	ParseURI(*RequestParseURI) (*ResponseParseURI, error)
	// Retrieves entries from the address book.
	GetAddressBook(*RequestGetAddressBook) (*ResponseGetAddressBook, error)
	// Add an entry to the address book.
	AddAddressBook(*RequestAddAddressBook) (*ResponseAddAddressBook, error)
	// Delete an entry from the address book.
	DeleteAddressBook(*RequestDeleteAddressBook) error
	// Refresh a wallet after openning.
	Refresh(*RequestRefresh) (*ResponseRefresh, error)
	// Rescan the blockchain for spent outputs.
	RescanSpent() error
	// Start mining in the Monero daemon.
	StartMining(*RequestStartMining) error
	// Stop mining in the Monero daemon.
	StopMining() error
	// Get a list of available languages for your wallet's seed.
	GetLanguages() (*ResponseGetLanguages, error)
	// Create a new wallet. You need to have set the argument "–wallet-dir" when launching monero-wallet-rpc to make this work.
	CreateWallet(*RequestCreateWallet) error
	// Restores a wallet from a given wallet address, view key, and optional spend key.
	GenerateFromKeys(*RequestGenerateFromKeys) (*ResponseGenerateFromKeys, error)
	// Open a wallet. You need to have set the argument "–wallet-dir" when launching monero-wallet-rpc to make this work.
	OpenWallet(*RequestOpenWallet) error
	// Create and open a wallet on the RPC server from an existing mnemonic phrase and close the currently open wallet.
	RestoreDeterministicWallet(*RequestRestoreDeterministicWallet) (*ResponseRestoreDeterministicWallet, error)
	// Close the currently opened wallet, after trying to save it.
	CloseWallet() error
	// Change a wallet password.
	ChangeWalletPassword(*RequestChangeWalletPassword) error
	// Check if a wallet is a multisig one.
	IsMultisig() (*ResponseIsMultisig, error)
	// Prepare a wallet for multisig by generating a multisig string to share with peers.
	PrepareMultisig() (*ResponsePrepareMultisig, error)
	// Make a wallet multisig by importing peers multisig string.
	MakeMultisig(*RequestMakeMultisig) (*ResponseMakeMultisig, error)
	// Export multisig info for other participants.
	ExportMultisigInfo() (*ResponseExportMultisigInfo, error)
	// Import multisig info from other participants.
	ImportMultisigInfo(*RequestImportMultisigInfo) (*ResponseImportMultisigInfo, error)
	// Turn this wallet into a multisig wallet, extra step for N-1/N wallets.
	FinalizeMultisig(*RequestFinalizeMultisig) (*ResponseFinalizeMultisig, error)
	// Sign a transaction in multisig.
	SignMultisig(*RequestSignMultisig) (*ResponseSignMultisig, error)
	// Submit a signed multisig transaction.
	SubmitMultisig(*RequestSubmitMultisig) (*ResponseSubmitMultisig, error)
	// Get RPC version Major & Minor integer-format, where Major is the first 16 bits and Minor the last 16 bits.
	GetVersion() (*ResponseGetVersion, error)
}

// New returns a new monero-wallet-rpc client.
func New(cfg Config) Client {
	cl := &client{
		addr:    cfg.Address,
		headers: cfg.CustomHeaders,
	}
	if cfg.Transport == nil {
		cl.httpcl = http.DefaultClient
	} else {
		cl.httpcl = &http.Client{
			Transport: cfg.Transport,
		}
	}
	return cl
}

type client struct {
	httpcl  *http.Client
	addr    string
	headers map[string]string
}

// Helper function
func (c *client) do(method string, in, out interface{}) error {
	payload, err := json2.EncodeClientRequest(method, in)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.addr, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if c.headers != nil {
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := c.httpcl.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	// in theory this is only done to catch
	// any monero related errors if
	// we are not expecting any data back
	if out == nil {
		v := &json2.EmptyResponse{}
		return json2.DecodeClientResponse(resp.Body, v)
	}
	return json2.DecodeClientResponse(resp.Body, out)
}

// Methods
func (c *client) GetBalance(req *RequestGetBalance) (resp *ResponseGetBalance, err error) {
	err = c.do("get_balance", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetAddress(req *RequestGetAddress) (resp *ResponseGetAddress, err error) {
	err = c.do("get_address", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetAddressIndex(req *RequestGetAddressIndex) (resp *ResponseGetAddressIndex, err error) {
	err = c.do("get_address_index", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CreateAddress(req *RequestCreateAddress) (resp *ResponseCreateAddress, err error) {
	err = c.do("create_address", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) LabelAddress(req *RequestLabelAddress) (err error) {
	err = c.do("label_address", req, nil)
	if err != nil {
		return err
	}
	return
}

func (c *client) ValidateAddress(req *RequestValidateAddress) (resp *ResponseValidateAddress, err error) {
	err = c.do("validate_address", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetAccounts(req *RequestGetAccounts) (resp *ResponseGetAccounts, err error) {
	err = c.do("get_accounts", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CreateAccount(req *RequestCreateAccount) (resp *ResponseCreateAccount, err error) {
	err = c.do("create_account", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) LabelAccount(req *RequestLabelAccount) (err error) {
	err = c.do("label_account", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GetAccountTags() (resp *ResponseGetAccountTags, err error) {
	err = c.do("get_account_tags", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) TagAccounts(req *RequestTagAccounts) (err error) {
	err = c.do("tag_accounts", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) UntagAccounts(req *RequestUntagAccounts) (err error) {
	err = c.do("untag_accounts", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) SetAccountTagDescription(req *RequestSetAccountTagDescription) (err error) {
	err = c.do("set_account_tag_description", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GetHeight() (resp *ResponseGetHeight, err error) {
	err = c.do("get_height", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) Transfer(req *RequestTransfer) (resp *ResponseTransfer, err error) {
	err = c.do("transfer", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) TransferSplit(req *RequestTransferSplit) (resp *ResponseTransferSplit, err error) {
	err = c.do("transfer_split", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SignTransfer(req *RequestSignTransfer) (resp *ResponseSignTransfer, err error) {
	err = c.do("sign_transfer", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SubmitTransfer(req *RequestSubmitTransfer) (resp *ResponseSubmitTransfer, err error) {
	err = c.do("submit_transfer", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SweepDust(req *RequestSweepDust) (resp *ResponseSweepDust, err error) {
	err = c.do("sweep_dust", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SweepAll(req *RequestSweepAll) (resp *ResponseSweepAll, err error) {
	err = c.do("sweep_all", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SweepSingle(req *RequestSweepSingle) (resp *ResponseSweepSingle, err error) {
	err = c.do("sweep_single", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) RelayTx(req *RequestRelayTx) (resp *ResponseRelayTx, err error) {
	err = c.do("relay_tx", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) Store() (err error) {
	err = c.do("store", nil, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GetPayments(req *RequestGetPayments) (resp *ResponseGetPayments, err error) {
	err = c.do("get_payments", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetBulkPayments(req *RequestGetBulkPayments) (resp *ResponseGetBulkPayments, err error) {
	err = c.do("get_bulk_payments", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) IncomingTransfers(req *RequestIncomingTransfers) (resp *ResponseIncomingTransfers, err error) {
	err = c.do("incoming_transfers", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) QueryKey(req *RequestQueryKey) (resp *ResponseQueryKey, err error) {
	err = c.do("query_key", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) MakeIntegratedAddress(req *RequestMakeIntegratedAddress) (resp *ResponseMakeIntegratedAddress, err error) {
	err = c.do("make_integrated_address", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SplitIntegratedAddress(req *RequestSplitIntegratedAddress) (resp *ResponseSplitIntegratedAddress, err error) {
	err = c.do("split_integrated_address", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) StopWallet() (err error) {
	err = c.do("stop_wallet", nil, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) RescanBlockchain() (err error) {
	err = c.do("rescan_blockchain", nil, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) SetTxNotes(req *RequestSetTxNotes) (err error) {
	err = c.do("set_tx_notes", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GetTxNotes(req *RequestGetTxNotes) (resp *ResponseGetTxNotes, err error) {
	err = c.do("get_tx_notes", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SetAttribute(req *RequestSetAttribute) (err error) {
	err = c.do("set_attribute", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GetAttribute(req *RequestGetAttribute) (resp *ResponseGetAttribute, err error) {
	err = c.do("get_attribute", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetTxKey(req *RequestGetTxKey) (resp *ResponseGetTxKey, err error) {
	err = c.do("get_tx_key", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CheckTxKey(req *RequestCheckTxKey) (resp *ResponseCheckTxKey, err error) {
	err = c.do("check_tx_key", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetTxProof(req *RequestGetTxProof) (resp *ResponseGetTxProof, err error) {
	err = c.do("get_tx_proof", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CheckTxProof(req *RequestCheckTxProof) (resp *ResponseCheckTxProof, err error) {
	err = c.do("check_tx_proof", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetSpendProof(req *RequestGetSpendProof) (resp *ResponseGetSpendProof, err error) {
	err = c.do("get_spend_proof", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CheckSpendProof(req *RequestCheckSpendProof) (resp *ResponseCheckSpendProof, err error) {
	err = c.do("check_spend_proof", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetReserveProof(req *RequestGetReserveProof) (resp *ResponseGetReserveProof, err error) {
	err = c.do("get_reserve_proof", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CheckReserveProof(req *RequestCheckReserveProof) (resp *ResponseCheckReserveProof, err error) {
	err = c.do("check_reserve_proof", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetTransfers(req *RequestGetTransfers) (resp *ResponseGetTransfers, err error) {
	err = c.do("get_transfers", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetTransferByTxID(req *RequestGetTransferByTxID) (resp *ResponseGetTransferByTxID, err error) {
	err = c.do("get_transfer_by_txid", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) Sign(req *RequestSign) (resp *ResponseSign, err error) {
	err = c.do("sign", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) Verify(req *RequestVerify) (resp *ResponseVerify, err error) {
	err = c.do("verify", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ExportOutputs() (resp *ResponseExportOutputs, err error) {
	err = c.do("export_outputs", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ImportOutputs(req *RequestImportOutputs) (resp *ResponseImportOutputs, err error) {
	err = c.do("import_outputs", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ExportKeyImages() (resp *ResponseExportKeyImages, err error) {
	err = c.do("export_key_images", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ImportKeyImages(req *RequestImportKeyImages) (resp *ResponseImportKeyImages, err error) {
	err = c.do("import_key_images", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) MakeURI(req *RequestMakeURI) (resp *ResponseMakeURI, err error) {
	err = c.do("make_uri", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ParseURI(req *RequestParseURI) (resp *ResponseParseURI, err error) {
	err = c.do("parse_uri", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetAddressBook(req *RequestGetAddressBook) (resp *ResponseGetAddressBook, err error) {
	err = c.do("get_address_book", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) AddAddressBook(req *RequestAddAddressBook) (resp *ResponseAddAddressBook, err error) {
	err = c.do("add_address_book", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) DeleteAddressBook(req *RequestDeleteAddressBook) (err error) {
	err = c.do("delete_address_book", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) Refresh(req *RequestRefresh) (resp *ResponseRefresh, err error) {
	err = c.do("refresh", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) RescanSpent() (err error) {
	err = c.do("rescan_spent", nil, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) StartMining(req *RequestStartMining) (err error) {
	err = c.do("start_mining", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) StopMining() (err error) {
	err = c.do("stop_mining", nil, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GetLanguages() (resp *ResponseGetLanguages, err error) {
	err = c.do("get_languages", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) CreateWallet(req *RequestCreateWallet) (err error) {
	err = c.do("create_wallet", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) GenerateFromKeys(req *RequestGenerateFromKeys) (resp *ResponseGenerateFromKeys, err error) {
	err = c.do("generate_from_keys", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}
func (c *client) OpenWallet(req *RequestOpenWallet) (err error) {
	err = c.do("open_wallet", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) RestoreDeterministicWallet(req *RequestRestoreDeterministicWallet) (resp *ResponseRestoreDeterministicWallet, err error) {
	err = c.do("restore_deterministic_wallet", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}
func (c *client) CloseWallet() (err error) {
	err = c.do("close_wallet", nil, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) ChangeWalletPassword(req *RequestChangeWalletPassword) (err error) {
	err = c.do("change_wallet_password", &req, nil)
	if err != nil {
		return err
	}
	return
}
func (c *client) IsMultisig() (resp *ResponseIsMultisig, err error) {
	err = c.do("is_multisig", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) PrepareMultisig() (resp *ResponsePrepareMultisig, err error) {
	err = c.do("prepare_multisig", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) MakeMultisig(req *RequestMakeMultisig) (resp *ResponseMakeMultisig, err error) {
	err = c.do("make_multisig", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ExportMultisigInfo() (resp *ResponseExportMultisigInfo, err error) {
	err = c.do("export_multisig_info", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) ImportMultisigInfo(req *RequestImportMultisigInfo) (resp *ResponseImportMultisigInfo, err error) {
	err = c.do("import_multisig_info", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) FinalizeMultisig(req *RequestFinalizeMultisig) (resp *ResponseFinalizeMultisig, err error) {
	err = c.do("finalize_multisig", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SignMultisig(req *RequestSignMultisig) (resp *ResponseSignMultisig, err error) {
	err = c.do("sign_multisig", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) SubmitMultisig(req *RequestSubmitMultisig) (resp *ResponseSubmitMultisig, err error) {
	err = c.do("submit_multisig", &req, &resp)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) GetVersion() (resp *ResponseGetVersion, err error) {
	err = c.do("get_version", nil, &resp)
	if err != nil {
		return nil, err
	}
	return
}
