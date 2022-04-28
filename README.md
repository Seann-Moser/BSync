# BSync
Download Custom maps for BeatSaber from [bsaber](https://bsaber.com/)

- a standalone executable that syncs/installs the following:
  - Bookmarked maps
  - Top Maps
  - Curated
  - Ranked
  - Recent Maps
  - Playlists

## Using
- Download it from [Releases](https://github.com/Seann-Moser/BSync/releases)
  - Currently, it is only built for windows
  - But this repo can be built locally for any os if cloned locally and built using go
  - `go build .`

### Standalone
- Extract .exe file from .zip file
- Double-click on .exe
  - Note windows will say this is not safe to run, but you can read through this repo if you are unsure
- You will have to anwser a few questions like:
  - What is your BSaber username?
  - Where is beat saber installed? - Note if this is in the default C: directory you can just press enter
#### Optional
- This will generate the following 2 files and start the download process:
  - `user_config.json`
  - `songs_config.json`
- This stores the default configuration for running the application
- You can easily add new playlists or links to download by adding the following to the `songs_config.json`

```json
    {
        "url": "https://bsaber.com/songs/top/?time=30-days",
        "amount": 40,
        "min_rating": 0.5,
        "difficulty_levels": null
    }
```
- `url` - The link you wish to download
- `amount` - the number of songs you wish to download from this collection
  - `-1` will download all songs in that category, strongly recommend this for bookmarked and small playlists only
- `min_rating` - this is the lowest rating song you wish to download. only values between `0 - 1`
  - `0.5` rating will have to have at least equal likes to dislikes on bsaber 
  - `1.0` rating will require the song to have 100% likes on the song to download
  - `0.0` rating will download every song regardless of the ratings
- `difficulty_levels` - (WIP) will be able to set the difficulies that you want installed
  - ex: ["expert","expert+"]
  - this avoids downloading easier songs


### CLI
- There is a CLI tool for downloading songs
```
.\BSync.exe -h
```
#### Commands
Note: if you have a custom beat saber path you will have to use the `-p` flag with the path to beat saber
ex:
```
-p "D:/Steam/steamapps/common/Beat Saber/Beat Saber_Data/CustomLevels/"
```

`song-search` - will search for the given value on bsaber and download songs associated to that search field
```
.\BSync.exe song-search -s "Test" -a 40 -r 0.6
```
 - This command will download 40 songs for `Test` that have a rating greater or equal to 0.6


`song` - 

```
.\BSync.exe song -b "https://bsaber.com/songs/curated/?recommended=true" -a 100 -r 0.5
```