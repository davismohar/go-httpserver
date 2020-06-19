# go-httpserver
Basic multithreaded http server implemented in go.

User accounts are stored in a SQLite database. When a user logs in, they are issued a JWT (Javascript Web Token).
Any requests to access any files within the private/ directory require the user to present a valid JWT token.

Directory layout:
server-root/
  private/ --Access to this directory requires a valid JWT Token
    privateHome.html
    secret.txt
  index.html
  login.html
  createAccount.html
    
