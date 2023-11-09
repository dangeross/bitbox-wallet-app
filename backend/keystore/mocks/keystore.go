// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/keystore"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/signing"
	"github.com/ethereum/go-ethereum/core/types"
	"sync"
)

// Ensure, that KeystoreMock does implement keystore.Keystore.
// If this is not the case, regenerate this file with moq.
var _ keystore.Keystore = &KeystoreMock{}

// KeystoreMock is a mock implementation of keystore.Keystore.
//
//	func TestSomethingThatUsesKeystore(t *testing.T) {
//
//		// make and configure a mocked keystore.Keystore
//		mockedKeystore := &KeystoreMock{
//			CanSignMessageFunc: func(code coin.Code) bool {
//				panic("mock out the CanSignMessage method")
//			},
//			CanVerifyAddressFunc: func(coinMoqParam coin.Coin) (bool, bool, error) {
//				panic("mock out the CanVerifyAddress method")
//			},
//			CanVerifyExtendedPublicKeyFunc: func() bool {
//				panic("mock out the CanVerifyExtendedPublicKey method")
//			},
//			ExtendedPublicKeyFunc: func(coinMoqParam coin.Coin, absoluteKeypath signing.AbsoluteKeypath) (*hdkeychain.ExtendedKey, error) {
//				panic("mock out the ExtendedPublicKey method")
//			},
//			RootFingerprintFunc: func() ([]byte, error) {
//				panic("mock out the RootFingerprint method")
//			},
//			SignBTCMessageFunc: func(message []byte, keypath signing.AbsoluteKeypath, scriptType signing.ScriptType) ([]byte, error) {
//				panic("mock out the SignBTCMessage method")
//			},
//			SignETHMessageFunc: func(message []byte, keypath signing.AbsoluteKeypath) ([]byte, error) {
//				panic("mock out the SignETHMessage method")
//			},
//			SignETHTypedMessageFunc: func(chainID uint64, data []byte, keypath signing.AbsoluteKeypath) ([]byte, error) {
//				panic("mock out the SignETHTypedMessage method")
//			},
//			SignETHWalletConnectTransactionFunc: func(chainID uint64, tx *types.Transaction, keypath signing.AbsoluteKeypath) ([]byte, error) {
//				panic("mock out the SignETHWalletConnectTransaction method")
//			},
//			SignTransactionFunc: func(ifaceVal interface{}) error {
//				panic("mock out the SignTransaction method")
//			},
//			SupportsAccountFunc: func(coinInstance coin.Coin, meta interface{}) bool {
//				panic("mock out the SupportsAccount method")
//			},
//			SupportsCoinFunc: func(coinInstance coin.Coin) bool {
//				panic("mock out the SupportsCoin method")
//			},
//			SupportsMultipleAccountsFunc: func() bool {
//				panic("mock out the SupportsMultipleAccounts method")
//			},
//			SupportsUnifiedAccountsFunc: func() bool {
//				panic("mock out the SupportsUnifiedAccounts method")
//			},
//			TypeFunc: func() keystore.Type {
//				panic("mock out the Type method")
//			},
//			VerifyAddressFunc: func(configuration *signing.Configuration, coinMoqParam coin.Coin) error {
//				panic("mock out the VerifyAddress method")
//			},
//			VerifyExtendedPublicKeyFunc: func(coinMoqParam coin.Coin, configuration *signing.Configuration) error {
//				panic("mock out the VerifyExtendedPublicKey method")
//			},
//		}
//
//		// use mockedKeystore in code that requires keystore.Keystore
//		// and then make assertions.
//
//	}
type KeystoreMock struct {
	// CanSignMessageFunc mocks the CanSignMessage method.
	CanSignMessageFunc func(code coin.Code) bool

	// CanVerifyAddressFunc mocks the CanVerifyAddress method.
	CanVerifyAddressFunc func(coinMoqParam coin.Coin) (bool, bool, error)

	// CanVerifyExtendedPublicKeyFunc mocks the CanVerifyExtendedPublicKey method.
	CanVerifyExtendedPublicKeyFunc func() bool

	// ExtendedPublicKeyFunc mocks the ExtendedPublicKey method.
	ExtendedPublicKeyFunc func(coinMoqParam coin.Coin, absoluteKeypath signing.AbsoluteKeypath) (*hdkeychain.ExtendedKey, error)

	// RootFingerprintFunc mocks the RootFingerprint method.
	RootFingerprintFunc func() ([]byte, error)

	// SignBTCMessageFunc mocks the SignBTCMessage method.
	SignBTCMessageFunc func(message []byte, keypath signing.AbsoluteKeypath, scriptType signing.ScriptType) ([]byte, error)

	// SignETHMessageFunc mocks the SignETHMessage method.
	SignETHMessageFunc func(message []byte, keypath signing.AbsoluteKeypath) ([]byte, error)

	// SignETHTypedMessageFunc mocks the SignETHTypedMessage method.
	SignETHTypedMessageFunc func(chainID uint64, data []byte, keypath signing.AbsoluteKeypath) ([]byte, error)

	// SignETHWalletConnectTransactionFunc mocks the SignETHWalletConnectTransaction method.
	SignETHWalletConnectTransactionFunc func(chainID uint64, tx *types.Transaction, keypath signing.AbsoluteKeypath) ([]byte, error)

	// SignTransactionFunc mocks the SignTransaction method.
	SignTransactionFunc func(ifaceVal interface{}) error

	// SupportsAccountFunc mocks the SupportsAccount method.
	SupportsAccountFunc func(coinInstance coin.Coin, meta interface{}) bool

	// SupportsCoinFunc mocks the SupportsCoin method.
	SupportsCoinFunc func(coinInstance coin.Coin) bool

	// SupportsMultipleAccountsFunc mocks the SupportsMultipleAccounts method.
	SupportsMultipleAccountsFunc func() bool

	// SupportsUnifiedAccountsFunc mocks the SupportsUnifiedAccounts method.
	SupportsUnifiedAccountsFunc func() bool

	// TypeFunc mocks the Type method.
	TypeFunc func() keystore.Type

	// VerifyAddressFunc mocks the VerifyAddress method.
	VerifyAddressFunc func(configuration *signing.Configuration, coinMoqParam coin.Coin) error

	// VerifyExtendedPublicKeyFunc mocks the VerifyExtendedPublicKey method.
	VerifyExtendedPublicKeyFunc func(coinMoqParam coin.Coin, configuration *signing.Configuration) error

	// calls tracks calls to the methods.
	calls struct {
		// CanSignMessage holds details about calls to the CanSignMessage method.
		CanSignMessage []struct {
			// Code is the code argument value.
			Code coin.Code
		}
		// CanVerifyAddress holds details about calls to the CanVerifyAddress method.
		CanVerifyAddress []struct {
			// CoinMoqParam is the coinMoqParam argument value.
			CoinMoqParam coin.Coin
		}
		// CanVerifyExtendedPublicKey holds details about calls to the CanVerifyExtendedPublicKey method.
		CanVerifyExtendedPublicKey []struct {
		}
		// ExtendedPublicKey holds details about calls to the ExtendedPublicKey method.
		ExtendedPublicKey []struct {
			// CoinMoqParam is the coinMoqParam argument value.
			CoinMoqParam coin.Coin
			// AbsoluteKeypath is the absoluteKeypath argument value.
			AbsoluteKeypath signing.AbsoluteKeypath
		}
		// RootFingerprint holds details about calls to the RootFingerprint method.
		RootFingerprint []struct {
		}
		// SignBTCMessage holds details about calls to the SignBTCMessage method.
		SignBTCMessage []struct {
			// Message is the message argument value.
			Message []byte
			// Keypath is the keypath argument value.
			Keypath signing.AbsoluteKeypath
			// ScriptType is the scriptType argument value.
			ScriptType signing.ScriptType
		}
		// SignETHMessage holds details about calls to the SignETHMessage method.
		SignETHMessage []struct {
			// Message is the message argument value.
			Message []byte
			// Keypath is the keypath argument value.
			Keypath signing.AbsoluteKeypath
		}
		// SignETHTypedMessage holds details about calls to the SignETHTypedMessage method.
		SignETHTypedMessage []struct {
			// ChainID is the chainID argument value.
			ChainID uint64
			// Data is the data argument value.
			Data []byte
			// Keypath is the keypath argument value.
			Keypath signing.AbsoluteKeypath
		}
		// SignETHWalletConnectTransaction holds details about calls to the SignETHWalletConnectTransaction method.
		SignETHWalletConnectTransaction []struct {
			// ChainID is the chainID argument value.
			ChainID uint64
			// Tx is the tx argument value.
			Tx *types.Transaction
			// Keypath is the keypath argument value.
			Keypath signing.AbsoluteKeypath
		}
		// SignTransaction holds details about calls to the SignTransaction method.
		SignTransaction []struct {
			// IfaceVal is the ifaceVal argument value.
			IfaceVal interface{}
		}
		// SupportsAccount holds details about calls to the SupportsAccount method.
		SupportsAccount []struct {
			// CoinInstance is the coinInstance argument value.
			CoinInstance coin.Coin
			// Meta is the meta argument value.
			Meta interface{}
		}
		// SupportsCoin holds details about calls to the SupportsCoin method.
		SupportsCoin []struct {
			// CoinInstance is the coinInstance argument value.
			CoinInstance coin.Coin
		}
		// SupportsMultipleAccounts holds details about calls to the SupportsMultipleAccounts method.
		SupportsMultipleAccounts []struct {
		}
		// SupportsUnifiedAccounts holds details about calls to the SupportsUnifiedAccounts method.
		SupportsUnifiedAccounts []struct {
		}
		// Type holds details about calls to the Type method.
		Type []struct {
		}
		// VerifyAddress holds details about calls to the VerifyAddress method.
		VerifyAddress []struct {
			// Configuration is the configuration argument value.
			Configuration *signing.Configuration
			// CoinMoqParam is the coinMoqParam argument value.
			CoinMoqParam coin.Coin
		}
		// VerifyExtendedPublicKey holds details about calls to the VerifyExtendedPublicKey method.
		VerifyExtendedPublicKey []struct {
			// CoinMoqParam is the coinMoqParam argument value.
			CoinMoqParam coin.Coin
			// Configuration is the configuration argument value.
			Configuration *signing.Configuration
		}
	}
	lockCanSignMessage                  sync.RWMutex
	lockCanVerifyAddress                sync.RWMutex
	lockCanVerifyExtendedPublicKey      sync.RWMutex
	lockExtendedPublicKey               sync.RWMutex
	lockRootFingerprint                 sync.RWMutex
	lockSignBTCMessage                  sync.RWMutex
	lockSignETHMessage                  sync.RWMutex
	lockSignETHTypedMessage             sync.RWMutex
	lockSignETHWalletConnectTransaction sync.RWMutex
	lockSignTransaction                 sync.RWMutex
	lockSupportsAccount                 sync.RWMutex
	lockSupportsCoin                    sync.RWMutex
	lockSupportsMultipleAccounts        sync.RWMutex
	lockSupportsUnifiedAccounts         sync.RWMutex
	lockType                            sync.RWMutex
	lockVerifyAddress                   sync.RWMutex
	lockVerifyExtendedPublicKey         sync.RWMutex
}

