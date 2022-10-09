package service

import (
	"data-generator-mock/models"
	"data-generator-mock/pkg/logging"
	"data-generator-mock/pkg/repository"
	"errors"
	"math/rand"
	"time"
)

type DataGeneratorService struct {
	logger logging.Logger
	repo   repository.DataGenerator
}

func NewDataGeneratorService(logger logging.Logger, repo repository.DataGenerator) *DataGeneratorService {
	rand.Seed(time.Now().Unix())
	return &DataGeneratorService{
		logger: logger,
		repo:   repo,
	}
}

func (s *DataGeneratorService) GetByRv(rv int64) ([]models.Employee, error) {
	employees, err := s.repo.GetEmployees()
	if err != nil {
		return nil, err
	}

	if rv != 0 {
		for index, value := range employees {
			if value.Rv == rv {
				return employees[index+1:], err
			}
		}
	} else {
		return employees, nil
	}

	return nil, errors.New("rv does not exist")
}

func (s *DataGeneratorService) GenSlice(n int) error {
	data := make([]models.Employee, 0)
	for i := 0; i < n; i++ {
		data = append(data, genNewEmployee())
	}

	err := s.repo.SaveSlice(data)
	if err != nil {
		return err
	}

	return nil
}
func genNewEmployee() models.Employee {
	// Для красоты
	var maxValue int64 = 70000

	employee := models.Employee{
		//Если решил поменять, глянь сюда
		//https://dev.to/0xbf/sqlite-integer-can-save-up-to-signed-64-bit-with-golang-4i8a
		CardNumber: rand.Int63n(maxValue),
		EmployeeID: rand.Int63n(maxValue),
		Rv:         rand.Int63n(100 * maxValue),
		IsDeleted:  false,
	}

	return employee
}

func (s *DataGeneratorService) GenNewEmployee() error {
	//rand.Seed(time.Now().Unix())
	const MaxUint = ^uint(0)
	const MinUint = 0
	const MaxInt = int(MaxUint >> 1)
	const MinInt = -MaxInt - 1

	// Для красоты
	var maxValue int64 = 70000

	employee := models.Employee{
		//Если решил поменять, глянь сюда
		//https://dev.to/0xbf/sqlite-integer-can-save-up-to-signed-64-bit-with-golang-4i8a
		CardNumber: rand.Int63n(maxValue),
		EmployeeID: rand.Int63n(maxValue),
		Rv:         rand.Int63n(100 * maxValue),
		IsDeleted:  false,
	}

	err := s.repo.AddEmployee(employee)
	if err != nil {
		return err
	}

	return nil
}
