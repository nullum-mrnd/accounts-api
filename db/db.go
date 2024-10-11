package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Db struct {
	*pgx.Conn
}

type Account struct {
	UserID   int    `db:"id"`
	Username string `db:"username"`
	Phone    string `db:"phone"`
}

func NewDB(login, passwd, server, table string) Db {
	conn, err := pgx.Connect(context.Background(), "postgres://"+login+":"+passwd+"@"+server+"/"+table+"?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return Db{conn}
}

func (d Db) GetAccounts() ([]Account, error) {
	query := `
		SELECT * FROM accounts
	`
	rows, err := d.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error Querying the Table")
		return nil, err
	}
	defer rows.Close()
	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Account])
	return accounts, err
}

func (d Db) GetAccountByID(AccountID int) (*Account, error) {
	query := `
		SELECT * FROM accounts WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id": AccountID,
	}
	row := d.QueryRow(context.Background(), query, args)

	var acc Account
	err := row.Scan(&acc.UserID, &acc.Username, &acc.Phone)

	if err != nil {
		fmt.Println("Error Retrieving Account")
		return nil, err
	}

	return &acc, nil
}

func (d Db) CreateAccount(username, phone string) (*Account, error) {
	query := `
		INSERT INTO accounts (username, phone) VALUES (@username, @phone) RETURNING id
	`
	args := pgx.NamedArgs{
		"username": username,
		"phone":    phone,
	}

	var id int
	err := d.QueryRow(context.Background(), query, args).Scan(&id)

	fmt.Println("id: ", id)
	if err != nil {
		fmt.Println("Error Inserting Account")
		fmt.Println(err)
		return nil, err
	}
	return &Account{UserID: id, Username: username, Phone: phone}, nil
}

func (d Db) UpdateAccount(AccountID int, acc Account) error {
	query := `
		UPDATE accounts
		SET username = @username, phone = @phone
		WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id":       AccountID,
		"username": acc.Username,
		"phone":    acc.Phone,
	}
	_, err := d.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Println("Error Updating Account")
		return err
	}
	return nil
}

func (d Db) DeleteAccount(AccountID int) error {
	query := `
		DELETE FROM accounts WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id": AccountID,
	}

	_, err := d.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Println("Error Deleting Account")
		return err
	}
	return nil
}
