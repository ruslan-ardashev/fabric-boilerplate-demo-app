// General utilities for chaincode
package data

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"build-chaincode/utils"
	"encoding/json"
	"github.com/pkg/errors"
)

var logger = shim.NewLogger("data")

func main() {
	logger.SetLevel(shim.LogInfo)
}

var indexStrings = map[string]string{
	"User":             "_users",
    "Thing":             "_things",
    "MTO":			"_mtos",
}

// Interface for saveable objects
type BlockchainItemer interface {
	GetId() string // TODO: just make everything have an Id field.
	SetId(string)
	GetIndexStr() string
	GetIdPrefix() string
}

type User struct {
	Id           string   `json:"id"` //Used to register with CA
	Salt         string   `json:"salt"`
	Hash         string   `json:"hash"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Things      []string `json:"things"`
	Address      string   `json:"address"`
	PhoneNumber  string   `json:"phoneNumber"`
	EmailAddress string   `json:"emailAddress"`
	Role         int64    `json:"role"`
}

// us
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


// funcs
func (u User) GetId() string         { return u.Id }
func (u User) SetId(id string)       { u.Id = id }
func (u User) GetIdPrefix() string { return "u" }
func (u User) GetIndexStr() string { return indexStrings[utils.GetTypeName(u)] }

type Thing struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

func (t Thing) GetId() string         { return t.Id }
func (t Thing) SetId(id string)       { t.Id = id }
func (t Thing) GetIdPrefix() string { return "u" }
func (t Thing) GetIndexStr() string { return indexStrings[utils.GetTypeName(t)] }

var ExampleStatus = map[string]bool{
	"CREATED":        	true,
	"UPDATE_REQUIRED":      true,
	"VALIDATION_REQUIRED":  true,
	"EXPIRED": 		true,
	"ACCEPTED":        	true,
	"DECLINED":	 	true,
}

var Roles = map[string]int64{
	"role1":  1,
	"role2":    2,
}

/*
	Public functions

*/

func GetIndexString(objectName string) string {
	idxStr := indexStrings[objectName]
	logger.Debugf("Returning indexstring for %v: %v", objectName, idxStr)
	return idxStr
}


// Save an object to the blockchain (and to the index). Generate an id if it doesn't exist yet.
func Save(stub shim.ChaincodeStubInterface, object BlockchainItemer) error {
	id := object.GetId()
	indexString := object.GetIndexStr()
	idPrefix := object.GetIdPrefix()

	if indexString == "" {
		return errors.New("Indexstring not found")
	}

	if id == "" {
		idBytes, err := utils.CreateId(stub, indexString, idPrefix)
		if err != nil {
			return errors.New("Could not create id")
		}
		id = string(idBytes)
		object.SetId(id)
	}
	return utils.Put(stub, object, indexString, id)
}

// Reset all index strings.
func ResetIndexes(stub shim.ChaincodeStubInterface) error {
	indexes := indexStrings
	logger.Infof("indexes: %v", indexes)
	for _, v := range indexes {
		// Marshal the index
		emptyIndex := make(map[string]bool)

		empty, err := json.Marshal(emptyIndex)
		if err != nil {
			return errors.New("Error marshalling")
		}
		err = stub.PutState(v, empty)

		if err != nil {
			return errors.New("Error deleting index")
		}
		logger.Debugf("Delete with success from ledger: " + v)
	}

	// save MTOs initially here
	var mtosKeyString = GetIndexString("MTO")
	var mtos []MTO
	var err error

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
	err = stub.PutState(mtosKeyString, jsonMTOsAsBytes)
	if err != nil {
		return err
	}

	return nil
}