// CanSignMessage calls CanSignMessageFunc.
func (mock *KeystoreMock) CanSignMessage(code coin.Code) bool {
	if mock.CanSignMessageFunc == nil {
		panic("KeystoreMock.CanSignMessageFunc: method is nil but Keystore.CanSignMessage was just called")
	}
	callInfo := struct {
		Code coin.Code
	}{
		Code: code,
	}
	mock.lockCanSignMessage.Lock()
	mock.calls.CanSignMessage = append(mock.calls.CanSignMessage, callInfo)
	mock.lockCanSignMessage.Unlock()
	return mock.CanSignMessageFunc(code)
}

// CanSignMessageCalls gets all the calls that were made to CanSignMessage.
// Check the length with:
//
//	len(mockedKeystore.CanSignMessageCalls())
func (mock *KeystoreMock) CanSignMessageCalls() []struct {
	Code coin.Code
} {
	var calls []struct {
		Code coin.Code
	}
	mock.lockCanSignMessage.RLock()
	calls = mock.calls.CanSignMessage
	mock.lockCanSignMessage.RUnlock()
	return calls
}

// CanVerifyAddress calls CanVerifyAddressFunc.
func (mock *KeystoreMock) CanVerifyAddress(coinMoqParam coin.Coin) (bool, bool, error) {
	if mock.CanVerifyAddressFunc == nil {
		panic("KeystoreMock.CanVerifyAddressFunc: method is nil but Keystore.CanVerifyAddress was just called")
	}
	callInfo := struct {
		CoinMoqParam coin.Coin
	}{
		CoinMoqParam: coinMoqParam,
	}
	mock.lockCanVerifyAddress.Lock()
	mock.calls.CanVerifyAddress = append(mock.calls.CanVerifyAddress, callInfo)
	mock.lockCanVerifyAddress.Unlock()
	return mock.CanVerifyAddressFunc(coinMoqParam)
}

