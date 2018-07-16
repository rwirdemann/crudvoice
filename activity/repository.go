package activity

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	activities map[int]domain.Activity
}

func NewRepository() *Repository {
	r := Repository{activities: make(map[int]domain.Activity)}
	r.Create(domain.Activity{Name: "Programming", UserId: "f19f6ac3-f0a1-4e6d-bde3-a3d8be0f91fe"})
	r.Create(domain.Activity{Name: "Projektmanagement", UserId: "f19f6ac3-f0a1-4e6d-bde3-a3d8be0f91fe"})
	r.Create(domain.Activity{Name: "Testing", UserId: "f19f6ac3-f0a1-4e6d-bde3-a3d8be0f91fe"})
	return &r
}

func (r *Repository) Create(activity domain.Activity) domain.Activity {
	activity.Id = r.nextId()
	r.activities[activity.Id] = activity
	return activity
}

func (r *Repository) ById(id int) domain.Activity {
	return r.activities[id]
}

func (r *Repository) nextId() int {
	nextId := 1
	for _, v := range r.activities {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}
func (r *Repository) All() []domain.Activity {
	var activites []domain.Activity
	for _, a := range r.activities {
		activites = append(activites, a)
	}
	return activites
}
