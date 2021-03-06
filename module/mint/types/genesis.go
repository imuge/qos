package types

import (
	"errors"
	"github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
	"time"
)

// 创世状态
type GenesisState struct {
	InflationPhrases InflationPhrases `json:"inflation_phrases"`  // 通胀阶段
	FirstBlockTime   int64            `json:"first_block_time"`   // UTC().UNIX()
	AppliedQOSAmount types.BigInt     `json:"applied_qos_amount"` // 总流通QOS
	TotalQOSAmount   types.BigInt     `json:"total_qos_amount"`   // QOS总量
}

func NewGenesisState(inflationPhrases InflationPhrases, firstBlockTime int64, appliedQOSAmount, totalQOSAmount types.BigInt) GenesisState {
	return GenesisState{
		InflationPhrases: inflationPhrases,
		FirstBlockTime:   firstBlockTime,
		AppliedQOSAmount: appliedQOSAmount,
		TotalQOSAmount:   totalQOSAmount,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultInflationPhrases(), time.Now().Unix(), types.ZeroInt(), types.NewInt(qtypes.DefaultTotalUnitQOSQuantity))
}

func ValidateGenesis(gs GenesisState) error {
	if !gs.TotalQOSAmount.GT(types.ZeroInt()) {
		return errors.New("total_qos_amount must be positive")
	}

	if !(gs.AppliedQOSAmount.Add(gs.InflationPhrases.TotalAmount())).Equal(gs.TotalQOSAmount) {
		return errors.New("total_qos_amount must equals applied_qos_amount + sum(total_amount in inflation_phrases)")
	}

	return gs.InflationPhrases.Valid()
}
