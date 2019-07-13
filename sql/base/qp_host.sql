CREATE DATABASE  IF NOT EXISTS `qp_host` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `qp_host`;
-- MySQL dump 10.13  Distrib 5.6.17, for osx10.6 (i386)
--
-- Host: 103.53.124.238    Database: qp_host
-- ------------------------------------------------------
-- Server version	5.7.25-0ubuntu0.16.04.2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admin_auth_rule`
--

DROP TABLE IF EXISTS `admin_auth_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `admin_auth_rule` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '方法',
  `title` varchar(45) NOT NULL COMMENT '方法描述',
  `type` tinyint(1) NOT NULL DEFAULT '1',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `condition` char(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=176 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_auth_rule`
--

LOCK TABLES `admin_auth_rule` WRITE;
/*!40000 ALTER TABLE `admin_auth_rule` DISABLE KEYS */;
INSERT INTO `admin_auth_rule` VALUES (1,'Index/index','首页',0,1,'1'),(2,'Orders/index','订单列表',1,1,'1'),(3,'Orders/index2','已支付订单',1,1,'1'),(4,'Orders/index3','未支付订单列表',1,1,'1'),(5,'Agsale/index','代理售卡列表',1,1,'1'),(6,'Adsale/index','管理员售卡列表',1,1,'1'),(7,'Tcash/index','所有列表',2,1,'1'),(8,'Tcash/index2','申请列表',2,1,'1'),(9,'Tcash/index3','已处理列表',2,1,'1'),(10,'Tcash/index4','已拒绝列表',2,1,'1'),(11,'Commission/index','提成价格编辑',2,1,'1'),(12,'Product/index','产品列表',3,1,'1'),(13,'Buycard/index','售卡列表',3,1,'1'),(14,'Sendag/index','发给代理',3,1,'1'),(15,'Sendpaly/index','发给玩家',3,1,'1'),(16,'Agent/index','代理列表',3,1,'1'),(17,'Atype/index','代理类别列表',4,1,'1'),(18,'Openter/index','跑马灯',6,1,'1'),(19,'Play/index','玩家列表',3,1,'1'),(20,'Admin/index','后台用户管理',9,1,'1'),(21,'Auth/index','后台权限管理',9,1,'1'),(22,'Commission/index','提成价格查看',2,1,'1'),(23,'Product/add','产品添加',3,1,'1'),(24,'Product/edit','产品修改',3,1,'1'),(25,'Product/del','产品删除',3,1,'1'),(26,'Buycard/add','售卡列表添加',3,1,'1'),(27,'Buycard/edit','售卡列表修改',3,1,'1'),(28,'Buycard/del','售卡列表删除',3,1,'1'),(29,'Sendag/card','代理加卡',3,1,'1'),(30,'Sendag/card2','代理惩罚',3,1,'1'),(31,'Sendpaly/card','玩家加卡',3,1,'1'),(32,'Sendpaly/card2','玩家惩罚',3,1,'1'),(33,'Agent/add','代理添加',3,1,'1'),(34,'Agent/del','代理删除',3,1,'1'),(35,'Agent/edit','代理编辑',3,1,'1'),(36,'Sendpaly/faka','玩家确定加卡',3,1,'1'),(37,'Sendag/faka','代理确定加卡',3,1,'1'),(38,'Sendpaly/faka2','玩家确定惩罚',3,1,'1'),(39,'Sendag/faka2','玩家确定惩罚',3,1,'1'),(40,'Atype/add','代理类别添加',1,1,'1'),(41,'Atype/edit','代理类别修改',1,1,'1'),(42,'Atype/del','代理类别删除',1,1,'1'),(43,'Openter/edit','跑马灯编辑',1,1,'1'),(44,'Admin/add','后台用户添加',9,1,'1'),(45,'Admin/edit','后台用户编辑',9,1,'1'),(46,'Admin/del','后台用户删除',9,1,'1'),(47,'Auth/add','后台权限添加',9,1,'1'),(48,'Auth/edit','后台权限修改',9,1,'1'),(49,'Auth/del','后台权限删除',1,1,'1'),(50,'Sendag/card3','代理转移',1,1,'1'),(51,'Tcash/tixian','用户提现',1,1,'1'),(52,'Tcash/tx_money','确认用户提现',1,1,'1'),(53,'Tcash/refuse','拒绝用户提现',1,1,'1'),(54,'Agent/password','修改代理密码',1,1,'1'),(55,'Agent/index2','查看代理的下级个数',1,1,'1'),(56,'Commission/edit','提成价格编辑',1,1,'1'),(57,'Notice/index','代理公告列表',9,1,'1'),(58,'Notice/add','代理公告新增',9,1,'1'),(59,'Notice/edit','代理公告编辑',9,1,'1'),(60,'Notice/del','代理公告删除',9,1,'1'),(61,'Notice/bobao','代理公告播报/停用',9,1,'1'),(62,'Index/index2','1',1,1,'1'),(63,'Xiaoshou/index','销售记录',1,1,'1'),(64,'Match/index','游戏比赛系统',1,1,'1'),(65,'Match/add','游戏比赛系统添加',1,1,'1'),(66,'Match/edit','游戏比赛系统修改',1,1,'1'),(67,'Match/del','游戏比赛系统删除',1,1,'1'),(68,'UserMake/index','用户构成',1,1,'1'),(69,'Upload/index','图片轮播首页',10,1,'1'),(70,'Upload/add','图片添加',10,1,'1'),(71,'Upload/del','图片删除',10,1,'1'),(72,'Api/getUrl','获取轮播图',1,1,'1'),(73,'Game/roomDismiss','显示解散房间',6,1,'1'),(74,'Game/doRoomDismiss','执行房间解散',1,1,'1'),(75,'Game/searchPlayer','查找玩家状态',1,1,'1'),(76,'Agent/modCommission','修改代理提成比例',1,1,'1'),(77,'TCash/ManualTCash','手动提现',2,1,'1'),(78,'TCash/searchAgent','查询代理信息',1,1,'1'),(79,'TCash/doTcash','执行提现操作',1,1,'1'),(80,'Agent/agentAuth','授权为代理',4,1,'1'),(81,'Agent/searchAgent','查看代理授权状态',1,1,'1'),(82,'Agent/moveAgent','查看转移代理',4,1,'1'),(83,'Agent/doMoveAgent','执行代理转移',1,1,'1'),(84,'Play/goldlogs','玩家金币记录',1,1,'1'),(85,'Agent/resetPass','重置提现密码',1,1,'1'),(86,'Play/addBlack','封禁玩家',1,1,'1'),(87,'Play/removeBlack','解封玩家',1,1,'1'),(88,'Play/searchGoldLogs','查询玩家金币记录',7,1,'1'),(89,'Play/topupLogs','查询玩家充值到账记录',1,1,'1'),(90,'Orders/searchOrder','查找订单',1,1,'1'),(91,'Game/bzwdealwin','设置数值',6,1,'1'),(92,'Chart/index','报表首页',8,1,'1'),(93,'Chart/goldCost','金币消耗报表',8,1,'1'),(94,'Play/gift','显示给玩家发礼包',7,1,'1'),(95,'Play/getGift','执行给玩家发礼包',1,1,'1'),(96,'Index/goldLogs','显示金币统计',1,1,'1'),(97,'Index/cardLogs','显示房卡统计',1,1,'1'),(98,'Index/tradeLogs','显示交易统计',1,1,'1'),(99,'tcash/mtcashLogs','显示手动提现列表',2,1,'1'),(100,'Play/gameLogs','显示玩家战况',7,1,'1'),(101,'Play/searchGameLogs','查询玩家战况',7,1,'1'),(102,'Index/goldLogsTest','测试',1,1,'1'),(103,'Upload/updateOrder','更新宣传资料排序',10,1,'1'),(104,'Chart/addedUser','统计每日新增用户',8,1,'1'),(105,'SpecAgent/authSpec','授权特殊代理账号',5,1,'1'),(106,'SpecAgent/index','特殊代理列表',5,1,'1'),(107,'Chart/timeUser','统计时段在线人数',8,1,'1'),(108,'Chart/goldCostByGame','统计游戏金币消耗',8,1,'1'),(109,'Exchange/index','推广额兑换金币列表',1,1,'1'),(110,'Exchange/listData','推广额兑换金币数据',1,1,'1'),(111,'Play/searchWzq','查询五子棋上下分',7,1,'1'),(112,'Exchange/index2','金币兑换申请',1,1,'1'),(113,'Exchange/shenhe','金币兑换审核',1,1,'1'),(114,'Exchange/index3','金币提现',1,1,'1'),(115,'Exchange/gold','金币提现审核',1,1,'1'),(116,'Sendpaly/card3','玩家金币明细查看',1,1,'1'),(117,'Agent/account','官方认证帐号管理',1,1,'1'),(118,'Agent/account_add','官方认证帐号新增',1,1,'1'),(119,'Agent/account_del','官方认证帐号删除',1,1,'1'),(120,'Agent/account_edit','官方认证帐号编辑',1,1,'1'),(121,'SpecAgent/butedit','特殊代理操作',1,1,'1'),(122,'Play/guanli','特殊帐号管理',1,1,'1'),(123,'Game/bzwdealwin2','自定义数值管理',1,1,'1'),(124,'Robot/index','机器人管理',1,1,'1'),(125,'Robot/del','机器人删除',1,1,'1'),(126,'Robot/edit','机器人编辑',1,1,'1'),(127,'Robot/add','机器人新增',1,1,'1'),(128,'Robot2/index','机器人新增',1,1,'1'),(129,'Fish/index','捕鱼设置',1,1,'1'),(130,'Fish/index2','炮设置',1,1,'1'),(131,'Fish/index3','鱼设置',1,1,'1'),(132,'Fish/edit','捕鱼设置',1,1,'1'),(133,'Fish/edit2','炮设置',1,1,'1'),(134,'Fish/edit2','鱼设置',1,1,'1'),(135,'Game/bzwdealwin3','奖池设置',1,1,'1'),(136,'index/robot','机器人统计',1,1,'1'),(137,'SpecAgent/show','代理下级查看',1,1,'1'),(138,'Index/cardLogs2','系统抽水统计',1,1,'1'),(139,'Dolog/index','后台操作日志',1,1,'1'),(140,'Dolog/index2','后台登录日志',1,1,'1'),(141,'index/cardLogs3','系统盈利统计',1,1,'1'),(142,'Doname/index','域名设置',1,1,'1'),(143,'SpecAgent/butedit','特殊代理删除',1,1,'1'),(144,'Doname/add','ip白名单添加',1,1,'1'),(145,'Doname/edit','ip白名单修改',1,1,'1'),(146,'Doname/del','ip白名单删除',1,1,'1'),(147,'Game/switchs','金币游戏开关',1,1,'1'),(148,'Game/switchs2','房卡游戏开关',1,1,'1'),(149,'Xiaoshou/gold','金币场金额',1,1,'1'),(150,'Xiaoshou/card','房卡套餐',1,1,'1'),(151,'Xiaoshou/cft','成富通支付',1,1,'1'),(152,'Game/bzwdealwin4','筹码管理',1,1,'1'),(153,'ZhuanJin/index','赚金说明',1,1,'1'),(154,'ZhuanJin/add','赚金添加',1,1,'1'),(155,'ZhuanJin/edit','赚金说明修改',1,1,'1'),(156,'ZhuanJin/del','赚金说明删除',1,1,'1'),(157,'ZhuanJin/jiesuan','结算设置',1,1,'1'),(158,'ZhuanJin/sxf','手续费设置',1,1,'1'),(159,'Index/getTXlists','提现消息',1,1,'1'),(160,'SpecAgent/newindex','特殊代理列表',1,1,'1'),(161,'SpecAgent/newauthSpec','授权特殊代理',1,1,'1'),(162,'SpecAgent/newbutedit','大区经理修改比例',1,1,'1'),(163,'SpecAgent/showAgent','大区经理下级查看',1,1,'1'),(164,'SpecAgent/selectAgent','大区代理下级查找',1,1,'1'),(165,'SpecAgent/AgentDetail','大区代理下级详细流水',1,1,'1'),(166,'Fish2/index','李逵劈鱼-鱼设置',1,1,'1'),(167,'Fish2/index2','李逵劈鱼-炮设置',1,1,'1'),(168,'Fish2/index3','李逵劈鱼-捕鱼参数',1,1,'1'),(169,'Fish2/edit','李逵劈鱼-鱼设置',1,1,'1'),(170,'Fish2/edit2','李逵劈鱼-炮设置',1,1,'1'),(171,'Fish2/edit2','李逵劈鱼-捕鱼参数设置',1,1,'1'),(172,'Xiaoshou/dsf','支付第三方选择支付第三方选择',1,1,'1'),(173,'Xiaoshou/kefu','在线客服添加',1,1,'1'),(174,'Robot2/game','机器人奖池',1,1,'1'),(175,'index/getonline','在线人数',1,1,'1');
/*!40000 ALTER TABLE `admin_auth_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `admin_auth_type`
--

DROP TABLE IF EXISTS `admin_auth_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `admin_auth_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `auth` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_auth_type`
--

LOCK TABLES `admin_auth_type` WRITE;
/*!40000 ALTER TABLE `admin_auth_type` DISABLE KEYS */;
INSERT INTO `admin_auth_type` VALUES (1,'总管理员权限','174,173,172,171,170,169,168,167,166,165,164,163,162,161,160,159,158,157,156,155,154,153,152,151,150,149,148,147,146,145,144,143,142,141,140,139,138,137,136,135,134,133,132,131,130,129,128,127,126,125,124,123,122,121,120,119,118,117,116,115,114,113,112,111,110,109,108,107,106,105,104,103,102,101,100,99,98,97,96,95,94,93,92,91,90,89,88,87,86,85,84,83,82,81,80,79,78,77,76,75,74,73,72,71,70,69,68,67,66,65,64,63,62,61,60,59,58,57,56,55,54,53,52,51,50,49,48,47,46,45,44,43,42,41,40,39,38,37,36,35,34,33,32,31,30,29,28,27,26,25,24,23,22,21,20,19,18,17,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1,175'),(2,'客服+金币管理权限','111,110,109,103,101,100,99,95,94,90,89,88,87,86,85,84,83,82,81,80,79,78,77,76,75,74,73,72,71,70,69,68,67,66,65,64,63,62,61,60,59,58,57,56,55,54,53,52,51,50,43,42,41,40,39,38,37,36,35,34,33,32,31,30,29,28,27,26,25,24,23,22,21,20,19,18,17,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1'),(4,'管理员','111,110,109,103,101,100,99,95,94,90,89,88,87,86,85,84,83,82,81,80,79,78,77,76,75,74,73,72,71,70,69,68,67,66,65,64,63,62,61,60,59,58,57,56,55,54,53,52,51,50,43,42,41,40,39,38,37,36,35,34,33,32,31,30,29,28,27,26,25,24,23,22,19,18,17,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1'),(5,'初级客服','111,110,109,103,101,100,99,95,94,90,89,88,87,86,85,84,83,82,81,80,79,78,77,76,75,74,73,72,71,70,69,68,67,66,65,64,63,62,61,60,57,54,53,52,51,50,37,36,35,34,33,32,31,30,29,28,27,26,25,24,23,22,19,18,17,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1'),(6,'数据分析员','111,110,109,108,107,104,103,101,100,99,98,97,96,95,89,88,87,86,85,84,83,82,81,80,79,78,77,76,75,74,73,72,71,70,69,68,67,66,65,64,63,62,61,60,59,58,57,56,55,54,53,52,51,50,43,42,41,40,35,34,33,32,31,30,29,28,27,26,25,24,23,22,19,18,15,14,13,12,6,3,2,1'),(7,'11','135,128,127,126,125,124,91,89,88,87,86,84,32,31,30,29,4,3,2,1');
/*!40000 ALTER TABLE `admin_auth_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aglog`
--

DROP TABLE IF EXISTS `aglog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `aglog` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  `account` varchar(255) DEFAULT NULL,
  `agid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aglog`
--

LOCK TABLES `aglog` WRITE;
/*!40000 ALTER TABLE `aglog` DISABLE KEYS */;
/*!40000 ALTER TABLE `aglog` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `agtype`
--

DROP TABLE IF EXISTS `agtype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `agtype` (
  `agid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `level` varchar(255) DEFAULT NULL,
  `auth` varchar(255) DEFAULT NULL,
  `pro` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`agid`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `agtype`
--

LOCK TABLES `agtype` WRITE;
/*!40000 ALTER TABLE `agtype` DISABLE KEYS */;
INSERT INTO `agtype` VALUES (1,'总代理','100',NULL,''),(2,'高级代理','200',NULL,''),(3,'普通代理','300',NULL,''),(9,'初级代理','400',NULL,''),(11,'低级代理','500',NULL,'');
/*!40000 ALTER TABLE `agtype` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `aguser`
--

DROP TABLE IF EXISTS `aguser`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `aguser` (
  `account` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `uname` blob NOT NULL,
  `time` datetime NOT NULL,
  `tel` varchar(255) NOT NULL DEFAULT '0',
  `agid` varchar(255) NOT NULL DEFAULT '0',
  `beizhu` varchar(255) NOT NULL DEFAULT '0',
  `weixin` varchar(255) NOT NULL DEFAULT '0',
  `zhifubao` varchar(255) NOT NULL DEFAULT '0',
  `pid` int(11) NOT NULL DEFAULT '0',
  `status` int(4) NOT NULL DEFAULT '1',
  `jifen` varchar(255) NOT NULL DEFAULT '0',
  `sjifen` varchar(255) NOT NULL DEFAULT '0',
  `paytype` int(11) NOT NULL DEFAULT '0',
  `product` varchar(255) NOT NULL DEFAULT '0',
  `ptype` varchar(255) NOT NULL DEFAULT '0',
  `url` varchar(255) NOT NULL DEFAULT '0',
  `cash_password` varchar(255) NOT NULL DEFAULT '0',
  `bili` varchar(255) NOT NULL,
  PRIMARY KEY (`account`),
  KEY `account` (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `aguser`
--

LOCK TABLES `aguser` WRITE;
/*!40000 ALTER TABLE `aguser` DISABLE KEYS */;
/*!40000 ALTER TABLE `aguser` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `baoming`
--

DROP TABLE IF EXISTS `baoming`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `baoming` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `tel` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `baoming`
--

LOCK TABLES `baoming` WRITE;
/*!40000 ALTER TABLE `baoming` DISABLE KEYS */;
/*!40000 ALTER TABLE `baoming` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bind_players`
--

DROP TABLE IF EXISTS `bind_players`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bind_players` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `uid` int(11) NOT NULL DEFAULT '-1',
  `code` int(11) NOT NULL DEFAULT '-1' COMMENT '邀请码',
  `agent_id` int(11) NOT NULL DEFAULT '-1' COMMENT '代理id',
  `bind_time` datetime NOT NULL COMMENT '绑定时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`) USING BTREE,
  KEY `agent_id` (`agent_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bind_players`
--

LOCK TABLES `bind_players` WRITE;
/*!40000 ALTER TABLE `bind_players` DISABLE KEYS */;
/*!40000 ALTER TABLE `bind_players` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `buycard`
--

DROP TABLE IF EXISTS `buycard`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `buycard` (
  `bid` int(11) NOT NULL AUTO_INCREMENT,
  `number` int(11) DEFAULT NULL,
  `money` varchar(255) DEFAULT NULL,
  `time` varchar(225) DEFAULT NULL,
  PRIMARY KEY (`bid`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `buycard`
--

LOCK TABLES `buycard` WRITE;
/*!40000 ALTER TABLE `buycard` DISABLE KEYS */;
INSERT INTO `buycard` VALUES (2,222,'200','2017-07-22 17:04:48'),(5,588,'500','2017-07-23 19:04:54');
/*!40000 ALTER TABLE `buycard` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cash`
--

DROP TABLE IF EXISTS `cash`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cash` (
  `cash_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `uname` varchar(255) DEFAULT NULL,
  `tel` varchar(255) DEFAULT NULL,
  `weixin` varchar(255) DEFAULT NULL,
  `jifen` varchar(255) DEFAULT NULL,
  `money` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL COMMENT '0 申请状态；1完成状态；-1拒绝状态',
  `cash_time` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  `paytype` varchar(255) DEFAULT NULL,
  `zhifubao` varchar(255) DEFAULT NULL,
  `ptype` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`cash_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cash`
--

LOCK TABLES `cash` WRITE;
/*!40000 ALTER TABLE `cash` DISABLE KEYS */;
/*!40000 ALTER TABLE `cash` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `commission`
--

DROP TABLE IF EXISTS `commission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `commission` (
  `com_id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) DEFAULT NULL,
  `aname` varchar(11) DEFAULT NULL,
  `pid` int(255) DEFAULT NULL,
  `pname` varchar(255) DEFAULT NULL,
  `pnum` int(11) DEFAULT NULL,
  `gid` int(11) DEFAULT NULL,
  `gname` varchar(255) DEFAULT NULL,
  `gnum` int(11) DEFAULT NULL,
  PRIMARY KEY (`com_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `commission`
--

LOCK TABLES `commission` WRITE;
/*!40000 ALTER TABLE `commission` DISABLE KEYS */;
INSERT INTO `commission` VALUES (1,3,'普通代理',0,'管理员',0,NULL,NULL,0),(2,3,'普通代理',1,'总代理',40,0,'管理员',0),(3,3,'普通代理',2,'高级代理',20,0,'管理员',0),(4,3,'普通代理',2,'高级代理',20,1,'总代理',20),(5,2,'高级代理',1,'总代理',20,0,'管理员',0),(6,2,'高级代理',0,'管理员',0,NULL,NULL,NULL),(7,1,'总代理',0,'管理员',0,NULL,NULL,NULL);
/*!40000 ALTER TABLE `commission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `config` (
  `id` int(11) NOT NULL,
  `name` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '平台名称',
  `ip` varchar(255) COLLATE utf8_bin NOT NULL COMMENT 'ip',
  `url` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '域名',
  `zfurl` varchar(255) COLLATE utf8_bin DEFAULT NULL COMMENT '转发域名',
  `card` int(11) NOT NULL COMMENT '1 金币  2 房卡  3 房卡+金币',
  `gold` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '1 推广额   2  金币',
  `account` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '特殊帐号',
  `path` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '后台路径',
  `path2` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '推广路径',
  `lurl` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '推广跳转地址',
  `down` varchar(255) COLLATE utf8_bin NOT NULL COMMENT '下载路径',
  `zf` int(255) NOT NULL COMMENT '是否转发 ：1 是 2不是',
  `zfpath` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config`
--

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;
INSERT INTO `config` VALUES (1,'巅峰娱乐','39.108.141.246','dfyl.z8w2.cn','dfyl1.z8w2.cn',1,'2','123!@#cnm','qp_host/','qp_ht','39.108.141.246','qp_down',1,'qp_ht');
/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_bzw`
--

DROP TABLE IF EXISTS `count_cost_bzw`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_bzw` (
  `bydate` date NOT NULL,
  `totalcost` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_bzw`
--

LOCK TABLES `count_cost_bzw` WRITE;
/*!40000 ALTER TABLE `count_cost_bzw` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_bzw` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_kwx`
--

DROP TABLE IF EXISTS `count_cost_kwx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_kwx` (
  `bydate` date NOT NULL,
  `kwx` bigint(20) NOT NULL,
  `psz` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_kwx`
--

LOCK TABLES `count_cost_kwx` WRITE;
/*!40000 ALTER TABLE `count_cost_kwx` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_kwx` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_niuniu`
--

DROP TABLE IF EXISTS `count_cost_niuniu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_niuniu` (
  `bydate` date NOT NULL,
  `totalcost` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_niuniu`
--

LOCK TABLES `count_cost_niuniu` WRITE;
/*!40000 ALTER TABLE `count_cost_niuniu` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_niuniu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_psz`
--

DROP TABLE IF EXISTS `count_cost_psz`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_psz` (
  `bydate` date NOT NULL,
  `totalcost` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_psz`
--

LOCK TABLES `count_cost_psz` WRITE;
/*!40000 ALTER TABLE `count_cost_psz` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_psz` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_ptj`
--

DROP TABLE IF EXISTS `count_cost_ptj`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_ptj` (
  `bydate` date NOT NULL,
  `totalcost` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_ptj`
--

LOCK TABLES `count_cost_ptj` WRITE;
/*!40000 ALTER TABLE `count_cost_ptj` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_ptj` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_ttz_mp`
--

DROP TABLE IF EXISTS `count_cost_ttz_mp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_ttz_mp` (
  `bydate` date NOT NULL,
  `totalcost` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_ttz_mp`
--

LOCK TABLES `count_cost_ttz_mp` WRITE;
/*!40000 ALTER TABLE `count_cost_ttz_mp` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_ttz_mp` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_cost_ttz_sp`
--

DROP TABLE IF EXISTS `count_cost_ttz_sp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_cost_ttz_sp` (
  `bydate` date NOT NULL,
  `totalcost` bigint(20) NOT NULL,
  PRIMARY KEY (`bydate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_cost_ttz_sp`
--

LOCK TABLES `count_cost_ttz_sp` WRITE;
/*!40000 ALTER TABLE `count_cost_ttz_sp` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_cost_ttz_sp` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_goldcost`
--

DROP TABLE IF EXISTS `count_goldcost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_goldcost` (
  `byhour` datetime NOT NULL,
  `totalcost` int(30) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`byhour`),
  KEY `c_time` (`byhour`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_goldcost`
--

LOCK TABLES `count_goldcost` WRITE;
/*!40000 ALTER TABLE `count_goldcost` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_goldcost` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_time_user`
--

DROP TABLE IF EXISTS `count_time_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_time_user` (
  `c_time` datetime NOT NULL COMMENT '统计时间',
  `num` int(50) unsigned NOT NULL DEFAULT '0' COMMENT '当前在线人数',
  PRIMARY KEY (`c_time`),
  KEY `c_time` (`c_time`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='统计每日每小时在线人数';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_time_user`
--

LOCK TABLES `count_time_user` WRITE;
/*!40000 ALTER TABLE `count_time_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_time_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `count_topup`
--

DROP TABLE IF EXISTS `count_topup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `count_topup` (
  `byhour` datetime NOT NULL,
  `totalmoney` decimal(30,2) unsigned NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`byhour`),
  KEY `byhour` (`byhour`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `count_topup`
--

LOCK TABLES `count_topup` WRITE;
/*!40000 ALTER TABLE `count_topup` DISABLE KEYS */;
/*!40000 ALTER TABLE `count_topup` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dolog`
--

DROP TABLE IF EXISTS `dolog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dolog` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `txt` text COLLATE utf8_bin NOT NULL,
  `uid` int(11) NOT NULL,
  `name` varchar(255) COLLATE utf8_bin NOT NULL,
  `time` varchar(225) COLLATE utf8_bin NOT NULL,
  `ip` varchar(255) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dolog`
--

LOCK TABLES `dolog` WRITE;
/*!40000 ALTER TABLE `dolog` DISABLE KEYS */;
/*!40000 ALTER TABLE `dolog` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `game`
--

DROP TABLE IF EXISTS `game`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `game` (
  `gid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `status` varchar(255) NOT NULL,
  PRIMARY KEY (`gid`)
) ENGINE=InnoDB AUTO_INCREMENT=300001 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game`
--

LOCK TABLES `game` WRITE;
/*!40000 ALTER TABLE `game` DISABLE KEYS */;
INSERT INTO `game` VALUES (1,'房卡牛牛',''),(6,'房卡斗地主',''),(10,'房卡十点半',''),(17,'房卡推筒子',''),(19,'房卡拼天九',''),(24,'房卡三公',''),(36,'房卡扫雷',''),(51,' 八张清','1'),(65,'房卡牛元帅',''),(77,'五子棋',''),(78,'房卡跑得快',''),(10000,'卡五星','0'),(20000,'炸金花','0'),(30000,'金币牛元帅','0'),(40000,'豹子王','1'),(50000,'金币拼天九','0'),(60000,'百人推筒子','1'),(70000,'金币场推筒子','0'),(80000,'金币跑得快','0'),(90000,'神仙夺宝','1'),(100001,'龙虎斗','1'),(110001,'一夜暴富','1'),(120000,'摇塞子','1'),(130001,'赛马','1'),(140000,'幸运数字','1'),(160000,'龙珠夺宝','1'),(170000,'地穴探宝','1'),(190000,'疯狂翻牌机','1'),(200000,'名品汇','1'),(210000,'红黑大战','1'),(220000,'腾讯龙虎斗','1'),(230000,'百人牛牛','1'),(240000,'鱼虾蟹','1'),(250000,'捕鱼','1'),(260000,'百家乐','1'),(270000,'水浒传','1'),(280000,'李逵劈鱼','1'),(290000,'金币场斗地主','1'),(300000,'红包扫雷','1');
/*!40000 ALTER TABLE `game` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gf_account`
--

DROP TABLE IF EXISTS `gf_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gf_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uname` varchar(255) DEFAULT NULL,
  `qq` varchar(255) DEFAULT NULL,
  `wx` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=28 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gf_account`
--

LOCK TABLES `gf_account` WRITE;
/*!40000 ALTER TABLE `gf_account` DISABLE KEYS */;
INSERT INTO `gf_account` VALUES (27,'1','2','3','1');
/*!40000 ALTER TABLE `gf_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gold_detail`
--

DROP TABLE IF EXISTS `gold_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gold_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `gametype1` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `gametype2` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=43 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gold_detail`
--

LOCK TABLES `gold_detail` WRITE;
/*!40000 ALTER TABLE `gold_detail` DISABLE KEYS */;
INSERT INTO `gold_detail` VALUES (1,'卡五星','10000','20000'),(2,'拼三张','20000','30000'),(3,'牛牛','30000','40000'),(4,'豹子王','40000','50000'),(5,'拼天九','50000','60000'),(6,'(百人)推筒子','60000','70000'),(7,'(单人)推筒子','70000','80000'),(8,'跑得快','80000','90000'),(9,'神仙夺宝','90000','100000'),(10,'龙虎斗','100000','110000'),(11,'一夜暴富','110000','120000'),(12,'翻牌机','180000','180000'),(13,'捕鱼','250000','260000'),(14,'鱼虾蟹','240000','250000'),(15,'红黑大战','210000','220000'),(16,'百家乐','260000','270000'),(17,'翻牌机','190000','200000'),(18,'名品汇','200000','210000'),(19,'赛马','130000','140000'),(20,'单双','140000','150000'),(21,'龙珠夺宝','160000','170000'),(22,'地穴探宝','170000','180000'),(23,'五子棋','77','77'),(24,'百人牛牛','230000','230000'),(25,'存入银行','-10','-10'),(26,'银行取出','-11','-11'),(27,'赠送','-5','-5'),(28,'后台赠送/后台扣除','-1','-1'),(29,'用户充值','-2','-2'),(30,'推广额兑换','-3','-3'),(31,'提现返还/提现申请','-4','-4'),(32,'腾讯分分彩','220000','220000'),(33,'水浒传','270000','280000'),(34,'摇塞子','120000','130000'),(35,'百人牛牛','230000','240000'),(36,'转盘','-8','-8'),(37,'转盘','-9','-9'),(38,'救济金','-12','-12'),(39,'斗地主','290000','300000'),(40,'红包扫雷','300000','310000'),(41,'金花包房','1000','1001'),(42,'李逵劈鱼','280000','290000');
/*!40000 ALTER TABLE `gold_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gold_robot`
--

DROP TABLE IF EXISTS `gold_robot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gold_robot` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `status` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `gametype` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=13 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gold_robot`
--

LOCK TABLES `gold_robot` WRITE;
/*!40000 ALTER TABLE `gold_robot` DISABLE KEYS */;
INSERT INTO `gold_robot` VALUES (1,'百家乐','1','260000'),(2,'百人牛牛','1','230000'),(3,'百人推筒子','1','60000'),(4,'豹子王','1','40000'),(5,'红包扫雷','1','300000'),(6,'龙虎斗','1','100001'),(7,'单双','1','140000'),(8,'鱼虾蟹','1','240000'),(9,'名品汇','1','200000'),(10,'牛元帅','1','30000'),(11,'跑得快','1','80000'),(12,'拼三张','1','20000');
/*!40000 ALTER TABLE `gold_robot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `img_upload`
--

DROP TABLE IF EXISTS `img_upload`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `img_upload` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '图片ID',
  `img_url` varchar(80) NOT NULL COMMENT '图片地址',
  `list_order` tinyint(10) unsigned NOT NULL COMMENT '排序',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8 COMMENT='轮播图片';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `img_upload`
--

LOCK TABLES `img_upload` WRITE;
/*!40000 ALTER TABLE `img_upload` DISABLE KEYS */;
/*!40000 ALTER TABLE `img_upload` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `imgurl`
--

DROP TABLE IF EXISTS `imgurl`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `imgurl` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `imgurl` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=680 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `imgurl`
--

LOCK TABLES `imgurl` WRITE;
/*!40000 ALTER TABLE `imgurl` DISABLE KEYS */;
INSERT INTO `imgurl` VALUES (14,'http://thirdwx.qlogo.cn/mmopen/vi_32/4hItYW8gvaXk7je5EjhGfOfEuKkCbu5icfV2m8PrcibUlauYye1EFzGuWW4JUY9I9bOIEC0Uf77ibUjLNpIof88rQ/132'),(15,'http://thirdwx.qlogo.cn/mmopen/vi_32/WaKoo1LCqEvazFWBhRBkWLcu4TpRG1YB8X5ZrTicIHryY7ViajOBfCyKDll40sPa0S5cU8w6g8VhSF48YBQiccibrA/132'),(16,'http://thirdwx.qlogo.cn/mmopen/vi_32/8aZxWBXDicESic7ZNNdolA43QspibW8vI5Rb9sxMGQTxQxrvG5a3yQgaKSOl6ATkSSNicEYe7ntqcFWibIeZzx3QV9w/132'),(17,'http://thirdwx.qlogo.cn/mmopen/vi_32/jSicks21ko67cfcc7PTiaxN1mRSIibgvLxh2k5WovibIymy9X3EFvSeeAHUKaV8Om3AlpOUehCqgZ9bpTwOEflwI4Q/132'),(18,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTL0tPwpUVQT51gtQJSQZ6p9ibhoErrnVXyiaUic8WKhib6tGhicerGJznX7bRgtGBC83A4mOuttHU2u3bQ/132'),(19,'http://thirdwx.qlogo.cn/mmopen/vi_32/ArOtUJfIKf6uu6xRAhwqibGoNeAfiak4HytJ9QVOGxgHeiclJib7yt3FZKmeDibEoDMlFascycsEnRGYRJtoh22AN9Q/132'),(20,'http://thirdwx.qlogo.cn/mmopen/vi_32/Mg64wACyxmgutoH5gDSmFiaEdKO0DjqaAAY33GMvmSp1LkmBbJOPjqAxwNl5JfkBPnvhmrx3tATWQhA0eUYzC5A/132'),(21,'http://thirdwx.qlogo.cn/mmopen/vi_32/N6EsZ7ATgl4D2j294lrYyibwz91IbzibbPrOsyFgmoKMianSKoGUw2OMA8UkJhLtjaGDdQdl3v8HLUnib5mLRibd36A/132'),(22,'http://thirdwx.qlogo.cn/mmopen/vi_32/icBibib5VGsQmsAGHiaicbIlpeMK8M47fPzl2n4ptHx8vpica8DONmEUlanDakfIuyJIMvfKPRORYGnDEl3BASiaJNFMw/132'),(24,'http://thirdwx.qlogo.cn/mmopen/vi_32/xkn2Cj1Qb4ia3xdSGulDgNFGNhMDp19biawL76HfWOQgF7B6Oj4hS1nHmT5mz42AbRL2yQLiblyd1uH4XUw5rcYNg/132'),(25,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKiaVde2marotXakiceicFeJ2GI4CiaoJb4kn3bYwczfZVL9lLbAz5wTM5YOad6KRxG6ujiaVzoSbLbDZw/132'),(26,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ereZ2nEftpwz2qs2moHlstct65B0fbRjDy9MI02ZOfOr4bJLJXkVo0WCBlXibPicbZNk6snybjvUrjA/132'),(27,'http://thirdwx.qlogo.cn/mmopen/vi_32/gZ5mBEwArcAgd1ycUhtG2ia8icHEVLcGKOGhibnskIGrF7B3zkibP8oTLxibELHQOxq8I1VTiaWFPrvrO9SsZeXicIvDA/132'),(28,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqDqGW8nFywP99YTSKHsYicX2LfMeJy0uQfoSNue5swG9zMdby0rPYz9x5FlNJSoD0GP8ibmSgoEXBQ/132'),(30,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqfS8Bib697qALExiaGq7MC0Tlz6hflpjDwz6G1ASbJXjtowJzsNFyAgm8XvQCKiaK9ldOUia9ibO8dHSA/132'),(31,'http://thirdwx.qlogo.cn/mmopen/vi_32/ebDwVlG74Q4WWaGffu2D5o94xarI19ibzic2nicv0tta1tiby8n3HqTawpegVqqm4gNpjw6Rs3Das3QZJqib0UexUsA/132'),(33,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqniculX4ysFaCZBpzibYRttuOR9VtXCXvrWiawicDFZADlibdcViaibZz0CZQAZBCVn72dgcUcRSNylFUKg/132'),(34,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epn3Bu3p5Qux6B3kMtfia9TA9ngfZ4cgT4JayYxhWWxRYcUUEQTtkjyoibXF3uLCqKY4paSwJ0M9AeA/132'),(35,'http://thirdwx.qlogo.cn/mmopen/vi_32/4xw0nhwpj2ib0cxvakskCEwjIQzbJL8Abh3aYyXiaBeT3Ntw2YaNZkAERfLgbtg4LkF6eIwZ8KY2za91wEq0kFqw/132'),(36,'http://thirdwx.qlogo.cn/mmopen/vi_32/nribDiatSjlg1zsn8icILQmHglPo74OdaDoNlKsiaRqXtgoibGUc0UT9d8I6Vqs2cQVTDFGsJ9h4w54wazmRTZ1tibng/132'),(37,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ17d6ibLYne9icKiadmxGzh0V8u1B3YSFqhEXQiayY5GtNR7QUB9angttwuWzaSteGGmug4c1YSpIreA/132'),(38,'http://thirdwx.qlogo.cn/mmopen/vi_32/ECk5u9dOKo71YT3sYQqw92waJWvdFtf11RoCEQW4BQZ70V4pFZGS8xYeD93Ps6hvVl5dIMeeLlR88fHCBpQibSA/132'),(39,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIsCumQjY5FiaAJlPHYYShhbricS0WkWmzfueHJw95UAOlLiaUEMl54NQ58VNFZF59GJdY29AYqrvkkA/132'),(40,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTI6zpHWIiaWSoT0dLR38rlJzbIrXou50iaR2iceQiaysTjYWWc2nZuxwlQdSmWnYenomxk8dny8GMicmyQ/132'),(41,'http://thirdwx.qlogo.cn/mmopen/vi_32/wTBD7GbNDbSGzkhrXgINLWNP8nhdDiajSxHnAkAS0icFM9xWSib0PoHhrFzNFUyDAsgQej3eQEuVAkTeAIZbkA5QQ/132'),(42,'http://thirdwx.qlogo.cn/mmopen/vi_32/TUp5skxNkMmlSssA9kcMV7z89iaJ0z4ee7ADbNhukuQe0r1aFZEg3cVEZa5NbBbqfHdQ2RUlCNMibeUibXRr9CO7g/132'),(43,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKAXY2brIaKwfH1TSRdAomW0Os4zd8S8wib9Duq5osLJ4gNZF2OT12ib4FWTp4XF0BWd3xNa49bzIPw/132'),(44,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTL8nR7kBpJIFaEOKrSOqk4uBO0tzeGKGsUZomkdFhWFSjbPnKpxuianWtUd00GET55xExYZVwK0J3A/132'),(45,'http://thirdwx.qlogo.cn/mmopen/vi_32/CSFPDnXAaHdfCNcDnicfPpFiba84O9ZYKibTuvhDTcLibt5v2PIyBjYFyDAibDAhJSHs7x4ba4fmETowqrbyJWJeicJw/132'),(46,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJXwkccAK33P0qfziaN2FZjwSw2cNib49jMFR5ichkcOe7iaqhJBHEQL5SfG3SKN04weaiakHFedJsJibPg/132'),(47,'http://thirdwx.qlogo.cn/mmopen/vi_32/mQP4cxeZGbCUJOyub1jHHibtQlSfibkQKhmwW95iaqBE1dg368WPtuAIxGE4niahSv8CcMx11wPGUTFmuJ2GsLf22w/132'),(48,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epUUde5bWCY2FZkIj4YL2piaE2RfUf6cGQagL7NQPJWqMT3Z7SwLxF6oEqDibdibrjMswibLn25d3VfgA/132'),(49,'http://thirdwx.qlogo.cn/mmopen/vi_32/CgbH0tCo7ozEHBaXNiaM2YiaAU6DTk4sKvPVFV9pOOT6j2S2JibJBJJOZ5y4tIQrExSMsGYO6sSD904fmkhNLHJ7A/132'),(50,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTI61jAFJXGuHib1jVUPm60SNtq3a8JBQIXibO8jDQjBt3A2DXe4nOAh4x5ANBnU8jUsYNIlcBvmibRgg/132'),(52,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIogenCZWJtZeHajE10Liaic3jDsLVY3Ia6H6nsN8Mao2cjicYztXP4VqWYd8BBAJELdmtXibxAliaXPDg/132'),(53,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIyqibpScLw6fTxA4g48a7su4QibeTQ8TYV5VSmcuKWYxYFEBZ4y7h55sLbXZTwJ3e18N0Nxg3f4qsg/132'),(55,'http://thirdwx.qlogo.cn/mmopen/vi_32/w4FoHuvj2qGyBGt8N5uicoQTTgKUD9XrH5raeDMV657Q61WOm4vNpQqjnh5RlSJCIJqs3diakqOWibRZGvMV5plTA/132'),(56,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoib30YCWy46mJXdg6iaGmc5icShKImOLVFRHa9Fyrwz5VgMrkibraZicuTlTnZm489OQiapUcw2b9cbbTw/132'),(57,'http://thirdwx.qlogo.cn/mmopen/vi_32/OvdhvW1kly34WygYOTZVrMPQDYCicYbJdawuIhK0fF98J5NCicicRAYgHqolSgn9fNxsduezargcciaOWuZw1vrvdw/132'),(58,'http://thirdwx.qlogo.cn/mmopen/vi_32/qyj3gpgbpwYiaShCzmzviapKdibmqVTIX828ausFJT5GialrfW1M0lz1JMP508NfUDeiahS7bZYdVicicVUSMVUEnWOgw/132'),(59,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epUdrH6Tfk8S4Z4ZNTGicoZFbBiaMqlLB2SjTjaQAhsISkALF1kEbMx9R6vqA6xb75iaXx3Xjmt3yumg/132'),(60,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTL592iaRam8yp0FFmicKsJFNNcjWNHicbNuTRicxvGS1gDgrRoGNGib4X8boHx6F5PQKcbhFX9gzMKI0icQ/132'),(61,'http://thirdwx.qlogo.cn/mmopen/vi_32/ouACJY5vpoSXPbM2jIjdrSiaBpVAt3XgMkIW5R26EjaOrAbdpZWiaePcvg6mZA33zFYkv5KggHzjvgkxouibPxBLQ/132'),(62,'http://thirdwx.qlogo.cn/mmopen/vi_32/2KwSPLVCORokAupXuMGOdia1pLWkbiaPCpm9g1QibkMpcic1QkTZDOaTgksWYflQ8giaqibJt3b79c2w87NKsw7Ushiag/132'),(63,'http://thirdwx.qlogo.cn/mmopen/vi_32/fYicaGpRc97MXU4esrPAO6XM6fPZjlI38gIibDqA6pEGrGIsCbjCib0iae0WCEicTatHjzT6eUicbFXfuiaHiaibu4J8qkg/132'),(64,'http://thirdwx.qlogo.cn/mmopen/vi_32/prWpwVd7UZzTMJ9KbibbrvQZ2iaIjjVCA1uhIHAWcC57EzPab2JjFwX6lZGdnuUicTpicIL9IMs1p0abVnczFpUMzQ/132'),(65,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIgc4fnIibZCmDPH5mmNic5fCPDukBsVtTw0Vnh4qByr6Ss7eGBVhESREoAzHiaHGqFQ1Ympibic50ShicA/132'),(66,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIwIc0CicwicJUxhLycibFtVib6HwPxzIh9rFB7wk4I3O0lKeoqEoibrxzcibC0RYkEOwcS9k0yvgT6cqDQ/132'),(67,'http://thirdwx.qlogo.cn/mmopen/vi_32/XuAYswwfiaeHqeRa1cXzz8Vibw4zmjQn4TQ3DQMnmWd3Jm9FxOeffYkpnH1QqzrmOcNfkr1GxjJ0KlfhFOGDdRibQ/132'),(68,'http://thirdwx.qlogo.cn/mmopen/vi_32/hcRLck95z1TX5R2ug7qicjiaINdwWSaSM79SDHLYWwMw7h3ccw6kpm57CcibNIEKSBtJnkGBCk50dqVZf9CiayG2uQ/132'),(69,'http://thirdwx.qlogo.cn/mmopen/vi_32/UHFUH62JmqwICU9EwyrkdQ9PaVsytWYVOv4zdOmaWe8SD4EtnicsXhnshRDf4Df3mS1aib94VkibuzUpeLXNSLxEQ/132'),(70,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoc8sReia18fc79thIzx8orSJs7ICS6bWe0vDy6jh9vldTtSlEbsn9O2utIyIDsGHd6wicKAnt4m05g/132'),(71,'http://thirdwx.qlogo.cn/mmopen/vi_32/2twV5EOILQLNicSyXlHgK3rIvE7ShSZU54ibp9OyiaACX839LBCica3GjmrAyy1Ge65m53fgmx5JiangLqMomEwE2ew/132'),(72,'http://thirdwx.qlogo.cn/mmopen/vi_32/e5VxicHXibh9ctGkhgiaVJVriardz5Tib5xoKRXgn3K1rYmCmQfiby0zLQZvU6O4znX0xrpoWR5m8CibZBVQ4w98v3TAw/132'),(73,'http://thirdwx.qlogo.cn/mmopen/vi_32/iaA3KELAf1taTVV03BCsvdgQpfEBml7mzdvUV8dt56WMtAcv3RsIm8bzpyvwHpBhA9SPgjZxLRibVL5dGXPgLuzQ/132'),(74,'http://thirdwx.qlogo.cn/mmopen/vi_32/xU7Pickem2iaRtIKiaBfVIA0ic7KDLy8xKwnw0h0s9UQhqJJSkfd1l5E6jqXZTM8ucgN1icwicNZG3bP2PxLGd3C69zg/132'),(75,'http://thirdwx.qlogo.cn/mmopen/vi_32/wAKc9N4y0JjPLfQSWxQV4lCqGDVvftNpNY13SoCGic2pwcEVOgVOAJSh47Qemlw8xcb2JYQjZDGH0meokIWmtFw/132'),(76,'http://thirdwx.qlogo.cn/mmopen/vi_32/mDf7g6pmYhgFYJqRlRsiciaguLCNI8beO2A9uJnzPlHeDhSQLibtw4xpbCywibGrPFvzpsfxWz6jhJnNpNC9Tgm5gw/132'),(77,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLghPQ8Qz0HUaPFkhk7B3E7pFR8LpB0Nkz5RRCpvyZoYRgVIgvy1pqKHLnRDibQzPlm2pKX91FK98A/132'),(78,'http://thirdwx.qlogo.cn/mmopen/vi_32/GVenxYicdvicibia7Fj6LSl4jlAHCdaWPyIoVFUiaIS2838cmxSA1NbApQonicCp4ibpoenalcH5fp0iapIHB3EIBEafMQ/132'),(79,'http://thirdwx.qlogo.cn/mmopen/vi_32/Sh7zQibC1xQbC1SmkgGF0YHDhYeMulNGtxs4cwdnKYABPRr5AQlAtmMEUaRAVtPBznb2z18ImU2HjicKXw6fTDBQ/132'),(80,'http://thirdwx.qlogo.cn/mmopen/vi_32/gFsre4I5goqZu6tf5I27ddnesgFAX0On37194d3GW3lvEN2XicFzzZfHibk78OYeUeWp6CIPpicR8HWhTOtnBqciag/132'),(81,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqniculX4ysFaHVp80eGtia1CF22mDNwaMnic5mtvpU1YqMCoW9ibJkaXqGDA0PWeAGJ2fs4Tv4hBY6yw/132'),(82,'http://thirdwx.qlogo.cn/mmopen/vi_32/nZqoatYrvfaWsZqqtTWlg0yYHfSgRPHxg0kpqZicVYzIaaGLHn7N2lAMI8fGv0yhJ6MZCDagSBjR8IOBIcSHrhg/132'),(83,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibhJPEo0GV3rJsTsrj8hOk2XPdrrWnGXFMEMTS27zLGzI7ndE7wicXVSVv8iaFPAIFHicBGALb1Qh6NWSpz0k2vXeQ/132'),(84,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJvB4Q9Of28J28aV3QDXVSm7aCUnkN6IjFwCVicveU6DB8ddpWufPxNrPeK02P4eodA8hz3zJNwE4Q/132'),(85,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKyzvPDITEkIjKxhICc0iaPoAxOkf52zaVWl6JWpVvs8RTrB7oM5K5tP0ASFD4UicZg3k6E5XyxR7CA/132'),(86,'http://thirdwx.qlogo.cn/mmopen/vi_32/MWxgLoYcJsBjK9P4iaeevrdIIDhBFsQ19Stll9cpMw5MAUWVMzPT9uRIgblKcXofLwhZkftJlah2xRID3dNabEA/132'),(87,'http://thirdwx.qlogo.cn/mmopen/vi_32/yNrnSzMdUkOiaOrP8Tblhpq0ibDVciaUeE37xtYDmupE2GxGAm5oHBkPJ9MQ3jg50ySsibO1W2xIibex10VMhOQ8WRw/132'),(88,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eohkClhYlsyiabWicIKckBickGF9zGbRZh6iarkZZG1vnicJicdtbTYibpEjsEC83I5gu7jNBcwOtzPXflQg/132'),(89,'http://thirdwx.qlogo.cn/mmopen/vi_32/fSXX3w2koF1fyINEoTRiaWEsUbMQSUsspwHmoSd3HG9n9fms4BTcWvls6UoV8JrFI5uxqTO0v5yI2zt5M7tBGSw/132'),(90,'http://thirdwx.qlogo.cn/mmopen/vi_32/iaTZIDLpWcfn8JShzfN5e89tWC1oeCjfKVKFkHelIcdqGOxDTaVWVLv55UicDHsEXNPorNX2GeQorZXf74wNPAxQ/132'),(91,'http://thirdwx.qlogo.cn/mmopen/vi_32/FydBiarFhvwDiaJickXJib4ukuqVPJLuyShNwxc7P5POYZqAU3mQLxiaF3qZJcwQwFSuUiaety8DCS9UoYVtXMvc6AEQ/132'),(92,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLvMso0GxKt2ZNdibY9wCPurhKJSttOUV1jK3fmB0hnnAR3DOllDMDIodtDFiaFfmMWotMLUxF6M5uw/132'),(93,'http://thirdwx.qlogo.cn/mmopen/vi_32/mPYePe5UoZQMia8hdGF0PWcx4Ye0iaCtvnpdE4597mkEiacJApEZwUhfBsg5PXoAfsVFSQ0qdB4HmfKd8jibElo0OQ/132'),(94,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK6JRpRMJLHZpIeia8I0TibxXuNzHDibaVk9yLy4hKXR4OFUlPvgCILAUmqKC2orN3TEzv1K2gINTHsA/132'),(97,'http://thirdwx.qlogo.cn/mmopen/vi_32/jV5cY8Cia7ibqmu3C6mtZoUhU84BrBGHHDETZewrDwmqonEbHKs6EtPswB73WG8SapsSgZhe5f8lWghLgyVryDng/132'),(98,'http://thirdwx.qlogo.cn/mmopen/vi_32/0gk21OaaTBGdjnXPSwwz2ksU8G48dm60dLEYbJCy7pZL2u3VYgrYW92LwIphIibv5YlN4NGJGlNrySJfBicXLoicQ/132'),(99,'http://thirdwx.qlogo.cn/mmopen/vi_32/38JaoO4ZWoaBrjiakEFVOx9djDBjx2sZKCwzrw5wHuTlKPckia140kDzHXzV2WjL3qBZPqTOicujRMnRMrK1PMr0A/132'),(100,'http://thirdwx.qlogo.cn/mmopen/vi_32/BDdFWltia48hiaFFNjnbsZcArpfU9d9bm590KmUVh8amUL7QPI5SUUvETTx8PibRsDlkicSXXOZWebic1tgMbhz2iamQ/132'),(101,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIia6JywZv309tQPiaANVK5ib3cmHtUbxnbbUFSP09rGLoAiaZFUqFALia4HsiaSxQicDw2jg2PQrKK75ucQ/132'),(102,'http://thirdwx.qlogo.cn/mmopen/vi_32/jb4Ov0UXOYq7pABSujSoh2hsAmqIg1Nol8fLGibOwxXxLpiaIxTFjE4JZ4nCYCGmcmsuqxngGOSyzsslwd5GZXCQ/132'),(103,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTILzhXtvu3iaEPU4vU37fI5Dve0UQf0KoCtllibhnzDU3sddpJHhphnk46Q1SvtVtczZib5dfD1co77g/132'),(104,'http://thirdwx.qlogo.cn/mmopen/vi_32/TVCk2MAggd5l4MrjKXGsKkHICJGo8eTeFiavlFgvxbJpdM8kunlJXwBiasCLIUAgkBuXBDp720zBicUNA6TibnKs7A/132'),(105,'http://thirdwx.qlogo.cn/mmopen/vi_32/lBE5saHColvkX4fpw9bdS229M5vAwMWY7FHicPoqkOWuzo9EhxVRfVVVUo8ZficmRYQZ3lhib1sf7F7kUar2iaLb5w/132'),(106,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epZpp1BTStA51I5UjwHiatHMjiaM8TBH812LuKwzEV0CbSrVicBTfNt0MTNB1z75syJT1BEKE3tY0low/132'),(107,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epXUe7Aj6M6ACBED8OvwWfbKE1QSODlOlTRaH2V2P3icteVSoCaZiaxf7pMJKEWfibrwl0PaDZy6KfaQ/132'),(108,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLS9O0bSBOnVanBBpcAcX1GbkaXGHz3ibldIaSFkhklibg8UnvtO1wTpWiaTo3dNAdVYm3YYmKkekBeQ/132'),(110,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIs2r0zgqD67uLb3TQfT5PHmuGAh9ib8FicAoyDUV1yNYAhoyL3gMOuiadz9g8Rt1dcuQyibVTCxAVA7g/132'),(111,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLjGXj4djLsibKe4FTyHehxobeP9ASjf6jQ4IV5PfVdgJS7Olfml7TRLddh3fteBiaXnzOtTQZ87wicg/132'),(112,'http://thirdwx.qlogo.cn/mmopen/vi_32/jVgmicj4L5niaC2hojHm9onsN4uVETUka5VwpibsUmGjBnUicFeUHqicX8kewGTOmV5ERmwtunKeBCWQcZ4U27ldgug/132'),(113,'http://thirdwx.qlogo.cn/mmopen/vi_32/ia9JP3bYdNeaYiaIpXyjbEGPMCYZW1nEiaNWfHXgGn3e6RibZmmVHfkdHp8AujNjyed9mkFC0WYhS1UBz00zkptT6Q/132'),(114,'http://thirdwx.qlogo.cn/mmopen/vi_32/yAuIVrgu3on1SRY1N61pGyFsBiaLQeEXFXuSibQIfwhoUnNQ8706LBJGIMDYo5FbXSpJSaRljT5EBWYTVicV7pFbA/132'),(115,'http://thirdwx.qlogo.cn/mmopen/vi_32/LXyb0jRQunknu7hjSicAh1jnJR2icuL6k2kNoo2l8V44HPsqyTmia6HEK5NqMqoCu7SIpIw1QADgVOIYE1e5rIkaQ/132'),(116,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJdmZZoAoRFX85854dVA9wLicHZORxH8IK0meoLbjP8ql3F9JYGehwXzNX6JY5xljckYI0dqCaFf0Q/132'),(117,'http://thirdwx.qlogo.cn/mmopen/vi_32/4iaKLibHRojFIIrlpaD74GaH8SDFaKqoZQRjkzUPYl6ah6N7xy5eCLL3YVDOFoSkYmo7icZuYAoBUUkXoibhRkQ3eg/132'),(118,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKEibyh9M276E4VCADUMXpATJ8LhhN1PaIWP7vJCGEV0393oY1Q9v7PO4BtNickIdPDYnCasPUk0Nibg/132'),(119,'http://thirdwx.qlogo.cn/mmopen/vi_32/JUjLTowl61Xo5P1RJR1jg6zibiczHGIHfWiauYwERAFeDHM4kaf5vdvYLj8vEkW3KhKN5ViatpYj2BFbFnPx3yP7yg/132'),(120,'http://thirdwx.qlogo.cn/mmopen/vi_32/mGfgwCK8RuXBd38V4VCX1UFXSvBJibaV13hq3N2P67S6D3V7CwtDm0KXsicqzYYl3oaJibMsWGaQd1egsmr8sibHjA/132'),(121,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJtS4Tb8iaPDZwzrX53hHAzDXHESADLhd4FOq2sjtTTSIaqgcTymmBqZiax3akQ0C0ckTC44n7fxQcw/132'),(122,'http://thirdwx.qlogo.cn/mmopen/vi_32/PUjQaIStY5PoRUo2ubp9Old7zYT9wdm0vGZC5Hy9appXM2Al8nsFS3a6yTdW0e3yuXNcX8a0HOIuYZRpocDjzg/132'),(123,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKdVqBynhdb7xyCpGnKzX0j62ib1ibwibaLDESJUuWHFqM7pn9CybUuUicOoypgWS6heiba54BNhXP3VSg/132'),(124,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIfwha2lRtPQXofQh0P5CccHxOR1j6Walichn9GrDRW3fov9bRsH330zQ9qdg2TE5viamiaudGsyIkJg/132'),(125,'http://thirdwx.qlogo.cn/mmopen/vi_32/03J1ERdSFyIePVrTafueIcNvRStP3USs0DTgen4xzD8DL22q6cC8HPoNeMnuibDbzckvAlKzyyeswTA8hXTezHQ/132'),(126,'http://thirdwx.qlogo.cn/mmopen/vi_32/cEiaj0hvn7ADP10y7hYXIn74CoEe7Y7cdp9JrJy7OIwbVC0mDN3t1MjlOCfcssw7AyJDDZMkMQwBuia4elmyo82w/132'),(127,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKDfy7icjzBBU8Y3B8J7GON2hrAl3s3v8N6Z2tC9ent0faDA4J83JfMzDIqnhNsgDPWfPB3cjbya3g/132'),(128,'http://thirdwx.qlogo.cn/mmopen/vi_32/fNYicSzF7Da8vf4N5Sib6JQ2aWtox9stmJZ87IFeY2QIvt9eR4ibkTbkn93JUxQw1KF6IJeEMg2Y3T3emk3UKRiaag/132'),(129,'http://thirdwx.qlogo.cn/mmopen/vi_32/kUsyb796TRWibeicLmBhVLwdDyIu0VYicCeSj55PQ4tAoeM7hfwJlLorbNcCn5Wqt0nLIiaoBzibLTZv7pc9AWiaTjLQ/132'),(130,'http://thirdwx.qlogo.cn/mmopen/vi_32/UsoUwsaPTIRtsFHDgufYBSDXuDXN3MibBuokjoEGxVyK8Eg2vValkP981HgFUS8b50Moe16zK7ceeMbFzY0jbFA/132'),(131,'http://thirdwx.qlogo.cn/mmopen/vi_32/aTyKHBL0tQEYOSGcE9FPX76zKicNby5VUTzpEKwkE4Vn8VUF0mOhXEu7ibjwtQKX3fGWP9icSeX7nYscmz7Eo4vfw/132'),(133,'http://thirdwx.qlogo.cn/mmopen/vi_32/2RpVI6pp4nZVBHVBicIdbyiaCNvuib1BsN5eosRo0NbFnic9BkM6nicnyUnajCUhzxkj7UqBOyj0Xic47Ba7ob4V7oyg/132'),(134,'http://thirdwx.qlogo.cn/mmopen/vi_32/WfibIFib7ToSicF8F6VoCD2qy81PoXhtHZ6KnibanQecsRx3SOiapFXME1ujWfzxVj7Ijfp4A0zMBaKqSAK2PjrslgQ/132'),(135,'http://thirdwx.qlogo.cn/mmopen/vi_32/mDrA2JFVTib7zfjvM3l1sZAib1PSDrgI3TDy3tuhofAM0zj9CJWS6TKicPUwu2FdPYpJZr5hJl1GdEBz7be7FIxPw/132'),(136,'http://thirdwx.qlogo.cn/mmopen/vi_32/ve2q95VsuIudceO7tibedsuQMcqETggqT1fRSk1eia8nWMBe3o9wqYiabkHstXjb6oBp8E5ibfvcw6wqq1Adticd7Ug/132'),(137,'http://thirdwx.qlogo.cn/mmopen/vi_32/AGYjjsBCksWnkOyich8CSpibOJPLKQ8l8hlRwD4Ria51GlepZYRBd8FRXmsnRVCuDA9rPibKELRPEoW1ms52zTMcFA/132'),(138,'http://thirdwx.qlogo.cn/mmopen/vi_32/gA21GibREhPsh8D1fjiaDDiaTYo5xbgzfkKz77aJ3BjQzu9vUficA74sIdoVJSPAKvd0etn7uEk1nCZqPjJDETL2IA/132'),(139,'http://thirdwx.qlogo.cn/mmopen/vi_32/xkn2Cj1Qb4ia3xdSGulDgNFGNhMDp19biawL76HfWOQgF7B6Oj4hS1nHmT5mz42AbRL2yQLiblyd1uH4XUw5rcYNg/132'),(140,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJvOvv0EK3jXiccc0xeYsMuQ3cSbdLQiceDNqdIDFcT135ibwtv0yr6ic7zsM39FCJTtPQJK9IfjiavXmw/132'),(141,'http://thirdwx.qlogo.cn/mmopen/vi_32/DWxOrLcdgR1hzz3ZCPaS18icbosmQYLTksIGPjtCHGe9MlotuKOl1mRJ8OEsfzO15muljrDRLgULJaCcz3Xib9ww/132'),(142,'http://thirdwx.qlogo.cn/mmopen/vi_32/Ih9XPzJBYd9wZrn45ChEtfibH0gCq9tjiaJnBZbEvCy8EEkqRr0aahMtkl87Z0IwO8Wiacb7ZKWX3cs7JNdrEgcoA/132'),(143,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epVKZdpcfPAfLS8ktdlx29ibbVxLBicicLYGpPnLFZ06Mhhc6YHQTtJkZlHEv9kSwTPBd0KO1icR3QWlg/132'),(144,'http://thirdwx.qlogo.cn/mmopen/vi_32/y2fSgcKaBzJFhqibe0epxxvhG4p0iaMJYtAXMWV6hOSbNArdLKAOlW2FZMxLWXCrVXAJJ549ZSnvOgkMqA6icfOOg/132'),(145,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erMoLzHzwWpxcibOuibhzGJOGHictZlrgQPolqBEDClkkVMiax6kxxGxHFia2lyECNY5DHpM8C7fu5VdaQ/132'),(146,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqmGc2DEugETIobx8VtcYiaSXCXFveKic71bN8cJWYG25pUEue69mKzOkOv0xjriaze5aOOmvzZ8XVxQ/132'),(147,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKQdlJmgN7I9wAEWAW4TlCTqfibLsFeFUTx3YMaqzWNq8MsykSpj9P9peKhvQVRSlajKxwkt7R40XA/132'),(148,'http://thirdwx.qlogo.cn/mmopen/vi_32/fVcfia1y4ASC5iaMhRtTaWwTeKjyHEN6kJpBECbUHd7kfyaIsQP5fxRXdibibBGxuMreuaZNYjPmj0mJ2Ehcw6E6og/132'),(149,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJnv4h4j5tWy8CItQQYibjtx825UT20VESOGvNKLUHfozcHGtFwe4Qw92iaVlDuwFGpjKFwQhxAwsibw/132'),(150,'http://thirdwx.qlogo.cn/mmopen/vi_32/d77uULtfOMcCy7hG6l5wbLbAlwC8Q2CWM69V6JyQicQgTVnmDnMD3TzLDiblbUU83EEFsbb4jb3btMP3qvPupfew/132'),(151,'http://thirdwx.qlogo.cn/mmopen/vi_32/VsrhoYCAOiaRDZSSNGJozuVnlKGicTaAQIRY1Cy3NgibIp1Z4a4G7fArEQ7WzpwYqgFoPMLFYCQAy9WbUBlG4W6KQ/132'),(152,'http://thirdwx.qlogo.cn/mmopen/vi_32/Hic7wPkKbMibdQQWvb2QaCj3MnicRibceCrZb8uPOSJgsW88yaxC079k3LyyLWia5Q3gtcldEoEh6UnHZRGrPtB9Nnw/132'),(153,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erWDXriajStd3LzIffrLahVLXuiaQy5PWuWOPILib7aPPgNiaxp60lhURUO21lgLjAWPoGhlhJib7CyYsw/132'),(154,'http://thirdwx.qlogo.cn/mmopen/vi_32/MG8dibtSbFCydG3EesBicQQKjXk5WrQmqopVHkUkLevntFAHkTAWibSpB4fibFz64P5tHkSzAknPAsoW6Uo7x7tWSA/132'),(155,'http://thirdwx.qlogo.cn/mmopen/vi_32/PQdDfb4ibeNU1M7C70viaYZBFdnIWBqhHOpBSWlHyC6VWR2Z670aCAEuOjxQhibuN1UwCicM1ddV0EqFeb5ibgiaCrDQ/132'),(156,'http://thirdwx.qlogo.cn/mmopen/vi_32/Rgic7sMsr3TKQhRnzNib9KShdXtwTdG82sKnwR3IJekYkmklbd14WiaRE4khhurwXibafQvjZNo8aSqXVCxoibTNlvQ/132'),(157,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erDpysty5AnNpQ1hsSxmv09SzBVmJdCyxo5R0pOcOfGyVRYGwSichYlEuwmEgCKC15SR1juVA34qkg/132'),(158,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erpOibcUMk0mUuXECIYSiaBAHsIibxibyoBY4HC4YxseCNQF4vOcRJzZOhTS3mJ48fQHemM1ibeKbL3aCA/132'),(159,'http://thirdwx.qlogo.cn/mmopen/vi_32/4DgXefgtM25hVTr33mo9ndBhA2VENHO63YQLyNejL6VzbuGBcjEjlJ10ibCQj7Eh7vLssaervuJh8ib3QtpwEQKw/132'),(160,'http://thirdwx.qlogo.cn/mmopen/vi_32/zkuKOeu7zmFVhibULjztQdJicia0Q9gGArnXvvkIsNrOwhNxWN5UCPN6P9Oapfop6Ht4QcdfAFQG2ZPBCrV3FlZ7g/132'),(161,'http://thirdwx.qlogo.cn/mmopen/vi_32/nopKVT1CMEia2bzJXOQVK519yibzzUAIlNpnIXXaCA3qZJxd7GO0ymYvTCU16icGmtwbunmR2W5sz4jQunGNyIbEQ/132'),(162,'http://thirdwx.qlogo.cn/mmopen/vi_32/dZj08ckjV0JUfOLa4sibpAqtM5HdCwjog2BTK26A2GXadFZb4ZK9qBCgR6YBqVqWG3DnHRHRnaxA7wxzfYOmkkw/132'),(163,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoJEicNhAbz08avI4gcr4fSlOgGqoRcKw7UGZmmvTAial4JvWMnMbWVuJ8C8N3be7OuFz2XPTDP6fDQ/132'),(164,'http://thirdwx.qlogo.cn/mmopen/vi_32/drib6oDOX7YknXOvhZk2CNhDn7MtzFsNibBDNcAjsD8vEnFyMUxN3TyWZkIh68icGicicKlsWZmGDtlZWhx6skybWKw/132'),(165,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKovibrgjqfh7R02toROzxMoth8cAr0tVm4dl1y9zMuaP56E3ibGtVs60ECSWFnzogtic6L2ibRCzhbeg/132'),(166,'http://thirdwx.qlogo.cn/mmopen/vi_32/rFt97LxB1t2gPrI1wBZEQ3on7HkmrzfotPXNEAm8iaj4CjMTYiabks8UZcibk97uXaP2pmBaobjU74KhAGPPYicBug/132'),(167,'http://thirdwx.qlogo.cn/mmopen/vi_32/79DMe5szskxeDCVicDmLsyTQujjjEnAB0duzejdRdUp3VelvbxNibX3qOQzTqbcWb7yr1Zl4HdV3009hHKAw5QBQ/132'),(168,'http://thirdwx.qlogo.cn/mmopen/vi_32/P0vTtA3s0a0icQul5ASePcPILQtLV2dEt6LrvFMpZg7t4ibZ1UG5l8XwnBYv8PL0QhJpzKat4LxrRPESNibQTQmNg/132'),(169,'http://thirdwx.qlogo.cn/mmopen/vi_32/iaibZazqWjUILk7k8LrLQuvF223GLKo80PRHgoYFfRSm75tT5kwVZpBOD7FwtibxwVM2ILf4Q1EY5E7A1GOzOjdxA/132'),(170,'http://thirdwx.qlogo.cn/mmopen/vi_32/e81rOgic4PjTeM1J6ejST5W5ic722yy1FNmUBx6ibsJwhOE5AL7A0xmEDI6rGA8RJJCmatoMWETEufxNqORZygD3A/132'),(171,'http://thirdwx.qlogo.cn/mmopen/vi_32/3of3SUYAHS7BjEb5k8ib6EyYp0RPkXH6iavt3V3eRgiclv3IeI65W5sDVoFrddHKlDQicqVrUsbMRKBibRFMMibibXoYg/132'),(172,'http://thirdwx.qlogo.cn/mmopen/vi_32/tOusqnSufuBBSxIibVSBWwGwYnzXmrEqWaqTf8w5AuPZ0KdA2poUEWuLKOnAYjIrnsxdicE6IT8kgxVLEv5OoCmw/132'),(173,'http://thirdwx.qlogo.cn/mmopen/vi_32/p7AwUO1w70CToDBrFhh7zeCc5BrgLAT7hwpjhP8y4PgiceMZCmgQTJCibGSIYwtn2r6S5Zb1eibWrkcDc8v6MVOeA/132'),(174,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKKKybyzQ8OKMzgzz4diaD2oYzuxibDvrVxT5f15eDxP8gDc0n3n0S1WIckz02aUw2oopK7OVYIyFIg/132'),(175,'http://thirdwx.qlogo.cn/mmopen/vi_32/WFmPw6MMlwOU32QKIGj1Ypx9hoiatEdG3KPOfD9cAtVAt2J0XniaIjy2HdT5SOjibxqDGicFNBkRq8KPbSSRZ2GEkQ/132'),(176,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epc9o1HmU46tcsSNPwYUfFqFRMbBNiafJ1dDRCUyf468kHgnIibzOgbxpVm7UCfmfOswwAH5d4435rA/132'),(177,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKAXY2brIaKwfH1TSRdAomW0Os4zd8S8wib9Duq5osLJ4gNZF2OT12ib4FWTp4XF0BWd3xNa49bzIPw/132'),(178,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTI643alLPLydkt1UXS1iblict84ASobHJiaXqmXyTNBFoAmeUiatWnGy2MbotPTmtk1whkk1EiauU1p9oQ/132'),(179,'http://thirdwx.qlogo.cn/mmopen/vi_32/Rp2h30FJgyTnUibC7Z6KUXrISA9xSEs8kcicHzs7RDbkicPJA4xdJZ2EiaKEibCr8tjKEBmmsfzKic2aicwic48zKV4tww/132'),(180,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJicBk21hWeQ2RCibbBG74Ip0GQD4vgpmQp9E2KLUBmmjZuLdpvRFCTF1miacxy6LUoOZxWSjAxVpETw/132'),(181,'http://thirdwx.qlogo.cn/mmopen/vi_32/aCWtLLkwJpg8MNicKuPzPRq7S5zYUF79lJMrYgWEOczVic6oGqg7XfVh6FcQgiboe8TOpfBFGR6RZPpcCeVLcZauA/132'),(182,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJmUtE6n3uSAaS07MrIJRMGEhln8VxfFzfNyhz6XMIQBVYdy80VKBulT5qPyObFGSXhS3NotpEw4A/132'),(185,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqIge8r3PBicUgulOzVicUd0kFiauoiaiaSgpbaGuQ3zTHA0ksDs26Fglq75LnKDKicu6LVSS4alf0hPT6A/132'),(186,'http://thirdwx.qlogo.cn/mmopen/vi_32/u7EzgibbtV0o61zoyoicALwAGgfQ3toJnc2cTicFL3H77LiabgKFJMvOtl9geGEOC3de4iackKuK6NRaGR2FTOhmSww/132'),(187,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erkvfU9twW2oSqr8VibfhgmgDC0qGZnjPj3Qf4q99OeotuWlm4icfGxqO7gxrm6jyYrKAqDBVqmMnPg/132'),(188,'http://thirdwx.qlogo.cn/mmopen/vi_32/Mg64wACyxmgutoH5gDSmFvhnWAHSwH5ETyI2whXHHGRQibmR4eJHibb8Uj1e9gy5x6EHg2OlYJKNWoDjxO5q3RRA/132'),(189,'http://thirdwx.qlogo.cn/mmopen/vi_32/7JNAyofwhjF6RAhez1aqhntsHn0yQdTAxn9MfRHQnvrpIUk8XuMn4stfINzAvzL1OaQJiae3WIsCKcck8D2rdicg/132'),(190,'http://thirdwx.qlogo.cn/mmopen/vi_32/21Nc3iayu2DkPYZTOWpRILjpiaXZ95O9lDgAQaSs96MalxHbicIVVffY2tnKkomwdQc6ticiaSyTWa2icQA2ibFJhYNyQ/132'),(191,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIZUFiaoMYMjyHic9dmNiabktiaIVicd62wHqoNYKyXPqh9FZKFDc5PRdfcXUOOkz6VRY5JicriadHqSzVWQ/132'),(192,'http://thirdwx.qlogo.cn/mmopen/vi_32/Lu0976RicsV1N5RgLdMMyvtxffiaeDZLU8B6NvmVhchuSnbNc3icjCcYHty9SxJdUibE8Z5MHMDG1axU4IribVJ6o2A/132'),(193,'http://thirdwx.qlogo.cn/mmopen/vi_32/zfNdGQdY6Ie0iclRkTU8vplFdhmeT9NkDP1JUEptribK1OMgRKPmKZxrsVJrSxjibA7fL3l4ds8k8sqibrCTD4QNpQ/132'),(195,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLO5yYsOSNmlEljfS2aE660aKbA9ftdqfDOFv1XUfMweg1UQwuS7FgBKKpVyH9ciaTJYYcBWqR2rNA/132'),(196,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ermBl1GHRhv7lxF7K2om3iay8MK2liazjLcGgpmZg7U5teG78tbUX7uLSaVRzpQHVAnGrBwBR4BjbbQ/132'),(197,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLNYrxNUJar9mJF7K6heo4ia3UtMqC0ubHkS0BMolDD0ibnFgb9hQdHHNH2eGrYw0iaul6Re1zGN7TWg/132'),(198,'http://thirdwx.qlogo.cn/mmopen/vi_32/8aZxWBXDicESic7ZNNdolA43QspibW8vI5Rb9sxMGQTxQxrvG5a3yQgaKSOl6ATkSSNicEYe7ntqcFWibIeZzx3QV9w/132'),(199,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKSgAPIgPicPHZt219VzSh4v7sd4G9EwWU7drsicLfYEWp2ibicX3icJ9p9jVnoRfEVMwG2IicCDZvIIx8A/132'),(200,'http://thirdwx.qlogo.cn/mmopen/vi_32/fF7icBE5ZFnuQN3kxdqTrx65tibDj44ibGCH7ruQtpJ9lARFWVGIQVPEibgIoJObgdLoPFTNwSsOVGK7CA8icibjCPWg/132'),(201,'http://thirdwx.qlogo.cn/mmopen/vi_32/UBicVUlOZib7Ta1SEFvXDpRCL5QCSkFBk1BXD0W9RK5ryy1uAqib5ZE9icYH2EVeYaNuUmenYC1qK4Vib9j87OicT3Eg/132'),(202,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTI61jAFJXGuHicLcl7xl1vjLjEyBL1w89KOqefqL9d3Nm8lj9q7kL6j8fKqhbExOX8nQtU1KHNU13g/132'),(203,'http://thirdwx.qlogo.cn/mmopen/vi_32/d6YPf2O7pEFb4AjDwiabhfWQkG2XexHBribibaSjNcbj1kQAu63U28AgBmHgfIw1HaWk9WwjZ5ALMv67kIwS35vNw/132'),(204,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK5uMW9Pn4PoCmmibmnLr2DdibH249ZU5Ijov6aK9IGmTdt51LdQ6tNibB78tLicGSKGEJqyEktzmScOQ/132'),(205,'http://thirdwx.qlogo.cn/mmopen/vi_32/AE4uicgrbBlVXBwJeDEogTM0dy7iaLJ4BfSJjxTX8vcT8NATfwIBx562SaCjSVE9Xiaib14ZpxYzKQcAoBPHIIjXzA/132'),(206,'http://thirdwx.qlogo.cn/mmopen/vi_32/0eeD7OItdo29tJSssViayj4US4WEmxADKrW3IaemJ9jv2oibaZ1HTYJSbM7oIZ5fibMLo9IiadT1EZibyf0LfcHdO8w/132'),(207,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJYXopQ9OQjVa8zh5bvNy3R1tFLuhdX3pvbYVDNt6T5qNfPJlsGXatOTEgeotTOc44QoGcYuZmcIQ/132'),(208,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKUgeDsk1C7TkPxHyc9iaenSLbxs9pHpQ0fiat7BZBenuD9tePEcD8WKrOGTsgGJL4fD3QB44tBJGug/132'),(209,'http://thirdwx.qlogo.cn/mmopen/vi_32/lbpZI5quY0ibkicqc9DrBplibDuNbnyDl2bsTNuv5IrNNkJpXOvj08MqQqez6icLD9yf1aiad5YYhxT2EmULtcGibWSA/132'),(210,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibZpGtIBiaQ8siboJA9gibP7fk9Ml088Q2ia3zdBIibMiaScvpCVspXG18uIQRXosPG1hrQiaxBFG525O6yzkgDL339vdQ/132'),(211,'http://thirdwx.qlogo.cn/mmopen/vi_32/wTBD7GbNDbSGzkhrXgINLWNP8nhdDiajSxHnAkAS0icFM9xWSib0PoHhrFzNFUyDAsgQej3eQEuVAkTeAIZbkA5QQ/132'),(212,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKbVPzJ2iab2mh4D3ALAbNwaFs2X7Uj3uE6wOTzY87ZEQ1OwXuAT8YwTMv5yd41oogCwfH2BtQXtBw/132'),(213,'http://thirdwx.qlogo.cn/mmopen/vi_32/0w9DAnn4Ydv43b6K7OHIOeYlEQicStZ0JXiciaJOJyeQJSeIW2samZaIDPfy0ObibaXqpZQNPUcEjtNIBzicXLrhL0Q/132'),(214,'http://thirdwx.qlogo.cn/mmopen/vi_32/fMo0YMCeSzwFkhbTX3nBJN1fXLEMA6DrrcWYqKS60E79dVxxTIpSJRI5H46GTLIJbTvDpAkH2icH29ZUd1BEk2w/132'),(215,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTL8nR7kBpJIFaEOKrSOqk4uBO0tzeGKGsUZomkdFhWFSjbPnKpxuianWtUd00GET55xExYZVwK0J3A/132'),(217,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erGP7iaAkPYdZThicZmaOQgzbwrmjiaibnqG71hAHsgedAib5Iz2j9kTtDI8oA5eKTAG6bF4r2T38EtWHw/132'),(218,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqCLKAwJWsNRJicPWRyuzicNB64YgdzREazYnIXFaeXHxdhn6wfj89NXg9QhfDKVPzn8Qy6x69HhNZA/132'),(219,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJibwhxFRFAGUBdZ4gejRrLQicP4FgZ8JhByt2icfkJkFszibRKRPmZI1UlNuAmk7d9bWvTzfCyIpjUZA/132'),(220,'http://thirdwx.qlogo.cn/mmopen/vi_32/b3XCB1kvXvskuZmbucX61ZTUXOB2eu2lwmGZsMyZk23pypyvOFz2Np0YxtNI01m1CCyibyxJfQZK7S3Ph0EaIhA/132'),(221,'http://thirdwx.qlogo.cn/mmopen/vi_32/GmyqJEkCFcPtoRQic2Eu0YgackHPROnusrn58XMetVdewwHE05HS2BmKKlVzRXPK1VbwWogslI1chD1M1UYJy4w/132'),(222,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKDran5ruSkebSPuwoiamY6sicjHb4oftOM00CsTFMbLtaQt5PpHLx7ib0OHocUcz9wkNMpP6xv5dPVg/132'),(223,'http://thirdwx.qlogo.cn/mmopen/vi_32/Z6XlHHTohpT4UcDzmFyIhbe020TcBO1m85H7ibrMC2Y5efTGYNVo6SOe2oSzgUFAkSl5w4l9sUFDyZyqicwyDibyg/132'),(224,'http://thirdwx.qlogo.cn/mmopen/vi_32/t7ZF0KpibFibDWt6Kf7ORQT0jwFjaxcOd08wfJkmd5lfvPKVldVB6qK0uC8buXWyH2xYibfpNXw74oluhYmGqquicQ/132'),(225,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJp0JzSialEHSBSBTBY8lPQh0lNhAJHbse1UeEjhlC4vVPJRRg6DbliaibBe8qFKvibLtd4y7ONpmwrKg/132'),(226,'http://thirdwx.qlogo.cn/mmopen/vi_32/pVSxzicoxgJvuyAPCrUVq7MGA67ZicsibzZPDZeFDWb4lMOjUvZBaexc8GtR8s0t8TiaAYDxHeFBviaGKBERXZW37Yg/132'),(227,'http://thirdwx.qlogo.cn/mmopen/vi_32/fhw1kXXdpyEF9xAwXYWicqtDS4Tj5VJeaNFU7NBOXgkQhkuhP42PKdKNlVE5MibCKDib9q6FqtXSJArydG7Y3gqHw/132'),(228,'http://thirdwx.qlogo.cn/mmopen/vi_32/FnAia4oNWqvshsFsRyme4CbNczZ3ZCbPQjXMsnZ8tictgAmzOGhEaQ40Gh0Il3kzLYuTp951EB0PlDQxvichWzXmA/132'),(229,'http://thirdwx.qlogo.cn/mmopen/vi_32/TVCk2MAggd5l4MrjKXGsKkHICJGo8eTeFiavlFgvxbJpdM8kunlJXwBiasCLIUAgkBuXBDp720zBicUNA6TibnKs7A/132'),(230,'http://thirdwx.qlogo.cn/mmopen/vi_32/OEEibIp8btr4m7X4ANpWiaVLGInicTrcr5zscuDJAqnUrkRiaLtnL9OaqqdD5kF1AhXmvTvX0A4MXeYGHjMusfDgTg/132'),(231,'http://thirdwx.qlogo.cn/mmopen/vi_32/seiaKVnrGiagNty2iaoTKgT0icrnGs5zsh3JiaLpJLqb7IS6nmJibulfibBlOibYoMmQMGV3CX0ObWqOw9PicGfmrcjh1qA/132'),(232,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJZp2K63d3IPFibbqibhN5icC8S4r75m4QBDuvasdt2dP6JMN8vybISibVE397MxIYtOt7VEjz1sW1zhg/132'),(233,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epn3Bu3p5Qux6B3kMtfia9TA9ngfZ4cgT4JayYxhWWxRYcUUEQTtkjyoibXF3uLCqKY4paSwJ0M9AeA/132'),(234,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIsAzkEtj99kj8LibfTfheXl8FtWN7krDOz6KhsWV0BFmln803BRlobG3nKQhODSUCLwRuicJRzaAxA/132'),(235,'http://thirdwx.qlogo.cn/mmopen/vi_32/bmmAtCsSEdaLWftjVjELUCExYCzPbI5D9BZcamXDqibtuvJiaiaNCPicLjLkBUeZacDVE098Mm09bs27MFbjRqmuaQ/132'),(236,'http://thirdwx.qlogo.cn/mmopen/vi_32/623ZEbAgRZe3100PjvGutP9LCdTEua0IdVmmLAaiaPjGmstscjVqX8DFZlraGtjAdI8TWEzvWceT2rPs1XjV00A/132'),(237,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eq41NH5aTp1C9zRGLSE11vWRd1s0ttTaJJJxiaXxBFzSWD5Jhcl3pAeianDiayYrGWNicb6tRdYC5m46Q/132'),(238,'http://thirdwx.qlogo.cn/mmopen/vi_32/AGYjjsBCksWnkOyich8CSpibOJPLKQ8l8hlRwD4Ria51GlepZYRBd8FRXmsnRVCuDA9rPibKELRPEoW1ms52zTMcFA/132'),(239,'http://thirdwx.qlogo.cn/mmopen/vi_32/Ria1H2nHh5b5sobSiak05fcyyVjHFwO2ZFgXbjyibImx6rn0OGsKERJC1SyciajE1WiaoLMn9oKzvicOb1kCEQ0pyOzg/132'),(240,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJHb2F9p8yC8PUfpFVd7ibpgDGowRiaEN306HiaSicqJlosIU9ugiaEuoQB0OgC2Hv273lfcEjsCiaI26WQ/132'),(241,'http://thirdwx.qlogo.cn/mmopen/vi_32/zKa5tkYWEmG2DdY9zOUCkhYDQZQkes2BePDqK6CYNCFXX25IEJ4WPkVib5LM2tryRh91twTWYaeJvuN4q50Kj7Q/132'),(242,'http://thirdwx.qlogo.cn/mmopen/vi_32/spxR3zcTh4Npic6ciaZAjxJHIbRlzeQDedibLduZj2x1VkBoG0VvTbZiaw7wWMmVLP4obHFbjncDDJmSON8FlOVblA/132'),(243,'http://thirdwx.qlogo.cn/mmopen/vi_32/ajNVdqHZLLB4iatkL6DnkLb00WPTq7FMZpLEtEmZNWpsib0suFUF75iaYjxCLQZNGp1k8YWAt6mW4IX9BLPYB55WQ/132'),(244,'http://thirdwx.qlogo.cn/mmopen/vi_32/iaUyongIyPglcyE0BU6j5ia7sttTbYfApbznTdOMUfNkfAkXqqxBcKsROP9pM4kuQ4Cp8e7e1qZZQnY6LSibbpuFw/132'),(245,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erpOibcUMk0mUuXECIYSiaBAHsIibxibyoBY4HC4YxseCNQF4vOcRJzZOhTS3mJ48fQHemM1ibeKbL3aCA/132'),(246,'http://thirdwx.qlogo.cn/mmopen/vi_32/XsTvvZYTZ0kiaYu84ekx0PqV9Maib7Vvwpzgo0jVng2b8tasqIXMqVB5Bu1hsM310EzL8nWpfQUTvbtTaEAgo4cg/132'),(247,'http://thirdwx.qlogo.cn/mmopen/vi_32/4RdOss6AW94mB9J22MmKBuCjuonjtWXTuA6632f4a1VsLVJbX1FakDSq2mG9UOWiaQMYZUCDzn5STibqOp4VnibnA/132'),(248,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJUYxvQUXvzgJMOicvKxs6xN7Uz3Ww5UsgicYB4ZAk711ricI1VVoMuibmMib9NCA1r2mE1WgFdXwgnT2g/132'),(249,'http://thirdwx.qlogo.cn/mmopen/vi_32/loVyFiaRialvIjY8Wlrica0mfp8cy7Ak8kYvWicD2NheszZEticof8LgLRzfFYQoibhfuFESricOq4Zsic3xecG3PWQicicw/132'),(250,'http://thirdwx.qlogo.cn/mmopen/vi_32/sCJTzefc6f5zJ6vDyr4vb6yTuiceHE5DibSUdEpc7Q5brMgfjlMB6mL0Y25jpYZ0zwUdNLmT6LoYMxoUUCjIxAHw/132'),(251,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJtdNoovBpLtlgicGSbnLcJiak8WWYmNUO7rsEIPBF5DEicibVorRaA5kTyZZE4arLpLWKgeiceLVUSTog/132'),(252,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoNFv6gVJApyZv49Iv29s3d3FvOibRrdlmBBiceFFS5EFeZWwZT5yv7JIictxQm1DYViaMcX3GIS1ibiaKg/132'),(253,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLfrbMvhKQYh380J7T4dlPGCpdY6gSKOPKBBgMxguhZbkaG3TqGf3Une1BibVx3NhfJabibJltbxNBw/132'),(254,'http://thirdwx.qlogo.cn/mmopen/vi_32/VsrhoYCAOiaRDZSSNGJozuVnlKGicTaAQIRY1Cy3NgibIp1Z4a4G7fArEQ7WzpwYqgFoPMLFYCQAy9WbUBlG4W6KQ/132'),(255,'http://thirdwx.qlogo.cn/mmopen/vi_32/R7wx1jz4EQ2P2noKA9jCBAicV7P1pzhmM2iaEibQqyfR5uBcoMoTK4fKic3bUnSvJx66XKNllIAteO86SSqaA0Tj7g/132'),(256,'http://thirdwx.qlogo.cn/mmopen/vi_32/7OrAKK3DscXBkeicbibQNmf1XQOXTLd6pibx6arDPwoOGwpUhbmrhNvJtVGz7iaibibYyTTwUUHDBff0pRk87ibia5xQEQ/132'),(257,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJtS4Tb8iaPDZwzrX53hHAzDXHESADLhd4FOq2sjtTTSIaqgcTymmBqZiax3akQ0C0ckTC44n7fxQcw/132'),(258,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIgDicRj9LAvWbfCykj1jicC9pRicPb3u8eOV9p9l3lXgWTzSzhuL7Momhen67HxsUmibySsQSBCQ1knA/132'),(259,'http://thirdwx.qlogo.cn/mmopen/vi_32/nribDiatSjlg1zsn8icILQmHglPo74OdaDoNlKsiaRqXtgoibGUc0UT9d8I6Vqs2cQVTDFGsJ9h4w54wazmRTZ1tibng/132'),(261,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erzPLg86hKDTKCPZWK0wcSmrnP7fNLIxa6rrcRJTb4g7pA3N5TDK7GNHfOt9Mziayj8LxQHJlLrCrg/132'),(262,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJV2PU4QW5xMeN5BCww9jq0ZSicKfHJicpbpJjNATV8sxcypwhnsBSZIGcAbDnX1ReLWEe9UG6uEGww/132'),(264,'http://thirdwx.qlogo.cn/mmopen/vi_32/ju46qadn02dSj1sodcuNFhoUW35Ht4g6yyZSoxOL2fMQjInPEm2P8r33oX5ibsZIAEn0kTePtnWo0vOibklHVAAg/132'),(266,'http://thirdwx.qlogo.cn/mmopen/vi_32/JLFyKBwQpFVd7bsGVIV3ebstru1giaibXKl36LskKgjibmgFVPtzUD15H057yS5mlx1myRCmstxXcjZTotWoD3ic8g/132'),(267,'http://thirdwx.qlogo.cn/mmopen/vi_32/uNltEvlzQFiaquzWl4iad39IFwjE0cdSweCkGpbiaTOuo8axaoQV12ZXJC4sTdlPAKUBYVHg6swgR9SyXk1jQw6qg/132'),(270,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJVJor0bmCEibHzXL1grsrmoL02VkictN1XF1iciaaxWAPiaUx6OpxlxluGdKzxKlLGGicvk091WXX7YjHA/132'),(271,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibUHcArzu3ibpUEyiayHTEIs8JqN6O8BJ3d1SsW34PticqApfG38hQwFMYFok9QtZmyvhwicLDLNJRDnBp5U7CHiaoOQ/132'),(272,'http://thirdwx.qlogo.cn/mmopen/vi_32/d5oPlvPIZgZJqjHvhATNo82gbkGKerLUOvGlsva8R7ayKnxz5XqMCx1l1mdXc7JF6KfbeeniaBnwQrfibsS61SSA/132'),(273,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erkvfU9twW2oSqr8VibfhgmgDC0qGZnjPj3Qf4q99OeotuWlm4icfGxqO7gxrm6jyYrKAqDBVqmMnPg/132'),(274,'http://thirdwx.qlogo.cn/mmopen/vi_32/ic2PcAKfUe7rWupEzKPPDb6gjzbaw9pBJZbic2zk4ibFZ3CXtfKP6px1lE2I5vBkaQB809m5FovhzxAcC7kWB5icdg/132'),(275,'http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaEKBENQekdV3e5zicic4OGW8icX3QbmPKpoqnBtDPdFC8ODZMJNEDoL8Qkws4UXB78iaicsXdTQUvP5kZicA/132'),(276,'http://thirdwx.qlogo.cn/mmopen/vi_32/8wviaS5q43dHibmiatHhCMqFtRnp0HJwdf4Kysaib8icLFQeCsFLShuJUFnEWGe2rFMlAL6tTajNkySfVhWd92OXicpg/132'),(277,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqAXcbicN2gQux4jHV15eNbq5GHD2YzQa6Pk0PcU66untRly5xKtYCOpmzyiauiaIAeH4Wtib6GUHXKCQ/132'),(278,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKkn1gUjlnRNwZIye8xDEPwRKG9KHZm7pGdIicccKGwjibCyzkNP7wibmUceMicK1D2ta7JK37JWmmCIQ/132'),(279,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKDran5ruSkebSPuwoiamY6sicjHb4oftOM00CsTFMbLtaQt5PpHLx7ib0OHocUcz9wkNMpP6xv5dPVg/132'),(280,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIicy7icDpyoJMI7I5beQ3BUUAQvibnZEAGJxbMpHEELCwxwcfn6m5bb8Ff4YuTFibfpBuIdowLQsGic6A/132'),(281,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTI8ricZ4yIgSmYCNFibAQW7XT5ic9SJYzTmEJUS4OegJshfIichjC5fJiaK46bg5ibdxruvz4ybkicfUXgnA/132'),(282,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJR0FvmMkb4ehZlSiaRicHRxvur0acRpLkqRicqUfsibrGl6Dicib7QdAeUm7GVsyW67sckwLhDBIHz6h5w/132'),(283,'http://thirdwx.qlogo.cn/mmopen/vi_32/icjdkEFketrQaM3PbhqMfWL1SdygJcXw5oGOqDVmOd7n4xkWjtib8FGibMibics3icjeq2nicibiaWjfKibogkf9ehBOM5xQ/132'),(284,'http://thirdwx.qlogo.cn/mmopen/vi_32/PUjQaIStY5PoRUo2ubp9Old7zYT9wdm0vGZC5Hy9appXM2Al8nsFS3a6yTdW0e3yuXNcX8a0HOIuYZRpocDjzg/132'),(285,'http://thirdwx.qlogo.cn/mmopen/vi_32/n0ZlKEibGGGvqCxKUnPMlicF3216BibI0iannNiboNWl4NSmd3CA6lX0gvtqF15Sj3tI9aVarYceKwICXUTx2narJhw/132'),(286,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKAXY2brIaKwfH1TSRdAomW0Os4zd8S8wib9Duq5osLJ4gNZF2OT12ib4FWTp4XF0BWd3xNa49bzIPw/132'),(287,'http://thirdwx.qlogo.cn/mmopen/vi_32/ib4mUbE8A5cP49E55WG7nCel7LGOuMERtgaiawRdslnQXLrvKhFvictODHIl5OpdFKOXKiaOJdj7n9ob8QPUBeicbkQ/132'),(288,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoc8sReia18fc79thIzx8orSJs7ICS6bWe0vDy6jh9vldTtSlEbsn9O2utIyIDsGHd6wicKAnt4m05g/132'),(289,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ6dSAon2tWBW6LnXolbzKFjC6aOibWPZ8cw8CUq0ibWag7iaBrP6ls1Q55QOmjdxyeqBLiadwzR5A4yA/132'),(290,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK6JRpRMJLHZpIeia8I0TibxXuNzHDibaVk9yLy4hKXR4OFUlPvgCILAUmqKC2orN3TEzv1K2gINTHsA/132'),(291,'http://thirdwx.qlogo.cn/mmopen/vi_32/FnAia4oNWqvshsFsRyme4CbNczZ3ZCbPQjXMsnZ8tictgAmzOGhEaQ40Gh0Il3kzLYuTp951EB0PlDQxvichWzXmA/132'),(292,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLnxjgOV4zxSIBt7RIYI42L6eia9P8wmQTD72tmeQWD7BSKFfkSwB0gw1fmMls7YGTgU42reALWdibQ/132'),(293,'http://thirdwx.qlogo.cn/mmopen/vi_32/F0Z2PfbHxlzbx4YABGvMaeqBgTr3dSibqcFjM4N1YcUudIpkicic3rYaVeqKXtHNTrVzpFt7sMSLhfysoNHCQtRXg/132'),(294,'http://thirdwx.qlogo.cn/mmopen/vi_32/ugBoV3VerYOXRI2SZbtfNA5qStc2fmMMeyeFibOgkAlWz8AsXoWSaFwF5bOwFazc7iaA2ZJPR3mfbniaRtyrT9sng/132'),(295,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIYGTjVzw0sOQHEcD5NKRlfvunW0ky9W11icGuVicib45xd1ZkUCKXcmibLYXTvnoCEyyWEDBWdQYNwcQ/132'),(296,'http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaELQDoCsxO8vZ0z7wbnPO1NOn8FibratTQAGaPucRYVqCUBAGr3uWGrWiczicyLpWlibsUC330wC5LicJyQ/132'),(297,'http://thirdwx.qlogo.cn/mmopen/vi_32/5I48SZF3GINibduVSRpgibe1No5IxVNt5kqjfXouc5YT8BOfEibXanhT23xcicSozTeQXPdiaOB32HF6KfeK9kqHCgQ/132'),(298,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIY1iaOVYkrLvOyJn7yotN0aFHCeic6vegiaRWBX64xTyb4eBQ6F6ZtiaA5gE1LM4Y8Zw0GECpiaWvv5uA/132'),(299,'http://thirdwx.qlogo.cn/mmopen/vi_32/rQOn22bNV0kr6sicFWFOYrUuWAwKJH69rPnkTeeQ0baAcljPSaXCy0sjic8kJKtDju2mHUcoVx1gfn8YT5XFSlIw/132'),(300,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ercP6CicF6jtLedYGXWouW9hDcVYCarWA0vquezPF1ZtzLhUMsHHCb0KNEibOnlj1PHlWPiaDP5EQiayA/132'),(301,'http://thirdwx.qlogo.cn/mmopen/vi_32/jKVBpmgWtFSngnic8yJsibzfH6Y4sqhetlI3EnBGRZYtcZXuyZs5WBm0fuoTxBeAoa4c2qkyuNNKiasTcZ4Xia3exQ/132'),(302,'http://thirdwx.qlogo.cn/mmopen/vi_32/lkZliaQoMbJibLvVCe8oaSK1lwhicu46koibxlJHZ7UBwP6iarIadMlqoib7NyiafVBjn6ynPKg43MP0PpmU6XdzzfNsA/132'),(303,'http://thirdwx.qlogo.cn/mmopen/vi_32/Cp6oMCj4n8cjPh1OgenA8zFf4fBPLVQ2WwBjicC1bl2rMbuJt9AyOEsTwcwcLkLEDURbJibic4ep3knp1yaBkJvVQ/132'),(304,'http://thirdwx.qlogo.cn/mmopen/vi_32/3RgiccjdnWJDKkfiaf78dDoMWnwZy36KRl3QMCsSb1JPMt4KDsce2fvp4WO6ibpjbqT9S7HKJicWuYBh76qNtiadRRA/132'),(305,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ep48ibf0qzEOGz35wLoqynILACJcSPKYXZrF7TL8NrsmNeHyh6rIyT85oHINYN0OjhOOa0osQZoAKQ/132'),(306,'http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaEIiajH5R2vjMpNHWrU5wqzUNJefiaqYLaiaTusbzv81SC1CHxswskW6bQ4ExzvVFtbapiayt0Zn0iaeaOQ/132'),(307,'http://thirdwx.qlogo.cn/mmopen/vi_32/j7R1I8LzLicpicMZOCnOP6A2gr6A7MrEpJEaz8qPibcU3D1OM8Bgb6XRSvwxWj09HzSwkic1eicYvNup6VicKcNFic2UA/132'),(308,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erglUjTQZuJWo9S5ibSlrQZTjiaaVaia8emS3TXGmYH3PicbOHog5tYKogEJIgHXnt6gYA592D6Wk5SmA/132'),(309,'http://thirdwx.qlogo.cn/mmopen/vi_32/H0CjOmJIdsd0VcWEibyQ9mTug5wugXMY3a5QnibFInFFxSXTJzJECCibW1Aic5AeNbhOuqGCQyNfnJibVHmY9CUesaw/132'),(310,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eo4Ticyn81ZjRMU5dkLQ8ApgWzkAwl9Lq30ElbDSIicNNlcjKAacy5JlCo71WcnRhRiaewNicK5icku5CA/132'),(311,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLqqHPKKz5EFiabSpIAd1dMjdZbtcxTuPZic9IibgkASPI10mcNMTSu6HbMjfBx03zEOibCnpwuPicNiblg/132'),(312,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q4TjjcZ6cd9zrurpPZyd7MsjJj5w4cQUcVQLb9AiaOGrruxVBemxbmzS6woyIG0V6tTyLH2Bk6WlElPJWVJoGZQ/132'),(313,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLns60AgRWAiad4SVjSLVPJwaq8uicMMsaF1sctrqkFzDcSADWrdvGxmSIdibp5ib0jKMCF9wQQDlX47g/132'),(314,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTITQCicTqU4wdHJc1ZOgQWnnuaIwiae4CfQqBwDPw1EHVJcW1D87eIptudt3nhVEd8J6XXgibsreicZUw/132'),(315,'http://thirdwx.qlogo.cn/mmopen/vi_32/Hn6sbEaJ35kibJgL0beUk6SgkH1ibicgbIzhaq18apk4V5YHAj2FIUhawf7R7mmRuqDAuTJrCbu9oAuXbxIy1MHKQ/132'),(316,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKSHc5u8Uhpcy5xQLkgwIABwAmXPWGl753v3cE1n23mPR0Z2DPzV0rVJ1pYBHjD75HMpVGx7Nnia1A/132'),(317,'http://thirdwx.qlogo.cn/mmopen/vi_32/sJ6Hhicib6NdskJYYkzauEo5dGxZB5IbDEzwGdPKB2G646SDYicXYToK1806LibDXljgnN12qHMDMiasicf5AxFwneLQ/132'),(318,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLJUicgP739INuVMQmclq9rpwibK4y1lZ6q8BBLfDbVcj21qpyuKUcIKDT48UL581p6iasbnsYll7KGQ/132'),(319,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJxiaKKtg60fDib8E8DjKJdOzM0hYumIFoKYlr0uK18lUop8DUu7hEyyo7Nw8s7eeLV7NcVYkWlyfrA/132'),(320,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eobMvNqQVpt8rQNpQmYV5IiaFcoU8a9Q1FL3pX3KTNibdbJHq2rO0Elvke2MvGvrzc0y3mjCfWN5vicw/132'),(321,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK216wGiaozibxeQCQMH0WDQFTDVm62FtyyBy6DffuRXucuQlOOssJKia9BevqnMo7cB9Ux9jXxICEGg/132'),(323,'http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaEIEtdElNicmpODBZic6UQRqBa2LynVpvkuiaiamHU1n9eB5eY3s5TzXQ1yvcczLo7ibTMRFnVyaAFcWOibA/132'),(324,'http://thirdwx.qlogo.cn/mmopen/vi_32/2ONHiatITMbWAjFr06odBx94dXibcdhaXQXhoxH3cTP8DUtgUW3iaARUtWDWLMOXT2sQ39Xkcp0W5nkwV0UOeLWRg/132'),(325,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJz4XOEMBowRsRAxvoo8EliaXeof8wscIaQdev1QlZLQqticsGUkZrolNwYscmViaE1mNSjk1QcLQOyA/132'),(326,'http://thirdwx.qlogo.cn/mmopen/vi_32/THa3U7iaCDf63tAdL2bZA0OCfdp6LJAPAP5GOqribUgzGGs3N0GRWbgribicaaIYBRKMk66CNcklDHmcibPXICbCYZA/132'),(327,'http://thirdwx.qlogo.cn/mmopen/vi_32/gWlLTgnfq2BQU9fG0zYsic82iaVhVuHkAwWs19Nkia83xp4TWHKBe4XuRic1iad2QQB1TNzJ7Pn8ZwWWCAlhibnEypdA/132'),(328,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoSoSicJGujTjlDsvhEIaoBZSnu4bzsdVxmNfXJPYuL2OvFUyOYh3gX8mKysnWML86lsicjIenpuMrw/132'),(329,'http://thirdwx.qlogo.cn/mmopen/vi_32/kq1tvhNib6icLibZxvOZaYRg6J5sjyCFvlwm7QVSKPnjZm2TqDpENGZf3BhGev2RPyC70jISwL4qp3EG7O9Y5l1icA/132'),(330,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eq0icoHSN71mWiaHSEu2WcGUGTZGhAs1kK5WDdRAjc2ibXCoUGnMBUEqWTAGib4hnJucQnGyW4mI7EoBA/132'),(331,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqVtI8M9EXo6DhuUxYGsrHuEMHWQzF8iaPUIZmUjpmmRvDO2E9pGRHfboaXpQVzlLr6cicibABQZiabTw/132'),(332,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqvgWAp2B7xLesunysTdzibFFbgziaicib5ueCb3MOJMOQFjBolEhsekpqfUqdCHB1rbeYAYFujaq4wqA/132'),(333,'http://thirdwx.qlogo.cn/mmopen/vi_32/FW5he0nv0CwevibY8pme2jtuAlv7ur28icZ3YrLjC0y2lBkfOCD7poZ6JyzGAOlcHj19ibStf16DTA6CfVfjMkpAg/132'),(334,'http://thirdwx.qlogo.cn/mmopen/vi_32/tYmKURJG9TtsVzLPZpsah71HWdafDicYz70OfZ9ZbcqAwicv0T6Rm5VcsYNYQ4fdiaxjh7znBuqY4tHM4HyVsBPag/132'),(335,'http://thirdwx.qlogo.cn/mmopen/vi_32/kqicIkIbWGHEO5Ju37q3PN6iayrpelNZRcMt57QNYNHR6ZFBVhldkym7wkqfVnIFo60KhcjW07vnsBZS7AHlSaeA/132'),(336,'http://thirdwx.qlogo.cn/mmopen/vi_32/9vvKMsgWTFwwdnrV3n17IJXuytj4Iib3ydCwtIgY0luPCh4icTb4bHYPO89kUsib4oAicjfGsGA0hibeVNTOazhiadxQ/132'),(337,'http://thirdwx.qlogo.cn/mmopen/vi_32/UnuJkUibyaYAUxVn84hjQFVicDOyPl8p963UrmjBbBEyHEZ4uRKAngHPVDfhJDbColGNeZ3beIsRXia1ibtkeHxKOA/132'),(338,'http://thirdwx.qlogo.cn/mmopen/vi_32/SIJxEcpKygFdESxBDQBUibIvV7RhKhyOe8WbXKNHWqtBqzlxoicerVROM9LRVq0vzL4qOhDArRqGfL86zEtFouLQ/132'),(339,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLqheicXFEJ48YsyzoOm2J5Mx0cJPVBsGqG59ldUu0BFIE9tNm0mVQBB2K0npcDos4a6d7hCHZ37pA/132'),(340,'http://thirdwx.qlogo.cn/mmopen/vi_32/pqQtK3FNnwIdaZ9BPteicSFMnqbKHdzCpP09KzfbJsNFbNX2WAPlOr2SOKp1TZgOzrCRPFS0iaP8CUt5ARHGMiatw/132'),(341,'http://thirdwx.qlogo.cn/mmopen/vi_32/2iagtqw7iaHzW6SzvMBJZjc5vmB8LPHweSz6sGJKG3XczWppVRQJlt5JZmDCasyOK1ibAZoB2xlDkkGXzvJsqLBUg/132'),(342,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKGiaLDlazWNmGdMTyY1iah84M6jkNWVLNexibibBSF2EhsTrXaCn6EyerMApnfg055qpfUAuL4oH9L4Q/132'),(343,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ0aica3JGr2n9Eqk8XAy7iaANaPx1cnR5iahy0cwCbMibHEUicia9Fm72E2kbunpQIIE0dQicZKtWcd7Jpg/132'),(345,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIvR70LhzRNoibL6pkhHOeQHicSLcpR09hY2DwMiaTKV7D9NmnIRNjKNQibUDsWLGtYoBUdicDYXLsoA0g/132'),(346,'http://thirdwx.qlogo.cn/mmopen/vi_32/ZHYm6KE9JxyYzzlqumncclZX8wo8M7ZJSzT1r4dRYQiadOmvUWfumc8J4ZiadFrIQb0jMKUBGcGVTFiaZVozic00VQ/132'),(347,'http://thirdwx.qlogo.cn/mmopen/vi_32/Z3AsrdliaNur2K1PVwicgzuCfA2ia1gD6AoXC5t0VnSQCfGbEuAufIvWp2GkcicraibGn2LoTcpqICGf8tibLS03WTNQ/132'),(348,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKibZ0bicHvMeXDT14SJtNPbkqlV2e3ymLCia5U6wobmFla0KvQSo3JqxE0sQrxTbwx780BN6PlkCklQ/132'),(349,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqcWicUnjAOQahWNNcay1acCnXLwyaRIYRtgIonIog13c7h3eK6fj4oK4zTHPELyOnLVqYibiaNMfllw/132'),(350,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJHHaTr2hrYjNibdIXuWM4md9Ps4goPBrXiafA3W8ricDWsLhic3D3lE0EvTsD13kMTC17HbfbPtrOnGA/132'),(351,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLWsWvF41GQ13Qz2CHsqSL05UjmovPqnlFCAEUia1ghv33LdDx7ibG5McQM2azkOfKU7OmOefM0LOAg/132'),(352,'http://thirdwx.qlogo.cn/mmopen/vi_32/SgibQMDUib98Ng26bG2k3ib7zIv3yMicIUwAmZ6QaZHvof8IiboUgRRdIYGib1bHvcLh6yF6cjU9qyeNruoECoKwibgHA/132'),(354,'http://thirdwx.qlogo.cn/mmopen/vi_32/H9SaYoRYKzzIhLlJb12TPLPcMGibWLK15HRYaDxPiaBx0ZnzNSYic6Suvt5IdvUkdmiaEBuvtibSf4CubIYoeJ7qHWw/132'),(355,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erLaN9qeEj70CnRGGW4MhUiaq3nHarU8piaOjhfCsMHQuy7G9ibQBFer2Ol4wbryZUZiamPx0icsqcmwgA/132'),(356,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLe00c7whq5qBfEJa1ibX1MvEFZeiadvKrXsSFhwuC9ypuITDmjhia18J9l9VH4fGH9OpHM7iasEEzBYg/132'),(357,'http://thirdwx.qlogo.cn/mmopen/vi_32/J2l8U5P52OPzlbbNeWE3u6lENkGmrt5TNmr3ialOibtoX0VXxYhsftQsxNlWQwJN0gQBfGwX0x9aOAjZcKFrcXaQ/132'),(358,'http://thirdwx.qlogo.cn/mmopen/vi_32/zwQUck9t9JW3qViaFZ5jEZ6dtu8vRsrYLFqdEtNIv2Ltxfxgbia5p5aTQdTibwBHLYVdoXGGuOVV1zibPfCBdbykVQ/132'),(359,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKMEKUjwLw9tvgQUiap8CqupQiacuTLdF1yR3OLAQicpWhGPWzuaIsF3kj5WBoL89Spt4upwFicUTJGdw/132'),(360,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKdN7M9b1uiawibhzuyRGennpE1yfKYCiaib2RM92aRdXnorP7sXrSGNpG5wkmnibzLrCraqc3lfK6ibG8w/132'),(362,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoC8XNwwQeZ5ZoE1z7sDn2p8DewQ7YP2YEXbWOavbx0V1ycKP1iafMXAiaMpoYZaym9bhlulGbmneuA/132'),(363,'http://thirdwx.qlogo.cn/mmopen/vi_32/XXvAic6jMmTicUCRIFY3pF1qqSYAQH5ArgWYDoWKSVwSw1TuP7yfxib9BOBhnDd8Mia5ibRz4GMurIiaOnbpHzptFKpA/132'),(364,'http://thirdwx.qlogo.cn/mmopen/vi_32/HuhDHAdgrvFbibC6BBAsVWyrtib4EYlSppnFicdIhBDCibh46XYmugZ0I7UUmajLwUY7l0Y9PAQZrIIFf6BOwk08qg/132'),(365,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJVVp1vegWFsbfoc4icibHuO1BaPibvLC701vzArPWZrb1icfQ6Ggzr2sZoiaTexa0ZJ7WYxiboCbpFfVBg/132'),(366,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIxPCR6O1ZYtScW0yib9USzetuHGYUicXohiciak4M9AMuDTuCUhjeDvib37SLFKGFJibUsvwoIxfufH6hg/132'),(367,'http://thirdwx.qlogo.cn/mmopen/vi_32/Tk2YVFwdoqJgAPQEVoB6TZyhBXXAJXviau2DUfjnHfQnhiaMPM8RTtf0qiaedE6x5EwnAwSCEZumv3kbsLOJ9BOug/132'),(368,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJuklibZZ1FsCAmjBcPVnrEg3SjSkfJdK15M5ytg5UOt6uMFM6RanqWoYIahhgmqicbxFZBRhrcCzibw/132'),(369,'http://thirdwx.qlogo.cn/mmopen/vi_32/qEqIbnu5lVL1ujNSHsZNuOIibkrzn0ucetHkKykCHHRJfvBoQqrdmN5rALW0cBEkNRpy9ctrK5KnzRSkqVQVAvA/132'),(370,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epUe4ysb690v7k1BjukvwOYWdiaaQ8pVhSAaHDsyJguUPibJyKmQYeMxKYdzBN35RhS02Wav7bkIC4A/132'),(371,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erZ3PmicAc73Wps1z6jp65lgRwTXg1BHVScDFxmLrPWkPv9r4f7bWHo3h4JG008lVBSsCGicst2sZnA/132'),(372,'http://thirdwx.qlogo.cn/mmopen/vi_32/9ibEpF2L8OAmwEYQ8yYCleRsXjMfFjCK1SraRqUfFAKry0eiaSVhwP3KHIpUb8GibFH5JdG023Gqwsck6rJeiaMNRw/132'),(373,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ergw7PXHzPWClYDZT00J16mqicfXnjQ69UdFJ3T2aU89NOj4LggBZVOUicevBSyrQ4J65LytdAIGTcA/132'),(374,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKAichLuVU7Hd0m54yqZbbPiaVlztxjKXJZzxHb9vZ1XkSqFdaGb7MkibnSKxrf1uKyVs7cAW39OcMOw/132'),(375,'http://thirdwx.qlogo.cn/mmopen/vi_32/5AXYPaoPLf1zhFbhTtEGLoAxhdg0UHKzJ6HyBhAibSfzVyK2xmlhZOzJwMTujjBSzhCVWWSgxmIHWlFneZcEbDA/132'),(376,'http://thirdwx.qlogo.cn/mmopen/vi_32/sfMKiaHzSHkIOicAsSQYcNJzccTAjsm3fY20rmAPq5bSeVumc3EV0iaTb6UStLibHibeZUQSUZOqkaokticuHtx9uSMA/132'),(377,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIV2wGkGic0BZheClRWyib2XZwIYbiaf4CZbwRFmqYrIDajklw72RhmsCV2mpeNqicfibRJHAnb8UaTiapQ/132'),(378,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eooNu1LRicC3Z2nTnamt5e9cqg0slibLKKiagsext3vXSozxnFSm4m1GlBcHTrpwbibafoSjYChBPibA8w/132'),(379,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epIbuQ6HLic8tZgaKsqKsT3C2fx8Z818GxTAF4AqPPqNj5vlE2AfDpoyew70R7BywicBEq9hysibOficA/132'),(380,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erhsOZ2YaD3LicdxBOicjW6Q6RaXneU0xTgajq0rej0ajbtZuYIiaYiabPOfMRYTN4Kzm7WZ11VBZZKZg/132'),(381,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoeVOXyicCx3qV1e1of2H8kLEXC0icQzcqlptDHG5OWPUnHyH2ZoiaZicEErrtCNicqlFAZ1qqamQoUpVQ/132'),(382,'http://thirdwx.qlogo.cn/mmopen/vi_32/zm6ZyC4y9Q5C70s2000kQbzIgsTh7vNdks6kpiaoE9NjWgXiaFI3ocFzhNFiaMibrHDjE1BMCicaiagQUfXXuRBtF0bg/132'),(383,'http://thirdwx.qlogo.cn/mmopen/vi_32/mWuPTibJqc9ia8K6DZRsgFgqpFib8icL3GDia2XCN6SJuFp4MuwhAURdYq5onRticGAXBJYgmKb0fwyp1wE1JAbgu1gw/132'),(384,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epEvPzdSj47oy8qxpqvJH14WL6P2TALhNbbaMpXSrmJibLE4uHjzhrCYZhLhBdPKPYicBSO9HLmQXaQ/132'),(386,'http://thirdwx.qlogo.cn/mmopen/vi_32/5GpAg6Kefnww7o1w4nIkHCpPVwo9zXBjyicicTQVzUpeZJ70RpvStY7vO0WYiacahGMGDHJbHV4mxGpEr6j3pmSAQ/132'),(387,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLhdtAxzov83DORuVQOojnr99CRkUxQXK4ax8P1sbOeictvGDwG8xZnY0rfupBiaox7n1hfcPFKTo9A/132'),(388,'http://thirdwx.qlogo.cn/mmopen/vi_32/KINOHbjUT7qlWCibwsyl5fsxXsRw1aeMpBe81OqmknrcibM1obB4wtHpPZLQ9SmkWWDem7SZMmiczYKrAicaTkGCdQ/132'),(389,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKEeC8LWBTkpgic228ysib66vNwNyfHMRYlWWH9LoiaYWqTT6H5bSe9icxt7Fmj0Y8UlvaWgZZLZibGSaw/132'),(390,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLmf1ww0iaDCvfvoYMAoWCtG0fM0pOcnK3zUiaAPZpNyCe8Q0WbQROUK06icSbPWj9xlTIptUwG2N4PQ/132'),(391,'http://thirdwx.qlogo.cn/mmopen/vi_32/AEekWjiaOUicHhznfHf1xTYhiaNCOGUqV0B4LrLlTzoUzbImVicLcLibaG2Hm6ialia8wRnjCjO7QfNLViaJicGrMZHicgFA/132'),(392,'http://thirdwx.qlogo.cn/mmopen/vi_32/PnBQ4583W058rYSXbZJWgHPBX2C8nHSjs6SeOugbicnBkQQhdMHib3rrXKeGtxJO1MSIpuibBI1ibGicwuia8DQ8Bq3g/132'),(393,'http://thirdwx.qlogo.cn/mmopen/vi_32/ehJPbjToeTLIz4s4IoicMQiawH6MSU370gRvicFo5GgO0ZdERGRPYyFlPSIfUM2U1VHuiayiasibDibX7EEricegyricylw/132'),(394,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ep8q1Kf0I8r64O4FeIaicF16cwjg6jzkK5SfvvVSTpUibXyCQXOL5tcrotzDn6ib5iaBIXh6eVD4Gia61Q/132'),(395,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83ercusHOGicobJK87yGC2533U3555VX1fKq3S32Gb2hVb1GRJbEuU1amsezoqoCZZ18t2T848hZy2sQ/132'),(396,'http://thirdwx.qlogo.cn/mmopen/vi_32/TiaAr0X15dBNdkSSorCVkUXSYdfMhnwxHMRWOJGY11jyCyTNPje7CpmJ7I50gwG8mm7TNdofyI0m92gunIicZ0qg/132'),(397,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJnlNPwhZlicia6DDPFZpYNiaDYEic0hCcvud1BIBDKQTCBiarWSvqQuSTJ7IKLlCbtXsjQuKPAYeb4dXA/132'),(398,'http://thirdwx.qlogo.cn/mmopen/vi_32/cQDaiaSmBnusX5CrN3VdQYEMhn1hFvaK6D2wQy4ibNtDTOvdTwordX7SO9QSeX3zLS9rJu0ajX70xrGbsibpG4ZGw/132'),(399,'http://thirdwx.qlogo.cn/mmopen/vi_32/6nKqGicPDkVuZEQ54zocR0gKoVKYaGYzbibr6yznPEx2AlHCjNsYLfEBicQM3bRz5eQ7QbWBP3OMKq6pAfcjiaRR4w/132'),(400,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibosBt5Yh5ETibn0eRBvvHMAtib5abaUdEiaJy7fUsfB4BiatMap1qERDa18F9Qa5ohtsNtvMDm8692ux9AQ0TliaGpg/132'),(401,'http://thirdwx.qlogo.cn/mmopen/vi_32/poHJtnGAaXg6aJ608gwACavxgLvDDYZCCJxR6mtlibiaDge8bpdYIuCzDaZqg40Pcic7NN9VzUeiafwW6tDrJEViboA/132'),(402,'http://thirdwx.qlogo.cn/mmopen/vi_32/AB6K3UOLysYfy5HKuNiamQzrA8QIpvZ9cTKD6HsULgt3nDArljnblicsw3o8Wibzypq4falxfWZgZxOzLibSiadznLA/132'),(403,'http://thirdwx.qlogo.cn/mmopen/vi_32/3jTeyaBpFdZP4yxbfRu6erpqXzDEwYxgbAQ7eTVa0NaWy4pZdicp36Edicg9P54iceBMPKmllZTKXZrW6icOYB9eicQ/132'),(404,'http://thirdwx.qlogo.cn/mmopen/vi_32/ThV3PHkBB0gibnOG1bVTCjJYbuZH5iaTsdfIXEmRxlKchPbkzVO5nhvkZFwDLJhXogwqovbAvkMTtGJwU0WGLiahw/132'),(405,'http://thirdwx.qlogo.cn/mmopen/vi_32/GkozGoqB0ka1px6to6NCQpqOxykrhHvPMBSgGaQ9uaceFFiarHoCUfME9F9P7HPW3jfsN5cFDxSr4KXnwINEyiaw/132'),(406,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ5MUYIZAaFTibeeMcs3kDwBxakEnxqIkt8lyT3VWKFGm3OIxvohJj5CcFPrLsBibcWS9Uwuu9dPiciaA/132'),(407,'http://thirdwx.qlogo.cn/mmopen/vi_32/cabLXAUXiavWW7hM0KhqJCGUXlLzUPh4lfPexN4KobJaAeDzIPRmvQicQHhjvmrAou6dAvN5ibGibF0EucPUAP85Lw/132'),(408,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epNxyM4Via1sYSmZVkcpxw6Y9ichhV4srJKcHF6TUCl5RCmibhMFoWLe89vj9lmLtrdyWicCvNlzRx2MQ/132'),(409,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoaWpUxAbFml5rIKaiaA7Q0S5TMQvR6pE7qpYSMXk3EibjqHSBrn5ZKVctGTnuPgBlNEd07NDdIGLxQ/132'),(410,'http://thirdwx.qlogo.cn/mmopen/vi_32/ia9EUiaA9CoKicZx1NMdlAtKN2uwiaDh8fZKPvwAVsaQgzz7kAPAZicHMticl6Bu7ahXtkuFaXgOt62q5LNpHu9anMlw/132'),(411,'http://thirdwx.qlogo.cn/mmopen/vi_32/JG2TxdsmpNqibBRnHJqcYBA0RFChSjTac6OEiamKibqHTVmhMMOqDzKJXSyk0zvRiazn3ibK0QC039yXcuXWEY45C1A/132'),(413,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLjGXj4djLsibCJopUIDebtobBXjzibfWQ8jZNC6szU2b2aCw0pJ9htMzXk5nLEwKWfv30yQstuYbGA/132'),(414,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKia6yoCoJFiacib3E83sQBic4wJIXn9JGibGhL3P725xrc5Rib74ZNJ1fS0SzgIWPia9RlqryjPjB75f0cQ/132'),(415,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIV0zSb0u5CUqkRDlkWAhw1nG0al9d1p1COZluRKic95viaesTCS6uoTjcw8Rg4KOwMPzTnwLEIwPug/132'),(416,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ4mmPbWgib795OT25vUpHqkrekTNrwGJp5KkSP6o1WTTBt5GWFGqQ2iaHWxmrVY0A009wFsDbnFxtg/132'),(417,'http://thirdwx.qlogo.cn/mmopen/vi_32/HMQoQAZ1IibXQD9P9bwicc6GtB9XWvefpyWenpVM5deCmzKfpfTkCBqsO4Ajr0s1cibf5rt50DYRN60nCuM2TbWMg/132'),(418,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJnlNPwhZlicia165JawQZib9fm0jPmjciaicJg6SiaF9YflKicVyk5wVvsLg5icxzAS1IxoDTC0dqPZWFI9Q/132'),(419,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKOz2PzU4u9We2PiauwfXHHxfWDHQiaLYHapeDicEpzBicx21ekaicicKPG6EypuX4ibkiatGMXPtic1mO5xxA/132'),(420,'http://thirdwx.qlogo.cn/mmopen/vi_32/IuaVeEnmiamqL4gYFc3hcf6SHIO7TIib3FrWHHZZicoVBLfLgrHuohnFtWKUgIlA4CWGeMqk8mNiavw8uQ7CRKUsPg/132'),(421,'http://thirdwx.qlogo.cn/mmopen/vi_32/YVrJua7icBttVkLb0QdbeA6AeQrMGJzSnW7pzxXn5h3ceTHklxhQNX9SjZNsPErx9u6fntOV5rWEjP1ZyDH1saA/132'),(422,'http://thirdwx.qlogo.cn/mmopen/vi_32/o7DgAtib5jiclfwXT7yrhJ6bibctZrAIwFLc2RuuweWWfmKX4x04qiaAuP1DTRpdaoVbicp9P1q5dtHO5Rib9qsyicG6A/132'),(423,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJo8BTDYxuT2oM7QgIWf4TvVUHAo1H6CPNqD1VsslsEpnu7O7NvLU4xhSpGVezKzVGBNTrAOqNSlQ/132'),(424,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLztsXpzUmY688Or2ZWH8ky7ZqTgYnPibXoCWoR9UibjMr74v6BGBOD4yheUd4onc0cLbXmSZXEE3ng/132'),(425,'http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaELpuTibxSsGxgLpkib14l4Ktic5eGjpP8IibpVvWfahFHiblibmNj05BGKw1mCEFnvWEHQq1QBTVO6VuymA/132'),(426,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eokDdtx2kEzaHUwib3bJiar0gLe3c4sb0wTVQUX6M3hHB5Gd8G5ibZmu0joqfXFIicUrOUXG2w8CDxyyg/132'),(427,'http://thirdwx.qlogo.cn/mmopen/vi_32/mtNLaz0EVdbOcYewo38atbDm4cQxqgP3PW3pKxRSffeycicGyhUgQ0YZr41icDTnXDdNdMiavW36eGfEGhic47Lr0w/132'),(428,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erfpRgswq6t3NoCah83IbmDHIsq2TET8icyzT4nBVAibKf1AEKBW8du3JBfdrdfzrgO6FPXKO1Xbpxg/132'),(429,'http://thirdwx.qlogo.cn/mmopen/vi_32/xy6Ch51UJZFQgT1ibz2cPprLB339LZbkbyJFWuTWiclB2hg57dmjTHdthaAaaZ8ngCFkRRYbxfv7CK0dFeOP2KNw/132'),(430,'http://thirdwx.qlogo.cn/mmopen/vi_32/46AUkKhEV2aibb118DFbj7G0erXeEUPLFGHVBOEgUo7bHibP3ib8EPoGauu20b67UkkxCNPutCXTEBibQ2w08LxBWg/132'),(431,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIAY8nicWdB83AhjNd8OZgqAXBeLqTyVXENQWGKozarSYXksqELics8jnTxicQ9VRic5klOmicCPtQvkHA/132'),(432,'http://thirdwx.qlogo.cn/mmopen/vi_32/xU7Pickem2iaRdgpYsV7h7RqNQID8unPicqonArBicC5GZ7P3QlSybgoMvRDYruqGYEpSHua9kRPm9nVNAnzgBemGA/132'),(433,'http://thirdwx.qlogo.cn/mmopen/vi_32/6xvQfHU7RNhewNRUrLz5YVVHjSqa5UxichMB6c7GcPKX6LQTsFeM2apaZvme5nO0Xom2qt078UyXqzJunRmUhicA/132'),(434,'http://thirdwx.qlogo.cn/mmopen/vi_32/tDLLR9ibymggLQJPHBkQl7LXxQsA9EopjtBTvUbRbhgAvgpUIfH6wpwPHKP2giaz8rPxsjtc8ibCOicJ6oTqDvnS0A/132'),(435,'http://thirdwx.qlogo.cn/mmopen/vi_32/666ZxTwTYHVVSXhev8dwQ0jVuX0Xg52dVIHgYLhTgl6pA5AzSicdKGQ1IQwAoNCazsVvlVUWdSia6hURUU9FRUaw/132'),(436,'http://thirdwx.qlogo.cn/mmopen/vi_32/vqE4ucX4AxNXUOY5I0DVAm3TeYc2EHmEchLlfJLpZMWxRulmkvqlqNR95WsbQ85FjKbvVNnfE9b8FsGLAxbWTw/132'),(437,'http://thirdwx.qlogo.cn/mmopen/vi_32/ySAH6vk0pVFUvllf9Glmvj1usdib2icp3x599CeGvb013rEPUMzfIKibptibibWcGWmibbFwX7cKo058zIqos6pEKH6w/132'),(438,'http://thirdwx.qlogo.cn/mmopen/vi_32/hUaBF6WJ8TP0DpmTb9IguOoSckjG7DjsNcqezZ2shXiaCZ5z07n3q36wqiaalE5Y1lRKQw1Mtibfn2AURmiaIGw5qg/132'),(439,'http://thirdwx.qlogo.cn/mmopen/vi_32/rI2SBjuggic2CXsC95v2etHibibXAzmJ4jwLTRoxMFicR5jFAJZpvwIucowhj4ficurowLK6Ogryw509CR7X2FWK5rw/132'),(441,'http://thirdwx.qlogo.cn/mmopen/vi_32/xU7Pickem2iaSZwHJiaFhrPF59Lb4WheJ9qwzEwlx2VJPMNofL66EjboibaicfzK2yf4k8ns8UH3ZvAK3A0CZ27vSsQ/132'),(442,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJZPdAzRMUVtnChUI2LKvr9AGibQS7Em975ib6B6JuMuJQK1bQwlia2cUiclBGE418ugibGBxBmUDUe3EA/132'),(443,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ5nKRcTbnl3EHOSIUl3dwe4jIicicshxlKhZY1rQQwsdWXSSkFmRibc2CbsDkavefAib5ESsexdYrp1w/132'),(444,'http://thirdwx.qlogo.cn/mmopen/vi_32/3picyLt7jbGEViarOxCuaQiaFuclYpMkR2uwRCOzg6FpqqQhTEOb7ia91tBmZEkpbSVWM45POyibt1wV9fWx1H5wXvg/132'),(445,'http://thirdwx.qlogo.cn/mmopen/vi_32/TnI7DmXKdY78aMOEzjR1IdZxaxIz8FsuVDpBprHfica6B9xTOibDyHlkm20B1e026HDXJEZtweBBtJlEY0LECWyg/132'),(446,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK40cFN9aJ7V0aoIibNJLTtxyGVRuPVibGJa2gLrHFS5kT7kibibg9hn6K5vFoV4oL4SmXftc8HBtkfkQ/132'),(447,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLeicPbxQPEqQB9AeDCMzueeRRa99l9D3X84R0VhubvWOoulK5350Y1QLHWszkCUTBeHqpBIODUPwA/132'),(448,'http://thirdwx.qlogo.cn/mmopen/vi_32/bGKp0JfRCQUcaQ1PQGh6wmp6gcfGBMWtad0TgIwicWcOiaELLwxxgTaZcMIzF6gUyJibf9nFdJ6m1hSziacdnaSgBw/132'),(449,'http://thirdwx.qlogo.cn/mmopen/vi_32/2J6XELMPIaoz1By3EuMenyZ4JPt4urqUjzBG3ko58ewPuBP6h5f3S7y7jJ8YpB8x0zCMl0vfrP5GuYyIrWGojg/132'),(450,'http://thirdwx.qlogo.cn/mmopen/vi_32/B0j9PT51ib1UdIwL1mgLklUJ5ktcWmPKYLpnSJS0hX8ZObVN7y7e22s4UsvfcibPheGE4snODBMJibia5RCQfzWh7Q/132'),(451,'http://thirdwx.qlogo.cn/mmopen/vi_32/YnU5tUCuAbia3q7GZw0kR8p2mxk7iaiaoeHH6qJjRBcOkLu2KkpqhKtLbND24v4gRQaWrSIlUUrZpbXok94jsOHqA/132'),(452,'http://thirdwx.qlogo.cn/mmopen/vi_32/1icjMlTYGP1ib1DPlic38E5h5l05mjnC1Eck5ehqOtIo4ufREySib4r1xHiaG8qAYOUtT8Y1PIBojwiasYxI85dnBiaicw/132'),(453,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epwbggNqcssJrA3Pq1Q4MQazXYhefhQ2KSNfoT92AcGdib29DMu87Vv2qycIccKoxooBMe4BoxSHgg/132'),(454,'http://thirdwx.qlogo.cn/mmopen/vi_32/eian5Yx7AsnVD3lLoLOrwicvExEJVsc7vibkCgNWm6CrwlLzM9cPOeD1HpXfAzhPnXegiawbwt0EcOLhC3KPUhnvSg/132'),(455,'http://thirdwx.qlogo.cn/mmopen/vi_32/lXDxSJbHWkxNEeGFibMRzJCMLhvNAfbUNywU2py3ZDSqoQwYXiaPMU0RRrAJTia7YLwJSIibmicwbGAPnicqKh5TlsLQ/132'),(456,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLEj4mdS1kmFpHkPe7vjkooczCs9u6iaItyb0QdiaVvicVKGHKoF1HNUM0n7iboyLib2b6bLm5CFwJMibOw/132'),(458,'http://thirdwx.qlogo.cn/mmopen/vi_32/IjFRWDxE5DXMicRYSDxOzhPp2PKNd2aicpIp0AiapJpzsic86RGgfciaUa5lw930qgoslK4ia7OhtdGnnB50wia2ib386w/132'),(459,'http://thirdwx.qlogo.cn/mmopen/vi_32/YcHOHzKPrhvO6J47T8teoWfcORcUsSrJo4YkTmeAzaaOhja3micxaYIbllzpRYrfQcgQHt8Aib4icj5WKXDoCEaicg/132'),(460,'http://thirdwx.qlogo.cn/mmopen/vi_32/YuGfLMSuZKBgK8DdCs0AhVu4BOESuVQrZvvUzpD5EEcH9Saboia4wAYtJdictsiaTAcaXglKx2dXTea42JkqlYPvg/132'),(461,'http://thirdwx.qlogo.cn/mmopen/vi_32/pj9S5BibuYrwEWT3UdwVTt6Ny7RTF9zItlicJLYOT7aHviaBrpqPbCfLicgIvKqfbicuqdlbjPypgB8N1NXXXQiapSNA/132'),(463,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eokgyO47b4wSU3QmD03Mygvib527yrhWiauMArzOlicTRBTvZn4qgvfmDLWtlZbib3MKdbqgCicNnE2dDQ/132'),(464,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibMtICQgRHswUesJjCIF06AcHBSL34OcMrPEl0vuQLzSaOKX0GetscoQz3XuWJyCch0czPGapIdGV5XLC7plkqA/132'),(465,'http://thirdwx.qlogo.cn/mmopen/vi_32/lWfF1X3ty2aAOibAPUMHLUxeBpXOhm0EibnRWsBON1YlHcjcSCC4CpvHibGApM1GTL6tamKn3CbxOGoyAh2VtDf4g/132'),(466,'http://thirdwx.qlogo.cn/mmopen/vi_32/75maa0hHibUt8icTL4XlQAFcw0XnCwE293J5lwicuYFpS7R1p4OrPDGuWU9iaLbyC9oqDsydianW3gJM1l5y7WCxteQ/132'),(467,'http://thirdwx.qlogo.cn/mmopen/vi_32/SLyfSBr698OwFacNricibmqxgaXjVBuG3SmPTZe7ibKw1NdwPibLsrDmiaN0C8x5sIEc3CACNVxogQRB7tyV0YeCyMg/132'),(468,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLE02ruTM7PzCFu9FldMiarE7M6KdCZZU3ySawqKMmQLf8g8Jbz2fxEYzN8R402pp9r8MGVvjsbSNA/132'),(469,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLjDSxkI8OsCEUWBxiaIiaXicwf4EARTOIeVsicXVxgEGSqCd1pUqf10ia536VrJw15N2HODWGzmmBefYw/132'),(470,'http://thirdwx.qlogo.cn/mmopen/vi_32/n8ruraDmIKRk618bziae79noyjaR3LYM07Wz1fM17Ricu7eYdnfK9wibU9E02gwzx82MBktK7HgTTLibrXE69TRqWg/132'),(471,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqFibyxKg2tO3Y7FmadMdD1ia5aBPNm8nNBrJfXZp4icythdBxln1SwgyBFO3drMxkhjXKMTyDmOWYZQ/132'),(472,'http://thirdwx.qlogo.cn/mmopen/vi_32/I78W0KbnOUibkJJY50jfjKsIQtvmtAWM3MGwBSzZRXI8zicX6fPXKCwXo1rKHqrD5HAIr2sKUj9D59H9jXp0ibJAQ/132'),(475,'http://thirdwx.qlogo.cn/mmopen/vi_32/wVJLKmibyYdq2tFFQVB1szEuib1AGIPdclJibMDEGYGvQQaYUDfI310nibpkhAlF6PPFSqLYG1jaIGYic4umQdDibyhg/132'),(476,'http://thirdwx.qlogo.cn/mmopen/vi_32/UkRDkeaAHfcex8fLBa6uMFgUxJSoVBanq0bytmuWQ1LQr5mYyHqGib8FgSjpOJpaKic15GXkNmUC0YVICssg7oxw/132'),(477,'http://thirdwx.qlogo.cn/mmopen/vi_32/KHdD5nYzLBBoWydjDyUdBOy47nE5XVug6ghduNxIdoXXYyuLGZdw4E6mC1qx1WN2ic7Jiab8D2LWCCY7CZzfibxTQ/132'),(478,'http://thirdwx.qlogo.cn/mmopen/vi_32/ZPcXGLcgGQ3APiavDZwq0UI5ibr3XnjfpnFDiaNDCeg2zskTy5MucWHMeE2VmkOraMqFT1uoRHoXA2zZkAOUbeLiaA/132'),(479,'http://thirdwx.qlogo.cn/mmopen/vi_32/SfbxicfFfx4glgibhJFDTQuzjicg5kviacu2cK6zMkt4edj8BibFpgY1uTqfe9iadl3xvYpH7MqPxbddU4outdTt4vJA/132'),(480,'http://thirdwx.qlogo.cn/mmopen/vi_32/fUusoEAYhk2bwMhK9pGRP4ibTyEDiazMqDJV00vq9sVwIvIxk87G5nXM5LzofFUwMiaAknC4GS9HwvTJddaCSILJQ/132'),(481,'http://thirdwx.qlogo.cn/mmopen/vi_32/bvhRGNgQ3oKE7qqS3ibBibsQNoMx6t9HXt9xrymeczrWzFmOJYvzichexjZiaQkxy25lqMibwnzsaEHEib2bwkicN7ppQ/132'),(482,'http://thirdwx.qlogo.cn/mmopen/vi_32/mDf7g6pmYhgFYJqRlRsicianqQrcRmntn0P3Y4tfqJ2QRNYSfmzcstDxFwZFSuE8K9zeXoMzY1uWUGxcW7ZpJWdg/132'),(483,'http://thirdwx.qlogo.cn/mmopen/vi_32/9wVIQxDlyq82zJLdHxTkWZDiayticjKSPn5G1xxecfSiaj30jBe1TYLuWEkysuiax4F9vP5O3PibyuG2FGK0TP5LSOQ/132'),(484,'http://thirdwx.qlogo.cn/mmopen/vi_32/KcAOtfW0nF4mfEJr1tosuib3uibQ4qh2qWqzSyzVI3EYOojpUpIZ5xvUdV6Amorb35icUzicFGLCgkv31vcAy0lfBA/132'),(485,'http://thirdwx.qlogo.cn/mmopen/vi_32/KzQzkohK53IUMBhyubIn47ZQCRX7xPFyV4iaIYNXOTp1w5LjvgwfFiav8M7n8FE5W5Q6U8Ih2RrgOcAdFOJyYwpA/132'),(486,'http://thirdwx.qlogo.cn/mmopen/vi_32/efCKKhbGFqjcK7qLvLOoxtdFibhyHw0Pm59n1JLEO3OZ1Qic3y82hbHuT5oSRCyoqibia2DBZrkRa2HSalyIsVCib8A/132'),(487,'http://thirdwx.qlogo.cn/mmopen/vi_32/EtteKwKuceQkmW0nOphv1KN7U5YE11N5EHBkQibLvBlSM4gq9CmG3BZGkZc1uDz6LqAAU7UbicycU2ytWFGfrwicw/132'),(488,'http://thirdwx.qlogo.cn/mmopen/vi_32/t6LPsCbJ0TY7BOqbwBP7pZaspsnak7qa0s3jShbmm7z7EpQkMqia9xxSUcvPeUY7cr2ORL35N6nqkTncSr1G6xA/132'),(489,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKS9Amw3mLkRZuYndCyzFQ4NF42OykUdxb6VLmJaWJurbfrvicPaqrcKPghmvGKRG9iaXmLwbTqstmg/132'),(490,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqSyRWoZhhsGdQXf4eDWZj63KtBdgib76y66icklOmP6rcHvgZ4fL9KXHvzFG12XhlkIvTIJqbdhPzA/132'),(491,'http://thirdwx.qlogo.cn/mmopen/vi_32/4N9jXHTzwQ6j2ibhzvx0c0XIYn7kTFTzR9YqicPSwk0VWTHumORGQzaLJ55qJBtTbibcpTWibM8y2XBxBKWoW86mQA/132'),(492,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eo5JM3JiasDk12dTsoNibvv5G0eQ0icOVLMuuxxEuOhmunAsibBqsZ4Tjt6gJGNOwpicAZuNlfXOc187AA/132'),(493,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLsOFgxUeol7xIvhOz3oDyrHibEGNKB4KD7lBI1VOa4yAhKnhvH4Tibnic0MwlGZib0UadblQCMzjklMw/132'),(494,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJKbEyJm9acyTVkquIZpp32CHkZGItvk9YkKnqGiaZIAp79sg5jKh0jqLxibfibDuia9BsOfpsC48hHIA/132'),(495,'http://thirdwx.qlogo.cn/mmopen/vi_32/fanyL5xEPIHnKqiaMKYXzz0QQLlUELEtUokmI8giavSWfemY9QKl4ib7dicYeLXQRdAkfO7icx4jJrGVKrV72IlVctw/132'),(496,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibganewlE0daiaOBzMMHFTGH1YzLDSq5PiaVFpMgqUleicHNjBgzDUVOpsxy6cibNxSs7UanPxL9K793CRkGbxBA51A/132'),(497,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJru7MZse3ErQmEBMVQRq283wpx52JoPxZiaqD30YWJ02B0s7KBqibacGHxAiaHfgjkictoAibVZTpyKYA/132'),(498,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJZVnjSGfMqhBjznpBia48HNQtSLGzXvhFhfWkBUze2EAfb98MkpfHry0kHY4j4GjVkBFuxAu8vrbQ/132'),(499,'http://thirdwx.qlogo.cn/mmopen/vi_32/0kQccCmosff8xddGYZ1DGWqwC3oBS6mBa7oDn9HkibHbiaOc13nAWiaWV5u0XFNUjwLH8XH1gpRVGmc6WmVRA6TrA/132'),(500,'http://thirdwx.qlogo.cn/mmopen/vi_32/ia0xleERpheU2THnRbqRLtFeV9YbCDdokufO7icKh82iblqIoV8TBY8AvzmibxHIMONJge3C4neempm2QQ2SMuQocQ/132'),(501,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ1JlQIMbdMQs5w368LGBr6wibf44EX4WNNXiaycSYh8uGMXxLDtlw6IbeNKQicd6wtblYPScuGibibMMQ/132'),(502,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoeVOXyicCx3qWU1bWNeDx62LT4iaS1oDW9BuYNK1AazmOHnw7vDZvbkjZDJKIEvnibdv3gelGhmWg4A/132'),(503,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIoDuPApzrXbFNnFCqmWqUC66R26738sBkthTGvwdMgFpyOxM1jGyh87O0UuNibEbfoP9QnHWpnbAA/132'),(504,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epdTfp7dHO7pE8ORHjHMwIFa21hVQnEedI28cn8aQmcD07V4PyKlcx37RPoBn3BSqODGicEVWbe4kA/132'),(505,'http://thirdwx.qlogo.cn/mmopen/vi_32/d5oPlvPIZgZJqjHvhATNoyqiaFwKWcMTWkkVDib9fvksKgIia4TQPOBJ70zMxERNmxz5kvh5R69sZFow9JOmeK6aA/132'),(506,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eomSnwDxoXKdnicVPcN9A1jcwXlLlXu8SU6PZ3JXGMMgwf17tCo1menoOvuT6DHqQjMZ7bk0tp0UBg/132'),(507,'http://thirdwx.qlogo.cn/mmopen/vi_32/ORajVMkHYO2FV765VjQ46bIPVYWcrY2M3ym6TMtvDvYcFlq97o6SXh2mgibIeOYvn9dEV8ZiaHoAiaTc0SPDFGMAw/132'),(508,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK9wvXdaLiaeyfDaibicsdeoSYO5ZZDw8dJzKKckVic9DUfcIKSI5BOxTp3EDpMdliamLDt1G184FmF4oA/132'),(509,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJ5nKRcTbnl3B4FP7X9DicvicdzD4ZS9sG9Ctz44VukKHCsPmnUotFIR3XichFwpiaZhNRWyUUibRB9EhA/132'),(510,'http://thirdwx.qlogo.cn/mmopen/vi_32/SlgQnpM94Trrpx6TSOgtwVG7oNYbEWPxxiaZjulH0ibNMvu4Oa57vKuCss4mcHCId4wVAjBVwkpTicaz8NqNor22Q/132'),(511,'http://thirdwx.qlogo.cn/mmopen/vi_32/icGsEvnThVFKCa9JdqvEBXchK2xIb1vOpyicM41EKss0ibX35bibyfSmu1w3n9hibFMbuxhl5AO08AsRZam7W1voQNw/132'),(513,'http://thirdwx.qlogo.cn/mmopen/vi_32/fvVIWsIkUDZtualQWGlpZRvYA8icaKia3oQkhVhYGw7XIq0iagjic26kTtsFscdLQ59hso49gE3k4J4keFel2tg97g/132'),(514,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLdvKClLiaIFlWUYJbEjwPI3tpBRtsCz7anU6HKqKxDu706LQFP7OnMWSyRicksIpbtUUqQZCEaNicYw/132'),(516,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJE33c8iarHlX6Xuua9OMsV5q7HIGM1kZFe5P24iaV1MSBq328icl0PbyA4FwNe303x1XkvU4oibeMDAg/132'),(517,'http://thirdwx.qlogo.cn/mmopen/vi_32/57ibsjNMHMmrj4PapVMqcc5zkiabgUkGVIKYJhuic254aOuzXOhGAe5WnARmjzcRktJNECCaQudJEoC89gaLQWSaw/132'),(518,'http://thirdwx.qlogo.cn/mmopen/vi_32/UBicVUlOZib7Ta1SEFvXDpRCicEibxyl5nuaupuzP6JsaLybXSbUxzeOX8JMhF3JCjAvWXCYRUQicnZjaP9f32skY4w/132'),(519,'http://thirdwx.qlogo.cn/mmopen/vi_32/F6WIe8BQ2bBOZA5yorxny6icgHW64k2A3RzZVRGvth9icVIAeicUcrtBSEnhMYmk7qwiayUA5VmF3rXKvU1CibmPTUg/132'),(520,'http://thirdwx.qlogo.cn/mmopen/vi_32/Aicl17fL1sjRALtTAibPL2bZw7niaibvx6jyeWYJrNCakLuibxfqjGedUDGCKsmicAymbWeXLSia7J2gK9AXf1dOz2asw/132'),(521,'http://thirdwx.qlogo.cn/mmopen/vi_32/3v5HfC0Czic0zdRMdbvJWibFR0F2pmRFhcgWuiaicUqP2Pq8VDOnZAvgiaUJJhWwEwnVAENokiamfyLFbCuAYt4N9ybQ/132'),(522,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIQObCXrNjNVOMY5ib8dWpRDXIEsD1730icQPFPctsUYicxmmeHGdfs6ssxkhm4zaliaO4r8XbWHd7kUw/132'),(523,'http://thirdwx.qlogo.cn/mmopen/vi_32/a8ibq87l4H8Jia4l2lgibEsmcicOCY8ZgJOEcxibPdYs3ysRd8uXica3Dhs0Cg3Q3ic4HXwia0oy7rvGo0PSPaz7go6QPw/132'),(524,'http://thirdwx.qlogo.cn/mmopen/vi_32/nZqoatYrvfaWsZqqtTWlgiciccKhmKB2ajbEhmQuon0jOyTWxIWoOGia88ewfM7TYYb7T3ljDS3jFjdSRXiaQwVngg/132'),(525,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIzrssqNyZARcc6bSkdRvAhxrCSagKyPlusQfAnWevEPibibPibWjHpsvyKTtMQ7xFhiaVkehjUCWgWnA/132'),(526,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLtzg1GeFxtFFYicSkWZNeBFwuCCGbUIEfv4WMlBp0b9N9ECUSymfBwru6ujrzicIWlWOFb7oNw4GMg/132'),(527,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKpuwV5RWVialuuQS5wsWmicmsmaVXpBnVESuUdszAe7z4JdWA2dc4h3sk8sMEjy0iaSSicNBtm2G8PSQ/132'),(528,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKGHDsL1nEmpe6fXENRxP1Os7AjKiaCqibGdBDUaqcyMNOXtAqic1FwaHROzAtCBJ6T2cKJx6bIUjhZQ/132'),(529,'http://thirdwx.qlogo.cn/mmopen/vi_32/QD9GIbqK8mx2zu6GbHQQEfZMOiaOzjWD5e7fyAjqyFvvG8a0g5JKpJYshicKLhbqLvvaqlsx8lrviavzVosYfn4ug/132'),(530,'http://thirdwx.qlogo.cn/mmopen/vi_32/rVBU51zia5L1T3zmWsoiaZ70Vblj6X63kqzfd7UKB1XzHfib25mtpK9icer2IkicjrPCLhZk4Jcia1XMglA7gV1gI8Ng/132'),(531,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKfFk0I54QjIKoNgQSH6vGlp8VbRSyDFUEzv7MKHnDicrWOJtnibEpmb0Dtkd1pRUB3wLmiasTJWDAXA/132'),(532,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJTYft5jlcXcme2kCXia7hibfPMnhicicGJY7L9PF21SOicLNQNKZtRsXzelJDWbGUj2dWebJSPA6R7LTg/132'),(533,'http://thirdwx.qlogo.cn/mmopen/vi_32/QRQLDQ4ajMia2zEYMsaoHluaj9UQTljrlA9RpCD89bs9x0s8EJiccpyAP7xXEO11eicjxbAIZ5FyhN0fR0DGxiaDicg/132'),(534,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKRyh0IE6iaMMqBpnG9rolhdWHAVdHaYvnXZWtIC5JGmcA8N9fiaR0O131iabHZkHmpO2HibOGlnjiaqXA/132'),(535,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erQjAdTZzZcUQaO1L1VX4BO2QqePcicXdibwhHEDLoBDfroMAssibMIBXEdURPSXqXFic5CzQ7vrTUumg/132'),(536,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKFnqNtXPgibib2pKLTibclj31icvDVgV0JcfVA2LscicW9804MPicWIwpFAwfVTdfawly3fiaibibjQKuvlEA/132'),(537,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJYsPjUr1yItJ2BexRGlMRDEbonuTXQHLeP6EXoAL8Tr1Vviah5tQQlYndZccZX5wR8jKBYm8z5WlQ/132'),(538,'http://thirdwx.qlogo.cn/mmopen/vi_32/ajNVdqHZLLBYlPyjKWco4kmicC8sGbWgO8WjxXbZoPgR89zLHfia8z0H943gkbJ6ibPMibgW1XciawFn8XV954iaIDEg/132'),(539,'http://thirdwx.qlogo.cn/mmopen/vi_32/dq5eL9bYznzncj85gPgzh5nBibcn2p3s4FoIkYZk1KpnSDcicmW9JXV6dMtR0PMYt8ARcDKqvKicriavppxFJ41PXQ/132'),(540,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoVXR0l4EW5aiavIiaLPNSxRzElhSEEiaRSLeyLEiavL3Llgn98R08mmPh2dzceu2aFcVvnLSibGZOPHNw/132'),(541,'http://thirdwx.qlogo.cn/mmopen/vi_32/cPvfvS7X9KFNDMqfl6xuiahyVjvZ48BVMCVAzLK7FWcL8JOlYx5wzSW16ORyGtatOzGWjCOLDOWliaM1WibugPlUg/132'),(542,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJtPgFFxupgtBtrruukDzsqE7coATE26t2ibMNT6ntHTr5c5gTZIqqDvQ6Bu0DcDVNeq2YvPyNsMgA/132'),(543,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLlxkLumIbol5mehwFPpAX397FD5bWwthBLv8JZZwteErkoIBRK3sEErH4Pqg3fTsM28iaYcNAILFg/132'),(544,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epc9o1HmU46taBZ6fltqEpmNE0LenWC4FMFIcj5FuEGaMCWVTe5pskCPgLIBLVTSV6Mh3ujIboIyw/132'),(545,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLmiaXFTaXFhoYTbk7InM9HoRicP8Z0ibQPQnI0axKtGpROfxia6uQXGibYWia3RHW09a5JudHCNdTftHAA/132'),(546,'http://thirdwx.qlogo.cn/mmopen/vi_32/qrKXAdaFFZGBqrAwibSbGnpJumpF2coe6z0rQmRoftAaKELSHCz7Mo9mZJMk1lqiaUFX0GaiaicqlzVoStaQ5Pa9DA/132'),(548,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLUn043LbLa8cDx5Z6qxa3HXhgKsPmic1juJUjVibWyd412GEfurn5Cb0GxeXfSyr1b3epHsm5yfRuw/132'),(550,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK4hFItvjWHia6cbKl5iaBia1icf1t4uNustve4Xrkl3iaQnHvM1ZGPt2icK0WwfouAAhIBvsk2O44RQhAQ/132'),(551,'http://thirdwx.qlogo.cn/mmopen/vi_32/Z9RUmx1HgtTebhPBaNguJbx5ej3HpR56ZWznhsLHV6P1Zniay7812XKWFW23iczxQ2J4Gf43AW7bCVFuHwhXTgQQ/132'),(552,'http://thirdwx.qlogo.cn/mmopen/vi_32/53C8nSicia0mMFx8qicsJtUib4nqUFfkbuzWjqQ2ex3opoJic5vX85dRZy7yfzWJgGmtwjb8GyxiaibnicYuHhrWJE8wrw/132'),(554,'http://thirdwx.qlogo.cn/mmopen/vi_32/rBkcY7hCxeLUiaURGvib6WhyxWPZUVaoiaV5NCHr56pulYUUvtWPrpGR31OPCvOAAaw5icibcVk7PhT6KOwmlCgcRvA/132'),(555,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqhrDiacFGFjWHiaKvGBdWns9eMoria4BXmv6Pe3hDYaNyNuhWibKksaMS1ba1o1Th7Kh8nQSiaqcviaWpQ/132'),(556,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqdkJhJVJN0gmdgCaa5xuKNRJMrugcVvq0CoyZyboL8NauEU004yvJKKibekqZeZhtmb06MpAfCCPg/132'),(557,'http://thirdwx.qlogo.cn/mmopen/vi_32/IwSShqvZ3wHms5mVlr4O3mlwJkhsykQa0iaOz595NHc0VeEV0RyoaJBU71Aalichibokg2USlGKqZyHQwwWhcSHfQ/132'),(558,'http://thirdwx.qlogo.cn/mmopen/vi_32/gFYcdDGS8mgRd3bwfabeGwkydP06Yk2lzo0m5AzrQZpmlTyrYPdxUzfpyoyGdCG6jGWG0peWztM3IDMibPYXChQ/132'),(559,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJqNViaBOvibDFVrbyT980SmXY3Y0eibeGNKjDABy4vpsuyCIVStFLHbq8r5qHgZibf9RDYfS2Bmibu8lw/132'),(560,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJpBl6xwMmnsPdj7jMAqEjiaEySpsxtTqos6nICpVc6XCC7HhjicSZUb4OYLXMqbhSa5o6QWzE4czIg/132'),(561,'http://thirdwx.qlogo.cn/mmopen/vi_32/YpDMOHzseU15H5asujx7YNkI751lN3BseIg2iaj0cpUHFHMX8XkjvL6YZt4xFD3HibicmYrmXbKSwjozf7Vr1PZ7A/132'),(562,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIEeWyibgWuBI9vrE0CNz0C1QgZmzKMupPp3kqFUMOGdrXtlBEGMbvPhZtKDiaBw5eiaGXyO7uk8EyXg/132'),(563,'http://thirdwx.qlogo.cn/mmopen/vi_32/0DaDVia5NZicvjibuNuzBY14OJiakxLYfyT9f1SOaeUOr19E7icczGxjiczETq6JqCjwibiak9AvWxiaticckpvZDceLVBDQ/132'),(564,'http://thirdwx.qlogo.cn/mmopen/vi_32/LlalvK2CKgpfztyyyAdjMc7zfXdZ9cpkbuhpWKL7LBib9j92oPCKRRGf5YJ1zNQRmK8Dic16xNbqjnjzCF5yA9hg/132'),(566,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKxODAAicBe8OicqYQ3PrtrLJ2aKxyutYvdq0JAcjyZzFwaSVQNaibMNvibQVsYYicgKNXBxb2ZukVG6WA/132'),(567,'http://thirdwx.qlogo.cn/mmopen/vi_32/TVCk2MAggd5l4MrjKXGsKkFE2QvVNwxcibH5q7KYuQn8s7fPP7R7HcVrqW7vfYic24SYnibGYIN3O4NNSKOtMJPbQ/132'),(568,'http://thirdwx.qlogo.cn/mmopen/vi_32/PqFWAMTRsSIBBjFbYFEBca9icB63WdoAwBdY1S7rQwdAN5kbWHiavsozvQC36fWgjRUjKnKpYSbNsRwuicdTLDC2A/132'),(569,'http://thirdwx.qlogo.cn/mmopen/vi_32/nIhZDFQdTfbib6EnQgeM2iceRJoXFTXUFyodfDmWywbWD8sG7mk19QUFhQFicCFibL0cmB5EpHUxFpnXFsR8gbnyQg/132'),(570,'http://thirdwx.qlogo.cn/mmopen/vi_32/hwcarhqiaT49v7mNo6FLDJKrZ29AvmHekDibic7D3VySFrvhVjW7WiapaYvp3LAw0uuooV7gt2xsZr1odoXRCgu9ibQ/132'),(571,'http://thirdwx.qlogo.cn/mmopen/vi_32/RWib8FCKWu7d78sU05kCKAUccFCItGakHK5Wsu4ZcXfgaQXnZBGeKUzQu86OLibZX7re0xQcKMObuqJwp8ojD6gw/132'),(572,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJq24S8KnicnAVa9BDyNMEGMsWf1nFClub8SR2Fc1pDbZQWncBIgtUusmk0AQpOMhwLk0EcvnQ8WGQ/132'),(573,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK0tc2XvCxkhwdkib5icLNSHtZT0lF6QicBCWOIM9Jp2zbB4RyIXXnwWp5K7qsSrbUJfbMotpn6KPIVA/132'),(574,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erlzscoVPqwC0zTk45J6u1PqjxibQvlcowCgheysxibCuZ5CJJVVmKQu7p1xU0yppy45NT0qg0xLe0A/132'),(575,'http://thirdwx.qlogo.cn/mmopen/vi_32/icoiaJxxG1Bb9VuNckqhqlJsVnwvpBIBdxz0AmB2KMfBnqdlEAXTtMV546SmYZib7OlURsw8a4YPV0DwnrWHbbjMA/132'),(576,'http://thirdwx.qlogo.cn/mmopen/vi_32/0BfcHk1dDdFdERVvAHqoyhzlp52CZvVLKeFn6iaz1HsATjRFN8DVupM4KIibLXtArdtkL13sZY5TNllc2WyNiaB9g/132'),(577,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLjczEniatkHahh0499uIIMq1Y2miaCqX86iaoJZ4u9A9zC7XgibHEicp3KialwbOIHfH2RWVnvCaawsLYA/132'),(578,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqF0gWIkrgdu2cqoOicCoVfaRibkLOv0d457U2WXOpFR9iczjG0rYbn48eUw2MasuWr2ZfPtasuhv3ug/132'),(579,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLMJ6Psxmv3iacGODyV47GXNdxg3ibq0HaY9Bv6yf0fBZ4gFsUIw2nLwfbxN62kyUIb4xxHUa51Kshg/132'),(580,'http://thirdwx.qlogo.cn/mmopen/vi_32/iaNIJJ0bKzoJgLiaG678Far4BvK85Liar3wsAbANVwVKknFjlAaicms9gqKic8zKib59PX7Eb4DuZaH6iceEZH8EaqbgA/132'),(581,'http://thirdwx.qlogo.cn/mmopen/vi_32/7fUaIecV2jI2LJJ3kWfFuiaMxyxhMEARaiaticoa83rDyaWtfyasD9qhUhfee7B4GWNfjdOupy5CPdYc6SfkRtghg/132'),(582,'http://thirdwx.qlogo.cn/mmopen/vi_32/Ria1H2nHh5b5sobSiak05fciblhV5jmsAqEhEyWgOChicLAiaeicPVcglzSzgmdX4koCuokPopd4ZWgZibunPScONJNUQ/132'),(583,'http://thirdwx.qlogo.cn/mmopen/vi_32/5Xwt2RdSlXrFkYNU0jCMK43XzMiaHuHUwfvB6uQvHJmWm7kHUCoC9XAxr64HCPLSmbwLWmhJ2mtKolBicO1MicJCQ/132'),(585,'http://thirdwx.qlogo.cn/mmopen/vi_32/PIcpTugROXicDibpBGcvCu1cEplwuPajic2mRL2o1NgyArdcxywoaE6qW5H28Bz34GtfZBU4BsaQ8BojERjUExw9Q/132'),(586,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eoaR7bz8ScOpKZd0aS1z8ouAicZ7Io4JOAl282XMjkDo4cOwN7vjUzVzAIZCicic3eVIIGaGw0F4ohNA/132'),(587,'http://thirdwx.qlogo.cn/mmopen/vi_32/FKxk5kUKibUuWWRDHQeth0RiaricxMtIQFibdPHibYxxaKq7mEAWVlge7wKw0qhd6lVBEgETNLnsr9tKnMGQRNSmUHw/132'),(588,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTINqHCVmh2sRNwiahjFzfmsiakRPX43VrEVWrqLtkdDqDWjKloiaia9CkYN0Z8ymnLv1js8TBJows0tuw/132'),(589,'http://thirdwx.qlogo.cn/mmopen/vi_32/qILibBicF0jzeyu24eMiaChg3lcpDXJCZW7YJbqLKGibX9x37RGTA0lM98NgF768WwhyqoTo9ShAQwNJDHlNN5sicrA/132'),(590,'http://thirdwx.qlogo.cn/mmopen/vi_32/jpRia98bOeHiascfLpSVabcN84kGKicBEkjRFxwXSce0Y6Gxxhs5ZZCG0JvibakCNm2t8t8ajibfFLvecCkiaazHlsFg/132'),(591,'http://thirdwx.qlogo.cn/mmopen/vi_32/21IZUPVwbAeRjxuwEJxgej9wBwPKmRSplNPPdyEiatpJXxcX1O7icduhqVWxNzCiatZEibMLk3fAibW8h1z1rkO3Y6Q/132'),(592,'http://thirdwx.qlogo.cn/mmopen/vi_32/ZLIvnq3uj8eoFlOclwnW9Yd3j2Kibp3xUBKc7tJr09skbicFHyRZzicyY3Z6kyeHl69g3rzEOtDQVYvgriaD8icibu4g/132'),(593,'http://thirdwx.qlogo.cn/mmopen/vi_32/aRzndKDDP7LMicx73tSFk4Sqptcic6HkibIZibb1BTuypVptibA1am02XUU3gAA6ErOXqaN9JvjpZWgHKBfDzKuK0zw/132'),(594,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLjHnhx5CJORwxDyBaI6Nr34fhKDhZUGIKN6vFcajoSiaBU6qTWKIusUuoLrTw8ibFg0ibxNC5wKeHMw/132'),(595,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqT3pba9RQEXNyQoAnou9n4SB52LhEkaxQc81zemw6s607SjCIWzWXVrhjWYMOUZd0vytDjiaCp3xQ/132'),(596,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIp151ytZAU0Jro8peiaAJ8GqwUOU6TjaufGicia5xYRibAicokqg0vYC9rtkh12roewODEUq6UN1uzXeg/132'),(597,'http://thirdwx.qlogo.cn/mmopen/vi_32/1aBHA072vhiatzhbmv3X8lTRETmBLXAo61lCAfF46Cz6CiasKKibkgLaHTtttKMdOmxMiaQ0D7cV5AyEr5At483WAQ/132'),(598,'http://thirdwx.qlogo.cn/mmopen/vi_32/VU1o5V3aoSj2Rlibh0hsd5qvd6jobQD5jiaoh3xxCL3PECO8yHemDmMVQkL8nHO6okGNuGxU4WuDIOoV14sEoiajg/132'),(599,'http://thirdwx.qlogo.cn/mmopen/vi_32/PiajxSqBRaEJV5zRUoS5b32pXLrcRnYzxgsvibON7PmoQGO0tDpRfQxtjSA1QbK2DodWMobIohygQibicgW6eT0NCg/132'),(600,'http://thirdwx.qlogo.cn/mmopen/vi_32/3Vjaiapy7Y9ia1yjLxMzItziafB7LDCy7RLib32x22VLia7FD10QqjOlUy6RgZMG7Yd2Io2PIsKYUVjIf5Efo7YsGeg/132'),(601,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLmcYgOuzibf5eZyYeicu9v17OjUibPE7j2jvvWJW9duRzGibItrtgqa18BBlsMdlfibv72x5YSTtuu3mQ/132'),(602,'http://thirdwx.qlogo.cn/mmopen/vi_32/0mkvfeK2k1UQicK8wME7CIEITcfNLApwic5YDW1fwplXj0r2cHNibm1g092D5JeQoY8XIXjEudg99jDkMqlSWcSDg/132'),(603,'http://thirdwx.qlogo.cn/mmopen/vi_32/nibb7W6bx5xkDWH5lMcia7hSLryMloic89Ddbh5pKd7GyWibiaM5zvLibmugKickyQzAPqq6k7AJ2F9Yydclk0C56zIEw/132'),(604,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJpBPzjhMck6sibyGFG5fibNImJhXVdYqJ6FgeonUL3LtoDm2QgBKR5dd4b8pens3Ss4m2iaHv7iaP40w/132'),(605,'http://thirdwx.qlogo.cn/mmopen/vi_32/C4FJhRmmaXCDbNKE2IGN4Jxzb6zD6uQoxT0ic9q7Q1XR0vGRsq4NtsRR0S2N3vTeqlbcfch1u9caojwfebe47QA/132'),(606,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKgCI7dRTTO9UYEmY7ynWHpWo4icaKuPGAJKoErBhMqoVs3aiaknQl5ic16A5JfJHwrjxBtAo8j9PGng/132'),(607,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJz9nMiaiccbLiaalNxqAXiagkJCT3cuAApXy4Rzn9hu6z3d6XDxbAYqsrS9vBf2yfEdOL0dicicrwOZSXg/132'),(608,'http://thirdwx.qlogo.cn/mmopen/vi_32/ibYeqX5m4ia00Rq9Keb2yqUNeQfVohLaibjavIfKCERCkCjMP927aFCYib5UIG6kjSTibINsLohEdjdttemg95814qQ/132'),(609,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83eqcgMrVvicyRrsIibbKRaicgibJocoYp3xicELe5n5KIsEoM5dHibCdnAgZCz6E2ibqL2djEic5Y6Nl5V4rwg/132'),(611,'http://thirdwx.qlogo.cn/mmopen/vi_32/PbiaYPiaVDfATe8OWj60Noa6NRIFzZGnT0l6Rae5lfPRaaYQrT690JRVd3xhdZx6PSS0teAS8Y89wwYibstZv2q8Q/132'),(612,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLic5rPK9EkwqibJePZthGd3YvsxOo456jRXQIVQHPP5VVRyhpQFbhricQg5mDia9lJvK81wugricKyhIw/132'),(613,'http://thirdwx.qlogo.cn/mmopen/vi_32/XD5jW9kdEoDBnrstvGiaQeVfBALPaPWEufytcKKjjdG3ObYtuEJM17onvsPFRYq3Dw5iaNJk8usN3icriaNicU2y2Bg/132'),(615,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIpr8mtHBHZnB7kD3SzfZh0aaDS3THlxHo5tqic6L7LiaywJtfYooHVF7nLZvmHBLMGyCqYJk6OTH2g/132'),(616,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKoRc1VyzdCZJYVgJ5tcOrNIgO1CknOKY1JhoI8yI3L86LA4HNyrzssuozcFJ13yDlUWBmgLnLqMg/132'),(617,'http://thirdwx.qlogo.cn/mmopen/vi_32/EpQdpktmrSTPnjFc22qg7ibt4Ggpj1B9ue9WO0ncn2V3fLpuh8dPnBwyxP7CAtgyQghianRowByHrvoNeUDwQ8qw/132'),(618,'http://thirdwx.qlogo.cn/mmopen/vi_32/ILiaU8Lv8lAQBk1LX1jxEKgtECwNSRF8CCaMicBicHUj4icwSsJKgWVXbSXqYLVicHNVEFMBLxbW0y6aYOLsibWFdoGg/132'),(619,'http://thirdwx.qlogo.cn/mmopen/vi_32/uYZtqUVftMZQIKe6xREA9UsIKlYAlT2RUicYcsHNCYMIGJNEib8jh7IZxeSBgSj429icRaMFI1NNY3olfNIcibX95w/132'),(620,'http://thirdwx.qlogo.cn/mmopen/vi_32/rVicxdX4UQzRvNQEs83Odlicr5E9opXnfKTiaibcoCQVgh2TG6yLTSw76vCbgqxpI0dAcO0FibI83FHJTw5a1uTFIMA/132'),(621,'http://thirdwx.qlogo.cn/mmopen/vi_32/UEVqZKDCKVXJiazYbOM1A8Vt3yWSWcPOgC3JCWuMKxZIx8PvBibicvulC4pK8w0hea7GKnicG9QFW3HMJRLwdDZGDw/132'),(622,'http://thirdwx.qlogo.cn/mmopen/vi_32/s5B6KibN6ZzONZ4ALIKwS0abzhAlWVqqI4rBnHAqic57rJf37WHM1JvhGZscopTKCYW8j6I5rdDWicoCBJcbVYFPA/132'),(623,'http://thirdwx.qlogo.cn/mmopen/vi_32/cnNBSrNVDiajDT35qicL6Zu5jKcKXZZ2pb7BEMsnCiaZbgc7ibYjR4IguPj0EaShe56mw11FYsL2pXXl3piaYRlc8zg/132'),(624,'http://thirdwx.qlogo.cn/mmopen/vi_32/icmicKNJCrx0poaA2Ig8awpaBV0zBMFunddktABGBIGQpZz4m5OZWeHr44xKFibZSyh0cWcfQMMZDTibzBibOBb1icMg/132'),(625,'http://thirdwx.qlogo.cn/mmopen/vi_32/yS1FEFJpCOjy8JxhHxziaLlBgtKnB71tZkZNfjKhXG0asBxoHy3ssLiayhVSaXI8I8ZXyxJhsxNRwHItuQB1WLqQ/132'),(626,'http://thirdwx.qlogo.cn/mmopen/vi_32/nqYNgicG4xlNZ5SoppuzxIYOfDkFQXEdKW1jhrJVRqse2zxBticeDSCsichd3y5RrVS7qUKDHZYtsLgQ9ibBmPpaaA/132'),(627,'http://thirdwx.qlogo.cn/mmopen/vi_32/fNibkKjQ6MicVIwFazicDpMvwAPFNZJgRd37yicCtdNLZ41y8pqHjVIm4YaTicmWfpu5fIuJyGXU4JMeib1QDJbhRrqw/132'),(628,'http://thirdwx.qlogo.cn/mmopen/vi_32/YccnNtxC5ibZ9WMRnqnGIzPyXJFpk6GlHUWOyxAfQ0vic7NRibe253DBH4z4g03EljtbamkRTwWSvZRmqPT4E6DXA/132'),(629,'http://thirdwx.qlogo.cn/mmopen/vi_32/6CWbCHTCKeicicuhn5icfu6zE70YgDuwTY2uzGNyXoddGaPexxmaxZY2HF9f9jRRNqetjPwgHn60Dw5hLxzdqcVvw/132'),(630,'http://thirdwx.qlogo.cn/mmopen/vi_32/aJXh674FLUXaHtXZSkujj291wGs5ccdjjHQPqFhGhYYlfBX2snIicQwRgicF30UbJE9NtPTeEeDsbHVWmZgkmb3A/132'),(631,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTILhR0NtpB4fBfmSOOJ2iaXhHMAze2uNCicMVSuDHbTnVv2xCCgzrOyChU7nyo89de8OxuV5tjeNgDQ/132'),(632,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKzJVefyI5EZsAQpkkcTvWWYzPNM0VBUvIyOkE9h21zNL5n0Hia97LeQFkIkx2gjgxgDZ3SlgYTZzw/132'),(633,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLWUfXPxw93lviaEKmPUSTYgRdP7u3nia6BibPkmHQl2qBzibTLcspVPMDOV5tibKU01Aia3kibPaicEgTaDA/132'),(634,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJmPvRLwIAId3S4NRYHLCrfnPdIfLiaL2EJwYP5cexGKUicKHy2QG3LdPANO8SmTumjAiaXFM68NO3vQ/132'),(635,'http://thirdwx.qlogo.cn/mmopen/vi_32/yCzPjZbMyPLIpQHFTZ2ArrwicuK76rzkMRAPLhWTsxMtAHdOUiaicytEEEEs9hwWDpjfrWmwK7A5ZzibZKbuMBTq9Q/132'),(636,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLibGBic0GXS2icSybibQ1eTqsZsfEaA7z7LECfuicyI47bguTvhCuZRk0tswicV15k3nl92LWc1O87icvUw/132'),(637,'http://thirdwx.qlogo.cn/mmopen/vi_32/Y1hovqKDZbRWVTqv3yJW7vnWHxzZFT1BQUwdJMYvVLWcneMtwoJpcicjWMcFaf5cC03w2PjLTQCscJtib6WjUy2w/132'),(638,'http://thirdwx.qlogo.cn/mmopen/vi_32/RbGf9FKbOibiaibMdicGBnGoln6meou3CcSgFiaBnlLHXla12dcgRWeHPcYVhYUMUFMCib7cSLibeBvoDM2mzoibovqibHg/132'),(639,'http://thirdwx.qlogo.cn/mmopen/vi_32/ZgyJXQ2SNTycpFIZPAUyGMHaI8icWROJj3ic6hxa9AjmTX7oYCVouBo9icRibBricobJ6kG4OSBVA6YysNteOy4m6fQ/132'),(640,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLS1yT279AmyibwlR8tqUrAwjdCwtVeJjEwKibqDuDFNOBou6dttmyiayyDK83hXUkfTDfXc2MQ7Oyyg/132'),(641,'http://thirdwx.qlogo.cn/mmopen/vi_32/21ZRkZRa6bHZXT2aeGJSeTLZBdgqSmaWQLrZeyiaChldvoicco9cx2ntp3GQuhLuC8EQXfwNQvhGp0Zicicz7CxNIQ/132'),(642,'http://thirdwx.qlogo.cn/mmopen/vi_32/on5fkN5YibO0fJbXRrClUoFlPEZclibqkwIRqFHWvgrzMCQuwYwxjfIR87RjLh3sDdWyptk9sOS9WZlsq2S7iauaQ/132'),(643,'http://thirdwx.qlogo.cn/mmopen/vi_32/4H9ib5tj4hcLWVZQLanYIzcUb7vwFtOeLK0qnSKBw2KszvuW4tY4C2h2n3Bhtv4FicXYHIoFlRQrMsShPib9AxOdQ/132'),(644,'http://thirdwx.qlogo.cn/mmopen/vi_32/NK86XmGVibqAib7QhO1hyvtebbiaCmHGhnLjxibcBGSLOQqLZZNLuXRyv7v7ohgpCqbamlQ4qoRsUW7kSJSic8KVOsg/132'),(645,'http://thirdwx.qlogo.cn/mmopen/vi_32/SRkpL9OhczYQHITHeskvkKtzP4KrRt2Dh1FCYZqDtguIiaWrlGhblHlXrewS4UiaNhrTgyQed7qqU2ZojeVhwscA/132'),(646,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLibVefuyPnIft2k8IELl0Yz5EbgT1z2mcSqiapkdN238f6K8nokjkiafFzzyCmFibmibzJCxMlICfoS0g/132'),(647,'http://thirdwx.qlogo.cn/mmopen/vi_32/dqcQY1hdliatehhmVuheW7LviasRpIh25UyO9j9ibVBIDiar1WCURz45ZoOBzvIb9Fs9yNKSS9x3XufgG4KetGB1FQ/132'),(648,'http://thirdwx.qlogo.cn/mmopen/vi_32/CEWufvUZXsXHp5oe7aic47YFrSblqhAfQ2HK6bQRmPsbwFXmAPYg2os2GJJPfQur7iawErGicemRdGVuquDmtUXgg/132'),(649,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJru7MZse3ErQxHEC2k9uD8JHicu0LmGibslM5swOFdVwKUQkakiaLAWMsrALxxRicVmSVRk7gYiaahicqg/132'),(650,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83epSFDWnmSqHTS9AS5feu8jWFmG8QrAiaDcpYHOlCV5yfDnPoVJL810yyqoJYzktQcqFrUicwb2l0edA/132'),(651,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTL0pJLGq22zXibVykqRvOVqyDXpuARtTV8bLgSqu3njf9L5oMhtCC81d1k3XAPjDzUbJef55gxib9Kg/132'),(652,'http://thirdwx.qlogo.cn/mmopen/vi_32/UfPbtqrHSQYuRkh1VUkcgvyWWziaJ2ZdaFDibCQKpiaR5ezvlpU3xmmdj7oZkxanibSXGZ486XFNWHa6lQ16m2QXwg/132'),(653,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLSFSE9vNiaI2pZWqWOTtK5Teah4VROp9J0sibZvAuTvJEyCgrq8MwNTn7d0eLvhNc0xNBxqI9kyTEg/132'),(654,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJWSh8ibBBLAZxgIlrOIICnCPKXyOGksD5Jucy5PpwYgfKbbS1mRKW42DwEFq1PhjJsxZSxRefeMiaw/132'),(655,'http://thirdwx.qlogo.cn/mmopen/vi_32/g6OWDLveYspW1KKol1EtZKGqdNy6D5yAibViaMaibQJTqkbIA2OmLlJDjypwNzhxOc9FsIlJcRVbicnibExPx77Wvrw/132'),(656,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKcgAyqufI4rvQMF0ndUpaBicn4f4Q6shd9VPQIclfS8ErJD7MW4iaub2JP5xM6iaGovuNzKoQJJYHMw/132'),(657,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erKay7VOFxBXWHzCxsjO96JlR3rtev6VxHIXhWFhC0YqURwGNyww3Zf27gicMjho3Eoc2exrUy9s7A/132'),(658,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLApvozicUZ866lOJDZFZa5AVWicibibN0bUfmsYe5Hlziaia7VjwJMhxibjDRMQULicEOHn95dDTk0dFficRA/132'),(659,'http://thirdwx.qlogo.cn/mmopen/vi_32/zxqib90MVVzCOnic10EQBffzwGDzagInicKKjVtsxzS5PUjW3WDUTMYJDXce6YJ8wHjD7r54PBFiawmeHRe1noLcSA/132'),(660,'http://thirdwx.qlogo.cn/mmopen/vi_32/uPBIianicySIDZvEO561XMjabWoAzypJ82Rj2tvVkVorAyiaHJcWzoAflkXq6tP12rdw5onxicCNn8QcQa8tMWeHlg/132'),(661,'http://thirdwx.qlogo.cn/mmopen/vi_32/DUWLzREtTML8Oq7TtIfBeSCJUGSLib437e6SwlpOFA7qBjghkHTq7wibDtNYVJZmzBcPV8HzY5T3Rl7Bjia2FuSfA/132'),(662,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIrmemKcibibicgyJxZiaZEUz5XQVHvzbvfu7SiacNFwb47FvACKK2eib4fynibKE2hBS3R7OZO5hh4sULjg/132'),(663,'http://thirdwx.qlogo.cn/mmopen/vi_32/UFBlX5BwfmdbpsTvlqHfOUUCn3q1kGPMhiaGw5Rchqn94UfYiaiaj3JoxyADiaAzQWOJV88FYTK5KPbQA4NRbBhrhg/132'),(664,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIF9hhGyUEhEpphP2JSf2gZpf720r1y6TXHdwkY5l8VInzMDAteozEJFcn0gmNUr8cDGGlWSVZjibg/132'),(665,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLDoP4gTuOnTQlyymGvibKuvoxyzvYBFyESWFEdXiavn59uId1eTMZZYbp7q1mm13xx1JxYdh5UhCmg/132'),(666,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIlZVyciaHNnK6DtJ6ic0XO3sNrWb6O4WoCxZWgPaicjia0A6ZNU6lDDiazcbU7JZB672Z05df4iakIuVcg/132'),(667,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIoEwHBWIVu6Zwia6OKaEbR974NdD1ZNXvvO3DhicceNsr5KB1PkRC3j5NR6JsXquEmOzfK1XE4VFTg/132'),(668,'http://thirdwx.qlogo.cn/mmopen/vi_32/8ZaD4c7uRDdiawPic0qgSOeIUWhzZtmHn9evQs1gW4tufKVhbXu5Y32El90e1R4Sw7b1UkPb77q6p6Hxp2WNxWvw/132'),(669,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTJnwSibIASBJJey55cqU9P9GkbC1FGHyiag4AqHIHDebPqeYMY3RgPGFW1Sbn4X3TJUAOOeqW8GZg3w/132'),(671,'http://thirdwx.qlogo.cn/mmopen/vi_32/u1Opb74xrf2IrSCblQM8bSBicHhpwViaCqgpryMYEibkyibWagqprFc3V1SEHxBdXhUicOJMs4jQLCVYBG3PYbJRnMw/132'),(672,'http://thirdwx.qlogo.cn/mmopen/vi_32/EVAowwWDAUyZiaK52zUMqiboCw2ibOiaQcJsXRVZOfIMib7kMRPeq9OvaiayEI0vEVLZnokYvzAicicrq6BKe2ibDZJoAZg/132'),(673,'http://thirdwx.qlogo.cn/mmopen/vi_32/2RpVI6pp4nZVBHVBicIdbyvXz89iaKhN08jIlfVeKMmOxAJ1yk9Zd9H7icgslDJ7MVjdiaBa773IibEIibQMOpibgwlSA/132'),(674,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTKEwmfyptibj23gDohL2pic3cZQwJETv27Lx6lGhic52YLDF0uppiaLPtgYXIcd9vaib7ybliclEnEoBSzA/132'),(675,'http://thirdwx.qlogo.cn/mmopen/vi_32/uFTEoCxwd1icByJZNf6cYwVo85Rh1ksRF4KyXlNuzAxknGt2cZK3lYUhMcy5kQodK0GmkOLkjmyGfeoRQ8qj5YQ/132'),(676,'http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTLv4gm6QsTByQCOiaoJagNHk5RJWSic0HbKV1KuOcIPg2BxdelXI6ibpVXkWAdjibrY7O45pnuUSw6Vog/132'),(677,'http://thirdwx.qlogo.cn/mmopen/vi_32/DYAIOgq83erFqADeLibFgV8iaLexTgE4bNrfIpFIwKPH0Kr2SlVS7j10xMNAWe17jAApOYvicc46iaKFTuvaPPNlIA/132'),(678,'http://thirdwx.qlogo.cn/mmopen/vi_32/44VNP1Drzl08ZLHosBrvaHuwfml8HUicLIKDqxJTGkltNhAiceWk2iaULspJKsj052xGIdFtiaarWogaJH49iaKGt2Q/132'),(679,'http://thirdwx.qlogo.cn/mmopen/vi_32/01zoZxMDQ97kRWVkbPHw1zfibpSXHuVNUkoXiaqZLPWZJicoDKXIal0wP6YYNuiaXUZO8pU3zvB1khvvMDVtiaRN7cQ/132');
/*!40000 ALTER TABLE `imgurl` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jifen`
--

DROP TABLE IF EXISTS `jifen`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `jifen` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  `daili_id` int(11) DEFAULT NULL,
  `daili_name` varchar(255) DEFAULT NULL,
  `jifen` varchar(255) DEFAULT NULL,
  `time` varchar(225) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jifen`
--

LOCK TABLES `jifen` WRITE;
/*!40000 ALTER TABLE `jifen` DISABLE KEYS */;
/*!40000 ALTER TABLE `jifen` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `kefu`
--

DROP TABLE IF EXISTS `kefu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `kefu` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `kefu`
--

LOCK TABLES `kefu` WRITE;
/*!40000 ALTER TABLE `kefu` DISABLE KEYS */;
INSERT INTO `kefu` VALUES (1,'http://www.baidu.com');
/*!40000 ALTER TABLE `kefu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log`
--

DROP TABLE IF EXISTS `log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  `account` varchar(255) DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log`
--

LOCK TABLES `log` WRITE;
/*!40000 ALTER TABLE `log` DISABLE KEYS */;
/*!40000 ALTER TABLE `log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `login`
--

DROP TABLE IF EXISTS `login`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `login` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` varchar(255) DEFAULT NULL,
  `game` varchar(255) DEFAULT NULL,
  `ftime` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `login`
--

LOCK TABLES `login` WRITE;
/*!40000 ALTER TABLE `login` DISABLE KEYS */;
/*!40000 ALTER TABLE `login` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `matchs`
--

DROP TABLE IF EXISTS `matchs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `matchs` (
  `mid` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) DEFAULT NULL,
  `stime` varchar(255) DEFAULT NULL,
  `etime` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`mid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `matchs`
--

LOCK TABLES `matchs` WRITE;
/*!40000 ALTER TABLE `matchs` DISABLE KEYS */;
/*!40000 ALTER TABLE `matchs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `matchscore`
--

DROP TABLE IF EXISTS `matchscore`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `matchscore` (
  `jid` int(11) NOT NULL AUTO_INCREMENT,
  `mid` int(11) DEFAULT NULL,
  `gid` int(11) DEFAULT NULL,
  `uid` int(11) DEFAULT NULL,
  `score` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`jid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `matchscore`
--

LOCK TABLES `matchscore` WRITE;
/*!40000 ALTER TABLE `matchscore` DISABLE KEYS */;
/*!40000 ALTER TABLE `matchscore` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `move_agent_log`
--

DROP TABLE IF EXISTS `move_agent_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `move_agent_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `moved_id` int(11) DEFAULT NULL COMMENT '被移动者ID',
  `receive_id` int(11) DEFAULT NULL COMMENT '接收者ID',
  `move_time` datetime DEFAULT NULL COMMENT '转移时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `move_agent_log`
--

LOCK TABLES `move_agent_log` WRITE;
/*!40000 ALTER TABLE `move_agent_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `move_agent_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mtcash_log`
--

DROP TABLE IF EXISTS `mtcash_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mtcash_log` (
  `id` int(20) NOT NULL AUTO_INCREMENT,
  `agid` int(20) DEFAULT NULL COMMENT '代理ID',
  `nickname` blob COMMENT '提现者昵称',
  `amount` int(25) DEFAULT NULL COMMENT '提现金额',
  `tcash_time` varchar(225) DEFAULT NULL COMMENT '操作时间',
  `operator` varchar(255) DEFAULT NULL COMMENT '操作者',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mtcash_log`
--

LOCK TABLES `mtcash_log` WRITE;
/*!40000 ALTER TABLE `mtcash_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `mtcash_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notice`
--

DROP TABLE IF EXISTS `notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notice` (
  `notice_id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `content` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`notice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notice`
--

LOCK TABLES `notice` WRITE;
/*!40000 ALTER TABLE `notice` DISABLE KEYS */;
/*!40000 ALTER TABLE `notice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `openter`
--

DROP TABLE IF EXISTS `openter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `openter` (
  `openter_id` int(11) NOT NULL AUTO_INCREMENT,
  `uname` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `tel` varchar(255) DEFAULT NULL,
  `oid` int(11) DEFAULT NULL,
  `beizhu` varchar(255) DEFAULT NULL,
  `weixin` varchar(255) DEFAULT NULL,
  `zhifubao` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  `pid` int(11) DEFAULT NULL,
  `jifen` varchar(255) DEFAULT NULL,
  `account` varchar(255) DEFAULT NULL,
  `sjifen` varchar(255) DEFAULT NULL,
  `paytype` int(11) DEFAULT NULL,
  `product` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`openter_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `openter`
--

LOCK TABLES `openter` WRITE;
/*!40000 ALTER TABLE `openter` DISABLE KEYS */;
/*!40000 ALTER TABLE `openter` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order`
--

DROP TABLE IF EXISTS `order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `order` (
  `order_id` int(11) NOT NULL AUTO_INCREMENT,
  `number` varchar(225) NOT NULL,
  `userid` int(11) NOT NULL,
  `num` varchar(255) NOT NULL,
  `money` varchar(255) NOT NULL,
  `ptype` int(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `paytime` varchar(255) NOT NULL,
  `time` varchar(255) NOT NULL COMMENT '下单时间',
  `paytype` varchar(255) NOT NULL,
  `nickname` blob NOT NULL COMMENT '昵称',
  `pay_desc` int(255) NOT NULL DEFAULT '0' COMMENT '1为微信支付 2为同乐支付 3为威富通支付 4.畅付云',
  PRIMARY KEY (`order_id`),
  KEY `order_no` (`number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order`
--

LOCK TABLES `order` WRITE;
/*!40000 ALTER TABLE `order` DISABLE KEYS */;
/*!40000 ALTER TABLE `order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `otype`
--

DROP TABLE IF EXISTS `otype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `otype` (
  `oid` int(11) NOT NULL AUTO_INCREMENT,
  `ip` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  `status` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`oid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `otype`
--

LOCK TABLES `otype` WRITE;
/*!40000 ALTER TABLE `otype` DISABLE KEYS */;
/*!40000 ALTER TABLE `otype` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `period`
--

DROP TABLE IF EXISTS `period`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `period` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `num` int(11) DEFAULT NULL,
  `ftime` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `period`
--

LOCK TABLES `period` WRITE;
/*!40000 ALTER TABLE `period` DISABLE KEYS */;
/*!40000 ALTER TABLE `period` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pmd`
--

DROP TABLE IF EXISTS `pmd`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pmd` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(225) NOT NULL,
  `adcode` varchar(225) NOT NULL,
  `content` varchar(225) NOT NULL,
  `time` varchar(225) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 CHECKSUM=1 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pmd`
--

LOCK TABLES `pmd` WRITE;
/*!40000 ALTER TABLE `pmd` DISABLE KEYS */;
/*!40000 ALTER TABLE `pmd` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product`
--

DROP TABLE IF EXISTS `product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product` (
  `pid` int(11) NOT NULL AUTO_INCREMENT,
  `pname` varchar(255) DEFAULT NULL,
  `price` decimal(10,2) DEFAULT NULL,
  `dec` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`pid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product`
--

LOCK TABLES `product` WRITE;
/*!40000 ALTER TABLE `product` DISABLE KEYS */;
INSERT INTO `product` VALUES (1,'房卡',1.00,''),(2,'金币',1.00,'');
/*!40000 ALTER TABLE `product` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rebate`
--

DROP TABLE IF EXISTS `rebate`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rebate` (
  `reid` int(11) NOT NULL AUTO_INCREMENT,
  `agids` int(255) DEFAULT NULL,
  `agid` int(11) DEFAULT NULL,
  `ratuo` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`reid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rebate`
--

LOCK TABLES `rebate` WRITE;
/*!40000 ALTER TABLE `rebate` DISABLE KEYS */;
/*!40000 ALTER TABLE `rebate` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sale`
--

DROP TABLE IF EXISTS `sale`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sale` (
  `sale_id` int(11) NOT NULL AUTO_INCREMENT,
  `send_id` int(11) DEFAULT NULL,
  `sender` varchar(255) DEFAULT NULL,
  `send_type` varchar(255) DEFAULT NULL,
  `receive_id` int(11) DEFAULT NULL,
  `receiver` blob,
  `receive_type` varchar(255) DEFAULT NULL,
  `product_id` int(11) DEFAULT NULL,
  `number` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `time` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`sale_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sale`
--

LOCK TABLES `sale` WRITE;
/*!40000 ALTER TABLE `sale` DISABLE KEYS */;
/*!40000 ALTER TABLE `sale` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `send_agent_log`
--

DROP TABLE IF EXISTS `send_agent_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `send_agent_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) DEFAULT NULL COMMENT '1为房卡 2为金币',
  `receive_id` int(11) DEFAULT NULL COMMENT '接收者ID',
  `send_time` varchar(225) DEFAULT NULL COMMENT '操作时间',
  `num` int(11) NOT NULL DEFAULT '0' COMMENT '发送数量',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `send_agent_log`
--

LOCK TABLES `send_agent_log` WRITE;
/*!40000 ALTER TABLE `send_agent_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `send_agent_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `send_gift_log`
--

DROP TABLE IF EXISTS `send_gift_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `send_gift_log` (
  `id` int(255) NOT NULL AUTO_INCREMENT,
  `gift_id` int(10) DEFAULT NULL COMMENT '礼包类型 1为1888金币',
  `uid` int(20) DEFAULT NULL COMMENT '领取用户id',
  `uip` varchar(255) DEFAULT NULL COMMENT '领取用户ip',
  `on_time` varchar(225) DEFAULT NULL COMMENT '领取时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `send_gift_log`
--

LOCK TABLES `send_gift_log` WRITE;
/*!40000 ALTER TABLE `send_gift_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `send_gift_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shop_card`
--

DROP TABLE IF EXISTS `shop_card`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shop_card` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `card` varchar(255) COLLATE utf8_bin NOT NULL,
  `money` varchar(255) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shop_card`
--

LOCK TABLES `shop_card` WRITE;
/*!40000 ALTER TABLE `shop_card` DISABLE KEYS */;
INSERT INTO `shop_card` VALUES (1,'10','11'),(2,'50','50'),(3,'100','100'),(4,'200','200');
/*!40000 ALTER TABLE `shop_card` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shop_cft`
--

DROP TABLE IF EXISTS `shop_cft`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shop_cft` (
  `id` int(11) NOT NULL,
  `uid` varchar(225) COLLATE utf8_bin NOT NULL,
  `key` varchar(255) COLLATE utf8_bin NOT NULL,
  `notify_url` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shop_cft`
--

LOCK TABLES `shop_cft` WRITE;
/*!40000 ALTER TABLE `shop_cft` DISABLE KEYS */;
INSERT INTO `shop_cft` VALUES (1,'123','23',NULL);
/*!40000 ALTER TABLE `shop_cft` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shop_dsf`
--

DROP TABLE IF EXISTS `shop_dsf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shop_dsf` (
  `id` int(11) NOT NULL,
  `type` int(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shop_dsf`
--

LOCK TABLES `shop_dsf` WRITE;
/*!40000 ALTER TABLE `shop_dsf` DISABLE KEYS */;
INSERT INTO `shop_dsf` VALUES (1,3);
/*!40000 ALTER TABLE `shop_dsf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shop_gold`
--

DROP TABLE IF EXISTS `shop_gold`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shop_gold` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `gold` varchar(255) COLLATE utf8_bin NOT NULL,
  `type` int(2) NOT NULL COMMENT '1:支付宝 2 微信  3QQ  4 银联',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=34 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shop_gold`
--

LOCK TABLES `shop_gold` WRITE;
/*!40000 ALTER TABLE `shop_gold` DISABLE KEYS */;
INSERT INTO `shop_gold` VALUES (1,'10',1),(2,'50',1),(3,'100',1),(4,'200',1),(5,'300',1),(6,'400',1),(7,'500',1),(8,'0',1),(9,'10',2),(10,'50',2),(11,'100',2),(12,'200',2),(13,'300',2),(14,'400',2),(15,'500',2),(16,'600',2),(17,'20',3),(18,'50',3),(19,'100',3),(20,'0',3),(21,'0',3),(22,'0',3),(23,'500',3),(24,'600',3),(25,'10',4),(26,'0',4),(27,'0',4),(28,'0',4),(29,'0',4),(30,'0',4),(31,'0',4),(32,'0',4),(33,'',0);
/*!40000 ALTER TABLE `shop_gold` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_address_city`
--

DROP TABLE IF EXISTS `t_address_city`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_address_city` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` char(6) NOT NULL COMMENT '城市编码',
  `name` varchar(40) NOT NULL COMMENT '城市名称',
  `provinceCode` char(6) NOT NULL COMMENT '所属省份编码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=329 DEFAULT CHARSET=utf8 COMMENT='城市信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_address_city`
--

LOCK TABLES `t_address_city` WRITE;
/*!40000 ALTER TABLE `t_address_city` DISABLE KEYS */;
INSERT INTO `t_address_city` VALUES (1,'110100','北京市','110000'),(2,'1102xx','北京下属县','1100xx'),(3,'120100','天津市','120000'),(4,'1202xx','天津下属县','1200xx'),(5,'130100','石家庄市','130000'),(6,'130200','唐山市','130000'),(7,'130300','秦皇岛市','130000'),(8,'130400','邯郸市','130000'),(9,'130500','邢台市','130000'),(10,'130600','保定市','130000'),(11,'130700','张家口市','130000'),(12,'130800','承德市','130000'),(13,'130900','沧州市','130000'),(14,'131000','廊坊市','130000'),(15,'131100','衡水市','130000'),(16,'140100','太原市','140000'),(17,'140200','大同市','140000'),(18,'140300','阳泉市','140000'),(19,'140400','长治市','140000'),(20,'140500','晋城市','140000'),(21,'140600','朔州市','140000'),(22,'140700','晋中市','140000'),(23,'140800','运城市','140000'),(24,'140900','忻州市','140000'),(25,'141000','临汾市','140000'),(26,'141100','吕梁市','140000'),(27,'150100','呼和浩特市','150000'),(28,'150200','包头市','150000'),(29,'150300','乌海市','150000'),(30,'150400','赤峰市','150000'),(31,'150500','通辽市','150000'),(32,'150600','鄂尔多斯市','150000'),(33,'150700','呼伦贝尔市','150000'),(34,'150800','巴彦淖尔市','150000'),(35,'150900','乌兰察布市','150000'),(36,'152200','兴安盟','150000'),(37,'152500','锡林郭勒盟','150000'),(38,'152900','阿拉善盟','150000'),(39,'210100','沈阳市','210000'),(40,'210200','大连市','210000'),(41,'210300','鞍山市','210000'),(42,'210400','抚顺市','210000'),(43,'210500','本溪市','210000'),(44,'210600','丹东市','210000'),(45,'210700','锦州市','210000'),(46,'210800','营口市','210000'),(47,'210900','阜新市','210000'),(48,'211000','辽阳市','210000'),(49,'211100','盘锦市','210000'),(50,'211200','铁岭市','210000'),(51,'211300','朝阳市','210000'),(52,'211400','葫芦岛市','210000'),(53,'220100','长春市','220000'),(54,'220200','吉林市','220000'),(55,'220300','四平市','220000'),(56,'220400','辽源市','220000'),(57,'220500','通化市','220000'),(58,'220600','白山市','220000'),(59,'220700','松原市','220000'),(60,'220800','白城市','220000'),(61,'222400','延边朝鲜族自治州','220000'),(62,'230100','哈尔滨市','230000'),(63,'230200','齐齐哈尔市','230000'),(64,'230300','鸡西市','230000'),(65,'230400','鹤岗市','230000'),(66,'230500','双鸭山市','230000'),(67,'230600','大庆市','230000'),(68,'230700','伊春市','230000'),(69,'230800','佳木斯市','230000'),(70,'230900','七台河市','230000'),(71,'231000','牡丹江市','230000'),(72,'231100','黑河市','230000'),(73,'231200','绥化市','230000'),(74,'232700','大兴安岭地区','230000'),(75,'310100','上海市','310000'),(76,'3102xx','上海下属县','3100xx'),(77,'320100','南京市','320000'),(78,'320200','无锡市','320000'),(79,'320300','徐州市','320000'),(80,'320400','常州市','320000'),(81,'320500','苏州市','320000'),(82,'320600','南通市','320000'),(83,'320700','连云港市','320000'),(84,'320800','淮安市','320000'),(85,'320900','盐城市','320000'),(86,'321000','扬州市','320000'),(87,'321100','镇江市','320000'),(88,'321200','泰州市','320000'),(89,'321300','宿迁市','320000'),(90,'330100','杭州市','330000'),(91,'330200','宁波市','330000'),(92,'330300','温州市','330000'),(93,'330400','嘉兴市','330000'),(94,'330500','湖州市','330000'),(95,'330600','绍兴市','330000'),(96,'330700','金华市','330000'),(97,'330800','衢州市','330000'),(98,'330900','舟山市','330000'),(99,'331000','台州市','330000'),(100,'331100','丽水市','330000'),(101,'340100','合肥市','340000'),(102,'340200','芜湖市','340000'),(103,'340300','蚌埠市','340000'),(104,'340400','淮南市','340000'),(105,'340500','马鞍山市','340000'),(106,'340600','淮北市','340000'),(107,'340700','铜陵市','340000'),(108,'340800','安庆市','340000'),(109,'341000','黄山市','340000'),(110,'341100','滁州市','340000'),(111,'341200','阜阳市','340000'),(112,'341300','宿州市','340000'),(113,'341400','巢湖市','340000'),(114,'341500','六安市','340000'),(115,'341600','亳州市','340000'),(116,'341700','池州市','340000'),(117,'341800','宣城市','340000'),(118,'350100','福州市','350000'),(119,'350200','厦门市','350000'),(120,'350300','莆田市','350000'),(121,'350400','三明市','350000'),(122,'350500','泉州市','350000'),(123,'350600','漳州市','350000'),(124,'350700','南平市','350000'),(125,'350800','龙岩市','350000'),(126,'350900','宁德市','350000'),(127,'360100','南昌市','360000'),(128,'360200','景德镇市','360000'),(129,'360300','萍乡市','360000'),(130,'360400','九江市','360000'),(131,'360500','新余市','360000'),(132,'360600','鹰潭市','360000'),(133,'360700','赣州市','360000'),(134,'360800','吉安市','360000'),(135,'360900','宜春市','360000'),(136,'361000','抚州市','360000'),(137,'361100','上饶市','360000'),(138,'370100','济南市','370000'),(139,'370200','青岛市','370000'),(140,'370300','淄博市','370000'),(141,'370400','枣庄市','370000'),(142,'370500','东营市','370000'),(143,'370600','烟台市','370000'),(144,'370700','潍坊市','370000'),(145,'370800','济宁市','370000'),(146,'370900','泰安市','370000'),(147,'371000','威海市','370000'),(148,'371100','日照市','370000'),(149,'371200','莱芜市','370000'),(150,'371300','临沂市','370000'),(151,'371400','德州市','370000'),(152,'371500','聊城市','370000'),(153,'371600','滨州市','370000'),(154,'371700','荷泽市','370000'),(155,'410100','郑州市','410000'),(156,'410200','开封市','410000'),(157,'410300','洛阳市','410000'),(158,'410400','平顶山市','410000'),(159,'410500','安阳市','410000'),(160,'410600','鹤壁市','410000'),(161,'410700','新乡市','410000'),(162,'410800','焦作市','410000'),(163,'410900','濮阳市','410000'),(164,'411000','许昌市','410000'),(165,'411100','漯河市','410000'),(166,'411200','三门峡市','410000'),(167,'411300','南阳市','410000'),(168,'411400','商丘市','410000'),(169,'411500','信阳市','410000'),(170,'411600','周口市','410000'),(171,'411700','驻马店市','410000'),(172,'420100','武汉市','420000'),(173,'420200','黄石市','420000'),(174,'420300','十堰市','420000'),(175,'420500','宜昌市','420000'),(176,'420600','襄樊市','420000'),(177,'420700','鄂州市','420000'),(178,'420800','荆门市','420000'),(179,'420900','孝感市','420000'),(180,'421000','荆州市','420000'),(181,'421100','黄冈市','420000'),(182,'421200','咸宁市','420000'),(183,'421300','随州市','420000'),(184,'422800','恩施土家族苗族自治州','420000'),(185,'429000','省直辖行政单位','420000'),(186,'430100','长沙市','430000'),(187,'430200','株洲市','430000'),(188,'430300','湘潭市','430000'),(189,'430400','衡阳市','430000'),(190,'430500','邵阳市','430000'),(191,'430600','岳阳市','430000'),(192,'430700','常德市','430000'),(193,'430800','张家界市','430000'),(194,'430900','益阳市','430000'),(195,'431000','郴州市','430000'),(196,'431100','永州市','430000'),(197,'431200','怀化市','430000'),(198,'431300','娄底市','430000'),(199,'433100','湘西土家族苗族自治州','430000'),(200,'440100','广州市','440000'),(201,'440200','韶关市','440000'),(202,'440300','深圳市','440000'),(203,'440400','珠海市','440000'),(204,'440500','汕头市','440000'),(205,'440600','佛山市','440000'),(206,'440700','江门市','440000'),(207,'440800','湛江市','440000'),(208,'440900','茂名市','440000'),(209,'441200','肇庆市','440000'),(210,'441300','惠州市','440000'),(211,'441400','梅州市','440000'),(212,'441500','汕尾市','440000'),(213,'441600','河源市','440000'),(214,'441700','阳江市','440000'),(215,'441800','清远市','440000'),(216,'441900','东莞市','440000'),(217,'442000','中山市','440000'),(218,'445100','潮州市','440000'),(219,'445200','揭阳市','440000'),(220,'445300','云浮市','440000'),(221,'450100','南宁市','450000'),(222,'450200','柳州市','450000'),(223,'450300','桂林市','450000'),(224,'450400','梧州市','450000'),(225,'450500','北海市','450000'),(226,'450600','防城港市','450000'),(227,'450700','钦州市','450000'),(228,'450800','贵港市','450000'),(229,'450900','玉林市','450000'),(230,'451000','百色市','450000'),(231,'451100','贺州市','450000'),(232,'451200','河池市','450000'),(233,'451300','来宾市','450000'),(234,'451400','崇左市','450000'),(235,'460100','海口市','460000'),(236,'460200','三亚市','460000'),(237,'469000','省直辖县级行政单位','460000'),(238,'500100','重庆市','500000'),(239,'5002xx','重庆下属县','5000xx'),(240,'5003xx','重庆下属市','5000xx'),(241,'510100','成都市','510000'),(242,'510300','自贡市','510000'),(243,'510400','攀枝花市','510000'),(244,'510500','泸州市','510000'),(245,'510600','德阳市','510000'),(246,'510700','绵阳市','510000'),(247,'510800','广元市','510000'),(248,'510900','遂宁市','510000'),(249,'511000','内江市','510000'),(250,'511100','乐山市','510000'),(251,'511300','南充市','510000'),(252,'511400','眉山市','510000'),(253,'511500','宜宾市','510000'),(254,'511600','广安市','510000'),(255,'511700','达州市','510000'),(256,'511800','雅安市','510000'),(257,'511900','巴中市','510000'),(258,'512000','资阳市','510000'),(259,'513200','阿坝藏族羌族自治州','510000'),(260,'513300','甘孜藏族自治州','510000'),(261,'513400','凉山彝族自治州','510000'),(262,'520100','贵阳市','520000'),(263,'520200','六盘水市','520000'),(264,'520300','遵义市','520000'),(265,'520400','安顺市','520000'),(266,'522200','铜仁地区','520000'),(267,'522300','黔西南布依族苗族自治州','520000'),(268,'522400','毕节地区','520000'),(269,'522600','黔东南苗族侗族自治州','520000'),(270,'522700','黔南布依族苗族自治州','520000'),(271,'530100','昆明市','530000'),(272,'530300','曲靖市','530000'),(273,'530400','玉溪市','530000'),(274,'530500','保山市','530000'),(275,'530600','昭通市','530000'),(276,'530700','丽江市','530000'),(277,'530800','思茅市','530000'),(278,'530900','临沧市','530000'),(279,'532300','楚雄彝族自治州','530000'),(280,'532500','红河哈尼族彝族自治州','530000'),(281,'532600','文山壮族苗族自治州','530000'),(282,'532800','西双版纳傣族自治州','530000'),(283,'532900','大理白族自治州','530000'),(284,'533100','德宏傣族景颇族自治州','530000'),(285,'533300','怒江傈僳族自治州','530000'),(286,'533400','迪庆藏族自治州','530000'),(287,'540100','拉萨市','540000'),(288,'542100','昌都地区','540000'),(289,'542200','山南地区','540000'),(290,'542300','日喀则地区','540000'),(291,'542400','那曲地区','540000'),(292,'542500','阿里地区','540000'),(293,'542600','林芝地区','540000'),(294,'610100','西安市','610000'),(295,'610200','铜川市','610000'),(296,'610300','宝鸡市','610000'),(297,'610400','咸阳市','610000'),(298,'610500','渭南市','610000'),(299,'610600','延安市','610000'),(300,'610700','汉中市','610000'),(301,'610800','榆林市','610000'),(302,'610900','安康市','610000'),(303,'611000','商洛市','610000'),(304,'620100','兰州市','620000'),(305,'620200','嘉峪关市','620000'),(306,'620300','金昌市','620000'),(307,'620400','白银市','620000'),(308,'620500','天水市','620000'),(309,'620600','武威市','620000'),(310,'620700','张掖市','620000'),(311,'620800','平凉市','620000'),(312,'620900','酒泉市','620000'),(313,'621000','庆阳市','620000'),(314,'621100','定西市','620000'),(315,'621200','陇南市','620000'),(316,'622900','临夏回族自治州','620000'),(317,'623000','甘南藏族自治州','620000'),(318,'630100','西宁市','630000'),(319,'632100','海东地区','630000'),(320,'632200','海北藏族自治州','630000'),(321,'632300','黄南藏族自治州','630000'),(322,'632500','海南藏族自治州','630000'),(323,'632600','果洛藏族自治州','630000'),(324,'632700','玉树藏族自治州','630000'),(325,'632800','海西蒙古族藏族自治州','630000'),(326,'640100','银川市','640000'),(327,'640200','石嘴山市','640000'),(328,'640300','吴忠市','640000');
/*!40000 ALTER TABLE `t_address_city` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_address_province`
--

DROP TABLE IF EXISTS `t_address_province`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_address_province` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` char(6) NOT NULL COMMENT '省份编码',
  `name` varchar(40) NOT NULL COMMENT '省份名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8 COMMENT='省份信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_address_province`
--

LOCK TABLES `t_address_province` WRITE;
/*!40000 ALTER TABLE `t_address_province` DISABLE KEYS */;
INSERT INTO `t_address_province` VALUES (1,'110000','北京市'),(2,'120000','天津市'),(3,'130000','河北省'),(4,'140000','山西省'),(5,'150000','内蒙古自治区'),(6,'210000','辽宁省'),(7,'220000','吉林省'),(8,'230000','黑龙江省'),(9,'310000','上海市'),(10,'320000','江苏省'),(11,'330000','浙江省'),(12,'340000','安徽省'),(13,'350000','福建省'),(14,'360000','江西省'),(15,'370000','山东省'),(16,'410000','河南省'),(17,'420000','湖北省'),(18,'430000','湖南省'),(19,'440000','广东省'),(20,'450000','广西壮族自治区'),(21,'460000','海南省'),(22,'500000','重庆市'),(23,'510000','四川省'),(24,'520000','贵州省'),(25,'530000','云南省'),(26,'540000','西藏自治区'),(27,'610000','陕西省'),(28,'620000','甘肃省'),(29,'630000','青海省'),(30,'640000','宁夏回族自治区'),(31,'650000','新疆维吾尔自治区'),(32,'710000','台湾省'),(33,'810000','香港特别行政区'),(34,'820000','澳门特别行政区');
/*!40000 ALTER TABLE `t_address_province` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_address_town`
--

DROP TABLE IF EXISTS `t_address_town`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_address_town` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` char(6) NOT NULL COMMENT '区县编码',
  `name` varchar(40) NOT NULL COMMENT '区县名称',
  `cityCode` char(6) NOT NULL COMMENT '所属城市编码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=437 DEFAULT CHARSET=utf8 COMMENT='区县信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_address_town`
--

LOCK TABLES `t_address_town` WRITE;
/*!40000 ALTER TABLE `t_address_town` DISABLE KEYS */;
INSERT INTO `t_address_town` VALUES (1,'110101','东城区','110100'),(2,'110102','西城区','110100'),(3,'110103','崇文区','110100'),(4,'110104','宣武区','110100'),(5,'110105','朝阳区','110100'),(6,'110106','丰台区','110100'),(7,'110107','石景山区','110100'),(8,'110108','海淀区','110100'),(9,'110109','门头沟区','110100'),(10,'110111','房山区','110100'),(11,'110112','通州区','110100'),(12,'110113','顺义区','110100'),(13,'110114','昌平区','110100'),(14,'110115','大兴区','110100'),(15,'110116','怀柔区','110100'),(16,'110117','平谷区','110100'),(17,'110228','密云县','110100'),(18,'110229','延庆县','110100'),(19,'120101','和平区','120100'),(20,'120102','河东区','120100'),(21,'120103','河西区','120100'),(22,'120104','南开区','120100'),(23,'120105','河北区','120100'),(24,'120106','红桥区','120100'),(25,'120107','塘沽区','120100'),(26,'120108','汉沽区','120100'),(27,'120109','大港区','120100'),(28,'120110','东丽区','120100'),(29,'120111','西青区','120100'),(30,'120112','津南区','120100'),(31,'120113','北辰区','120100'),(32,'120114','武清区','120100'),(33,'120115','宝坻区','120100'),(34,'120221','宁河县','120100'),(35,'120223','静海县','120100'),(36,'120225','蓟　县','120100'),(37,'130101','市辖区','130100'),(38,'130102','长安区','130100'),(39,'130103','桥东区','130100'),(40,'130104','桥西区','130100'),(41,'130105','新华区','130100'),(42,'130107','井陉矿区','130100'),(43,'130108','裕华区','130100'),(44,'130121','井陉县','130100'),(45,'130123','正定县','130100'),(46,'130124','栾城县','130100'),(47,'130125','行唐县','130100'),(48,'130126','灵寿县','130100'),(49,'130127','高邑县','130100'),(50,'130128','深泽县','130100'),(51,'130129','赞皇县','130100'),(52,'130130','无极县','130100'),(53,'130131','平山县','130100'),(54,'130132','元氏县','130100'),(55,'130133','赵　县','130100'),(56,'130181','辛集市','130100'),(57,'130182','藁城市','130100'),(58,'130183','晋州市','130100'),(59,'130184','新乐市','130100'),(60,'130185','鹿泉市','130100'),(61,'130201','市辖区','130200'),(62,'130202','路南区','130200'),(63,'130203','路北区','130200'),(64,'130204','古冶区','130200'),(65,'130205','开平区','130200'),(66,'130207','丰南区','130200'),(67,'130208','丰润区','130200'),(68,'130223','滦　县','130200'),(69,'130224','滦南县','130200'),(70,'130225','乐亭县','130200'),(71,'130227','迁西县','130200'),(72,'130229','玉田县','130200'),(73,'130230','唐海县','130200'),(74,'130281','遵化市','130200'),(75,'130283','迁安市','130200'),(76,'130301','市辖区','130300'),(77,'130302','海港区','130300'),(78,'130303','山海关区','130300'),(79,'130304','北戴河区','130300'),(80,'130321','青龙满族自治县','130300'),(81,'130322','昌黎县','130300'),(82,'130323','抚宁县','130300'),(83,'130324','卢龙县','130300'),(84,'130401','市辖区','130400'),(85,'130402','邯山区','130400'),(86,'130403','丛台区','130400'),(87,'130404','复兴区','130400'),(88,'130406','峰峰矿区','130400'),(89,'130421','邯郸县','130400'),(90,'130423','临漳县','130400'),(91,'130424','成安县','130400'),(92,'130425','大名县','130400'),(93,'130426','涉　县','130400'),(94,'130427','磁　县','130400'),(95,'130428','肥乡县','130400'),(96,'130429','永年县','130400'),(97,'130430','邱　县','130400'),(98,'130431','鸡泽县','130400'),(99,'130432','广平县','130400'),(100,'130433','馆陶县','130400'),(101,'130434','魏　县','130400'),(102,'130435','曲周县','130400'),(103,'130481','武安市','130400'),(104,'130501','市辖区','130500'),(105,'130502','桥东区','130500'),(106,'130503','桥西区','130500'),(107,'130521','邢台县','130500'),(108,'130522','临城县','130500'),(109,'130523','内丘县','130500'),(110,'130524','柏乡县','130500'),(111,'130525','隆尧县','130500'),(112,'130526','任　县','130500'),(113,'130527','南和县','130500'),(114,'130528','宁晋县','130500'),(115,'130529','巨鹿县','130500'),(116,'130530','新河县','130500'),(117,'130531','广宗县','130500'),(118,'130532','平乡县','130500'),(119,'130533','威　县','130500'),(120,'130534','清河县','130500'),(121,'130535','临西县','130500'),(122,'130581','南宫市','130500'),(123,'130582','沙河市','130500'),(124,'130601','市辖区','130600'),(125,'130602','新市区','130600'),(126,'130603','北市区','130600'),(127,'130604','南市区','130600'),(128,'130621','满城县','130600'),(129,'130622','清苑县','130600'),(130,'130623','涞水县','130600'),(131,'130624','阜平县','130600'),(132,'130625','徐水县','130600'),(133,'130626','定兴县','130600'),(134,'130627','唐　县','130600'),(135,'130628','高阳县','130600'),(136,'130629','容城县','130600'),(137,'130630','涞源县','130600'),(138,'130631','望都县','130600'),(139,'130632','安新县','130600'),(140,'130633','易　县','130600'),(141,'130634','曲阳县','130600'),(142,'130635','蠡　县','130600'),(143,'130636','顺平县','130600'),(144,'130637','博野县','130600'),(145,'130638','雄　县','130600'),(146,'130681','涿州市','130600'),(147,'130682','定州市','130600'),(148,'130683','安国市','130600'),(149,'130684','高碑店市','130600'),(150,'130701','市辖区','130700'),(151,'130702','桥东区','130700'),(152,'130703','桥西区','130700'),(153,'130705','宣化区','130700'),(154,'130706','下花园区','130700'),(155,'130721','宣化县','130700'),(156,'130722','张北县','130700'),(157,'130723','康保县','130700'),(158,'130724','沽源县','130700'),(159,'130725','尚义县','130700'),(160,'130726','蔚　县','130700'),(161,'130727','阳原县','130700'),(162,'130728','怀安县','130700'),(163,'130729','万全县','130700'),(164,'130730','怀来县','130700'),(165,'130731','涿鹿县','130700'),(166,'130732','赤城县','130700'),(167,'130733','崇礼县','130700'),(168,'130801','市辖区','130800'),(169,'130802','双桥区','130800'),(170,'130803','双滦区','130800'),(171,'130804','鹰手营子矿区','130800'),(172,'130821','承德县','130800'),(173,'130822','兴隆县','130800'),(174,'130823','平泉县','130800'),(175,'130824','滦平县','130800'),(176,'130825','隆化县','130800'),(177,'130826','丰宁满族自治县','130800'),(178,'130827','宽城满族自治县','130800'),(179,'130828','围场满族蒙古族自治县','130800'),(180,'130901','市辖区','130900'),(181,'130902','新华区','130900'),(182,'130903','运河区','130900'),(183,'130921','沧　县','130900'),(184,'130922','青　县','130900'),(185,'130923','东光县','130900'),(186,'130924','海兴县','130900'),(187,'130925','盐山县','130900'),(188,'130926','肃宁县','130900'),(189,'130927','南皮县','130900'),(190,'130928','吴桥县','130900'),(191,'130929','献　县','130900'),(192,'130930','孟村回族自治县','130900'),(193,'130981','泊头市','130900'),(194,'130982','任丘市','130900'),(195,'130983','黄骅市','130900'),(196,'130984','河间市','130900'),(197,'131001','市辖区','131000'),(198,'131002','安次区','131000'),(199,'131003','广阳区','131000'),(200,'131022','固安县','131000'),(201,'131023','永清县','131000'),(202,'131024','香河县','131000'),(203,'131025','大城县','131000'),(204,'131026','文安县','131000'),(205,'131028','大厂回族自治县','131000'),(206,'131081','霸州市','131000'),(207,'131082','三河市','131000'),(208,'131101','市辖区','131100'),(209,'131102','桃城区','131100'),(210,'131121','枣强县','131100'),(211,'131122','武邑县','131100'),(212,'131123','武强县','131100'),(213,'131124','饶阳县','131100'),(214,'131125','安平县','131100'),(215,'131126','故城县','131100'),(216,'131127','景　县','131100'),(217,'131128','阜城县','131100'),(218,'131181','冀州市','131100'),(219,'131182','深州市','131100'),(220,'140101','市辖区','140100'),(221,'140105','小店区','140100'),(222,'140106','迎泽区','140100'),(223,'140107','杏花岭区','140100'),(224,'140108','尖草坪区','140100'),(225,'140109','万柏林区','140100'),(226,'140110','晋源区','140100'),(227,'140121','清徐县','140100'),(228,'140122','阳曲县','140100'),(229,'140123','娄烦县','140100'),(230,'140181','古交市','140100'),(231,'140201','市辖区','140200'),(232,'140202','城　区','140200'),(233,'140203','矿　区','140200'),(234,'140211','南郊区','140200'),(235,'140212','新荣区','140200'),(236,'140221','阳高县','140200'),(237,'140222','天镇县','140200'),(238,'140223','广灵县','140200'),(239,'140224','灵丘县','140200'),(240,'140225','浑源县','140200'),(241,'140226','左云县','140200'),(242,'140227','大同县','140200'),(243,'140301','市辖区','140300'),(244,'140302','城　区','140300'),(245,'140303','矿　区','140300'),(246,'140311','郊　区','140300'),(247,'140321','平定县','140300'),(248,'140322','盂　县','140300'),(249,'140401','市辖区','140400'),(250,'140402','城　区','140400'),(251,'140411','郊　区','140400'),(252,'140421','长治县','140400'),(253,'140423','襄垣县','140400'),(254,'140424','屯留县','140400'),(255,'140425','平顺县','140400'),(256,'140426','黎城县','140400'),(257,'140427','壶关县','140400'),(258,'140428','长子县','140400'),(259,'140429','武乡县','140400'),(260,'140430','沁　县','140400'),(261,'140431','沁源县','140400'),(262,'140481','潞城市','140400'),(263,'140501','市辖区','140500'),(264,'140502','城　区','140500'),(265,'140521','沁水县','140500'),(266,'140522','阳城县','140500'),(267,'140524','陵川县','140500'),(268,'140525','泽州县','140500'),(269,'140581','高平市','140500'),(270,'140601','市辖区','140600'),(271,'140602','朔城区','140600'),(272,'140603','平鲁区','140600'),(273,'140621','山阴县','140600'),(274,'140622','应　县','140600'),(275,'140623','右玉县','140600'),(276,'140624','怀仁县','140600'),(277,'140701','市辖区','140700'),(278,'140702','榆次区','140700'),(279,'140721','榆社县','140700'),(280,'140722','左权县','140700'),(281,'140723','和顺县','140700'),(282,'140724','昔阳县','140700'),(283,'140725','寿阳县','140700'),(284,'140726','太谷县','140700'),(285,'140727','祁　县','140700'),(286,'140728','平遥县','140700'),(287,'140729','灵石县','140700'),(288,'140781','介休市','140700'),(289,'140801','市辖区','140800'),(290,'140802','盐湖区','140800'),(291,'140821','临猗县','140800'),(292,'140822','万荣县','140800'),(293,'140823','闻喜县','140800'),(294,'140824','稷山县','140800'),(295,'140825','新绛县','140800'),(296,'140826','绛　县','140800'),(297,'140827','垣曲县','140800'),(298,'140828','夏　县','140800'),(299,'140829','平陆县','140800'),(300,'140830','芮城县','140800'),(301,'140881','永济市','140800'),(302,'140882','河津市','140800'),(303,'140901','市辖区','140900'),(304,'140902','忻府区','140900'),(305,'140921','定襄县','140900'),(306,'140922','五台县','140900'),(307,'140923','代　县','140900'),(308,'140924','繁峙县','140900'),(309,'140925','宁武县','140900'),(310,'140926','静乐县','140900'),(311,'140927','神池县','140900'),(312,'140928','五寨县','140900'),(313,'140929','岢岚县','140900'),(314,'140930','河曲县','140900'),(315,'140931','保德县','140900'),(316,'140932','偏关县','140900'),(317,'140981','原平市','140900'),(318,'141001','市辖区','141000'),(319,'141002','尧都区','141000'),(320,'141021','曲沃县','141000'),(321,'141022','翼城县','141000'),(322,'141023','襄汾县','141000'),(323,'141024','洪洞县','141000'),(324,'141025','古　县','141000'),(325,'141026','安泽县','141000'),(326,'141027','浮山县','141000'),(327,'141028','吉　县','141000'),(328,'141029','乡宁县','141000'),(329,'141030','大宁县','141000'),(330,'141031','隰　县','141000'),(331,'141032','永和县','141000'),(332,'141033','蒲　县','141000'),(333,'141034','汾西县','141000'),(334,'141081','侯马市','141000'),(335,'141082','霍州市','141000'),(336,'141101','市辖区','141100'),(337,'141102','离石区','141100'),(338,'141121','文水县','141100'),(339,'141122','交城县','141100'),(340,'141123','兴　县','141100'),(341,'141124','临　县','141100'),(342,'141125','柳林县','141100'),(343,'141126','石楼县','141100'),(344,'141127','岚　县','141100'),(345,'141128','方山县','141100'),(346,'141129','中阳县','141100'),(347,'141130','交口县','141100'),(348,'141181','孝义市','141100'),(349,'141182','汾阳市','141100'),(350,'150101','市辖区','150100'),(351,'150102','新城区','150100'),(352,'150103','回民区','150100'),(353,'150104','玉泉区','150100'),(354,'150105','赛罕区','150100'),(355,'150121','土默特左旗','150100'),(356,'150122','托克托县','150100'),(357,'150123','和林格尔县','150100'),(358,'150124','清水河县','150100'),(359,'150125','武川县','150100'),(360,'150201','市辖区','150200'),(361,'150202','东河区','150200'),(362,'150203','昆都仑区','150200'),(363,'150204','青山区','150200'),(364,'150205','石拐区','150200'),(365,'150206','白云矿区','150200'),(366,'150207','九原区','150200'),(367,'150221','土默特右旗','150200'),(368,'150222','固阳县','150200'),(369,'150223','达尔罕茂明安联合旗','150200'),(370,'150301','市辖区','150300'),(371,'150302','海勃湾区','150300'),(372,'150303','海南区','150300'),(373,'150304','乌达区','150300'),(374,'150401','市辖区','150400'),(375,'150402','红山区','150400'),(376,'150403','元宝山区','150400'),(377,'150404','松山区','150400'),(378,'150421','阿鲁科尔沁旗','150400'),(379,'150422','巴林左旗','150400'),(380,'150423','巴林右旗','150400'),(381,'150424','林西县','150400'),(382,'150425','克什克腾旗','150400'),(383,'150426','翁牛特旗','150400'),(384,'150428','喀喇沁旗','150400'),(385,'150429','宁城县','150400'),(386,'150430','敖汉旗','150400'),(387,'150501','市辖区','150500'),(388,'150502','科尔沁区','150500'),(389,'150521','科尔沁左翼中旗','150500'),(390,'150522','科尔沁左翼后旗','150500'),(391,'150523','开鲁县','150500'),(392,'150524','库伦旗','150500'),(393,'150525','奈曼旗','150500'),(394,'150526','扎鲁特旗','150500'),(395,'150581','霍林郭勒市','150500'),(396,'150602','东胜区','150600'),(397,'150621','达拉特旗','150600'),(398,'150622','准格尔旗','150600'),(399,'150623','鄂托克前旗','150600'),(400,'150624','鄂托克旗','150600'),(401,'150625','杭锦旗','150600'),(402,'150626','乌审旗','150600'),(403,'150627','伊金霍洛旗','150600'),(404,'150701','市辖区','150700'),(405,'150702','海拉尔区','150700'),(406,'150721','阿荣旗','150700'),(407,'150722','莫力达瓦达斡尔族自治旗','150700'),(408,'150723','鄂伦春自治旗','150700'),(409,'150724','鄂温克族自治旗','150700'),(410,'150725','陈巴尔虎旗','150700'),(411,'150726','新巴尔虎左旗','150700'),(412,'150727','新巴尔虎右旗','150700'),(413,'150781','满洲里市','150700'),(414,'150782','牙克石市','150700'),(415,'150783','扎兰屯市','150700'),(416,'150784','额尔古纳市','150700'),(417,'150785','根河市','150700'),(418,'150801','市辖区','150800'),(419,'150802','临河区','150800'),(420,'150821','五原县','150800'),(421,'150822','磴口县','150800'),(422,'150823','乌拉特前旗','150800'),(423,'150824','乌拉特中旗','150800'),(424,'150825','乌拉特后旗','150800'),(425,'150826','杭锦后旗','150800'),(426,'150901','市辖区','150900'),(427,'150902','集宁区','150900'),(428,'150921','卓资县','150900'),(429,'150922','化德县','150900'),(430,'150923','商都县','150900'),(431,'150924','兴和县','150900'),(432,'150925','凉城县','150900'),(433,'150926','察哈尔右翼前旗','150900'),(434,'150927','察哈尔右翼中旗','150900'),(435,'150928','察哈尔右翼后旗','150900'),(436,'150929','四子王旗','150900');
/*!40000 ALTER TABLE `t_address_town` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `userid` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `uname` varchar(255) NOT NULL,
  `tel` varchar(255) DEFAULT NULL,
  `wx` varchar(255) DEFAULT NULL,
  `auth` varchar(255) NOT NULL,
  `role` varchar(255) NOT NULL,
  `time` varchar(255) DEFAULT NULL,
  `beizhu` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'admin','123456','总管理','10086','ypdasds','1','1','2017-06-27','');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_notice`
--

DROP TABLE IF EXISTS `user_notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_notice` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `agent_id` int(11) DEFAULT NULL,
  `notice_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_notice`
--

LOCK TABLES `user_notice` WRITE;
/*!40000 ALTER TABLE `user_notice` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_notice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_show`
--

DROP TABLE IF EXISTS `user_show`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_show` (
  `userid` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `uname` varchar(255) NOT NULL,
  `tel` varchar(255) DEFAULT NULL,
  `wx` varchar(255) DEFAULT NULL,
  `auth` varchar(255) NOT NULL,
  `role` varchar(255) NOT NULL,
  `time` varchar(255) DEFAULT NULL,
  `beizhu` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_show`
--

LOCK TABLES `user_show` WRITE;
/*!40000 ALTER TABLE `user_show` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_show` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yjstatus`
--

DROP TABLE IF EXISTS `yjstatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yjstatus` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `quota` bigint(255) NOT NULL,
  `yj` int(11) NOT NULL,
  `level` varchar(255) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=13 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yjstatus`
--

LOCK TABLES `yjstatus` WRITE;
/*!40000 ALTER TABLE `yjstatus` DISABLE KEYS */;
INSERT INTO `yjstatus` VALUES (1,1000000,50,'会员级'),(2,3000000,60,'超级会员级'),(3,6000000,70,'代理级'),(4,10000000,80,'超级代理级'),(5,20000000,100,'总代理级'),(6,40000000,120,'超级总代理级'),(7,60000000,140,'股东级'),(8,80000000,160,'超级股东级'),(9,100000000,180,'总监级'),(10,0,200,'超级总监级');
/*!40000 ALTER TABLE `yjstatus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yjsxf`
--

DROP TABLE IF EXISTS `yjsxf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yjsxf` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `num` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yjsxf`
--

LOCK TABLES `yjsxf` WRITE;
/*!40000 ALTER TABLE `yjsxf` DISABLE KEYS */;
INSERT INTO `yjsxf` VALUES (1,10);
/*!40000 ALTER TABLE `yjsxf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yjtype`
--

DROP TABLE IF EXISTS `yjtype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yjtype` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` int(255) NOT NULL COMMENT '0 每周结算  1 每月 结算  2 每日结算',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yjtype`
--

LOCK TABLES `yjtype` WRITE;
/*!40000 ALTER TABLE `yjtype` DISABLE KEYS */;
INSERT INTO `yjtype` VALUES (1,2);
/*!40000 ALTER TABLE `yjtype` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yunying`
--

DROP TABLE IF EXISTS `yunying`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yunying` (
  `oid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `level` varchar(255) DEFAULT NULL,
  `auth` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`oid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yunying`
--

LOCK TABLES `yunying` WRITE;
/*!40000 ALTER TABLE `yunying` DISABLE KEYS */;
/*!40000 ALTER TABLE `yunying` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'qp_host'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-01-25 13:51:05
