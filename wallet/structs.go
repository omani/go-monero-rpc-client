package wallet

// Helper structs
type Destination struct {
	// Amount to send to each destination, in atomic units.
	Amount uint64 `json:"amount"`
	// Destination public address.
	Address string `json:"address"`
}
type SignedKeyImage struct {
	KeyImage  string `json:"key_image"`
	Signature string `json:"signature"`
}

// *** RPC STRUCTS ***
// GetBalance()
type RequestGetBalance struct {
	// Return balance for this account.
	AccountIndex uint64 `json:"account_index"`
	// (Optional) Return balance detail for those subaddresses.
	AddressIndices []uint64 `json:"address_indices"`
}
type ResponseGetBalance struct {
	// The total balance of the current monero-wallet-rpc in session.
	Balance uint64 `json:"balance"`
	// Unlocked funds are those funds that are sufficiently deep enough in the Monero blockchain to be considered safe to spend.
	UnlockedBalance uint64 `json:"unlocked_balance"`
	// True if importing multisig data is needed for returning a correct balance.
	MultisigImportNeeded bool `json:"multisig_import_needed"`
	// Array of subaddress information. Balance information for each subaddress in an account:
	PerSubaddress []struct {
		// Index of the subaddress in the account.
		AddressIndex uint64 `json:"address_index"`
		// Address at this index. Base58 representation of the public keys.
		Address string `json:"address"`
		// Balance for the subaddress (locked or unlocked).
		Balance uint64 `json:"balance"`
		// Unlocked balance for the subaddress.
		UnlockedBalance uint64 `json:"unlocked_balance"`
		// Label for the subaddress.
		Label string `json:"label"`
		// Number of unspent outputs available for the subaddress.
		NumUnspentOutputs uint64 `json:"num_unspent_outputs"`
		// Blocks to unlock
		BlocksToUnlock int64 `json:"blocks_to_unlock"`
	} `json:"per_subaddress"`
}

// GetAddress()
type RequestGetAddress struct {
	// Return subaddresses for this account.
	AccountIndex uint64 `json:"account_index"`
	// (Optional) List of subaddresses to return from an account.
	AddressIndex []uint64 `json:"address_index"`
}
type ResponseGetAddress struct {
	// The 95-character hex address string of the monero-wallet-rpc in session.
	Address string `json:"address"`
	// Array of addresses informations
	Addresses []struct {
		// The 95-character hex (sub)address string.
		Address string `json:"address"`
		// Label of the (sub)address
		Label string `json:"label"`
		// Index of the subaddress
		AddressIndex uint64 `json:"address_index"`
		// States if the (sub)address has already received funds
		Used bool `json:"used"`
	} `json:"addresses"`
}

// GetAddressIndex()
type RequestGetAddressIndex struct {
	// (Sub)address to look for.
	Address string `json:"address"`
}
type ResponseGetAddressIndex struct {
	// Subaddress informations
	Index struct {
		// Account index.
		Major uint64 `json:"major"`
		// Address index.
		Minor uint64 `json:"minor"`
	} `json:"index"`
}

// CreateAddress()
type RequestCreateAddress struct {
	// Create a new address for this account.
	AccountIndex uint64 `json:"account_index"`
	// (Optional) Label for the new address.
	Label string `json:"label"`
}
type ResponseCreateAddress struct {
	// Newly created address. Base58 representation of the public keys.
	Address string `json:"address"`
	// Index of the new address under the input account.
	AddressIndex uint64 `json:"address_index"`
}

// LabelAddress()
type RequestLabelAddress struct {
	// Subaddress index; JSON Object containing the major & minor address
	Index struct {
		// Account index for the subaddress.
		Major uint64 `json:"major"`
		// Index of the subaddress in the account.
		Minor uint64 `json:"minor"`
	} `json:"index"`
	// Label for the address.
	Label string `json:"label"`
}

// ValidateAddress()
type RequestValidateAddress struct {
	Address        string `json:"address"`
	AnyNetType     bool   `json:"any_net_type,omitempty"`
	AllowOpenAlias bool   `json:"allow_openalias,omitempty"`
}

type ResponseValidateAddress struct {
	Valid            bool   `json:"valid"`
	Integrated       bool   `json:"integrated"`
	Subaddress       bool   `json:"subaddress"`
	NetType          string `json:"nettype"`
	OpenAliasAddress string `json:"openalias_address"`
}

// GetAccounts()
type RequestGetAccounts struct {
	// (Optional) Tag for filtering accounts.
	Tag string `json:"tag"`
}
type ResponseGetAccounts struct {
	// Array of subaddress account information:
	SubaddressAccounts []struct {
		// Index of the account.
		AccountIndex uint64 `json:"account_index"`
		// Balance of the account (locked or unlocked).
		Balance uint64 `json:"balance"`
		// Base64 representation of the first subaddress in the account.
		BaseAddress string `json:"base_address"`
		// (Optional) Label of the account.
		Label string `json:"label"`
		// (Optional) Tag for filtering accounts.
		Tag string `json:"tag"`
		// Unlocked balance for the account.
		UnlockedBalance uint64 `json:"unlocked_balance"`
	} `json:"subaddress_accounts"`
	// Total balance of the selected accounts (locked or unlocked).
	TotalBalance uint64 `json:"total_balance"`
	// Total unlocked balance of the selected accounts.
	TotalUnlockedBalance uint64 `json:"total_unlocked_balance"`
}

