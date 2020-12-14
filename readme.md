# booking-server

## 開発日記

### 2020/12/13 

first draft.  

- token generation (sign with rsa512 private/public key)
- token payload encrypt
- don't use this in production
- use casbin support rbac
- token check on server side

**実装済み：**  
1. User authorization  
login by email and password and return access_token and refresh_token.

POST http://localhost:4000/api/v1/login  

```
curl -X POST "http://localhost:4000/api/v1/login" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"email\": \"booking\", \"password\": \"123456\"}"

Response:
{
  "code": 200,
  "error": "",
  "data": {
    "access_token": "eyJhbGciOiJSU...go_OhT3ti4BBLT1V7tM",
    "token_type": "bearer",
    "expires_in": 64800000,
    "refresh_token": "eyJhbGciOiJSU...eKUDzmip5ViA3sHBOo"
  }
}
```

2. Refresh token  
get new access token by using refresh token.

POST http://localhost:4000/api/v1/refresh  

```
curl -X POST "http://localhost:4000/api/v1/refresh" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"grant_type\": \"refresh_token\", \"refresh_token\": \"eyJhbGciOiJSU...eKUDzmip5ViA3sHBOo\"}"

Response:
{
  "code": 200,
  "error": "",
  "data": {
    "access_token": "eyJhbGciOiJSUz...BhEZ1nRLAKAW1UzVEC9aA",
    "token_type": "bearer",
    "expires_in": 64800000,
    "refresh_token": "eyJhbGciOiJS...PJ_m817XtmfZa2X0ROogA"
  }
}
```

## SSH Key Pair

- public.pem
- private.pem

```
1. Generating public/private rsa key pair.
ssh-keygen -m PEM -t rsa -b 4096 -C "booking@ditto"

2. Convert public key to PCKS8
ssh-keygen -f id_rsa.pub -e -M PKCS8 > public.pem

3. rename private key
cp -fp id_rsa private.pem
```

## Tools

generate go source code from mariadb table.

```
db2struct --host localhost -d bookingdb -t holidays --package models --struct Holiday -p --user booking --gorm
```
