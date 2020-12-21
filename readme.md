# booking-server

## 開発日記

[実装済みAPI](api.md)
  
### 2020/12/21
  
**実装済み**  

Schedule  
- POST /schedule スケジュールを作成します
- DELETE /schedule/{id} スケジュールを削除します
- PUT /schedule/{id}/{status} スケジュールの利用可否（有効・無効）
  
Menu
- GET /menu/schedule/{id} スケジュール情報を取得します
- GET /menu/{id}/schedules スケジュール情報（複数）を取得します
- POST /menu メニューを作成
- PUT /menu/{id} メニュー情報を更新します
- DELETE /menu/{id} メニューを削除します 
- PUT /menu/{id}/{status} メニューの利用可否（有効・無効）

Tenant
- GET /tenant/menus メニュー情報を取得（複数）を取得します
- GET /tenant/menu/{id} メニュー情報を取得します

---
Facility
- POST /facility 施設を作成
- PUT /facility/{id} 施設情報を更新します
- DELETE /facility/{id} 施設を削除します 
- PUT /facility/{id}/{status} 施設の利用可否（有効・無効）
  
Tenant
- GET /tenant/facilities 施設情報を取得（複数）を取得します
- GET /tenant/facility/{id} 施設情報を取得します

---
### 2020/12/20
  
**実装済み**

Class  
- POST /class クラス作成
- PUT /class/user/{id} テナントに追加・削除
- POST /class/user/{id} ユーザーを作成します
- GET /class/users/{id} 指定クラスのユーザー一覧を取得
- GET /class/users/{id}/detail 指定クラスのユーザー一覧（詳細）を取得

User 
- GET /user/classes ユーザーのクラス一覧を取得します
- PUT /user/class/{id} ユーザークラスの切り替え
  
  
エラー時の戻り値やり直し、統一しました。

```
//Response -
type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
```
  
Tenant
- POST /tenant/user ユーザーを作成します（テナント）
- DELETE /tenant/user/{id} ユーザーを削除します（テナント）
- GET /tenant/users ユーザー一覧を取得します（ログイン中）
- GET /tenant/users/detail /tenant/users/detail

---
### 2020/12/18
  
**実装済み**

テナントの管理機能
Tenant
- POST /tenant テナントの仮作成を行います
- GET /tenant/users ユーザー一覧を取得します（ログイン中）

Admin/tenant
- GET /admin/tenant テナント情報を取得します(admin)
- POST /admin/tenant テナントを新規作成します(admin)
- PUT /admin/tenant/{id} テナント情報を更新します(admin)
- DELETE /admin/tenant/{id} テナント情報を削除します(admin)

User
- GET /user/tenants ユーザーの所属テナント一覧を取得します
- PUT /user/tenans/{id} ユーザーの所属テナントをアクティブします
- GET /user/logout ログアウト
<br>
<br>

---
### 2020/12/17
  
**実装済み**
  
仮登録完了メールの送信  

Dict 辞書
- GET /dict 辞書情報を取得します
- POST /dict 辞書情報を新規作成します
- POST /dict/array 辞書情報を新規作成します（複数）
  
Account  
- POST /admin/user 作成後メール送信
- DELETE /account/user アカウントを削除します

---
### 2020/12/16

**実装済み**

User  
- GET /user/confirm Email確認

/user/confirm/e=email&code=a69nrd62itj4
```
SQL:
　update 
    accounts a,
    accounts_confirm ac 
  set 
    a.email_confirmed=1,
    ac.used=1 
  where a.id=ac.account_id
        and ac.used=0
        and ac.confirm_code=? 
        and ac.email=?"
        and TIME_TO_SEC(timediff(now(),ac.update_date))<=?
  
  メールアドレスと確認コードが一致する
  かつ、指定期間内のみ更新できます。
  更新完了後、該当確認コードを使用済みにする。
```

**確認コードの生成**  

Twitterのsnowflake方式で唯一のコードを生成します。  
twitter's snowflake -> generate id -> Base36()  

```
github.com/bwmarrin/snowflake
```

---
### 2020/12/15

**実装済み**

Admin  

- GET /admin/account/{email} アカウント情報取得します
- GET /admin/user/{email} ユーザー情報を取得します(admin)
- PUT /admin/user/{email} ユーザー情報を更新します
  
User  
  
- GET /user ログイン中ユーザー（自分）情報を取得します
- PUT /user ユーザー（自分）情報を更新します
- GET /user/account ログイン中ユーザー（自分）ログイン情報を取得します
- POST /user/login ログイン
- PUT /user/password ユーザー（自分）パスワードを更新します
- POST /user/refresh リフレッシュトークンを使って新しいアクセストークンを取得します
  
---
### 2020/12/13 

first draft.  

- token generation (sign with rsa512 private/public key)
- token payload encrypt
- don't use this in production
- use casbin support rbac
- token check on server side

**実装済み**  
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
db2struct --host localhost -d bookingdb -t holidays --package models --struct Holiday -p --user booking --gorm --json
```
