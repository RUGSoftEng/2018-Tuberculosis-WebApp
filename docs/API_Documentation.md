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
### Creating
#### Adding a note for a specified patient
**PUT** `/api/accounts/patients/{id}/notes` 

Include the note (in JSON) in the body of the request:
```json
{
	"note": "String: The message of the note",
	"created_at": "Date: YYYY-MM-DD"
}
```
### Retrieving
#### Retrieving all notes from a specified patient
**GET** `/api/accounts/patients/{id}/notes`

##### Returns:
List of notes in JSON format:
```json
{
	"note": "String: The message of the note",
	"created_at": "Date: YYYY-MM-DD"
}
```

### Updating
_To be added_
### Deleting
_To be added_
## Dosages
### Creating
_To be documented_
### Retrieving
#### Retrieving all dosages from a patient in an inverval
**GET** `/api/account/patients/{id}/dosages`

You also need to query the interval with the variables:
* `from=YYYY-MM-DD`
* `until=YYYY-MM-DD`

##### Returns
List of dosages in JSON format:
```json
{
	"date": "Date: YYYY-MM-DD",
	"intake_moment": "Time: HH-MM-SS",
	"amount": "Integer",
	"medicine": "String: Name of the medicine which needs to be taken",
	"taken": "Boolean"
}
```
### Updating
_To be added_
### Deleting
_To be added_
## Videos
### Creating
#### Adding a new video
**PUT** `/api/general/videos`

Include the video (in JSON) in the body of the request:
```json
{
	"topic": "String",
	"title": "String",
	"reference": "String: url"
}
```
### Retrieving
#### Retrieving all videos with a certain topic
**GET** `/api/general/videos/{topic}

##### Returns:
List of videos in JSON format:
```json
{
	"topic": "String",
	"title": "String",
	"reference": "String: url"
}
```
