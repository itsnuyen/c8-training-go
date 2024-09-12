## Simple Go C8 Example

### Start Locally or Against a Running SaaS Env

run app agains saas env (set env to `prod`:

```shell
env=prod go run cmd/main.go  
```

### Test Api and Start Processes

Api for Post request 
```shell
http://localhost:9999/api/v1/credit/start
```

Sample Body
```shell
{
  "customerId": "test",
  "orderTotal": 90,
  "productId": "3",
  "creditCard": {
    "cardNumber": "1234 5678",
    "cvc": "123",
    "expiryDate": "09/42"
  }
}

```

Example Post Request:

```shell
curl -X POST http://localhost:9999/api/v1/credit/start \
-H "Content-Type: application/json" \
-d '{
  "customerId": "test",
  "orderTotal": 90,
  "productId": "3",
  "creditCard": {
    "cardNumber": "1234 5678",
    "cvc": "123",
    "expiryDate": "09/42"
  }
}'
```