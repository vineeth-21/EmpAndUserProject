package dto

type User struct {
	Emailid      string `json:"emailid,omitempty" bson:"emailid,omitempty"`
	Password     string `json:"password,omitempty" bson:"password,omitempty"`
	PasswordHash []byte `json:"passwordhash,omitempty" bson:"passwordhash,omitempty"`
	Role         string `json:"role,omitempty" bson:"role,omitempty"`
	Status       string `json:"status,omitempty" bson:"status,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Count        int64  `json:"count,omitempty" bson:"count,omitempty"`
}
