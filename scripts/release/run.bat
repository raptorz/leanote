@echo off
setlocal

set "ROOT=%~dp0"
set "APP_LINK=%ROOT%runtime\github.com\pearlnote\pearlnote"

if not exist "%APP_LINK%" mklink /J "%APP_LINK%" "%ROOT%" >nul

"%ROOT%bin\pearlnote.exe" -importPath github.com/pearlnote/pearlnote -srcPath "%ROOT%runtime" -runMode prod %*
