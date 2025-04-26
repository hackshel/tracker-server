-- MySQL dump 10.13  Distrib 8.0.32, for Linux (x86_64)
--
-- Host: localhost    Database: tk
-- ------------------------------------------------------
-- Server version	8.0.32

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `tk_audit_logs`
--

DROP TABLE IF EXISTS `tk_audit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tk_audit_logs` (
  `log_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT ' ID',
  `user_id` bigint unsigned NOT NULL COMMENT ' ID',
  `operation_type` enum('UPLOAD','DOWNLOAD','DELETE','MODIFY','SHARE','LOGIN') NOT NULL,
  `target_torrent` bigint unsigned DEFAULT NULL COMMENT ' ID',
  `detail` json NOT NULL COMMENT 'IP',
  `timestamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'UTC ',
  `blockchain_hash` varchar(66) DEFAULT NULL,
  PRIMARY KEY (`log_id`),
  KEY `target_torrent` (`target_torrent`),
  CONSTRAINT `tk_audit_logs_ibfk_1` FOREIGN KEY (`target_torrent`) REFERENCES `tk_torrents` (`torrent_id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tk_audit_logs`
--

LOCK TABLES `tk_audit_logs` WRITE;
/*!40000 ALTER TABLE `tk_audit_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `tk_audit_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tk_hospital`
--

DROP TABLE IF EXISTS `tk_hospital`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tk_hospital` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `code` bigint NOT NULL DEFAULT '0',
  `hsp_name` varchar(255) NOT NULL DEFAULT '' COMMENT '医院的名字',
  `grade` int DEFAULT '0' COMMENT 'hospital rank or grade',
  `province_id` int DEFAULT NULL COMMENT '省id',
  `province_name` varchar(300) DEFAULT NULL COMMENT '省名称',
  `city_id` int DEFAULT NULL COMMENT '市id',
  `city_name` varchar(300) DEFAULT NULL COMMENT '市名称',
  `county_id` int DEFAULT NULL COMMENT '区县id',
  `county_name` varchar(300) DEFAULT NULL COMMENT '区县名称',
  `add_time` datetime DEFAULT NULL COMMENT '增加照表时间',
  `last_modify` datetime DEFAULT NULL COMMENT '最终修改状态时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tk_hospital`
--

LOCK TABLES `tk_hospital` WRITE;
/*!40000 ALTER TABLE `tk_hospital` DISABLE KEYS */;
INSERT INTO `tk_hospital` VALUES (8,1100001006641,'北京大学人民医院（北京大学第二临床医学院）',3,0,'',0,'',0,'','2025-04-26 17:39:40','2025-04-26 17:39:40'),(9,1100001006665,'中国医学科学院阜外医院',3,0,'',0,'',0,'','2025-04-26 17:41:43','2025-04-26 17:41:43'),(10,4201021030001,'武汉市中心医院',3,0,'',0,'',0,'','2025-04-26 17:43:23','2025-04-26 17:43:23'),(11,6599001102286,'乌鲁木齐市第一人民医院',3,0,'',0,'',0,'','2025-04-26 17:45:10','2025-04-26 17:45:10'),(12,6104001010007,'西藏民族大学附属医院',2,0,'',0,'',0,'','2025-04-26 17:46:14','2025-04-26 17:46:14');
/*!40000 ALTER TABLE `tk_hospital` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tk_peers`
--

DROP TABLE IF EXISTS `tk_peers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tk_peers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `torrent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'torrent table id ',
  `peer_id` varchar(255) NOT NULL COMMENT 'transmission peer_id hex',
  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'client ip address',
  `port` int unsigned NOT NULL DEFAULT '0' COMMENT 'client port',
  `uploaded` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'upload ',
  `downloaded` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'download ',
  `left_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT ' ',
  `is_seeder` enum('yes','no') NOT NULL DEFAULT 'no',
  `started` datetime DEFAULT NULL,
  `last_action` datetime DEFAULT NULL,
  `prev_action` datetime DEFAULT NULL,
  `connectable` enum('yes','no') NOT NULL DEFAULT 'yes',
  `agent` varchar(60) NOT NULL DEFAULT '' COMMENT 'client agent',
  `finished_at` datetime DEFAULT NULL,
  `downloadoffset` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'download offset',
  `uploadoffset` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'upload offset',
  `dl_speed` varchar(255) NOT NULL DEFAULT '',
  `passkey` varchar(255) NOT NULL DEFAULT '' COMMENT 'passkey ',
  `info_hash` varchar(255) NOT NULL COMMENT 'info sha1 16 ',
  PRIMARY KEY (`id`),
  KEY `torrent` (`torrent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tk_peers`
--

LOCK TABLES `tk_peers` WRITE;
/*!40000 ALTER TABLE `tk_peers` DISABLE KEYS */;
INSERT INTO `tk_peers` VALUES (1,11,'2d7142353035302d61636d7654596e5043457535','192.168.11.202',62131,0,0,0,'yes','2025-04-25 21:07:28','2025-04-26 20:34:03','2025-04-26 20:34:03','yes','','2025-04-26 20:34:03',0,0,'','2be1adfa1f3911f0aeeec3ae841ca01b','4a701fa02b95297da8b909c113194ecfc15fd1eb'),(2,11,'2d7142353034302d337e304b6b52752e38524c74','192.168.11.239',25114,0,1732048795,0,'yes','2025-04-25 21:07:55','2025-04-25 21:10:05','2025-04-25 21:10:05','yes','','2025-04-25 21:10:05',0,0,'','65c62b361f3911f08d495bd447b57075','4a701fa02b95297da8b909c113194ecfc15fd1eb'),(3,13,'2d7142353034302d695a5f42323172742a695130','192.168.11.239',25114,0,1204040441,0,'yes','2025-04-26 20:36:10','2025-04-26 20:37:46','2025-04-26 20:37:46','yes','','2025-04-26 20:37:46',0,0,'','09de96094ff84fe0bffe4100e9ab46a4','b61b99b9419a76fdcc1e044fe4bb0a057c081a1b');
/*!40000 ALTER TABLE `tk_peers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tk_pieces_info`
--

DROP TABLE IF EXISTS `tk_pieces_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tk_pieces_info` (
  `piece_id` varchar(64) NOT NULL COMMENT 'SHA256(piece_data)',
  `torrent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'ID ',
  `storage_nodes` json NOT NULL COMMENT ' IP  3 ',
  `status` enum('PENDING','VERIFIED','INVALID') DEFAULT 'PENDING' COMMENT 'PENDING/VERIFIED/INVALID',
  `encryption_algo` varchar(20) NOT NULL COMMENT 'SM4/AES-256',
  `key_id` varchar(36) NOT NULL COMMENT ' ID',
  PRIMARY KEY (`piece_id`),
  KEY `torrent_id` (`torrent_id`),
  CONSTRAINT `tk_pieces_info_ibfk_1` FOREIGN KEY (`torrent_id`) REFERENCES `tk_torrents` (`torrent_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tk_pieces_info`
--

LOCK TABLES `tk_pieces_info` WRITE;
/*!40000 ALTER TABLE `tk_pieces_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `tk_pieces_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tk_torrents`
--

DROP TABLE IF EXISTS `tk_torrents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tk_torrents` (
  `torrent_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `info_hash` varchar(255) NOT NULL COMMENT 'info sha1 16 ',
  `name` varchar(255) NOT NULL DEFAULT 'dcp name',
  `filename` varchar(500) NOT NULL DEFAULT 'dcp torrent file name',
  `dcp_uuid` varchar(255) NOT NULL COMMENT 'dcp cpl uuid',
  `dcp_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'dcp ',
  `piece_length` int NOT NULL DEFAULT '0',
  `pices_count` int NOT NULL DEFAULT '0' COMMENT 'dcp ',
  `added_time` datetime DEFAULT NULL,
  `dcp_type` enum('single','multi') NOT NULL DEFAULT 'single',
  `numfiles` smallint unsigned NOT NULL DEFAULT '0',
  `tracker_url` varchar(300) NOT NULL COMMENT 'announce',
  `f_sha1` varchar(100) NOT NULL COMMENT 'sha1',
  `seeders` mediumint unsigned NOT NULL DEFAULT '0',
  `leechers` mediumint unsigned NOT NULL DEFAULT '0',
  `times_completed` mediumint unsigned NOT NULL DEFAULT '0',
  `last_action` datetime DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'id',
  PRIMARY KEY (`torrent_id`),
  UNIQUE KEY `info_hash` (`info_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tk_torrents`
--

LOCK TABLES `tk_torrents` WRITE;
/*!40000 ALTER TABLE `tk_torrents` DISABLE KEYS */;
INSERT INTO `tk_torrents` VALUES (7,'1e6b1c2b1948e16604219e5cc098069ba7d93afe','ZhuaWaWa-2D_235_JP_134M_51_PTH_0706_OV.torrent','ZhuaWaWa-2D_235_JP_134M_51_PTH_0706_OV.torrent','UUID:702aa6d4-2cdc-4088-a29a-9d0956004174',154359440228,8388608,18402,'2025-04-22 20:42:46','single',16,'https://tracker.jiujinmax.com/api/v1/tracker/announce','6a9bf6713f43b474ef1070ce738bfdc0c6916662',0,0,0,'2025-04-22 20:42:46',1),(9,'8a19577fb5f690970ca43a57ff1011ae202244b8','ubuntu-25.04-desktop-amd64.iso.torrent','ubuntu-25.04-desktop-amd64.iso.torrent','Ubuntu CD releases.ubuntu.com',0,262144,23951,'2025-04-22 21:14:08','single',0,'https://tracker.jiujinmax.com/api/v1/tracker/announce','7643cd05d8e82823b6e7df7a4e04f268f1d705fe',0,0,0,'2025-04-22 21:14:08',1),(11,'4a701fa02b95297da8b909c113194ecfc15fd1eb','Dune2_TLR-2-4K-48-CNT-LED_S_EN-QMS-EN_51_CINITYLAB_100M_20250306_SMPTE_OV.torrent','Dune2_TLR-2-4K-48-CNT-LED_S_EN-QMS-EN_51_CINITYLAB_100M_20250306_SMPTE_OV.torrent','7777521e-d3d3-4e4a-9765-9b3246cb35b7',1732048795,8388608,207,'2025-04-23 21:12:19','single',6,'https://tracker.jiujinmax.com/api/v1/tracker/announce','2b49f5ca3bf3b5ddfc3df9a65d6b169b39361e82',3,1,2,'2025-04-25 21:10:04',1),(12,'05ee404be87c81d41e462e6aa0eb30718fa3bda3','ZhuiLuoDeShenPan-2D_185_JP_152M_51_YS_0119.torrent','ZhuiLuoDeShenPan-2D_185_JP_152M_51_YS_0119.torrent','a366adf5-b877-4aea-928d-705b91a7dd62',169870503313,8388608,20251,'2025-04-23 22:12:45','single',20,'https://tracker.jiujinmax.com/api/v1/tracker/announce','cd35bd95155c57719702055dcd2435ed0951de81',0,0,0,'2025-04-23 22:12:45',3),(13,'b61b99b9419a76fdcc1e044fe4bb0a057c081a1b','AMR_TLR-2-4K-120-CNT-LED_S_XX-XX_XX_CINITYLAB_100M_20250307_SMPTE_OV.torrent','AMR_TLR-2-4K-120-CNT-LED_S_XX-XX_XX_CINITYLAB_100M_20250307_SMPTE_OV.torrent','ecc67b20-374c-4ec1-bbcb-015842962a5d',1204040441,8388608,144,'2025-04-26 20:33:43','single',6,'http://tracker.cf-noc.work/api/v1/tracker/announce','af551ade401dfefde0bea46638cc6365a79ceae1',1,0,1,'2025-04-26 20:37:46',1);
/*!40000 ALTER TABLE `tk_torrents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tk_users`
--

DROP TABLE IF EXISTS `tk_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tk_users` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_name` varchar(255) NOT NULL DEFAULT '',
  `code` bigint NOT NULL DEFAULT '0',
  `role` enum('admin','doctor','researcher') NOT NULL COMMENT 'admin/doctor/researcher',
  `public_key` text NOT NULL COMMENT 'SM2 ',
  `access_level` int DEFAULT '1' COMMENT '1-55',
  `salt` varchar(50) NOT NULL DEFAULT '' COMMENT 'salt',
  `passwd` varchar(255) NOT NULL DEFAULT '' COMMENT 'hash',
  `last_login` timestamp NULL DEFAULT NULL,
  `passkey` varchar(100) DEFAULT '' COMMENT 'torrent key',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tk_users`
--

LOCK TABLES `tk_users` WRITE;
/*!40000 ALTER TABLE `tk_users` DISABLE KEYS */;
INSERT INTO `tk_users` VALUES (1,'wangxiaochen',1100001006665,'admin','',1,'JDJiJDEyJGxTS3l2M1pVd0RJZ2hxZU1ZdUI4cWU=','JDJiJDEyJGxTS3l2M1pVd0RJZ2hxZU1ZdUI4cWVBdDJYMlRrYllmZzA5cy5LemwweGFrVnpyVGxuZXhlMQ==','2025-04-21 11:13:50','2be1adfa1f3911f0aeeec3ae841ca01b'),(3,'zhaoqiang',1100001006665,'doctor','',1,'JDJiJDEyJG1iZnRCbUxqTUQ3MUN0N0xCSnhBWXU=','JDJiJDEyJG1iZnRCbUxqTUQ3MUN0N0xCSnhBWXV0M1ZHTFBRaDducXlWSUUwQTlGZXNadTRvb1JuLjVL','2025-04-21 11:14:17','65c62b361f3911f08d495bd447b57075'),(10,'liutao',6599001102286,'doctor','',1,'JDJhJDEwJFgveklLMXpHMDJxVWkzRWg0SThrRmU=','JDJhJDEwJFgveklLMXpHMDJxVWkzRWg0SThrRmVXLlBqcjdYalJSbXBjOTZXV2ZUMGU2WlU2UllSYWll','2025-04-26 07:54:25','a01284f02fa847949b81102c4c595019'),(11,'dainan',4201021030001,'doctor','',1,'JDJhJDEwJFJPYWtaUVlFRWY4RldkcTlCY0Ridk8=','JDJhJDEwJFJPYWtaUVlFRWY4RldkcTlCY0Ridk9PTEo1clZROVZUZHRhdmZLb0tDMXNONjhoZ0ZqWlpT','2025-04-26 12:26:47','09de96094ff84fe0bffe4100e9ab46a4'),(12,'zhaoqiang',1100001006665,'doctor','',1,'JDJhJDEwJFBvMEpQeGlpckZQOTQwOUNUOC5SSS4=','JDJhJDEwJFBvMEpQeGlpckZQOTQwOUNUOC5SSS5JRUg0ZkFOdFRUZWpOQkJqT1poRHZkM3VvSTRXVllD','2025-04-26 12:29:53','5807ae353dbe41e1a9a12a4c86532a3a');
/*!40000 ALTER TABLE `tk_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `user_id` bigint NOT NULL AUTO_INCREMENT,
  `user_name` varchar(191) DEFAULT NULL,
  `code` bigint DEFAULT NULL,
  `role` bigint DEFAULT NULL,
  `public_key` longtext,
  `access_level` bigint DEFAULT NULL,
  `salt` longtext,
  `password` longtext,
  `last_login` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uni_users_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-04-26 20:43:04
