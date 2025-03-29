## Notes: 

### For the database:

- Updated diagram and schema.sql to the current schema

- Added reactions_id to reactions table fo consistency and so that either post_id or comment_id can be null 

- Added category_id also to the bd and indexed so every category can be faster to trace

- Added constraints in username in user table
``` sql
CHECK (length(username) >= 3 AND length(username) <= 15) -- minimum length 3 and maximum length 15
CHECK (username GLOB '[a-zA-Z0-9_]*') --  ensures username has at least one alphanumeric or underscore character
CHECK (username NOT GLOB '*[^a-zA-Z0-9_]*'), -- ensures username doesn't contain any characters outside our allowed set
```

- Added constraint in hashed_password in user_auth table
``` sql
password_hash TEXT NOT NULL 
    CHECK (length(password_hash) = 60), -- bcrypt always provides a 60 character hash so we should expect this specific length
``` 

### Considerations

- Consider running the db schema from the schema.sql file and not create queries within an array:
    - `Better readability and maintainability:` SQL statements are easier to read, write, and modify when they're in a dedicated .sql file rather than embedded as strings in Go code.
    - `SQL syntax highlighting and validation:` Most code editors provide better syntax highlighting and error checking for SQL when it's in a .sql file, making it easier to spot and fix issues.
    - `SQL-specific tooling:` You can use SQL-specific tools (like formatters, linters, and database migration tools) directly on your schema file.
    - `Cleaner Go code:` Your Go database initialization becomes more concise and focused on the connection logic rather than SQL statements.
    - `Simplifies complex schema management:` As your schema grows more complex, managing it in a separate file becomes increasingly beneficial.

- Consider moving AuthService to its own directory so we can add it to middleware instead of user-session repo