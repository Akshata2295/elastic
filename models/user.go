package models

type User struct {
	ID   int `json:"id", autoincrement`
	Name string `json:"name"`
	Age  int    `json:"age"`
}