package caller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Gate struct {
	MethodName        string
	AllowUnsigned     bool
	AllowArnsLike     []string
	AllowArnsMatching []string
}

const (
	header = "X-Amzn-Request-Context"
)

func Extract(w http.ResponseWriter, r *http.Request) (caller *string, pass bool) {
	var amzCtx events.APIGatewayV2HTTPRequestContext
	amzReqCtxHeader := r.Header.Get(header)

	err := json.Unmarshal([]byte(amzReqCtxHeader), &amzCtx)
	if err != nil {
		WriteUnauthorized(w, "failed to unmarshal X-Amzn-Request-Context header")
		return
	}

	if amzCtx.Authorizer == nil {
		WriteUnauthorized(w, "no Authorizer defined in X-Amzn-Request-Context")
		return
	}

	if amzCtx.Authorizer.IAM == nil {
		WriteUnauthorized(w, "no IAM defined in X-Amzn-Request-Context")
		return
	}

	if amzCtx.Authorizer.IAM.UserARN == "" {
		WriteUnauthorized(w, "no UserARN defined in X-Amzn-Request-Context")
		return
	}

	return &amzCtx.Authorizer.IAM.UserARN, true
}

func IsUnsigned(r *http.Request) (pass bool) {
	return r.Header.Get(header) == "" || r.Header.Get(header) == "null"
}

func ArnLike(w http.ResponseWriter, callerArn string, allowedArns []string) (pass bool) {
	for _, partialArn := range allowedArns {
		if strings.Contains(callerArn, partialArn) {
			return true
		}
	}
	WriteUnauthorized(w, fmt.Sprintf("caller %s is not authorized", callerArn))
	return false
}

func ArnMatch(w http.ResponseWriter, callerArn string, allowedArns []string) (pass bool) {
	for _, fullArn := range allowedArns {
		if callerArn == fullArn {
			return true
		}
	}
	WriteUnauthorized(w, fmt.Sprintf("caller %s is not authorized", callerArn))
	return false
}

func WriteUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unauthorized",
		"message": message,
	})
}
