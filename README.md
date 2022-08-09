# meal-planner

[![CI](https://github.com/roberto-aldera/meal-planner/actions/workflows/go.yml/badge.svg)](https://github.com/roberto-aldera/meal-planner/actions/workflows/go.yml)

An application to plan meals given a set of constraints.

## What is this?
Coming up with a balanced weekly meal plan can sometimes be quite tricky as we tend to fall into patterns and forget about alternative recipes we haven't made in a while.
This small application aims to take some of the burden out of this task.
Given a database of meals with a few characteristics defined, it'll run through a customised policy to find a reasonable set of 7 different dishes, avoiding duplicate food types or impractical plans as a human would if doing this task manually (maybe having pasta 5 times a week is not sustainable, or cooking a 3-hour recipe on a Monday evening).

## Configuration
You'll need to set up a sensible starting [configuration](https://github.com/roberto-aldera/meal-planner/blob/main/default_config.json) to begin with, that's been built to allow some flexibilty to accommodate certain preferences.
It's a very manual process at this stage (open up your favourite text editor and start setting those bools) but perhaps one day there'll be a front-end.
Some of the basics include setting specific days of the week as days you'd prefer a quicker dish (or a more complex one), any meals to exclude (perhaps you had them last week?), and a few other bits and bobs that aren't yet documented.

## Running the meal planner:
Set up [golang](https://go.dev/) so that you can run `.go` files directly, and then run from the command line:
```
go run main.go --config path/to/your/custom_config.json
```

## Running tests
To run all tests, navigate to the top-level directory and run:
```
go test ./... -v
```
for verbose output.
For coverage reporting in HTML format, run:
```
go test ./... -coverprofile cover.out && go tool cover -html=cover.out
```