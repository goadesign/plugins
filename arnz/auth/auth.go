package auth

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

const (
	header = "X-Amzn-Request-Context"
)

type Gate struct {
	MethodName        string
	AllowUnsigned     bool
	AllowArnsMatching []string
}

func IsUnsigned(r *http.Request) (pass bool) {
	return r.Header.Get(header) == "" || r.Header.Get(header) == "null"
}

func Authenticate(w http.ResponseWriter, r *http.Request) (caller *string, pass bool) {
	var amzCtx events.APIGatewayV2HTTPRequestContext
	amzReqCtxHeader := r.Header.Get(header)

	if IsUnsigned(r) {
		WriteUnauthenticated(w, "caller not authenticated")
		return
	}

	err := json.Unmarshal([]byte(amzReqCtxHeader), &amzCtx)
	if err != nil {
		WriteUnauthenticated(w, "failed to unmarshal X-Amzn-Request-Context header")
		return
	}

	if amzCtx.Authorizer == nil {
		WriteUnauthenticated(w, "no Authorizer defined in X-Amzn-Request-Context")
		return
	}

	if amzCtx.Authorizer.IAM == nil {
		WriteUnauthenticated(w, "no IAM defined in X-Amzn-Request-Context")
		return
	}

	if amzCtx.Authorizer.IAM.UserARN == "" {
		WriteUnauthenticated(w, "no UserARN defined in X-Amzn-Request-Context")
		return
	}

	return &amzCtx.Authorizer.IAM.UserARN, true
}

func Authorize(w http.ResponseWriter, callerArn string, matchers []string) (pass bool) {
	for _, pattern := range matchers {
		re := regexp.MustCompile(pattern)
		if re.MatchString(callerArn) {
			return true
		}
	}
	WriteUnauthorized(w, "caller not authorized")
	return false
}

func WriteUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unauthorized",
		"message": message,
	})
}

func WriteUnauthenticated(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unauthenticated",
		"message": message,
	})
}
