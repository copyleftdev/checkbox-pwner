package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	apiEndpoint          = "https://onemillioncheckboxes.com/api/initial-state"
	socketURL            = "wss://onemillioncheckboxes.com/socket.io/?EIO=4&transport=websocket"
	defaultSleepDuration = 200 * time.Millisecond
	defaultReconnectWait = 5 * time.Second
	defaultNumWorkers    = 100
	defaultMaxRetries    = 3
	defaultBatchSize     = 50
)

type BitSet struct {
	bytes      []byte
	checkCount int
	mu         sync.Mutex
}

var (
	sleepDuration time.Duration
	reconnectWait time.Duration
	numWorkers    int
	maxRetries    int
	batchSize     int

	rootCmd = &cobra.Command{
		Use:   "checkbox-pwner",
		Short: "A tool to dominate the One Million Checkboxes site",
		Run:   run,
	}
)

func init() {
	rootCmd.Flags().DurationVar(&sleepDuration, "sleep-duration", defaultSleepDuration, "Sleep duration between operations")
	rootCmd.Flags().DurationVar(&reconnectWait, "reconnect-wait", defaultReconnectWait, "Wait time before reconnecting after failure")
	rootCmd.Flags().IntVar(&numWorkers, "num-workers", defaultNumWorkers, "Number of parallel workers")
	rootCmd.Flags().IntVar(&maxRetries, "max-retries", defaultMaxRetries, "Maximum number of retries for each batch")
	rootCmd.Flags().IntVar(&batchSize, "batch-size", defaultBatchSize, "Number of checkboxes to process in each batch")

	// Configure logrus for better logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewBitSet(base64String string, count int) *BitSet {
	binaryData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		logrus.Fatalf("❌ Failed to decode base64 string: %v", err)
	}
	return &BitSet{
		bytes:      binaryData,
		checkCount: count,
	}
}

func (b *BitSet) Get(index int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	byteIndex := index / 8
	bitOffset := 7 - (index % 8)
	return (b.bytes[byteIndex] & (1 << bitOffset)) != 0
}

func (b *BitSet) Set(index int, value bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	byteIndex := index / 8
	bitOffset := 7 - (index % 8)
	if value {
		b.bytes[byteIndex] |= 1 << bitOffset
	} else {
		b.bytes[byteIndex] &^= 1 << bitOffset
	}
}

type InitialState struct {
	FullState string `json:"full_state"`
	Count     int    `json:"count"`
}

func fetchInitialState() (*BitSet, error) {
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch initial state: %w", err)
	}
	defer resp.Body.Close()

	var state InitialState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, fmt.Errorf("❌ Failed to decode initial state: %w", err)
	}

	return NewBitSet(state.FullState, state.Count), nil
}

func checkBatch(bitset *BitSet, start, end int, conn *wsutil.Writer) error {
	logrus.Infof("🔄 Checking batch from %d to %d", start, end)
	for i := start; i <= end; i += batchSize {
		batchEnd := i + batchSize
		if batchEnd > end {
			batchEnd = end
		}
		var batch []string
		for j := i; j < batchEnd; j++ {
			if !bitset.Get(j) {
				batch = append(batch, fmt.Sprintf(`{"index": %d}`, j))
				bitset.Set(j, true)
			}
		}
		if len(batch) > 0 {
			msg := fmt.Sprintf(`42["toggle_bits", [%s]]`, strings.Join(batch, ","))
			for retries := 0; retries < maxRetries; retries++ {
				logrus.Debugf("📤 Sending batch message: %s", msg)
				if err := wsutil.WriteClientText(conn, []byte(msg)); err != nil {
					logrus.Errorf("❌ Failed to send batch message: %v", err)
					time.Sleep(sleepDuration)
				} else {
					logrus.Infof("✅ Checked checkboxes from %d to %d", i, batchEnd-1)
					break
				}
			}
		}
		time.Sleep(sleepDuration)
	}
	return nil
}

func worker(start, end int, wg *sync.WaitGroup) {
	defer wg.Done()
	logrus.Infof("👷‍♂️ Worker started for range %d to %d", start, end)
	for {
		bitset, err := fetchInitialState()
		if err != nil {
			logrus.Errorf("❌ Failed to fetch initial state: %v", err)
			time.Sleep(reconnectWait)
			continue
		}

		conn, _, _, err := ws.Dial(context.Background(), socketURL)
		if err != nil {
			logrus.Errorf("❌ Failed to connect to WebSocket: %v", err)
			time.Sleep(reconnectWait)
			continue
		}
		writer := wsutil.NewWriter(conn, ws.StateClientSide, ws.OpText)
		if err := checkBatch(bitset, start, end, writer); err != nil {
			logrus.Errorf("❌ Error in checkBatch: %v", err)
			conn.Close()
			time.Sleep(reconnectWait)
			continue
		}
		conn.Close()
	}
}

func run(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i*(1000000/numWorkers), (i+1)*(1000000/numWorkers)-1, &wg)
	}

	logrus.Infof("🚀 Started checking checkboxes with %d workers and %s sleep duration...", numWorkers, sleepDuration)
	wg.Wait()
}

func main() {
	Execute()
}
