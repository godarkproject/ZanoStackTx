package main

import (
	"bytes"
	"encoding/json"
	mongodb "github.com/godarkproject/ZanoStackTx/pkg/storage/mongodb/read"
	mongodb2 "github.com/godarkproject/ZanoStackTx/pkg/storage/mongodb/update"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"time"
)

type GetPaymentsRes struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		LastItemIndex int `json:"last_item_index"`
		Pi            struct {
			Balance              int64 `json:"balance"`
			CurentHeight         int   `json:"curent_height"`
			TransferEntriesCount int   `json:"transfer_entries_count"`
			TransfersCount       int   `json:"transfers_count"`
			UnlockedBalance      int64 `json:"unlocked_balance"`
		} `json:"pi"`
		TotalTransfers int `json:"total_transfers"`
		Transfers      []struct {
			Amount          int64  `json:"amount"`
			Comment         string `json:"comment"`
			EmployedEntries struct {
				Receive []struct {
					Amount  int64  `json:"amount"`
					AssetId string `json:"asset_id"`
					Index   int    `json:"index"`
				} `json:"receive"`
			} `json:"employed_entries"`
			Fee            int64  `json:"fee"`
			Height         int    `json:"height"`
			IsIncome       bool   `json:"is_income"`
			IsMining       bool   `json:"is_mining"`
			IsMixing       bool   `json:"is_mixing"`
			IsService      bool   `json:"is_service"`
			PaymentId      string `json:"payment_id"`
			ServiceEntries []struct {
				Body        string `json:"body"`
				Flags       int    `json:"flags"`
				Instruction string `json:"instruction"`
				ServiceId   string `json:"service_id"`
			} `json:"service_entries"`
			ShowSender   bool `json:"show_sender"`
			Subtransfers []struct {
				Amount   int64  `json:"amount"`
				AssetId  string `json:"asset_id"`
				IsIncome bool   `json:"is_income"`
			} `json:"subtransfers"`
			Timestamp             int    `json:"timestamp"`
			TransferInternalIndex int    `json:"transfer_internal_index"`
			TxBlobSize            int    `json:"tx_blob_size"`
			TxHash                string `json:"tx_hash"`
			TxType                int    `json:"tx_type"`
			UnlockTime            int    `json:"unlock_time"`
		} `json:"transfers"`
	} `json:"result"`
}

func clearScreen(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func init() {

	// clear screen
	switch runtime.GOOS {
	case "darwin":
		clearScreen("clear")
	case "linux":
		clearScreen("clear")
	case "windows":
		clearScreen("cmd", "/c", "cls")
	default:
		clearScreen("clear")
	}
}

func getEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

func monitorTx() {
	// Access the loaded environment variables
	mongoUri := getEnvVar("MONGO_URI_DEV")

	jsonBody := `
		{
		  "jsonrpc": "2.0",
		  "id": 0,
		  "method": "get_recent_txs_and_info",
		  "params": {
			"offset": 0,
			"update_provision_info": true,
			"exclude_mining_txs": true,
			"count": 100,
			"order": "FROM_END_TO_BEGIN",
			"exclude_unconfirmed": false
		  }
		}`

	request, err := http.NewRequest("POST", "http://127.0.0.1:11214/json_rpc", bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		log.Println("error 1 with wallet POST")
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Println("error 2 with wallet POST")
	}
	defer func(Body io.ReadCloser) {

	}(res.Body)

	body, _ := io.ReadAll(res.Body)
	data := GetPaymentsRes{}
	_ = json.Unmarshal(body, &data)

	for _, transfer := range data.Result.Transfers {
		confirmations := int64(data.Result.Pi.CurentHeight) - int64(transfer.Height)

		log.Println(transfer.TxHash)

		if confirmations < 10 && transfer.PaymentId != "" && transfer.IsIncome {
			log.Printf("Transaction confirming for %d $ZANO, %d confirmations left.\n", transfer.Height, 10-confirmations)
		}

		if confirmations >= 10 && transfer.PaymentId != "" && transfer.IsIncome {

			// fetch user details
			user, err := mongodb.FetchUser(mongoUri, transfer.PaymentId)
			if err == nil {
				var userTxHashes []string
				for _, hash := range user.ZanoDeposits {
					userTxHashes = append(userTxHashes, hash.TxHash)
				}

				if !slices.Contains(userTxHashes, transfer.TxHash) {
					mongodb2.AddTx(mongoUri, transfer.TxHash, transfer.Amount, user.ID)

					newBalance := user.Balance + transfer.Amount
					updated, err := mongodb2.UpdateBalance(mongoUri, newBalance, user.ID)
					if err != nil {
						panic(err)
					}

					log.Printf("balance updated: %v \n", updated)
				}
			} else {
				log.Println("no user exists with payment id")
			}
		}

	}
}

func main() {
	log.Println("Monitoring incoming transactions")
	for {
		monitorTx()
		time.Sleep(60 * time.Second)
	}
}
