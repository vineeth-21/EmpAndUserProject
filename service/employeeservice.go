package service

import (
	"errors"
	"log"
	"net/http"
	"test/db"
	"test/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *RPCService) AddEmployee(r *http.Request, args *dto.Employee, result *dto.Response) error {
	userrole := r.Context().Value("userRole")
	doj, err := time.Parse("02-01-2006", args.Doj)
	if err != nil {
		return err

	}
	log.Println(doj)
	if userrole == "user" {
		if args.EmpID != "" {
			var emp dto.Employee
			emp.EmpID = args.EmpID
			emp.Doj = ""
			emp.Designation = args.Designation
			emp.Experience = args.Experience
			emp.ProjectName = args.ProjectName
			emp.Salary = args.Salary
			emp.DateOfJoining = doj
			err := db.InsertOne(&emp)
			if err != nil {
				return err
			}
			*result = dto.Response{Result: "employee added successfuly"}
		} else {
			return errors.New("emailid should not be empty")
		}
	} else {
		return errors.New("only user can add employee")
	}
	return nil
}
func (t *RPCService) GetEmployee(r *http.Request, args *dto.Employee, result *[]dto.Employee) error {
	var employee []dto.Employee
	err := db.FindAll(&employee, bson.M{}, bson.M{})
	if err != nil {
		return err
	}
	*result = employee
	return nil
}
func (t *RPCService) EmployeePagination(r *http.Request, args *dto.Employee, result *dto.Response) error {
	log.Println("entered into employee pagination")
	var employee []dto.Employee
	cur, err := db.FindAllPagination(&employee, bson.M{}, int64(args.Page), int64(args.Size), bson.M{})
	if err != nil {
		return err
	}
	*result = dto.Response{Employees: employee, Count: cur}
	return nil
}
func (t *RPCService) GetEmployeeById(r *http.Request, args *dto.Employee, result *dto.Employee) error {
	userrole := r.Context().Value("userRole")
	if userrole == "user" {
		var emp dto.Employee
		err := db.Find(&emp, bson.M{"empid": args.EmpID})
		if err != nil {
			log.Println(err)
		}
		*result = emp
	} else {
		return errors.New("ONLY USER CAN GET THE EMPLOYEE DETAILS")
	}
	return nil
}
func (t *RPCService) DeleteEmpById(r *http.Request, args *dto.Employee, result *dto.Response) error {
	userrole := r.Context().Value("userRole")
	if userrole == "user" {
		log.Println(args.EmpID)
		err := db.Delete(args, bson.M{"empid": args.EmpID})
		if err != nil {
			return err
		}
		*result = dto.Response{Result: "employee deleted successfully"}
	} else {
		return errors.New("USER ONLY CAN DELETE EMPLOYEE")
	}
	return nil
}
func (t *RPCService) UpdateEmployee(r *http.Request, args *dto.Employee, result *dto.Response) error {
	log.Println("entered into update employee")
	userrole := r.Context().Value("userRole")
	if userrole == "user" {
		err := db.Update(args, bson.M{"empid": args.EmpID}, bson.M{"$set": args})
		if err != nil {
			return nil
		}
		*result = dto.Response{Result: "UPDATED SUCCESSFULLY"}
	} else {
		return errors.New("only user can update employee")
	}
	return nil

}
