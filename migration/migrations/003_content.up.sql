INSERT INTO `user` VALUES 
(1,'admin','Admin','admin@domain.local'),
(2,'editor','Editor','editor@domain.local'),
(3,'user','User','user@domain.local');
INSERT INTO `user_login` VALUES 
(1,1,'3c36418d496a3112db5ed34cafabf266','',1,1),
(2,2,'3c36418d496a3112db5ed34cafabf266','',1,1),
(3,3,'3c36418d496a3112db5ed34cafabf266','',1,1);
INSERT INTO `user_permission` VALUES 
(1,1,1,1,1),
(2,2,0,1,0),
(3,3,0,0,1);
INSERT INTO `menu` VALUES 
(1,0,'single','Home','/','home','',1,1),
(2,0,'multi','Pages','','list','id_auth',11,1),
(3,2,'','UserInfo','/user/info','id-card-o','id_auth',1,1),
(4,2,'','User Page','/user/user','lock','id_user',4,1),
(5,2,'','Editor Page','/user/editor','lock','id_editor',7,1),
(6,0,'multi','System','','gears','id_admin',31,1),
(7,6,'','Users','/system/users','users','id_admin',1,1),
(8,6,'','Menu','/system/menu','list','id_admin',4,1);
