@echo off
setlocal
set ARG=%1

if "%ARG%"=="" (
    echo [Error] Please provide an argument: 0 = frontend only, 1 = backend only, 2 = full restart
    exit /b 1
)

echo Restart type: %ARG%

:: Restart Frontend Only
if "%ARG%"=="0" (
    echo [Frontend] Building frontend...
    pushd src\frontend
    call npm run build
    popd
    echo [Frontend] Build complete.
    echo [Frontend] Restarting Docker...
    docker-compose down
    docker-compose up --build -d
    echo [Frontend] Restart complete.
    goto end
)

:: Restart Backend Only
if "%ARG%"=="1" (
    echo [Backend] Tidy Go modules...
    pushd src\backend
    call go mod tidy
    popd

    echo [Backend] Restarting Docker...
    docker-compose down
    docker-compose up --build -d
    echo [Backend] Restart complete.
    goto end
)

:: Restart Both
if "%ARG%"=="2" (
    echo [Frontend] Building frontend...
    pushd src\frontend
    call npm run build
    popd

    echo [Backend] Tidy Go modules...
    pushd src\backend
    call go mod tidy
    popd

    echo [Docker] Restarting full stack...
    docker-compose down
    docker-compose up --build -d
    echo [System] Full restart complete.
    goto end
)

:: Invalid argument
echo [Error] Invalid argument: use 0, 1, or 2
exit /b 1

:end
exit /b 0
