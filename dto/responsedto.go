package dto

type JwtKey struct {
	Key string
}
type Response struct {
	Result    string
	Role      string
	Employee  *Employee
	Employees []Employee
	User      *User
	Count     int64
}
