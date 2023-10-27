-- MySQL dump 10.13  Distrib 8.0.32, for Win64 (x86_64)
--
-- Host: localhost    Database: finance_db
-- ------------------------------------------------------
-- Server version	8.0.19

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `customers`
--

DROP TABLE IF EXISTS `customers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `customers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nik` varchar(16) NOT NULL,
  `full_name` varchar(100) NOT NULL,
  `legal_name` varchar(100) NOT NULL,
  `birth_place` varchar(100) NOT NULL,
  `birth_date` date NOT NULL DEFAULT '1990-01-20',
  `salary` int NOT NULL DEFAULT '0',
  `photo_customer` varchar(150) NOT NULL,
  `photo_selfie` varchar(150) NOT NULL,
  `created_at` timestamp NOT NULL,
  `created_by` varchar(100) NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `nik_UNIQUE` (`nik`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customers`
--

LOCK TABLES `customers` WRITE;
/*!40000 ALTER TABLE `customers` DISABLE KEYS */;
INSERT INTO `customers` VALUES (1,'1234567890000001','Cristiano Ronaldo','Cristiano Ronaldo','Bandung','1985-06-01',10000000,'https://imgur.com/GALVWWt','https://imgur.com/GALVWWt','2023-10-25 17:00:00','SYSTEM',NULL,NULL),(2,'1234567890000002','Lionel Messi','Lionel Messi','Jakarta','1987-07-01',15000000,'https://imgur.com/Bije5ga','https://imgur.com/Bije5ga','2023-10-25 17:00:00','SYSTEM',NULL,NULL),(4,'1111222233334444','David Beckham','David Joseph William Beckham','LONDON','1975-12-28',750000,'https://imgur.com/GALVWWt','https://imgur.com/GALVWWt','2023-10-26 09:19:11','SYSTEM','2023-10-26 09:38:27','SYSTEM');
/*!40000 ALTER TABLE `customers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tenors`
--

DROP TABLE IF EXISTS `tenors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tenors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `customer_id` int NOT NULL,
  `one_month` int NOT NULL DEFAULT '0',
  `two_month` int NOT NULL DEFAULT '0',
  `three_month` int NOT NULL DEFAULT '0',
  `four_month` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tenors`
--

LOCK TABLES `tenors` WRITE;
/*!40000 ALTER TABLE `tenors` DISABLE KEYS */;
INSERT INTO `tenors` VALUES (1,1,500000,1000000,1500000,2000000),(2,2,1000000,2000000,3000000,4000000),(3,4,500000,0,0,0);
/*!40000 ALTER TABLE `tenors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transactions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `contract_number` varchar(50) NOT NULL,
  `otr` int unsigned NOT NULL DEFAULT '0',
  `admin_fee` int unsigned NOT NULL DEFAULT '0',
  `installment_amount` int unsigned NOT NULL DEFAULT '0',
  `amount_of_interest` decimal(6,2) NOT NULL DEFAULT '0.00',
  `asset_name` varchar(150) NOT NULL,
  `customer_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
INSERT INTO `transactions` VALUES (1,'0123456789',15000000,25000,12,0.10,'Iphone 15 PRO Max',4,'2023-10-27 06:39:17','2023-10-27 06:39:17');
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-10-27 13:43:47
