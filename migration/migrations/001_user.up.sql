CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `login` varchar(50) NOT NULL,
  `fullname` varchar(50) NOT NULL,
  `email` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE TABLE `user_login` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user` int NOT NULL,
  `password` varchar(50) NOT NULL,
  `session` varchar(50) NOT NULL,
  `session_date` int NOT NULL,
  `enabled` int NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE TABLE `user_permission` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user` int NOT NULL,
  `id_admin` int NOT NULL,
  `id_editor` int NOT NULL,
  `id_user` int NOT NULL,
  PRIMARY KEY (`id`)
);