// CreateAccount()
type RequestCreateAccount struct {
	// (Optional) Label for the account.
	Label string `json:"label"`
}
type ResponseCreateAccount struct {
	// Index of the new account.
	AccountIndex uint64 `json:"account_index"`
	// Address for this account. Base58 representation of the public keys.
	Address string `json:"address"`
}

// LabelAccount()
type RequestLabelAccount struct {
	// Apply label to account at this index.
	AccountIndex uint64 `json:"account_index"`
	// Label for the account.
	Label string `json:"label"`
}

// GetAccountTags()
type ResponseGetAccountTags struct {
	// Array of account tag information:
	AccountTags []struct {
		// Filter tag.
		Tag string `json:"tag"`
		// Label for the tag.
		Label string `json:"label"`
		// List of tagged account indices.
		Accounts []uint64 `json:"accounts"`
	} `json:"account_tags"`
}

// TagAccounts()
type RequestTagAccounts struct {
	// Tag for the accounts.
	Tag string `json:"tag"`
	// Tag this list of accounts.
	Accounts []uint64 `json:"accounts"`
}

// UntagAccounts()
type RequestUntagAccounts struct {
	// Remove tag from this list of accounts.
	Accounts []uint64 `json:"accounts"`
}

// SetAccountTagDescription()
type RequestSetAccountTagDescription struct {
	// Set a description for this tag.
	Tag string `json:"tag"`
	// Description for the tag.
	Description string `json:"description"`
}

// GetHeight()
type ResponseGetHeight struct {
	// The current monero-wallet-rpc's blockchain height. If the wallet has been offline for a long time, it may need to catch up with the daemon.
	Height uint64 `json:"height"`
}

// Transfer()
type RequestTransfer struct {
	// Array of destinations to receive XMR:
	Destinations []*Destination `json:"destinations"`
	// (Optional) Transfer from this account index. (Defaults to 0)
	AccountIndex uint64 `json:"account_index"`
	// (Optional) Transfer from this set of subaddresses. (Defaults to empty - all indices)
	SubaddrIndices []uint64 `json:"subaddr_indices"`
	// Set a priority for the transaction. Accepted Values are: 0-3 for: default, unimportant, normal, elevated, priority.
	Priority Priority `json:"priority"`
	// Number of outputs from the blockchain to mix with (0 means no mixing).
	Mixing uint64 `json:"mixin"`
	// (Optional) Number of outputs to mix in the transaction (this output + N decoys from the blockchain).
	RingSize uint64 `json:"ring_size,omitempty"`
	// Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint64 `json:"unlock_time"`
	// (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymentID string `json:"payment_id"`
	// (Optional) Return the transaction key after sending.
	GetTxKey bool `json:"get_tx_key"`
	// (Optional) If true, the newly created transaction will not be relayed to the monero network. (Defaults to false)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	// (Optional) Return the transaction as hex string after sending (Defaults to false)
	GetTxHex bool `json:"get_tx_hex,omitempty"`
	// (Optional) Return the metadata needed to relay the transaction. (Defaults to false)
	GetTxMetadata bool `json:"get_tx_metadata,omitempty"`
}
type ResponseTransfer struct {
	// Amount transferred for the transaction.
	Amount uint64 `json:"amount"`
	// Integer value of the fee charged for the txn.
	Fee uint64 `json:"fee"`
	// MultiTxSet multisig_txset - Set of multisig transactions in the process of being signed (empty for non-multisig).

	// Raw transaction represented as hex string, if get_tx_hex is true.
	TxBlob string `json:"tx_blob"`
	// String for the publically searchable transaction hash.
	TxHash string `json:"tx_hash"`
	// String for the transaction key if get_tx_key is true, otherwise, blank string.
	TxKey      string `json:"tx_key"`
	TxMetadata string `json:"tx_metadata"` // TxMetadata tx_metadata - Set of transaction metadata needed to relay this transfer later, if get_tx_metadata is true.

	// String. Set of unsigned tx for cold-signing purposes.
	UnsignedTxSet string `json:"unsigned_txset"`
}

