package main

import (
	"bytes"
	"encoding/json"
	mongodb "github.com/godarkproject/ZanoStackTx/pkg/storage/mongodb/read"
	mongodb2 "github.com/godarkproject/ZanoStackTx/pkg/storage/mongodb/update"
	"io"
	"log"
	"net/http"
	"slices"
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
			UnlockedBalance      int   `json:"unlocked_balance"`
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
			Fee             int64    `json:"fee"`
			Height          int      `json:"height"`
			IsIncome        bool     `json:"is_income"`
			IsMining        bool     `json:"is_mining"`
			IsMixing        bool     `json:"is_mixing"`
			IsService       bool     `json:"is_service"`
			PaymentId       string   `json:"payment_id"`
			RemoteAddresses []string `json:"remote_addresses"`
			ShowSender      bool     `json:"show_sender"`
			Subtransfers    []struct {
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

func main() {
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
			"exclude_unconfirmed": true
		  }
		}`

	request, err := http.NewRequest("POST", "http://127.0.0.1:11212/json_rpc", bytes.NewBuffer([]byte(jsonBody)))
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

		if confirmations >= 1 && transfer.PaymentId != "" {

			// fetch user details
			user, err := mongodb.FetchUser(transfer.PaymentId)
			if err != nil {
				log.Println("user doesnt exist")
			} else {

				var userTxHashes []string
				for _, hash := range user.ZanoDeposits {
					userTxHashes = append(userTxHashes, hash.TxHash)
				}

				if !slices.Contains(userTxHashes, transfer.TxHash) {
					mongodb2.AddTx(transfer.TxHash, transfer.Amount, user.ID)
				}
			}
		}
	}
}
