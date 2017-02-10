package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"os"
	"build-chaincode/query"
	"build-chaincode/invoke"
	"github.com/pkg/errors"
	"encoding/json"
	"strconv"
)

type SimpleChaincode struct {}

type Account struct {
	FirstName       	string `json:"firstName"`
	LastName		 	string `json:"lastName"`
	AccountNumber   	string `json:"accountNumber"`
	Balance			string `json:"balance"`
}

type MTO struct {
	Name            string `json:"name"`
	Accounts		 []Account `json:"accounts"`
}

var mtoIndexString = "_mtos"

var logger = shim.NewLogger("fabric-boilerplate")

// Init - is called when the chaincode is deployed
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Infof("~~ Running the Init we want. ~~")
	bytes, err := invoke.Init(stub, function, args)
	if err != nil { logger.Errorf("Error invoking Init: %v \n %v", function, err) }

	// save MTOs initially here
	var mtos []MTO

	// create sample MTOs and Banks here
	var mto0 = MTO{
		Name: "WesternUnion", 
		Accounts: []Account{
			{
				FirstName: "Bob",
				LastName: "Jenkins",
				AccountNumber: "123",
				Balance: "10000",
			},
			{
				FirstName: "Bob",
				LastName: "Leeroy",
				AccountNumber: "1234",
				Balance: "10000",
			},
		},
	}

	var mto1 = MTO{
		Name: "CambodianExchange", 
		Accounts: []Account{
			{
				FirstName: "Alice",
				LastName: "Jenkins",
				AccountNumber: "12345",
				Balance: "10000",
			},
			{
				FirstName: "Alice",
				LastName: "Leeroy",
				AccountNumber: "123456",
				Balance: "10000",
			},
		},
	}

	var bank0 = MTO{
		Name: "JPMorganChaseInternational", 
		Accounts: []Account{
			{
				FirstName: "John",
				LastName: "Jenkins",
				AccountNumber: "1234567",
				Balance: "10000",
			},
			{
				FirstName: "John",
				LastName: "Leeroy",
				AccountNumber: "12345678",
				Balance: "10000",
			},
		},
	}

	mtos = []MTO{
		mto0, mto1, bank0,
	}

	logger.Infof("created first mtos: %v", mtos)

	jsonMTOsAsBytes, _ := json.Marshal(mtos)

	// store array into "_mtos"
	err = stub.PutState(mtoIndexString, jsonMTOsAsBytes)
	if err != nil {
		return nil, err
	}
	
	return bytes, err

}

// Invoke - handles all the invoke functions
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "createAccount" {
		logger.Infof("~~ Running the createAccount we want. ~~")
		return t.createAccount(stub, args)
	} else if function == "transfer" {
		logger.Infof("~~ Running the transfer we want. ~~")
		return t.transfer(stub, args)
	} else {
		bytes, err := invoke.Invoke(stub, function, args)
		if err != nil { logger.Errorf("Error invoking %v: %v", function, err) }
		return bytes, err
	}
	
}

//==============================================================================================================================
//		Invoke Functions
//==============================================================================================================================
// createAccount
//	arg0 - string mto.Name
//	arg1 - string account.firstName
//	arg2 - string account.lastName
// 	arg3 - string account.number
//	arg4 - string account.balance
// func createAccount(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, args []string) ([]byte, error) {
func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Infof("~~ createAccount begin ~~")

	// // ecert, err := stub.GetState(name)
	// // if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	var mtos []MTO
	var err error

	mtosAsBytes, err := stub.GetState(mtoIndexString)

	if err != nil {
		return nil, errors.New("{\"Error\":\"Failed to get mtos for mtoIndexString.\"}")
	}

	json.Unmarshal(mtosAsBytes, &mtos)

	// // find mto with given name
	index := -1
	for i := 0; i < len(mtos); i++ {
		if mtos[i].Name == args[0] {
			index = i
		}
	}

	if index == -1 {
		return nil, errors.New("{\"Error\":\"Failed to get assets for provided user - user not found in users.\"}")
	}

	// // create account
	var account = Account{
		FirstName: args[1],
		LastName: args[2],
		AccountNumber: args[3],
		Balance: args[4],
	}

	// // add account to mto.accounts
	mtos[index].Accounts = append(mtos[index].Accounts, account)

	// // store mtos back in KVS
	mtosAsBytes, _ = json.Marshal(mtos)

	err = stub.PutState(mtoIndexString, mtosAsBytes)
	if err != nil {
		return nil, err
	}

	logger.Infof("~~ createAccount returning ~~")

	return nil, nil

}

