package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
// each query performs a single operation
// but in case of txn, we need a series of operations to be performed as a single unit
// so we need a way to execute a series of operations as a single unit

// Store interface defines the methods that a store must implement
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// this is called composition over inheritance
// we embedded queries into store, now all func methods from queries are available to store
// and also we can use all func methods from queries in store
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// below is the method that executes    a transaction
// it takes a context and a function that takes a queries pointer and returns an error
// it creates a new transaction and executes the function
// it then commits the transaction
// it then returns the error
// it then rolls back the transaction if there is an error
// it then returns the error
// execTx executes a function within a database transaction
// it takes a context and a function that takes a queries pointer and returns an er ror
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	//BeginTx tells the database “start a transaction.”
	// From now on, changes are temporary and isolated.
	// If the transaction is successful, we can commit it.
	if err != nil {
		return err
	}

	q := New(tx) // create a new queries pointer with the transaction
	// this is the queries pointer that will be used to execute the queries

	// execute the function with the queries pointer
	err = fn(q)
	//h tx-bound Queries” means - execTx calls your function once: fn(q). -
	// Inside that function, call as many q.* methods as you want;
	// they all use the same q.db (the same *sql.Tx). -
	// When fn returns: - if error → Rollback - if nil → Commit
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr) //incase rollback error also occurs, we return both errors
		}
		return err
	}
	return tx.Commit()
}

//execTx is unexported, so it can only be used within the package
//we shall provide an exported func for each particular transaction

// TransferTxParams contains the input parameters for the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// var txKey = struct{}{} // define a key for the transaction name in the context

// TransferTx performs a money transfer from one account to another
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey) //get value of txKey from the context

		// fmt.Println(txName, " Create Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// create entries for the transfer transaction

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, " Creating entry 2")
		// create entry for the to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, " get account 1")

		//replace get accout for update and updte account -> addAccountBalance singele query
		// get the from account for update
		// acc1, err := q.GetAccountsForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return fmt.Errorf("failed to get from account: %w", err)
		// }
		// fmt.Println(txName, "update account 1 balance")
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil

		//for amount
		//get account -> update account ka balance
		//need a locking mechanism here gng

		//ex

	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1, // add the amount to the account
	})
	if err != nil {
		return
		//return with no params is basicallt same as return account1, account2, err
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2, // add the amount to the account
	})
	if err != nil {
		return
	}
	return account1, account2, nil
	// return account1, account2, nil //return the accounts and no error
}
