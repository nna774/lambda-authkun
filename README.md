# lambda-authkun
lambdaでも[nginx_omniauth_adapter](https://github.com/sorah/nginx_omniauth_adapter)の認証使いたい！

## つかいかた

mainのlambdaを

```
Auth:
  Authorizers:
    Authkun:
      FunctionPayloadType: REQUEST
      FunctionArn: arn:aws:lambda:ap-northeast-1:486414336274:function:lambda-authkun
      Identity:
        Headers:
          - Cookie
        ReauthorizeEvery: 0
```

のようにAuthorizer type REQUESTにする。

そして認証を受ける側に

`http.HandleFunc("/_auth/callback", adapter.NewCallbackHandler("https://auth.dark-kuins.net/callback"))`
 
`/_auth/callback` を生やす。
