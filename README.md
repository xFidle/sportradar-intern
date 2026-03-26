# Sports Events Dashboard

## Overview
Sports Events Dashboard is a simple full-stack application for browsing and creating
sports events across multiple sports and competitions. It lets users browse upcoming
and finished matches, apply filters to quickly find relevant events, and view key match
details like participants, schedule, and results. The project is designed to present
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

When application is running, simply enter the following url:
```
http://localhost:5050
```

### Backend
Backend is a REST API written in `Go` mostly  using language essentials (`context` and `net/http` packages). 

#### External modules
- [chi](https://github.com/go-chi/chi) - lightweight and idiomatic HTTP router
  used to define API routes. It keeps requests handling simple and well-organized
  while remaining fully compatible with standard library interfaces.

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
business (services) and transport (HTTP handlers). It is a standard,
well-established design that promotes clear separation of concerns, making
different parts of the system independent and easier to modify or replace. 

To improve error handling and communication between the client and server, a
simple `httpx` library was introduced. It provides utility functions for writing
JSON responses and defines a set of common error types, making message exchange
more consistent and easier to manage.

Backend service prepends assets URLs stored in database with configured file
server address. This approach makes the storage is easily swappable - only the
environment variable defining the file server address needs to be updated.

#### Important decisions
To avoid N+1 query patterns (query in loop) data is fetched in batches and
assembled in services, which reduces database round-trips and keeps response
building predictable. Kept SQL explicit and centralized, then mapped results
into domain DTOs. This makes complex reads easier to reason about and avoids
hidden ORM behavior.

#### Potential improvements
- Strengthen error handling by introducing clearer error categories and
consistent wrapping. Currently, only a small set of generic errors (e.g.
`InvalidPayloadError`, `InternalFailureError`) is used. These provide
consistent responses, but sometimes is hard to distinguish between different
failure scenarios.
- Add automated testing (unit + integration) for different layers of
application to protect critical flows like filtering and event creation. 
- Introduce generated API docs (for example by using
[swagger](https://github.com/swaggo/swag)) to increase readability and
accessibility.

#### Endpoints
```
POST "/api/events"                         - creates new event
GET "/api/events/{event_id}"               - returns full details for single event
POST "/api/events/filter"                  - returns many events based on provided filter

GET "/sports/"                             - returns all supported sports
GET "/sports/{sport_id}/competitions"      - returns competitions associated with specified sport
GET "/sports/{sport_id}/teams"             - returns teams associated with specified sport
GET "/sports/{sport_id}/event-options"     - return competition/venue pairs for the specified sport

GET "/competitions/{competition_id}/teams" - returns teams associated with the specified competition
```

### Database

#### ERD
<img width="1708" height="1120" alt="image" src="https://github.com/user-attachments/assets/002d589b-4003-4ec4-9263-019ecc6c56a9" />

#### Overview
Database represents a multi-sport competition system designed to manage sports
events. It is mostly focused on team sports. Athletics sports could be
supported, but it would require custom logic on backend to calculate scores
properly and generalize player/team types. Schema focuses on keeping 3NF.

Featured relations:
- **Sports & Competitions** - define different sports and organize them into
competitions (leagues or tournaments), with teams assigned to each
competition.
- **Events** - represents individual games, including their schedule, venue,
and stage (e.g. group stage, quarterfinal).
- **Teams & Players** - teams are tied to a sport and location, while players
belong to teams and include basic profile data. While detailed event is retrieved
its current roster is also returned.
- **Participants & Scores** - events include participating teams, and scores are
tracked per participant supporting segmented scoring (e.g. half, quarter).
- **Locations** - venues are organized by cities and countries, with support
for mapping sports to playable venues (**playgrounds**) table.

#### Potential improvements
To simplify repeated query logic and keep API reads consistent read-optimized
SQL views could be created. By encapsulating complex joins into database views,
the  application layer becomes cleaner and less error-prone. What is more, I
think that query execution plans could be analyzed to identify and address 
potential performance bottlenecks.

### Fileserver
Fileserver introduces separation between backend and static assets storage. In
current design database only stores URLs to images, which are returned by API
and later used on frontend. To fetch images I used simple `python` script to
query external API ([sportsdb](https://www.thesportsdb.com/)) and store
results.

### Frontend

#### Mainpage
<img width="1476" height="794" alt="image" src="https://github.com/user-attachments/assets/0f135512-c394-4e74-b13f-6b6991284043" />

#### Event details
<img width="1340" height="1031" alt="image" src="https://github.com/user-attachments/assets/2ee6bb4b-011d-48ba-a0c0-b2145aa50380" />


#### Overview 
Frontend is a simple web application written in vanilla JavaScript, HTML and CSS. It allows
for browsing and adding new events.

#### Potential improvements
To reduce latency and backend load caching layer could be introduced. It would
store frequently requested API responses (for example those related to
filtering). Additionally, I think migration to modern, lightweight frontend
framework (like `vue` or `svelte`) could be considered to improve state
management and DOM manipulation. Managing plain JS code sometimes can
become inconvenient.
