@echo off
setlocal enabledelayedexpansion

set "PROJECT_NAME=bub"
set "VERSION=0.1.0"
set "BUILD_DIR=build"
set "COMPRESS=%COMPRESS%"

if not "%~1"=="" set "PROJECT_NAME=%~1"
if not "%~2"=="" set "VERSION=%~2"

echo.
echo ================================================
echo    Go Multi-Platform Build System
echo    Building: %PROJECT_NAME% v%VERSION%
echo ================================================
echo.

where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Go is not installed or not in PATH!
    pause
    exit /b 1
)

echo [SUCCESS] Go detected
go version
echo.

if not exist go.mod (
    echo [ERROR] go.mod not found!
    echo Please run this script from your Go project root.
    pause
    exit /b 1
)

if exist %BUILD_DIR% (
    echo [WARNING] Cleaning existing build directory...
    rmdir /s /q %BUILD_DIR% 2>nul
)
mkdir %BUILD_DIR% 2>nul
echo [SUCCESS] Created build directory
echo.

echo [INFO] Downloading dependencies...
go mod download
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to download dependencies
    pause
    exit /b 1
)
echo [SUCCESS] Dependencies downloaded
echo.

set /a SUCCESS_COUNT=0
set /a FAIL_COUNT=0
set /a TOTAL_COUNT=15

echo [INFO] Building for %TOTAL_COUNT% platforms...
echo.
echo ------------------------------------------------

call :build_target "windows" "amd64" ".exe"
call :build_target "windows" "arm64" ".exe"
call :build_target "linux" "amd64" ""
call :build_target "linux" "arm" ""
call :build_target "linux" "arm64" ""
call :build_target "darwin" "amd64" ""
call :build_target "darwin" "arm64" ""
call :build_target "android" "arm64" ""

echo ------------------------------------------------
echo.
echo ================================================
echo Build Summary:
echo ================================================
echo [SUCCESS] Successful: %SUCCESS_COUNT% / %TOTAL_COUNT%
if %FAIL_COUNT% gtr 0 (
    echo [FAILED] Failed: %FAIL_COUNT% / %TOTAL_COUNT%
)

set "TOTAL_SIZE=0"
for %%F in (%BUILD_DIR%\*) do (
    set /a TOTAL_SIZE+=%%~zF
)
set /a TOTAL_SIZE_MB=%TOTAL_SIZE% / 1048576
echo [INFO] Total size: %TOTAL_SIZE_MB% MB
echo [INFO] Output directory: %BUILD_DIR%
echo.

echo [INFO] Built binaries:
for %%F in (%BUILD_DIR%\*) do (
    set "FILE_SIZE=%%~zF"
    set /a FILE_SIZE_MB=!FILE_SIZE! / 1048576
    if !FILE_SIZE_MB! equ 0 (
        set /a FILE_SIZE_KB=!FILE_SIZE! / 1024
        echo   * %%~nxF ^(!FILE_SIZE_KB! KB^)
    ) else (
        echo   * %%~nxF ^(!FILE_SIZE_MB! MB^)
    )
)
echo.

if %SUCCESS_COUNT% equ %TOTAL_COUNT% (
    echo [SUCCESS] ALL BUILDS COMPLETED SUCCESSFULLY!
) else (
    if %SUCCESS_COUNT% gtr 0 (
        echo [WARNING] Some builds failed, but %SUCCESS_COUNT% succeeded
    ) else (
        echo [ERROR] ALL BUILDS FAILED!
        pause
        exit /b 1
    )
)

echo.
set /p CREATE_RELEASE="Create release archive? (y/n): "
if /i "%CREATE_RELEASE%"=="y" (
    set "RELEASE_NAME=%PROJECT_NAME%-v%VERSION%-all-platforms.zip"
    echo [INFO] Creating release archive: !RELEASE_NAME!
    
    where powershell >nul 2>nul
    if !ERRORLEVEL! equ 0 (
        powershell -command "Compress-Archive -Path '%BUILD_DIR%\*' -DestinationPath '!RELEASE_NAME!' -Force"
        if !ERRORLEVEL! equ 0 (
            for %%F in (!RELEASE_NAME!) do set "RELEASE_SIZE=%%~zF"
            set /a RELEASE_SIZE_MB=!RELEASE_SIZE! / 1048576
            echo [SUCCESS] Release archive created: !RELEASE_NAME! ^(!RELEASE_SIZE_MB! MB^)
        ) else (
            echo [ERROR] Failed to create release archive
        )
    ) else (
        echo [WARNING] PowerShell not available for compression
    )
)

echo.
echo [SUCCESS] BUILD PROCESS COMPLETE!
echo.
pause
exit /b 0

:build_target
set "GOOS=%~1"
set "GOARCH=%~2"
set "EXT=%~3"
set "OUTPUT_NAME=%PROJECT_NAME%-%GOOS%-%GOARCH%-v%VERSION%%EXT%"
set "OUTPUT_PATH=%BUILD_DIR%\%OUTPUT_NAME%"

echo   [BUILD] %GOOS%/%GOARCH% -^> %OUTPUT_NAME%

set CGO_ENABLED=0
go build -ldflags="-s -w -X main.Version=%VERSION%" -o "%OUTPUT_PATH%" 2>nul

if %ERRORLEVEL% equ 0 (
    if exist "%OUTPUT_PATH%" (
        for %%F in ("%OUTPUT_PATH%") do set "FILE_SIZE=%%~zF"
        set /a FILE_SIZE_MB=!FILE_SIZE! / 1048576
        if !FILE_SIZE_MB! equ 0 (
            set /a FILE_SIZE_KB=!FILE_SIZE! / 1024
            echo   [OK] Built successfully ^(!FILE_SIZE_KB! KB^)
        ) else (
            echo   [OK] Built successfully ^(!FILE_SIZE_MB! MB^)
        )
        
        if /i "%COMPRESS%"=="true" (
            where powershell >nul 2>nul
            if !ERRORLEVEL! equ 0 (
                set "ZIP_NAME=%OUTPUT_NAME%.zip"
                set "ZIP_PATH=%BUILD_DIR%\!ZIP_NAME!"
                powershell -command "Compress-Archive -Path '%OUTPUT_PATH%' -DestinationPath '!ZIP_PATH!' -Force" 2>nul
                if !ERRORLEVEL! equ 0 (
                    for %%F in ("!ZIP_PATH!") do set "ZIP_SIZE=%%~zF"
                    set /a ZIP_SIZE_MB=!ZIP_SIZE! / 1048576
                    if !ZIP_SIZE_MB! equ 0 (
                        set /a ZIP_SIZE_KB=!ZIP_SIZE! / 1024
                        echo        Compressed: !ZIP_NAME! ^(!ZIP_SIZE_KB! KB^)
                    ) else (
                        echo        Compressed: !ZIP_NAME! ^(!ZIP_SIZE_MB! MB^)
                    )
                )
            )
        )
        
        set /a SUCCESS_COUNT+=1
    ) else (
        echo   [FAIL] Build failed
        set /a FAIL_COUNT+=1
    )
) else (
    echo   [FAIL] Build failed
    set /a FAIL_COUNT+=1
)
echo.
goto :eof