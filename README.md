# Anigo
Anime list tracking software written in go

With anigo you can keep track of your animelist, and sync it with MAL (please refrain from submitting your local list to MAL, as it will delete your scores)

## Setup

You'll have to create an 'anigo' directory in your ~/.config/ folder or equivalent

## Usage

List 'Watching' shows (-C for completed or -P for planned too)
```
anigo -W
```

Add an anime with title "Little Witch Academia (TV)", 25 total chapters and Completed status to the list
```
anigo add -t "Little Witch Academia (TV)" -c 25 -C
```

Search for anime named "Cowboy Bebop" AND is completed
```
anigo [query] -T "Cowboy Bebop" -C
```

## Operation List

* add, -A: Add an entry to the database
* del, delete, -R: Remove one or more entries matching the criteria
* edit, -E: Edit all entries mathcing the criteria. To set values use --set-<value>
* pull: Pull your entries from MyAnimeList.net
* push: Push your entries to MyAnimeList.net (not recommended)
* search, query, -Q: Do a search (default if no operation is found)
* --spit: spit autocomplete from template

## Options List

|Command|Description|
|-----|-----|
|-P, -W, -C, -D | Set status as Plan To Watch, Watching, Completed and Dropped respectively|
|--chapters, -c|Set total number of chapters|
|--debug|Enable debug logging|
|--set-{chapters, completed, status, title}|Using the edit operation, change the target property of ALL matching anime to target value|
|--title, -t|Set anime title|
|--user, -u|Set MAL user|
|--password, -p|Set MAL password|
|--watched, -w|Set number of watched chapters|
