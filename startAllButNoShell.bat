@echo off

start /B powershell go run %CD%/apps/app/start.go
start /B powershell go run %CD%/apps/user/start.go
start /B powershell go run %CD%/apps/video/start.go
start /B powershell go run %CD%/apps/social/start.go
start /B powershell go run %CD%/apps/interaction/start.go

pause