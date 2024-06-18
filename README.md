# Capyclick
## A simple cross platform Capybara clicker game

## Features

- Leveling system
- Mandarin rain event
- 3 types of capybaras
- Audio level control (keyboard only)
- Responsive to window size change rendering
- Mouse and touch input controls
- Save files

## Flags

- `-silent` -> All console messages will not be outputted
- `-version` -> Prints version information and exits
- `-saveFiles` -> Saves all game progress and window parameters to separate files. Progress will be imported from these files as well if the flag is present (false by default in order for web to work out of the box)


## Build

You'll need Go installed (and `make` if you're on *nix system for an easier build).

`cd` in `src` directory, run `go build`, now you should have a single binary for your platform which contains the whole game with all resources embedded. Highly portable!

or

Run `make` in the root directory of repository. You will get the binary in the `bin/desktop` directory. Makefile also allows you to easily cross compile to WASM with `make web` and Windows along with your platform with `make desktop`. To compile both for web and for desktop - run `make cross`. 

## License

AGPLv3