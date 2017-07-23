package model

import "gopkg.in/mgo.v2/bson"

//Plan 态度
type Plan struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	AuthorId    bson.ObjectId `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	First       int           `json:"first,omitempty" bson:"first,omitempty"`
	Second      int           `json:"second,omitempty" bson:"second,omitempty"`
	F2          int           `json:"f2,omitempty" bson:"f2,omitempty"`
	F3          int           `json:"f3,omitempty" bson:"f3,omitempty"`
	Level       string        `json:"level,omitempty" bson:"level,omitempty"`
	Feel        string        `json:"feel,omitempty" bson:"feel,omitempty"`
	Weeks       []Week        `json:"week,omitempty" bson:"week,omitempty"`
	CreateTime  int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	IsRecommend bool          `json:"is_recommend,omitempty" bson:"is_recommend,omitempty"`
}

type Week struct {
	ID   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Days []Day         `json:"days,omitempty" bson:"days,omitempty"`
}

type Day struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	DayActions []DayAction   `json:"day_actions,omitempty" bson:"day_actions,omitempty"`
	Goal       string        `json:"goal,omitempty" bson:"goal,omitempty"`
}

type DayAction struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Action     Action        `json:"action,omitempty" bson:"action,omitempty"`
	GroupTime  int           `json:"group_time,omitempty" bson:"group_time,omitempty"`
	Time       int           `json:"time,omitempty" bson:"time,omitempty"`
	NotifyTime int           `json:"notify_time,omitempty" bson:"notify_time,omitempty"`
}
