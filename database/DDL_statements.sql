CREATE TABLE account(

  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,

  name VARCHAR(200) NOT NULL,

  username VARCHAR(200) NOT NULL,

  pass_hash VARCHAR(200) NOT NULL,

  role VARCHAR(10) NOT NULL, 

  api_token VARCHAR(64)

);

CREATE TABLE physician(

  id INT NOT NULL PRIMARY KEY,

  email VARCHAR(200) NOT NULL,

  token VARCHAR(8) NOT NULL,

  FOREIGN KEY (id) REFERENCES account(id)

);

CREATE TABLE patient(

  id INT NOT NULL PRIMARY KEY,

  physician_id INT NOT NULL,

  FOREIGN KEY (id) REFERENCES account(id),

  FOREIGN KEY (physician_id) REFERENCES physician(id)

);

CREATE TABLE note(

  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,

  patient_Id INT NOT nULL, 

  question TEXT NOT NULL,

  day DATE NOT NULL,

  FOREIGN KEY (patient_Id) REFERENCES patient(id)

);

CREATE TABLE medicine(

   id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,

   med_name VARCHAR(50) NOT NULL	

);

CREATE TABLE dosage(

  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  
  amount TINYINT NOT NULL,

  patient_id INT NOT NULL,

  medicine_id INT NOT NULL,

  day DATE NOT NULL,

  intake_time TIME NOT NULL,

  FOREIGN KEY (patient_id) REFERENCES patient(id),

  FOREIGN KEY (medicine_id) REFERENCES medicine(id)

);

CREATE TABLE movie(
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE faq(
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  question TEXT NOT NULL 
);
