CREATE TABLE `user` (
	`id`  INT(11) NOT NULL AUTO_INCREMENT COMMENT 'name of Hello',
	`name`  VARCHAR(128) NOT NULL DEFAULT "" COMMENT 'name of Hello',
	`age`  INT(11) NULL DEFAULT "0" COMMENT 'age of Hello',
	`sex`  TINYINT(1) UNSIGNED NOT NULL DEFAULT "0" COMMENT '',
	`foo_bar`  INT(11) NULL DEFAULT "0" COMMENT 'fooBar',
	`create_at`  BIGINT(20) NOT NULL DEFAULT "0" COMMENT 'create_at',
	`update_at`  BIGINT(20) NULL DEFAULT "0" COMMENT 'update_at',
    UNIQUE KEY `UserNameUK`(`name`),
    PRIMARY KEY(`id`)
) ENGINE=innodb AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT '';
CREATE INDEX `UserAgeSexIDX` ON `user` (`age`,`sex`);