# Functional Requirement

The Functional Requirement document captures one set of prescribed properties or activities that the system must implement to satify the project goals.

Each functional requiremnt should be:

- Unitary: Address one and only one thing.
- Complete: Contain all the information in one doument with nothing missing.
- Consistent: Non-contradictory with respect to itself and other requirements.
- Atomic: Separate requrments should be contained in separate documents.
- Traceable: Thats why you are using `mdd` to capture your requirements :-).
- Current: Has not been made obsolete by the passage of time.
- Unanbiguous: Where possible language should be clear and simple with no technical jargon. It expresses facts not subjective opinions.
- Importance: Representaing the stakeholders importance.
- Vefifyable: The implementation must be able to bde deomnstrated, and tested.


**When you create a new document, delete from this line to the top of the document, and alter the example sections below to suit your situation.**

# TODO Place your  Functional Requirement title here eg: 'User login'

The website presents a login screen containing:

- username field, a single line text field.
- password field, a single line password field which hides the usets entry from view.
- 'login' button

## Login suceeds, happy path

Valid login attempts **shall** move the  the user to the 'logged in' state

## Login fails wrong password

Invalid login attempts **shall** leave user in the 'logged out' state

## Record login fails

Invalid login attempts **shall** be recorded in the database

# Standown period

Three sucessive invalid login attempts for any user within a 3 minute period **shall** prevent the user from logging in for 30 mins

