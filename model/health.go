package model

import (
	"database/sql"
	"fmt"
	"log"
)

type Service struct {
	ID     int
	Name   string        `json:"name"`
	Checks []HealthCheck `json:"health_checks"`
	Status int           `json:"status"` // 0 - Operational, 1 - Degradation, 2 - Downtime, 3 - Unknown
	Tags   []string      `json:"tags"`
	Feed   []Outage      `json:"feed"`
}

func (service Service) GetDescription() (string, []string) {
	return service.Name, service.Tags
}

func (service Service) GetShortService() Service {
	service.Checks = []HealthCheck{}
	return service
}

type ServiceDB struct {
	DB *sql.DB
}

func (svc ServiceDB) All() ([]Service, error) {
	var result []Service
	rows, err := svc.DB.Query("SELECT name, status FROM public.system")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var service Service
		err := rows.Scan(&service.Name, &service.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, service)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (svc ServiceDB) Create(service Service) (int, error) {
	sqlStatement := `
		INSERT INTO public.system(name, status, parentID)
		VALUES($1, $2, $3)
		RETURNING id`
	id := 0
	err := svc.DB.QueryRow(sqlStatement, service.Name, service.Status).Scan(&id)
	if err != nil {
		log.Println(err.Error())
		return id, err
	}
	return id, nil
}

func (svc ServiceDB) Get(title string) (Service, error) {
	var service Service
	sqlStatement := `SELECT * FROM public.system WHERE name=$1`
	err := svc.DB.QueryRow(sqlStatement, title).Scan(&service.ID, &service.Name, &service.Status)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = fmt.Errorf("Service with title %v doesn't exist", title)
		}
		log.Println(err.Error())
		return service, err
	}
	return service, nil
}

func (svc ServiceDB) Delete(title string) (bool, error) {
	sqlStatement := `DELETE FROM public.system WHERE name=$1`
	res, err := svc.DB.Exec(sqlStatement, title)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	rows, err := res.RowsAffected()
	if rows == 1 {
		return true, err
	}
	return false, err
}

func (svc ServiceDB) Update(title string, updatedSvc Service) (bool, error) {
	sqlStatement := `UPDATE public.system SET
		name = COALESCE($1, name),
		status = COALESCE($2, status),
		parentID = COALESCE($3, parentID)
		WHERE name=$4`
	res, err := svc.DB.Exec(sqlStatement, updatedSvc.Name, updatedSvc.Status, title)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	rows, err := res.RowsAffected()
	if rows == 1 {
		return true, err
	}
	return false, err
}

type HealthCheck struct {
	Status       int           `json:"status"` // 0 - OK, 1 - WARN, 2 - CRITICAL, 3 - UNKNOWN
	Title        string        `json:"title"`
	Details      string        `json:"details"`
	Priority     int           `json:"priority"`
	Relationship []HealthCheck `json:"rels"`
}

func (check *HealthCheck) SetStatus(status int) error {
	possibleValues := []int{0, 1, 2, 3}
	for _, value := range possibleValues {
		if status == value {
			check.Status = status
			return nil
		}
	}
	err := fmt.Errorf("Status %v is wrong, possible values are", possibleValues)
	return err
}

type Outage struct {
	Summary string `json:"summary"`
	Details string `json:"details"`
	Start   string `json:"scheduled_begin"`
	End     string `json:"schedule_end"`
}
