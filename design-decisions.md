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

## Averages calculations

- assumption: both integers and fractionals are valid for representing the required averages

- integers have a degree of validity of a float/double etc, as a streamer never gets a view > 0 && < 1

> <b>Outcome</b>: For simplicity in consuming the information, calculate and display these figures as integers

