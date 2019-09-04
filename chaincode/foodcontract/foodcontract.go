package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type FoodContract struct {
}

type food struct {
	OrderId                string
	FoodId                 string
	ConsumerId             string
	ManufactureId          string
	WholesalerId           string
	RetailerId             string
	LogisticsId            string
	Status                 string
	FoodName               string
	Origin                 string
	RawFoodProcessDate     string
	ManufactureName        string
	ManufactureProcessDate string
	WholesaleName          string
	WholesaleProcessDate   string
	LogisticsProvider      string
	ShippingProcessDate    string
	RetailName             string
	RetailProcessDate      string
	OrderPrice             int
	ShippingPrice          int
	DeliveryDate           string
}

func (t *FoodContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return setupFoodSupplyChainOrder(stub)
}

func (t *FoodContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createRawFood" {
		return t.createRawFood(stub, args)
	} else if function == "manufactureProcessing" {

		return t.manufactureProcessing(stub, args)
	} else if function == "wholesalerDistribute" {

		return t.wholesalerDistribute(stub, args)
	} else if function == "initiateShipment" {

		return t.initiateShipment(stub, args)
	} else if function == "deliverToRetail" {

		return t.deliverToRetail(stub, args)
	} else if function == "completeOrder" {

		return t.completeOrder(stub, args)
	} else if function == "query" {

		return t.query(stub, args)
	}

	return shim.Error("Invalid function name")
}

func setupFoodSupplyChainOrder(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	orderId := args[0]
	//consumerId := args[1]
	//orderPrice, _ := strconv.Atoi(args[2])
	///shippingPrice, _ := strconv.Atoi(args[3])
	foodContract := food{
		OrderId: orderId,
		//ConsumerId:    consumerId,
		//OrderPrice:    orderPrice,
		//ShippingPrice: shippingPrice,
		Status: "商品產生"}

	foodBytes, _ := json.Marshal(foodContract)
	stub.PutState(foodContract.OrderId, foodBytes)

	return shim.Success(nil)
}

func (f *FoodContract) createRawFood(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	foodBytes, _ := stub.GetState(orderId)
	fd := food{}
	json.Unmarshal(foodBytes, &fd)

	if fd.Status == "商品產生" {
		fd.Origin = args[1]
		fd.FoodName = args[2]
		currentts := time.Now()
		fd.RawFoodProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "原物料已取得"
	} else {
		fmt.Printf("Order not initiated yet")
	}

	foodBytes, _ = json.Marshal(fd)
	stub.PutState(orderId, foodBytes)

	return shim.Success(nil)
}

func (f *FoodContract) manufactureProcessing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	foodBytes, err := stub.GetState(orderId)
	fd := food{}
	err = json.Unmarshal(foodBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "原物料已取得" {
		fd.ManufactureName = args[1]
		currentts := time.Now()
		fd.ManufactureProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "加工處理中"
	} else {
		fd.Status = "Error"
		fmt.Printf("Raw food not initiated yet")
	}

	foodBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, foodBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
func (f *FoodContract) wholesalerDistribute(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	foodBytes, err := stub.GetState(orderId)
	fd := food{}
	err = json.Unmarshal(foodBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "加工處理中" {
		fd.WholesalerId = "Wholesaler_1"
		fd.WholesaleName = args[1]
		currentts := time.Now()
		fd.WholesaleProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "批發處理中"
	} else {
		fd.Status = "Error"
		fmt.Printf("Manufacture not initiated yet")
	}

	foodBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, foodBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (f *FoodContract) initiateShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	foodBytes, err := stub.GetState(orderId)
	fd := food{}
	err = json.Unmarshal(foodBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "批發處理中" {
		fd.LogisticsId = "LogisticsId_1"
		fd.LogisticsProvider = args[1]
		currentts := time.Now()
		fd.ShippingProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "運送處理中"
	} else {
		fmt.Printf("Wholesaler not initiated yet")
	}

	foodBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, foodBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (f *FoodContract) deliverToRetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	foodBytes, err := stub.GetState(orderId)
	fd := food{}
	err = json.Unmarshal(foodBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "運送處理中" {
		fd.RetailerId = "Retailer_1"
		fd.RetailName = args[1]
		currentts := time.Now()
		fd.RetailProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "已上架"

	} else {
		fmt.Printf("Shipment not initiated yet")
	}

	foodBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, foodBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (f *FoodContract) completeOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	foodBytes, err := stub.GetState(orderId)
	fd := food{}
	err = json.Unmarshal(foodBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}
	orderPrice, _ := strconv.Atoi(args[2])
	if fd.Status == "已上架" {
		currentts := time.Now()
		fd.DeliveryDate = currentts.Format("2006-01-02 15:04:05")
		fd.ConsumerId = args[1]
		fd.OrderPrice = orderPrice
		fd.Status = "已售出"
	} else {
		fmt.Printf("Retailer not initiated yet")
	}

	foodBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, foodBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (f *FoodContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ENIITY string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected ENIITY Name")
	}

	ENIITY = args[0]
	Avalbytes, err := stub.GetState(ENIITY)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil order for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)
}

func main() {

	err := shim.Start(new(FoodContract))
	if err != nil {
		fmt.Printf("Error creating new Food Contract: %s", err)
	}
}
