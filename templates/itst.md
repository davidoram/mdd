# Inspection test

The Inspection test documents a manual inspection test to cover a specific requirement. It links that requirements and the inspection test scripts that cover them.

**When you create a new document, delete from this line to the top of the document, and alter the example sections below to suit your situation.**

# TODO Place your Inspection test title here eg: 'Test login flow'

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
