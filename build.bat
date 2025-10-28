@echo off
setlocal enabledelayedexpansion

REM Cross-compilation script for Go binaries
REM Usage: build.bat [binary-name]

set BINARY_NAME=%1
if "%BINARY_NAME%"=="" set BINARY_NAME=qrl-ledger-recovery

echo Building %BINARY_NAME% for all platforms...
echo.

REM Create output directory
if not exist "builds" mkdir builds

REM Linux builds
echo Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o builds/%BINARY_NAME%-linux-amd64 .
if %ERRORLEVEL% NEQ 0 (echo Failed: Linux AMD64 & goto :error)

set GOARCH=arm64
go build -o builds/%BINARY_NAME%-linux-arm64 .
if %ERRORLEVEL% NEQ 0 (echo Failed: Linux ARM64 & goto :error)

set GOARCH=arm
set GOARM=7
go build -o builds/%BINARY_NAME%-linux-arm .
if %ERRORLEVEL% NEQ 0 (echo Failed: Linux ARM & goto :error)

REM Windows builds
echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o builds/%BINARY_NAME%-windows-amd64.exe .
if %ERRORLEVEL% NEQ 0 (echo Failed: Windows AMD64 & goto :error)

set GOARCH=arm64
go build -o builds/%BINARY_NAME%-windows-arm64.exe .
if %ERRORLEVEL% NEQ 0 (echo Failed: Windows ARM64 & goto :error)

REM macOS builds
echo Building for macOS...
set GOOS=darwin
set GOARCH=amd64
go build -o builds/%BINARY_NAME%-darwin-amd64 .
if %ERRORLEVEL% NEQ 0 (echo Failed: macOS AMD64 & goto :error)

set GOARCH=arm64
go build -o builds/%BINARY_NAME%-darwin-arm64 .
if %ERRORLEVEL% NEQ 0 (echo Failed: macOS ARM64 & goto :error)

REM FreeBSD builds
echo Building for FreeBSD...
set GOOS=freebsd
set GOARCH=amd64
go build -o builds/%BINARY_NAME%-freebsd-amd64 .
if %ERRORLEVEL% NEQ 0 (echo Failed: FreeBSD AMD64 & goto :error)

echo.
echo ========================================
echo All builds completed successfully!
echo Binaries are in the 'builds' directory
echo ========================================
goto :end

:error
echo.
echo ========================================
echo Build failed! Check the error above.
echo ========================================
exit /b 1

:end
endlocal