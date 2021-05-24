# [spotify-views](http://www.spotify-views.com)
Exposing some functionality that is provided by the [Spotify](https://www.spotify.com/us/) [API](https://developer.spotify.com/documentation/web-api/)

## Installation
run: `./bootstrap.sh`

Note: Docker is optional
#### With Docker
- This will pull a mysql image and then configure it as needed
- The image uses a volume so data is persisted beyond the life of one container

#### Without Docker
- Dependencies:
- mysql v8.0.23
- go v1.15.7