insert into workgroups (name) values ('Managers and Non-billable - ANZ');
insert into workgroups (name) values ('Architects - ANZ');
insert into workgroups (name) values ('Retainer - ANZ');
insert into workgroups (name) values ('PMO - ANZ');
insert into workgroups (name) values ('Server - Australia');

insert into locations (name) values ('ACT');
insert into locations (name) values ('NSW');
insert into locations (name) values ('QLD-NT');
insert into locations (name) values ('SA');
insert into locations (name) values ('WA');
insert into locations (name) values ('VIC-TAS');
insert into locations (name) values ('NZL');

insert into job_titles (title) values ('Director');
insert into job_titles (title) values ('Senior Manager');
insert into job_titles (title) values ('Manager');
insert into job_titles (title) values ('Senior Project Manager');
insert into job_titles (title) values ('Project Manager');
insert into job_titles (title) values ('Associate Project Manager II');
insert into job_titles (title) values ('Associate Project Manager I');
insert into job_titles (title) values ('Staff Consulting Architect');
insert into job_titles (title) values ('Staff Consultant');
insert into job_titles (title) values ('Consulting Architect');
insert into job_titles (title) values ('Senior Consultant');
insert into job_titles (title) values ('Consultant');
insert into job_titles (title) values ('Associate Consultant II');
insert into job_titles (title) values ('Associate Consultant I');

insert into employment_types (type) values ('Full time');
insert into employment_types (type) values ('Contractor');
insert into employment_types (type) values ('Sub-Contractor');

insert into resources (id, name, email, job_title_id, workgroup_id, location_id, type_id, manager_id)
values (123456, 'Jane Doe', 'jane@doe.com',
  (select t.id from job_titles t where title = 'Director'),
  (select w.id from workgroups w where w.name = 'Managers and Non-billable - ANZ'),
  (select l.id from locations l where l.name = 'ACT'),
  (select typ.id from employment_types typ where typ.type = 'Full time'),
  123456
);