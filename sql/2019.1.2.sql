CREATE TABLE `qp`.`signconfig` (
  `id` INT NOT NULL,
  `icon` INT NOT NULL,
  `money` INT NOT NULL,
  PRIMARY KEY (`id`));
  
  CREATE TABLE `qp`.`sign` (
  `uid` BIGINT NOT NULL,
  `index` INT NOT NULL,
  `time` BIGINT NOT NULL,
  PRIMARY KEY (`uid`));

INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('1', '1', '100');
INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('2', '1', '200');
INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('3', '2', '300');
INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('4', '2', '400');
INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('5', '3', '600');
INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('6', '3', '800');
INSERT INTO `qp`.`signconfig` (`id`, `icon`, `money`) VALUES ('7', '4', '1000');