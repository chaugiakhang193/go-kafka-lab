// Command keypartition produces keyed and unkeyed records to one topic and
// prints the partition and offset the broker assigned to each record.
//
// The client uses the franz-go default partitioner. Records with a key are
// placed by hashing the key. Records without a key are placed by the
// partitioner's own strategy.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	brokers := flag.String("brokers", "localhost:9092", "comma-separated bootstrap servers")
	topic := flag.String("topic", "lab.products", "topic to produce to, created beforehand")
	hold := flag.Duration("hold", 0, "keep the client open this long after producing")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Split(*brokers, ",")...),
		kgo.DefaultProduceTopic(*topic),
	)
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	defer client.Close()

	produceCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := produceAll(produceCtx, client, labRecords()); err != nil {
		log.Fatal(err)
	}

	if *hold > 0 {
		fmt.Printf("holding the client open for %s\n", *hold)
		select {
		case <-ctx.Done():
		case <-time.After(*hold):
		}
	}
}

// labRecords returns three rounds of keys A, B and C interleaved, followed by
// six records without a key.
func labRecords() []*kgo.Record {
	var records []*kgo.Record
	for round := 1; round <= 3; round++ {
		for _, key := range []string{"A", "B", "C"} {
			records = append(records, newRecord(key, fmt.Sprintf("%s%d", key, round)))
		}
	}
	for i := 1; i <= 6; i++ {
		records = append(records, newRecord("", fmt.Sprintf("n%d", i)))
	}
	return records
}

// newRecord leaves Key nil for an empty key. A non-nil empty slice is still a
// key and would be hashed like any other.
func newRecord(key, value string) *kgo.Record {
	r := &kgo.Record{Value: []byte(value)}
	if key != "" {
		r.Key = []byte(key)
	}
	return r
}

// produceAll sends the records one at a time and waits for each result. The
// printed lines follow the send order.
func produceAll(ctx context.Context, client *kgo.Client, records []*kgo.Record) error {
	fmt.Printf("%-6s %-6s %-9s %s\n", "key", "value", "partition", "offset")
	for _, r := range records {
		if err := client.ProduceSync(ctx, r).FirstErr(); err != nil {
			return fmt.Errorf("produce %q: %w", r.Value, err)
		}
		fmt.Printf("%-6s %-6s %-9d %d\n", displayKey(r.Key), r.Value, r.Partition, r.Offset)
	}
	return nil
}

func displayKey(key []byte) string {
	if key == nil {
		return "(nil)"
	}
	return string(key)
}