// CanVerifyAddressCalls gets all the calls that were made to CanVerifyAddress.
// Check the length with:
//
//	len(mockedKeystore.CanVerifyAddressCalls())
func (mock *KeystoreMock) CanVerifyAddressCalls() []struct {
	CoinMoqParam coin.Coin
} {
	var calls []struct {
		CoinMoqParam coin.Coin
	}
	mock.lockCanVerifyAddress.RLock()
	calls = mock.calls.CanVerifyAddress
	mock.lockCanVerifyAddress.RUnlock()
	return calls
}

// CanVerifyExtendedPublicKey calls CanVerifyExtendedPublicKeyFunc.
func (mock *KeystoreMock) CanVerifyExtendedPublicKey() bool {
	if mock.CanVerifyExtendedPublicKeyFunc == nil {
		panic("KeystoreMock.CanVerifyExtendedPublicKeyFunc: method is nil but Keystore.CanVerifyExtendedPublicKey was just called")
	}
	callInfo := struct {
	}{}
	mock.lockCanVerifyExtendedPublicKey.Lock()
	mock.calls.CanVerifyExtendedPublicKey = append(mock.calls.CanVerifyExtendedPublicKey, callInfo)
	mock.lockCanVerifyExtendedPublicKey.Unlock()
	return mock.CanVerifyExtendedPublicKeyFunc()
}

