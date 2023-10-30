# Eventstore

eventstore is a library where I try out new stuff related to an eventstore as a single point of truth.

At the moment I'm writing this it's not meant to evolve to a service. If you want to use it as a framework in your app I'm happy to share my thoughts with you.

My initial intention of this repo is to showcase ideas i develop durring days, months or even years working with eventsourcing and CQRS on the [groundbreaking IAM called ZITADEL](https://github.com/caos/zitadel).

The current development relatest to the [subject based messaging pattern](https://docs.nats.io/nats-concepts/subjects) of [NATS](https://nats.io). I enjoy to read what these people are doing.

## Why?

I love eventsourcing. It helps me a lot durring engineering processes because it is crystal clear that everything I decide now won't change. I can change the future but i can't change the past. Happily, thanks to GDPR the internet has to forget my data. With this fact in mind as an application developer it must be possible to manipulate or forget the past without loss of an activity stream.

It makes development hard because you have to think of what you do before you start doing it, a definition of an event can evolve durring time but you're not able to enrich information to an event after it happened.

## Ideas

Some ideas which probably will be implemented. The list is unordered and some points might never be implemented.

- [x] allow multiple filters in eventstore.Filter
  - [ ] ~~storage: provide an Optimize method to simpify queries~~
  - [ ] ~~maybe two layers of optimizations would be more useful. First in eventstore to collect filters and one in storage optimized on it's internal data structures.~~
- [ ] additional storage types
  - [ ] sql (crdb) storage
  - [ ] file storage
- [ ] memory: optimize tree
  - [ ] self balanced
  - [ ] check out different tree styles
- [x] testing suite
- [ ] fuzzy testing with go1.18
- [x] Think of an option to register event types to return the concrete type instead of the `Event`-struct (Event would change to interface)
- [ ] Subscriber: add the possibility to listen to message queues
  - [ ] Pub/Sub (NATS, ...)
  - [ ] (Web-)hook
- [ ] Publisher: add the possibility to push events
  - [ ] Third party tools (NATS, ...)
  - [ ] Specifications (MQTT, ...)
  - [ ] No dependencies
    - [ ] Webhooks
    - [ ] GRPC streams
