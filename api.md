## 実装済みAPI

| class | api path | method | memo |
| ---- | ---- | ---- | ---- |
| user | /user | GET | ログイン中ユーザー（自分）情報を取得します |
| | /user | PUT | ユーザー（自分）情報を更新します<br> {"name":"ユーザー","address":"住所"} |
| | /user/account | GET | ログイン中ユーザー（自分）ログイン情報を取得します |
| | /user/class/{id} | PUT | ユーザークラスを切り替えします<br>複数クラスに所属の場合、アクティブのクラスが１つしかありません。 |
| | /user/classes | GET | ユーザーのクラス一覧取得を取得します |
| | /user/confirm | GET | ユーザーEmailの確認 |
| | /user/login | POST | ログイン |
| | /user/logout | GET | ログアウト |
| | /user/password | PUT | パスワードを変更します |
| | /user/refresh | リフレッシュトークン<br>新しいアクセストークンを発行します |
| | /user/tenants | GET | ユーザーの所属テナント一覧を取得します |
| | /user/tenants/{id} | PUT | ユーザーの所属テナントをアクティブします<br>複数テナントに所属の場合、アクティブのテナントが１つしかありません。 |


  
- User  
![User](images/api_user.png)

- Tenant 
![Tenant](images/api_tenant.png)

- Holidays
![Holiday](images/api_holidays.png)

- Dict
![Dict](images/api_dict.png)

- Admin
![Admin](images/api_admin.png)