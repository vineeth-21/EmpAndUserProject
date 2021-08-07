package service

import (
	"errors"
	"log"
	"net/http"
	"test/db"
	"test/dto"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (t *RPCService) AddUser(r *http.Request, args *dto.User, result *dto.Response) error {
	userrole := r.Context().Value("userRole")
	log.Println(userrole)
	if userrole == "admin" {
		password := args.Password
		hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 6)
		if err != nil {
			return err
		}
		log.Println("hashpassword", hashpassword)
		var addUser dto.User
		addUser.Emailid = args.Emailid
		addUser.Password = ""
		addUser.PasswordHash = hashpassword
		addUser.Role = args.Role
		addUser.Status = args.Status
		addUser.Name = args.Name
		if args.Emailid != "" {
			err := db.InsertOne(addUser)
			if err != nil {
				return err
			}
			*result = dto.Response{Result: "added sucessfully"}
		} else {
			return errors.New("emailid is empty")
		}

	} else {
		return errors.New("admin can only add user")
	}
	return nil
}
func (t *RPCService) GetUser(r *http.Request, args *dto.User, result *[]dto.User) error {
	log.Println("entered into GetUser")
	var user []dto.User
	err := db.FindAll(&user, bson.M{}, bson.M{})
	if err != nil {
		return err
	}
	*result = user
	return nil

}
func (t *RPCService) GetUserByEmail(r *http.Request, args *dto.User, result *dto.User) error {
	log.Println("entered into GetUser By Email")
	var user dto.User
	err := db.Find(&user, bson.M{"emailid": args.Emailid})
	if err != nil {
		return err
	}
	*result = user
	return err
}

func (t *RPCService) DeleteUserByEmail(r *http.Request, args *dto.User, result *dto.Response) error {
	userrole := r.Context().Value("userRole")
	if userrole == "admin" {
		err := db.Delete(args, bson.M{"emailid": args.Emailid})
		if err != nil {
			return err

		}
		*result = dto.Response{Result: "DELETED SUCCESSFULLY"}
	} else {
		return errors.New("admin can only delete user")
	}

	return nil
}
func (t *RPCService) UpdateUser(r *http.Request, args *dto.User, result *dto.Response) error {
	userrole := r.Context().Value("userRole")
	if userrole == "admin" {
		err := db.Update(args, bson.M{"emailid": args.Emailid}, bson.M{"$set": args})
		if err != nil {
			return err
		}
		*result = dto.Response{Result: "updated successfully"}

	} else {
		return errors.New("admin only can update user")
	}
	return nil
}
