# Proxx Test

## Run 
(Compiled executables bundled)

on Linux (tested on Ubuntu 22.04):
```
    ./build/proxx
```

on Windows 64
```
    ./build/proxx.exe
```

on OSX (not tested)
```
    ./build/proxx_osx
```


## Gameplay

Game runs in the terminal, UI is textual using unicode sumbols,
best to have a fixed font on your terminal

Just run for Linux:
```
    ./build/proxx
```
or another executable for you platform


## Build: Linux
Prerequisite: golang 1.19 installed

On Linux, to build for all platforms, run

```
    ./build.sh
```

executables will be in the ./build folder

## Build: on other platforms

do
```
    go build -o build/<name_your_executable> ./src
```

## Algorithms

Free space revealing algorithm appeared a bit challenging from the logic side,
as desision was to take relatively efficient non-recursive flood-fill algorithm,
and then adapt it to opening free areas with hole-adjacent cells (with counters),
so the main challenge was to provide IsFillable() predicate comprising
free cells and counter cells from the border to the same "color" class (in terms of flood-fill)

## References used
https://en.wikipedia.org/wiki/Flood_fill

https://lodev.org/cgtutor/floodfill.html


## Bonus: bug screens
fixed_bugs/*.png - some flood-fill/IsFillable() bugs I had to mitigate during development 

## Copyright

Denys Volokhovskyi <rotaryden.gmail.com>