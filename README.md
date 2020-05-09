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

さらに、`GatewayResponses`として、

```
ACCESS_DENIED:
  StatusCode: 302
  ResponseParameters:
    Headers:
      Location: context.authorizer.location
  ResponseTemplates:
    application/json: $context.authorizer.msg
```

を設定します。
また、API Gatewayの制約で`UNAUTHORIZED`の時にcontextを返したりできない？ 気がするので、悲しい気持ちになります。これが誤っていたら教えてください……。

そして認証を受ける側に

`import "github.com/nna774/lambda-authkun/adapter"` して

`http.HandleFunc("/_auth/callback", adapter.NewCallbackHandler("https://auth.dark-kuins.net/callback"))`

`/_auth/callback` を生やす。

## user info

omniauth_adapterの `/test` の返す以下の情報

- x-ngx-omniauth-provider
- x-ngx-omniauth-user
- x-ngx-omniauth-info

の3つをそれぞれ

- context.authorizer.provider
- context.authorizer.user
- context.authorizer.info

に入れて返します。
