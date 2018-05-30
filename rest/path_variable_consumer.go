package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

type PathVariableConsumer struct {
	pathVarName string
}

func NewPathVariableConsumer(pathVarName string) PathVariableConsumer {
	return PathVariableConsumer{pathVarName: pathVarName}
}

func (c PathVariableConsumer) Consume(request interface{}) interface{} {
	vars := mux.Vars(request.(*http.Request))
	return vars[c.pathVarName]
}
