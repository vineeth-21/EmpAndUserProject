package db

import (
	"context"
	"errors"
	"log"
	"reflect"

	dto "test/dto"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

var Client *mongo.Client
var C map[string]*mongo.Collection
var (
	NoPtr = errors.New("you must pass in a pointer")
)
var database string

func Connect() error {

	database = viper.GetString("MongoDBName")
	C = make(map[string]*mongo.Collection)
	Ctx := context.Background()
	defer Ctx.Done()
	var err error
	Client, err = mongo.Connect(Ctx, options.Client().ApplyURI(viper.GetString("MongoUrl")))
	if err != nil {
		return err
	}
	err = Client.Ping(Ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}
func Init() {
	coll := Col(&dto.User{})

	CreateUniqueIndex(coll, "emailid", nil)
	cols := Col(&dto.Employee{})
	CreateUniqueIndex(cols, "empid", nil)
}
func CreateUniqueIndex(col *mongo.Collection, field string, partial interface{}) error {
	log.Println("createUniqueIndex called")
	opts := options.Index().SetUnique(true)
	if partial != nil {
		opts.SetPartialFilterExpression(partial)
	}
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: opts,
	}
	ctx := context.Background()
	defer ctx.Done()
	_, err := col.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Printf("error in creating index:[&s] : %v\n", field, err)
	}
	return err
}

func CreateDefaultUser() int64 {
	log.Println("default user called....")
	var userexists dto.User
	userCount, err := Count(&userexists, bson.M{})
	if err != nil {
		log.Println(err)
	}
	log.Println("number of users in usercollection is :", userCount)
	password, _ := bcrypt.GenerateFromPassword([]byte("112233"), 5)
	var defaultUser dto.User
	defaultUser.Emailid = "defaultuser@gmail.com"
	defaultUser.Password = ""
	defaultUser.PasswordHash = password
	defaultUser.Role = "admin"
	defaultUser.Status = "active"
	defaultUser.Name = "default user"
	if userCount == 0 {
		errs := InsertOne(defaultUser)
		if errs != nil {
			log.Println(err)
		}
	}
	return userCount
}
func InsertOne(e interface{}) error {
	ctx := context.Background()
	defer ctx.Done()
	cl := Col(e)
	_, err := cl.InsertOne(ctx, e)
	if err != nil {
		return err
	}
	return nil
}
func Find(res interface{}, q bson.M) error {
	if !IsPtr(res) {
		return NoPtr
	}
	coll := Col(res)
	ctx := context.Background()
	defer ctx.Done()
	err := coll.FindOne(ctx, q).Decode(res)
	return err
}
func FindAll(res interface{}, q bson.M, sort bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	if !IsPtr(res) {
		return NoPtr
	}
	opts := options.Find()
	opts.SetSort(sort)
	coll := Col(res)
	cursor, err := coll.Find(ctx, q, opts)
	if err != nil {
		return err
	}
	cursor.All(ctx, res)
	return err

}
func FindAllPagination(res interface{}, q bson.M, page int64, size int64, sort bson.M) (int64, error) {
	ctx := context.Background()
	defer ctx.Done()
	if !IsPtr(res) {
		return 0, NoPtr
	}
	findOptions := options.Find()
	var skip int64 = 0
	if page > 1 {
		skip = size * (page - 1)
	}
	log.Println(skip)
	log.Println("pagesize in pagination", page)
	findOptions.SetLimit(size)
	findOptions.SetSkip(skip)
	findOptions.SetSort(sort)
	coll := Col(res)
	cursor, err := coll.Find(context.Background(), q, findOptions)
	if err != nil {
		return 0, err
	}
	var count int64
	count, err = coll.CountDocuments(ctx, q)
	if err != nil {
		return 0, err
	}
	err = cursor.All(ctx, res)
	if err != nil {

		return 0, err
	}
	return count, err
}
func Delete(res interface{}, q bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	if !IsPtr(res) {
		return NoPtr
	}
	if q == nil || len(q) == 0 {
		return errors.New("qury annot be nil")
	}
	coll := Col(res)
	_, err := coll.DeleteOne(ctx, q)
	return err
}
func Update(res interface{}, q bson.M, set bson.M) error {
	ctx := context.Background()
	defer ctx.Done()
	coll := Col(res)
	err := coll.FindOneAndUpdate(ctx, q, set).Decode(res)
	return err

}
func Count(res interface{}, q bson.M) (int64, error) {
	if !IsPtr(res) {
		return 0, NoPtr
	}
	coll := Col(res)
	ctx := context.Background()
	defer ctx.Done()
	c, err := coll.CountDocuments(ctx, q)
	if err != nil {
		log.Println(err)
	}
	return c, err
}
func Col(e interface{}) *mongo.Collection {
	cname := typeName(e)
	res, ok := C[cname]
	if !ok {
		db := Client.Database(database)
		log.Printf("Type:%s", cname)
		r2 := db.Collection(cname)

		res = r2
		C[cname] = r2

	}
	return res

}
func typeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	if isSlice(t) {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	return t.Name()
}
func isSlice(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Slice
}
func IsPtr(i interface{}) bool {
	return reflect.ValueOf(i).Kind() == reflect.Ptr
}
