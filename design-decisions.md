# Design Decisions Log

---

## Use External Client Library to Make Requests to the Twitch.tv API?

- There is no official client library provided by Twitch for Go.
- There are unofficial libraries available.

> **Outcome**: As this is a technical test and no official library exists, implementing the API requests manually provides more transparency and control over the process, which better demonstrates engineering fundamentals.

---

## Use of Generics in the Helix Client

### Decision

Generics are used in the Helix client's HTTP execution layer over alternatives such as using `interface{}` or `any`.

### Rationale

- **Type Safety**  
  Using generics ensures that only expected response types are handled at compile-time, reducing the likelihood of runtime errors or misinterpretation of data structures.

- **Eliminates Runtime Type Assertions**  
  Avoiding `interface{}` removes the need for manual type assertions or reflection, which can introduce bugs and add unnecessary complexity.

- **Improved Performance**  
  Generics allow data to be decoded directly into concrete types without intermediate structures or interface conversions, reducing overhead.

- **Clarity and Maintainability**  
  The use of generics makes the client easier to reason about and extend, with clear boundaries on what types are expected where. This improves readability and supports safer future changes.

> **Outcome**: Use Go generics with a constrained union of allowed response types to ensure type-safe, efficient, and maintainable request handling.

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

## "Title of the most viewed video along with its view count"

The technical specification includes a requirement for:  
> "Title of the most viewed video along with its view count"

This was presented as a single bullet point, and it was unclear whether the intended response format should:

- Combine both values into a single string (e.g. `"Most viewed: <title> (views: <count>)"`), or  
- Return them as distinct values.

### Rationale

- Returning a combined string obscures the view count and makes it harder to consume in downstream processes (e.g., UI, analytics, testing).
- Returning `title` and `view_count` as separate fields within a nested struct preserves semantic clarity.
- Grouping these fields under a `most_viewed_video` object clearly communicates their relationship.
- This structure also allows for future extensibility — if more metadata about the most viewed video is needed later (e.g., duration, URL), it can be added to the struct without disrupting the response shape.

> **Outcome**: Return a `most_viewed_video` object with `title` and `view_count` as separate fields.