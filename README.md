# goRDStest

--
CREATE DATABASE test DEFAULT CHARACTER SET utf8;
create table test.member(id int primary key, name varchar(8),password varchar(100));
ALTER TABLE test.member MODIFY id INT AUTO_INCREMENT;
--
