@echo off
REM 全开窗口版
start cmd /c go run %CD%/apps/app/start.go
start cmd /c go run %CD%/apps/user/user.go
start cmd /c go run %CD%/apps/video/video.go
start cmd /c go run %CD%/apps/social/start.go
start cmd /c go run %CD%/apps/interaction/interaction.go

pause