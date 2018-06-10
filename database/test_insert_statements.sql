insert into Accounts(name, username, pass_hash, role, api_token) values(
       'John',
       'Jboy',
       'passhash',
       'physician',
       '00token00'
);

insert into Physicians(id, email, token) values(
       1,
       'john@email.com',
       '11tokn'
);

insert into Accounts(name, username, pass_hash, role, api_token) values(
       'Patient',
       'Username',
       '$2a$14$ILMXoKunbxlXmXco12NNi.MGNaIFL6nSDj25XXs7bpMxhLfCnjtfW',
       'patient',
       ''
);

insert into Patients(id, physician_id) values(
       2,
       1
);

insert into Medicines(med_name) values(
       'Rifampicin'
);

insert into Medicines(med_name) values(
       'Isoniazid'
);

insert into Medicines(med_name) values(
       'Pyrazinamide'
);

insert into Dosages(patient_id, medicine_id, amount, intake_interval_start, intake_interval_end) values (
       2, 1, 12, ADDTIME(CURTIME(), '00:10:00'), ADDTIME(CURTIME(), '02:10:00')
);

insert into Dosages(patient_id, medicine_id, amount, intake_interval_start, intake_interval_end) values (
       2, 2, 3, ADDTIME(CURTIME(), '01:09:00'), ADDTIME(CURTIME(), '02:09:00')
);

insert into Dosages(patient_id, medicine_id, amount, intake_interval_start, intake_interval_end) values (
       2, 3, 6, ADDTIME(CURTIME(), '10:10:00'), ADDTIME(CURTIME(), '11:10:00')
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
       2, 'This is a test note', '2018-04-10'
);

insert into Notes(patient_Id, question, day) values(
       2, 'This is another test note', '2018-04-11'
);


insert into Notes(patient_Id, question, day) values(
       2, 'A NOTE', '2018-02-20'
);
