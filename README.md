# go-httpserver
Basic multithreaded http server implemented in go.

User accounts are stored in a SQLite database. When a user logs in, they are issued a JWT (Javascript Web Token).
Any requests to access any files within the private/ directory require the user to present a valid JWT token.

Template directory layout:  
server-root/  
&nbsp;&nbsp;&nbsp;&nbsp;index.html  
&nbsp;&nbsp;&nbsp;&nbsp;login.html  
&nbsp;&nbsp;&nbsp;&nbsp;createAccount.html  
&nbsp;&nbsp;&nbsp;&nbsp;private/ --Access to this directory requires a valid JWT Token  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;privateHome.html  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;secret.txt  

    
