# logstv

Logstv is an open source logging platform for twitch.tv (NOT ENDORSED BY TWITCH.TV).
Most of the stuff is still in development.

# TODO 

- verify sort. Currently I rely on cassandra default sort which sorts by partition key, the timestamp is in that so should be fine
- use versioning for go dependencies
- think about deployment. Travis? Jenkins? Docker? K8?
- write a frontend