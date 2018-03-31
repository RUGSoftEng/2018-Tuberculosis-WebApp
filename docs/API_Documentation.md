# API Documentation

## Overview:
* [Physicians](#physicians)
* [Patients](#patients)
* [Notes](#notes)
* [Dosages](#dosages)
* [Videos](#videos)


## Physicians
### Creating
_To be documented_
### Retrieving
_To be added_
### Updating
_To be documented_
### Deleting
_To be documented_
## Patients
### Creating
_To be documented_
### Retrieving
_To be added_
### Updating
_To be documented_
### Deleting
_To be documented_
## Notes
### JSON
```json
{
	"note": "String: The message of the note",
	"created_at": "Date: YYYY-MM-DD"
}
```
### Creating
#### Adding a note for a specified patient
|---------------|-------------------------------------|
| **Type**      | PUT                                 |
| **Path**      | `/api/accounts/patients/{id}/notes` |
| **Variables** | id: Integer                         |
| **Body**      | JSON: Note                          |
| **Return**    | -                                   |
### Retrieving
#### Retrieving all notes from a specified patient
|---------------|-------------------------------------|
| **Type**      | GET                                 |
| **Path**      | `/api/accounts/patients/{id}/notes` |
| **Variables** | id: Integer                         |
| **Body**      | -                                   |
| **Return**    | List of JSON Notes                  |
### Updating
_To be added_
### Deleting
_To be added_
## Dosages
### JSON
```json
{
	"date": "Date: YYYY-MM-DD",
	"intake_moment": "Time: HH-MM-SS",
	"amount": "Integer",
	"medicine": "String: Name of the medicine which needs to be taken",
	"taken": "Boolean"
}
```
### Creating
|---------------|---------------------------------------|
| **Type**      | PUT                                   |
| **Path**      | `/api/accounts/patients/{id}/dosages` |
| **Variables** | id: Integer                           |
| **Body**      | JSON Dosage                           |
| **Return**    | -                                     |
### Retrieving
#### Retrieving all dosages from a patient in an inverval
|---------------|---------------------------------------|
| **Type**      | GET                                   |
| **Path**      | `/api/accounts/patients/{id}/dosages` |
| **Variables** | id: Integer                           |
|               | ?from: YYYY-MM-DD                     |
|               | ?until: YYYY-MM-DD                    |
| **Body**      | -                                     |
| **Return**    | List of JSON Dosages                  |
### Updating
_To be added_
### Deleting
_To be added_
## Videos
### JSON
```json
{
	"topic": "String",
	"title": "String",
	"reference": "String: url"
}
```
### Creating
#### Adding a new video
|---------------|-----------------------|
| **Type**      | PUT                   |
| **Path**      | `/api/general/videos` |
| **Variables** | -                     |
| **Body**      | JSON Video            |
| **Return**    | -                     |
### Retrieving
#### Retrieving all videos with a certain topic
|---------------|-------------------------------|
| **Type**      | GET                           |
| **Path**      | `/api/general/videos/{topic}` |
| **Variables** | topic: String                 |
| **Body**      | -                             |
| **Return**    | List of JSON Videos           |
