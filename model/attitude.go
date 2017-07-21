package model

import "gopkg.in/mgo.v2/bson"

//Attitude 态度
type Attitude struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	AuthorId   bson.ObjectId `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Name       string        `json:"name,omitempty" bson:"name,omitempty"`
	Desc       string        `json:"desc,omitempty" bson:"symptom,omitempty"`
	MainImg    string        `json:"main_img,omitempty" bson:"main_img,omitempty"`
	IconImg    string        `json:"icon_img,omitempty" bson:"icon_img,omitempty"`
	CreateTime int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
}
