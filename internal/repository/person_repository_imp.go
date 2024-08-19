package repository

import "fileserver/internal/model"

type PersonRepositoryImp struct{}

var _ PersonRepository = (*PersonRepositoryImp)(nil)

func (p *PersonRepositoryImp) FindAll() ([]*model.Person, error) {
	persons := make([]*model.Person, 0)
	persons = append(persons, &model.Person{
		ID:   1,
		Name: "John",
		Age:  23,
	})
	persons = append(persons, &model.Person{
		ID:   2,
		Name: "Jane",
		Age:  24,
	})
	return persons, nil
}
