package cmd

import (
	"dag-balances/domain/be"
	dos "dag-balances/domain/opensearch"
	"dag-balances/infrastructure/opensearch"
	"dag-balances/pkg/color"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(balancesCmd)

	balancesCmd.Flags().String("opensearch", "", "Opensearch url")
	balancesCmd.Flags().Int64("from", 0, "Starting ordinal")
	balancesCmd.Flags().Int64("to", 0, "Ending ordinal")
	balancesCmd.Flags().Int("workers", 10, "No. of parallel workers")
	balancesCmd.Flags().Bool("dry-run", false, "When true, works in read-only mode")

	balancesCmd.MarkFlagRequired("opensearch")
	balancesCmd.MarkFlagRequired("from")
	balancesCmd.MarkFlagRequired("to")
}

func CheckOrdinal(ordinal int64, client dos.OpenSearch, dryRun bool) error {
	if ordinal%10000 == 0 || Verbose {
		fmt.Printf("%s\n", color.Ize(color.Yellow, fmt.Sprintf("Ordinal: %d", ordinal)))
	}

	txs, err := client.GetTransactions(ordinal)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		address := tx.Source
		snapshotHash := tx.SnapshotHash
		balance, err := client.GetBalance(address, ordinal)

		if balance == nil && err == nil {
			fmt.Printf("%s\n", color.Ize(color.Red, fmt.Sprintf("Balance does not exist on ordinal: %d for address: %s", ordinal, address)))

			if !dryRun {
				err := client.PutBalance(be.Balance{
					Address:         address,
					Balance:         0,
					SnapshotOrdinal: ordinal,
					SnapshotHash:    snapshotHash,
				})
				if err != nil {
					fmt.Printf("%s\n", color.Ize(color.Red, fmt.Sprintf("--- err: %s", err)))
				} else {
					fmt.Printf("%s\n", color.Ize(color.Green, fmt.Sprintf("Fixed balance on ordinal: %d for address: %s", ordinal, address)))
				}
			}

		}
	}
	return nil
}

func worker(jobChan <-chan int64, wg *sync.WaitGroup, client dos.OpenSearch, dryRun bool) {
	defer wg.Done()

	for job := range jobChan {
		err := CheckOrdinal(job, client, dryRun)
		if err != nil {
			fmt.Printf("%s\n", color.Ize(color.Red, fmt.Sprintf("--- err: %s", err)))
		}
	}
}

var balancesCmd = &cobra.Command{
	Use:   "fix-balances",
	Short: "Fix wrongly indexed zero balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("opensearch")
		client := opensearch.GetClient(url)

		from, _ := cmd.Flags().GetInt64("from")
		to, _ := cmd.Flags().GetInt64("to")
		workers, _ := cmd.Flags().GetInt("workers")

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		var wg sync.WaitGroup
		jobChan := make(chan int64)

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go worker(jobChan, &wg, client, dryRun)
		}

		var ordinal int64
		for ordinal = from; ordinal <= to; ordinal++ {
			jobChan <- ordinal
		}

		close(jobChan)
		wg.Wait()

		return nil
	},
}
