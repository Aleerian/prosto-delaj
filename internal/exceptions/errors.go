package exceptions

import (
	"errors"
	"fmt"
)

const (
	Success ErrorCode = iota
	Create
	ResourceCreated
	NoContent
	BadRequest
	Unauthorized
	Forbidden
	BadHeader
	Conflict
	ResourceAlreadyExist
	InternalServerError
	ResourceDeleted
	UnregisteredAccount
	InvalidAccountStatus
	ConfirmationTimeout
	ConfirmCodeNotMatched
	AccountMustConfirmed
	NotFound
	ResourceInTrash
)

const (
	errorServer = "server error"
)

type ErrorCode int

var EmptyResponse interface{}

func (e ErrorCode) String() string {
	return [...]string{
		"success",
		"create",
		"resourceCreated",
		"noContent",
		"badRequest",
		"unauthorized",
		"forbidden",
		"badHeader",
		"conflict",
		"resourceAlreadyExist",
		"internalServerError",
		"resourceDeleted",
		"accountIsNotRegistered",
		"invalidAccountStatus",
		"confirmationCodeIsTimeOut",
		"confirmationCodeIsNotMatched",
		"accountMustBeConfirmedBefore",
		"notFound",
		"resourceInTrash",
	}[e]
}

func (e ErrorCode) Message() interface{} {
	return [...]interface{}{
		"success",
		"create",
		"resourceDeleted",
		EmptyResponse,
		"badRequest",
		"unauthorized",
		"forbidden",
		"badHeader",
		"conflict",
		"resourceAlreadyExist",
		"internalServerError",
		"resourceDeleted",
		"accountIsNotRegistered",
		"invalidAccountStatus",
		"confirmationCodeIsTimeOut",
		"confirmationCodeIsNotMatched",
		"accountMustBeConfirmedBefore",
		"notFound",
		"resourceInTrash",
	}[e]
}

func (e ErrorCode) CustomMessage(text string) interface{} {
	baseMessage := [...]interface{}{
		"success",
		"create",
		"resourceDeleted",
		EmptyResponse,
		"badRequest",
		"unauthorized",
		"forbidden",
		"badHeader",
		"conflict",
		"resourceAlreadyExist",
		"internalServerError",
		"resourceDeleted",
		"accountIsNotRegistered",
		"invalidAccountStatus",
		"confirmationCodeIsTimeOut",
		"confirmationCodeIsNotMatched",
		"accountMustBeConfirmedBefore",
		"notFound",
		"resourceInTrash",
	}[e]

	return fmt.Sprintf(`%s: %s`, baseMessage, text)
}

var ErrorServer = errors.New(errorServer)
