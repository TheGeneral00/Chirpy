# Cocumentation for user event logging

## Contents
- [Core idea](#core-idea)
- [Database integration](#database-integration)
- [Essential functions](#essential-functions)
  - [Function A](#function-a)
  - [Function B](#function-b)
- [Missing functionality](#missing-functionality)
  - [Feature X](#feature-x)
  - [Feature Y](#feature-y)
- [References](#references)

### Core idea

This module aims to solve the user event logging for this server implementation. It is supposed to be only dependent on the Postgres table related to user events, so it should be independent of the purpose of the server itself. Certain statistical computations are implemented in C++ for performance reasons.

### Database Integration

The `user_events` module stores and retrieves events in a Postgres database. It relies on a single table, `user_events`, which tracks user activity. This module is independent of the rest of the server, aside from the foreign key reference to the `users` table.

#### Schema

| Parameter      | Type       | Description                        |
|----------------|-----------|------------------------------------|
| id             | Serial     | Primary key                        |
| user_id        | UUID       | References `users.id` (foreign key) |
| method         | Text       | HTTP method                        |
| method_details | Text       | Additional information on metadata (JSON) |
| created_at     | Timestamp  | Timestamp of the method log        |

#### Queries and Expected Parameters

| Query Name | Type | Description | Parameters |
|------------|------|-------------|------------|
| `InsertEvent` | `:exec` | Inserts a new user event into `user_events`, automatically setting `created_at` to the current timestamp. Returns the inserted row. | `$1: UUID` (user_id), `$2: Text` (method), `$3: Text` (method_details) |
| `ResetEvents` | `:exec` | Truncates the `user_events` table and resets the ID sequence. | none |
| `GetEvents` | `:many` | Retrieves all events, ordered by newest first (`created_at DESC`). | none |
| `GetEventCount` | `:one` | Returns the total number of events in the table. | none |
| `GetEventsInTimeWindow` | `:many` | Retrieves events whose `created_at` timestamp is within a specified time range. | `$1: Timestamp` (start), `$2: Timestamp` (end) |
| `GetEventsByUser` | `:many` | Retrieves all events for a specific user, ordered by newest first. | `$1: UUID` (user_id) |
| `GetEventsByAction` | `:many` | Retrieves all events with a specific `action` type. | `$1: Text` (method) |
| `GetEventsByEndpoint` | `:many` | Retrieves events targeting a specific endpoint (extracted from JSON `method_details`). | `$1: Text` (endpoint) |
| `CountEventsByUser` | `:one` | Counts all events performed by a specific user. | `$1: UUID` (user_id) |
| `CountEventsByAction` | `:one` | Counts all events of a specific `action` type. | `$1: Text` (method) |
| `CountEventsByIP` | `:one` | Counts all events originating from a specific IP (extracted from JSON `method_details`). | `$1: Text` (IP address) |
| `GetEventsByIP` | `:many` | Retrieves all events from a specific IP address, ordered by newest first. | `$1: Text` (IP address) |
| `GetLatestEvents` | `:many` | Retrieves the most recent events, limited by a given number. | `$1: Int` (limit) |
| `GetLatestEventsByUser` | `:many` | Retrieves the most recent events for a specific user, limited by a given number. | `$1: UUID` (user_id), `$2: Int` (limit) |

#### Notes

- `user_id` has a foreign key constraint referencing the `users` table.
- `method_details` stores structured data (JSON) with extra information such as IP and endpoint.
- Most queries order results by `created_at DESC` to provide the latest events first.
- Parameter types correspond to Go function arguments when generating with `sqlc`.

### Essential functions

#### Function A 

#### Function B 

### Missing functionality

In this section functions and features to come will be listed and discussed. A rough description and an aproximated time of implementation will be given.

#### Feature A 

#### Feature B 

### Notes

### References 

