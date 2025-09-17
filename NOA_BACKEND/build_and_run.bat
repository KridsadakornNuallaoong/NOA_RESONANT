@echo off
go build -o NOA.exe ./cmd/app
if %errorlevel% neq 0 exit /b %errorlevel%
NOA.exe