// CanVerifyExtendedPublicKeyCalls gets all the calls that were made to CanVerifyExtendedPublicKey.
// Check the length with:
//
//	len(mockedKeystore.CanVerifyExtendedPublicKeyCalls())
func (mock *KeystoreMock) CanVerifyExtendedPublicKeyCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCanVerifyExtendedPublicKey.RLock()
	calls = mock.calls.CanVerifyExtendedPublicKey
	mock.lockCanVerifyExtendedPublicKey.RUnlock()
	return calls
}

// ExtendedPublicKey calls ExtendedPublicKeyFunc.
func (mock *KeystoreMock) ExtendedPublicKey(coinMoqParam coin.Coin, absoluteKeypath signing.AbsoluteKeypath) (*hdkeychain.ExtendedKey, error) {
	if mock.ExtendedPublicKeyFunc == nil {
		panic("KeystoreMock.ExtendedPublicKeyFunc: method is nil but Keystore.ExtendedPublicKey was just called")
	}
	callInfo := struct {
		CoinMoqParam    coin.Coin
		AbsoluteKeypath signing.AbsoluteKeypath
	}{
		CoinMoqParam:    coinMoqParam,
		AbsoluteKeypath: absoluteKeypath,
	}
	mock.lockExtendedPublicKey.Lock()
	mock.calls.ExtendedPublicKey = append(mock.calls.ExtendedPublicKey, callInfo)
	mock.lockExtendedPublicKey.Unlock()
	return mock.ExtendedPublicKeyFunc(coinMoqParam, absoluteKeypath)
}

// ExtendedPublicKeyCalls gets all the calls that were made to ExtendedPublicKey.
// Check the length with:
//
//	len(mockedKeystore.ExtendedPublicKeyCalls())
func (mock *KeystoreMock) ExtendedPublicKeyCalls() []struct {
	CoinMoqParam    coin.Coin
	AbsoluteKeypath signing.AbsoluteKeypath
} {
	var calls []struct {
		CoinMoqParam    coin.Coin
		AbsoluteKeypath signing.AbsoluteKeypath
	}
	mock.lockExtendedPublicKey.RLock()
	calls = mock.calls.ExtendedPublicKey
	mock.lockExtendedPublicKey.RUnlock()
	return calls
}

