package customer

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	customers map[int]domain.Customer
}

func NewRepository() *Repository {
	r := Repository{customers: make(map[int]domain.Customer)}
	c := domain.Customer{Name: "3skills", UserId: "f19f6ac3-f0a1-4e6d-bde3-a3d8be0f91fe"}
	c = r.Create(c)

	project := domain.Project{Id: 1, Name: "Instantfoo.com", CustomerId: c.Id}
	c.Projects = []domain.Project{project}
	r.Update(c)

	return &r
}

func (r *Repository) Create(customer domain.Customer) domain.Customer {
	customer.Id = r.nextId()
	r.customers[customer.Id] = customer
	return customer
}

func (r *Repository) Update(customer domain.Customer) domain.Customer {
	r.customers[customer.Id] = customer
	return customer
}

func (r *Repository) ById(id int) domain.Customer {
	return r.customers[id]
}

func (r *Repository) All() []domain.Customer {
	var customers []domain.Customer
	for _, c := range r.customers {
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
