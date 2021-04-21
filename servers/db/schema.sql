create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(320) not null unique,
    pass_hash varbinary(72) not null,
    username varchar(255) not null unique,
    first_name varchar(64) not null,
    last_name varchar(128) not null,
    photo_url varchar(128) not null
);

create table if not exists userLog (
    id int not null auto_increment primary key,
    userID int not null,
    inTime datetime not null,
    clientIP varchar(15) not null,
    foreign key (userID) references users(id)
)