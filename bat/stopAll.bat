@echo off
wmic process where "name='start.exe'" delete
pause