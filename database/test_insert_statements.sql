insert into account(name, username, pass_hash, role, api_token) values(
       "John",
       "Jboy",
       "passhash",
       "tester",
       "00token00"
);

insert into physician(id, email, token) values(
       1,
       "john@email.com",
       "11tokn"
);

insert into account(name, username, pass_hash, role, api_token) values(
       "Jerome",
       "Jthesecond",
       "passhasher",
       "tester",
       "00tok"
);

insert into patient(id, physician_id) values(
       2,
       1
);

insert into medicine(med_name) values(
       "Highly experimental pills"
);

insert into medicine(med_name) values(
       "Safe pills"
);

insert into medicine(med_name) values(
       "Normal pills"
);

insert into dosage(amount, patient_id, medicine_id, day, intake_time) values(
       12,
       2,
       1,
       ADDDATE(CURDATE(), 30),
       ADDTIME(CURTIME(), "00:10:00")
);

insert into dosage(amount, patient_id, medicine_id, day, intake_time) values(
       3,
       2,
       2,
       ADDDATE(CURDATE(), 10),
       ADDTIME(CURTIME(), "01:09:00")
);
insert into dosage(amount, patient_id, medicine_id, day, intake_time) values(
       12,
       2,
       1,
       ADDDATE(CURDATE(), 60),
       ADDTIME(CURTIME(), "00:10:00")
);

insert into dosage(amount, patient_id, medicine_id, day, intake_time) values(
       1,
       2,
       3,
       ADDDATE(CURDATE(), 50),
       ADDTIME(CURTIME(), "10:10:00")
);

insert into note(patient_Id, question, day) values(
       2, "This is a test note", "2018-04-10"
);

insert into note(patient_Id, question, day) values(
       2, "This is another test note", "2018-04-11"
);


insert into note(patient_Id, question, day) values(
       2, "A NOTE", "2018-02-20"
);
