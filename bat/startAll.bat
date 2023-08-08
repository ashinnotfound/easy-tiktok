@echo off
REM 全开窗口版
cd ..
start powershell go run %CD%/apps/user/start.go
start powershell go run %CD%/apps/video/start.go
start powershell go run %CD%/apps/social/start.go
start powershell go run %CD%/apps/interaction/start.go
start  powershell go run %CD%/apps/app/start.go
pause

