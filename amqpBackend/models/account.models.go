package models

type CollAccount struct {
	ID            string   `json:"_id,omitempty" bson:"_id"`
	Username      string   `json:"Username,omitempty" bson:"Username"`
	Email         string   `json:"Email,omitempty" bson:"Email"`
	Password      string   `json:"Password,omitempty" bson:"Password"`
	IsVerifyEmail bool     `json:"isVerifyEmail,omitempty" bson:"isVerifyEmail"`
	Profile       []string `json:"Profile,omitempty" bson:"Profile"`
	CourseId      []int    `json:"CourseId,omitempty" bson:"CourseId"`
}
