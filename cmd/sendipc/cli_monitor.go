package main

import (
	"fmt"
	"github.com/icon-project/rewardcalculator/common"

	"github.com/icon-project/rewardcalculator/common/ipc"
	"github.com/icon-project/rewardcalculator/rewardcalculator"
)



func (cli *CLI) query(conn ipc.Connection, address string) *rewardcalculator.ResponseQuery {
	var addr common.Address
	resp := new(rewardcalculator.ResponseQuery)

	addr.SetString(address)

	conn.SendAndReceive(msgQuery, cli.id, &addr, resp)
	fmt.Printf("QUERY command get response: %s\n", Display(resp))

	return resp
}