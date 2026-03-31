# Delta for Infrastructure: Gitlab Configuration

## ADDED Requirements

### Requirement: Gitlab Token Configurator
The system MUST provide a configurator for the `gitlab-token-config` ID that updates global Composer and NPM settings with a Gitlab personal access token.

### Requirement: Composer Global Auth
The Gitlab configurator MUST update the user's `~/.composer/auth.json` file with the Gitlab token for `gitlab.com`.

#### Scenario: Update Composer auth.json
- GIVEN a valid Gitlab token
- WHEN the Gitlab configurator is executed
- THEN it MUST ensure `~/.composer` directory exists
- AND it MUST add or update the `gitlab-token` entry for `gitlab.com` in `auth.json`

### Requirement: NPM Global Auth
The Gitlab configurator MUST update the user's `~/.npmrc` file with the Gitlab token for `gitlab.com`.

#### Scenario: Update npmrc
- GIVEN a valid Gitlab token
- WHEN the Gitlab configurator is executed
- THEN it MUST append or update `//gitlab.com/api/v4/packages/npm/:_authToken=TOKEN` in `~/.npmrc`