// RootFingerprint calls RootFingerprintFunc.
func (mock *KeystoreMock) RootFingerprint() ([]byte, error) {
	if mock.RootFingerprintFunc == nil {
		panic("KeystoreMock.RootFingerprintFunc: method is nil but Keystore.RootFingerprint was just called")
	}
	callInfo := struct {
	}{}
	mock.lockRootFingerprint.Lock()
	mock.calls.RootFingerprint = append(mock.calls.RootFingerprint, callInfo)
	mock.lockRootFingerprint.Unlock()
	return mock.RootFingerprintFunc()
}

// RootFingerprintCalls gets all the calls that were made to RootFingerprint.
// Check the length with:
//
//	len(mockedKeystore.RootFingerprintCalls())
func (mock *KeystoreMock) RootFingerprintCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockRootFingerprint.RLock()
	calls = mock.calls.RootFingerprint
	mock.lockRootFingerprint.RUnlock()
	return calls
}

// SignBTCMessage calls SignBTCMessageFunc.
func (mock *KeystoreMock) SignBTCMessage(message []byte, keypath signing.AbsoluteKeypath, scriptType signing.ScriptType) ([]byte, error) {
	if mock.SignBTCMessageFunc == nil {
		panic("KeystoreMock.SignBTCMessageFunc: method is nil but Keystore.SignBTCMessage was just called")
	}
	callInfo := struct {
		Message    []byte
		Keypath    signing.AbsoluteKeypath
		ScriptType signing.ScriptType
	}{
		Message:    message,
		Keypath:    keypath,
		ScriptType: scriptType,
	}
	mock.lockSignBTCMessage.Lock()
	mock.calls.SignBTCMessage = append(mock.calls.SignBTCMessage, callInfo)
	mock.lockSignBTCMessage.Unlock()
	return mock.SignBTCMessageFunc(message, keypath, scriptType)
}

// SignBTCMessageCalls gets all the calls that were made to SignBTCMessage.
// Check the length with:
//
//	len(mockedKeystore.SignBTCMessageCalls())
func (mock *KeystoreMock) SignBTCMessageCalls() []struct {
	Message    []byte
	Keypath    signing.AbsoluteKeypath
	ScriptType signing.ScriptType
} {
	var calls []struct {
		Message    []byte
		Keypath    signing.AbsoluteKeypath
		ScriptType signing.ScriptType
	}
	mock.lockSignBTCMessage.RLock()
	calls = mock.calls.SignBTCMessage
	mock.lockSignBTCMessage.RUnlock()
	return calls
}

// SignETHMessage calls SignETHMessageFunc.
func (mock *KeystoreMock) SignETHMessage(message []byte, keypath signing.AbsoluteKeypath) ([]byte, error) {
	if mock.SignETHMessageFunc == nil {
		panic("KeystoreMock.SignETHMessageFunc: method is nil but Keystore.SignETHMessage was just called")
	}
	callInfo := struct {
		Message []byte
		Keypath signing.AbsoluteKeypath
	}{
		Message: message,
		Keypath: keypath,
	}
	mock.lockSignETHMessage.Lock()
	mock.calls.SignETHMessage = append(mock.calls.SignETHMessage, callInfo)
	mock.lockSignETHMessage.Unlock()
	return mock.SignETHMessageFunc(message, keypath)
}

// SignETHMessageCalls gets all the calls that were made to SignETHMessage.
// Check the length with:
//
//	len(mockedKeystore.SignETHMessageCalls())
func (mock *KeystoreMock) SignETHMessageCalls() []struct {
	Message []byte
	Keypath signing.AbsoluteKeypath
} {
	var calls []struct {
		Message []byte
		Keypath signing.AbsoluteKeypath
	}
	mock.lockSignETHMessage.RLock()
	calls = mock.calls.SignETHMessage
	mock.lockSignETHMessage.RUnlock()
	return calls
}

