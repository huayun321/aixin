package model

import "gopkg.in/mgo.v2/bson"

//News 资讯
type News struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AuthorId      bson.ObjectId `json:"author_id" bson:"author_id,omitempty"`
	Content       string        `json:"content" bson:"content,omitempty"`
	Image         string        `json:"image,omitempty" bson:"image,omitempty"`
	IsPublished   bool          `json:"is_published" bson:"is_published,omitempty"`
	PublishTime   int64         `json:"publish_time,omitempty" bson:"publish_time,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	BookmarkCount int           `json:"bookmark_count" bson:",omitempty"`
	Author        User          `json:"author" bson:",omitempty"`
	Comments      []Comment     `json:"comments"  bson:",omitempty"`
	Position      int           `json:"position"`
}
