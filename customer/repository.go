package customer

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	customers map[int]domain.Customer
}

func NewRepository() *Repository {
	r := Repository{customers: make(map[int]domain.Customer)}
	c := domain.Customer{UserId: "f19f6ac3-f0a1-4e6d-bde3-a3d8be0f91fe"}
	r.Create(c)
	return &r
}

func (r *Repository) Create(customer domain.Customer) domain.Customer{
	customer.Id = r.nextId()
	r.customers[customer.Id] = customer
	return customer
}

func (r *Repository) ById(id int) domain.Customer {
	return r.customers[id]
}

func (r *Repository) ByUserId(userId string)[]domain.Customer {
	var customers []domain.Customer
	for _, c:= range r.customers {
		customers = append(customers, c)
	}
	return customers
}

func (r *Repository) nextId() int {
	nextId := 1
	for _, v := range r.customers {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

