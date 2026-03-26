# Sports Events Dashboard

## Overview
Sports Events Dashboard is a simple full-stack application for browsing and creating
sports events across multiple sports and competitions. It lets users browse upcoming
and finished matches, apply filters to quickly find relevant events, and view key match
details like participants, schedule, and resutls. The project is designed to present
sports data in clear, user-friendly format.

## Architecture
The application consists of four main services:
- **Backend**: REST API
- **Database**: PostgreSQL with seeded sports data
- **File server**: Static asset hosting for team and competition badges
- **Frontend**: Web client for browsing and creating events

All services are orchestrated with `Docker Compose`, so setup is straightforward.
To start the application, run:
```bash
docker compose up
```

### Backend
Backend is a REST API written in `Go` mostly  using langugage essentials (`context` and `net/http` packages). 

#### External modules
- [chi](https://github.com/go-chi/chi) - lightweight and idiomatic HTTP router
  used to define API routes. It keeps requests handling simple and well-organized
  while remaining fully compatibile with standard library interfaces.

- [sqlc](https://github.com/sqlc-dev/sqlc) - a database access tool that
  generates type-safe Go code directly from SQL queries. It allows writing plain
  SQL while still getting clean, strongly typed repository methods. I chose it because
  I prefer writing SQL manually rather than using ORMs like GORM.

- [validator](https://github.com/go-playground/validator) - a package for validating
  marshalled structs based on field tags. It’s used to ensure that all required
  DTO fields are provided, while any additional validation logic is handled
  manually.

- [copier](https://github.com/jinzhu/copier) - a utility that provides
  field-to-field copying based on matching names. It is used to map values between
  different layers of application, reducing the need for explicit mapping.

#### Layers
Application is structured into three layers - data access (repository),
bussiness (services) and transport (HTTP handlers). It is a standard,
well-established design that promtes clear separation of concerns, making
different parts of the system indepnedent and easier to modify or replace. 

To improve error handling and communication between the client and server, a
simple `httpx` library was introduced. It provides utility functions for writing
JSON responses and defines a set of common error types, making message exchange
more consistent and easier to manage.

Backend service prepends assets URLs stored in database with configured file
server address. This approach makes the storage is easily swapabble - only the
environment variable defining the file server address needs to be updated.

#### Important decisions
To avoid N+1 query patterns (query in loop) data is fetched in batches and
assembled in services, which reduces database round-trips and keeps response
building predictable. Kept SQL explicit and centralized, then mapped results
into domain DTOs. This makes complex reads easier to reason about and avoids
hidden ORM behavior.

#### Potential improvements
error handling

testing

api docs (swagger)


#### Endpoints
```
POST "/api/events"                         - creates new event
GET "/api/events{event_id}"                - returns full details for single event
POST "/api/events/filter"                  - returns many events based on provided filter

GET "/sports/"                             - returns all supported sports
GET "/sports/{sport_id}/competitions"      - returns competitions associated with specified sport
GET "/sports/{sport_id}/teams"             - returns teams associated with specifed sport
GET "/sports/{sport_id}/event-options"     - return 

GET "/competitions/{competition_id}/teams" - 
```

### Database

#### ERD
<img width="1708" height="1120" alt="image" src="https://github.com/user-attachments/assets/002d589b-4003-4ec4-9263-019ecc6c56a9" />


####

#### Potential improvements
indexes

### Fileserver
separation between assets bucket, backend and db (db stores only urls)

### Frontend

#### Mainpage
<img width="2560" height="1293" alt="image" src="https://github.com/user-attachments/assets/56d51b85-d3c2-49bd-acb3-602728d88511" />
#### Potential improvements
caching

modern framework for easier DOM and state manipulation

