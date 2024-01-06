# Hevonen (he·vo·nen)

Simplify some of the organization and tasks a riding school has to do on a daily basis.

The project might be overly complex for the actual use case and scale, but since it is a project for me in my spare time I wanted to use stuff I'm not able to use at work and try out implementations which are not possible at work.

## Features

### Planned
- admin area
  - user managment
  - auth managment
- member management
- working hours tracking
- riding plan


### User Managment

The registartion and login process is currenlty implemented as own service. In a later stage this should probably be changed and a solution like (ory)[https://www.ory.sh/] should be used. For now this was the easier solution.

To be able to switch out the IDP, the services itself only check a token. The token validation itslef is in a shared module. This should enable an easy switch of the used IDP later.

## Architecture/Design Decision

Hevonen is build using a service oriented architecture. Each service is responsible for a specific part of the applicaiton. Communications between service can be asyncronus via a pub/sub mechanism (using Redpanda) or it can call the service directly using a client from the services client module.

### Structure of a service
The service modules should follow the same folder strucutre.

- `main.go` - start a the service
- `server/` - contains the setup for the server. This way, the application can also be started as a single service
- `services/` - service implementaions used by the handlers
- `db/` - database repositries used by services
- `db/migrations/` - sql migration scripts for the database
- `client/` - client implementation to be used by other moduels to communicate with this service.


## Tech Stack

- [golang](https://go.dev/) - just wanted to get try out a new language for me and golang has a low barrier to get started
- [echo](https://echo.labstack.com/) - popular and with a lot of documentation. Why not.
- [templ](https://templ.guide/) - typesafe templates. Yes, I like type safety!
- [htmx](https://htmx.org/) - jump on the new hype train!
- postgreSQL - stable and established db with good services for a free tier. Don't need any fancy unproven stuff for the DB for now.
- [tailwindCSS](https://tailwindcss.com/) - keept it to the standard since I'm not a frontend dev
- Web Components - use the html standard instead of a framework


## Setup

- generate certs `go run /usr/local/go/src/crypto/tls/generate_cert.go --host localhost`
- generate templates `templ generate`
- (install air)[https://github.com/cosmtrek/air]
- (start ory proxy `ory proxy --project <project.slug> http://localhost:4443`)


### Requirements
- go verion >= 1.21
- [install tern](https://github.com/jackc/tern) for database migrations
- [install templ](https://templ.guide/quick-start/installation)


## Resources
### Icons
- favicon by [](https://freeicons.io/profile/417342)https://freeicons.io/filled-lineo-sport-28826/equestrian-horse-racing-jockey-horseback-riding-horse-icon-1024284