// TransferSplit()
type RequestTransferSplit struct {
	// Array of destinations to receive XMR:
	Destinations []*Destination `json:"destinations"`
	// (Optional) Transfer from this account index. (Defaults to 0)
	AccountIndex uint64 `json:"account_index"`
	// (Optional) Transfer from this set of subaddresses. (Defaults to empty - all indices)
	SubaddrIndices []uint64 `json:"subaddr_indices"`
	// Number of outputs from the blockchain to mix with (0 means no mixing).
	Mixin uint64 `json:"mixin"`
	// (Optional) Sets ringsize to n (mixin + 1).
	RingSize uint64 `json:"ring_size,omitempty"`
	// Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint64 `json:"unlock_time"`
	// (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymendID string `json:"payment_id"`
	// (Optional) Return the transaction keys after sending.
	GetxKeys bool `json:"get_tx_keys"`
	// Set a priority for the transactions. Accepted Values are: 0-3 for: default, unimportant, normal, elevated, priority.
	Priority Priority `json:"priority"`
	// (Optional) If true, the newly created transaction will not be relayed to the monero network. (Defaults to false)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	// (Optional) Return the transactions as hex string after sending
	GetTxHex bool `json:"get_tx_hex,omitempty"`
	// True to use the new transaction construction algorithm, defaults to false.
	NewAlgorithm bool `json:"new_algorithm"`
	// (Optional) Return list of transaction metadata needed to relay the transfer later.
	GetTxMetadata bool `json:"get_tx_metadata,omitempty"`
}
type ResponseTransferSplit struct {
	// The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// The transaction keys for every transaction.
	TxKeyList []string `json:"tx_key_list"`
	// The amount transferred for every transaction.
	AmountList []uint64 `json:"amount_list"`
	// The amount of fees paid for every transaction.
	FeeList []uint64 `json:"fee_list"`
	// The tx as hex string for every transaction.
	TxBlobList []string `json:"tx_blob_list"`
	// List of transaction metadata needed to relay the transactions later.
	TxMetadataList []string `json:"tx_metadata_list"`
	// The set of signing keys used in a multisig transaction (empty for non-multisig).
	MultisigTxSet string `json:"multisig_txset"`
	// Set of unsigned tx for cold-signing purposes.
	UnsignedTxSet string `json:"unsigned_txset"`
}

// SignTransfer()
type RequestSignTransfer struct {
	// Set of unsigned tx returned by "transfer" or "transfer_split" methods.
	UnsighnedxSet string `json:"unsigned_txset"`
	// (Optional) If true, return the raw transaction data. (Defaults to false)
	ExportRaw bool `json:"export_raw,omitempty"`
}
type ResponseSignTransfer struct {
	// Set of signed tx to be used for submitting transfer.
	SignedTxSet string `json:"signed_txset"`
	// The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// The tx raw data of every transaction.
	TxRawList []string `json:"tx_raw_list"`
}

// SubmitTransfer()
type RequestSubmitTransfer struct {
	// Set of signed tx returned by "sign_transfer"
	TxDataHex string `json:"tx_data_hex"`
}
type ResponseSubmitTransfer struct {
	// The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
}

// SweepDust()
type RequestSweepDust struct {
	// (Optional) Return the transaction keys after sending.
	GetTxKeys bool `json:"get_tx_keys"`
	// (Optional) If true, the newly created transaction will not be relayed to the monero network. (Defaults to false)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	// (Optional) Return the transactions as hex string after sending. (Defaults to false)
	GetTxHex bool `json:"get_tx_hex,omitempty"`
	// (Optional) Return list of transaction metadata needed to relay the transfer later. (Defaults to false)
	GetTxMetadata bool `json:"get_tx_metadata,omitempty"`
}
type ResponseSweepDust struct {
	// The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// The transaction keys for every transaction.
	TxKeyList []string `json:"tx_key_list"`
	//  The amount transferred for every transaction.
	AmountList []uint64 `json:"amount_list"`
	//  The amount of fees paid for every transaction.
	FeeList []uint64 `json:"fee_list"`
	// The tx as hex string for every transaction.
	TxBlobList []string `json:"tx_blob_list"`
	// List of transaction metadata needed to relay the transactions later.
	TxMetadataList []string `json:"tx_metadata_list"`
	// The set of signing keys used in a multisig transaction (empty for non-multisig).
	MultisigTxSet string `json:"multisig_txset"`
	// Set of unsigned tx for cold-signing purposes.
	UnsignedTxSet string `json:"unsigned_txset"`
}

// SweepAll()
type RequestSweepAll struct {
	//  Destination public address.
	Address string `json:"address"`
	//  Sweep transactions from this account.
	AccountIndex uint64 `json:"account_index"`
	//  (Optional) Sweep from this set of subaddresses in the account.
	SubaddrIndices []uint64 `json:"subaddr_indices"`
	//  (Optional) Sweep from all subaddresses in the account (default: false).
	SubaddrIndicesAll bool `json:"subaddr_indices_all"`
	//  (Optional) Priority for sending the sweep transfer, partially determines fee.
	Priority Priority `json:"priority"`
	//  Number of outputs from the blockchain to mix with (0 means no mixing).
	Mixin uint64 `json:"mixin"`
	//  (Optional) Sets ringsize to n (mixin + 1).
	RingSize uint64 `json:"ring_size,omitempty"`
	//  Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint64 `json:"unlock_time"`
	//  (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymentID string `json:"payment_id"`
	//  (Optional) Return the transaction keys after sending.
	GetTxKeys bool `json:"get_tx_keys"`
	//  (Optional) Include outputs below this amount.
	BelowAmount uint64 `json:"below_amount"`
	//  (Optional) If true, do not relay this sweep transfer. (Defaults to false)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	//  (Optional) return the transactions as hex encoded string. (Defaults to false)
	GetTxHex bool `json:"get_tx_hex,omitempty"`
	//  (Optional) return the transaction metadata as a string. (Defaults to false)
	GetTxMetadata bool `json:"get_tx_metadata,omitempty"`
}
type ResponseSweepAll struct {
	// The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// The transaction keys for every transaction.
	TxKeyList []string `json:"tx_key_list"`
	// The amount transferred for every transaction.
	AmountList []uint64 `json:"amount_list"`
	// The amount of fees paid for every transaction.
	FeeList []uint64 `json:"fee_list"`
	// The tx as hex string for every transaction.
	TxBlobList []string `json:"tx_blob_list"`
	// List of transaction metadata needed to relay the transactions later.
	TxMetadataList []string `json:"tx_metadata_list"`
	// Set of signing keys used in a multisig transaction (empty for non-multisig).
	MultisigTxSet string `json:"multisig_txset"`
	// Set of unsigned tx for cold-signing purposes.
	UnsignedTxSet string `json:"unsigned_txset"`
}

