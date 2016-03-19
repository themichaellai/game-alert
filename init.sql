create table games (id integer primary key, home string, away string, status string, positionId string);
create index positionId_index on games (positionId);
create table events (id integer primary key, gameId integer, datetime integer, status string);
