#!/bin/bash
# 移除环境变量
unset Mytool

current_path=$PATH
current_dir=$(pwd)
export PATH=$(echo $current_path | sed "s|:$current_dir||")

echo "Mytool removed"
echo "PATH updated: $PATH"