package mapper

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
)

// 统计投票信息，返回提议结果，投票结果，投票验证节点集合以及质押处理类型
func Tally(ctx context.Context, mapper *Mapper, proposal gtypes.Proposal) (passes gtypes.ProposalResult, tallyResults gtypes.TallyResult, validators map[string]bool, deductOption gtypes.DeductOption) {
	results := make(map[gtypes.VoteOption]btypes.BigInt)
	results[gtypes.OptionYes] = btypes.ZeroInt()
	results[gtypes.OptionAbstain] = btypes.ZeroInt()
	results[gtypes.OptionNo] = btypes.ZeroInt()
	results[gtypes.OptionNoWithVeto] = btypes.ZeroInt()
	// total power for voting this proposal
	totalVotingPower := btypes.ZeroInt()
	// total voting power int the whole network
	totalSystemPower := btypes.ZeroInt()
	// voted validators
	validators = make(map[string]bool)
	// current validators int the whole network
	currValidators := make(map[string]stake.Validator)

	// 统计当前验证节点信息
	sm := stake.GetMapper(ctx)
	iterator := sm.IteratorValidatorByVoterPower(false)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.ValAddress(key[9:])
		if validator, exists := sm.GetValidator(valAddr); exists {
			currValidators[validator.GetValidatorAddress().String()] = validator
			totalSystemPower = totalSystemPower.Add(validator.BondTokens)
		}
	}

	// 统计投票人对应验证节点信息，一个验证节点任意delegator投过票，则算这个验证节点参与了投票
	votesIterator := mapper.GetVotes(proposal.ProposalID)
	defer votesIterator.Close()
	for ; votesIterator.Valid(); votesIterator.Next() {
		vote := &gtypes.Vote{}
		mapper.DecodeObject(votesIterator.Value(), vote)

		sm.IterateDelegationsInfo(vote.Voter, func(delegation stake.Delegation) {
			if _, ok := currValidators[delegation.ValidatorAddr.String()]; ok {
				totalVotingPower = totalVotingPower.Add(delegation.Amount)
				results[vote.Option] = results[vote.Option].Add(delegation.Amount)
				validators[delegation.ValidatorAddr.String()] = true
			}
		})
	}

	params := mapper.GetLevelParams(ctx, proposal.GetProposalLevel())
	tallyResults = gtypes.TallyResult{
		Yes:        results[gtypes.OptionYes],
		Abstain:    results[gtypes.OptionAbstain],
		No:         results[gtypes.OptionNo],
		NoWithVeto: results[gtypes.OptionNoWithVeto],
	}

	// If no one votes, proposal fails
	if totalVotingPower.Equal(results[gtypes.OptionAbstain]) {
		return gtypes.REJECT, tallyResults, validators, gtypes.DepositDeductNone
	}

	//if more than `Quorum` of voters abstain, proposal fails
	if types.NewDecFromInt(totalVotingPower.Div(totalSystemPower)).LT(params.Quorum) {
		return gtypes.REJECT, tallyResults, validators, gtypes.DepositDeductPart
	}

	// If more than `Veto` of voters veto, proposal fails
	if types.NewDecFromInt(results[gtypes.OptionNoWithVeto]).Quo(types.NewDecFromInt(totalVotingPower)).GT(params.Veto) {
		return gtypes.REJECTVETO, tallyResults, validators, gtypes.DepositDeductAll
	}

	// If more than `Threshold` of non-abstaining voters vote Yes, proposal passes
	if types.NewDecFromInt(results[gtypes.OptionYes]).Quo(types.NewDecFromInt(totalVotingPower.Sub(results[gtypes.OptionAbstain]))).GT(params.Threshold) {
		return gtypes.PASS, tallyResults, validators, gtypes.DepositDeductNone
	}

	// If more than `Threshold` of non-abstaining voters vote No, proposal fails

	return gtypes.REJECT, tallyResults, validators, gtypes.DepositDeductNone
}
