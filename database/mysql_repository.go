package database

import "github.com/rwirdemann/restvoice/domain"

type MySQLRepository struct {
}

func NewMySQLRepository() *MySQLRepository {
	r := MySQLRepository{}
	return &r
}

func (MySQLRepository) Invoices() []domain.Invoice {
	i1 := domain.Invoice{Id: 1, Customer: "Libri GmbH"}
	return []domain.Invoice{i1}
}