// transfer
// transfers balance from 
// arg0 - (Source MTO).Name
// arg1 - (Source account).Number
// arg2 - balance to transfer from source account
// arg3 - (Destination MTO).Name
// arg4 - (Destination account).Number
// if successful, add functionality to return new balances as strings [source balance, destination balance]
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Infof("~~ Invoking transfer ~~")

	var mtos []MTO
	var err error

	mtosAsBytes, err := stub.GetState(mtoIndexString)

	if err != nil {
		logger.Infof("~~ Failed to find mtoIndexString mtos. ~~")
		return nil, errors.New("{\"Error\":\"Failed to get mtos for mtoIndexString.\"}")
	}

	json.Unmarshal(mtosAsBytes, &mtos)

	// ecert, err := stub.GetState(name)
	// if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	var sourceMTOName = args[0]
	var sourceAccountNumber = args[1]
	var balanceToTransfer = args[2]
	var destinationMTOName = args[3]
	var destinationAccountNumber = args[4]

	// // find MTOs to transfer between
	sourceIndex := -1
	destinationIndex := -1
	for i := 0; i < len(mtos); i++ {

		if mtos[i].Name == sourceMTOName {
			logger.Infof("~~ found sourceIndex at %v ~~", i)
			sourceIndex = i
		}

		// support both ifs, not a bug
		// transfers between accounts at self MTO
		if mtos[i].Name == destinationMTOName {
			logger.Infof("~~ found destinationIndex at %v ~~", i)
			destinationIndex = i
		} 

	}

	if (sourceIndex == -1) || (destinationIndex == -1) {
		logger.Infof("~~ Failed to find both source and destination MTO. ~~")
		return nil, errors.New("{\"Error\":\"Failed to get both MTOs to transfer between.\"}")
	}

	// get account with number from source
	var sourceMTOAccounts = mtos[sourceIndex].Accounts
	var sourceAccountIndex = -1

	for i := 0; i < len(sourceMTOAccounts); i++ {
		if sourceMTOAccounts[i].AccountNumber == sourceAccountNumber {
			logger.Infof("~~ found sourceAccountIndex at %v ~~", i)
			sourceAccountIndex = i
		}
	}

	if sourceAccountIndex == -1 {
		logger.Infof("~~ Failed to get source MTO account given specified account number. ~~")
		return nil, errors.New("{\"Error\":\"Failed to get source MTO account given specified account number.\"}")
	}

	var sourceAccount = sourceMTOAccounts[sourceAccountIndex]

	logger.Infof("~~ got source account. index %v ~~", sourceAccountIndex)

	// get account with number from destination
	var destinationMTOAccounts = mtos[destinationIndex].Accounts
	var destinationAccountIndex = -1

	for i := 0; i < len(destinationMTOAccounts); i++ {
		if destinationMTOAccounts[i].AccountNumber == destinationAccountNumber {
			destinationAccountIndex = i
		}
	}

	if destinationAccountIndex == -1 {
		logger.Infof("~~ Failed to get destination MTO account given specified account number. ~~")
		return nil, errors.New("{\"Error\":\"Failed to get destination MTO account given specified account number.\"}")
	}

	var destinationAccount = destinationMTOAccounts[destinationAccountIndex]

	logger.Infof("~~ got destination account. index %v ~~", destinationAccountIndex)

	// get balances from source, destination, args
	sourceBalance, err0 := strconv.ParseFloat(sourceAccount.Balance, 64)
	destinationBalance, err1 := strconv.ParseFloat(destinationAccount.Balance, 64)
	balanceToTransferFloat, err2 := strconv.ParseFloat(balanceToTransfer, 64)

	logger.Infof("~~ sourceBalance: %v ~~", sourceBalance)
	logger.Infof("~~ destinationBalance: %v ~~", destinationBalance)
	logger.Infof("~~ balanceToTransferFloat: %v ~~", balanceToTransferFloat)

	// check for errors in conversion
	if err0 != nil {
		logger.Infof("~~ failed parsing sourceBalance ~~")
		return nil, err0
	} else if err1 != nil {
		logger.Infof("~~ failed parsing destinationBalance ~~")
		return nil, err1
	} else if err2 != nil {
		logger.Infof("~~ failed parsing balanceToTransferFloat ~~")
		return nil, err2
	}

	// check for sufficient balance in sourceBalance
	var resultingSourceBalance = sourceBalance - balanceToTransferFloat

	logger.Infof("~~ resultingSourceBalance: %v ~~", resultingSourceBalance)

	if resultingSourceBalance < 0 {
		logger.Infof("~~ Insufficient source funds to initiate transfer. ~~")
		return nil, errors.New("{\"Error\":\"Insufficient source funds to initiate transfer.\"}")
	}

	// transfer
	sourceBalance = sourceBalance - balanceToTransferFloat
	destinationBalance = destinationBalance + balanceToTransferFloat

	logger.Infof("~~ new sourceBalance: %v ~~", sourceBalance)
	logger.Infof("~~ new destinationBalance: %v ~~", destinationBalance)

	sourceBalanceString := strconv.FormatFloat(sourceBalance, 'f', -1, 64)
	destinationBalanceString := strconv.FormatFloat(destinationBalance, 'f', -1, 64)

	logger.Infof("~~ sourceBalanceString: %v ~~", sourceBalanceString)
	logger.Infof("~~ destinationBalanceString: %v ~~", destinationBalanceString)

	// store balances to accounts, accounts to mtos
	mtos[sourceIndex].Accounts[sourceAccountIndex].Balance = sourceBalanceString
	mtos[destinationIndex].Accounts[destinationAccountIndex].Balance = destinationBalanceString

	logger.Infof("~~ new mtos: %v ~~", mtos)

	// store mtos back in KVS
	mtosAsBytes, _ = json.Marshal(mtos)

	err = stub.PutState(mtoIndexString, mtosAsBytes)

	if err != nil {
		logger.Infof("~~ error putting state back ~~")
		return nil, err
	}

	logger.Infof("~~ put state back succes!")

	// only reach if success
	returnBalances := []string{sourceBalanceString, destinationBalanceString}

	logger.Infof("~~ Completing invoking transfer, returnBalances: %v", returnBalances)

	return json.Marshal(returnBalances)

}

