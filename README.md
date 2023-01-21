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

The free space revealing algorithm appeared a bit challenging from the logic point,
as decision was to try to adapt efficient non-recursive flood-fill algorithm for  
opening free areas with border of hole-adjacent cells (with counters),

so the main challenge was to provide the IsFillable() predicate classifying
free cells and bordering counter cells to the same "color" class (in terms of flood-fill)

## References used
https://en.wikipedia.org/wiki/Flood_fill

https://lodev.org/cgtutor/floodfill.html


## Bonus: bug screens
screenshots/fixed_bug*.png - some flood-fill/IsFillable() bugs I had to mitigate during development 

## Copyright

Denys Volokhovskyi <rotaryden.gmail.com>