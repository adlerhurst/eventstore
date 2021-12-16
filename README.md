# Eventstore

eventstore is a library where I try out new stuff related to an eventstore as a single point of truth.

At the moment I'm writing this it's not meant to evolve to a service. If you want to use it as a framework in your service I'm happy to share my thaughts with you.

My initial intention of this repo is to showcase ideas i develop durring days, months or even years working with eventsourcing and CQRS on the [groundbreaking IAM called ZITADEL](https://github.com/caos/zitadel).

The current development relatest to the [subject based messaging pattern](https://docs.nats.io/nats-concepts/subjects) of [NATS](https://nats.io). I enjoy to read what these people are doing.

## Why?

I love eventsourcing. It helps me a lot durring engineering processes because it is crystal clear that everything I decide now won't change. I can change the future but i can't change the past. Happily, thanks to GDPR the internet has to forget my data. With this fact in mind as an application developer it must be possible to manipulate the past without loss of an activity stream.

It makes development hard because you have to think of what you de before you start doing it, a definition of an event can evolve durring time but I'm not able to enrich information to an event after it happened.
