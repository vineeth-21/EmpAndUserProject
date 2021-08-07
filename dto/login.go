package dto

type Login struct {
	Email    string `json:"emailid,omitempty" bson:"emailid,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}
