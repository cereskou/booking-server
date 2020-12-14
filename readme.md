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
    "access_token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc5MTU1MzIsIm5hbWUiOiJjODBiZTllZjliMmNlNmZkNTZhNTY2YmZkMjAwZmY4ZmZkNmQ1ZGQ4MjQ5OTcxY2Y4ZGMxZjFlNmQxMDE2MWRkYWI1MGZmMGVlZDM0OTllMTM1MGEzYyJ9.mrzLlzB0N92N_z2Z44nBsuOP7C2EuwgrxDFf9Czwbmw1SO43yr8AvcHHnyer9vADYJmOxaY4-zae7u1ZfIXRGxOSnWurZ-L7neK8Fb0rXgt0MIWR-xMSzJL1pal1IelPBjf5_gkTDbC2PWPqETgPj0BhN7iB9xLoZ1ig_hNepeCA2EU00d8l8IHlAcTlRXj0dJ_TV8O_mGptd5694SRm7B0oH7wLWDl8ok0Lh0CWprvm...EFRctSy5K3rxYBOT_88evszgBVUD4iyl-Nil-QtbzgA9HHNUZbbYAa9CiVrJykbCoqrHp5VEVU6gLwiW8hEoFyreU9B7jAGjLNIiYIimPT_bjcyRbhKJne7YcHdv8lPT1RXnEobZ3DKnLJPMQuq2TDvA0sYjn_1NedIKMYIDugFdM3CSQMH3SZ9obgo_OhT3ti4BBLT1V7tM",
    "token_type": "bearer",
    "expires_in": 64800000,
    "refresh_token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc5MzcxMzIsInN1YiI6IjlmOWEwOTg5YjlkZmQ5NmM0NDczZmY1NmNlMDFhYjg1Mzg2NDVlMGQ0NmFmZmI2YWQ0MjMxOGFjMmVkODRmYzY5ZDBlMTVlNjI4NDgwMWNkMmY2NTU0In0.Q7oLW6dx6blJPXXccd8Uf7d6I8i6iDpwftqHpb2-f3RegUyhZZhXsm7BajVaXmEYXADQmRm4liJOMaResliwwkJmN--59cBRjQVJ0pEM4Otljpoc8EBYzpAbjmuiQmiVSqx...dQ9HXaoYOKq_Cqfdyp0Y1AfoYjBthFMyymmvWKY5dDJ8k2vCoKTH-NGGZRwv51JQzVwuKwKHFjN1_UzkR3thFkJipeiuZ9e6LnzH8OXgqHuV-GlyA9LQNxE7x4ojwFi4OWDxAbSWCdaSP31j2NI2MEHXjiFXcUrbQjpH7Behcq1bpVsEf8h-VnwxqwXQVNZqLbXpYcqY267ATeKUDzmip5ViA3sHBOo"
  }
}
```

2. Refresh token  
get new access token by using refresh token.

POST http://localhost:4000/api/v1/refresh  

```
curl -X POST "http://localhost:4000/api/v1/refresh" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"grant_type\": \"refresh_token\", \"refresh_token\": \"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc5MzcxMzIsInN1YiI6IjlmOWEwOTg5YjlkZmQ5NmM0NDczZmY1NmNlMDFhYjg1Mzg2NDVlMGQ0NmFmZmI2YWQ0MjMxOGFjMmVkODRmYzY5ZDBlMTVlNjI4NDgwMWNkMmY2NTU0In0.Q7oLW6dx6blJPXXccd8Uf7d6I8i6iDpwftqHpb2-f3RegUyhZZhXsm7BajVaXmEYXADQmRm4liJOMaResliwwkJmN--59cBRjQVJ0pEM4Otljpoc8EBYzpAbjmuiQmiVSqx...dQ9HXaoYOKq_Cqfdyp0Y1AfoYjBthFMyymmvWKY5dDJ8k2vCoKTH-NGGZRwv51JQzVwuKwKHFjN1_UzkR3thFkJipeiuZ9e6LnzH8OXgqHuV-GlyA9LQNxE7x4ojwFi4OWDxAbSWCdaSP31j2NI2MEHXjiFXcUrbQjpH7Behcq1bpVsEf8h-VnwxqwXQVNZqLbXpYcqY267ATeKUDzmip5ViA3sHBOo\"}"

Response:
{
  "code": 200,
  "error": "",
  "data": {
    "access_token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc5MTU3OTIsIm5hbWUiOiI1ZmI2NzZjMDhhNmY4ZDMyMmVmN2U3Y2QxZjY3OGIzMDI5NzllMmYwYmJlNWE4MmJhNWVjNDVmMzY0MmU3NzI5MGZhMzU0M2NkMjZmODA4ODA1Yzc4MSJ9.RCStnIre2CtVCA4wuuq5xi_N4Wq5gv6TDCFmUquco-QHHjVc9vIXFVTtvZK1CKE4cC4T2-EDYq3eEBtLyqouTIVXZ-ysfW7nxtDAQvxQlvNB6-AdN6nIWjTIoiwl4O_xPGTDX...Nvzjn7qFCz2kk7SuVleafnPnwFKbJ0xgKb51zZuzNEvcxV6I64SKvYMLQs5JP53y5kyb7WEFgTJOQ4Tj32o6aR-e8NECtu9aOkZDgmpouzFikfNG6Mod_o2QNXAU4ftWFHWxNGCzRRvIyWx3HkKKu8wUMbJhfYjVgFv5NS_8M60RZgNO-EcLxI65ZWwafKhE5PaHeS97AJeeRIVgnWI7iUmnF3MBqNi0yRFCIMS-gyfgtvqs2LmP3DjcLYWppUGt4phVMWDylr9Heq8p3KqwEYQqry72WAsqBidci3-XbywBhEZ1nRLAKAW1UzVEC9aA",
    "token_type": "bearer",
    "expires_in": 64800000,
    "refresh_token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDc5MzczOTIsInN1YiI6ImQ3YjIyMjYxNDZkZDc0ZDEzZGJjOTk3YTA1MDNlMjJjZDEwOTAyM2MzNWViNjM4MjJiMDkwYTRkMGVmMGQ4ZTMwOTE4Y2Y0NWUzMzFiOWE5NzdmN2FkIn0.uM5jPlkQA480NL3K-3fKcm_AR...XJdNS3nd_LBiKoiA8JlOSXE0wIy5bcONvarNWBPKs3lY4-IyqgodBCFDUo1P8g7EZ-kri4ktExTZKwddP35XubTSoMiMEOEZoFxZP2TC6Dm0SRqQgDY9TktF8VmUiZanS2zgAamRJ4-WvsDvXhs2ZbBtKaz31ElcB7FfEMnRtYHM85b5Ps63axaRq6uly5-h1PJ_m817XtmfZa2X0ROogA"
  }
}
```

## Tools

generate go source code from mariadb table.

```
db2struct --host localhost -d bookingdb -t holidays --package models --struct Holiday -p --user booking --gorm
```
