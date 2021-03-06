package datatype

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

func TestInteger(ctx *testframework.TestFrameworkContext) bool {
	code := "00C56B5A616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestInteger GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestInteger",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestInteger DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestInteger WaitForGenerateBlock error:%s", err)
		return false
	}
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		codeAddress,
		[]interface{}{},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestInteger InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, 10)
	if err != nil {
		ctx.LogError("TestInteger test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main()
    {
        return 10;
    }
}
*/
