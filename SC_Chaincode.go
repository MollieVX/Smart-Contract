/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//structure for vehicles
type Vehicle struct {
	ID     string `json:"id"`
	Balance  int    `json:"balance"`
}

//structure for services
type Service struct {
	Type     string `json:"type"`
	Cost  int    `json:"cost"`
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("No arguments required")
	}

	//Creating 3 vehicles
	var vehicle1, vehicle2, vehicle3 Vehicle
	vehicle1.ID = "1"
	vehicle1.Balance = 1000
	vehicle2.ID = "2"
	vehicle2.Balance = 1000
	vehicle3.ID = "3"
	vehicle3.Balance = 1000

	b, err := json.Marshal(vehicle1)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for vehicle 1")
	}

	err = stub.PutState("1", b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(vehicle2)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for vehicle 2")
	}

	err = stub.PutState("2", b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(vehicle3)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for vehicle 3")
	}

	err = stub.PutState("3", b)
	if err != nil {
		return nil, err
	}

	//Creating 3 services
	var service1, service2, service3 Service
	service1.Type = "Wash"
	service1.Cost = 15
	service2.Type = "Parking"
	service2.Cost = 20
	service3.Type = "Toll"
	service3.Cost = 10

	b, err = json.Marshal(service1)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for service 1")
	}

	err = stub.PutState("Wash", b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(service2)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for service 2")
	}

	err = stub.PutState("Parking", b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(service3)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for service 3")
	}

	err = stub.PutState("Toll", b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "avail_service"{
		return t.avail_service(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "checkBalance" { //read a variable
		return t.checkBalance(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//Function to transfer funds from one account to another
func (t *SimpleChaincode) avail_service(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var vehicle Vehicle
	var service Service

	vehicleAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	err = json.Unmarshal(vehicleAsBytes, &vehicle)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of accFrom")
	}

	serviceAsBytes, err := stub.GetState(args[1])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[1] + "\"}"
		return nil, errors.New(jsonResp)
	}
	err = json.Unmarshal(serviceAsBytes, &service)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of accTo")
	}

	if vehicle.Balance < service.Cost {
		return nil, errors.New("Insufficient Balance")
	}

	vehicle.Balance = vehicle.Balance - service.Cost

	b, err := json.Marshal(vehicle)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for account 1")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//Function to check balances
func (t *SimpleChaincode) checkBalance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var vehicle Vehicle
	var balance int

	vehicleAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	err = json.Unmarshal(vehicleAsBytes, &vehicle)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of accCurrent")
	}

	balance = vehicle.Balance
	balanceAsBytes, err := json.Marshal(balance)
	if err != nil {
		return nil, errors.New("Failed to marshal balance")
	}

	return balanceAsBytes, nil
}
