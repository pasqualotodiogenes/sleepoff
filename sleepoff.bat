@echo off
REM Launcher auxiliar para executar o app a partir da pasta atual
title sleepoff
cd /d "%~dp0"
sleepoff.exe %*
