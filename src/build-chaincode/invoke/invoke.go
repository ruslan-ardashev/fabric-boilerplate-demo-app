package invoke
import (
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"build-chaincode/data"
	"strconv"
)

var logger = shim.NewLogger("invoke")

func main() {
	logger.SetLevel(shim.LogDebug)
}

var Functions = map[string]func(shim.ChaincodeStubInterface,[]string)([]byte, error) {
    "add_user": add_user,
    "add_thing": add_thing,
    "add_test_data": add_test_data,
    "createAccount": createAccount,
    "transfer": transfer,
}

// Invoke function.
func Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Infof("-- Invoking function %v with args %v", function, args)

	if function == "init" {
		return Init(stub, "init", args)
    } else {
        return Functions[function](stub,args)
    }

	return nil, errors.New("Received unknown invoke function name")
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================

func Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Infof("Deployed chaincode.")

	return nil, data.ResetIndexes(stub)
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
func createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// // ecert, err := stub.GetState(name)
	// // if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	var mtosKeyString = data.GetIndexString("MTO")
	var mtos []data.MTO
	var err error

	mtosAsBytes, err := stub.GetState(mtosKeyString)

	if err != nil {
		return nil, errors.New("{\"Error\":\"Failed to get mtos for mtosKeyString.\"}")
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
	var account = data.Account{
		FirstName: args[1],
		LastName: args[2],
		AccountNumber: args[3],
		Balance: args[4],
	}

	// // add account to mto.accounts
	mtos[index].Accounts = append(mtos[index].Accounts, account)

	// // store mtos back in KVS
	mtosAsBytes, _ = json.Marshal(mtos)

	err = stub.PutState(mtosKeyString, mtosAsBytes)
	if err != nil {
		return nil, err
	}

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
func transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// ecert, err := stub.GetState(name)
	// if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	var mtosKeyString = data.GetIndexString("MTO")
	var mtos []data.MTO
	var err error

	var sourceMTOName = args[0]
	var sourceAccountNumber = args[1]
	var balanceToTransfer = args[2]
	var destinationMTOName = args[3]
	var destinationAccountNumber = args[4]

	mtosAsBytes, err := stub.GetState(mtosKeyString)

	if err != nil {
		return nil, errors.New("{\"Error\":\"Failed to get mtos for mtosKeyString.\"}")
	}

	json.Unmarshal(mtosAsBytes, &mtos)

	// // find MTOs to transfer between
	sourceIndex := -1
	destinationIndex := -1
	for i := 0; i < len(mtos); i++ {
		if mtos[i].Name == sourceMTOName {
			sourceIndex = i
		} else if mtos[i].Name == destinationMTOName {
			destinationIndex = i
		} 
	}

	if (sourceIndex == -1) || (destinationIndex == -1) {
		return nil, errors.New("{\"Error\":\"Failed to get both MTOs to transfer between.\"}")
	}

	// get account with number from source
	var sourceMTOAccounts = mtos[sourceIndex].Accounts
	var sourceAccountIndex = -1

	for i := 0; i < len(sourceMTOAccounts); i++ {
		if sourceMTOAccounts[i].AccountNumber == sourceAccountNumber {
			sourceAccountIndex = i
		}
	}

	if sourceAccountIndex == -1 {
		return nil, errors.New("{\"Error\":\"Failed to get source MTO account given specified account number.\"}")
	}

	var sourceAccount = sourceMTOAccounts[sourceAccountIndex]

	// get account with number from destination
	var destinationMTOAccounts = mtos[destinationIndex].Accounts
	var destinationAccountIndex = -1

	for i := 0; i < len(destinationMTOAccounts); i++ {
		if destinationMTOAccounts[i].AccountNumber == destinationAccountNumber {
			destinationAccountIndex = i
		}
	}

	if destinationAccountIndex == -1 {
		return nil, errors.New("{\"Error\":\"Failed to get destination MTO account given specified account number.\"}")
	}

	var destinationAccount = destinationMTOAccounts[destinationAccountIndex]

	// get balances from source, destination, args
	sourceBalance, err0 := strconv.ParseFloat(sourceAccount.Balance, 64)
	destinationBalance, err1 := strconv.ParseFloat(destinationAccount.Balance, 64)
	balanceToTransferFloat, err2 := strconv.ParseFloat(balanceToTransfer, 64)

	// check for errors in conversion
	if err0 != nil {
		return nil, err0
	} else if err1 != nil {
		return nil, err1
	} else if err2 != nil {
		return nil, err2
	}

	// check for sufficient balance in sourceBalance
	var resultingSourceBalance = sourceBalance - balanceToTransferFloat

	if resultingSourceBalance < 0 {
		return nil, errors.New("{\"Error\":\"Insufficient source funds to initiate transfer.\"}")
	}

	// transfer
	sourceBalance = sourceBalance - balanceToTransferFloat
	destinationBalance = destinationBalance + balanceToTransferFloat

	sourceBalanceString := strconv.FormatFloat(sourceBalance, 'f', -1, 64)
	destinationBalanceString := strconv.FormatFloat(destinationBalance, 'f', -1, 64)

	// store balances to accounts, accounts to mtos
	mtos[sourceIndex].Accounts[sourceAccountIndex].Balance = sourceBalanceString
	mtos[destinationIndex].Accounts[destinationAccountIndex].Balance = destinationBalanceString

	// store mtos back in KVS
	mtosAsBytes, _ = json.Marshal(mtos)

	err = stub.PutState(mtosKeyString, mtosAsBytes)
	if err != nil {
		return nil, err
	}

	// only reach if success
	returnBalances := []string{sourceBalanceString, destinationBalanceString}

	return json.Marshal(returnBalances)

}


func add_test_data(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var usersIndex = args[0]
    var thingsIndex = args[1]

    var users []data.User
    err := json.Unmarshal([]byte(usersIndex), &users)
	if err != nil { return nil, err }
    for _,user := range users {
        data.Save(stub, user)
    }

    var things []data.Thing
    err = json.Unmarshal([]byte(thingsIndex), &things)
	if err != nil { return nil, err }
    for _,thing := range things {
        data.Save(stub, thing)
    }
    return nil,err
}

// args 0 is the caller id (not anymore needed in fabric v. 0.6)
func add_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var item data.User
	err := json.Unmarshal([]byte(args[1]), &item)
	if err != nil { return nil, err }
	return nil, data.Save(stub, item)
}
func add_thing(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var item data.Thing
	err := json.Unmarshal([]byte(args[1]), &item)
	if err != nil { return nil, err }
	return nil, data.Save(stub, item)
}
