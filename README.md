# gator

This program requires both Postgres and Go to be installed.

After installing both Postgres and Go, use the following command to install gator:  
 `go install gator`

The program requires a config file to work properly. This config file should be placed in the root of the user's home directory. The name of the config file should be .gatorconfig.json and it should have the following contents:  
{  
&nbsp; &nbsp; "db_url": "postgres://postgres@localhost:5432/gator?sslmode=disable",  
&nbsp; &nbsp; "current_user_name": ""  
}  

The dB URL should be updated to the Postgres database you created for the 

The program accepts the following commands:
```
 - register <name>      Creates a new user with name <name> and switches to the new user.
 - login <name>         Switch the current user to <name>.
 - users                List all users.
 - addfeed <name> <url> Add a feed with name <name>. The data for the feed will be pulled from URL <url>.
 - feeds                List all feeds that have been added and which user added the feed.
 - follow <url>         Follow the feed with URL <url>. The feed has to exist in the list of existing feeds.
 - following            List all the feeds the current user is following.
 - unfollow <url>       Unfollow a feed.
 - browse <n>           Display the <n> most recent posts from all feeds the user is following. <n> defaults to 2 if no <n> is given.
```
