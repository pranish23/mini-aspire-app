package db

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
	"github.com/pranish23/mini-aspire-app/models"
)

type LoanDB struct {
	Conn *memdb.MemDB
}

func Init() (*LoanDB, error) {

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"loan": &memdb.TableSchema{
				Name: "loan",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "LoanID"},
					},
					"customer_id": &memdb.IndexSchema{
						Name:    "customer_id",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "CustomerID"},
					},
				},
			},
		},
	}

	// Create a new memdb database
	dbConn, err := memdb.NewMemDB(schema)
	if err != nil {
		fmt.Println("Failed to create memdb:", err)
		return nil, err
	}
	return &LoanDB{Conn: dbConn}, nil
}

func (ldb *LoanDB) Upsert(loan *models.Loan) error {
	txn := ldb.Conn.Txn(true)
	defer txn.Abort()
	if err := txn.Insert("loan", loan); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (ldb *LoanDB) ReadByID(col string, val string) ([]*models.Loan, error) {
	var loans []*models.Loan
	txn := ldb.Conn.Txn(false)
	defer txn.Abort()

	iter, err := txn.Get("loan", col, val)
	if err != nil {
		return nil, err
	}

	for obj := iter.Next(); obj != nil; obj = iter.Next() {
		l := obj.(*models.Loan)
		loans = append(loans, l)
	}
	if len(loans) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return loans, nil
}
