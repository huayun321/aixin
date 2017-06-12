package model

import "gopkg.in/mgo.v2/bson"

//News 资讯
type News struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AuthorId      bson.ObjectId `json:"author_id" bson:"author_id,omitempty"`
	Title         string        `json:"title" bson:"title,omitempty"`
	Content       string        `json:"content" bson:"content,omitempty"`
	Image         string        `json:"image,omitempty" bson:"image,omitempty"`
	Position      int           `json:"position" bson:"position,omitempty"`
	IsPublished   bool          `json:"is_published" bson:"is_published,omitempty"`
	PublishTime   int64         `json:"publish_time,omitempty" bson:"publish_time,omitempty"`
	UnPublishTime int64         `json:"un_publish_time,omitempty" bson:"un_publish_time,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Author        User          `json:"author, omitempty" bson:",omitempty"`
	Comments      []Comment     `json:"comments"  bson:",omitempty"`
}
