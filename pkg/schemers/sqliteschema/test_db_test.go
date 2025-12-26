package sqliteschema

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func createTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 in-memory db: %w", err)
	}

	queries := []string{
		`CREATE TABLE Country (
			ID string PRIMARY KEY,
			Name string NOT NULL,
			NativeName string NOT NULL,
			Currency string NOT NULL,
			UNIQUE (Name),
			UNIQUE (NativeName)
		)`,
		`CREATE TABLE Currency (
			ID string PRIMARY KEY,
			Symbol string
		)`,
		`CREATE TABLE User (
			UserID int PRIMARY KEY,
			MainCountryID string,
			DefaultCurrency string,
			FOREIGN KEY (MainCountryID) REFERENCES Country(ID),
			FOREIGN KEY (DefaultCurrency) REFERENCES Currency(ID)
		)`,
		`CREATE TABLE Shop (
			ShopID int PRIMARY KEY,
			Name string NOT NULL
		)`,
		`CREATE TABLE [Order] (
			ShopID int,
			OrderID int,
			Currency string NOT NULL,
			CountryID string NOT NULL,
			Total decimal NOT NULL,
			PRIMARY KEY (ShopID, OrderID),
			FOREIGN KEY (Currency) REFERENCES Currency(ID),
			FOREIGN KEY (CountryID) REFERENCES Country(ID)
		)`,
		`CREATE TABLE Product (
			ProductID int PRIMARY KEY,
			Name string NOT NULL
		)`,
		`CREATE TABLE OrderDetails (
			ShopID int,
			OrderID int,
			ProductID int NOT NULL,
			Currency string NOT NULL,
			Price decimal NOT NULL,
			Quantity decimal NOT NULL,
			Total decimal NOT NULL,
			PRIMARY KEY (ShopID, OrderID, ProductID),
			FOREIGN KEY (ShopID, OrderID) REFERENCES [Order](ShopID, OrderID),
			FOREIGN KEY (ProductID) REFERENCES Product(ProductID),
			FOREIGN KEY (Currency) REFERENCES Currency(ID)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return nil, fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return db, nil
}
