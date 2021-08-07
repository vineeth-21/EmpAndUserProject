package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"test/db"
	"test/dto"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type RPCService struct{}

func (t *RPCService) Login(r *http.Request, args *dto.Login, result *dto.JwtKey) error {
	log.Println("entered into login api....")
	email := args.Email
	password := args.Password
	log.Println("from login args password", password)
	log.Println("from login args", email, password)
	var user dto.User
	err := db.Find(&user, bson.M{"emailid": email})
	if err != nil {
		log.Println(err)
	}
	passwordEncryption := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if passwordEncryption != nil {
		log.Println("password incorrect")
		return passwordEncryption
	} else {
		log.Println("password correct")
		token, err := CreateToken(user.Role)
		if err != nil {
			return err
		}
		log.Println("jwtKey in login", token)
		*result = dto.JwtKey{Key: token}
	}
	return nil
}

var mySigningKey = []byte("secretMessage")

func CreateToken(userrole string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["userrole"] = userrole
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("entered into  middleware...")
		req, _ := ioutil.ReadAll(r.Body)
		log.Println(req)
		var request dto.RPCRequest
		if err := json.Unmarshal(req, &request); err != nil {
			log.Println(err)
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(req))
		log.Println(r.Body)
		log.Println(request)
		log.Println(request.Method)
		if request.Method == "rpcService.Login" {
			next.ServeHTTP(w, r)
		} else {
			key := r.Header.Get("Authorization")
			log.Println("key in middleware", key)
			credential := strings.Split(key, " ")
			result, err := jwt.Parse(credential[1], func(token *jwt.Token) (interface{}, error) {
				return mySigningKey, nil
			})
			log.Println(result)
			log.Println(result.Claims)
			claims := result.Claims.(jwt.MapClaims)
			role := claims["userrole"]
			log.Println(role)
			if err == nil && result.Valid {
				log.Println("Valid Token")
				ctx := context.WithValue(r.Context(), "userRole", role)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				log.Println("token invalid")
			}

		}
	})
}
