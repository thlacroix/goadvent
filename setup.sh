# setup.sh is a script to be sourced (either manually or in your (bash|zsh)rc file)

# Automatically export AOC_REPO is not set already
if [ -z "$AOC_REPO" ]; then export AOC_REPO="$(dirname $0)"; fi

# AOC daily setup function available when sourcing this file. Usage: aoc YEAR DAY
# Example usage: aoc 2020 7
# It will:
# * Move in the root of this git repo based on AOC_REPO env var
# * Create a base folder based on the year an day provided
# * Copy the template in this folder
# * Get the input from AOC if you provide AOC_SESSION in the env (extracted from a browser cookie)
# * Open the root of this repo in your VISUAL editor (if set)
# * Move in the new folder and run the code to validate the setup
function aoc {
    year="$1"
    day="$2"

    if [ -z "$AOC_REPO" ]; then echo "AOC_REPO is not set"; return 1; fi
    echo "Moving in git repoitory in $AOC_REPO"
    cd "$AOC_REPO"
    base_folder="$year/day$(printf %02d $day)"
    echo "Creating base folder $base_folder"
    mkdir -p "$base_folder"
    if [ -f "$base_folder/main.go" ]; then echo "Day $base_folder is already setup"; return 1; fi
    echo "Copying template to $base_folder/main.go"
    cp template.go "$base_folder/main.go"
    if [ -n "$AOC_SESSION" ]; then echo "Getting input from AOC"; curl -s "https://adventofcode.com/$year/day/$day/input" --cookie "session=$AOC_SESSION" > "$base_folder/input.txt"; fi
    if [ -n "$VISUAL" ]; then echo "Opening git repo with $VISUAL"; $VISUAL .; fi
    echo "Moving in new folder"
    cd "$base_folder"
    echo "Simple go run to validate setup"
    go run .
}