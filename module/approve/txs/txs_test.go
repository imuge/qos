package txs

import (
	"github.com/QOSGroup/qos/module/bank"
	"testing"

	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/types"
	"github.com/QOSGroup/qos/module/qsc"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var testFromAddr = btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
var testToAddr = btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())

func genTestApprove() types.Approve {
	return types.Approve{
		From: testFromAddr,
		To:   testToAddr,
		QOS:  btypes.NewInt(100),
		QSCs: qtypes.QSCs{
			{
				Name:   "QSTAR",
				Amount: btypes.NewInt(100),
			},
		},
	}
}

func TestValidateData(t *testing.T) {
	ctx := defaultContext()
	approve := genTestApprove()
	require.NotNil(t, ValidateData(ctx, approve))

	bank.GetMapper(ctx).SetAccount(qtypes.NewQOSAccountWithAddress(testFromAddr))
	require.NotNil(t, ValidateData(ctx, approve))

	saveQSCInfo(ctx, "QSTAR")
	require.Nil(t, ValidateData(ctx, approve))

	approve.QSCs = append(approve.QSCs, &qtypes.QSC{
		Name:   "QSTAR",
		Amount: btypes.NewInt(100),
	})
	require.NotNil(t, ValidateData(ctx, approve))

	approve.QSCs = qtypes.QSCs{
		{
			Name:   "QOS",
			Amount: btypes.NewInt(100),
		},
	}
	require.NotNil(t, ValidateData(ctx, approve))
}

func genTestAccount(addr btypes.AccAddress) *qtypes.QOSAccount {
	return &qtypes.QOSAccount{
		BaseAccount: bacc.BaseAccount{
			AccountAddress: addr,
			Publickey:      nil,
			Nonce:          0,
		},
		QOS: btypes.NewInt(100),
		QSCs: qtypes.QSCs{
			{
				Name:   "QSTAR",
				Amount: btypes.NewInt(100),
			},
		},
	}
}

func genApproveCancelTx() TxCancelApprove {
	return TxCancelApprove{
		From: btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		To:   btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
	}
}

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	approveMapper := mapper.NewApproveMapper()
	approveMapper.SetCodec(Cdc)
	approveKey := approveMapper.GetStoreKey()
	mapperMap[types.MapperName] = approveMapper

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(Cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	qscMapper := qsc.NewMapper()
	qscMapper.SetCodec(Cdc)
	qscKey := qscMapper.GetStoreKey()
	mapperMap[qsc.MapperName] = qscMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(qscKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func saveQSCInfo(ctx context.Context, qscName string) {
	qscMapper := qsc.GetMapper(ctx)
	qscMapper.SaveQsc(&qsc.Info{
		Name: qscName,
	})
}

func defaultContextWithQSC() context.Context {
	ctx := defaultContext()
	saveQSCInfo(ctx, "QSTAR")
	return ctx
}

func TestTxApproveCreate_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()

	tx := TxCreateApprove{
		genTestApprove(),
	}
	require.NotNil(t, tx.ValidateData(ctx))

	bank.GetMapper(ctx).SetAccount(qtypes.NewQOSAccountWithAddress(testFromAddr))
	require.Nil(t, tx.ValidateData(ctx))

	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(tx.Approve)

	require.NotNil(t, tx.ValidateData(ctx))
}

func TestTxApproveCreate_Exec(t *testing.T) {
	ctx := defaultContext()

	tx := TxCreateApprove{
		genTestApprove(),
	}
	result, cross := tx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approveMapper := mapper.GetMapper(ctx)
	approve, exists := approveMapper.GetApprove(tx.From, tx.To)
	require.True(t, exists)
	require.True(t, tx.Approve.Equals(approve))
}

func TestTxApproveIncrease_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()
	bank.GetMapper(ctx).SetAccount(qtypes.NewQOSAccountWithAddress(testFromAddr))

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	increaseTx := TxIncreaseApprove{
		genTestApprove(),
	}
	require.NotNil(t, increaseTx.ValidateData(ctx))

	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(createTx.Approve)

	require.Nil(t, increaseTx.ValidateData(ctx))
}

func TestTxApproveIncrease_Exec(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	increaseTx := TxIncreaseApprove{
		genTestApprove(),
	}

	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(createTx.Approve)

	result, cross := increaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(createTx.From, createTx.To)
	require.True(t, exists)
	require.True(t, createTx.Approve.Plus(increaseTx.QOS, increaseTx.QSCs).Equals(approve))
}

