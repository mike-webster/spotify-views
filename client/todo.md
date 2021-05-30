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

## Handle 429s
- We're getting a lot of 429s to the point that it's disrupting development
    - Sure, it makes quite a few calls, but I can't even get recommendations to
      work before it 429s.
- Options
    - Mock calls in development
        - This would potentially make it harder to find bugs until prod
    - Stand up request cache
        - I tried this before with Redis before I had my own server stood up.
            - At the time, negatives were: paying for a host and paying for storage
        - Redis alternatives?