package accounts

import (
	"github.com/ledgerwatch/turbo-geth/crypto"
	"math/big"

	"bytes"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/rlp"
)

type ExtAccount struct {
	Nonce   uint64
	Balance *big.Int
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce       uint64
	Balance     *big.Int
	Root        common.Hash // merkle root of the storage trie
	CodeHash    []byte
	StorageSize *uint64
}

type accountWithoutStorage struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
}

const (
	accountSizeWithoutData            = 1
	minAccountSizeWithRootAndCodeHash = 60
)

var emptyCodeHash = crypto.Keccak256(nil)
var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

func (a *Account) Encode(enableStorageSize bool) ([]byte, error) {
	var toEncode interface{}

	if a.IsEmptyCodeHash() && a.IsEmptyRoot() {
		if (a.Balance == nil || a.Balance.Sign() == 0) && a.Nonce == 0 {
			return []byte{byte(192)}, nil
		}

		toEncode = new(ExtAccount).
			fill(a).
			setDefaultBalance()
	} else {
		acc := newAccountCopy(a)
		toEncode = acc

		if !enableStorageSize || acc.StorageSize == nil {
			toEncode = &accountWithoutStorage{
				Nonce:    acc.Nonce,
				Balance:  acc.Balance,
				Root:     acc.Root,
				CodeHash: acc.CodeHash,
			}
		}
	}

	return rlp.EncodeToBytes(toEncode)
}

func (a *Account) EncodeRLP(enableStorageSize bool) ([]byte, error) {
	acc := newAccountCopy(a)
	toEncode := interface{}(acc)

	if !enableStorageSize {
		toEncode = &accountWithoutStorage{
			Nonce:    acc.Nonce,
			Balance:  acc.Balance,
			Root:     acc.Root,
			CodeHash: acc.CodeHash,
		}
	}

	return rlp.EncodeToBytes(toEncode)
}

func (a *Account) Decode(enc []byte) error {
	switch encodedLength := len(enc); {
	case encodedLength == 0:

	case encodedLength == accountSizeWithoutData:
		a.Balance = new(big.Int)
		a.setCodeHash(emptyCodeHash)
		a.Root = emptyRoot

	case encodedLength < minAccountSizeWithRootAndCodeHash:
		var extData ExtAccount
		if err := rlp.DecodeBytes(enc, &extData); err != nil {
			return err
		}

		a.fillFromExtAccount(extData)
	default:
		dataWithoutStorage := &accountWithoutStorage{}
		err := rlp.DecodeBytes(enc, dataWithoutStorage)
		if err == nil {
			a.fillAccountWithoutStorage(dataWithoutStorage)
			return nil
		}

		if err.Error() != "rlp: input list has too many elements for accounts.accountWithoutStorage" {
			return err
		}

		dataWithStorage := &Account{}
		if err := rlp.DecodeBytes(enc, &dataWithStorage); err != nil {
			return err
		}

		a.fill(dataWithStorage)
	}

	return nil
}

func Decode(enc []byte) (*Account, error) {
	if len(enc) == 0 {
		return nil, nil
	}

	acc := new(Account)
	err := acc.Decode(enc)
	return acc, err
}

func newAccountCopy(srcAccount *Account) *Account {
	return new(Account).
		fill(srcAccount).
		setDefaultBalance().
		setDefaultCodeHash().
		setDefaultRoot()
}

func (a *Account) fill(srcAccount *Account) *Account {
	a.Root = srcAccount.Root

	a.CodeHash = make([]byte, len(srcAccount.CodeHash))
	copy(a.CodeHash, srcAccount.CodeHash)

	a.setDefaultBalance()
	a.Balance.Set(srcAccount.Balance)

	a.Nonce = srcAccount.Nonce

	if srcAccount.StorageSize != nil {
		a.StorageSize = new(uint64)
		*a.StorageSize = *srcAccount.StorageSize
	}

	return a
}

func (a *Account) setCodeHash(codeHash []byte) {
	a.CodeHash = make([]byte, len(codeHash))
	copy(a.CodeHash, codeHash)
}

func (a *Account) fillAccountWithoutStorage(srcAccount *accountWithoutStorage) *Account {
	a.Root = srcAccount.Root

	a.setCodeHash(srcAccount.CodeHash)

	a.setDefaultBalance()
	a.Balance.Set(srcAccount.Balance)

	a.Nonce = srcAccount.Nonce

	a.StorageSize = nil

	return a
}

func (a *Account) fillFromExtAccount(srcExtAccount ExtAccount) *Account {
	a.Nonce = srcExtAccount.Nonce

	a.setDefaultBalance()
	a.Balance.Set(srcExtAccount.Balance)

	a.CodeHash = emptyCodeHash

	a.Root = emptyRoot

	return a
}

func (a *Account) setDefaultBalance() *Account {
	if a.Balance == nil {
		a.Balance = new(big.Int)
	}

	return a
}

func (a *Account) setDefaultCodeHash() *Account {
	if a.IsEmptyCodeHash() {
		a.CodeHash = emptyCodeHash
	}

	return a
}

func (a *Account) setDefaultRoot() *Account {
	if a.IsEmptyRoot() {
		a.Root = emptyRoot
	}

	return a
}

func (a *Account) IsEmptyCodeHash() bool {
	return a.CodeHash == nil || bytes.Equal(a.CodeHash[:], emptyCodeHash)
}

func (a *Account) IsEmptyRoot() bool {
	return a.Root == emptyRoot || a.Root == common.Hash{}
}

func (extAcc *ExtAccount) fill(srcAccount *Account) *ExtAccount {
	extAcc.setDefaultBalance()
	extAcc.Balance.Set(srcAccount.Balance)

	extAcc.Nonce = srcAccount.Nonce

	return extAcc
}

func (extAcc *ExtAccount) setDefaultBalance() *ExtAccount {
	if extAcc.Balance == nil {
		extAcc.Balance = new(big.Int)
	}

	return extAcc
}
