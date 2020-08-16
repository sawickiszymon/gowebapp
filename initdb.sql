CREATE KEYSPACE IF NOT EXISTS cass WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};
USE cass;
-- THIS KEYSPACE SHOULD MATCH cassandraEnv file keyspace for it to use the same table as gocql
CREATE TABLE IF NOT EXISTS Email
(
    email        text,
    title        text,
    content      text,
    magic_number int,
    PRIMARY KEY (email, magic_number, content)
);CREATE INDEX  IF NOT EXISTS mNumber_index ON Email(magic_number);


INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 10);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 11);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 12);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 13);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 14);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 15);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 16);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 17);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 18);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 19);
INSERT INTO Email (email, title, content, magic_number) VALUES ('sz.sawicki1@gmail.com', 'what', 'ever', 20);