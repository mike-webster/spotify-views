# TODO ideas

## Implement Word Cloud
- This should be pretty easy, since it's already a client -> service call
- Create the new component, but remember to move the route to the api group too
- Adjust listing drop down handler
- This seems it may be broken at the service layer... This will take some troubleshooting

## Implement Library
- Is this actually helpful in any way?
    - If not, could it be?

## Allow Toggling Artist/Track For Recommendations

## Build Out Recommendations Algorithm
- This is the main reason, remember

## Delete Old Service Templates

## Deploy React App To Production
- This will cause an outage, since the server is currently sitting on 3000
    - We'll need to move it to 3001
    - Probably need to adjust nginx
- Add site enabled for the new app