// SweepSingle()
type RequestSweepSingle struct {
	// Destination public address.
	Address string `json:"address"`
	// Sweep transactions from this account.
	AccountIndex uint64 `json:"account_index"`
	// (Optional) Sweep from this set of subaddresses in the account.
	SubaddrIndices []uint64 `json:"subaddr_indices"`
	// (Optional) Priority for sending the sweep transfer, partially determines fee.
	Priority Priority `json:"priority"`
	// Number of outputs from the blockchain to mix with (0 means no mixing).
	Mixin uint64 `json:"mixin"`
	// (Optional) Sets ringsize to n (mixin + 1).
	RingSize uint64 `json:"ring_size,omitempty"`
	// Number of blocks before the monero can be spent (0 to not add a lock).
	UnlockTime uint64 `json:"unlock_time"`
	// (Optional) Random 32-byte/64-character hex string to identify a transaction.
	PaymentID string `json:"payment_id"`
	// (Optional) Return the transaction keys after sending.
	GetxKeys bool `json:"get_tx_keys"`
	// Key image of specific output to sweep.
	KeyImage string `json:"key_image"`
	// (Optional) Include outputs below this amount.
	BelowAmount uint64 `json:"below_amount"`
	// (Optional) If true, do not relay this sweep transfer. (Defaults to false)
	DoNotRelay bool `json:"do_not_relay,omitempty"`
	// (Optional) return the transactions as hex encoded string. (Defaults to false)
	GetTxHex bool `json:"get_tx_hex,omitempty"`
	// (Optional) return the transaction metadata as a string. (Defaults to false)
	GetTxMetadata bool `json:"get_tx_metadata,omitempty"`
}
type ResponseSweepSingle struct {
	// The tx hashes of every transaction.
	TxHashList []string `json:"tx_hash_list"`
	// The transaction keys for every transaction.
	TxKeyList []string `json:"tx_key_list"`
	// The amount transferred for every transaction.
	AmountList []uint64 `json:"amount_list"`
	// The amount of fees paid for every transaction.
	FreeList []uint64 `json:"fee_list"`
	// The tx as hex string for every transaction.
	TxBlobList []string `json:"tx_blob_list"`
	// List of transaction metadata needed to relay the transactions later.
	TxMetadataList []string `json:"tx_metadata_list"`
	// The set of signing keys used in a multisig transaction (empty for non-multisig).
	MultisigTxSet string `json:"multisig_txset"`
	// Set of unsigned tx for cold-signing purposes.
	UnsignedTxSet string `json:"unsigned_txset"`
}

// RelayTx()
type RequestRelayTx struct {
	// Transaction metadata returned from a transfer method with get_tx_metadata set to true.
	Hex string `json:"hex"`
}
type ResponseRelayTx struct {
	// String for the publically searchable transaction hash.
	TxHash string `json:"tx_hash"`
}

// GetPayments()
type RequestGetPayments struct {
	// Payment ID used to find the payments (16 characters hex).
	PaymentID string `json:"payment_id"`
}
type ResponseGetPayments struct {
	// list of payments
	Payments []struct {
		// Payment ID matching the input parameter.
		PaymentID string `json:"payment_id"`
		// Transaction hash used as the transaction ID.
		TxHash string `json:"tx_hash"`
		// Amount for this payment.
		Amount uint64 `json:"amount"`
		// Height of the block that first confirmed this payment.
		BlockHeight uint64 `json:"block_height"`
		// Time (in block height) until this payment is safe to spend.
		UnlockTime uint64 `json:"unlock_time"`
		// Subaddress index:
		SubaddrIndex struct {
			// Account index for the subaddress.
			Major uint64 `json:"major"`
			// Index of the subaddress in the account.
			Minor uint64 `json:"minor"`
		} `json:"subaddr_index"`
		// Address receiving the payment; Base58 representation of the public keys.
		Address string `json:"address"`
	} `json:"payments"`
}

// GetBulkPayments()
type RequestGetBulkPayments struct {
	// Payment IDs used to find the payments (16 characters hex).
	PaymentIDs []string `json:"payment_ids"`
	// The block height at which to start looking for payments.
	MinBlockHeight uint64 `json:"min_block_height"`
}
type ResponseGetBulkPayments struct {
	// List of payments
	Payments []struct {
		// Payment ID matching one of the input IDs.
		PaymentID string `json:"payment_id"`
		// Transaction hash used as the transaction ID.
		TxHash string `json:"tx_hash"`
		// Amount for this payment.
		Amount uint64 `json:"amount"`
		// Height of the block that first confirmed this payment.
		BlockHeight uint64 `json:"block_height"`
		// Time (in block height) until this payment is safe to spend.
		UnlockTime uint64 `json:"unlock_time"`
		// Subaddress index:
		SubaddrIndex struct {
			// Account index for the subaddress.
			Major uint64 `json:"major"`
			// Index of the subaddress in the account.
			Minor uint64 `json:"minor"`
		} `json:"subaddr_index"`
		// Address receiving the payment; Base58 representation of the public keys.
		Address string `json:"address"`
	} `json:"payments"`
}

