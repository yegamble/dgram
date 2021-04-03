package wallet

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/converter"
	"github.com/iotaledger/iota.go/transaction"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/pebbe/zmq4"
	"log"
	"math/big"
	"strings"
)

type IotaWallet struct {
	seed       string `json:"seed"`
	privateKey string `json:"privateKey"`
}

var iotaAPI *api.API
var version = "undefined"
var node = "https://nodes.thetangle.org:443"

const letters = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ" //pool of letters to generate IOTA seed

func init() {
	var err error
	iotaAPI, err = api.ComposeAPI(api.HTTPClientSettings{URI: node})
	if err != nil {
		panic(err)
		return
	}
}

func CheckTransactions() error {

	client, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return err
	}

	// Make sure the connection is closed after stopping the program
	defer client.Close()

	// Connect to a Tangle node's address
	client.Connect(node)

	// Subscribe to both tx and sn (confirmed tx) topics
	client.SetSubscribe("tx")
	client.SetSubscribe("sn")

	// Keep looping for messages
	for {
		msg, _ := client.RecvMessage(0)
		for _, str := range msg {

			// Split the fields by the space character
			words := strings.Fields(str)

			if words[0] == "tx" {
				fmt.Println("New transaction: ", words[1])
			}
			if words[0] == "sn" {
				fmt.Println("Confirmed transaction: ", words[2], "for milestone", words[1])
			}
		}

	}
}

func RetrieveLastBundle(hash string) (bundle.Bundle, error) {

	tryteHash, err := trinary.NewTrytes(hash)
	if err != nil {
		return nil, err
	}

	bundle, err := iotaAPI.GetBundle(tryteHash)
	if err != nil {
		// handle error
		return nil, err
	}
	return bundle, nil
}

func UpdateProfileAddress(address string, seed string, profile string) (bundle.Bundle, string, error) {

	const depth = 3
	const minimumWeightMagnitude = 14

	trinaryAddress, _ := trinary.NewTrytes(address)
	trinarySeed, _ := trinary.NewTrytes(seed)
	trinaryProfile, _ := converter.ASCIIToTrytes(profile)
	trinaryTag, _ := trinary.NewTrytes("DGRAM")

	transfers := bundle.Transfers{
		{
			Address: trinaryAddress,
			Value:   0,
			Message: trinaryProfile,
			Tag:     trinaryTag,
		},
	}

	//// create inputs for the transfer
	//inputs := []api.Input{
	//	{
	//		// must be 90 trytes long (inlcude the checksum)
	//		Address:  trinaryAddress,
	//		Security: consts.SecurityLevelLow,
	//		KeyIndex: 0,
	//		Balance:  0,
	//	},
	//}

	// we don't need to set the security level or timestamp in the options because we supply
	// the input and remainder addresses.
	prepTransferOpts := api.PrepareTransfersOptions{}

	// prepare the transfer by creating a bundle with the given transfers and inputs.
	// the result are trytes ready for PoW.
	trytes, err := iotaAPI.PrepareTransfers(trinarySeed, transfers, prepTransferOpts)
	if err != nil {
		// handle error
		return nil, "", err
	}

	resultBundle, err := iotaAPI.SendTrytes(trytes, depth, minimumWeightMagnitude)
	if err != nil {
		return nil, "", err
	}

	err = bundle.ValidBundle(resultBundle)
	if err != nil {
		return nil, "", err
	}

	return resultBundle, bundle.TailTransactionHash(resultBundle), nil
}

func GetTransactionJSON(transactionHash string) (string, error) {

	bundle, err := iotaAPI.GetBundle(transactionHash)
	if err != nil {
		return "", err
	}

	bundleHashArray := make([]trinary.Hash, len(bundle))
	var signatureMessage bytes.Buffer

	for i := range bundle {
		bundleHashArray = append(bundleHashArray, bundle[i].SignatureMessageFragment)
		signatureMessage.Write([]byte(bundle[i].SignatureMessageFragment))
	}

	resultTryte := strings.TrimSpace(strings.Join(bundleHashArray[:], ""))
	resultAscii, err := converter.TrytesToASCII(resultTryte)

	resultNullRemoved := bytes.Trim([]byte(resultAscii), "\x00")

	if err != nil {
		return "", err
	}

	fmt.Println(resultAscii)

	return strings.TrimSpace(string(resultNullRemoved)), nil
}

func GenerateNewWallet() (string, string, string) {
	newWallet := new(IotaWallet)
	newWallet.seed, _ = GenerateRandomSeed()

	address, err := GenerateNewAddress(newWallet.seed)
	if err != nil {
		log.Fatal(err)
		return "", "", ""
	}

	return newWallet.seed, address, ""
}

// i req: addresses, The addresses of which to get the bundles of.
// i: inclusionState, Whether to set the persistence field on the transactions.
// o: Bundles, The bundles gathered of the given addresses.
// o: error, Returned for invalid parameters and internal errors.
func GetBundlesFromAddresses(inputAddress string) (bundle.Bundles, error) {

	address, _ := trinary.NewTrytes(inputAddress)

	addresses := trinary.Hashes{
		address,
	}

	bundles, err := iotaAPI.GetBundlesFromAddresses(addresses)
	if err != nil {
		// handle error
		return nil, err
	}
	return bundles, nil
}

// i req: query, The object defining the transactions to search for.
// o: Hashes, The Hashes of the query result.
// o: error, Returned for invalid query objects and internal errors.
func GetTransactions(transaction []string) (transaction.Transactions, error) {

	txHashes, err := iotaAPI.FindTransactionObjects(api.FindTransactionsQuery{
		Approvees: []trinary.Trytes{
			"DJDMZD9G9VMGR9UKMEYJWYRLUDEVWTPQJXIQAAXFGMXXSCONBGCJKVQQZPXFMVHAAPAGGBMDXESTZ9999",
		},
	})

	if err != nil {
		// handle error
		return nil, err
	}

	return txHashes, nil
}

func GenerateNewAddress(seed string) (string, error) {
	addr, err := iotaAPI.GetNewAddress(seed, api.GetNewAddressOptions{Index: 0})
	if err != nil {
		// handle error
		log.Fatal(err)
		return "", err
	}
	return addr[0], nil
}

func generateRandomInts(n int) ([]int64, error) {
	ints := make([]int64, n)

	for i := range ints {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(27))

		if err != nil {
			return nil, err
		}

		ints[i] = randomInt.Int64()
	}

	return ints, nil
}

func intToCharByte(i int64) byte {
	return byte(letters[i])
}

func GenerateRandomSeed() (string, error) {
	ints, err := generateRandomInts(81)

	if err != nil {
		return "", err
	}

	token := make([]byte, 81)

	for i, x := range ints {
		token[i] = intToCharByte(x)
	}

	return string(token), nil
}
