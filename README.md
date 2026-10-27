# Kaelly-registrar

[![CI](https://github.com/kaellybot/kaelly-registrar/actions/workflows/ci.yml/badge.svg)](https://github.com/kaellybot/kaelly-registrar/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kaellybot/kaelly-registrar/branch/main/graph/badge.svg)](https://codecov.io/gh/kaellybot/kaelly-registrar) 

Application to update Discord commands written in Go

Each run registers the commands of one game, set by `GAME`: 1 (DOFUS), 2 (DOFUS Touch) or 3 (DOFUS Retro).
`CLIENT_ID` and `TOKEN` must belong to that game's Discord application. Registering replaces the application's
whole command set, so in production the registrar refuses a `CLIENT_ID` that isn't the game's own application
(see `models/constants/games.go`).
