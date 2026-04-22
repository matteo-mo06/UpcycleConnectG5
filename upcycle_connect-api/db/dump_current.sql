-- MySQL dump 10.13  Distrib 8.0.43, for Win64 (x86_64)
--
-- Host: localhost    Database: upcycle
-- ------------------------------------------------------
-- Server version	8.0.43

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
-- Table structure for table `advice`
--

DROP TABLE IF EXISTS `advice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `advice` (
  `id_advice` char(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `tag` varchar(120) DEFAULT NULL,
  `date_advice` date DEFAULT NULL,
  PRIMARY KEY (`id_advice`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `advice`
--

LOCK TABLES `advice` WRITE;
/*!40000 ALTER TABLE `advice` DISABLE KEYS */;
INSERT INTO `advice` VALUES ('30000003-0000-4000-8000-000000000000','Transformer un jean usé en sac à main','couture','2026-01-15'),('30000004-0000-4000-8000-000000000000','Réparer une chaise avec des chutes de bois','menuiserie','2026-02-03'),('30000005-0000-4000-8000-000000000000','Customiser un luminaire avec des bouteilles','deco','2026-02-20'),('30000006-0000-4000-8000-000000000000','Créer un jardin vertical avec des palettes','jardinage','2026-03-05'),('30000007-0000-4000-8000-000000000000','Peindre ses meubles sans ponçage préalable','peinture','2026-03-18'),('30000008-0000-4000-8000-000000000000','Fabriquer un bac à compost maison en bois','jardinage','2026-04-01'),('30000009-0000-4000-8000-000000000000','Customiser de vieilles baskets avec peinture textile','mode','2026-04-10'),('30000010-0000-4000-8000-000000000000','Réparer une fermeture éclair glissière','couture','2026-04-15'),('TEST-ADV-0001','Test advice title 1','TEST_TAG','2026-01-01'),('TEST-ADV-0002','Test advice title 2','TEST_TAG','2026-02-01');
/*!40000 ALTER TABLE `advice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `announcement`
--

