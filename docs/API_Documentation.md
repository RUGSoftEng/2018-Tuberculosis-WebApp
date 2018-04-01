# API Documentation

## Overview:
* [Physicians](#physicians)
* [Patients](#patients)
* [Notes](#notes)
* [Dosages](#dosages)
* [Videos](#videos)


## Physicians
### Creating
#### Creating a new physician account

| **Type**      | PUT                                 |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/physicians`          |
| **Variables** | -                                   |
| **Body**      | JSON: Physician                     |
| **Return**    | -                                   |

### Retrieving
_To be added_
### Updating
#### Modifying an existing physician account

| **Type**      | POST                                |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/physicians/{id}`     |
| **Variables** | id: Integer                         |
| **Body**      | JSON: Physician                     |
| **Return**    | -                                   |

### Deleting
#### Deleting a physician account

| **Type**      | DELETE                              |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/physicians/{id}`     |
| **Variables** | id: Integer                         |
| **Body**      | -                                   |
| **Return**    | -                                   |

## Patients
### Creating
#### Creating a new patient account

| **Type**      | PUT                                 |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/patients  `          |
| **Variables** | token: String                       |
| **Body**      | JSON: Patient                       |
| **Return**    | -                                   |

### Retrieving
_To be added_
### Updating
#### Modifying an exixting  patient account

| **Type**      | POST                                |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/patients/{id}`       |
| **Variables** | id: Integer                         |
| **Body**      | JSON: Patient                       |
| **Return**    | -                                   |

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

| **Type**      | PUT                                 |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/patients/{id}/notes` |
| **Variables** | id: Integer                         |
| **Body**      | JSON: Note                          |
| **Return**    | -                                   |
### Retrieving
#### Retrieving all notes from a specified patient
| **Type**      | GET                                 |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/patients/{id}/notes` |
| **Variables** | id: Integer                         |
| **Body**      | -                                   |
| **Return**    | List of JSON Notes                  |
### Updating
_To be added_
### Deleting
#### Deleting a patient account

| **Type**      | DELETE                              |
|---------------|-------------------------------------|
| **Path**      | `/api/accounts/patients/{id}`       |
| **Variables** | id: Integer                         |
| **Body**      | -                                   |
| **Return**    | -                                   |

## Dosages
### JSON

**ScheduledDosage**
```json
{
	"dosage": "JSON Dosage",
	"date": "Date: YYYY-MM-DD",
	"taken": "Boolean"
}
```

**Dosage**
```json
{
	"intake_moment": "Time: HH-MM-SS",
	"amount": "Integer",
	"medicine": "JSON Medicine"	
}
```

**Medicine**
```json
{
	"name": "String"
}
```
### Creating
#### Adding a new dosage for a specified patient
| **Type**      | PUT                                   |
|---------------|---------------------------------------|
| **Path**      | `/api/accounts/patients/{id}/dosages` |
| **Variables** | id: Integer                           |
| **Body**      | JSON Dosage                           |
| **Return**    | -                                     |
### Retrieving
#### Retrieving all scheduled dosages from a patient in an inverval
| **Type**      | GET                                   |
|---------------|---------------------------------------|
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
| **Type**      | PUT                   |
|---------------|-----------------------|
| **Path**      | `/api/general/videos` |
| **Variables** | -                     |
| **Body**      | JSON Video            |
| **Return**    | -                     |
### Retrieving
#### Retrieving all videos with a certain topic
| **Type**      | GET                           |
|---------------|-------------------------------|
| **Path**      | `/api/general/videos/{topic}` |
| **Variables** | topic: String                 |
| **Body**      | -                             |
