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
	IsPublished   bool          `json:"is_published" bson:"is_published"`
	PublishTime   int64         `json:"publish_time,omitempty" bson:"publish_time,omitempty"`
	UnPublishTime int64         `json:"un_publish_time,omitempty" bson:"un_publish_time,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Author        User          `json:"author, omitempty" bson:",omitempty"`
	NComments     []NComment    `json:"comments"  bson:",omitempty"`
}

//NComment 用户回复
type NComment struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	NewsID      bson.ObjectId `json:"news_id" bson:"news_id,omitempty"`
	AuthorID    bson.ObjectId `json:"author_id" bson:"author_id,omitempty"`
	Content     string        `json:"content" bson:"content,omitempty"`
	ReferenceID bson.ObjectId `json:"reference_id,omitempty" bson:"reference_id,omitempty"`
	CreateTime  int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Author      User          `json:"author" bson:"author,omitempty"`
	Comments    []Comment     `json:"comments" bson:",omitempty"`
}