// SignETHTypedMessage calls SignETHTypedMessageFunc.
func (mock *KeystoreMock) SignETHTypedMessage(chainID uint64, data []byte, keypath signing.AbsoluteKeypath) ([]byte, error) {
	if mock.SignETHTypedMessageFunc == nil {
		panic("KeystoreMock.SignETHTypedMessageFunc: method is nil but Keystore.SignETHTypedMessage was just called")
	}
	callInfo := struct {
		ChainID uint64
		Data    []byte
		Keypath signing.AbsoluteKeypath
	}{
		ChainID: chainID,
		Data:    data,
		Keypath: keypath,
	}
	mock.lockSignETHTypedMessage.Lock()
	mock.calls.SignETHTypedMessage = append(mock.calls.SignETHTypedMessage, callInfo)
	mock.lockSignETHTypedMessage.Unlock()
	return mock.SignETHTypedMessageFunc(chainID, data, keypath)
}

// SignETHTypedMessageCalls gets all the calls that were made to SignETHTypedMessage.
// Check the length with:
//
//	len(mockedKeystore.SignETHTypedMessageCalls())
func (mock *KeystoreMock) SignETHTypedMessageCalls() []struct {
	ChainID uint64
	Data    []byte
	Keypath signing.AbsoluteKeypath
} {
	var calls []struct {
		ChainID uint64
		Data    []byte
		Keypath signing.AbsoluteKeypath
	}
	mock.lockSignETHTypedMessage.RLock()
	calls = mock.calls.SignETHTypedMessage
	mock.lockSignETHTypedMessage.RUnlock()
	return calls
}

// SignETHWalletConnectTransaction calls SignETHWalletConnectTransactionFunc.
func (mock *KeystoreMock) SignETHWalletConnectTransaction(chainID uint64, tx *types.Transaction, keypath signing.AbsoluteKeypath) ([]byte, error) {
	if mock.SignETHWalletConnectTransactionFunc == nil {
		panic("KeystoreMock.SignETHWalletConnectTransactionFunc: method is nil but Keystore.SignETHWalletConnectTransaction was just called")
	}
	callInfo := struct {
		ChainID uint64
		Tx      *types.Transaction
		Keypath signing.AbsoluteKeypath
	}{
		ChainID: chainID,
		Tx:      tx,
		Keypath: keypath,
	}
	mock.lockSignETHWalletConnectTransaction.Lock()
	mock.calls.SignETHWalletConnectTransaction = append(mock.calls.SignETHWalletConnectTransaction, callInfo)
	mock.lockSignETHWalletConnectTransaction.Unlock()
	return mock.SignETHWalletConnectTransactionFunc(chainID, tx, keypath)
}

// SignETHWalletConnectTransactionCalls gets all the calls that were made to SignETHWalletConnectTransaction.
// Check the length with:
//
//	len(mockedKeystore.SignETHWalletConnectTransactionCalls())
func (mock *KeystoreMock) SignETHWalletConnectTransactionCalls() []struct {
	ChainID uint64
	Tx      *types.Transaction
	Keypath signing.AbsoluteKeypath
} {
	var calls []struct {
		ChainID uint64
		Tx      *types.Transaction
		Keypath signing.AbsoluteKeypath
	}
	mock.lockSignETHWalletConnectTransaction.RLock()
	calls = mock.calls.SignETHWalletConnectTransaction
	mock.lockSignETHWalletConnectTransaction.RUnlock()
	return calls
}

// SignTransaction calls SignTransactionFunc.
func (mock *KeystoreMock) SignTransaction(ifaceVal interface{}) error {
	if mock.SignTransactionFunc == nil {
		panic("KeystoreMock.SignTransactionFunc: method is nil but Keystore.SignTransaction was just called")
	}
	callInfo := struct {
		IfaceVal interface{}
	}{
		IfaceVal: ifaceVal,
	}
	mock.lockSignTransaction.Lock()
	mock.calls.SignTransaction = append(mock.calls.SignTransaction, callInfo)
	mock.lockSignTransaction.Unlock()
	return mock.SignTransactionFunc(ifaceVal)
}

// SignTransactionCalls gets all the calls that were made to SignTransaction.
// Check the length with:
//
//	len(mockedKeystore.SignTransactionCalls())
func (mock *KeystoreMock) SignTransactionCalls() []struct {
	IfaceVal interface{}
} {
	var calls []struct {
		IfaceVal interface{}
	}
	mock.lockSignTransaction.RLock()
	calls = mock.calls.SignTransaction
	mock.lockSignTransaction.RUnlock()
	return calls
}

