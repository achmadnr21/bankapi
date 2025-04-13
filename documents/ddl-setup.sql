SET TIME ZONE 'Asia/Jakarta';

create schema auth;
create schema profile;
create schema account;
create schema transfer;

drop table transfer.transactions;
drop table transfer.transaction_fees;
drop table account.accounts;
drop table profile.users;
drop table profile.user_roles;
drop table account.branches;
drop table account.account_types;
drop table account.currencies;

create table profile.user_roles(
	id char(3) unique primary key,
	name varchar(255) not null,
	can_add_user boolean default false,
	can_add_branches boolean default false,
	can_add_account_types boolean default false,
	can_add_currencies boolean default false,
	can_add_accounts boolean default false,
	created_at  TIMESTAMP DEFAULT NOW()
);

insert into profile.user_roles(id, name, can_add_user, can_add_branches, can_add_account_types,  can_add_currencies, can_add_accounts )
values
('OWN', 'Owner', true, true, true, true, true),
('MNG', 'Manager', true, true, true, true, true),
('CSO', 'Customer Service Officer', true, false, false, false, true),
('USR', 'User', false, false, false, false, false);
select * from profile.user_roles;     

create table profile.users(
	id uuid unique default gen_random_uuid() primary key,
	employee_id  UUID null,
	nik varchar(16) unique not null,
	role_id char(3) not null default 'USR',
	full_name   VARCHAR(255) NOT NULL,
	username	VARCHAR(64) unique not null,
    email       VARCHAR(255) UNIQUE NOT NULL,
    phone       VARCHAR(20) UNIQUE NOT NULL,
    password_hash    TEXT NOT NULL,  -- Disimpan dalam bentuk hash
    created_at  TIMESTAMP DEFAULT NOW(),
    foreign key (role_id) references profile.user_roles(id),
    foreign key (employee_id) references profile.users(id)
);

insert into profile.users(nik, role_id, full_name, username, email, phone, password_hash)
values('3522150101010004','MNG', 'admin manager','admin123', 'rudy@gmail.com', '6285234123124', '$2a$12$TFe.3.s5X/EdhktJDEud3uNpHMqjO6NgdGNDRUhuc0dvkTPuKC5VK' );

select * from profile.users;

SELECT id, employee_id, role_id, full_name, username, email, phone, password_hash, created_at FROM profile.users WHERE username = 'achmadnr9';

create table account.branches(
	id SERIAL unique not null primary key,
	name varchar(255) not null,
	address text not null,
	created_at  TIMESTAMP DEFAULT NOW()
);
--delete from account.branches;
insert into account.branches(id, name, address)
values
(1, 'Jakarta Branch', 'Gedung BRI 1, Jl. Jendral Sudirman Kav. 44-46, Jakarta'),
(11, 'Bojonegoro', 'Jl. D. I. Panjaitan');
SELECT setval(pg_get_serial_sequence('account.branches', 'id'), (SELECT MAX(id) FROM account.branches));
select * from account.branches b ;


create table account.account_types(
	id SERIAL unique not null primary key,
	name varchar(255) not null,
	can_etoll boolean default false,
	can_indomart boolean default false,
	can_steam boolean default false,
	created_at  TIMESTAMP DEFAULT NOW()
);

insert into account.account_types(id, name, can_etoll, can_indomart, can_steam)
values
(110, 'Tabungan BriTama', true, true, true);
SELECT setval(pg_get_serial_sequence('account.account_types', 'id'), (SELECT MAX(id) FROM account.account_types));
select * from account.account_types;


create table account.currencies(
	id 			char(3) unique not null primary key,
	name 		varchar(255) not null,
	rate_in_idr DECIMAL(15, 2) not null,
	created_at  TIMESTAMP DEFAULT NOW()
);

insert into account.currencies(id, name, rate_in_idr)
values
('IDR', 'Indonesia Rupiah', 1.0);
select * from account.currencies;

-- TABEL SEQUENCE PER KOMB CABANG + TIPE
CREATE TABLE account.account_sequences (
    branch_id INTEGER NOT NULL,
    account_type_id INTEGER NOT NULL,
    last_sequence INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (branch_id, account_type_id),
    FOREIGN KEY (branch_id) REFERENCES account.branches(id),
    FOREIGN KEY (account_type_id) REFERENCES account.account_types(id)
);

CREATE TABLE account.accounts (
    id           	UUID unique PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      	UUID NOT null,
    employee_id  	UUID null,
    branch_id    	integer not null,
    account_type_id integer not null,
    sequence_number integer not null,
    account_number 	VARCHAR(20) UNIQUE NOT NULL,
    pin_hash 		varchar(255) not null,
    currency_id     VARCHAR(3) NOT NULL DEFAULT 'IDR',
    balance      	DECIMAL(15,2) DEFAULT 0.00,
    created_at   	TIMESTAMP DEFAULT NOW(),
    foreign key (user_id) references profile.users(id) ON DELETE CASCADE,
    foreign key (employee_id) references profile.users(id),
    foreign key (branch_id) references account.branches(id),
    foreign key (account_type_id) references account.account_types(id),
    foreign key (currency_id) references account.currencies(id)
);

select * from account.accounts;

--schema transfer and stuff 
drop type transfer.status;
create TYPE transfer.status AS ENUM ('pending', 'completed', 'failed');
drop type transfer.transaction_type;
create TYPE transfer.transaction_type AS ENUM ('deposit', 'transfer', 'withdraw');

create table transfer.transaction_fees(
	id char(3) unique primary key,
	name varchar(255) not null,
	fee decimal(15,2) default 0,
	created_at TIMESTAMP DEFAULT NOW()
);
insert into transfer.transaction_fees(id, name, fee)
values
('INT', 'Same Bank TF', 0.0),
('BIF', 'BI-Fast', 2500.0);
select * from transfer.transaction_fees;

create table transfer.transactions(
	id UUID unique PRIMARY KEY DEFAULT gen_random_uuid(),
	transaction_type transfer.transaction_type not null,
	transaction_fee char(3) not null,
	account_from uuid not null,
	account_to uuid not null,
	ref_no varchar(20),
	currency_id char(3) not null,
	amount decimal(15,2) not null,
	status transfer.status default 'pending',
	created_at TIMESTAMP DEFAULT NOW(),
	foreign key (transaction_fee) references transfer.transaction_fees(id),
	foreign key (account_from) references account.accounts(id),
	foreign key (account_to) references account.accounts(id),
	foreign key (currency_id) references account.currencies(id)
);



