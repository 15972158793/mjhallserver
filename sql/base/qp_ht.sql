CREATE DATABASE  IF NOT EXISTS `qp_ht` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `qp_ht`;
-- MySQL dump 10.13  Distrib 5.6.17, for osx10.6 (i386)
--
-- Host: 103.53.124.238    Database: qp_ht
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
-- Table structure for table `exchange_log`
--

DROP TABLE IF EXISTS `exchange_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `exchange_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `orderid` bigint(20) NOT NULL,
  `uid` int(11) NOT NULL COMMENT '代理id',
  `gold` bigint(20) NOT NULL COMMENT '兑换金币',
  `score` bigint(20) NOT NULL,
  `time` datetime NOT NULL,
  `goldtime` datetime NOT NULL,
  `ispay` tinyint(2) NOT NULL DEFAULT '0' COMMENT '扣除推广额状态 0为未扣 1为已扣除',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '金币到账状态 0为未到账 1为已到账',
  PRIMARY KEY (`id`),
  KEY `orderid` (`orderid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exchange_log`
--

LOCK TABLES `exchange_log` WRITE;
/*!40000 ALTER TABLE `exchange_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `exchange_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_admin`
--

DROP TABLE IF EXISTS `fa_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_admin` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
  `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(30) NOT NULL DEFAULT '' COMMENT '密码盐',
  `avatar` varchar(100) NOT NULL DEFAULT '' COMMENT '头像',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `loginfailure` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '失败次数',
  `logintime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '登录时间',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `token` varchar(59) NOT NULL DEFAULT '' COMMENT 'Session标识',
  `status` varchar(30) NOT NULL DEFAULT 'normal' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='管理员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_admin`
--

LOCK TABLES `fa_admin` WRITE;
/*!40000 ALTER TABLE `fa_admin` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_admin_log`
--

DROP TABLE IF EXISTS `fa_admin_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_admin_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '管理员ID',
  `username` varchar(30) NOT NULL DEFAULT '' COMMENT '管理员名字',
  `url` varchar(100) NOT NULL DEFAULT '' COMMENT '操作页面',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '日志标题',
  `content` text NOT NULL COMMENT '内容',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP',
  `useragent` varchar(255) NOT NULL DEFAULT '' COMMENT 'User-Agent',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `name` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='管理员日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_admin_log`
--

LOCK TABLES `fa_admin_log` WRITE;
/*!40000 ALTER TABLE `fa_admin_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_admin_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_agent_user`
--

DROP TABLE IF EXISTS `fa_agent_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_agent_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '代理ID',
  `agid` int(11) NOT NULL COMMENT '游戏ID=代理ID',
  `open_id` varchar(50) NOT NULL COMMENT '微信open_id',
  `union_id` varchar(50) NOT NULL COMMENT '微信union_id',
  `top_group` varchar(200) NOT NULL DEFAULT '' COMMENT '上级ID组',
  `score` float(20,2) NOT NULL DEFAULT '0.00' COMMENT '累计获得积分',
  `t_score` int(20) NOT NULL DEFAULT '0' COMMENT '已经提取积分',
  `deepin` int(5) NOT NULL DEFAULT '3' COMMENT '向上返利深度',
  `card` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '房卡',
  `password` varchar(50) NOT NULL COMMENT '提现密码',
  `nickname` blob NOT NULL,
  `head` varchar(255) NOT NULL DEFAULT '',
  `rating` varchar(255) NOT NULL DEFAULT '35,10,5,5,5,5,5,5,5,5' COMMENT '返利比例',
  `level` int(5) NOT NULL DEFAULT '0',
  `add_time` int(50) NOT NULL DEFAULT '0' COMMENT '代理新增时的时间戳',
  `todaygold` float(20,2) NOT NULL DEFAULT '0.00' COMMENT '今日下级玩家创造的活跃额',
  `yestodaygold` float(20,2) NOT NULL DEFAULT '0.00',
  `todaytime` bigint(20) NOT NULL DEFAULT '0',
  `parent` int(11) NOT NULL DEFAULT '0',
  `name` varchar(45) NOT NULL DEFAULT '',
  `alipay` varchar(45) NOT NULL DEFAULT '',
  `aliname` varchar(45) NOT NULL DEFAULT '',
  `bankcard` varchar(45) NOT NULL DEFAULT '',
  `bankname` varchar(45) NOT NULL DEFAULT '',
  `phone` varchar(45) NOT NULL DEFAULT '',
  `allcost` bigint(20) NOT NULL DEFAULT '0',
  `allbills` bigint(20) NOT NULL DEFAULT '0',
  `daybills` bigint(20) NOT NULL DEFAULT '0',
  `weekbills` bigint(20) NOT NULL DEFAULT '0',
  `monthbills` bigint(20) NOT NULL DEFAULT '0',
  `timebills` bigint(20) NOT NULL DEFAULT '0',
  `commission` int(11) NOT NULL DEFAULT '0',
  `t_commission` int(11) NOT NULL DEFAULT '0',
  `bills1` bigint(20) NOT NULL DEFAULT '0',
  `bills2` bigint(20) NOT NULL DEFAULT '0',
  `timecommission` bigint(20) NOT NULL DEFAULT '0',
  `areascale` int(11) NOT NULL DEFAULT '0',
  `areascore` float NOT NULL DEFAULT '0',
  `areatscore` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `index_name` (`agid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='代理表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_agent_user`
--

LOCK TABLES `fa_agent_user` WRITE;
/*!40000 ALTER TABLE `fa_agent_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_agent_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_attachment`
--

DROP TABLE IF EXISTS `fa_attachment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_attachment` (
  `id` int(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '物理路径',
  `imagewidth` varchar(30) NOT NULL DEFAULT '' COMMENT '宽度',
  `imageheight` varchar(30) NOT NULL DEFAULT '' COMMENT '宽度',
  `imagetype` varchar(30) NOT NULL DEFAULT '' COMMENT '图片类型',
  `imageframes` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '图片帧数',
  `filesize` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '文件大小',
  `mimetype` varchar(30) NOT NULL DEFAULT '' COMMENT 'mime类型',
  `extparam` varchar(255) NOT NULL DEFAULT '' COMMENT '透传数据',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建日期',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `uploadtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上传时间',
  `storage` enum('local','upyun','qiniu') NOT NULL DEFAULT 'local' COMMENT '存储位置',
  `sha1` varchar(40) NOT NULL DEFAULT '' COMMENT '文件 sha1编码',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='附件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_attachment`
--

LOCK TABLES `fa_attachment` WRITE;
/*!40000 ALTER TABLE `fa_attachment` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_attachment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_auth_group`
--

DROP TABLE IF EXISTS `fa_auth_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_auth_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父组别',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '组名',
  `rules` text NOT NULL COMMENT '规则ID',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='分组表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_auth_group`
--

LOCK TABLES `fa_auth_group` WRITE;
/*!40000 ALTER TABLE `fa_auth_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_auth_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_auth_group_access`
--

DROP TABLE IF EXISTS `fa_auth_group_access`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_auth_group_access` (
  `uid` int(10) unsigned NOT NULL COMMENT '会员ID',
  `group_id` int(10) unsigned NOT NULL COMMENT '级别ID',
  UNIQUE KEY `uid_group_id` (`uid`,`group_id`),
  KEY `uid` (`uid`),
  KEY `group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限分组表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_auth_group_access`
--

LOCK TABLES `fa_auth_group_access` WRITE;
/*!40000 ALTER TABLE `fa_auth_group_access` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_auth_group_access` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_auth_rule`
--

DROP TABLE IF EXISTS `fa_auth_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_auth_rule` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `type` enum('menu','file') NOT NULL DEFAULT 'file' COMMENT 'menu为菜单,file为权限节点',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父ID',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '规则名称',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '规则名称',
  `icon` varchar(50) NOT NULL DEFAULT '' COMMENT '图标',
  `condition` varchar(255) NOT NULL DEFAULT '' COMMENT '条件',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `ismenu` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否为菜单',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `weigh` int(10) NOT NULL DEFAULT '0' COMMENT '权重',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE,
  KEY `pid` (`pid`),
  KEY `weigh` (`weigh`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='节点表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_auth_rule`
--

LOCK TABLES `fa_auth_rule` WRITE;
/*!40000 ALTER TABLE `fa_auth_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_auth_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_bind_players`
--

DROP TABLE IF EXISTS `fa_bind_players`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_bind_players` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '绑定ID',
  `uid` int(11) NOT NULL COMMENT '游戏ID',
  `agid` int(11) NOT NULL COMMENT '代理ID',
  `bind_time` datetime NOT NULL COMMENT '绑定时间',
  `score1` float(20,2) NOT NULL,
  `score2` float(20,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_name` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='绑定玩家';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_bind_players`
--

LOCK TABLES `fa_bind_players` WRITE;
/*!40000 ALTER TABLE `fa_bind_players` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_bind_players` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_category`
--

DROP TABLE IF EXISTS `fa_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_category` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父ID',
  `type` varchar(30) NOT NULL DEFAULT '' COMMENT '栏目类型',
  `name` varchar(30) NOT NULL DEFAULT '',
  `nickname` varchar(50) NOT NULL DEFAULT '',
  `flag` set('hot','index','recommend') NOT NULL DEFAULT '',
  `image` varchar(100) NOT NULL DEFAULT '' COMMENT '图片',
  `keywords` varchar(255) NOT NULL DEFAULT '' COMMENT '关键字',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `diyname` varchar(30) NOT NULL DEFAULT '' COMMENT '自定义名称',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `weigh` int(10) NOT NULL DEFAULT '0' COMMENT '权重',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`),
  KEY `weigh` (`weigh`,`id`),
  KEY `pid` (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='分类表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_category`
--

LOCK TABLES `fa_category` WRITE;
/*!40000 ALTER TABLE `fa_category` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_changelevel_log`
--

DROP TABLE IF EXISTS `fa_changelevel_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_changelevel_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `operator_id` int(11) DEFAULT NULL COMMENT '操作者ID',
  `agid` int(11) DEFAULT NULL COMMENT '被修改推广员id',
  `change_time` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_changelevel_log`
--

LOCK TABLES `fa_changelevel_log` WRITE;
/*!40000 ALTER TABLE `fa_changelevel_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_changelevel_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_config`
--

DROP TABLE IF EXISTS `fa_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL DEFAULT '' COMMENT '变量名',
  `group` varchar(30) NOT NULL DEFAULT '' COMMENT '分组',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '变量标题',
  `tip` varchar(100) NOT NULL DEFAULT '' COMMENT '变量描述',
  `type` varchar(30) NOT NULL DEFAULT '' COMMENT '类型:string,text,int,bool,array,datetime,date,file',
  `value` text NOT NULL COMMENT '变量值',
  `content` text NOT NULL COMMENT '变量字典数据',
  `rule` varchar(100) NOT NULL DEFAULT '' COMMENT '验证规则',
  `extend` varchar(255) NOT NULL DEFAULT '' COMMENT '扩展属性',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_config`
--

LOCK TABLES `fa_config` WRITE;
/*!40000 ALTER TABLE `fa_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_getgift_log`
--

DROP TABLE IF EXISTS `fa_getgift_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_getgift_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '领取ID',
  `gift_id` int(11) DEFAULT NULL COMMENT '礼包ID',
  `uid` int(11) DEFAULT NULL COMMENT '玩家ID',
  `on_time` datetime DEFAULT NULL COMMENT '领取时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_getgift_log`
--

LOCK TABLES `fa_getgift_log` WRITE;
/*!40000 ALTER TABLE `fa_getgift_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_getgift_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_gift`
--

DROP TABLE IF EXISTS `fa_gift`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_gift` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '礼包ID',
  `giftname` varchar(100) DEFAULT NULL COMMENT '礼包名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_gift`
--

LOCK TABLES `fa_gift` WRITE;
/*!40000 ALTER TABLE `fa_gift` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_gift` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_order`
--

DROP TABLE IF EXISTS `fa_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `goods_name` varchar(20) NOT NULL COMMENT '商品名称',
  `goods_num` varchar(10) NOT NULL DEFAULT '0' COMMENT '商品数量',
  `game_id` varchar(255) NOT NULL DEFAULT '' COMMENT '玩家ID',
  `out_trade_no` varchar(50) NOT NULL DEFAULT '0' COMMENT '订单号',
  `total_fee` varchar(255) NOT NULL DEFAULT '0' COMMENT '订单总价',
  `order_time` datetime NOT NULL COMMENT '下单时间',
  `paytime` datetime NOT NULL COMMENT '支付时间',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '支付状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_order`
--

LOCK TABLES `fa_order` WRITE;
/*!40000 ALTER TABLE `fa_order` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_rating`
--

DROP TABLE IF EXISTS `fa_rating`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_rating` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `score_rating` varchar(255) DEFAULT NULL COMMENT '积分返利比例组',
  `topup_rating` varchar(255) DEFAULT NULL COMMENT '后台充值返利比例组',
  `shopbuy_rating` varchar(255) DEFAULT NULL COMMENT '商城返利比例组',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_rating`
--

LOCK TABLES `fa_rating` WRITE;
/*!40000 ALTER TABLE `fa_rating` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_rating` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_step`
--

DROP TABLE IF EXISTS `fa_step`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_step` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `score_num` bigint(20) DEFAULT '0' COMMENT '积分表查询进度',
  `topup_num` int(50) DEFAULT '0' COMMENT '后台充值表查询进度',
  `shopbuy_num` int(50) DEFAULT '0' COMMENT '商城充值表查询进度',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_step`
--

LOCK TABLES `fa_step` WRITE;
/*!40000 ALTER TABLE `fa_step` DISABLE KEYS */;
INSERT INTO `fa_step` VALUES (1,0,0,0);
/*!40000 ALTER TABLE `fa_step` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_test`
--

DROP TABLE IF EXISTS `fa_test`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_test` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `admin_id` int(10) NOT NULL COMMENT '管理员ID',
  `category_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '分类ID(单选)',
  `category_ids` varchar(100) NOT NULL COMMENT '分类ID(多选)',
  `week` enum('monday','tuesday','wednesday') NOT NULL COMMENT '星期(单选):monday=星期一,tuesday=星期二,wednesday=星期三',
  `flag` set('hot','index','recommend') NOT NULL DEFAULT '' COMMENT '标志(多选):hot=热门,index=首页,recommend=推荐',
  `genderdata` enum('male','female') NOT NULL DEFAULT 'male' COMMENT '性别(单选):male=男,female=女',
  `hobbydata` set('music','reading','swimming') NOT NULL COMMENT '爱好(多选):music=音乐,reading=读书,swimming=游泳',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `image` varchar(100) NOT NULL DEFAULT '' COMMENT '图片',
  `images` varchar(1500) NOT NULL DEFAULT '' COMMENT '图片组',
  `attachfile` varchar(100) NOT NULL DEFAULT '' COMMENT '附件',
  `keywords` varchar(100) NOT NULL DEFAULT '' COMMENT '关键字',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `city` varchar(100) NOT NULL DEFAULT '' COMMENT '省市',
  `price` float(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '价格',
  `views` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '点击',
  `startdate` date DEFAULT NULL COMMENT '开始日期',
  `activitytime` datetime DEFAULT NULL COMMENT '活动时间(datetime)',
  `year` year(4) DEFAULT NULL COMMENT '年',
  `times` time DEFAULT NULL COMMENT '时间',
  `refreshtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '刷新时间(int)',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `weigh` int(10) NOT NULL DEFAULT '0' COMMENT '权重',
  `switch` tinyint(1) NOT NULL DEFAULT '0' COMMENT '开关',
  `status` enum('normal','hidden') NOT NULL DEFAULT 'normal' COMMENT '状态',
  `state` enum('0','1','2') NOT NULL DEFAULT '1' COMMENT '状态值:0=禁用,1=正常,2=推荐',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='测试表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_test`
--

LOCK TABLES `fa_test` WRITE;
/*!40000 ALTER TABLE `fa_test` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_test` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_wechat_autoreply`
--

DROP TABLE IF EXISTS `fa_wechat_autoreply`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_wechat_autoreply` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '标题',
  `text` varchar(100) NOT NULL DEFAULT '' COMMENT '触发文本',
  `eventkey` varchar(50) NOT NULL DEFAULT '' COMMENT '响应事件',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='微信自动回复表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_wechat_autoreply`
--

LOCK TABLES `fa_wechat_autoreply` WRITE;
/*!40000 ALTER TABLE `fa_wechat_autoreply` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_wechat_autoreply` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_wechat_config`
--

DROP TABLE IF EXISTS `fa_wechat_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_wechat_config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '配置名称',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '配置标题',
  `value` text NOT NULL COMMENT '配置值',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='微信配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_wechat_config`
--

LOCK TABLES `fa_wechat_config` WRITE;
/*!40000 ALTER TABLE `fa_wechat_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_wechat_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_wechat_context`
--

DROP TABLE IF EXISTS `fa_wechat_context`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_wechat_context` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `openid` varchar(64) NOT NULL DEFAULT '',
  `type` varchar(30) NOT NULL DEFAULT '' COMMENT '类型',
  `eventkey` varchar(64) NOT NULL DEFAULT '',
  `command` varchar(64) NOT NULL DEFAULT '',
  `message` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
  `refreshtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后刷新时间',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `openid` (`openid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='微信上下文表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_wechat_context`
--

LOCK TABLES `fa_wechat_context` WRITE;
/*!40000 ALTER TABLE `fa_wechat_context` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_wechat_context` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fa_wechat_response`
--

DROP TABLE IF EXISTS `fa_wechat_response`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fa_wechat_response` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '资源名',
  `eventkey` varchar(128) NOT NULL DEFAULT '' COMMENT '事件',
  `type` enum('text','image','news','voice','video','music','link','app') NOT NULL DEFAULT 'text' COMMENT '类型',
  `content` text NOT NULL COMMENT '内容',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `createtime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updatetime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `status` varchar(30) NOT NULL DEFAULT '' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `event` (`eventkey`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='微信资源表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fa_wechat_response`
--

LOCK TABLES `fa_wechat_response` WRITE;
/*!40000 ALTER TABLE `fa_wechat_response` DISABLE KEYS */;
/*!40000 ALTER TABLE `fa_wechat_response` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gamelogs`
--

DROP TABLE IF EXISTS `gamelogs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gamelogs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(10) DEFAULT NULL,
  `gold` int(10) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gamelogs`
--

LOCK TABLES `gamelogs` WRITE;
/*!40000 ALTER TABLE `gamelogs` DISABLE KEYS */;
/*!40000 ALTER TABLE `gamelogs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gold`
--

DROP TABLE IF EXISTS `gold`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gold` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL,
  `uname` blob,
  `money` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `alipay` varchar(255) DEFAULT NULL,
  `bankcard` varchar(255) DEFAULT NULL,
  `bankname` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `applytime` varchar(225) DEFAULT NULL,
  `paytime` varchar(225) DEFAULT NULL,
  `dec` varchar(255) DEFAULT NULL,
  `aliname` varchar(225) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gold`
--

LOCK TABLES `gold` WRITE;
/*!40000 ALTER TABLE `gold` DISABLE KEYS */;
/*!40000 ALTER TABLE `gold` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `parent_cost_log`
--

DROP TABLE IF EXISTS `parent_cost_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parent_cost_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `parent` int(11) NOT NULL,
  `score` decimal(20,4) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `parent` (`parent`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `parent_cost_log`
--

LOCK TABLES `parent_cost_log` WRITE;
/*!40000 ALTER TABLE `parent_cost_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `parent_cost_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `parent_log`
--

DROP TABLE IF EXISTS `parent_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parent_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `gold` int(11) NOT NULL,
  `parent` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `parent` (`parent`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `parent_log`
--

LOCK TABLES `parent_log` WRITE;
/*!40000 ALTER TABLE `parent_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `parent_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `parent_lost_bills`
--

DROP TABLE IF EXISTS `parent_lost_bills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parent_lost_bills` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `gold` int(11) NOT NULL,
  `parent` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `parent_lost_bills`
--

LOCK TABLES `parent_lost_bills` WRITE;
/*!40000 ALTER TABLE `parent_lost_bills` DISABLE KEYS */;
/*!40000 ALTER TABLE `parent_lost_bills` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `parent_win_bills`
--

DROP TABLE IF EXISTS `parent_win_bills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parent_win_bills` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `gold` int(11) NOT NULL,
  `parent` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `parent_win_bills`
--

LOCK TABLES `parent_win_bills` WRITE;
/*!40000 ALTER TABLE `parent_win_bills` DISABLE KEYS */;
/*!40000 ALTER TABLE `parent_win_bills` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quyu`
--

DROP TABLE IF EXISTS `quyu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quyu` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `amount` varchar(255) COLLATE utf8_bin NOT NULL,
  `status` varchar(255) COLLATE utf8_bin NOT NULL,
  `time` varchar(225) COLLATE utf8_bin NOT NULL,
  `gold` varchar(255) COLLATE utf8_bin NOT NULL,
  `orderid` varchar(225) COLLATE utf8_bin NOT NULL,
  `goldtime` varchar(225) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quyu`
--

LOCK TABLES `quyu` WRITE;
/*!40000 ALTER TABLE `quyu` DISABLE KEYS */;
/*!40000 ALTER TABLE `quyu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quyu_log`
--

DROP TABLE IF EXISTS `quyu_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quyu_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `orderid` bigint(20) NOT NULL,
  `uid` int(11) NOT NULL COMMENT '代理id',
  `gold` bigint(20) NOT NULL COMMENT '兑换金币',
  `score` bigint(20) NOT NULL,
  `time` datetime NOT NULL,
  `goldtime` datetime NOT NULL,
  `ispay` tinyint(2) NOT NULL DEFAULT '0' COMMENT '扣除推广额状态 0为未扣 1为已扣除',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '金币到账状态 0为未到账 1为已到账',
  PRIMARY KEY (`id`),
  KEY `orderid` (`orderid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quyu_log`
--

LOCK TABLES `quyu_log` WRITE;
/*!40000 ALTER TABLE `quyu_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `quyu_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `score_log`
--

DROP TABLE IF EXISTS `score_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `score_log` (
  `id` int(20) NOT NULL AUTO_INCREMENT COMMENT '积分日志ID',
  `operator_id` int(20) NOT NULL DEFAULT '0' COMMENT '影响积分玩家ID',
  `agid` int(20) NOT NULL DEFAULT '0' COMMENT '积分变化的代理ID',
  `change_score` float NOT NULL DEFAULT '0' COMMENT '积分变化',
  `save_time` datetime NOT NULL COMMENT '积分变化时间',
  `action` tinyint(4) NOT NULL DEFAULT '0' COMMENT '操作类型(1为积分返利,2为后台充值,3为商城充值,-1为提现)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `score_log`
--

LOCK TABLES `score_log` WRITE;
/*!40000 ALTER TABLE `score_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `score_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `send_card_log`
--

DROP TABLE IF EXISTS `send_card_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `send_card_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `agid` int(11) NOT NULL DEFAULT '0' COMMENT '代理id',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `card_num` int(11) NOT NULL COMMENT '发卡数量',
  `status` int(11) NOT NULL COMMENT '发卡状态',
  `send_time` datetime NOT NULL,
  `method` int(5) DEFAULT NULL COMMENT '发卡方式 1为发给玩家 2为发给代理',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='发卡记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `send_card_log`
--

LOCK TABLES `send_card_log` WRITE;
/*!40000 ALTER TABLE `send_card_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `send_card_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sign`
--

DROP TABLE IF EXISTS `sign`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sign` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `time` varchar(225) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sign`
--

LOCK TABLES `sign` WRITE;
/*!40000 ALTER TABLE `sign` DISABLE KEYS */;
/*!40000 ALTER TABLE `sign` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_score_log`
--

DROP TABLE IF EXISTS `t_score_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_score_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `time` datetime NOT NULL,
  `score` int(11) NOT NULL,
  `order` varchar(128) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_score_log`
--

LOCK TABLES `t_score_log` WRITE;
/*!40000 ALTER TABLE `t_score_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_score_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `yongjin_log`
--

DROP TABLE IF EXISTS `yongjin_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `yongjin_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `orderid` bigint(20) NOT NULL,
  `uid` int(11) NOT NULL COMMENT '代理id',
  `gold` bigint(20) NOT NULL COMMENT '兑换金币',
  `score` bigint(20) NOT NULL,
  `time` datetime NOT NULL,
  `goldtime` datetime NOT NULL,
  `ispay` tinyint(2) NOT NULL DEFAULT '0' COMMENT '扣除推广额状态 0为未扣 1为已扣除',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '金币到账状态 0为未到账 1为已到账',
  PRIMARY KEY (`id`),
  KEY `orderid` (`orderid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `yongjin_log`
--

LOCK TABLES `yongjin_log` WRITE;
/*!40000 ALTER TABLE `yongjin_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `yongjin_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'qp_ht'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-01-25 13:41:26
