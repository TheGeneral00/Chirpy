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
- [Notes](#notes)
- [References](#references)

### Core idea

This module aims to solve the user event logging for this server implementation. It is supposed to be only dependent on the Postgres table related to user events, so it should be independent of the purpose of the server itself. 

### Database integration

The server is running a Postgres database. As a result of this, the user events module will be based on a Postgres database. The explicit sql schema and queries can be found in the sql dir. Schema are named according to the goose migration rules.

**User events schema**

|Parameter      |Type       |Description        |
|---------------|-----------|-------------------|
|id             |Serial     |Primary Key        |
|user_id        |UUID       |Primary Key and ref|
|               |           |rences user id     |
|methode        |Text       |http methode       |
|methode_details|Text       |additional informat|
|               |           |ion on meta data   |
|created_at     |Timestamp  |time the methode wa|
|               |           |s logged           |

### Essential functions

#### Function A 

#### Function B 

### Missing functionality

In this section functions and features to come will be listed and discussed. A rough description and an aproximated time of implementation will be given.

#### Feature A 

#### Feature B 

### Notes

### References 

