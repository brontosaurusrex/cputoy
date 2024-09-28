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

Have golang enviroment installed, then:

    go mod init cputoy
    ./build

might just work.

## run

./cputoy

## exit

ctrl+c

## Limits

This will only run on Linux and similar enviroments which can provide proper /proc/stat data.  
Cputoy was mostly written by chatGPT with some minor help from /me.
