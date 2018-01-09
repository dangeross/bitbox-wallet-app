package maketx

import (
	"errors"
	"sort"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/txsort"
	"github.com/shiftdevices/godbb/util/errp"
)

type byValue struct {
	outPoints []wire.OutPoint
	outputs   map[wire.OutPoint]*wire.TxOut
}

func (p *byValue) Len() int { return len(p.outPoints) }
func (p *byValue) Less(i, j int) bool {
	return p.outputs[p.outPoints[i]].Value < p.outputs[p.outPoints[j]].Value
}
func (p *byValue) Swap(i, j int) { p.outPoints[i], p.outPoints[j] = p.outPoints[j], p.outPoints[i] }

func coinSelection(
	minAmount btcutil.Amount,
	outputs map[wire.OutPoint]*wire.TxOut,
) (btcutil.Amount, []wire.OutPoint, error) {
	outPoints := []wire.OutPoint{}
	for outPoint := range outputs {
		outPoints = append(outPoints, outPoint)
	}
	sort.Sort(sort.Reverse(&byValue{outPoints, outputs}))
	selectedOutPoints := []wire.OutPoint{}
	outputsSum := btcutil.Amount(0)

	for _, outPoint := range outPoints {
		if outputsSum >= minAmount {
			break
		}
		selectedOutPoints = append(selectedOutPoints, outPoint)
		outputsSum += btcutil.Amount(outputs[outPoint].Value)
	}
	if outputsSum < minAmount {
		return 0, nil, errp.New("insufficient funds")
	}
	return outputsSum, selectedOutPoints, nil
}

// NewTxSpendAll creates a transaction which spends all available unspent outputs.
func NewTxSpendAll(
	spendableOutputs map[wire.OutPoint]*wire.TxOut,
	outputPkScript []byte,
	feePerKb btcutil.Amount) (btcutil.Amount, *wire.MsgTx, []wire.OutPoint, error) {

	selectedOutPoints := []wire.OutPoint{}
	inputs := []*wire.TxIn{}
	outputsSum := btcutil.Amount(0)
	for outPoint, output := range spendableOutputs {
		outPoint := outPoint // avoid reference reuse due to range loop
		selectedOutPoints = append(selectedOutPoints, outPoint)
		outputsSum += btcutil.Amount(output.Value)
		inputs = append(inputs, wire.NewTxIn(&outPoint, nil, nil))
	}
	output := wire.NewTxOut(0, outputPkScript)
	txSize := EstimateSerializeSize(len(selectedOutPoints), []*wire.TxOut{output}, false)
	maxRequiredFee := FeeForSerializeSize(feePerKb, txSize)
	if outputsSum < maxRequiredFee {
		return 0, nil, nil, errp.New("insufficient funds for fee")
	}
	output = wire.NewTxOut(int64(outputsSum-maxRequiredFee), outputPkScript)
	unsignedTransaction := &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     inputs,
		TxOut:    []*wire.TxOut{output},
		LockTime: 0,
	}
	txsort.InPlaceSort(unsignedTransaction)
	return btcutil.Amount(output.Value), unsignedTransaction, selectedOutPoints, nil
}

// NewTx creates a transaction from a set of unspent outputs, targeting an output value. A subset of
// the unspent outputs is selected to cover the needed amount. A change output is added if needed.
func NewTx(
	spendableOutputs map[wire.OutPoint]*wire.TxOut,
	output *wire.TxOut,
	feePerKb btcutil.Amount,
	getChangePKScript func() ([]byte, error),
) (*wire.MsgTx, []wire.OutPoint, error) {
	targetAmount := btcutil.Amount(output.Value)
	outputs := []*wire.TxOut{output}
	estimatedSize := EstimateSerializeSize(1, outputs, true)
	targetFee := FeeForSerializeSize(feePerKb, estimatedSize)

	for {
		selectedOutputsSum, selectedOutPoints, err := coinSelection(
			targetAmount+targetFee,
			spendableOutputs,
		)
		if err != nil {
			return nil, nil, err
		}

		txSize := EstimateSerializeSize(len(selectedOutPoints), outputs, true)
		maxRequiredFee := FeeForSerializeSize(feePerKb, txSize)
		if selectedOutputsSum-targetAmount < maxRequiredFee {
			targetFee = maxRequiredFee
			continue
		}

		inputs := make([]*wire.TxIn, len(selectedOutPoints))
		for i, outPoint := range selectedOutPoints {
			inputs[i] = wire.NewTxIn(&outPoint, nil, nil)
		}
		unsignedTransaction := &wire.MsgTx{
			Version:  wire.TxVersion,
			TxIn:     inputs,
			TxOut:    outputs,
			LockTime: 0,
		}
		changeAmount := selectedOutputsSum - targetAmount - maxRequiredFee
		if changeAmount != 0 && !IsDustAmount(changeAmount, P2PKHPkScriptSize, feePerKb) {
			changePKScript, err := getChangePKScript()
			if err != nil {
				return nil, nil, err
			}
			if len(changePKScript) > P2PKHPkScriptSize {
				return nil, nil, errors.New("fee estimation requires change " +
					"scripts no larger than P2PKH output scripts")
			}
			changeOutput := wire.NewTxOut(int64(changeAmount), changePKScript)
			unsignedTransaction.TxOut = append(unsignedTransaction.TxOut, changeOutput)
		}
		txsort.InPlaceSort(unsignedTransaction)
		return unsignedTransaction, selectedOutPoints, nil
	}
}
