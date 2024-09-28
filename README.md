# cputoy

print cpu usage in Linux with some vertical bars:

    |||||||||||
    |||||
    |||||||||||||||||||
    |||||||||||||||||||
    |||||||||||
    |||||||
    |||
    ||||||||||||||||||||

## build

Have glang enviroment installed

    go mod init cputoy
    ./build

## run

./cputoy

## exit

ctrl+c

## Limits

This will only run on Linux and similar enviroments which can provide proper /proc/stat data.  
Cputoy was mostly written by chatGPT with some minor help from /me.
