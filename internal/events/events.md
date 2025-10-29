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

This module aims to solve the user event logging for this server implementation. It is supposed to be only dependent on the Postgres table related to user events, so it should be independent of the purpose of the server itself. 


### Database Integration

The `user_events` module stores and retrieves events in a Postgres database. It relies on a single table, `user_events`, which tracks user activity. This module is independent of the rest of the server, aside from the foreign key reference to the `users` table.

#### Schema

| Parameter      | Type       | Description                        |
|----------------|-----------|------------------------------------|
| id             | Serial     | Primary key                        |
| user_id        | UUID       | References `users.id` (foreign key) |
| method         | Text       | HTTP method                        |
| method_details | Text       | Additional information on metadata |
| created_at     | Timestamp  | Timestamp of the method log        |

#### Queries

The module provides several queries to manage and retrieve user events:

| Query Name | Type | Description |
|------------|------|-------------|
| `InsertEvent` | `:exec` | Inserts a new user event into `user_events`, automatically setting `created_at` to the current timestamp. Returns the inserted row. |
| `ResetEvents` | `:exec` | Truncates the `user_events` table and resets the ID sequence. |
| `GetEvents` | `:many` | Retrieves all events, ordered by newest first (`created_at DESC`). |
| `GetEventCount` | `:one` | Returns the total number of events in the table. |
| `GetEventsInTimeWindow` | `:many` | Retrieves events whose `created_at` timestamp is within a specified time range. |
| `GetEventsByUser` | `:many` | Retrieves all events for a specific user, ordered by newest first. |
| `GetEventsByAction` | `:many` | Retrieves all events with a specific `action` type. |
| `GetEventsByEndpoint` | `:many` | Retrieves events targeting a specific endpoint (extracted from JSON `method_details`). |
| `CountEventsByUser` | `:one` | Counts all events performed by a specific user. |
| `CountEventsByAction` | `:one` | Counts all events of a specific `action` type. |
| `CountEventsByIP` | `:one` | Counts all events originating from a specific IP (extracted from JSON `method_details`). |
| `GetEventsByIP` | `:many` | Retrieves all events from a specific IP address, ordered by newest first. |
| `GetLatestEvents` | `:many` | Retrieves the most recent events, limited by a given number. |
| `GetLatestEventsByUser` | `:many` | Retrieves the most recent events for a specific user, limited by a given number. |

#### Notes

- `user_id` has a foreign key constraint referencing the `users` table.
- `method_details` stores structured data (JSON) with extra information such as IP and endpoint.
- Most queries order results by `created_at DESC` to provide the latest events first.


### Essential functions

#### Function A 

#### Function B 

### Missing functionality

In this section functions and features to come will be listed and discussed. A rough description and an aproximated time of implementation will be given.

#### Feature A 

#### Feature B 

### Notes

### References 

