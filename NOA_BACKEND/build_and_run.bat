@echo off
go build -o NOA.exe main.go
if %errorlevel% neq 0 exit /b %errorlevel%
NOA.exe