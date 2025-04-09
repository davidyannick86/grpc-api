package models

type Teacher struct {
	Id        string `json:"id,omitempty" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Class     string `json:"class" bson:"class"`
	Subject   string `json:"subject" bson:"subject"`
}
