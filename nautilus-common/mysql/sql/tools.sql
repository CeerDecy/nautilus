CREATE TABLE `ai_tools`
(
    id          bigint primary key auto_increment COMMENT 'id',
    name        varchar(255) not null COMMENT 'tool name',
    description text         not null COMMENT 'tool description',
    strict      char(1) default 0 COMMENT 'strict',
    parameters  text         not null COMMENT 'json parameters'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='ai tools';