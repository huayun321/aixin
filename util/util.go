package util

import "github.com/unrolled/render"

//Ren render
var Ren *render.Render

func init() {
	Ren = render.New(render.Options{IndentJSON: true, StreamingJSON: true, IsDevelopment: true})
}
