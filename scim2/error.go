package scim2

import "errors"

var (
	ErrBadRequest = errors.New("request is unparsable or violates schema")
	ErrUnauthorized = errors.New("authorization failure as header is invalid or missing")
	ErrForbidden = errors.New("operation not permitted based on the supplied authorization")
	ErrNotFound = errors.New("specified resource or endpoint does not exist")
	ErrConflict = errors.New("the specified version number does match")
	ErrPreconditionFailed = errors.New("failed to update as resource has changed on the server")
	ErrPayloadTooLarge = errors.New("max operations or max payload size exceeded")
	ErrInternalError = errors.New("internal error")
	ErrNotImplemented = errors.New("service provider does not support the requested operation")
	ErrBadRequestInvalidFilter = errors.New("the specified filter syntax was invalid or not supported")
	ErrBadRequestTooMany = errors.New("the specified filter yields manz more results than the server is willing to process")
	ErrBadRequestUniqueness = errors.New("one or more of the attribute values are already in use or are reserverd")
	ErrBadRequestMutability = errors.New("the attempted modification is not compatible with the target attribute's mutability")
	ErrBadRequestInvalidSyntax = errors.New("the request body message structure was invalid")
	ErrBadRequestInvalidPath = errors.New("the path attribute was invalid or malformed")
	ErrBadRequestNoTarget = errors.New("the specified path did not yield an attribute that could be operated on")
	ErrBadRequestInvalidValue = errors.New("a required value was missing or not compatible with the operation")
	ErrBadRequestInvalidVers = errors.New("the specified SCIM protocol version is not supported")
	ErrBadRequestSensitive = errors.New("the specified request cannot be completed due to the passing of sensitive information in a request uri")
	)
