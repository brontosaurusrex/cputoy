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

Have golang (=>1.22) enviroment installed, then:

    go mod init cputoy
    ./build

might just work.

## run

./cputoy

## fake cpu stress

    # apt install stress
    stress -i 12 -c 3
    # should give you some moving bars

## exit

ctrl+c

## Limits

This will only run on Linux and similar enviroments which can provide proper /proc/stat data.  
Cputoy was mostly written by chatGPT with some minor help from /me.
