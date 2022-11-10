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

#### Releasing the API
- you can get the date by running `date +%Y%m%d-%H%M`
`MASTER_KEY="test_master_key" GO_ENV=production RELEASE_DATE=20221110-0031 make release_api`
    - This will take care of building and moving, but we have a manual restart step
- ssh to server, navigate to live directory, kill running server `lsof -i tcp:3001`, start new server `startapi`

#### Releasing the client
- you can get the date by running `date +%Y%m%d-%H%M`
`RELEASE_DATE=20221110-0031 make release_client`
- After this, it's good to go. no more steps needed.