package model

import "gopkg.in/mgo.v2/bson"

//Article 用户帖子
type News struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AuthorId      bson.ObjectId `json:"author_id" bson:"author_id,omitempty"`
	Content       string        `json:"content" bson:"content,omitempty"`
	Image        string      `json:"image,omitempty" bson:"image,omitempty"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted,omitempty"`
	DeleteTime    int64         `json:"delete_time,omitempty" bson:"delete_time,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	BookmarkCount int           `json:"bookmark_count" bson:",omitempty"`
	Author        User          `json:"author" bson:",omitempty"`
	Comments      []Comment     `json:"comments"  bson:",omitempty"`
}