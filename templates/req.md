# Requirement

The website presents a login screen with username and password, along with the 'login' action.

## Login happy path

Valid login attempts **shall** move the  the user to the 'logged in' state

## Login fails wrong password

Invalid login attempts **shall** leave user in the 'logged out' state

## Record login fails

Invalid login attempts **shall** be recorded in the database

# Standown period

Three sucessive invalid login attempts for any user within a 3 minute period **shall** prevent the user from logging in for 30 mins

