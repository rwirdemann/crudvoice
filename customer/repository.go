package customer

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	cutomers map[int]domain.Customer
}


