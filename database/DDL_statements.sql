CREATE TABLE Accounts(
       id		INT		NOT NULL PRIMARY KEY AUTO_INCREMENT,
       name		VARCHAR(200) 	NOT NULL,
       username		VARCHAR(200) 	NOT NULL  UNIQUE,
       pass_hash	VARCHAR(200) 	NOT NULL,
       role		VARCHAR(10) 	NOT NULL, 
       api_token	VARCHAR(64)
);

CREATE TABLE Physicians(
       id    INT		NOT NULL PRIMARY KEY,
       email VARCHAR(200) 	NOT NULL,
       token VARCHAR(8) 	NOT NULL,
       FOREIGN KEY (id) REFERENCES Accounts(id)
);

CREATE TABLE Patients(
       id		INT NOT NULL PRIMARY KEY,
       physician_id  	INT NOT NULL,
       FOREIGN KEY (id)			REFERENCES Accounts(id),
       FOREIGN KEY (physician_id) 	REFERENCES Physicians(id)
);

CREATE TABLE Notes(
       id		INT	NOT NULL PRIMARY KEY AUTO_INCREMENT,
       patient_id    	INT 	NOT NULL,
       question		TEXT 	NOT NULL,
       day		DATE 	NOT NULL,
       FOREIGN KEY (patient_Id) REFERENCES Patients(id)
);

CREATE TABLE Medicines(
       id	INT		NOT NULL PRIMARY KEY AUTO_INCREMENT,
       med_name	VARCHAR(50) 	NOT NULL	
);

CREATE TABLE Dosages(
       id		INT	NOT NULL PRIMARY KEY AUTO_INCREMENT,
       patient_id	INT	NOT NULL,
       medicine_id	INT	NOT NULL,
       amount		TINYINT	NOT NULL,
       intake_time	TIME	NOT NULL,
       FOREIGN KEY (patient_id)     REFERENCES Patients(id),
       FOREIGN KEY (medicine_id)    REFERENCES Medicines(id)
);

CREATE TABLE ScheduledDosages(
       dosage		INT	NOT NULL,
       day		DATE 	NOT NULL,
       taken		BOOLEAN NOT NULL,
       PRIMARY KEY (dosage, day)
);

CREATE TABLE Videos(
       id		INT		NOT NULL PRIMARY KEY AUTO_INCREMENT,
       topic		VARCHAR(255) 	NOT NULL,
       title		VARCHAR(255) 	NOT NULL,
       reference	VARCHAR(255) 	NOT NULL
);

CREATE TABLE Quizzes(
       video		INT		NOT NULL,
       question		VARCHAR(255)	NOT NULL,
       answers		TEXT		NOT NULL,
       FOREIGN KEY (video) REFERENCES Videos(id),
       PRIMARY KEY (video, question)
);

CREATE TABLE FAQ(
       id	INT	NOT NULL PRIMARY KEY AUTO_INCREMENT,
       question TEXT 	NOT NULL,
       answer	TEXT 	NOT NULL 
);

