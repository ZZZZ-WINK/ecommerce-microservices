@echo off
echo Generating proto files...

:: Check if protoc is installed
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: protoc not found. Please install protoc first.
    exit /b 1
)

:: Check if go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: go not found. Please install Go first.
    exit /b 1
)

:: Install required go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

:: Parse command line arguments
set SERVICE=
set CLEAN=

:parse_args
if "%~1"=="" goto :run
if /i "%~1"=="-service" (
    set SERVICE=-service=%~2
    shift
    shift
    goto :parse_args
)
if /i "%~1"=="-clean" (
    set CLEAN=-clean
    shift
    goto :parse_args
)
shift
goto :parse_args

:run
:: Run generate script
go run generate.go %SERVICE% %CLEAN%

echo Generation completed! 