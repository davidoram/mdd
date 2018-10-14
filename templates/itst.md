# Inspection test

## Login happy path

Steps

- Navigate to the website
- Click 'login'
- Enter username `fred.smith` & password `secret`

Test pass

- Shows banner `Welcome fred`
- User enters login state


## Login fails wrong password

Steps

- Navigate to the website
- Click 'login'
- Enter username `fred.smith` & password `wrongone`

Test pass

- Shows error message `Try again`
- User remiains in logged out state

## Standown period

Steps

- Navigate to the website
- Click 'login'
- Enter username `fred.smith` & password `wrongone`, four times within (3 mins)

Test pass

- Shows error message `Locked out, try again in a while`
- User remiains in logged out state
