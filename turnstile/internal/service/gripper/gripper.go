package gripper

import (
	"turnstile/internal/service"
	"turnstile/internal/service/infrastructure"
)

type DataGripper struct {
	client     *infrastructure.DataClient
	empRepo    service.EmployeeRepo
	rvFileName string
}

func New(client *infrastructure.DataClient, empRepo service.EmployeeRepo, rvFileName string) *DataGripper {
	return &DataGripper{
		client:     client,
		empRepo:    empRepo,
		rvFileName: rvFileName,
	}
}

func (dg *DataGripper) LoadData() error {
	// Прочитали rv
	rv, err := readRv(dg.rvFileName)
	if err != nil {
		return err
	}

	// Получили данные
	employees, err := dg.client.GetData(rv)
	if err != nil {
		return err
	}

	// Вставляем данные в бд
	//for _, value := range employees["data"] {
	//	err = dg.empRepo.Save(value)
	//	if err != nil {
	//		return err
	//	}
	//}

	// Вставляем данные в бд по-нормальному
	err = dg.empRepo.SaveSlice(employees["data"])
	if err != nil {
		return err
	}

	// Записали последнее rv в файл
	emp := employees["data"]
	if len(emp) > 0 {
		err = writeRv(dg.rvFileName, emp[len(emp)-1].Rv)
		if err != nil {
			return err
		}
	}

	return nil
}