func TestTxApproveDecrease_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()
	bank.GetMapper(ctx).SetAccount(qtypes.NewQOSAccountWithAddress(testFromAddr))

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	decreaseTx := TxDecreaseApprove{
		genTestApprove(),
	}
	require.NotNil(t, decreaseTx.ValidateData(ctx))

	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(createTx.Approve)

	require.Nil(t, decreaseTx.ValidateData(ctx))

	decreaseTx.QOS = btypes.NewInt(100)
	require.Nil(t, decreaseTx.ValidateData(ctx))

	decreaseTx.QOS = btypes.NewInt(110)
	require.NotNil(t, decreaseTx.ValidateData(ctx))
}

func TestTxApproveDecrease_Exec(t *testing.T) {
	ctx := defaultContext()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	decreaseTx := TxDecreaseApprove{
		genTestApprove(),
	}
	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(createTx.Approve)

	result, cross := decreaseTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	_, exists := approveMapper.GetApprove(createTx.From, createTx.To)
	require.False(t, exists)
	require.False(t, createTx.Approve.Minus(decreaseTx.QOS, decreaseTx.QSCs).IsPositive())
	require.True(t, createTx.Approve.Minus(decreaseTx.QOS, decreaseTx.QSCs).IsNotNegative())
}

func TestTxApproveUse_ValidateData(t *testing.T) {
	ctx := defaultContextWithQSC()

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.NotNil(t, useTx.ValidateData(ctx))

	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(createTx.Approve)
	require.NotNil(t, useTx.ValidateData(ctx))

	accountMapper := bank.GetMapper(ctx)
	accountMapper.SetAccount(genTestAccount(useTx.From))
	accountMapper.SetAccount(genTestAccount(useTx.To))

	require.Nil(t, useTx.ValidateData(ctx))

	useTx.QOS = btypes.NewInt(110)
	require.NotNil(t, useTx.ValidateData(ctx))

}

func TestTxApproveUse_GetSigner(t *testing.T) {
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.Equal(t, useTx.GetSigner(), []btypes.AccAddress{useTx.To})
}

func TestTxApproveUse_GetGasPayer(t *testing.T) {
	useTx := TxUseApprove{
		genTestApprove(),
	}
	require.Equal(t, useTx.GetGasPayer(), useTx.To)
}

func TestTxApproveUse_Exec(t *testing.T) {
	ctx := defaultContext()
	bank.GetMapper(ctx).SetAccount(qtypes.NewQOSAccountWithAddress(testFromAddr))

	createTx := TxCreateApprove{
		genTestApprove(),
	}
	useTx := TxUseApprove{
		genTestApprove(),
	}
	accountMapper := bank.GetMapper(ctx)
	accountMapper.SetAccount(genTestAccount(useTx.From))
	accountMapper.SetAccount(genTestAccount(useTx.To))

	approveMapper := mapper.GetMapper(ctx)
	approveMapper.SaveApprove(createTx.Approve)

	result, cross := useTx.Exec(ctx)
	require.Nil(t, cross)
	require.Equal(t, result.Code, btypes.CodeOK)

	approve, exists := approveMapper.GetApprove(useTx.From, useTx.To)
	require.False(t, exists)
	require.False(t, createTx.Minus(useTx.QOS, useTx.QSCs).Equals(approve))

}

func TestTxApproveCancel_ValidateData(t *testing.T) {
	ctx := defaultContext()
	createTx := TxCreateApprove{
		genTestApprove(),
	}
	cancelTx := TxCancelApprove{
		createTx.From,
		createTx.To,
	}
	require.NotNil(t, cancelTx.ValidateData(ctx))

	mapper := mapper.GetMapper(ctx)
	mapper.SaveApprove(createTx.Approve)

	require.Nil(t, cancelTx.ValidateData(ctx))
}

func TestTxApproveCancel_Exec(t *testing.T) {
	ctx := defaultContext()
	createTx := TxCreateApprove{
		genTestApprove(),
	}
	cancelTx := TxCancelApprove{
		createTx.From,
		createTx.To,
	}

	mapper := mapper.GetMapper(ctx)
	mapper.SaveApprove(createTx.Approve)

	result, _ := cancelTx.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

}

func TestTxApproveCancel_GetSigner(t *testing.T) {
	cancelTx := genApproveCancelTx()
	require.Equal(t, cancelTx.GetSigner(), []btypes.AccAddress{cancelTx.From})
}

func TestTxApproveCancel_GetGasPayer(t *testing.T) {
	cancelTx := genApproveCancelTx()
	require.Equal(t, cancelTx.GetGasPayer(), cancelTx.From)
}

func TestTxApproveCancel_CalcGas(t *testing.T) {
	cancelTx := genApproveCancelTx()
	require.Equal(t, cancelTx.CalcGas(), btypes.NewInt(0))
}

func TestTxApproveCancel_GetSignData(t *testing.T) {
	cancelTx := genApproveCancelTx()
	ret := Cdc.MustMarshalBinaryBare(cancelTx)
	require.Equal(t, cancelTx.GetSignData(), ret)
}
