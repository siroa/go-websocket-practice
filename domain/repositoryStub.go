package domain

import "websocket/utils/errs"

type RepositoryStub struct {
	hoge string
}

func (r RepositoryStub) saveMessage(msg Message) *errs.AppError {
	// ToDo: impl

	return nil
}

func NewRepositoryStub() RepositoryStub {

	return RepositoryStub{}
}
