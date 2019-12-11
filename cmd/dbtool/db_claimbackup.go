package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/icon-project/rewardcalculator/common"
	"github.com/icon-project/rewardcalculator/common/db"
	"github.com/icon-project/rewardcalculator/core"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func queryClaimBackupDB(input Input) (err error) {
	if input.path == "" {
		fmt.Println("Enter dbPath")
		return errors.New("invalid db path")
	}

	if input.height == 0 {
		err = printDB(input.path, util.BytesPrefix([]byte(db.PrefixClaim)), printClaimBackup)
	} else {
		prefix := core.MakeIteratorPrefix(db.PrefixClaim, input.height, nil, 0)
		err = printDB(input.path, prefix, printClaimBackup)
	}
	return
}

func isManageKey(key []byte) bool {
	return len(key) == len(db.PrefixManagement) && bytes.Equal(key, []byte(db.PrefixManagement))
}

func printClaimBackupInfo(value []byte) error {
	var cbInfo core.ClaimBackupInfo
	err := cbInfo.SetBytes(value)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", cbInfo.String())
	return nil
}

func printClaimBackup(key []byte, value []byte) (err error) {
	if isManageKey(key) {
		return printClaimBackupInfo(value)
	}

	if claim, e := newClaimFromBackup(key, value); e != nil {
		return e
	} else {
		fmt.Printf("Key(%s), Value(%s)\n", core.ClaimBackupKeyString(key), claim.String())
		return nil
	}
}

func newClaimFromBackup(key []byte, value []byte) (*core.Claim, error) {
	if isManageKey(key) {
		return nil, nil
	}
	if claim, err := core.NewClaimFromBytes(value); err != nil {
		fmt.Printf("Failed to make claim instance %v\n", err)
		return nil, err
	} else {
		claim.Address = *common.NewAddress(key[core.BlockHeightSize:])
		return claim, nil
	}
}
