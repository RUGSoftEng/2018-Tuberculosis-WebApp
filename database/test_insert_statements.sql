insert into Accounts(name, username, pass_hash, role, api_token) values(
       "John",
       "Jboy",
       "passhash",
       "tester",
       "00token00"
);

insert into Physicians(id, email, token) values(
       1,
       "john@email.com",
       "11tokn"
);

insert into Accounts(name, username, pass_hash, role, api_token) values(
       "Jerome",
       "Jthesecond",
       "passhasher",
       "tester",
       "00tok"
);

insert into Patients(id, physician_id) values(
       2,
       1
);

insert into Medicines(med_name) values(
       "Highly experimental pills"
);

insert into Medicines(med_name) values(
       "Safe pills"
);

insert into Medicines(med_name) values(
       "Normal pills"
);

insert into Dosages(amount, patient_id, medicine_id, day, intake_time) values(
       12,
       2,
       1,
       ADDDATE(CURDATE(), 30),
       ADDTIME(CURTIME(), "00:10:00")
);

insert into Dosages(amount, patient_id, medicine_id, day, intake_time) values(
       3,
       2,
       2,
       ADDDATE(CURDATE(), 10),
       ADDTIME(CURTIME(), "01:09:00")
);
insert into Dosages(amount, patient_id, medicine_id, day, intake_time) values(
       12,
       2,
       1,
       ADDDATE(CURDATE(), 60),
       ADDTIME(CURTIME(), "00:10:00")
);

insert into Dosages(amount, patient_id, medicine_id, day, intake_time) values(
       1,
       2,
       3,
       ADDDATE(CURDATE(), 50),
       ADDTIME(CURTIME(), "10:10:00")
);

insert into Notes(patient_Id, question, day) values(
       2, "This is a test note", "2018-04-10"
);

insert into Notes(patient_Id, question, day) values(
       2, "This is another test note", "2018-04-11"
);


insert into Notes(patient_Id, question, day) values(
       2, "A NOTE", "2018-02-20"
);

insert into Videos(topic, title, reference) values(
       "Go Tutorials",
       "Golang REST API With Mux",
       "https://www.youtube.com/watch?v=SonwZ6MF5BE"
);

insert into Videos(topic, title, reference) values(
       "Go Tutorials",
       "Introduction - Go Lang Practical Programming Tutorial p.1",
       "https://www.youtube.com/watch?v=G3PvTWRIhZA"
);

insert into Videos(topic, title, reference) values(
       "Linux",
       "The Complete Linux Course: Beginner to Power User",
       "https://www.youtube.com/watch?v=wBp0Rb-ZJak"
);
