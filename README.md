# Proxx Test
Proxx Test project

## Run 
Compiled executables may be run like this

on Linux (tested on Ubuntu 22.04):
```
    ./build/proxx
```

on Windows 64
```
    build\proxx.exe
```

on OSX (not tested)
```
    ./build/proxx_osx
```

## Gameplay

Game runs in the terminal, UI is textual.

Game asks questions about settings and then about turns.

The first question - whether you want to play Unicode variant -
this is good looking on Ubuntu terminal at least.
(screenshots made from Ubunut terminal with Unicode version)
For Windows - there were issues with Unicode symbols, it is better to play ASCII.

Anyway, best to have a fixed font in your terminal

screenshots/gameplay*.png - some gameplay samples.

## Build: Linux
Prerequisite: golang 1.19 should be installed

On Linux, to build for ALL platforms, run

```
    ./build.sh all
```

executables will be in the ./build folder

## Build: on other platforms
from the project root, do
```
    go build -o build/<name_your_executable> ./src
```

## Note on algorithms

The free space revealing algorithm appeared a bit challenging from the logic point,
as decision was to try to adapt an efficient non-recursive flood-fill algorithm for  
opening free areas with border of hole-adjacent cells (with counters),

so the main challenge was to provide the IsFillable() predicate classifying
free cells and bordering counter cells to the same "color" class (in terms of flood-fill)

For this purpose, isFirstCell parameter was added, to distinguish whether this is a first free cell
trying to be opened, as this appeared to be critical for rejecting non-connected regions
(see additional comments in the code) 

## References used
https://en.wikipedia.org/wiki/Flood_fill

https://lodev.org/cgtutor/floodfill.html


## Bug history screenshots
screenshots/fixed_bug*.png - some flood-fill/IsFillable() and index arithmetic bugs 
I had to mitigate during development 

## Copyright

Denys Volokhovskyi <rotaryden.gmail.com>