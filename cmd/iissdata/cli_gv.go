package main

import (
	"fmt"
	"github.com/icon-project/rewardcalculator/common/db"
	"github.com/icon-project/rewardcalculator/core"
)

func (cli *CLI) governanceVariable(blockHeight uint64, incentive uint64, reward uint64,
	mainPRepCount uint64, subPRepCount uint64) {
	fmt.Printf("Start set header of IISS data DB.\n")

	bucket, _ := cli.DB.GetBucket(db.PrefixIISSGV)

	gv := new(core.IISSGovernanceVariable)
	gv.BlockHeight = blockHeight
	gv.IncentiveRep = incentive
	gv.RewardRep = reward
	gv.MainPRepCount = mainPRepCount
	gv.SubPRepCount = subPRepCount

	value, _ := gv.Bytes()
	bucket.Set(gv.ID(), value)

	fmt.Printf("Add governance variable: ID: %+v, %s\n", gv.ID(), gv.String())
}
