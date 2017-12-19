//************************************************************************//
// Code generated with aliaser, DO NOT EDIT.
//
// Aliased DSL Constants
//************************************************************************//

package design

import (
	design "goa.design/goa/http/design"
)

const (
	// FormatCIDR designates
	FormatCIDR = design.FormatCIDR
	// FormatDateTime designates values that follow RFC3339
	FormatDateTime = design.FormatDateTime
	// FormatUUID designates values that follow RFC4122
	FormatUUID = design.FormatUUID
	// FormatEmail designates values that follow RFC5322
	FormatEmail = design.FormatEmail
	// FormatHostname designates
	FormatHostname = design.FormatHostname
	// FormatIPv4 designates values that follow RFC2373 IPv4
	FormatIPv4 = design.FormatIPv4
	// FormatIPv6 designates values that follow RFC2373 IPv6
	FormatIPv6 = design.FormatIPv6
	// FormatIP designates values that follow RFC2373 IPv4 or IPv6
	FormatIP = design.FormatIP
	// FormatMAC designates
	FormatMAC = design.FormatMAC
	// FormatRegexp designates
	FormatRegexp = design.FormatRegexp
	// FormatURI designates
	FormatURI = design.FormatURI
	// FormatRFC1123 designates values that follow RFC1123
	FormatRFC1123 = design.FormatRFC1123
)

const (
	// DefaultView is the name of the default result type view.
	DefaultView = design.DefaultView
)

const (
	// BooleanKind represents a boolean.
	BooleanKind = design.BooleanKind
	// IntKind represents a signed integer.
	IntKind = design.IntKind
	// Int32Kind represents a signed 32-bit integer.
	Int32Kind = design.Int32Kind
	// Int64Kind represents a signed 64-bit integer.
	Int64Kind = design.Int64Kind
	// UIntKind represents an unsigned integer.
	UIntKind = design.UIntKind
	// UInt32Kind represents an unsigned 32-bit integer.
	UInt32Kind = design.UInt32Kind
	// UInt64Kind represents an unsigned 64-bit integer.
	UInt64Kind = design.UInt64Kind
	// Float32Kind represents a 32-bit floating number.
	Float32Kind = design.Float32Kind
	// Float64Kind represents a 64-bit floating number.
	Float64Kind = design.Float64Kind
	// StringKind represents a JSON string.
	StringKind = design.StringKind
	// BytesKind represent a series of bytes (binary data).
	BytesKind = design.BytesKind
	// ArrayKind represents a JSON array.
	ArrayKind = design.ArrayKind
	// ObjectKind represents a JSON object.
	ObjectKind = design.ObjectKind
	// MapKind represents a JSON object where the keys are not known in
	// advance.
	MapKind = design.MapKind
	// UserTypeKind represents a user type.
	UserTypeKind = design.UserTypeKind
	// ResultTypeKind represents a result type.
	ResultTypeKind = design.ResultTypeKind
	// AnyKind represents an unknown type.
	AnyKind = design.AnyKind
)

const (
	// Boolean is the type for a JSON boolean.
	Boolean = design.Boolean
	// Int is the type for a signed integer.
	Int = design.Int
	// Int32 is the type for a signed 32-bit integer.
	Int32 = design.Int32
	// Int64 is the type for a signed 64-bit integer.
	Int64 = design.Int64
	// UInt is the type for an unsigned integer.
	UInt = design.UInt
	// UInt32 is the type for an unsigned 32-bit integer.
	UInt32 = design.UInt32
	// UInt64 is the type for an unsigned 64-bit integer.
	UInt64 = design.UInt64
	// Float32 is the type for a 32-bit floating number.
	Float32 = design.Float32
	// Float64 is the type for a 64-bit floating number.
	Float64 = design.Float64
	// String is the type for a JSON string.
	String = design.String
	// Bytes is the type for binary data.
	Bytes = design.Bytes
	// Any is the type for an arbitrary JSON value (interface{} in Go).
	Any = design.Any
)

