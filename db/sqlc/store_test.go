package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// key for context to identify the transaction name

	//now time to run n concurrent transfer transactions
	n := 10
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)
	amount := int64(10) // fixed amount for simplicity, can be randomized if needed
	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx-%d", i+1)
		go func() {
			// run transfer transaction
			ctx := context.Background()
			res, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- res
		}()
	}

	// Only launch n goroutines using the fixed amount above
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		//check if transfer exists in the database
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		//check if from entry exists in the database
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)

		//check if the to entry exists in the database (fix)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO : check account balance
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println("fromAccount balance:", fromAccount.Balance)
		fmt.Println("toAccount balance:", toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)     // diff should be same
		require.True(t, diff1 > 0)         //should be a positive number
		require.True(t, diff1%amount == 0) //should be a multiple of amount

		k := int(diff1 / amount)
		require.True(t, (k >= 1 && k <= n)) //should be an integer between 1 and n
		require.NotContains(t, existed, k)  //should not be repeated
		existed[k] = true                   //mark as existed

	}
	//check final update balance
	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Print("Final account balances:\n")
	fmt.Printf("Account 1 ID: %d, Balance: %d\n", updateAccount1.ID, updateAccount1.Balance)
	fmt.Printf("Account 2 ID: %d, Balance: %d\n", updateAccount2.ID, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// key for context to identify the transaction name

	//now time to run n concurrent transfer transactions
	n := 10
	errs := make(chan error, n)
	amount := int64(10) // fixed amount for simplicity, can be randomized if needed
	for i := 0; i < n; i++ {
		FromAccountID := account1.ID
		ToAccountID := account2.ID

		//toggle
		if i%2 == 0 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID
		}
		// txName := fmt.Sprintf("tx-%d", i+1)
		go func() {
			// run transfer transaction
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	//also there are same 5 transactions to and fro from accounts,
	//thus final balances must remain the same

	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Print("Final account balances:\n")
	fmt.Printf("Account 1 ID: %d, Balance: %d\n", updateAccount1.ID, updateAccount1.Balance)
	fmt.Printf("Account 2 ID: %d, Balance: %d\n", updateAccount2.ID, updateAccount2.Balance)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)

}
