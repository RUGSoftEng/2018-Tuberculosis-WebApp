insert into Accounts(name, username, pass_hash, role, api_token) values(
       "John",
       "Jboy",
       "passhash",
       "physician",
       "00token00"
);

insert into Physicians(id, email, token) values(
       1,
       "john@email.com",
       "11tokn"
);

insert into Accounts(name, username, pass_hash, role, api_token) values(
       "Patient",
       "Username",
       "$2a$14$ILMXoKunbxlXmXco12NNi.MGNaIFL6nSDj25XXs7bpMxhLfCnjtfW",
       "patient",
       ""
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

insert into Dosages(patient_id, medicine_id, amount, intake_time) values (
       2, 1, 12, ADDTIME(CURTIME(), "00:10:00")
);

insert into Dosages(patient_id, medicine_id, amount, intake_time) values (
       2, 2, 3, ADDTIME(CURTIME(), "01:09:00")
);

insert into Dosages(patient_id, medicine_id, amount, intake_time) values (
       2, 3, 6, ADDTIME(CURTIME(), "10:10:00")
);

insert into ScheduledDosages (dosage, day, taken) values(
       1, ADDDATE(CURDATE(), 30), FALSE
);

insert into ScheduledDosages (dosage, day, taken) values(
       1, ADDDATE(CURDATE(), 60), FALSE
);

insert into ScheduledDosages (dosage, day, taken) values(
       2, ADDDATE(CURDATE(), 10), FALSE
);

insert into ScheduledDosages (dosage, day, taken) values(
       2, ADDDATE(CURDATE(), 20), FALSE
);

insert into ScheduledDosages (dosage, day, taken) values(
       3, ADDDATE(CURDATE(), 50), FALSE
);

insert into ScheduledDosages (dosage, day, taken) values(
       3, ADDDATE(CURDATE(), 1), TRUE
);

insert into ScheduledDosages (dosage, day, taken) values(
       3, ADDDATE(CURDATE(), 0), TRUE
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
       "GoTutorials",
       "Golang REST API With Mux",
       "https://www.youtube.com/watch?v=SonwZ6MF5BE"
);

insert into Videos(topic, title, reference) values(
       "GoTutorials",
       "Introduction - Go Lang Practical Programming Tutorial p.1",
       "https://www.youtube.com/watch?v=G3PvTWRIhZA"
);

insert into Videos(topic, title, reference) values(
       "Linux",
       "The Complete Linux Course: Beginner to Power User",
       "https://www.youtube.com/watch?v=wBp0Rb-ZJak"
);

insert into Videos(topic, title, reference) values(
       "Tuberculosis",
       "What happens to me and my body when I have tuberculosis?",
       "https://youtu.be/fXiXGRlvH70"
);

insert into Quizzes(video, question, answers) values(
       4,
       "Tuberculosis is caused by",
       "A bacteria:A virus:Fungi"
);

insert into FAQ(question, answer) values (
       "What is tuberculosis?",
       "Tuberculosis is a disease, caused by the bacterium called Mycobacterium tuberculosis. It can cause lung problems, like a dry cough, but it can also manifest with other symptoms."
);

insert into FAQ(question, answer) values (
       "How do you get tuberculosis?",
       "Tuberculosis is spread via the air. You can get it for example. when someone coughs or talks to you. The bacteria are present in small droplets in the air and you could inhale those droplets."
);
