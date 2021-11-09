CREATE TABLE `menu` (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent` int NOT NULL,
  `menu` varchar(30) NOT NULL,
  `name` varchar(30) NOT NULL,
  `url` varchar(30) NOT NULL,
  `icon` varchar(30) NOT NULL,
  `permission` varchar(30) NOT NULL,
  `ordering` int NOT NULL,
  `enabled` int NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
);