CREATE DATABASE  IF NOT EXISTS `qp` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `qp`;
-- MySQL dump 10.13  Distrib 5.6.17, for osx10.6 (i386)
--
-- Host: 103.53.124.238    Database: qp
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
-- Table structure for table `account`
--

DROP TABLE IF EXISTS `account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `openid` text NOT NULL,
  `wyid` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=154245 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `account`
--

LOCK TABLES `account` WRITE;
/*!40000 ALTER TABLE `account` DISABLE KEYS */;
/*!40000 ALTER TABLE `account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `alms`
--

DROP TABLE IF EXISTS `alms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alms` (
  `id` bigint(20) NOT NULL,
  `num` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `alms`
--

LOCK TABLES `alms` WRITE;
/*!40000 ALTER TABLE `alms` DISABLE KEYS */;
/*!40000 ALTER TABLE `alms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `black`
--

DROP TABLE IF EXISTS `black`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `black` (
  `id` bigint(20) NOT NULL,
  `createtime` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `black`
--

LOCK TABLES `black` WRITE;
/*!40000 ALTER TABLE `black` DISABLE KEYS */;
/*!40000 ALTER TABLE `black` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `club`
--

DROP TABLE IF EXISTS `club`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `club` (
  `uid` bigint(20) NOT NULL,
  `info` blob NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `club`
--

LOCK TABLES `club` WRITE;
/*!40000 ALTER TABLE `club` DISABLE KEYS */;
/*!40000 ALTER TABLE `club` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `clubmgr`
--

DROP TABLE IF EXISTS `clubmgr`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `clubmgr` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `info` blob NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10001 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clubmgr`
--

LOCK TABLES `clubmgr` WRITE;
/*!40000 ALTER TABLE `clubmgr` DISABLE KEYS */;
/*!40000 ALTER TABLE `clubmgr` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `dial`
--

DROP TABLE IF EXISTS `dial`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dial` (
  `uid` bigint(20) NOT NULL,
  `num` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dial`
--

LOCK TABLES `dial` WRITE;
/*!40000 ALTER TABLE `dial` DISABLE KEYS */;
/*!40000 ALTER TABLE `dial` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `exchange`
--

DROP TABLE IF EXISTS `exchange`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `exchange` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `name` blob NOT NULL,
  `gold` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  `state` int(11) NOT NULL,
  `dec` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `exchange`
--

LOCK TABLES `exchange` WRITE;
/*!40000 ALTER TABLE `exchange` DISABLE KEYS */;
/*!40000 ALTER TABLE `exchange` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invite`
--

DROP TABLE IF EXISTS `invite`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invite` (
  `uid` bigint(20) NOT NULL,
  `value` blob NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invite`
--

LOCK TABLES `invite` WRITE;
/*!40000 ALTER TABLE `invite` DISABLE KEYS */;
/*!40000 ALTER TABLE `invite` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_agent`
--

DROP TABLE IF EXISTS `log_agent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_agent` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `agent` int(11) NOT NULL,
  `gametype` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_agent`
--

LOCK TABLES `log_agent` WRITE;
/*!40000 ALTER TABLE `log_agent` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_agent` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_agentbills`
--

DROP TABLE IF EXISTS `log_agentbills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_agentbills` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `gold` int(11) NOT NULL,
  `gametype` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_agentbills`
--

LOCK TABLES `log_agentbills` WRITE;
/*!40000 ALTER TABLE `log_agentbills` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_agentbills` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_agentgold`
--

DROP TABLE IF EXISTS `log_agentgold`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_agentgold` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `gold` int(11) NOT NULL,
  `gametype` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_agentgold`
--

LOCK TABLES `log_agentgold` WRITE;
/*!40000 ALTER TABLE `log_agentgold` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_agentgold` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_base`
--

DROP TABLE IF EXISTS `log_base`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_base` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_base`
--

LOCK TABLES `log_base` WRITE;
/*!40000 ALTER TABLE `log_base` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_base` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_bills`
--

DROP TABLE IF EXISTS `log_bills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_bills` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `gametype` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_bills`
--

LOCK TABLES `log_bills` WRITE;
/*!40000 ALTER TABLE `log_bills` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_bills` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_bzw`
--

DROP TABLE IF EXISTS `log_bzw`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_bzw` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `gold` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  `gametype` int(11) NOT NULL DEFAULT '40000',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_bzw`
--

LOCK TABLES `log_bzw` WRITE;
/*!40000 ALTER TABLE `log_bzw` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_bzw` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_client`
--

DROP TABLE IF EXISTS `log_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_client` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `gametype` int(11) NOT NULL,
  `next` text NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_client`
--

LOCK TABLES `log_client` WRITE;
/*!40000 ALTER TABLE `log_client` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_client` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_ddz`
--

DROP TABLE IF EXISTS `log_ddz`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_ddz` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_ddz`
--

LOCK TABLES `log_ddz` WRITE;
/*!40000 ALTER TABLE `log_ddz` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_ddz` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_eathot`
--

DROP TABLE IF EXISTS `log_eathot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_eathot` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_eathot`
--

LOCK TABLES `log_eathot` WRITE;
/*!40000 ALTER TABLE `log_eathot` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_eathot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_gc`
--

DROP TABLE IF EXISTS `log_gc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_gc` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_gc`
--

LOCK TABLES `log_gc` WRITE;
/*!40000 ALTER TABLE `log_gc` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_gc` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_gold`
--

DROP TABLE IF EXISTS `log_gold`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_gold` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `dec` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `index_uid` (`uid`),
  KEY `index_dec` (`dec`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_gold`
--

LOCK TABLES `log_gold` WRITE;
/*!40000 ALTER TABLE `log_gold` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_gold` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_kwx`
--

DROP TABLE IF EXISTS `log_kwx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_kwx` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_kwx`
--

LOCK TABLES `log_kwx` WRITE;
/*!40000 ALTER TABLE `log_kwx` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_kwx` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_niuniu`
--

DROP TABLE IF EXISTS `log_niuniu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_niuniu` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_niuniu`
--

LOCK TABLES `log_niuniu` WRITE;
/*!40000 ALTER TABLE `log_niuniu` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_niuniu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_niuniujx`
--

DROP TABLE IF EXISTS `log_niuniujx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_niuniujx` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_niuniujx`
--

LOCK TABLES `log_niuniujx` WRITE;
/*!40000 ALTER TABLE `log_niuniujx` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_niuniujx` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_ptj`
--

DROP TABLE IF EXISTS `log_ptj`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_ptj` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_ptj`
--

LOCK TABLES `log_ptj` WRITE;
/*!40000 ALTER TABLE `log_ptj` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_ptj` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_room`
--

DROP TABLE IF EXISTS `log_room`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_room` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `p1` bigint(20) NOT NULL,
  `p2` bigint(20) NOT NULL,
  `p3` bigint(20) NOT NULL,
  `p4` bigint(20) NOT NULL,
  `p5` bigint(20) NOT NULL,
  `p6` bigint(20) NOT NULL,
  `ip1` varchar(45) NOT NULL,
  `ip2` varchar(45) NOT NULL,
  `ip3` varchar(45) NOT NULL,
  `ip4` varchar(45) NOT NULL,
  `ip5` varchar(45) NOT NULL,
  `ip6` varchar(45) NOT NULL,
  `win1` int(11) NOT NULL,
  `win2` int(11) NOT NULL,
  `win3` int(11) NOT NULL,
  `win4` int(11) NOT NULL,
  `win5` int(11) NOT NULL,
  `win6` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_room`
--

LOCK TABLES `log_room` WRITE;
/*!40000 ALTER TABLE `log_room` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_room` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_run`
--

DROP TABLE IF EXISTS `log_run`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_run` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_run`
--

LOCK TABLES `log_run` WRITE;
/*!40000 ALTER TABLE `log_run` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_run` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_saolei`
--

DROP TABLE IF EXISTS `log_saolei`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_saolei` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_saolei`
--

LOCK TABLES `log_saolei` WRITE;
/*!40000 ALTER TABLE `log_saolei` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_saolei` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_score`
--

DROP TABLE IF EXISTS `log_score`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_score` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `name` blob NOT NULL,
  `head` varchar(128) NOT NULL,
  `gameid` int(11) NOT NULL,
  `room` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  `score` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_score`
--

LOCK TABLES `log_score` WRITE;
/*!40000 ALTER TABLE `log_score` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_score` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_tenhalf`
--

DROP TABLE IF EXISTS `log_tenhalf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_tenhalf` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_tenhalf`
--

LOCK TABLES `log_tenhalf` WRITE;
/*!40000 ALTER TABLE `log_tenhalf` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_tenhalf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_ttz`
--

DROP TABLE IF EXISTS `log_ttz`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_ttz` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_ttz`
--

LOCK TABLES `log_ttz` WRITE;
/*!40000 ALTER TABLE `log_ttz` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_ttz` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_wzq`
--

DROP TABLE IF EXISTS `log_wzq`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_wzq` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid1` bigint(20) NOT NULL,
  `uid2` bigint(20) NOT NULL,
  `gold` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_wzq`
--

LOCK TABLES `log_wzq` WRITE;
/*!40000 ALTER TABLE `log_wzq` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_wzq` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_xyzjh`
--

DROP TABLE IF EXISTS `log_xyzjh`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_xyzjh` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_xyzjh`
--

LOCK TABLES `log_xyzjh` WRITE;
/*!40000 ALTER TABLE `log_xyzjh` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_xyzjh` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `log_zjh`
--

DROP TABLE IF EXISTS `log_zjh`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `log_zjh` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL,
  `num` int(11) NOT NULL,
  `info` text NOT NULL,
  `creation_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `log_zjh`
--

LOCK TABLES `log_zjh` WRITE;
/*!40000 ALTER TABLE `log_zjh` DISABLE KEYS */;
/*!40000 ALTER TABLE `log_zjh` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notice`
--

DROP TABLE IF EXISTS `notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notice` (
  `id` int(11) NOT NULL,
  `title` text NOT NULL,
  `date` text NOT NULL,
  `context` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notice`
--

LOCK TABLES `notice` WRITE;
/*!40000 ALTER TABLE `notice` DISABLE KEYS */;
/*!40000 ALTER TABLE `notice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `noticekwx`
--

DROP TABLE IF EXISTS `noticekwx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `noticekwx` (
  `id` int(11) NOT NULL,
  `title` text NOT NULL,
  `date` text NOT NULL,
  `context` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `noticekwx`
--

LOCK TABLES `noticekwx` WRITE;
/*!40000 ALTER TABLE `noticekwx` DISABLE KEYS */;
/*!40000 ALTER TABLE `noticekwx` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `oldkwx`
--

DROP TABLE IF EXISTS `oldkwx`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `oldkwx` (
  `openid` varchar(128) NOT NULL,
  `uid` varchar(128) NOT NULL,
  `bit` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `oldkwx`
--

LOCK TABLES `oldkwx` WRITE;
/*!40000 ALTER TABLE `oldkwx` DISABLE KEYS */;
/*!40000 ALTER TABLE `oldkwx` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `report`
--

DROP TABLE IF EXISTS `report`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL,
  `rid` bigint(20) NOT NULL,
  `type` int(11) NOT NULL COMMENT '0 斗',
  `dec` text NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `report`
--

LOCK TABLES `report` WRITE;
/*!40000 ALTER TABLE `report` DISABLE KEYS */;
/*!40000 ALTER TABLE `report` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `robot`
--

DROP TABLE IF EXISTS `robot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `robot` (
  `id` int(11) NOT NULL,
  `name` varchar(45) NOT NULL,
  `head` varchar(255) NOT NULL,
  `maxmoney` int(11) NOT NULL,
  `minmoney` int(11) NOT NULL,
  `interval` bigint(20) NOT NULL,
  `mode` int(11) NOT NULL,
  `sign` varchar(255) NOT NULL,
  `sex` int(11) NOT NULL,
  `ip` varchar(45) NOT NULL,
  `address` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `robot`
--

LOCK TABLES `robot` WRITE;
/*!40000 ALTER TABLE `robot` DISABLE KEYS */;
/*!40000 ALTER TABLE `robot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `share`
--

DROP TABLE IF EXISTS `share`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `share` (
  `uid` bigint(20) NOT NULL,
  `num` int(11) NOT NULL,
  `get` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `share`
--

LOCK TABLES `share` WRITE;
/*!40000 ALTER TABLE `share` DISABLE KEYS */;
/*!40000 ALTER TABLE `share` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sign`
--

DROP TABLE IF EXISTS `sign`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sign` (
  `uid` bigint(20) NOT NULL,
  `index` int(11) NOT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sign`
--

LOCK TABLES `sign` WRITE;
/*!40000 ALTER TABLE `sign` DISABLE KEYS */;
/*!40000 ALTER TABLE `sign` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `signconfig`
--

DROP TABLE IF EXISTS `signconfig`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `signconfig` (
  `id` int(11) NOT NULL,
  `icon` int(11) NOT NULL,
  `money` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `signconfig`
--

LOCK TABLES `signconfig` WRITE;
/*!40000 ALTER TABLE `signconfig` DISABLE KEYS */;
INSERT INTO `signconfig` VALUES (1,1,100),(2,1,200),(3,2,300),(4,2,400),(5,3,600),(6,3,800),(7,4,1000);
/*!40000 ALTER TABLE `signconfig` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `topbet`
--

DROP TABLE IF EXISTS `topbet`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `topbet` (
  `id` int(11) NOT NULL,
  `info` blob NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `topbet`
--

LOCK TABLES `topbet` WRITE;
/*!40000 ALTER TABLE `topbet` DISABLE KEYS */;
/*!40000 ALTER TABLE `topbet` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL,
  `value` blob NOT NULL,
  `createtime` bigint(20) NOT NULL,
  `gold` int(11) NOT NULL,
  `savegold` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userbase`
--

DROP TABLE IF EXISTS `userbase`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `userbase` (
  `uid` bigint(20) NOT NULL,
  `money` int(11) NOT NULL,
  `gem` int(11) NOT NULL,
  `charm` int(11) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userbase`
--

LOCK TABLES `userbase` WRITE;
/*!40000 ALTER TABLE `userbase` DISABLE KEYS */;
/*!40000 ALTER TABLE `userbase` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'qp'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-01-25 13:36:52
