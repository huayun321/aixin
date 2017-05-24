package model

import "gopkg.in/mgo.v2/bson"

//Article 用户帖子
type Article struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AuthorId      bson.ObjectId `json:"author_id" bson:"author_id,omitempty"`
	Content       string        `json:"content" bson:"content,omitempty"`
	Images        []string      `json:"images,omitempty" bson:"images,omitempty"`
	IsSelected    bool          `json:"is_selected" bson:"is_selected,omitempty"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted,omitempty"`
	SelectTime    int64         `json:"select_time,omitempty" bson:"select_time,omitempty"`
	DeleteTime    int64         `json:"delete_time,omitempty" bson:"delete_time,omitempty"`
	CreateTime    int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	FansCount     int64         `json:"fans_count,omitempty" bson:",omitempty"`
	CommentsCount int64         `json:"comments_count,omitempty" bson:",omitempty"`
	ViewCount     int64         `json:"view_count,omitempty" bson:",omitempty"`
}

//Comment 用户回复
type Comment struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author      User          `json:"author" bson:"author,omitempty"`
	Content     string        `json:"content" bson:"content,omitempty"`
	ReferenceID bson.ObjectId `json:"reference_id,omitempty" bson:"reference_id,omitempty"`
	CreateTime  int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
}

//Fan 喜欢的人
type Fan struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ArticleID bson.ObjectId `json:"article_id" bson:"article_id,omitempty"`
	UserID    bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`
}

//Bookmark 收藏的人
type Bookmark struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ArticleID bson.ObjectId `json:"article_id" bson:"article_id,omitempty"`
	UserID    bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`
}
