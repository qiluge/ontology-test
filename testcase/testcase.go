/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package testcase

import (
	"github.com/ontio/ontology-test/testcase/http"
	"github.com/ontio/ontology-test/testcase/smartcontract"
	"github.com/ontio/ontology-test/testcase/vm"
	"github.com/ontio/ontology-test/testframework"
	"math"
	"time"
)

//TestCase list
func init() {
	testframework.TFramework.SetBeforeCallback(BeforeTestCase)
	http.TestHttp()
	vm.TestVM()
	smartcontract.TestSmartContract()
}

func BeforeTestCase(ctx *testframework.TestFrameworkContext) {
	defAccount, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("GetDefaultAccount error:%s", err)
		ctx.FailNow()
		return
	}
	newAccount := ctx.NewAccount()
	amount := uint64(10000)
	balance, err := ctx.Ont.Rpc.GetBalance(defAccount.Address)
	if err != nil {
		ctx.LogError("GetBalance error:%s", err)
		ctx.FailNow()
		return
	}
	minONG := uint64(100000 * math.Pow10(9))
	if balance.Ong > minONG {
		ctx.LogInfo("Default account balance ont:%d ong:%d", balance.Ont, balance.Ong)
		return
	}
	if balance.Ont == 0 {
		ctx.LogWarn("Default Account balance = 0 ")
		return
	}
	if balance.Ont < amount {
		amount = balance.Ont
	}
	_, err = ctx.Ont.Rpc.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ont", defAccount, newAccount.Address, amount)
	if err != nil {
		ctx.LogError("Transfer error:%s", err)
		ctx.FailNow()
		return
	}
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return
	}
	unboundONG, err := ctx.Ont.Rpc.UnboundONG(defAccount.Address)
	if err != nil {
		ctx.LogError("UnboundONG error:%s", err)
		ctx.FailNow()
		return
	}
	_, err = ctx.Ont.Rpc.WithdrawONG(ctx.GetGasPrice(), ctx.GetGasLimit(), defAccount, unboundONG)
	if err != nil {
		ctx.LogError("WithdrawONG error:%s", err)
		ctx.FailNow()
		return
	}
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		ctx.FailNow()
		return
	}
	balance, err = ctx.Ont.Rpc.GetBalance(defAccount.Address)
	if err != nil {
		ctx.LogInfo("GetBalance error:%s", err)
		ctx.FailNow()
		return
	}
	ctx.LogInfo("Default account balance ont:%d ong:%d", balance.Ont, balance.Ong)
}