// IncomingTransfers()
type RequestIncomingTransfers struct {
	// "all": all the transfers, "available": only transfers which are not yet spent, OR "unavailable": only transfers which are already spent.
	TransferType string `json:"transfer_type"`
	// (Optional) Return transfers for this account. (defaults to 0)
	AccountIndex uint64 `json:"account_index"`
	// (Optional) Return transfers sent to these subaddresses.
	SubaddrIndices []uint64 `json:"subaddr_indices"`
	// (Optional) Enable verbose output, return key image if true.
	Verbose bool `json:"verbose"`
}
type ResponseIncomingTransfers struct {
	// list of transfers:
	Transfers struct {
		// Amount of this transfer.
		Amount uint64 `json:"amount"`
		// Mostly internal use, can be ignored by most users.
		GlobalIndex uint64 `json:"global_index"`
		// Key image for the incoming transfer's unspent output (empty unless verbose is true).
		KeyImage string `json:"key_image"`
		// Indicates if this transfer has been spent.
		Spent bool `json:"spent"`
		// Subaddress index for incoming transfer.
		SubaddrIndex uint64 `json:"subaddr_index"`
		// Several incoming transfers may share the same hash if they were in the same transaction.
		TxHash string `json:"tx_hash"`
		// Size of transaction in bytes.
		TxSize uint64 `json:"tx_size"`
	} `json:"transfers"`
}

// QueryKey()
type RequestQueryKey struct {
	// Which key to retrieve: "mnemonic" - the mnemonic seed (older wallets do not have one) OR "view_key" - the view key
	KeyType string `json:"key_type"`
}
type ResponseQueryKey struct {
	// The view key will be hex encoded, while the mnemonic will be a string of words.
	Key string `json:"key"`
}

// MakeIntegratedAddress()
type RequestMakeIntegratedAddress struct {
	// (Optional, defaults to primary address) Destination public address.
	StandardAddress string `json:"standard_address"`
	// (Optional, defaults to a random ID) 16 characters hex encoded.
	PaymentID string `json:"payment_id"`
}
type ResponseMakeIntegratedAddress struct {
	// The newly created integrated address
	IntegratedAddress string `json:"integrated_address"`
	// Hex encoded payment id
	PaymentID string `json:"payment_id"`
}

// SplitIntegratedAddress()
type RequestSplitIntegratedAddress struct {
	// Integrated address
	IntegratedAddress string `json:"integrated_address"`
}
type ResponseSplitIntegratedAddress struct {
	// States if the address is a subaddress
	IsSubaddress bool `json:"is_subaddress"`
	// Hex encoded payment id
	PaymentID string `json:"payment_id"`
	// Address of integrated address
	StandardAddress string `json:"standard_address"`
}

// SetTxNotes()
type RequestSetTxNotes struct {
	// Transaction ids
	TxIDs []string `json:"txids"`
	// Notes for the transactions
	Notes []string `json:"notes"`
}

// GetTxNotes()
type RequestGetTxNotes struct {
	// Transaction ids
	TxIDs []string `json:"txids"`
}
type ResponseGetTxNotes struct {
	// Notes for the transactions
	Notes []string `json:"notes"`
}

// SetAttribute()
type RequestSetAttribute struct {
	// Attribute name
	Key string `json:"key"`
	// Attribute value
	Value string `json:"value"`
}

// GetAttribute()
type RequestGetAttribute struct {
	// Attribute name
	Key string `json:"key"`
}
type ResponseGetAttribute struct {
	// Attribute value
	Value string `json:"value"`
}

// GetTxKey()
type RequestGetTxKey struct {
	// Transaction id.
	TxID string `json:"txid"`
}
type ResponseGetTxKey struct {
	// Transaction secret key.
	TxKey string `json:"tx_key"`
}

// CheckTxKey()
type RequestCheckTxKey struct {
	// Transaction id.
	TxID string `json:"txid"`
	// Transaction secret key.
	TxKey string `json:"tx_key"`
	// Destination public address of the transaction.
	Address string `json:"address"`
}
type ResponseCheckTxKey struct {
	// Number of block mined after the one with the transaction.
	Confirmations uint64 `json:"confirmations"`
	// States if the transaction is still in pool or has been added to a block.
	InPool bool `json:"in_pool"`
	// Amount of the transaction.
	Received uint64 `json:"received"`
}

// GetTxProof()
type RequestGetTxProof struct {
	// Transaction id.
	TxID string `json:"txid"`
	// Destination public address of the transaction.
	Address string `json:"address"`
	// (Optional) add a message to the signature to further authenticate the prooving process.
	Message string `json:"message"`
}
type ResponseGetTxProof struct {
	// Transaction signature.
	Signature string `json:"signature"`
}

