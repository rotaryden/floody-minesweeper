# Floody MinesWeeper
![Floody MinesWeeper](/images/title.png)

Minesweeper implementation

The motivation is to make a custom Go implementation of the well-known Minesweeper game,
using non-recursive flood-fill variation of the free area opening algorithm,
and number of data structures with generics.

Best playable in the Linux terminal with proper unicode support.

## Run from release 
Download latest binary for your OS from 
[Releases](https://github.com/rotaryden/floody-minesweeper/releases)
and run in a terminal with good unicode support (for unicode ui)

## Build and run: Linux
Prerequisite: latest golang should be installed
tested on Ubuntu 22.04

to build and run
```
    ./build.sh all
    ./build/floody-minesweeper
```

to build for ALL 3 platforms
```
    ./build.sh all
```
executables will be in the ./build folder


## Build and run: other platforms
from the project root, do
```
    go build -o build/<name_your_executable> ./src
    build/<name_your_executable>
```

## Gameplay

Game runs in the terminal, textual UI.

Game asks questions about UI settings and then about turns.

The first question - whether you want to play Unicode variant -
this is good looking on Ubuntu terminal at least.
(screenshots made from Ubunut terminal with Unicode version)

For Windows - there were issues with Unicode, it is better to play ASCII.

Best to have a fixed font in your terminal

screenshots/gameplay*.png - some gameplay samples.

## Note on algorithms

The free space revealing algorithm appeared a bit challenging from the logic point,
as decision was to try to adapt an efficient non-recursive flood-fill algorithm for  
opening free areas with border of mine-adjacent cells (with counters),

so the main challenge was to provide the IsFillable() predicate classifying
free cells and bordering counter cells to the same "color" class (in terms of flood-fill)

For this purpose, isFirstCell parameter was added, to distinguish whether this is a first free cell
trying to be opened, as this appeared to be critical for rejecting non-connected regions
(see additional comments in the code) 

## version2 fixes
- boundary fix for cell names , e.g. "b2" should be in scope of [a..width] [0...height]

## References used
https://en.wikipedia.org/wiki/Flood_fill

https://lodev.org/cgtutor/floodfill.html


## Bugs
screenshots/fixed_bug*.png - some flood-fill/IsFillable() and index arithmetic bugs 

Known bug 1 (believed to be fixed):  see screenshot known_bug__line_scan_not_reaching_left.png - 
    when line scan flood-fill algorithm is on the stage when searching for left-most end of the free line to be fileld, 
    proper isFillable predicate should be provided, standard one is not working, as cells to the right on this line stay closed,
    and predicate has a protextion from double opened areas. 

Possible bug 2: opening dis-joint free area via a diagonal free neighbour. Needs more investigations.

## Copyright

Denys Volokhovskyi <rotaryden.gmail.com>

MIT License
