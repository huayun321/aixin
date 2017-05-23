package model

import "gopkg.in/mgo.v2/bson"

//Article 用户帖子
type Article struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author        User          `json:"author" bson:"author,omitempty"`
	Content       string        `json:"content" bson:"content,omitempty"`
	Images        []string      `json:"images,omitempty" bson:"images,omitempty"`
	Comments      []Comment     `json:"comments,omitempty" bson:"comments,omitempty"`
	Fans          []User        `json:"fans,omitempty" bson:"fans,omitempty"`
	IsSelected    bool          `json:"is_selected" bson:"is_selected,omitempty"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted,omitempty"`
	SelectTime    int64         `json:"select_time,omitempty" bson:"select_time,omitempty"`
	DeleteTime    int64         `json:"delete_time,omitempty" bson:"delete_time,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	LastReplyTime int64         `json:"last_reply_time,omitempty" bson:"last_reply_time,omitempty"`
}

//Comment 用户回复
type Comment struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author        User          `json:"author" bson:"author,omitempty"`
	Content       string        `json:"content" bson:"content,omitempty"`
	Reference     Comment       `json:"images,omitempty" bson:"images,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
}