DROP TABLE IF EXISTS `announcement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `announcement` (
  `id_announcement` char(36) NOT NULL,
  `address_annoucement` varchar(255) DEFAULT NULL,
  `city` varchar(120) DEFAULT NULL,
  `postal` varchar(30) DEFAULT NULL,
  `description_annoucement` text,
  `availability_date` date NOT NULL,
  `price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `request` tinyint(1) NOT NULL DEFAULT '0',
  `state_annoucement` varchar(60) DEFAULT NULL,
  `id_category` char(36) DEFAULT NULL,
  `title_announcement` varchar(255) NOT NULL,
  `type_announcement` varchar(20) NOT NULL DEFAULT 'don',
  `condition_announcement` varchar(20) DEFAULT NULL,
  `id_buyer` char(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_announcement`),
  KEY `fk_announcement_category` (`id_category`),
  CONSTRAINT `fk_announcement_category` FOREIGN KEY (`id_category`) REFERENCES `category` (`id_category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `announcement`
--

LOCK TABLES `announcement` WRITE;
/*!40000 ALTER TABLE `announcement` DISABLE KEYS */;
INSERT INTO `announcement` VALUES ('c0000003-0000-4000-8000-000000000000','14 rue des Lilas','Paris','75020','Tissus coton de différentes couleurs, idéaux pour couture créative et patchwork.','2026-05-01',0.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000001','Lot de tissus coton à donner','don',NULL,NULL,NULL),('c0000004-0000-4000-8000-000000000000','3 impasse du Moulin','Bordeaux','33000','Dizaine de palettes EPAL récupérées sur chantier, propres et sans clou apparent.','2026-05-15',20.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000002','Palettes bois en bon état','don',NULL,NULL,NULL),('c0000005-0000-4000-8000-000000000000','27 avenue Gambetta','Lyon','69003','Cinq PC tours hors service, idéaux pour récupérer composants électroniques.','2026-06-01',5.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000003','Lot PC tours à démanteler','don',NULL,NULL,NULL),('c0000006-0000-4000-8000-000000000000','8 rue de la Paix','Lille','59000','Sacs en cuir véritable usagés, parfaits pour projets de maroquinerie upcyclée.','2026-06-10',15.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000004','Sacs cuir à transformer','don',NULL,NULL,NULL),('c0000007-0000-4000-8000-000000000000','55 boulevard Libération','Marseille','13001','Service à vaisselle dépareillé en céramique artisanale, plusieurs pièces disponibles.','2026-05-20',10.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000005','Service céramique dépareillé','don',NULL,NULL,NULL),('c0000008-0000-4000-8000-000000000000','2 allée des Roses','Nantes','44000','Composteur bois 300L légèrement fissuré, fonctionnel, à récupérer sur place.','2026-07-01',0.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000006','Composteur de jardin 300L','don',NULL,NULL,NULL),('c0000009-0000-4000-8000-000000000000','19 rue Oberkampf','Paris','75011','Table et 4 chaises restaurées, peinture chalk paint style industriel, très bon état.','2026-06-25',80.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000007','Table et chaises relookées','don',NULL,NULL,NULL),('c0000010-0000-4000-8000-000000000000','7 zone industrielle Nord','Strasbourg','67100','Chutes de métal variées acier et aluminium, idéales pour sculpture et soudure créative.','2026-07-15',5.00,0,'Active','aaaaaaaa-aaaa-aaaa-aaaa-000000000008','Chutes acier et aluminium','don',NULL,NULL,NULL),('TEST-ANN-0001','1 Test Street','TestCity','00001','Test announcement description','2026-01-01',100.00,0,'ACTIVE','aaaaaaaa-aaaa-aaaa-aaaa-000000000001','Test Announcement One','don',NULL,NULL,NULL),('TEST-ANN-0002','2 Test Avenue','TestTown','00002','Second test announcement','2026-02-01',200.00,1,'ACTIVE','aaaaaaaa-aaaa-aaaa-aaaa-000000000002','Test Announcement Two','don',NULL,NULL,NULL);
/*!40000 ALTER TABLE `announcement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `announcement_locker`
--

DROP TABLE IF EXISTS `announcement_locker`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `announcement_locker` (
  `id_announcement` char(36) NOT NULL,
  `id_locker` char(36) NOT NULL,
  PRIMARY KEY (`id_announcement`),
  UNIQUE KEY `id_locker` (`id_locker`),
  CONSTRAINT `fk_announcement_locker_announcement` FOREIGN KEY (`id_announcement`) REFERENCES `announcement` (`id_announcement`),
  CONSTRAINT `fk_announcement_locker_locker` FOREIGN KEY (`id_locker`) REFERENCES `locker` (`id_locker`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `announcement_locker`
--

LOCK TABLES `announcement_locker` WRITE;
/*!40000 ALTER TABLE `announcement_locker` DISABLE KEYS */;
INSERT INTO `announcement_locker` VALUES ('c0000003-0000-4000-8000-000000000000','d0000003-0000-4000-8000-000000000000'),('c0000004-0000-4000-8000-000000000000','d0000004-0000-4000-8000-000000000000'),('c0000005-0000-4000-8000-000000000000','d0000005-0000-4000-8000-000000000000'),('c0000006-0000-4000-8000-000000000000','d0000006-0000-4000-8000-000000000000'),('c0000007-0000-4000-8000-000000000000','d0000007-0000-4000-8000-000000000000'),('c0000008-0000-4000-8000-000000000000','d0000008-0000-4000-8000-000000000000'),('c0000009-0000-4000-8000-000000000000','d0000009-0000-4000-8000-000000000000'),('c0000010-0000-4000-8000-000000000000','d0000010-0000-4000-8000-000000000000'),('TEST-ANN-0001','TEST-LOCK-0001'),('TEST-ANN-0002','TEST-LOCK-0002');
/*!40000 ALTER TABLE `announcement_locker` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `category`
--

DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `category` (
  `id_category` char(36) NOT NULL,
  `name_category` varchar(100) NOT NULL,
  `description_category` text,
  PRIMARY KEY (`id_category`),
  UNIQUE KEY `name_category` (`name_category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category`
--

LOCK TABLES `category` WRITE;
/*!40000 ALTER TABLE `category` DISABLE KEYS */;
INSERT INTO `category` VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-000000000001','Couture & Textile','Retouches, créations et recyclage textile'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000002','Menuiserie & Bois','Réparation et création de meubles en bois'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000003','Électronique','Réparation d\'appareils électroniques'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000004','Maroquinerie','Réparation et création en cuir'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000005','Céramique & Poterie','Création et restauration céramique'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000006','Jardinage & Compostage','Compost, récupération déchets verts, jardinage écologique'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000007','Peinture & Décoration','Remise en peinture et relooking de mobilier usagé'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000008','Récupération Métaux','Travail des métaux et matériaux de récupération industrielle'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000009','Papeterie & Carton','Recyclage papier, carton et créations associées'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000010','Mode & Couture Upcyclée','Transformation de vêtements et accessoires usagés');
/*!40000 ALTER TABLE `category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `document`
--

DROP TABLE IF EXISTS `document`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `document` (
  `id_document` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `category` varchar(120) DEFAULT NULL,
  `link` varchar(2048) DEFAULT NULL,
  PRIMARY KEY (`id_document`),
  KEY `fk_document_user` (`id_user`),
  CONSTRAINT `fk_document_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `document`
--

LOCK TABLES `document` WRITE;
/*!40000 ALTER TABLE `document` DISABLE KEYS */;
INSERT INTO `document` VALUES ('TEST-DOC-0001','TEST-USER-0001','TEST_CATEGORY','https://example.com/test_doc1'),('TEST-DOC-0002','TEST-USER-0002','TEST_CATEGORY','https://example.com/test_doc2');
/*!40000 ALTER TABLE `document` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `event`
--

DROP TABLE IF EXISTS `event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `event` (
  `id_event` char(36) NOT NULL,
  `title_event` varchar(255) NOT NULL DEFAULT '',
  `description_event` text,
  `date_event` datetime DEFAULT NULL,
  `location_event` varchar(255) DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `price_cents` int NOT NULL DEFAULT '0',
  `id_creator` char(36) DEFAULT NULL,
  PRIMARY KEY (`id_event`),
  KEY `fk_event_creator` (`id_creator`),
  CONSTRAINT `fk_event_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `event`
--

LOCK TABLES `event` WRITE;
/*!40000 ALTER TABLE `event` DISABLE KEYS */;
INSERT INTO `event` VALUES ('e0000003-0000-4000-8000-000000000000','Atelier couture zéro déchet','Transformez vos vieux vêtements en pièces uniques grâce aux techniques de upcycling textile.','2026-05-10 10:00:00','Paris 11e, 14 rue de la Roquette',15,1000,'TEST-USER-0001'),('e0000004-0000-4000-8000-000000000000','Conférence upcycling industriel','Les nouvelles pratiques de récupération dans l\'industrie manufacturière française.','2026-05-18 14:00:00','Lyon 2e, 5 place Bellecour',80,0,'TEST-USER-0001'),('e0000005-0000-4000-8000-000000000000','Repair Café Bordeaux','Apportez vos objets cassés, repartez avec une seconde vie grâce à nos bénévoles réparateurs.','2026-06-01 09:00:00','Bordeaux Victoire, 3 cours Victor Hugo',30,0,'TEST-USER-0002'),('e0000006-0000-4000-8000-000000000000','Marché de la récup\'','Brocante spécialisée dans les matériaux de récupération : bois, métal, textiles, électronique.','2026-06-14 08:00:00','Nantes Bouffay, place du Bouffay',200,0,'TEST-USER-0002'),('e0000007-0000-4000-8000-000000000000','Formation menuiserie recyclée','Créez vos propres meubles à partir de palettes et bois de récupération, niveau débutant.','2026-07-05 09:30:00','Toulouse Compans, allée du Brienne',12,3000,'TEST-USER-0001'),('e0000008-0000-4000-8000-000000000000','Atelier poterie récupération','Travaillez l\'argile et redonnez vie aux céramiques brisées par kintsugi et assemblage.','2026-07-19 14:00:00','Strasbourg Krutenau, 8 quai Mullenheim',10,1500,'TEST-USER-0002'),('e0000009-0000-4000-8000-000000000000','Hackathon éco-conception','Prototypage en 48h de produits innovants à base de matériaux recyclés et de chutes industrielles.','2026-08-22 09:00:00','Paris La Défense, Grande Arche',50,500,'TEST-USER-0001'),('e0000010-0000-4000-8000-000000000000','Exposition upcycling artistique','Vernissage et vente d\'œuvres plastiques réalisées intégralement à partir de déchets ménagers.','2026-09-06 17:00:00','Montpellier Ecusson, 2 rue de la Loge',100,0,'TEST-USER-0002'),('TEST-EVENT-0001','',NULL,NULL,NULL,NULL,0,NULL),('TEST-EVENT-0002','',NULL,NULL,NULL,NULL,0,NULL);
/*!40000 ALTER TABLE `event` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `formation`
--

DROP TABLE IF EXISTS `formation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `formation` (
  `id_formation` char(36) NOT NULL,
  `id_creator` char(36) DEFAULT NULL,
  `id_formateur` char(36) DEFAULT NULL,
  PRIMARY KEY (`id_formation`),
  KEY `fk_formation_creator` (`id_creator`),
  KEY `fk_formation_formateur` (`id_formateur`),
  CONSTRAINT `fk_formation_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_formation_formateur` FOREIGN KEY (`id_formateur`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `formation`
--

LOCK TABLES `formation` WRITE;
/*!40000 ALTER TABLE `formation` DISABLE KEYS */;
INSERT INTO `formation` VALUES ('b0000003-0000-4000-8000-000000000000','TEST-USER-0001','TEST-USER-0002'),('b0000004-0000-4000-8000-000000000000','TEST-USER-0001','TEST-USER-0001'),('b0000005-0000-4000-8000-000000000000','TEST-USER-0002','TEST-USER-0001'),('b0000006-0000-4000-8000-000000000000','TEST-USER-0002','TEST-USER-0002'),('b0000007-0000-4000-8000-000000000000','TEST-USER-0001','TEST-USER-0002'),('b0000008-0000-4000-8000-000000000000','TEST-USER-0001','TEST-USER-0001'),('b0000009-0000-4000-8000-000000000000','TEST-USER-0002','TEST-USER-0002'),('b0000010-0000-4000-8000-000000000000','TEST-USER-0002','TEST-USER-0001'),('TEST-FORM-0001',NULL,NULL),('TEST-FORM-0002',NULL,NULL);
/*!40000 ALTER TABLE `formation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `locker`
--

DROP TABLE IF EXISTS `locker`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `locker` (
  `id_locker` char(36) NOT NULL,
  `access_code` varchar(120) DEFAULT NULL,
  `locker_number` int DEFAULT NULL,
  PRIMARY KEY (`id_locker`),
  UNIQUE KEY `locker_number_unique` (`locker_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locker`
--

LOCK TABLES `locker` WRITE;
/*!40000 ALTER TABLE `locker` DISABLE KEYS */;
INSERT INTO `locker` VALUES ('0c230337-3e84-11f1-ab91-005056c00001',NULL,6),('0c230546-3e84-11f1-ab91-005056c00001',NULL,7),('0c230599-3e84-11f1-ab91-005056c00001',NULL,8),('0c2305d8-3e84-11f1-ab91-005056c00001',NULL,9),('0c23060d-3e84-11f1-ab91-005056c00001',NULL,10),('0c230723-3e84-11f1-ab91-005056c00001',NULL,11),('0c23075c-3e84-11f1-ab91-005056c00001',NULL,12),('0c230794-3e84-11f1-ab91-005056c00001',NULL,13),('0c2307c1-3e84-11f1-ab91-005056c00001',NULL,14),('0c2307ee-3e84-11f1-ab91-005056c00001',NULL,15),('0c23081d-3e84-11f1-ab91-005056c00001',NULL,16),('0c23084e-3e84-11f1-ab91-005056c00001',NULL,17),('0c23087b-3e84-11f1-ab91-005056c00001',NULL,18),('0c2308aa-3e84-11f1-ab91-005056c00001',NULL,19),('0c2308d7-3e84-11f1-ab91-005056c00001',NULL,20),('4fade1b3-d6b6-4635-b505-0820f31b41f2',NULL,3),('59adc1b0-5760-425f-81b7-3e61755cf3b0',NULL,5),('90c96690-e4ff-4711-b900-89be2fe9dcdb',NULL,2),('c46811ae-8fda-4b2c-9757-fc4d4e56ef08',NULL,1),('d0000003-0000-4000-8000-000000000000','4821',NULL),('d0000004-0000-4000-8000-000000000000','7356',NULL),('d0000005-0000-4000-8000-000000000000','2947',NULL),('d0000006-0000-4000-8000-000000000000','6183',NULL),('d0000007-0000-4000-8000-000000000000','9074',NULL),('d0000008-0000-4000-8000-000000000000','3562',NULL),('d0000009-0000-4000-8000-000000000000','8419',NULL),('d0000010-0000-4000-8000-000000000000','5237',NULL),('fbe9708e-3919-4b15-b59f-767821715847',NULL,4),('TEST-LOCK-0001','CODETEST1',NULL),('TEST-LOCK-0002','CODETEST2',NULL);
/*!40000 ALTER TABLE `locker` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notification`
--

DROP TABLE IF EXISTS `notification`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notification` (
  `id_notification` char(36) NOT NULL,
  `view` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_notification`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notification`
--

LOCK TABLES `notification` WRITE;
/*!40000 ALTER TABLE `notification` DISABLE KEYS */;
INSERT INTO `notification` VALUES ('70000003-0000-4000-8000-000000000000',0),('70000004-0000-4000-8000-000000000000',1),('70000005-0000-4000-8000-000000000000',0),('70000006-0000-4000-8000-000000000000',1),('70000007-0000-4000-8000-000000000000',0),('70000008-0000-4000-8000-000000000000',0),('70000009-0000-4000-8000-000000000000',1),('70000010-0000-4000-8000-000000000000',0),('TEST-NOTIF-0001',0),('TEST-NOTIF-0002',1);
/*!40000 ALTER TABLE `notification` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payement`
--

DROP TABLE IF EXISTS `payement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `payement` (
  `id_payement` char(36) NOT NULL,
  `amount_cents` int NOT NULL,
  `currency` char(3) NOT NULL,
  `provider_payement` varchar(60) DEFAULT NULL,
  `provider_ref` varchar(255) DEFAULT NULL,
  `status_payement` varchar(60) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_payement`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payement`
--

LOCK TABLES `payement` WRITE;
/*!40000 ALTER TABLE `payement` DISABLE KEYS */;
INSERT INTO `payement` VALUES ('TEST-PAY-0001',1500,'EUR','TEST_PROVIDER','TEST_REF_1','PAID','2026-04-22 21:47:07','2026-04-22 21:47:07'),('TEST-PAY-0002',2500,'EUR','TEST_PROVIDER','TEST_REF_2','PENDING','2026-04-22 21:47:07',NULL);
/*!40000 ALTER TABLE `payement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payement_event`
--

DROP TABLE IF EXISTS `payement_event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `payement_event` (
  `id_payement` char(36) NOT NULL,
  `id_event` char(36) NOT NULL,
  PRIMARY KEY (`id_payement`,`id_event`),
  KEY `fk_payement_event_event` (`id_event`),
  CONSTRAINT `fk_payement_event_event` FOREIGN KEY (`id_event`) REFERENCES `event` (`id_event`),
  CONSTRAINT `fk_payement_event_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payement_event`
--

LOCK TABLES `payement_event` WRITE;
/*!40000 ALTER TABLE `payement_event` DISABLE KEYS */;
INSERT INTO `payement_event` VALUES ('TEST-PAY-0001','TEST-EVENT-0001'),('TEST-PAY-0002','TEST-EVENT-0002');
/*!40000 ALTER TABLE `payement_event` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payement_formation`
--

DROP TABLE IF EXISTS `payement_formation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `payement_formation` (
  `id_payement` char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  PRIMARY KEY (`id_payement`,`id_formation`),
  KEY `fk_payement_formation_formation` (`id_formation`),
  CONSTRAINT `fk_payement_formation_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`),
  CONSTRAINT `fk_payement_formation_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payement_formation`
--

LOCK TABLES `payement_formation` WRITE;
/*!40000 ALTER TABLE `payement_formation` DISABLE KEYS */;
INSERT INTO `payement_formation` VALUES ('TEST-PAY-0001','TEST-FORM-0001'),('TEST-PAY-0002','TEST-FORM-0002');
/*!40000 ALTER TABLE `payement_formation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payement_project`
--

DROP TABLE IF EXISTS `payement_project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `payement_project` (
  `id_payement` char(36) NOT NULL,
  `id_project` char(36) NOT NULL,
  PRIMARY KEY (`id_payement`,`id_project`),
  KEY `fk_payement_project_project` (`id_project`),
  CONSTRAINT `fk_payement_project_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`),
  CONSTRAINT `fk_payement_project_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payement_project`
--

LOCK TABLES `payement_project` WRITE;
/*!40000 ALTER TABLE `payement_project` DISABLE KEYS */;
INSERT INTO `payement_project` VALUES ('TEST-PAY-0001','TEST-PROJ-0001'),('TEST-PAY-0002','TEST-PROJ-0002');
/*!40000 ALTER TABLE `payement_project` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permission` (
  `id_permission` char(36) NOT NULL,
  `name_permission` varchar(100) NOT NULL,
  `domain` varchar(50) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_permission`),
  UNIQUE KEY `name_permission` (`name_permission`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
INSERT INTO `permission` VALUES ('0b1c2d3e-4f5a-4b6c-9d7e-8f9a0b1c2d3e','moderate_forum','forum','2026-04-22 21:47:08'),('0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b','manage_announcements','announcements','2026-04-22 21:47:08'),('1a2b3c4d-5e6f-4a7b-9c8d-0e1f2a3b3c4d','manage_payments','administration','2026-04-22 21:47:08'),('1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f','register_formation','formations','2026-04-22 21:47:08'),('1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a','create_event','events','2026-04-22 21:47:08'),('2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d','request_locker_deposit','lockers','2026-04-22 21:47:08'),('2c3d4e5f-6a7b-4c8d-ae0f-1a2b3c4d5e6f','manage_categories','administration','2026-04-22 21:47:08'),('2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a','manage_advice','conseils','2026-04-22 21:47:08'),('2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b','manage_projects','projects','2026-04-22 21:47:08'),('3f2d8a1b-4c5e-4f6a-9b0c-1d2e3f4a5b6c','manage_users','administration','2026-04-22 21:47:08'),('3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c','register_event','events','2026-04-22 21:47:08'),('4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d','create_post','forum','2026-04-22 21:47:08'),('4c5d6e7f-8a9b-4c0d-8e1f-2a3b4c5d6e7f','manage_lockers','lockers','2026-04-22 21:47:08'),('4f5a6b7c-8d9e-4f0a-8b1c-2d3e4f5a6b7c','send_notifications','administration','2026-04-22 21:47:08'),('5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e','manage_formations','formations','2026-04-22 21:47:08'),('5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a','create_announcement','announcements','2026-04-22 21:47:08'),('6b7c8d9e-0f1a-4b2c-9d3e-4f5a6b7c8d9e','manage_subscriptions','administration','2026-04-22 21:47:08'),('6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f','create_advice','conseils','2026-04-22 21:47:08'),('6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a','create_project','projects','2026-04-22 21:47:08'),('73875cfa-b134-4b6d-9a46-43ea4fe9b500','validate_deposit','lockers','2026-04-22 21:47:08'),('7a8b9c0d-1e2f-4a3b-8c4d-5e6f7a8b9c0d','manage_roles','administration','2026-04-22 21:47:08'),('7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b','manage_events','events','2026-04-22 21:47:08'),('7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c','buy_announcement','announcements','2026-04-22 21:47:08'),('8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e','pickup_locker','lockers','2026-04-22 21:47:08'),('8e9f0a1b-2c3d-4e5f-9a6b-7c8d9e0f1a2b','manage_professional_requests','administration','2026-04-22 21:47:08'),('8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c','create_topic','forum','2026-04-22 21:47:08'),('9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d','create_formation','formations','2026-04-22 21:47:08'),('9c0d1e2f-3a4b-4c5d-8e6f-7a8b9c0d1e2f','manage_documents','administration','2026-04-22 21:47:08');
/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `post`
--

DROP TABLE IF EXISTS `post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post` (
  `id_post` char(36) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_post`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post`
--

LOCK TABLES `post` WRITE;
/*!40000 ALTER TABLE `post` DISABLE KEYS */;
INSERT INTO `post` VALUES ('50000003-0000-4000-8000-000000000000',NULL),('50000004-0000-4000-8000-000000000000',NULL),('50000005-0000-4000-8000-000000000000',NULL),('50000006-0000-4000-8000-000000000000',NULL),('50000007-0000-4000-8000-000000000000',NULL),('50000008-0000-4000-8000-000000000000',NULL),('50000009-0000-4000-8000-000000000000',NULL),('50000010-0000-4000-8000-000000000000',NULL),('TEST-POST-0001',NULL),('TEST-POST-0002',NULL);
/*!40000 ALTER TABLE `post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `professional_request`
--

DROP TABLE IF EXISTS `professional_request`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `professional_request` (
  `id_request` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `processed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_request`),
  KEY `fk_professional_request_user` (`id_user`),
  CONSTRAINT `fk_professional_request_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `professional_request`
--

LOCK TABLES `professional_request` WRITE;
/*!40000 ALTER TABLE `professional_request` DISABLE KEYS */;
/*!40000 ALTER TABLE `professional_request` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `project`
--

DROP TABLE IF EXISTS `project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `project` (
  `id_project` char(36) NOT NULL,
  `start_date_project` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `id_creator` char(36) DEFAULT NULL,
  PRIMARY KEY (`id_project`),
  KEY `fk_project_creator` (`id_creator`),
  CONSTRAINT `fk_project_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `project`
--

LOCK TABLES `project` WRITE;
/*!40000 ALTER TABLE `project` DISABLE KEYS */;
INSERT INTO `project` VALUES ('f0000003-0000-4000-8000-000000000000','2026-05-01','2026-08-31',8,'TEST-USER-0001'),('f0000004-0000-4000-8000-000000000000','2026-06-01','2026-11-30',12,'TEST-USER-0001'),('f0000005-0000-4000-8000-000000000000','2026-04-15','2026-07-15',6,'TEST-USER-0002'),('f0000006-0000-4000-8000-000000000000','2026-07-01','2026-12-31',20,'TEST-USER-0002'),('f0000007-0000-4000-8000-000000000000','2026-03-01','2026-06-30',5,'TEST-USER-0001'),('f0000008-0000-4000-8000-000000000000','2026-08-01','2026-10-31',10,'TEST-USER-0001'),('f0000009-0000-4000-8000-000000000000','2026-05-15','2026-09-15',15,'TEST-USER-0002'),('f0000010-0000-4000-8000-000000000000','2026-09-01','2027-01-31',25,'TEST-USER-0002'),('TEST-PROJ-0001','2026-01-01','2026-06-01',10,NULL),('TEST-PROJ-0002','2026-02-01','2026-07-01',20,NULL);
/*!40000 ALTER TABLE `project` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `report`
--

DROP TABLE IF EXISTS `report`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `report` (
  `id_report` char(36) NOT NULL,
  `id_user` char(36) DEFAULT NULL,
  `id_announcement` char(36) DEFAULT NULL,
  `id_topic` char(36) DEFAULT NULL,
  `id_post` char(36) DEFAULT NULL,
  `reason` text,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `id_reported_user` char(36) DEFAULT NULL,
  `resolved_by` char(36) DEFAULT NULL,
  `action_taken` varchar(50) DEFAULT NULL,
  `resolved_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_report`),
  KEY `fk_report_user` (`id_user`),
  KEY `fk_report_announcement` (`id_announcement`),
  KEY `fk_report_topic` (`id_topic`),
  KEY `fk_report_post` (`id_post`),
  KEY `fk_report_reported_user` (`id_reported_user`),
  KEY `fk_report_resolved_by` (`resolved_by`),
  CONSTRAINT `fk_report_announcement` FOREIGN KEY (`id_announcement`) REFERENCES `announcement` (`id_announcement`),
  CONSTRAINT `fk_report_post` FOREIGN KEY (`id_post`) REFERENCES `post` (`id_post`),
  CONSTRAINT `fk_report_reported_user` FOREIGN KEY (`id_reported_user`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_report_resolved_by` FOREIGN KEY (`resolved_by`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_report_topic` FOREIGN KEY (`id_topic`) REFERENCES `topic` (`id_topic`),
  CONSTRAINT `fk_report_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `report`
--

LOCK TABLES `report` WRITE;
/*!40000 ALTER TABLE `report` DISABLE KEYS */;
INSERT INTO `report` VALUES ('80000001-0000-4000-8000-000000000000','TEST-USER-0001','c0000003-0000-4000-8000-000000000000','40000003-0000-4000-8000-000000000000','50000003-0000-4000-8000-000000000000','Contenu trompeur sur l\'état du produit','pending','2026-03-01 09:00:00',NULL,NULL,NULL,NULL),('80000002-0000-4000-8000-000000000000','TEST-USER-0002','c0000003-0000-4000-8000-000000000000','40000004-0000-4000-8000-000000000000','50000004-0000-4000-8000-000000000000','Prix incorrect indiqué','resolved','2026-03-02 10:30:00',NULL,NULL,NULL,NULL),('80000003-0000-4000-8000-000000000000','TEST-USER-0001','c0000004-0000-4000-8000-000000000000','40000005-0000-4000-8000-000000000000','50000005-0000-4000-8000-000000000000','Doublon d\'annonce existante','pending','2026-03-05 14:00:00',NULL,NULL,NULL,NULL),('80000004-0000-4000-8000-000000000000','TEST-USER-0002','c0000005-0000-4000-8000-000000000000','40000006-0000-4000-8000-000000000000','50000006-0000-4000-8000-000000000000','Matériel potentiellement dangereux non signalé','resolved','2026-03-08 11:15:00',NULL,NULL,NULL,NULL),('80000005-0000-4000-8000-000000000000','TEST-USER-0001','c0000006-0000-4000-8000-000000000000','40000007-0000-4000-8000-000000000000','50000007-0000-4000-8000-000000000000','Description inexacte de l\'état du produit','pending','2026-03-10 16:20:00',NULL,NULL,NULL,NULL),('80000006-0000-4000-8000-000000000000','TEST-USER-0002','c0000007-0000-4000-8000-000000000000','40000008-0000-4000-8000-000000000000','50000008-0000-4000-8000-000000000000','Annonce à caractère commercial non autorisé','pending','2026-03-12 08:45:00',NULL,NULL,NULL,NULL),('80000007-0000-4000-8000-000000000000','TEST-USER-0001','c0000008-0000-4000-8000-000000000000','40000009-0000-4000-8000-000000000000','50000009-0000-4000-8000-000000000000','Contenu hors sujet de la plateforme','pending','2026-03-15 13:00:00',NULL,NULL,NULL,NULL),('80000008-0000-4000-8000-000000000000','TEST-USER-0002','c0000009-0000-4000-8000-000000000000','40000010-0000-4000-8000-000000000000','50000010-0000-4000-8000-000000000000','Propos inappropriés dans la description','resolved','2026-03-18 09:30:00',NULL,NULL,NULL,NULL),('80000009-0000-4000-8000-000000000000','TEST-USER-0001','c0000010-0000-4000-8000-000000000000','40000003-0000-4000-8000-000000000000','50000003-0000-4000-8000-000000000000','Spam et liens externes suspects','pending','2026-03-20 15:10:00',NULL,NULL,NULL,NULL),('80000010-0000-4000-8000-000000000000','TEST-USER-0002','c0000003-0000-4000-8000-000000000000','40000004-0000-4000-8000-000000000000','50000004-0000-4000-8000-000000000000','Usurpation d\'identité d\'un artisan connu','resolved','2026-03-22 11:00:00',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `report` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role` (
  `id_role` char(36) NOT NULL,
  `name_role` varchar(100) NOT NULL,
  PRIMARY KEY (`id_role`),
  UNIQUE KEY `name_role` (`name_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role`
--

LOCK TABLES `role` WRITE;
/*!40000 ALTER TABLE `role` DISABLE KEYS */;
INSERT INTO `role` VALUES ('77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','admin'),('11111111-1111-1111-1111-000000000002','artisan'),('11111111-1111-1111-1111-000000000004','salarie'),('TEST-ROLE-0001','TEST_ADMIN'),('TEST-ROLE-0002','TEST_USER'),('11111111-1111-1111-1111-000000000003','user');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_permission`
--

DROP TABLE IF EXISTS `role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permission` (
  `id` char(36) NOT NULL,
  `role_id` char(36) NOT NULL,
  `permission_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_role_permission` (`role_id`,`permission_id`),
  KEY `fk_rp_permission` (`permission_id`),
  CONSTRAINT `fk_rp_permission` FOREIGN KEY (`permission_id`) REFERENCES `permission` (`id_permission`),
  CONSTRAINT `fk_rp_role` FOREIGN KEY (`role_id`) REFERENCES `role` (`id_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_permission`
--

LOCK TABLES `role_permission` WRITE;
/*!40000 ALTER TABLE `role_permission` DISABLE KEYS */;
INSERT INTO `role_permission` VALUES ('0bd46713-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b'),('0bd46ad6-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4634c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a'),('0bd46b6e-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd4667a-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a'),('0bd46908-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b'),('0bd46a3c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd46494-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd46858-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e'),('0bd4627e-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a'),('0bd46074-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f'),('0bd46531-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a'),('0bd46c10-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','73875cfa-b134-4b6d-9a46-43ea4fe9b500'),('0bd467a6-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b'),('0bd469a1-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd465d9-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd463ee-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d'),('0bd4e56a-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4e620-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd4e4b9-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd4e280-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd4dfce-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c'),('0bd4e406-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd4e34d-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd4a82c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4a8d5-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd4a57c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a'),('0bd4a625-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b'),('0bd4a778-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd4a377-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd4a2c8-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a'),('0bd4a1f6-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f'),('0bd4a422-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a'),('0bd49f8d-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c'),('0bd4a6cf-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd4a4ce-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd42ccf-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','0b1c2d3e-4f5a-4b6c-9d7e-8f9a0b1c2d3e'),('0bd425d5-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b'),('0bd42937-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','1a2b3c4d-5e6f-4a7b-9c8d-0e1f2a3b3c4d'),('0bd42e67-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4225e-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a'),('0bd42f02-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd42665-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2c3d4e5f-6a7b-4c8d-ae0f-1a2b3c4d5e6f'),('0bd4254a-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a'),('0bd42a5e-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b'),('0bd42c49-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','3f2d8a1b-4c5e-4f6a-9b0c-1d2e3f4a5b6c'),('0bd42ddd-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd4238b-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd428a7-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','4c5d6e7f-8a9b-4c0d-8e1f-2a3b4c5d6e7f'),('0bd42f91-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','4f5a6b7c-8d9e-4f0a-8b1c-2d3e4f5a6b7c'),('0bd4281a-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e'),('0bd421d1-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a'),('0bd42bbb-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','6b7c8d9e-0f1a-4b2c-9d3e-4f5a6b7c8d9e'),('0bd420f3-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f'),('0bd4241c-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a'),('0bd4301c-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','73875cfa-b134-4b6d-9a46-43ea4fe9b500'),('0bd42ae6-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','7a8b9c0d-1e2f-4a3b-8c4d-5e6f7a8b9c0d'),('0bd42792-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b'),('0bd41f11-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c'),('0bd42d54-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd429c6-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','8e9f0a1b-2c3d-4e5f-9a6b-7c8d9e0f1a2b'),('0bd424b3-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd422ef-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d'),('0bd42700-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','9c0d1e2f-3a4b-4c5d-8e6f-7a8b9c0d1e2f');
/*!40000 ALTER TABLE `role_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `schema_migrations`
--

LOCK TABLES `schema_migrations` WRITE;
/*!40000 ALTER TABLE `schema_migrations` DISABLE KEYS */;
INSERT INTO `schema_migrations` VALUES (17,0);
/*!40000 ALTER TABLE `schema_migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sub_sub_plan`
--

DROP TABLE IF EXISTS `sub_sub_plan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sub_sub_plan` (
  `id_subscription` char(36) NOT NULL,
  `id_plan` char(36) NOT NULL,
  PRIMARY KEY (`id_subscription`,`id_plan`),
  KEY `fk_sub_sub_plan_plan` (`id_plan`),
  CONSTRAINT `fk_sub_sub_plan_plan` FOREIGN KEY (`id_plan`) REFERENCES `subscription_plans` (`id_plan`),
  CONSTRAINT `fk_sub_sub_plan_subscription` FOREIGN KEY (`id_subscription`) REFERENCES `subscription` (`id_subscription`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sub_sub_plan`
--

LOCK TABLES `sub_sub_plan` WRITE;
/*!40000 ALTER TABLE `sub_sub_plan` DISABLE KEYS */;
INSERT INTO `sub_sub_plan` VALUES ('20000003-0000-4000-8000-000000000000','10000003-0000-4000-8000-000000000000'),('20000004-0000-4000-8000-000000000000','10000004-0000-4000-8000-000000000000'),('20000005-0000-4000-8000-000000000000','10000005-0000-4000-8000-000000000000'),('20000006-0000-4000-8000-000000000000','10000006-0000-4000-8000-000000000000'),('20000007-0000-4000-8000-000000000000','10000007-0000-4000-8000-000000000000'),('20000008-0000-4000-8000-000000000000','10000008-0000-4000-8000-000000000000'),('20000009-0000-4000-8000-000000000000','10000009-0000-4000-8000-000000000000'),('20000010-0000-4000-8000-000000000000','10000010-0000-4000-8000-000000000000'),('TEST-SUB-0001','TEST-PLAN-0001'),('TEST-SUB-0002','TEST-PLAN-0002');
/*!40000 ALTER TABLE `sub_sub_plan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subscription`
--

DROP TABLE IF EXISTS `subscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `subscription` (
  `id_subscription` char(36) NOT NULL,
  `start_timestamp` datetime DEFAULT NULL,
  `end_timestamp` datetime DEFAULT NULL,
  `cancelled` tinyint(1) NOT NULL DEFAULT '0',
  `cancelled_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_subscription`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subscription`
--

LOCK TABLES `subscription` WRITE;
/*!40000 ALTER TABLE `subscription` DISABLE KEYS */;
INSERT INTO `subscription` VALUES ('20000003-0000-4000-8000-000000000000','2026-01-01 00:00:00','2026-12-31 23:59:59',0,'2026-12-31 23:59:59'),('20000004-0000-4000-8000-000000000000','2026-02-01 00:00:00','2027-01-31 23:59:59',0,'2027-01-31 23:59:59'),('20000005-0000-4000-8000-000000000000','2025-11-01 00:00:00','2026-10-31 23:59:59',0,'2026-10-31 23:59:59'),('20000006-0000-4000-8000-000000000000','2026-03-15 00:00:00','2027-03-14 23:59:59',0,'2027-03-14 23:59:59'),('20000007-0000-4000-8000-000000000000','2025-06-01 00:00:00','2026-05-31 23:59:59',1,'2026-04-01 10:00:00'),('20000008-0000-4000-8000-000000000000','2026-04-01 00:00:00','2027-03-31 23:59:59',0,'2027-03-31 23:59:59'),('20000009-0000-4000-8000-000000000000','2025-12-01 00:00:00','2026-11-30 23:59:59',0,'2026-11-30 23:59:59'),('20000010-0000-4000-8000-000000000000','2026-01-15 00:00:00','2027-01-14 23:59:59',1,'2026-03-20 09:30:00'),('TEST-SUB-0001','2026-04-22 21:47:07',NULL,0,NULL),('TEST-SUB-0002','2026-04-22 21:47:07',NULL,0,NULL);
/*!40000 ALTER TABLE `subscription` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subscription_plans`
--

DROP TABLE IF EXISTS `subscription_plans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `subscription_plans` (
  `id_plan` char(36) NOT NULL,
  `price_cents` int NOT NULL,
  `interval_unit` varchar(30) NOT NULL,
  `interval_count` int NOT NULL DEFAULT '1',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id_plan`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subscription_plans`
--

LOCK TABLES `subscription_plans` WRITE;
/*!40000 ALTER TABLE `subscription_plans` DISABLE KEYS */;
INSERT INTO `subscription_plans` VALUES ('10000003-0000-4000-8000-000000000000',299,'MONTH',1,1),('10000004-0000-4000-8000-000000000000',2990,'YEAR',1,1),('10000005-0000-4000-8000-000000000000',499,'MONTH',1,1),('10000006-0000-4000-8000-000000000000',4990,'YEAR',1,1),('10000007-0000-4000-8000-000000000000',799,'MONTH',1,1),('10000008-0000-4000-8000-000000000000',7990,'YEAR',1,1),('10000009-0000-4000-8000-000000000000',149,'MONTH',3,1),('10000010-0000-4000-8000-000000000000',1490,'YEAR',3,1),('TEST-PLAN-0001',999,'MONTH',1,1),('TEST-PLAN-0002',9999,'YEAR',1,1);
/*!40000 ALTER TABLE `subscription_plans` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ticket`
--

DROP TABLE IF EXISTS `ticket`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ticket` (
  `id_ticket` char(36) NOT NULL,
  `timestamp_ticket` datetime DEFAULT NULL,
  PRIMARY KEY (`id_ticket`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ticket`
--

LOCK TABLES `ticket` WRITE;
/*!40000 ALTER TABLE `ticket` DISABLE KEYS */;
INSERT INTO `ticket` VALUES ('60000003-0000-4000-8000-000000000000','2026-05-01 10:30:00'),('60000004-0000-4000-8000-000000000000','2026-05-10 14:15:00'),('60000005-0000-4000-8000-000000000000','2026-05-18 09:00:00'),('60000006-0000-4000-8000-000000000000','2026-06-02 16:45:00'),('60000007-0000-4000-8000-000000000000','2026-06-14 11:00:00'),('60000008-0000-4000-8000-000000000000','2026-07-05 09:30:00'),('60000009-0000-4000-8000-000000000000','2026-07-19 14:00:00'),('60000010-0000-4000-8000-000000000000','2026-08-22 10:00:00'),('TEST-TICKET-0001','2026-04-22 21:47:07'),('TEST-TICKET-0002','2026-04-22 21:47:07');
/*!40000 ALTER TABLE `ticket` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `topic`
--

DROP TABLE IF EXISTS `topic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `topic` (
  `id_topic` char(36) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `topic`
--

LOCK TABLES `topic` WRITE;
/*!40000 ALTER TABLE `topic` DISABLE KEYS */;
INSERT INTO `topic` VALUES ('40000003-0000-4000-8000-000000000000',NULL),('40000004-0000-4000-8000-000000000000',NULL),('40000005-0000-4000-8000-000000000000',NULL),('40000006-0000-4000-8000-000000000000',NULL),('40000007-0000-4000-8000-000000000000',NULL),('40000008-0000-4000-8000-000000000000',NULL),('40000009-0000-4000-8000-000000000000',NULL),('40000010-0000-4000-8000-000000000000',NULL),('TEST-TOPIC-0001',NULL),('TEST-TOPIC-0002',NULL);
/*!40000 ALTER TABLE `topic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `topic_post`
--

DROP TABLE IF EXISTS `topic_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `topic_post` (
  `id_topic` char(36) NOT NULL,
  `id_post` char(36) NOT NULL,
  PRIMARY KEY (`id_topic`,`id_post`),
  KEY `fk_topic_post_post` (`id_post`),
  CONSTRAINT `fk_topic_post_post` FOREIGN KEY (`id_post`) REFERENCES `post` (`id_post`),
  CONSTRAINT `fk_topic_post_topic` FOREIGN KEY (`id_topic`) REFERENCES `topic` (`id_topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `topic_post`
--

LOCK TABLES `topic_post` WRITE;
/*!40000 ALTER TABLE `topic_post` DISABLE KEYS */;
INSERT INTO `topic_post` VALUES ('40000003-0000-4000-8000-000000000000','50000003-0000-4000-8000-000000000000'),('40000004-0000-4000-8000-000000000000','50000004-0000-4000-8000-000000000000'),('40000005-0000-4000-8000-000000000000','50000005-0000-4000-8000-000000000000'),('40000006-0000-4000-8000-000000000000','50000006-0000-4000-8000-000000000000'),('40000007-0000-4000-8000-000000000000','50000007-0000-4000-8000-000000000000'),('40000008-0000-4000-8000-000000000000','50000008-0000-4000-8000-000000000000'),('40000009-0000-4000-8000-000000000000','50000009-0000-4000-8000-000000000000'),('40000010-0000-4000-8000-000000000000','50000010-0000-4000-8000-000000000000'),('TEST-TOPIC-0001','TEST-POST-0001'),('TEST-TOPIC-0002','TEST-POST-0002');
/*!40000 ALTER TABLE `topic_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id_user` char(36) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_user` varchar(255) NOT NULL,
  `first_name` varchar(150) DEFAULT NULL,
  `last_name` varchar(150) DEFAULT NULL,
  `upcycling_score` int NOT NULL DEFAULT '0',
  `premium` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(20) NOT NULL DEFAULT 'active',
  PRIMARY KEY (`id_user`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES ('7cd909a9-a27a-4849-a867-55e39b871c66','admin@upcycle.com','$2a$10$JHayhaAb0AYcARtBUWkyuuxf2fN64DydnGgbHSVPQfA60uvY5mo5K','Admin','Upcycle',0,0,'2026-04-22 21:47:08','active'),('TEST-USER-0001','test_user1@example.com','hashed_pwd_test1','Test','UserOne',10,0,'2026-04-22 21:47:07','active'),('TEST-USER-0002','test_user2@example.com','hashed_pwd_test2','Test','UserTwo',25,1,'2026-04-22 21:47:07','active');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_advice`
--

DROP TABLE IF EXISTS `user_advice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_advice` (
  `id_user` char(36) NOT NULL,
  `id_advice` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_advice`),
  KEY `fk_user_advice_advice` (`id_advice`),
  CONSTRAINT `fk_user_advice_advice` FOREIGN KEY (`id_advice`) REFERENCES `advice` (`id_advice`),
  CONSTRAINT `fk_user_advice_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_advice`
--

LOCK TABLES `user_advice` WRITE;
/*!40000 ALTER TABLE `user_advice` DISABLE KEYS */;
INSERT INTO `user_advice` VALUES ('TEST-USER-0001','30000003-0000-4000-8000-000000000000'),('TEST-USER-0001','30000004-0000-4000-8000-000000000000'),('TEST-USER-0001','30000005-0000-4000-8000-000000000000'),('TEST-USER-0001','30000006-0000-4000-8000-000000000000'),('TEST-USER-0002','30000007-0000-4000-8000-000000000000'),('TEST-USER-0002','30000008-0000-4000-8000-000000000000'),('TEST-USER-0002','30000009-0000-4000-8000-000000000000'),('TEST-USER-0002','30000010-0000-4000-8000-000000000000'),('TEST-USER-0001','TEST-ADV-0001'),('TEST-USER-0002','TEST-ADV-0002');
/*!40000 ALTER TABLE `user_advice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_announcement`
--

DROP TABLE IF EXISTS `user_announcement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_announcement` (
  `id_user` char(36) NOT NULL,
  `id_announcement` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_announcement`),
  KEY `fk_user_announcement_announcement` (`id_announcement`),
  CONSTRAINT `fk_user_announcement_announcement` FOREIGN KEY (`id_announcement`) REFERENCES `announcement` (`id_announcement`),
  CONSTRAINT `fk_user_announcement_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_announcement`
--

LOCK TABLES `user_announcement` WRITE;
/*!40000 ALTER TABLE `user_announcement` DISABLE KEYS */;
INSERT INTO `user_announcement` VALUES ('TEST-USER-0001','c0000003-0000-4000-8000-000000000000'),('TEST-USER-0001','c0000004-0000-4000-8000-000000000000'),('TEST-USER-0001','c0000005-0000-4000-8000-000000000000'),('TEST-USER-0001','c0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','c0000007-0000-4000-8000-000000000000'),('TEST-USER-0002','c0000008-0000-4000-8000-000000000000'),('TEST-USER-0002','c0000009-0000-4000-8000-000000000000'),('TEST-USER-0002','c0000010-0000-4000-8000-000000000000'),('TEST-USER-0001','TEST-ANN-0001'),('TEST-USER-0002','TEST-ANN-0002');
/*!40000 ALTER TABLE `user_announcement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_event_inscription`
--

DROP TABLE IF EXISTS `user_event_inscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_event_inscription` (
  `id_user` char(36) NOT NULL,
  `id_event` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_event`),
  KEY `fk_user_event_inscription_event` (`id_event`),
  CONSTRAINT `fk_user_event_inscription_event` FOREIGN KEY (`id_event`) REFERENCES `event` (`id_event`),
  CONSTRAINT `fk_user_event_inscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_event_inscription`
--

LOCK TABLES `user_event_inscription` WRITE;
/*!40000 ALTER TABLE `user_event_inscription` DISABLE KEYS */;
INSERT INTO `user_event_inscription` VALUES ('TEST-USER-0001','e0000003-0000-4000-8000-000000000000'),('TEST-USER-0002','e0000003-0000-4000-8000-000000000000'),('TEST-USER-0001','e0000004-0000-4000-8000-000000000000'),('TEST-USER-0002','e0000004-0000-4000-8000-000000000000'),('TEST-USER-0001','e0000005-0000-4000-8000-000000000000'),('TEST-USER-0002','e0000005-0000-4000-8000-000000000000'),('TEST-USER-0001','e0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','e0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','TEST-EVENT-0001'),('TEST-USER-0001','TEST-EVENT-0002');
/*!40000 ALTER TABLE `user_event_inscription` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_formation_inscription`
--

DROP TABLE IF EXISTS `user_formation_inscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_formation_inscription` (
  `id_user` char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_formation`),
  KEY `fk_user_formation_inscription_formation` (`id_formation`),
  CONSTRAINT `fk_user_formation_inscription_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`),
  CONSTRAINT `fk_user_formation_inscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_formation_inscription`
--

LOCK TABLES `user_formation_inscription` WRITE;
/*!40000 ALTER TABLE `user_formation_inscription` DISABLE KEYS */;
INSERT INTO `user_formation_inscription` VALUES ('TEST-USER-0001','b0000003-0000-4000-8000-000000000000'),('TEST-USER-0002','b0000003-0000-4000-8000-000000000000'),('TEST-USER-0001','b0000004-0000-4000-8000-000000000000'),('TEST-USER-0002','b0000004-0000-4000-8000-000000000000'),('TEST-USER-0001','b0000005-0000-4000-8000-000000000000'),('TEST-USER-0002','b0000005-0000-4000-8000-000000000000'),('TEST-USER-0001','b0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','b0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','TEST-FORM-0001'),('TEST-USER-0001','TEST-FORM-0002');
/*!40000 ALTER TABLE `user_formation_inscription` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_formation_inscription_payement`
--

DROP TABLE IF EXISTS `user_formation_inscription_payement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_formation_inscription_payement` (
  `id_user` char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  `id_payement` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_formation`,`id_payement`),
  KEY `fk_ufip_formation` (`id_formation`),
  KEY `fk_ufip_payement` (`id_payement`),
  CONSTRAINT `fk_ufip_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`),
  CONSTRAINT `fk_ufip_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`),
  CONSTRAINT `fk_ufip_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_formation_inscription_payement`
--

LOCK TABLES `user_formation_inscription_payement` WRITE;
/*!40000 ALTER TABLE `user_formation_inscription_payement` DISABLE KEYS */;
INSERT INTO `user_formation_inscription_payement` VALUES ('TEST-USER-0001','TEST-FORM-0001','TEST-PAY-0001'),('TEST-USER-0002','TEST-FORM-0002','TEST-PAY-0002');
/*!40000 ALTER TABLE `user_formation_inscription_payement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_notification`
--

DROP TABLE IF EXISTS `user_notification`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_notification` (
  `id_user` char(36) NOT NULL,
  `id_notification` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_notification`),
  KEY `fk_user_notification_notification` (`id_notification`),
  CONSTRAINT `fk_user_notification_notification` FOREIGN KEY (`id_notification`) REFERENCES `notification` (`id_notification`),
  CONSTRAINT `fk_user_notification_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_notification`
--

LOCK TABLES `user_notification` WRITE;
/*!40000 ALTER TABLE `user_notification` DISABLE KEYS */;
INSERT INTO `user_notification` VALUES ('TEST-USER-0001','70000003-0000-4000-8000-000000000000'),('TEST-USER-0001','70000004-0000-4000-8000-000000000000'),('TEST-USER-0001','70000005-0000-4000-8000-000000000000'),('TEST-USER-0001','70000006-0000-4000-8000-000000000000'),('TEST-USER-0002','70000007-0000-4000-8000-000000000000'),('TEST-USER-0002','70000008-0000-4000-8000-000000000000'),('TEST-USER-0002','70000009-0000-4000-8000-000000000000'),('TEST-USER-0002','70000010-0000-4000-8000-000000000000'),('TEST-USER-0001','TEST-NOTIF-0001'),('TEST-USER-0002','TEST-NOTIF-0002');
/*!40000 ALTER TABLE `user_notification` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_payement`
--

DROP TABLE IF EXISTS `user_payement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_payement` (
  `id_user` char(36) NOT NULL,
  `id_payement` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_payement`),
  KEY `fk_user_payement_payement` (`id_payement`),
  CONSTRAINT `fk_user_payement_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`),
  CONSTRAINT `fk_user_payement_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_payement`
--

LOCK TABLES `user_payement` WRITE;
/*!40000 ALTER TABLE `user_payement` DISABLE KEYS */;
INSERT INTO `user_payement` VALUES ('TEST-USER-0001','TEST-PAY-0001'),('TEST-USER-0002','TEST-PAY-0002');
/*!40000 ALTER TABLE `user_payement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_project_inscription`
--

DROP TABLE IF EXISTS `user_project_inscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_project_inscription` (
  `id_user` char(36) NOT NULL,
  `id_project` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_project`),
  KEY `fk_user_project_inscription_project` (`id_project`),
  CONSTRAINT `fk_user_project_inscription_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`),
  CONSTRAINT `fk_user_project_inscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_project_inscription`
--

LOCK TABLES `user_project_inscription` WRITE;
/*!40000 ALTER TABLE `user_project_inscription` DISABLE KEYS */;
INSERT INTO `user_project_inscription` VALUES ('TEST-USER-0001','f0000003-0000-4000-8000-000000000000'),('TEST-USER-0002','f0000003-0000-4000-8000-000000000000'),('TEST-USER-0001','f0000004-0000-4000-8000-000000000000'),('TEST-USER-0002','f0000004-0000-4000-8000-000000000000'),('TEST-USER-0001','f0000005-0000-4000-8000-000000000000'),('TEST-USER-0002','f0000005-0000-4000-8000-000000000000'),('TEST-USER-0001','f0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','f0000006-0000-4000-8000-000000000000'),('TEST-USER-0002','TEST-PROJ-0001'),('TEST-USER-0001','TEST-PROJ-0002');
/*!40000 ALTER TABLE `user_project_inscription` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_project_updown`
--

DROP TABLE IF EXISTS `user_project_updown`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_project_updown` (
  `id_user` char(36) NOT NULL,
  `id_project` char(36) NOT NULL,
  `updown` smallint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_user`,`id_project`),
  KEY `fk_user_project_updown_project` (`id_project`),
  CONSTRAINT `fk_user_project_updown_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`),
  CONSTRAINT `fk_user_project_updown_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_project_updown`
--

LOCK TABLES `user_project_updown` WRITE;
/*!40000 ALTER TABLE `user_project_updown` DISABLE KEYS */;
INSERT INTO `user_project_updown` VALUES ('TEST-USER-0001','f0000003-0000-4000-8000-000000000000',1),('TEST-USER-0001','f0000004-0000-4000-8000-000000000000',1),('TEST-USER-0001','f0000005-0000-4000-8000-000000000000',-1),('TEST-USER-0001','f0000006-0000-4000-8000-000000000000',1),('TEST-USER-0001','TEST-PROJ-0001',1),('TEST-USER-0002','f0000003-0000-4000-8000-000000000000',-1),('TEST-USER-0002','f0000004-0000-4000-8000-000000000000',1),('TEST-USER-0002','f0000005-0000-4000-8000-000000000000',1),('TEST-USER-0002','f0000006-0000-4000-8000-000000000000',-1),('TEST-USER-0002','TEST-PROJ-0002',-1);
/*!40000 ALTER TABLE `user_project_updown` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_role`
--

DROP TABLE IF EXISTS `user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_role` (
  `id_user` char(36) NOT NULL,
  `id_role` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_role`),
  KEY `fk_user_role_role` (`id_role`),
  CONSTRAINT `fk_user_role_role` FOREIGN KEY (`id_role`) REFERENCES `role` (`id_role`),
  CONSTRAINT `fk_user_role_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_role`
--

LOCK TABLES `user_role` WRITE;
/*!40000 ALTER TABLE `user_role` DISABLE KEYS */;
INSERT INTO `user_role` VALUES ('TEST-USER-0001','11111111-1111-1111-1111-000000000002'),('TEST-USER-0002','11111111-1111-1111-1111-000000000002'),('TEST-USER-0001','11111111-1111-1111-1111-000000000003'),('TEST-USER-0002','11111111-1111-1111-1111-000000000003'),('TEST-USER-0001','11111111-1111-1111-1111-000000000004'),('TEST-USER-0002','11111111-1111-1111-1111-000000000004'),('7cd909a9-a27a-4849-a867-55e39b871c66','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab'),('TEST-USER-0001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab'),('TEST-USER-0002','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab'),('TEST-USER-0001','TEST-ROLE-0001'),('TEST-USER-0002','TEST-ROLE-0002');
/*!40000 ALTER TABLE `user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_sanctions`
--

DROP TABLE IF EXISTS `user_sanctions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_sanctions` (
  `id_sanction` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `id_admin` char(36) NOT NULL,
  `id_report` char(36) DEFAULT NULL,
  `type` varchar(20) NOT NULL,
  `reason` text,
  `duration_days` int DEFAULT NULL,
  `expires_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_sanction`),
  KEY `fk_sanctions_user` (`id_user`),
  KEY `fk_sanctions_admin` (`id_admin`),
  KEY `fk_sanctions_report` (`id_report`),
  CONSTRAINT `fk_sanctions_admin` FOREIGN KEY (`id_admin`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_sanctions_report` FOREIGN KEY (`id_report`) REFERENCES `report` (`id_report`),
  CONSTRAINT `fk_sanctions_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_sanctions`
--

LOCK TABLES `user_sanctions` WRITE;
/*!40000 ALTER TABLE `user_sanctions` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_sanctions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_subscription`
--

DROP TABLE IF EXISTS `user_subscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_subscription` (
  `id_user` char(36) NOT NULL,
  `id_subscription` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_subscription`),
  KEY `fk_user_subscription_subscription` (`id_subscription`),
  CONSTRAINT `fk_user_subscription_subscription` FOREIGN KEY (`id_subscription`) REFERENCES `subscription` (`id_subscription`),
  CONSTRAINT `fk_user_subscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_subscription`
--

LOCK TABLES `user_subscription` WRITE;
/*!40000 ALTER TABLE `user_subscription` DISABLE KEYS */;
INSERT INTO `user_subscription` VALUES ('TEST-USER-0001','20000003-0000-4000-8000-000000000000'),('TEST-USER-0001','20000004-0000-4000-8000-000000000000'),('TEST-USER-0001','20000005-0000-4000-8000-000000000000'),('TEST-USER-0001','20000006-0000-4000-8000-000000000000'),('TEST-USER-0002','20000007-0000-4000-8000-000000000000'),('TEST-USER-0002','20000008-0000-4000-8000-000000000000'),('TEST-USER-0002','20000009-0000-4000-8000-000000000000'),('TEST-USER-0002','20000010-0000-4000-8000-000000000000'),('TEST-USER-0001','TEST-SUB-0001'),('TEST-USER-0002','TEST-SUB-0002');
/*!40000 ALTER TABLE `user_subscription` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_ticket`
--

DROP TABLE IF EXISTS `user_ticket`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_ticket` (
  `id_user` char(36) NOT NULL,
  `id_ticket` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_ticket`),
  KEY `fk_user_ticket_ticket` (`id_ticket`),
  CONSTRAINT `fk_user_ticket_ticket` FOREIGN KEY (`id_ticket`) REFERENCES `ticket` (`id_ticket`),
  CONSTRAINT `fk_user_ticket_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_ticket`
--

LOCK TABLES `user_ticket` WRITE;
/*!40000 ALTER TABLE `user_ticket` DISABLE KEYS */;
INSERT INTO `user_ticket` VALUES ('TEST-USER-0001','60000003-0000-4000-8000-000000000000'),('TEST-USER-0001','60000004-0000-4000-8000-000000000000'),('TEST-USER-0001','60000005-0000-4000-8000-000000000000'),('TEST-USER-0001','60000006-0000-4000-8000-000000000000'),('TEST-USER-0002','60000007-0000-4000-8000-000000000000'),('TEST-USER-0002','60000008-0000-4000-8000-000000000000'),('TEST-USER-0002','60000009-0000-4000-8000-000000000000'),('TEST-USER-0002','60000010-0000-4000-8000-000000000000'),('TEST-USER-0001','TEST-TICKET-0001'),('TEST-USER-0002','TEST-TICKET-0002');
/*!40000 ALTER TABLE `user_ticket` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_topic`
--

DROP TABLE IF EXISTS `user_topic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_topic` (
  `id_user` char(36) NOT NULL,
  `id_topic` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_topic`),
  KEY `fk_user_topic_topic` (`id_topic`),
  CONSTRAINT `fk_user_topic_topic` FOREIGN KEY (`id_topic`) REFERENCES `topic` (`id_topic`),
  CONSTRAINT `fk_user_topic_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_topic`
--

LOCK TABLES `user_topic` WRITE;
/*!40000 ALTER TABLE `user_topic` DISABLE KEYS */;
INSERT INTO `user_topic` VALUES ('TEST-USER-0001','40000003-0000-4000-8000-000000000000'),('TEST-USER-0001','40000004-0000-4000-8000-000000000000'),('TEST-USER-0001','40000005-0000-4000-8000-000000000000'),('TEST-USER-0001','40000006-0000-4000-8000-000000000000'),('TEST-USER-0002','40000007-0000-4000-8000-000000000000'),('TEST-USER-0002','40000008-0000-4000-8000-000000000000'),('TEST-USER-0002','40000009-0000-4000-8000-000000000000'),('TEST-USER-0002','40000010-0000-4000-8000-000000000000'),('TEST-USER-0001','TEST-TOPIC-0001'),('TEST-USER-0002','TEST-TOPIC-0002');
/*!40000 ALTER TABLE `user_topic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'upcycle'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-22 22:05:03
