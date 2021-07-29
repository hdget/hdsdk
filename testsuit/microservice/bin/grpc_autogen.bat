@echo off
setlocal enabledelayedexpansion

@REM obtain base dir which is the parent dir of bin
pushd..
SET BASE_DIR=%cd%
SET PROTO_DIR_NAME=pb
popd

@REM add "bin" to windows PATH
IF EXIST %BASE_DIR%\bin SET PATH=%PATH%;%BASE_DIR%\bin

@REM set root dir
set ROOT_DIR=%BASE_DIR%\%1
if [%1] == [] (
    set ROOT_DIR=%BASE_DIR%
)
set OUTPUT_DIR=%ROOT_DIR%\%PROTO_DIR_NAME%

mkdir %OUTPUT_DIR%
for /r %ROOT_DIR% %%i in (*.proto) do (
    set absdir=%%i
    for /f %%j in ('dir /b %%i') do (
        @REM replace the filename to be empty in abs dir string
        @REM please use ! in for loop and enabledelayedexpansion
        set PROTO_DIR=!absdir:%%j=!
        if exist !PROTO_DIR! (
            echo proto directory: !PROTO_DIR!
            echo output directory: !OUTPUT_DIR!
            protoc -I !PROTO_DIR! --gogofaster_out=plugins=grpc:!OUTPUT_DIR! %%i
            goto quit
        )
    )
)
:quit
