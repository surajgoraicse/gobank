package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

// constructor function for postgres store
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=surajgoraicse password=surajgoraicse dbname=gobank host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB Connection error ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Ping failed:", err)
		return nil, err
	}
	return &PostgresStore{db: db}, nil

}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account(
		id serial primary key ,
		first_name varchar(50),
		last_name varchar(50),
		number bigint not null unique,
		balance bigint ,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		fmt.Println("account table creation error")
	}
	return err
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := `
		insert into account (first_name, last_name, number, created_at)
		 values ($1, $2, $3, $4) returning id, created_at
	`
	_, err := s.db.Exec(query, a.FirstName, a.LastName, a.Number, a.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {

	// var account = new(Account)

	query := `select * from account where id = $1`
	row := s.db.QueryRow(query, id)
	// if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.CreatedAt, &account.Balance); err != nil {
	// 	return nil, err
	// }
	account, err := scanIntoAccount(row)
	if err != nil {
		return nil, err
	}

	fmt.Println("account from db : ", account)

	return account, nil

}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `select * from account;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// var accounts []*Account

	// for rows.Next() {
	// 	account := new(Account) // creates a Account type variable and returns it pointer
	// 	if err := rows.Scan(
	// 		&account.ID,
	// 		&account.FirstName,
	// 		&account.LastName,
	// 		&account.Number,
	// 		&account.CreatedAt,
	// 		&account.Balance,
	// 	); err != nil {
	// 		return nil, err
	// 	}

	// 	accounts = append(accounts, account)

	// }

	accounts, err := scanIntoAccounts(rows)
	if err != nil {
		return nil, err
	}

	return accounts, nil

}

// scan through row
func scanIntoAccount(row *sql.Row) (*Account, error) {

	account := new(Account)
	if err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.CreatedAt,
		&account.Balance); err != nil {
		return nil, err
	}
	return account, nil
}

// scan through rows
func scanIntoAccounts(rows *sql.Rows) ([]*Account, error) {

	var accounts []*Account

	for rows.Next() {
		account := new(Account)
		// creates a Account type variable and returns it pointer
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.CreatedAt,
			&account.Balance,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)

	}
	return accounts, nil
}
