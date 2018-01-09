package deterministicwallet

import (
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/txsort"

	"github.com/shiftdevices/godbb/deterministicwallet/addresses"
	"github.com/shiftdevices/godbb/deterministicwallet/transactions"
	"github.com/shiftdevices/godbb/util/errp"
)

// SignTransaction signs all inputs. It assumes all outputs spent belong to this wallet.
func SignTransaction(
	sign HDKeyStoreInterface,
	transaction *wire.MsgTx,
	previousOutputs []*transactions.TxOut,
) error {
	if len(previousOutputs) != len(transaction.TxIn) {
		panic("output/input mismatch; there needs to be exactly one output being spent ber input")
	}
	signatureHashes := [][]byte{}
	keyPaths := []string{}
	sigHashes := txscript.NewTxSigHashes(transaction)
	for index := range transaction.TxIn {
		spentOutput := previousOutputs[index]
		address := spentOutput.Address.(*addresses.Address)
		isSegwit, subScript := address.SigHashData()
		var signatureHash []byte
		if isSegwit {
			var err error
			signatureHash, err = txscript.CalcWitnessSigHash(
				subScript, sigHashes, txscript.SigHashAll, transaction, index, spentOutput.Value)
			if err != nil {
				return err
			}
		} else {
			var err error
			signatureHash, err = txscript.CalcSignatureHash(
				subScript, txscript.SigHashAll, transaction, index)
			if err != nil {
				return err
			}
		}

		signatureHashes = append(signatureHashes, signatureHash)
		keyPaths = append(keyPaths, spentOutput.Address.(*addresses.Address).KeyPath)
	}
	signatures, err := sign.Sign(signatureHashes, keyPaths)
	if err != nil {
		return err
	}
	if len(signatures) != len(transaction.TxIn) {
		panic("number of signatures doesn't match number of inputs")
	}
	for index, input := range transaction.TxIn {
		spentOutput := previousOutputs[index]
		address := spentOutput.Address.(*addresses.Address)
		signature := signatures[index]
		input.SignatureScript, input.Witness = address.InputData(signature)
	}
	// Sanity check: see if the created transaction is valid.
	if err := txValidityCheck(transaction, previousOutputs, sigHashes); err != nil {
		panic(err)
	}
	return nil
}

func txValidityCheck(transaction *wire.MsgTx, previousOutputs []*transactions.TxOut,
	sigHashes *txscript.TxSigHashes) error {
	if !txsort.IsSorted(transaction) {
		return errp.New("tx not bip69 conformant")
	}
	for index := range transaction.TxIn {
		engine, err := txscript.NewEngine(
			previousOutputs[index].PkScript,
			transaction,
			index,
			txscript.StandardVerifyFlags, nil, sigHashes, previousOutputs[index].Value)
		if err != nil {
			return errp.WithStack(err)
		}
		if err := engine.Execute(); err != nil {
			return errp.WithStack(err)
		}
	}
	return nil
}
