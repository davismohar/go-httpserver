# go-httpserver
Basic http server implemented in go.
All files requested are searched for in the server-root/ directory.
User accounts are stored in a SQLite database. When a user logs in, they are issued a JWT (Javascript Web Token).
Any requests to access any files within the private/ directory require the user to present a valid JWT token.
