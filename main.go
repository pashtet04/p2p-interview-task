package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type BlockLatest struct {
	BlockID struct {
		Hash          string `json:"hash"`
		PartSetHeader struct {
			Total int    `json:"total"`
			Hash  string `json:"hash"`
		} `json:"part_set_header"`
	} `json:"block_id"`
	Block struct {
		Header struct {
			Version struct {
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"version"`
			ChainID     string    `json:"chain_id"`
			Height      string    `json:"height"`
			Time        time.Time `json:"time"`
			LastBlockID struct {
				Hash          string `json:"hash"`
				PartSetHeader struct {
					Total int    `json:"total"`
					Hash  string `json:"hash"`
				} `json:"part_set_header"`
			} `json:"last_block_id"`
			LastCommitHash     string `json:"last_commit_hash"`
			DataHash           string `json:"data_hash"`
			ValidatorsHash     string `json:"validators_hash"`
			NextValidatorsHash string `json:"next_validators_hash"`
			ConsensusHash      string `json:"consensus_hash"`
			AppHash            string `json:"app_hash"`
			LastResultsHash    string `json:"last_results_hash"`
			EvidenceHash       string `json:"evidence_hash"`
			ProposerAddress    string `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs []interface{} `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence []interface{} `json:"evidence"`
		} `json:"evidence"`
		LastCommit struct {
			Height  string `json:"height"`
			Round   int    `json:"round"`
			BlockID struct {
				Hash          string `json:"hash"`
				PartSetHeader struct {
					Total int    `json:"total"`
					Hash  string `json:"hash"`
				} `json:"part_set_header"`
			} `json:"block_id"`
			Signatures []struct {
				BlockIDFlag      string    `json:"block_id_flag"`
				ValidatorAddress string    `json:"validator_address"`
				Timestamp        time.Time `json:"timestamp"`
				Signature        string    `json:"signature"`
			} `json:"signatures"`
		} `json:"last_commit"`
	} `json:"block"`
	SdkBlock struct {
		Header struct {
			Version struct {
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"version"`
			ChainID     string    `json:"chain_id"`
			Height      string    `json:"height"`
			Time        time.Time `json:"time"`
			LastBlockID struct {
				Hash          string `json:"hash"`
				PartSetHeader struct {
					Total int    `json:"total"`
					Hash  string `json:"hash"`
				} `json:"part_set_header"`
			} `json:"last_block_id"`
			LastCommitHash     string `json:"last_commit_hash"`
			DataHash           string `json:"data_hash"`
			ValidatorsHash     string `json:"validators_hash"`
			NextValidatorsHash string `json:"next_validators_hash"`
			ConsensusHash      string `json:"consensus_hash"`
			AppHash            string `json:"app_hash"`
			LastResultsHash    string `json:"last_results_hash"`
			EvidenceHash       string `json:"evidence_hash"`
			ProposerAddress    string `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs []interface{} `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence []interface{} `json:"evidence"`
		} `json:"evidence"`
		LastCommit struct {
			Height  string `json:"height"`
			Round   int    `json:"round"`
			BlockID struct {
				Hash          string `json:"hash"`
				PartSetHeader struct {
					Total int    `json:"total"`
					Hash  string `json:"hash"`
				} `json:"part_set_header"`
			} `json:"block_id"`
			Signatures []struct {
				BlockIDFlag      string    `json:"block_id_flag"`
				ValidatorAddress string    `json:"validator_address"`
				Timestamp        time.Time `json:"timestamp"`
				Signature        string    `json:"signature"`
			} `json:"signatures"`
		} `json:"last_commit"`
	} `json:"sdk_block"`
}

type NetInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Listening bool          `json:"listening"`
		Listeners []string      `json:"listeners"`
		NPeers    string        `json:"n_peers"`
		Peers     []interface{} `json:"peers"`
	} `json:"result"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	client = &http.Client{}

	cosmosApiEndpoint     = getEnv("COSMOS_API", "http://localhost:11317/")
	tendermintApiEndpoint = getEnv("TENDERMINT_API", "http://localhost:26657/")

	up = prometheus.NewDesc(
		prometheus.BuildFQName("", "", "up"),
		"Was the last API request successful.",
		nil, nil,
	)

	cosmos_latest_block_height = prometheus.NewDesc(
		prometheus.BuildFQName("", "", "cosmos_latest_block_height"),
		"The latest block id hash",
		nil, nil,
	)

	cosmos_latest_block_timestamp = prometheus.NewDesc(
		prometheus.BuildFQName("", "", "cosmos_latest_block_timestamp"),
		"Unsync node in ms",
		nil, nil,
	)

	cosmos_number_of_peers = prometheus.NewDesc(
		prometheus.BuildFQName("", "", "cosmos_number_of_peers"),
		"Number of peers",
		nil, nil,
	)
)

type Exporter struct {
	cosmosApiEndpoint string
}

func NewExporter(cosmosApiEndpoint string) *Exporter {
	return &Exporter{
		cosmosApiEndpoint: cosmosApiEndpoint,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- cosmos_latest_block_height
	ch <- cosmos_latest_block_timestamp
	ch <- cosmos_number_of_peers
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	s, err := e.GetLatestBlockHash()
	r, _ := strconv.ParseFloat(s, 64)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			cosmos_latest_block_height, prometheus.UntypedValue, 0,
		)
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		cosmos_latest_block_height, prometheus.UntypedValue, r,
	)

	t, _ := e.GetLatestBlockTime()
	f := float64(t)
	ch <- prometheus.MustNewConstMetric(
		cosmos_latest_block_timestamp, prometheus.UntypedValue, f,
	)

	p, _ := e.GetCountPeers()
	cp, _ := strconv.ParseFloat(p, 64)
	ch <- prometheus.MustNewConstMetric(
		cosmos_number_of_peers, prometheus.UntypedValue, cp,
	)

}

func CosmosApiReq(method string, url string) (BlockLatest, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result BlockLatest
	json.Unmarshal([]byte(body), &result)
	return result, err
}

func TendermintApiReq(method string, url string) (NetInfo, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result NetInfo
	json.Unmarshal([]byte(body), &result)
	return result, err
}

func (e *Exporter) GetLatestBlockHash() (string, error) {
	result, err := CosmosApiReq("GET", cosmosApiEndpoint+"cosmos/base/tendermint/v1beta1/blocks/latest")
	return result.Block.Header.Height, err
}

func (e *Exporter) GetLatestBlockTime() (int64, error) {
	result, err := CosmosApiReq("GET", cosmosApiEndpoint+"cosmos/base/tendermint/v1beta1/blocks/latest")
	unsynced_ms := time.Now().Unix() - result.Block.Header.Time.Unix()
	return unsynced_ms, err
}

func (e *Exporter) GetCountPeers() (string, error) {
	result, err := TendermintApiReq("GET", tendermintApiEndpoint+"net_info")
	return result.Result.NPeers, err
}

func main() {
	exporter := NewExporter(cosmosApiEndpoint)
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))

}

// curl -X 'GET' \
//   'http://localhost:8080/cosmos/base/tendermint/v1beta1/blocks/latest' \
//   -H 'accept: application/json'