// CheckTxProof()
type RequestCheckTxProof struct {
	// Transaction id.
	TxID string `json:"txid"`
	// Destination public address of the transaction.
	Address string `json:"address"`
	// (Optional) Should be the same message used in get_tx_proof.
	Message string `json:"message"`
	// Transaction signature to confirm.
	Signature string `json:"signature"`
}
type ResponseCheckTxProof struct {
	// Number of block mined after the one with the transaction.
	Confirmations uint64 `json:"confirmations"`
	// States if the inputs proves the transaction.
	Good bool `json:"good"`
	// States if the transaction is still in pool or has been added to a block.
	InPool bool `json:"in_pool"`
	// Amount of the transaction.
	Received uint64 `json:"received"`
}

// GetSpendProof()
type RequestGetSpendProof struct {
	// Transaction id.
	TxID string `json:"txid"`
	// (Optional) add a message to the signature to further authenticate the prooving process.
	Message string `json:"message"`
}
type ResponseGetSpendProof struct {
	// Spend signature.
	Signature string `json:"signature"`
}

// CheckSpendProof()
type RequestCheckSpendProof struct {
	// Transaction id.
	TxID string `json:"txid"`
	// (Optional) Should be the same message used in get_spend_proof.
	Message string `json:"message"`
	// Spend signature to confirm.
	Signature string `json:"signature"`
}
type ResponseCheckSpendProof struct {
	// States if the inputs proves the spend.
	Good bool `json:"good"`
}

// GetReserveProof()
type RequestGetReserveProof struct {
	// Proves all wallet balance to be disposable.
	All bool `json:"all"`
	// Specify the account from witch to prove reserve. (ignored if all is set to true)
	AccountIndex uint64 `json:"account_index"`
	// Amount (in atomic units) to prove the account has for reserve. (ignored if all is set to true)
	Amount uint64 `json:"amount"`
	// (Optional) add a message to the signature to further authenticate the prooving process.
	Message string `json:"message"`
}
type ResponseGetReserveProof struct {
	// Reserve signature.
	Signature string `json:"signature"`
}

// CheckReserveProof()
type RequestCheckReserveProof struct {
	// Public address of the wallet.
	Address string `json:"address"`
	// (Optional) Should be the same message used in get_reserve_proof.
	Message string `json:"message"`
	// Reserve signature to confirm.
	Signature string `json:"signature"`
}
type ResponseCheckReserveProof struct {
	// States if the inputs proves the reserve.
	Good bool `json:"good"`
}

// GetTransfers()
type RequestGetTransfers struct {
	// (Optional) Include incoming transfers.
	In bool `json:"in"`
	// (Optional) Include outgoing transfers.
	Out bool `json:"out"`
	// (Optional) Include pending transfers.
	Pending bool `json:"pending"`
	// (Optional) Include failed transfers.
	Failed bool `json:"failed"`
	// (Optional) Include transfers from the daemon's transaction pool.
	Pool bool `json:"pool"`
	// (Optional) Filter transfers by block height.
	FilterByHeight bool `json:"filter_by_height"`
	// (Optional) Minimum block height to scan for transfers, if filtering by height is enabled.
	MinHeight uint64 `json:"min_height"`
	// (Opional) Maximum block height to scan for transfers, if filtering by height is enabled (defaults to max block height).
	MaxHeight uint64 `json:"max_height,omitempty"`
	// (Optional) Index of the account to query for transfers. (defaults to 0)
	AccountIndex uint64 `json:"account_index"`
	// (Optional) List of subaddress indices to query for transfers. (Defaults to empty - all indices)
	SubaddrIndices []uint64 `json:"subaddr_indices"`
}
type Transfer struct {
	// Public address of the transfer.
	Address string `json:"address"`
	// Amount transferred.
	Amount uint64 `json:"amount"`
	// Number of block mined since the block containing this transaction (or block height at which the transaction should be added to a block if not yet confirmed).
	Confirmations uint64 `json:"confirmations"`
	// JSON objects containing transfer destinations:
	Destinations []*Destination `json:"destinations"`
	// True if the key image(s) for the transfer have been seen before.
	DoubleSpendSeen bool `json:"double_spend_seen"`
	// Transaction fee for this transfer.
	Fee uint64 `json:"fee"`
	// Height of the first block that confirmed this transfer (0 if not mined yet).
	Height uint64 `json:"height"`
	// Note about this transfer.
	Note string `json:"note"`
	// Payment ID for this transfer.
	PaymentID string `json:"payment_id"`
	// JSON object containing the major & minor subaddress index:
	SubaddrIndex struct {
		// Account index for the subaddress.
		Major uint64 `json:"major"`
		// Index of the subaddress under the account.
		Minor uint64 `json:"minor"`
	} `json:"subaddr_index"`
	// Estimation of the confirmations needed for the transaction to be included in a block.
	SuggestedConfirmationsThreshold uint64 `json:"suggested_confirmations_threshold"`
	// POSIX timestamp for when this transfer was first confirmed in a block (or timestamp submission if not mined yet).
	Timestamp uint64 `json:"timestamp"`
	// Transaction ID for this transfer.
	TxID string `json:"txid"`
	// Transfer type: "in/out/pending/failed/pool"
	Type string `json:"type"`
	// Number of blocks until transfer is safely spendable.
	UnlockTime uint64 `json:"unlock_time"`
}
type ResponseGetTransfers struct {
	// Array of transfers:
	In      []*Transfer `json:"in"`
	Out     []*Transfer `json:"out"`
	Pending []*Transfer `json:"pending"`
	Failed  []*Transfer `json:"failed"`
	Pool    []*Transfer `json:"pool"`
}