const (
	// StatusContinue refers to HTTP code 100 (RFC 7231, 6.2.1)
	StatusContinue = design.StatusContinue
	// StatusSwitchingProtocols refers to HTTP code 101 (RFC 7231, 6.2.2)
	StatusSwitchingProtocols = design.StatusSwitchingProtocols
	// StatusProcessing refers to HTTP code 102 (RFC 2518, 10.1)
	StatusProcessing = design.StatusProcessing
	// StatusOK refers to HTTP code 200 (RFC 7231, 6.3.1)
	StatusOK = design.StatusOK
	// StatusCreated refers to HTTP code 201 (RFC 7231, 6.3.2)
	StatusCreated = design.StatusCreated
	// StatusAccepted refers to HTTP code 202 (RFC 7231, 6.3.3)
	StatusAccepted = design.StatusAccepted
	// StatusNonAuthoritativeInfo refers to HTTP code 203 (RFC 7231, 6.3.4)
	StatusNonAuthoritativeInfo = design.StatusNonAuthoritativeInfo
	// StatusNoContent refers to HTTP code 204 (RFC 7231, 6.3.5)
	StatusNoContent = design.StatusNoContent
	// StatusResetContent refers to HTTP code 205 (RFC 7231, 6.3.6)
	StatusResetContent = design.StatusResetContent
	// StatusPartialContent refers to HTTP code 206 (RFC 7233, 4.1)
	StatusPartialContent = design.StatusPartialContent
	// StatusMultiStatus refers to HTTP code 207 (RFC 4918, 11.1)
	StatusMultiStatus = design.StatusMultiStatus
	// StatusAlreadyReported refers to HTTP code 208 (RFC 5842, 7.1)
	StatusAlreadyReported = design.StatusAlreadyReported
	// StatusIMUsed refers to HTTP code 226 (RFC 3229, 10.4.1)
	StatusIMUsed = design.StatusIMUsed
	// StatusMultipleChoices refers to HTTP code 300 (RFC 7231, 6.4.1)
	StatusMultipleChoices = design.StatusMultipleChoices
	// StatusMovedPermanently refers to HTTP code 301 (RFC 7231, 6.4.2)
	StatusMovedPermanently = design.StatusMovedPermanently
	// StatusFound refers to HTTP code 302 (RFC 7231, 6.4.3)
	StatusFound = design.StatusFound
	// StatusSeeOther refers to HTTP code 303 (RFC 7231, 6.4.4)
	StatusSeeOther = design.StatusSeeOther
	// StatusNotModified refers to HTTP code 304 (RFC 7232, 4.1)
	StatusNotModified = design.StatusNotModified
	// StatusUseProxy refers to HTTP code 305 (RFC 7231, 6.4.5)
	StatusUseProxy = design.StatusUseProxy
	// StatusTemporaryRedirect refers to HTTP code 307 (RFC 7231, 6.4.7)
	StatusTemporaryRedirect = design.StatusTemporaryRedirect
	// StatusPermanentRedirect refers to HTTP code 308 (RFC 7538, 3)
	StatusPermanentRedirect = design.StatusPermanentRedirect
	// StatusBadRequest refers to HTTP code 400 (RFC 7231, 6.5.1)
	StatusBadRequest = design.StatusBadRequest
	// StatusUnauthorized refers to HTTP code 401 (RFC 7235, 3.1)
	StatusUnauthorized = design.StatusUnauthorized
	// StatusPaymentRequired refers to HTTP code 402 (RFC 7231, 6.5.2)
	StatusPaymentRequired = design.StatusPaymentRequired
	// StatusForbidden refers to HTTP code 403 (RFC 7231, 6.5.3)
	StatusForbidden = design.StatusForbidden
	// StatusNotFound refers to HTTP code 404 (RFC 7231, 6.5.4)
	StatusNotFound = design.StatusNotFound
	// StatusMethodNotAllowed refers to HTTP code 405 (RFC 7231, 6.5.5)
	StatusMethodNotAllowed = design.StatusMethodNotAllowed
	// StatusNotAcceptable refers to HTTP code 406 (RFC 7231, 6.5.6)
	StatusNotAcceptable = design.StatusNotAcceptable
	// StatusProxyAuthRequired refers to HTTP code 407 (RFC 7235, 3.2)
	StatusProxyAuthRequired = design.StatusProxyAuthRequired
	// StatusRequestTimeout refers to HTTP code 408 (RFC 7231, 6.5.7)
	StatusRequestTimeout = design.StatusRequestTimeout
	// StatusConflict refers to HTTP code 409 (RFC 7231, 6.5.8)
	StatusConflict = design.StatusConflict
	// StatusGone refers to HTTP code 410 (RFC 7231, 6.5.9)
	StatusGone = design.StatusGone
	// StatusLengthRequired refers to HTTP code 411 (RFC 7231, 6.5.10)
	StatusLengthRequired = design.StatusLengthRequired
	// StatusPreconditionFailed refers to HTTP code 412 (RFC 7232, 4.2)
	StatusPreconditionFailed = design.StatusPreconditionFailed
	// StatusRequestEntityTooLarge refers to HTTP code 413 (RFC 7231, 6.5.11)
	StatusRequestEntityTooLarge = design.StatusRequestEntityTooLarge
	// StatusRequestURITooLong refers to HTTP code 414 (RFC 7231, 6.5.12)
	StatusRequestURITooLong = design.StatusRequestURITooLong
	// StatusUnsupportedResultType refers to HTTP code 415 (RFC 7231, 6.5.13)
	StatusUnsupportedResultType = design.StatusUnsupportedResultType
	// StatusRequestedRangeNotSatisfiable refers to HTTP code 416 (RFC 7233, 4.4)
	StatusRequestedRangeNotSatisfiable = design.StatusRequestedRangeNotSatisfiable
	// StatusExpectationFailed refers to HTTP code 417 (RFC 7231, 6.5.14)
	StatusExpectationFailed = design.StatusExpectationFailed
	// StatusTeapot refers to HTTP code 418 (RFC 7168, 2.3.3)
	StatusTeapot = design.StatusTeapot
	// StatusUnprocessableEntity refers to HTTP code 422 (RFC 4918, 11.2)
	StatusUnprocessableEntity = design.StatusUnprocessableEntity
	// StatusLocked refers to HTTP code 423 (RFC 4918, 11.3)
	StatusLocked = design.StatusLocked
	// StatusFailedDependency refers to HTTP code 424 (RFC 4918, 11.4)
	StatusFailedDependency = design.StatusFailedDependency
	// StatusUpgradeRequired refers to HTTP code 426 (RFC 7231, 6.5.15)
	StatusUpgradeRequired = design.StatusUpgradeRequired
	// StatusPreconditionRequired refers to HTTP code 428 (RFC 6585, 3)
	StatusPreconditionRequired = design.StatusPreconditionRequired
	// StatusTooManyRequests refers to HTTP code 429 (RFC 6585, 4)
	StatusTooManyRequests = design.StatusTooManyRequests
	// StatusRequestHeaderFieldsTooLarge refers to HTTP code 431 (RFC 6585, 5)
	StatusRequestHeaderFieldsTooLarge = design.StatusRequestHeaderFieldsTooLarge
	// StatusUnavailableForLegalReasons refers to HTTP code 451 (RFC 7725, 3)
	StatusUnavailableForLegalReasons = design.StatusUnavailableForLegalReasons
	// StatusInternalServerError refers to HTTP code 500 (RFC 7231, 6.6.1)
	StatusInternalServerError = design.StatusInternalServerError
	// StatusNotImplemented refers to HTTP code 501 (RFC 7231, 6.6.2)
	StatusNotImplemented = design.StatusNotImplemented
	// StatusBadGateway refers to HTTP code 502 (RFC 7231, 6.6.3)
	StatusBadGateway = design.StatusBadGateway
	// StatusServiceUnavailable refers to HTTP code 503 (RFC 7231, 6.6.4)
	StatusServiceUnavailable = design.StatusServiceUnavailable
	// StatusGatewayTimeout refers to HTTP code 504 (RFC 7231, 6.6.5)
	StatusGatewayTimeout = design.StatusGatewayTimeout
	// StatusHTTPVersionNotSupported refers to HTTP code 505 (RFC 7231, 6.6.6)
	StatusHTTPVersionNotSupported = design.StatusHTTPVersionNotSupported
	// StatusVariantAlsoNegotiates refers to HTTP code 506 (RFC 2295, 8.1)
	StatusVariantAlsoNegotiates = design.StatusVariantAlsoNegotiates
	// StatusInsufficientStorage refers to HTTP code 507 (RFC 4918, 11.5)
	StatusInsufficientStorage = design.StatusInsufficientStorage
	// StatusLoopDetected refers to HTTP code 508 (RFC 5842, 7.2)
	StatusLoopDetected = design.StatusLoopDetected
	// StatusNotExtended refers to HTTP code 510 (RFC 2774, 7)
	StatusNotExtended = design.StatusNotExtended
	// StatusNetworkAuthenticationRequired refers to HTTP code 511 (RFC 6585, 6)
	StatusNetworkAuthenticationRequired = design.StatusNetworkAuthenticationRequired
)
