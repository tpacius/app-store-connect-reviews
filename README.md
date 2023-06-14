## App Store Connect Reviews

To run this project, first enter the backend directory `backend/`

`$ cd backend/`

And start running the Go server

`$ go run main.go`

Next enter the frontend directory
`$ cd ../frontend/review-client`

And start the React client 

`$ npm run dev`


## TODOS/Known Issues/Future Plans
I ran into some complications trying to use time.Time as a function parameters while trying to refactor some functionality into seperate functions

Due to some time constraints, I did not get to add unit tests as planned

Other than adding unit tests, additional things that I would want to incorporate as enhancements:

- Taking in parameters to allowing a user to change the window from 48 hours to whatever they'd like
- Smarter logic to determine how many pages to fetch when calling the app store url
- Add a new struct for the persisted json to reduce some of the excessive nesting from the initial app store response json (which I assume is from Apple storing most of this data as XML?)