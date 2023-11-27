package storage

import (
	"context"
	"database/sql"
	"fmt"
	"go-db-adapter/main/app/domain"
)

// MySQL represents a MySQL datastore.
type MySQL struct {
	Name string
	db   *sql.DB
}

func (m *MySQL) AdapterName() string {
	m.Name = "MySQL"
	return m.Name
}

func (m *MySQL) Connect(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return m.db, nil
}

func (m *MySQL) Close() error {
	return m.db.Close()
}

func (m *MySQL) FindAll() ([]domain.Person, error) {
	rows, err := m.db.Query("SELECT * FROM your_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Person
	for rows.Next() {
		var data domain.Person
		if err := rows.Scan(&data.ID, &data.Name, &data.Age); err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *MySQL) Insert(data domain.Person) error {
	_, err := m.db.Exec("INSERT INTO Persons (ID, Name, Age) VALUES (?, ?, /* ... */)",
		data.ID, data.Name, data.Age)
	return err
}

func (m *MySQL) GetSession() (any, error) {
	return nil, fmt.Errorf("GetSession is not applicable for MySQL")
}

func (m *MySQL) WithTransaction(ctx context.Context, fn func(sessionContext any) error) error {
	return fmt.Errorf("WithTransaction is not applicable for MySQL")
}

//=====+++=========================== ADAPTER CODE

// MySQLAdapter represents a MySQLAdapter datastore adapter.
type MySQLAdapter struct {
	mysql *MySQL
}

func (a *MySQLAdapter) AdapterName() string {
	return a.mysql.AdapterName()
}

func (a *MySQLAdapter) Connect(dataSourceName string) (*sql.DB, error) {
	return a.mysql.Connect(dataSourceName)
}

// Close closes the MySQL database connection.
func (a *MySQLAdapter) Close() error {
	return a.mysql.Close()
}

// FindAll retrieves all data from the MySQL table.
func (a *MySQLAdapter) FindAll() ([]domain.Person, error) {
	return a.mysql.FindAll()
}

// Insert inserts data into the MySQL table.
func (a *MySQLAdapter) Insert(data domain.Person) error {
	return a.mysql.Insert(data)
}

// GetSession is not applicable for MySQL, so it returns nil.
func (a *MySQLAdapter) GetSession() (any, error) {
	return a.mysql.GetSession()
}

// WithTransaction is not applicable for MySQL, so it returns an error.
func (a *MySQLAdapter) WithTransaction(ctx context.Context, fn func(sessionContext any) error) error {
	return a.mysql.WithTransaction(ctx, fn)
}
