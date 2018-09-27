# logstv

Logstv is an open source logging platform for twitch.tv (NOT ENDORSED BY TWITCH.TV).
Most of the stuff is still in development.

# TODO 

- use versioning for go dependencies
- write a frontend
- jenkinsfile for easier deployments
- implement rq: get latest and oldest message by timestamp, generate a random date between latest and oldest --> then what?

# Helpful

- run provision.sh to provision all (currently 1) servers
- run `ansible-vault encrypt_string my_secret_123 --ask-vault-pass` to encrypt a secret
