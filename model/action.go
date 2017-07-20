package model

import "gopkg.in/mgo.v2/bson"

//Action 动作
type Action struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	AuthorId   bson.ObjectId `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Name       string        `json:"name,omitempty" bson:"name,omitempty"`
	Level      string        `json:"level,omitempty" bson:"level,omitempty"`     //难度
	Symptom    string        `json:"symptom,omitempty" bson:"symptom,omitempty"` //适应症状
	People     string        `json:"people,omitempty" bson:"people,omitempty"`
	Notice     string        `json:"notice,omitempty" bson:"notice,omitempty"`
	MainImg    string        `json:"main_img,omitempty" bson:"main_img,omitempty"`
	StepImg    []string      `json:"step_img,omitempty" bson:"step_img,omitempty"`
	Key        string        `json:"key,omitempty" bson:"key,omitempty"` //关键点
	CreateTime int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
}
