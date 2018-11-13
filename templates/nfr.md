# Non Functional Requirement

The Functional Requirement document specifies criteria that can be used to judge the operation of a system, rather than specific behaviors. Another term for non functional requirements is _Quality goals_

Each non functional requiremnt should be:

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

# TODO Place your Non Functional Requirement title here eg: 'No sensitive information will be logged'

The system contains the following list of sensitive information, which shall not appear in a logfile unless it has been replaced with a santitsed form, for example 'secret' might be replaced with '*****'

The sensitive fields are:

- User password
- API keys
- SSN
- Credit card numbers

When any field is sanitised for logging, the system shall replace it with a string containing 5 stars eg: '*****'

eg: Instead of logging this:

```
Username: Bob, Password: secresquirrel, Role: admin
```

it will log this:
```
Username: Bob, Password: *****, Role: admin
```
