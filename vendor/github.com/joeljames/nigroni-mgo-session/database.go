package nigronimgosession

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

type Database struct {
	dba DatabaseAccessor
}

func NewDatabase(databaseAccessor DatabaseAccessor) *Database {
	return &Database{databaseAccessor}
}

func (d *Database) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		reqSession := d.dba.Clone()
		defer reqSession.Close()
		ctx := d.dba.Set(request.Context(), request, reqSession)
		next(writer, request.WithContext(ctx))
	}
}
