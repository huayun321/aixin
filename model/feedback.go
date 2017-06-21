package model

import "gopkg.in/mgo.v2/bson"

type Feedback struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AuthorId   bson.ObjectId `json:"author_id" bson:"author_id,omitempty"`
	Content    string        `json:"content" bson:"content,omitempty"`
	Phone      string        `json:"phone" bson:"phone,omitempty"`
	IsTracked  bool          `json:"is_tracked" bson:"is_tracked"`
	IsReplied  bool          `json:"is_replied" bson:"is_replied"`
	Reply      Reply         `json:"reply" bson:"reply,omitempty"`
	TrackTime  int64         `json:"track_time,omitempty" bson:"track_time,omitempty"`
	CreateTime int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Author     User          `json:"author,omitempty" bson:",omitempty"`
}

type Reply struct {
	Author  User   `json:"author,omitempty" bson:",omitempty"`
	Content string `json:"content" bson:"content,omitempty"`
	CreateTime int64         `json:"create_time,omitempty" bson:"create_time,omitempty"`

}
