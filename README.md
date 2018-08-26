# logstv

Logstv is an open source logging platform for twitch.tv (NOT ENDORSED BY TWITCH.TV).
Most of the stuff is still in development.

# TODO 

- verify sort. Currently I rely on cassandra default sort which sorts by partition key, the timestamp is in that so should be fine
- use versioning for go dependencies
- Improve deployment script
- write a frontend

# Helpful

- run provision.sh to provision all (currently 1) servers
- run `ansible-vault encrypt_string my_secret_123 --ask-vault-pass` to encrypt a secret