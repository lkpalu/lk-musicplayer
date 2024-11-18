@echo off
REM 添加当前目录到系统PATH环境变量
setx PATH "%PATH%;%~dp0" /M

REM 检查MyTool环境变量是否存在，如果不存在则创建
if not defined MyTool /M (
    setx MyTool "%~dp0" /M
) else (
    REM 添加当前目录到系统MyTool环境变量
    setx MyTool "%MyTool%;%~dp0" /M
)

echo Successfully installed.