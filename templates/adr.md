# Architecture Decision Record

The Architectural Decision Record captures key technological, or design choices for your system. These choices are made in the context of particular technogical, political, and project forces.  Capture the forces in play, that have led to the decision.

**When you create a new document, delete from this line to the top of the document, and alter the example sections below to suit your situation.**

# TODO Place your Architecture Decision title here: eg 'Serverside framework decision'

## Context

This section describes the forces at play, including technological, political, social, and project local. These forces are probably in tension, and should be called out as such. The language in this section is value-neutral. It is simply describing facts.

eg:

- The development team must approve the selected framework.
- The target architecture is ABC cloud provider
- Any framework must support the Dynamic website and API components required by the application
- Due to the sensitive nature of this project a consideration must be made to frameworks that supportcurrent best practice secure coding techniques.

## Decision

This section describes our response to these forces. It is stated in full sentences, with active voice. "We will ..."

eg:

We will use the 'XYZ' Framework, v8.5.3 (or later).

- The development team has used this framework recently in the ABC and DEF projects.
- ABC cloud provider has native API support for the Framework see [ref](http::/abc.com/api/xyx)
- The XYZ Framework method.

## Status

A decision may be "proposed" if the project stakeholders haven't agreed with it yet, or "accepted" once it is agreed. If a later ADR changes or reverses a decision, it may be marked as "deprecated" or "superseded" with a reference to its replacement.

eg:

This decision is current 'proposed' pending approval from the Development team lead.

## Consequences

This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

eg:

**Pros**

- XYZ Framework uses permissive [MIT License](https://opensource.org/licenses/MIT), so we do not anticipate any ongoing  licencing costs.


**Cons**

- The high level nature of the Framework incurs significant runtime costs

