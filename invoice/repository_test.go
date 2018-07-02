package invoice

import (
	"testing"
	"github.com/rwirdemann/crudvoice/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvoice(t *testing.T) {
	r := NewRepository()
	i := domain.Invoice{Year: 2019, Month: 12}
	created := r.Create(i)
	reloaded, _ := r.FindById(created.Id)
	assert.Equal(t, 2019, reloaded.Year)
}