// GetTransferByTxID()
type RequestGetTransferByTxID struct {
	// Transaction ID used to find the transfer.
	TxID string `json:"txid"`
	// (Optional) Index of the account to query for the transfer.
	AccountIndex uint64 `json:"account_index,omitempty"`
}
type ResponseGetTransferByTxID struct {
	// JSON object containing payment information:
	Transfer Transfer `json:"transfer"`
}

// Sign()
type RequestSign struct {
	// Anything you need to sign.
	Data string `json:"data"`
}
type ResponseSign struct {
	// Signature generated against the "data" and the account public address.
	Signature string `json:"signature"`
}

// Verify()
type RequestVerify struct {
	// What should have been signed.
	Data string `json:"data"`
	// Public address of the wallet used to sign the data.
	Address string `json:"address"`
	// Signature generated by sign method.
	Signature string `json:"signature"`
}
type ResponseVerify struct {
	// True if signature is valid.
	Good bool `json:"good"`
}

// ExportOutputs()
type ResponseExportOutputs struct {
	// Wallet outputs in hex format.
	OutputsDataHex string `json:"outputs_data_hex"`
}

// ImportOutputs()
type RequestImportOutputs struct {
	// Wallet outputs in hex format.
	OutputsDataHex string `json:"outputs_data_hex"`
}
type ResponseImportOutputs struct {
	// Number of outputs imported.
	NumImported uint64 `json:"num_imported"`
}

// ExportKeyImages()
type ResponseExportKeyImages struct {
	// Array of signed key images:
	SignedKeyImages []struct {
		KeyImage  string `json:"key_image"`
		Signature string `json:"signature"`
	} `json:"signed_key_images"`
}

// ImportKeyImages()
type RequestImportKeyImages struct {
	// Array of signed key images:
	SignedKeyImages []*SignedKeyImage `json:"signed_key_images"`
}
type ResponseImportKeyImages struct {
	Height uint64 `json:"height"`
	// Amount (in atomic units) spent from those key images.
	Spent uint64 `json:"spent"`
	// Amount (in atomic units) still available from those key images.
	Unspent uint64 `json:"unspent"`
}

// MakeURI()
type RequestMakeURI struct {
	// Wallet address
	Address string `json:"address"`
	// (Optional) the integer amount to receive, in atomic units
	Amount uint64 `json:"amount"`
	// (Optional) 16 or 64 character hexadecimal payment id
	PaymentID string `json:"payment_id"`
	// (Optional) name of the payment recipient
	RecipientName string `json:"recipient_name"`
	// (Optional) Description of the reason for the tx
	TxDescription string `json:"tx_description"`
}
type ResponseMakeURI struct {
	// This contains all the payment input information as a properly formatted payment URI
	URI string `json:"uri"`
}

// ParseURI()
type RequestParseURI struct {
	// This contains all the payment input information as a properly formatted payment URI
	URI string `json:"uri"`
}
type ResponseParseURI struct {
	// JSON object containing payment information:
	URI struct {
		// Wallet address
		Address string `json:"address"`
		// Integer amount to receive, in atomic units (0 if not provided)
		Amount uint64 `json:"amount"`
		// 16 or 64 character hexadecimal payment id (empty if not provided)
		PaymentID string `json:"payment_id"`
		// Name of the payment recipient (empty if not provided)
		RecipientName string `json:"recipient_name"`
		// Description of the reason for the tx (empty if not provided)
		TxDescription string `json:"tx_description"`
	} `json:"uri"`
}

// GetAddressBook()
type RequestGetAddressBook struct {
	// Indices of the requested address book entries
	Entries []uint64 `json:"entries"`
}
type ResponseGetAddressBook struct {
	// Array of entries:
	Entries []struct {
		// Public address of the entry
		Address string `json:"address"`
		// Description of this address entry
		Description string `json:"description"`
		Index       uint64 `json:"index"`
		PaymentID   string `json:"payment_id"`
	} `json:"entries"`
}

// AddAddressBook()
type RequestAddAddressBook struct {
	Address string `json:"address"`
	// (Optional) string, defaults to "0000000000000000000000000000000000000000000000000000000000000000";
	PaymentID string `json:"payment_id"`
	// (Optional) string, defaults to "";
	Description string `json:"description"`
}
type ResponseAddAddressBook struct {
	// The index of the address book entry.
	Index uint64 `json:"index"`
}

// DeleteAddressBook()
type RequestDeleteAddressBook struct {
	// The index of the address book entry.
	Index uint64 `json:"index"`
}

// Refresh()
type RequestRefresh struct {
	// (Optional) The block height from which to start refreshing.
	StartHeight uint64 `json:"start_height,omitempty"`
}
type ResponseRefresh struct {
	// Number of new blocks scanned.
	BlocksFetched uint64 `json:"blocks_fetched"`
	// States if transactions to the wallet have been found in the blocks.
	ReceivedMoney bool `json:"received_money"`
}

