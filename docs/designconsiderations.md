# Design Considerations

1. Use JSON for sending and receiving data where appropriate
   - e.g. Receive JSON input to compute sum @ http://localhost:8080/files/sum
   - endpoints return JSON output to the caller
2. Use nouns instead of verbs for endpoints (e.g. /files and /files/sum)
3. Name collections with plural nouns (e.g. /files)
4. Respond with standard HTTP status code (e.g. 200, 400, 404, 500)
5. Use nesting on endpoints to show relationship between endpoints
   - /files/sum calculates the sum of inputs from files
6. Include API documentation both in /docs and README.md
7. Versioning API (as of now: v1/files)
8. Return HTTP 400 if client inputs invalid data into the request
9. Return HTTP 404 if the resource does not exist
