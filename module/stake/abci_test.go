package stake

import (
	"encoding/binary"
	"testing"
	"time"

	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestValidatorMapper(t *testing.T) {

	ctx := defaultContext()
	validatorMapper := mapper.GetMapper(ctx)

	validator := types.Validator{
		Description:     types.Description{Moniker: "test"},
		OperatorAddress: btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address()),
		Owner:           btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		ConsPubKey:      ed25519.GenPrivKey().PubKey(),
		BondTokens:      btypes.NewInt(500),
		Status:          types.Active,
		MinPeriod:       0,
		BondHeight:      1,
	}

	valAddr := validator.GetValidatorAddress()
	key := types.BuildValidatorKey(valAddr)
	validatorMapper.Set(key, validator)

	v, exists := validatorMapper.GetValidator(valAddr)
	require.Equal(t, true, exists)
	require.Equal(t, btypes.NewInt(500), v.GetBondTokens())
	require.Equal(t, true, v.IsActive())

	now := time.Now().UTC().Unix()
	for i := int64(0); i <= int64(100); i++ {
		addr := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())
		validatorMapper.Set(types.BuildInactiveValidatorKey(now+i, addr), i)
	}

	iter := validatorMapper.IteratorInactiveValidator(0, now+int64(20))
	defer iter.Close()

	i := 0
	for ; iter.Valid(); iter.Next() {
		i++
		k := iter.Key()
		cp := int64(binary.BigEndian.Uint64(k[1:9]))
		require.Equal(t, true, cp >= now)
		now = cp
	}
	require.Equal(t, 20, i)

	for i := int64(100); i <= int64(200); i++ {
		addr := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())
		validatorMapper.Set(types.BuildValidatorByVotePower(i, addr), 1)
	}

	descIter := validatorMapper.IteratorValidatorByVoterPower(false)
	defer descIter.Close()

	power := uint64(200)
	for ; descIter.Valid(); descIter.Next() {
		p := descIter.Key()
		cp := binary.BigEndian.Uint64(p[1:9])

		require.Equal(t, true, power >= cp)
		power = cp
	}

	ascIter := validatorMapper.IteratorValidatorByVoterPower(true)
	defer ascIter.Close()

	power = uint64(0)
	for ; ascIter.Valid(); ascIter.Next() {
		p := ascIter.Key()
		cp := binary.BigEndian.Uint64(p[1:9])

		require.Equal(t, true, power <= cp)
		power = cp
	}

}

func TestVoteInfoMapper(t *testing.T) {

	ctx := defaultContext()

	sm := mapper.GetMapper(ctx)

	addr := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	voteInfo := types.NewValidatorVoteInfo(1, 1, 1)

	sm.SetValidatorVoteInfo(addr, voteInfo)

	v, exists := sm.GetValidatorVoteInfo(addr)
	require.Equal(t, int64(1), v.StartHeight)

	sm.DelValidatorVoteInfo(addr)

	v, exists = sm.GetValidatorVoteInfo(addr)
	require.Equal(t, false, exists)

	for i := int64(0); i <= 10; i++ {
		sm.SetVoteInfoInWindow(addr, i, false)
	}

	vote := sm.GetVoteInfoInWindow(addr, 11)
	require.Equal(t, true, vote)

	vote = sm.GetVoteInfoInWindow(addr, 10)
	require.Equal(t, false, vote)

	sm.ClearValidatorVoteInfoInWindow(addr)

	vote = sm.GetVoteInfoInWindow(addr, 10)
	require.Equal(t, true, vote)

}

func defaultContext() context.Context {

	mapperMap := make(map[string]bmapper.IMapper)

	paramsMapper := params.NewMapper()
	mapperMap[params.MapperName] = paramsMapper

	accountMapper := account.NewAccountMapper(txs.Cdc, qtypes.ProtoQOSAccount)
	mapperMap[account.AccountMapperName] = accountMapper

	sm := mapper.NewMapper()
	sm.SetCodec(txs.Cdc)
	mapperMap[types.MapperName] = sm

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	for _, v := range mapperMap {
		cms.MountStoreWithDB(v.GetStoreKey(), btypes.StoreTypeIAVL, db)
	}
	cms.LoadLatestVersion()

	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}
