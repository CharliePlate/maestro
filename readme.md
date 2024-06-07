## Maestro

Not totally sure where this project is going to go yet, but the idea is to be a message broker that uses database watch streams to handle managing the queues

### Things to do now:

- \[x\] Basic Container Implementation
  - The goal is to eventually have more containers like a heap, but for now just a slice is fine
- \[x\] Mongo Change Stream Watcher
- \[x\] Protocol parsing
  - Decided to get really fancy here with `struct tags`. Probably overkill
- \[ \] Protocol Buffer Implementation
  - Initial thought it to have the protocol send some version number, content length and then the data as a protobuf.
- \[ \] Peer Subscribing/Unsubscribing
- \[ \] Queues sending data and receiving acknowledgements (probably some more protobuf work)
- \[ \] Message type that will probably be some fixed length so I know if someone is subbing, acking, ect.

### Building

There really is no building yet...

The docker-compose file will start up a Mongo instance that has a replica set so that change streams work. This is all for now