// SupportsAccount calls SupportsAccountFunc.
func (mock *KeystoreMock) SupportsAccount(coinInstance coin.Coin, meta interface{}) bool {
	if mock.SupportsAccountFunc == nil {
		panic("KeystoreMock.SupportsAccountFunc: method is nil but Keystore.SupportsAccount was just called")
	}
	callInfo := struct {
		CoinInstance coin.Coin
		Meta         interface{}
	}{
		CoinInstance: coinInstance,
		Meta:         meta,
	}
	mock.lockSupportsAccount.Lock()
	mock.calls.SupportsAccount = append(mock.calls.SupportsAccount, callInfo)
	mock.lockSupportsAccount.Unlock()
	return mock.SupportsAccountFunc(coinInstance, meta)
}

// SupportsAccountCalls gets all the calls that were made to SupportsAccount.
// Check the length with:
//
//	len(mockedKeystore.SupportsAccountCalls())
func (mock *KeystoreMock) SupportsAccountCalls() []struct {
	CoinInstance coin.Coin
	Meta         interface{}
} {
	var calls []struct {
		CoinInstance coin.Coin
		Meta         interface{}
	}
	mock.lockSupportsAccount.RLock()
	calls = mock.calls.SupportsAccount
	mock.lockSupportsAccount.RUnlock()
	return calls
}

// SupportsCoin calls SupportsCoinFunc.
func (mock *KeystoreMock) SupportsCoin(coinInstance coin.Coin) bool {
	if mock.SupportsCoinFunc == nil {
		panic("KeystoreMock.SupportsCoinFunc: method is nil but Keystore.SupportsCoin was just called")
	}
	callInfo := struct {
		CoinInstance coin.Coin
	}{
		CoinInstance: coinInstance,
	}
	mock.lockSupportsCoin.Lock()
	mock.calls.SupportsCoin = append(mock.calls.SupportsCoin, callInfo)
	mock.lockSupportsCoin.Unlock()
	return mock.SupportsCoinFunc(coinInstance)
}

// SupportsCoinCalls gets all the calls that were made to SupportsCoin.
// Check the length with:
//
//	len(mockedKeystore.SupportsCoinCalls())
func (mock *KeystoreMock) SupportsCoinCalls() []struct {
	CoinInstance coin.Coin
} {
	var calls []struct {
		CoinInstance coin.Coin
	}
	mock.lockSupportsCoin.RLock()
	calls = mock.calls.SupportsCoin
	mock.lockSupportsCoin.RUnlock()
	return calls
}

