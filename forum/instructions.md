# Instructions

## Register a new user:

curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"Password123!"}'

## Login

curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Password123!"}' \
  -c cookies.txt

## Visit HomePage:

curl -X GET http://localhost:8080/ \
  -b cookies.txt

## Logout

curl -X POST http://localhost:8080/api/auth/logout \
  -b cookies.txt