package model

import "gopkg.in/mgo.v2/bson"

//Plan 态度
type Plan struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	AuthorId    bson.ObjectId `json:"author_id,omitempty" bson:"author_id,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	First       Part           `json:"first,omitempty" bson:"first,omitempty"`
	Second      Part           `json:"second,omitempty" bson:"second,omitempty"`
	F2          Part           `json:"f2,omitempty" bson:"f2,omitempty"`
	F3          Part           `json:"f3,omitempty" bson:"f3,omitempty"`
	Level       string        `json:"level,omitempty" bson:"level,omitempty"`
	Feel        string        `json:"feel,omitempty" bson:"feel,omitempty"`
	Weeks       []Week        `json:"weeks,omitempty" bson:"weeks,omitempty"`
	CreateTime  int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	IsRecommend bool          `json:"is_recommend,omitempty" bson:"is_recommend,omitempty"`
	Desc        string        `json:"desc,omitempty" bson:"desc,omitempty"`
	Img        string        `json:"img,omitempty" bson:"img,omitempty"`
}

type Week struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Days     []Day         `json:"days,omitempty" bson:"days,omitempty"`
	Goal     Goal          `json:"goal,omitempty" bson:"goal,omitempty"`
	Attitudes []Attitude      `json:"attitudes,omitempty" bson:"attitudes,omitempty"`
}

type Goal struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Content string        `json:"content,omitempty" bson:"content,omitempty"`
	Mode    string        `json:"mode,omitempty" bson:"mode,omitempty"`
}

type Day struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	DayActions []DayAction   `json:"day_actions,omitempty" bson:"day_actions,omitempty"`
}

type DayAction struct {
	ID         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Action     Action        `json:"action,omitempty" bson:"action,omitempty"`
	Group      int           `json:"group,omitempty" bson:"group,omitempty"`
	Time       int           `json:"time,omitempty" bson:"time,omitempty"`
	NotifyTime int           `json:"notify_time,omitempty" bson:"notify_time,omitempty"`
}

type Part struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	Text string        `json:"text,omitempty" bson:"text,omitempty"`
}
