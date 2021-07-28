echo off
setlocal

rem obtain base dir which is the parent dir of bin
pushd..
SET BASE_DIR=%cd%
popd

rem add "bin" to windows PATH
IF EXIST %BASE_DIR%\bin SET PATH=%PATH%;%BASE_DIR%\bin

rem set DIR0=%BASE_DIR%\service
set PROTO_DIR=%BASE_DIR%\%1
set OUTPUT_DIR=%PROTO_DIR%\autogen

echo input directory: %PROTO_DIR%
echo output directory: %OUTPUT_DIR%

mkdir %OUTPUT_DIR%
for /F %%i in ('dir /b %PROTO_DIR%\*.proto') do (
    protoc -I %PROTO_DIR% --gogofaster_out=plugins=grpc:%OUTPUT_DIR% %PROTO_DIR%\%%i
)
