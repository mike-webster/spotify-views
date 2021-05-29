# TODO ideas

## Implement Word Cloud
- This should be pretty easy, since it's already a client -> service call
- Create the new component, but remember to move the route to the api group too
- Adjust listing drop down handler

## Implement Top Tracks
- make request
- move route in service
- new component
- adjust listing drop down handler

## Implement Top Artists
- make request
- move route in service
- new component
- adjust listing drop down handler

## Implement Top Genres
- make request
- move route in service
- new component
- adjust listing drop down handler

## Implement Library
- Is this actually helpful in any way?
    - If not, could it be?

## Allow Toggling Artist/Track For Recommendations

## Build Out Recommendations Algorithm
- This is the main reason, remember

## Deprecate Old Service Routes

## Delete Old Service Templates

## Deploy React App To Production
- This will cause an outage, since the server is currently sitting on 3000
    - We'll need to move it to 3001
    - Probably need to adjust nginx
- Add site enabled for the new app

## Monitoring For Errors
- Really need a way to get this back

## Handle Redirect Loop
- When `sv-authed` is set but the cookies are missing, we'll end up in a
  redirect loop between / and /discover because we think we're logged in
  but when we try to make a request we get an error.
  - Maybe be a little smarter on the redirect? Delete the cookie, if it exists,
    before sending the user back to /?