// Query - handles all the query functions
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "mtos" {
		return t.mtos(stub, args)
	} else if function == "accounts" {
		return t.accounts(stub, args)
	} else if function == "afterTransfer" {
		return t.afterTransfer(stub, args)
	} else {
		bytes, err := query.Query(stub, function, args)
		if err != nil { logger.Errorf("Error querying %v: %v", function, err) }
		return bytes, err
	}
}

func (t *SimpleChaincode) mtos(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var mtos []MTO
	var err error
	returnMtos := []string{}

	mtosAsBytes, err := stub.GetState(mtoIndexString)

	if err != nil {
		return nil, errors.New("{\"Error\":\"Failed to get mtos for mtosKeyString.\"}")
	}

	json.Unmarshal(mtosAsBytes, &mtos)

	// // find mto with given name
	for i := 0; i < len(mtos); i++ {
		returnMtos = append(returnMtos, mtos[i].Name)
	}

	return json.Marshal(returnMtos)

}

// arg[0] - mtoName
func (t *SimpleChaincode) accounts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var mtos []MTO
	var err error

	var sourceMTOName = args[0]

	mtosAsBytes, err := stub.GetState(mtoIndexString)

	if err != nil {
		return nil, errors.New("{\"Error\":\"Failed to get mtos for mtosKeyString.\"}")
	}

	json.Unmarshal(mtosAsBytes, &mtos)

	index := -1
	for i := 0; i < len(mtos); i++ {

		if mtos[i].Name == sourceMTOName {
			logger.Infof("~~ found index at %v ~~", i)
			index = i
		}

	}

	if index == -1 {
		logger.Infof("~~ Did not find MTO to create account for in mtos. ~~")
		return nil, errors.New("{\"Error\":\"Did not find MTO to create account for in mtos.\"}")
	}

	return json.Marshal(mtos[index])

}

// remove if can detect errors from invokes
func (t *SimpleChaincode) afterTransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	return stub.GetState(mtoIndexString)

}

// Main - starts up the chaincode
func main() {
	logger.SetLevel(shim.LogInfo)

	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)

	err := shim.Start(new(SimpleChaincode))
	if err != nil { logger.Errorf("Error starting chaincode:", err) }
}
