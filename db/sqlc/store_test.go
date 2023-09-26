package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfertx(t *testing.T) {
	store := NewStore(textBD)

	account1 := createRandomTestAccounts(t)
	account2 := createRandomTestAccounts(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTXParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)

		// _, err = store.GetTransfer(context.Background(), transfer.ID)
		// require.NoError(t, err)

	}
}
