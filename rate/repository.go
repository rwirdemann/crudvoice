package rate

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	rates map[int]map[int]domain.Rate
}

func NewRepository() *Repository {
	r := Repository{rates: make(map[int]map[int]domain.Rate)}

	// Programmierung
	r.Create(domain.Rate{ActivityId: 1, ProjectId: 1, Price: 60})
	r.Create(domain.Rate{ActivityId: 1, ProjectId: 1, Price: 60})

	// Projektmanagement
	r.Create(domain.Rate{ActivityId: 2, ProjectId: 1, Price: 50})
	r.Create(domain.Rate{ActivityId: 2, ProjectId: 2, Price: 50})

	// QA
	r.Create(domain.Rate{ActivityId: 3, ProjectId: 1, Price: 55})
	r.Create(domain.Rate{ActivityId: 3, ProjectId: 2, Price: 55})

	return &r
}

func (r *Repository) Create(rate domain.Rate) {
	if projectRates, ok := r.rates[rate.ProjectId]; ok {
		projectRates[rate.ActivityId] = rate
	} else {
		r.rates[rate.ProjectId] = make(map[int]domain.Rate)
		r.rates[rate.ProjectId][rate.ActivityId] = rate
	}
}

func (r *Repository) ByProjectIdAndActivityId(projectId int, activityId int) domain.Rate {
	return r.rates[projectId][activityId]
}