// StartMining()
type RequestStartMining struct {
	// Number of threads created for mining.
	ThreadsCount uint64 `json:"threads_count"`
	// Allow to start the miner in smart mining mode.
	DoBackgroundMining bool `json:"do_background_mining"`
	// Ignore battery status (for smart mining only)
	IgnoreBattery bool `json:"ignore_battery"`
}

// GetLanguages()
type ResponseGetLanguages struct {
	// List of available languages
	Languages []string `json:"languages"`
}

// CreateWallet()
type RequestCreateWallet struct {
	// Wallet file name.
	Filename string `json:"filename"`
	// (Optional) password to protect the wallet.
	Password string `json:"password"`
	// Language for your wallets' seed.
	Language string `json:"language"`
}

// GenerateFromKeys()
type RequestGenerateFromKeys struct {
	// (Optional) The block height to restore the wallet from. (Defaults to 0)
	RestoreHeight int64 `json:"restore_height"`
	// The wallet's file name on the RPC server.
	Filename string `json:"filename"`
	// The wallet's primary address.
	Address string `json:"address"`
	// (Optional - omit to create a view-only wallet) The wallet's private spend key.
	SpendKey string `json:"spendkey"`
	// The wallet's private view key.
	ViewKey string `json:"viewkey"`
	// The wallet's password.
	Password string `json:"password"`
	// (Optional) If true, save the current wallet before generating the new wallet. (Defaults to true)
	AutoSaveCurrent bool `json:"autosave_current"`
	// (Optional) Language for your wallets' seed. (Defaults is "English")
	Language bool `json:"language"`
}

// GenerateFromKeys()
type ResponseGenerateFromKeys struct {
	// The wallet's address.
	Address string `json:"address"`
	// Verification message indicating that the wallet was generated successfully and whether or not it is a view-only wallet.
	Info string `json:"info"`
}

// OpenWallet()
type RequestOpenWallet struct {
	// Wallet name stored in –wallet-dir.
	Filename string `json:"filename"`
	// (Optional) only needed if the wallet has a password defined.
	Password string `json:"password"`
}

// ChangeWalletPassword()
type RequestChangeWalletPassword struct {
	// (Optional) Current wallet password, if defined.
	OldPassword string `json:"old_password"`
	// (Optional) New wallet password, if not blank.
	NewPassword string `json:"new_password"`
}

// IsMultisig()
type ResponseIsMultisig struct {
	// States if the wallet is multisig
	Multisig bool `json:"multisig"`
	Ready    bool `json:"ready"`
	// Amount of signature needed to sign a transfer.
	Threshold uint64 `json:"threshold"`
	// Total amount of signature in the multisig wallet.
	Total uint64 `json:"total"`
}

// PrepareMultisig()
type ResponsePrepareMultisig struct {
	// Multisig string to share with peers to create the multisig wallet.
	MultisigInfo string `json:"multisig_info"`
}

// MakeMultisig()
type RequestMakeMultisig struct {
	// List of multisig string from peers.
	MultisigInfo []string `json:"multisig_info"`
	// Amount of signatures needed to sign a transfer. Must be less or equal than the amount of signature in multisig_info.
	Threshold uint64 `json:"threshold"`
	// Wallet password
	Password string `json:"password"`
}
type ResponseMakeMultisig struct {
	// Multisig wallet address.
	Address string `json:"address"`
	// Multisig string to share with peers to create the multisig wallet (extra step for N-1/N wallets).
	MultisigInfo string `json:"multisig_info"`
}

// ExportMultisigInfo()
type ResponseExportMultisigInfo struct {
	// Multisig info in hex format for other participants.
	Info string `json:"info"`
}

// ImportMultisigInfo()
type RequestImportMultisigInfo struct {
	// List of multisig info in hex format from other participants.
	Info []string `json:"info"`
}
type ResponseImportMultisigInfo struct {
	// Number of outputs signed with those multisig info.
	NOutputs uint64 `json:"n_outputs"`
}

// FinalizeMultisig()
type RequestFinalizeMultisig struct {
	// List of multisig string from peers.
	MultisigInfo []string `json:"multisig_info"`
	// Wallet password
	Password string `json:"password"`
}
type ResponseFinalizeMultisig struct {
	// Multisig wallet address.
	Address string `json:"address"`
}

// SignMultisig()
type RequestSignMultisig struct {
	// Multisig transaction in hex format, as returned by transfer under multisig_txset.
	TxDataHex string `json:"tx_data_hex"`
}
type ResponseSignMultisig struct {
	// Multisig transaction in hex format.
	TxDataHex string `json:"tx_data_hex"`
	// List of transaction Hash.
	TxHashList []string `json:"tx_hash_list"`
}

// SubmitMultisig()
type RequestSubmitMultisig struct {
	// Multisig transaction in hex format, as returned by sign_multisig under tx_data_hex.
	TxDataHex string `json:"tx_data_hex"`
}
type ResponseSubmitMultisig struct {
	// List of transaction Hash.
	TxHashList []string `json:"tx_hash_list"`
}

// GetVersion()
type ResponseGetVersion struct {
	// RPC version, formatted with Major * 2^16 + Minor (Major encoded over the first 16 bits, and Minor over the last 16 bits).
	Version uint64 `json:"version"`
}
