CREATE TABLE `category` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `alias` varchar(50) NOT NULL,
  `auth` int NOT NULL,
  `permission` varchar(250) NOT NULL,
  PRIMARY KEY (`id`)
);
CREATE TABLE `blog` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category` int NOT NULL,
  `name` varchar(250) NOT NULL,
  `author` int NOT NULL,
  `preview` text NOT NULL,
  `text` text NOT NULL,
  `datecreated` datetime(6) NOT NULL,
  `datechanged` datetime(6) NOT NULL,
  PRIMARY KEY (`id`)
);
INSERT INTO `menu` (`parent`, `menu`, `name`, `url`, `icon`, `permission`, `ordering`, `enabled`) VALUES 
('0', 'single', 'Blog', '/blog/', 'edit', '', '3', '1');
INSERT INTO `category` VALUES 
(1, 'Category 1', 'category1', 0, ''),
(2, 'Category 2', 'category2', 0, ''),
(3, 'Category 3', 'category3', 1, '');
INSERT INTO `blog` VALUES 
(1, 1, 'Record 1', 2, '<p>Preview 1</p>', '<p>Text 1</p>', '2021-11-01 11:11:11.1', '2021-11-01 11:11:11.1'),
(2, 2, 'Record 2', 2, '<p>Preview 2</p>', '<p>Text 2</p>', '2021-11-02 12:11:11.1', '2021-11-02 12:11:11.1'),
(3, 3, 'Record 3', 1, '<p>Preview 3</p>', '<p>Text 3</p>', '2021-11-03 13:11:11.1', '2021-11-03 13:11:11.1');
