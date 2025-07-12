# Design Decisions Log

---

## Use External Client Library to Make Requests to the Twitch.tv API?

- There is no official client library provided by Twitch for Go.
- There are unofficial libraries available.

> **Outcome**: As this is a technical test and no official library exists, implementing the API requests manually provides more transparency and control over the process, which better demonstrates engineering fundamentals.

---

## `helix/users` Returns an Array of User Data

- **Assumption**: The Helix API returns an array in case multiple usernames are queried or if partial/fuzzy matches are made. However, when querying with a specific `login` (username), we expect at most one exact match.

### Handling Multiple Returned Users

Two main approaches were considered:

1. **Use only the first returned user (and warn if others are present)**  
2. **Aggregate data across all returned users**

### Decision: Use Only the First Returned User

After reviewing the [Twitch Helix API documentation](https://dev.twitch.tv/docs/api/reference#get-users), the `login` query param is expected to return an exact match. The response is still wrapped in an array (likely for consistency across the API), but only one item should be present if the provided login exists.

> **Outcome**: We will use the first user in the returned array and log a warning if more than one user is returned. This is treated as an unexpected condition.

### Rationale

- The API expects exact login values; if the value is valid, only one user should be returned.
- Aggregating across users is semantically incorrect for our use case (stats are per streamer).
- Returning or processing multiple users could lead to misleading or invalid analytics.
- Defensive programming: logging multiple returns helps catch unexpected API behavior or misconfigurations.
- Simplicity and clarity in both implementation and results.

---

## Averages Calculations

- **Assumption**: Both integers and floating-point values are valid for representing the required averages.
- Integers are also valid, as a streamer never gets a view count between `0 < x < 1`.

> **Outcome**: For simplicity in consuming the information, calculate and display these figures as integers.
