# Design Decisions Log

## Use external client library to make requests to the twitch.tv API?

- There is no official client library provided by twitch for Go.

- There are unofficial libraries available.

> <b>Outcome</b>: As this is a technical test and no official library exists, implementing the API requests manually provides more transparency and control over the process, which better demonstrates engineering fundamentals.


## helix/users returns an array of user data. 

- Assumption: this handles returning multiple users where not exact match is found

- approaches to handle if client receives multiple user data items:

    - error or log a warning on array length > 1 , peform aggregation for first return user only
    - aggregate statistics for all users returned by the search query

- points to consider:

    - gather info from API docs and reference here before making a concrete decision!

> <b>Outcome</b>: still under review
