/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Bike struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryBike" {
		return s.queryBike(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createBike" {
		return s.createBike(APIstub, args)
	} else if function == "queryAllBikes" {
		return s.queryAllBikes(APIstub)
	} else if function == "changeBikeOwner" {
		return s.changeBikeOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryBike(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	bikeAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(bikeAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	bikes := []Bike{
		Bike{Make: "Honda", Model: "Shine", Colour: "blue", Owner: "Gowda"},
		Bike{Make: "BMW", Model: "F 700", Colour: "black", Owner: "George"},
		Bike{Make: "RoyalEnfield", Model: "Bullet 350", Colour: "black", Owner: "Bhaskar"},
		Bike{Make: "KTM", Model: "RC 200", Colour: "blue", Owner: "Darshan"},
		Bike{Make: "TVS", Model: "Apache", Colour: "blue", Owner: "Krishna"},
		Bike{Make: "Honda", Model: "205", Colour: "purple", Owner: "Raman"},
		Bike{Make: "Bajaj", Model: "Pulsar", Colour: "red", Owner: "Pradeep"},
		Bike{Make: "Yamaha", Model: "XSR 155", Colour: "violet", Owner: "Naveen"},
		Bike{Make: "Kawasaki", Model: "Ninja H2", Colour: "blue", Owner: "Raghav"},
		Bike{Make: "Hardly Davidson", Model: "Iron 883", Colour: "black", Owner: "Dinesh"},
	}

	i := 0
	for i < len(bikes) {
		fmt.Println("i is ", i)
		bikeAsBytes, _ := json.Marshal(bikes[i])
		APIstub.PutState("BIKE"+strconv.Itoa(i), bikeAsBytes)
		fmt.Println("Added", bikes[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createBike(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var bike = Bike{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

	bikeAsBytes, _ := json.Marshal(bike)
	APIstub.PutState(args[0], bikeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllBikes(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "BIKE0"
	endKey := "BIKE999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllBikes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeBikeOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	bikeAsBytes, _ := APIstub.GetState(args[0])
	bike := Bike {}

	json.Unmarshal(bikeAsBytes, &bike)
	bike.Owner = args[1]

	bikeAsBytes, _ = json.Marshal(bike)
	APIstub.PutState(args[0], bikeAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
