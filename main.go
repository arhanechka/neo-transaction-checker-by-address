package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/CityOfZion/neo-go-sdk/neo"
	"github.com/CityOfZion/neo-go-sdk/neo/models"
)

//from config.json
type Configuration struct {
	AccountAddress string
	Host           string
}

type ConfirmedTransactionsByAddress struct {
	vouts       []models.Vout
	transaction models.Transaction
	block       models.Block
}

var confirmedTransactions []ConfirmedTransactionsByAddress

func main() {
	nodeURI := "http://localhost:10332"
	//   nodeURI := "http://test1.cityofzion.io:8880"
	client := neo.NewClient(nodeURI)

	ok := client.Ping()
	if !ok {
		log.Fatal("Unable to connect to NEO node")
	}

	//1.Generate new address and send to client
	myAddress := GetNewAddress()
	SendNewAddress(myAddress)
	confirmedTransactions := make([]ConfirmedTransactionsByAddress, 0)
	//todo uncomment this block
	// bestBlockHash, err := client.GetBestBlockHash()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// currentBlock, err := client.GetBlockByHash(bestBlockHash)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(bestBlockHash)
	//only for test
	currentBlock, err := client.GetBlockByIndex(1825668)
	if err != nil {
		log.Fatal(err)
	}
	// get number of last block
	lastCheckedBlock := currentBlock.Index
	expectedSum := 0.13
	totalBlocksCount, err := client.GetBlockCount()
	sliceOfBlocks := make([]models.Block, 0)
	getAllLastBlocks(&client, lastCheckedBlock, totalBlocksCount, &sliceOfBlocks)

	var vouts []models.Vout

	checkAllTransactions(&sliceOfBlocks, myAddress, &vouts, &confirmedTransactions)
	for _, structure := range confirmedTransactions {
		fmt.Printf("transaction id = %v \n", structure.transaction.ID)
		fmt.Printf("block hash = %v \n", structure.block.Hash)
		fmt.Printf("block confirmations = %v \n", structure.block.Confirmations)
		fmt.Printf("vouts = %v \n", structure.vouts)
	}
	if ifSumEqual(&vouts, expectedSum) {
		fmt.Println("Winner!")
	}

}

func getAllLastBlocks(client *neo.Client, lastCheckedBlock int64, totalBlocksCount int64, sliceOfBlocks *[]models.Block) {
	//
	for i := lastCheckedBlock; i <= totalBlocksCount; i++ {
		newBlock, err := client.GetBlockByIndex(i)
		if err != nil {
			log.Fatal(err)
		}
		*sliceOfBlocks = append(*sliceOfBlocks, *newBlock)
	}
}

func checkAllTransactions(sliceOfBlocks *[]models.Block, myAddress string, vouts *[]models.Vout, confirmedTransactions *[]ConfirmedTransactionsByAddress) {
	for _, block := range *sliceOfBlocks {
		transactions := block.Transactions
		for _, trans := range transactions {
			if checkVouts(trans, myAddress, vouts) {
				newConfirmedStructure := addToStructure(trans, *vouts, block)
				*confirmedTransactions = append(*confirmedTransactions, newConfirmedStructure)
			}
		}
	}
}

func addToStructure(transaction models.Transaction, vouts []models.Vout, block models.Block) ConfirmedTransactionsByAddress {
	confirmedTransaction := ConfirmedTransactionsByAddress{}
	confirmedTransaction.block = block
	confirmedTransaction.transaction = transaction
	confirmedTransaction.vouts = vouts
	return confirmedTransaction
}
func checkVouts(transaction models.Transaction, address string, vouts *[]models.Vout) bool {
	ifvoutExist := false
	for _, vout := range transaction.Vout {
		//fmt.Println(vout.Address)
		if vout.Address == address {
			*vouts = append(*vouts, vout)
			ifvoutExist = true
		}
	}
	return ifvoutExist
}

func ifSumEqual(vouts *[]models.Vout, expectedSum float64) bool {
	var sum float64
	for _, vout := range *vouts {
		value, err := strconv.ParseFloat(vout.Value, 64)
		if err != nil {
			log.Fatal(err)
		}
		sum += value
	}
	return sum == expectedSum
}

func SendNewAddress(address string) {
	//todo send pay address to Front
}
