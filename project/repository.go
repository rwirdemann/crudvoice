package project

import "github.com/rwirdemann/crudvoice/domain"

type Repository struct {
	projects map[int]domain.Project
}

func NewRepository() *Repository {
	r := &Repository{projects: make(map[int]domain.Project)}
	r.Create(domain.Project{Id: 1, Name: "Picue"})
	return r
}

func (r *Repository) Create(project domain.Project) domain.Project {
	project.Id = r.nextId()
	r.projects[project.Id] = project
	return project
}

func (r *Repository) FindById(id int) (domain.Project, bool) {
	p, ok := r.projects[id]
	return p, ok
}

func (r *Repository) nextId() int {
	nextId := 1
	for _, v := range r.projects {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}