// SupportsMultipleAccounts calls SupportsMultipleAccountsFunc.
func (mock *KeystoreMock) SupportsMultipleAccounts() bool {
	if mock.SupportsMultipleAccountsFunc == nil {
		panic("KeystoreMock.SupportsMultipleAccountsFunc: method is nil but Keystore.SupportsMultipleAccounts was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSupportsMultipleAccounts.Lock()
	mock.calls.SupportsMultipleAccounts = append(mock.calls.SupportsMultipleAccounts, callInfo)
	mock.lockSupportsMultipleAccounts.Unlock()
	return mock.SupportsMultipleAccountsFunc()
}

// SupportsMultipleAccountsCalls gets all the calls that were made to SupportsMultipleAccounts.
// Check the length with:
//
//	len(mockedKeystore.SupportsMultipleAccountsCalls())
func (mock *KeystoreMock) SupportsMultipleAccountsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSupportsMultipleAccounts.RLock()
	calls = mock.calls.SupportsMultipleAccounts
	mock.lockSupportsMultipleAccounts.RUnlock()
	return calls
}

// SupportsUnifiedAccounts calls SupportsUnifiedAccountsFunc.
func (mock *KeystoreMock) SupportsUnifiedAccounts() bool {
	if mock.SupportsUnifiedAccountsFunc == nil {
		panic("KeystoreMock.SupportsUnifiedAccountsFunc: method is nil but Keystore.SupportsUnifiedAccounts was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSupportsUnifiedAccounts.Lock()
	mock.calls.SupportsUnifiedAccounts = append(mock.calls.SupportsUnifiedAccounts, callInfo)
	mock.lockSupportsUnifiedAccounts.Unlock()
	return mock.SupportsUnifiedAccountsFunc()
}

// SupportsUnifiedAccountsCalls gets all the calls that were made to SupportsUnifiedAccounts.
// Check the length with:
//
//	len(mockedKeystore.SupportsUnifiedAccountsCalls())
func (mock *KeystoreMock) SupportsUnifiedAccountsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSupportsUnifiedAccounts.RLock()
	calls = mock.calls.SupportsUnifiedAccounts
	mock.lockSupportsUnifiedAccounts.RUnlock()
	return calls
}

// Type calls TypeFunc.
func (mock *KeystoreMock) Type() keystore.Type {
	if mock.TypeFunc == nil {
		panic("KeystoreMock.TypeFunc: method is nil but Keystore.Type was just called")
	}
	callInfo := struct {
	}{}
	mock.lockType.Lock()
	mock.calls.Type = append(mock.calls.Type, callInfo)
	mock.lockType.Unlock()
	return mock.TypeFunc()
}

// TypeCalls gets all the calls that were made to Type.
// Check the length with:
//
//	len(mockedKeystore.TypeCalls())
func (mock *KeystoreMock) TypeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockType.RLock()
	calls = mock.calls.Type
	mock.lockType.RUnlock()
	return calls
}

// VerifyAddress calls VerifyAddressFunc.
func (mock *KeystoreMock) VerifyAddress(configuration *signing.Configuration, coinMoqParam coin.Coin) error {
	if mock.VerifyAddressFunc == nil {
		panic("KeystoreMock.VerifyAddressFunc: method is nil but Keystore.VerifyAddress was just called")
	}
	callInfo := struct {
		Configuration *signing.Configuration
		CoinMoqParam  coin.Coin
	}{
		Configuration: configuration,
		CoinMoqParam:  coinMoqParam,
	}
	mock.lockVerifyAddress.Lock()
	mock.calls.VerifyAddress = append(mock.calls.VerifyAddress, callInfo)
	mock.lockVerifyAddress.Unlock()
	return mock.VerifyAddressFunc(configuration, coinMoqParam)
}

// VerifyAddressCalls gets all the calls that were made to VerifyAddress.
// Check the length with:
//
//	len(mockedKeystore.VerifyAddressCalls())
func (mock *KeystoreMock) VerifyAddressCalls() []struct {
	Configuration *signing.Configuration
	CoinMoqParam  coin.Coin
} {
	var calls []struct {
		Configuration *signing.Configuration
		CoinMoqParam  coin.Coin
	}
	mock.lockVerifyAddress.RLock()
	calls = mock.calls.VerifyAddress
	mock.lockVerifyAddress.RUnlock()
	return calls
}

// VerifyExtendedPublicKey calls VerifyExtendedPublicKeyFunc.
func (mock *KeystoreMock) VerifyExtendedPublicKey(coinMoqParam coin.Coin, configuration *signing.Configuration) error {
	if mock.VerifyExtendedPublicKeyFunc == nil {
		panic("KeystoreMock.VerifyExtendedPublicKeyFunc: method is nil but Keystore.VerifyExtendedPublicKey was just called")
	}
	callInfo := struct {
		CoinMoqParam  coin.Coin
		Configuration *signing.Configuration
	}{
		CoinMoqParam:  coinMoqParam,
		Configuration: configuration,
	}
	mock.lockVerifyExtendedPublicKey.Lock()
	mock.calls.VerifyExtendedPublicKey = append(mock.calls.VerifyExtendedPublicKey, callInfo)
	mock.lockVerifyExtendedPublicKey.Unlock()
	return mock.VerifyExtendedPublicKeyFunc(coinMoqParam, configuration)
}

// VerifyExtendedPublicKeyCalls gets all the calls that were made to VerifyExtendedPublicKey.
// Check the length with:
//
//	len(mockedKeystore.VerifyExtendedPublicKeyCalls())
func (mock *KeystoreMock) VerifyExtendedPublicKeyCalls() []struct {
	CoinMoqParam  coin.Coin
	Configuration *signing.Configuration
} {
	var calls []struct {
		CoinMoqParam  coin.Coin
		Configuration *signing.Configuration
	}
	mock.lockVerifyExtendedPublicKey.RLock()
	calls = mock.calls.VerifyExtendedPublicKey
	mock.lockVerifyExtendedPublicKey.RUnlock()
	return calls
}
