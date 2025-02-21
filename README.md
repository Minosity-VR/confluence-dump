# confluence-dump
Dump a confluence wiki pages

This tool has been designed for Red Team exercise in the scenario of a stolen session token.

Blabla use this only when you are authorized to do so blabla and I take no responsibility blabla

## Usage

```bash
git clone https://github.com/Minosity-VR/confluence-dump.git
cd confluence-dump
go build
./confdump --cookie ey... --host company.atlassian.net --output /tmp/confDump
```

The required token is the `tenant.session.token` that you can find in a browser by opening the
developer tools and looking at any requests cookies.

There is often a lot of different key/values in the cookie, only the `tenant.session.token` is useful
