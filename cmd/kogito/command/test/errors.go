package test

import (
	"github.com/kiegroup/kogito-cloud-operator/cmd/kogito/command/errors"
	"log"
)

type testErrorHandler struct {
}

func NewTestErrorHandler() errors.ErrorHandler {
	return &testErrorHandler{}
}

func (t *testErrorHandler) HandleError(err error) {
	log.Print(err.Error())
	log.Panic(err.Error())
}
