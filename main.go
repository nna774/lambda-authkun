package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var authOrigin = os.Getenv("AuthOrigin")

func generatePolicy(principalID, effect, resource string, ctx map[string]interface{}) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	authResponse.Context = ctx
	return authResponse
}

func deny(event events.APIGatewayCustomAuthorizerRequestTypeRequest, why string) (events.APIGatewayCustomAuthorizerResponse, error) {
	return generatePolicy("user", "Deny", event.MethodArn, map[string]interface{}{
		"response": `{"mes": "denied!"}`,
		"location": "https://" + event.Headers["host"],
	}), nil
}
func initiate(event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	host := event.Headers["Host"]
	backTo := url.QueryEscape("https://" + host + event.Path)
	callback := url.QueryEscape("https://" + host + "/_auth/callback")
	return generatePolicy("user", "Deny", event.MethodArn, map[string]interface{}{
		"msg":      `{"msg": "initate"}`,
		"location": authOrigin + "/auth?back_to=" + backTo + "&callback=" + callback,
	}), nil
}

func authkun(_ context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	req, err := http.NewRequest(http.MethodGet, authOrigin+"/test", nil)
	if err != nil {
		return deny(event, err.Error())
	}
	original := "https://" + event.Headers["Host"] + event.Path
	if len(event.PathParameters) > 0 {
		original += "?"
		for k, v := range event.PathParameters {
			original += k + "=" + v + "&"
		}
		original = original[:len(original)-1] // last &
	}
	req.Header.Add("x-ngx-omniauth-original-uri", original)
	req.Header.Add("cookie", event.Headers["cookie"])
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return deny(event, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return initiate(event)
	}

	return generatePolicy("user", "Allow", event.MethodArn, map[string]interface{}{
		"provider": resp.Header.Get("x-ngx-omniauth-provider"),
		"user":     resp.Header.Get("x-ngx-omniauth-user"),
		"info":     strings.Join(resp.Header.Values("x-ngx-omniauth-info"), ""),
	}), nil
}

func main() {
	lambda.Start(authkun)
}
