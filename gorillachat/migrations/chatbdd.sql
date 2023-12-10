-- Adminer 4.8.1 MySQL 8.0.22-13 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `Rooms`;
CREATE TABLE `Rooms` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `name` varchar(255) NOT NULL,
                         `description` text NOT NULL,
                         `private` tinyint(1) DEFAULT '0',
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `Rooms` (`id`, `name`, `description`, `private`) VALUES
                                                                 (33,	'bistrot Birdies',	'Boire du café et jouer',	0),
                                                                 (35,	'Ragnarok',	'Parler de l\'aprés chez les Vikings',	0),
(36,	'école de nassim',	'Apprendre l\'arabe littéraire avec Nassim',	0),
                                                                 (37,	'cabinet des curiosités',	'Parler de choses surprenantes',	0),
                                                                 (39,	'Mariage',	'mariées',	0),
                                                                 (40,	'ramen',	'Cuisine Japonaise',	0),
                                                                 (41,	'Sénégale',	'Dakar',	0);

DROP TABLE IF EXISTS `User_Room`;
CREATE TABLE `User_Room` (
                             `id` int NOT NULL AUTO_INCREMENT,
                             `user_id` int NOT NULL,
                             `room_id` int NOT NULL,
                             PRIMARY KEY (`id`),
                             KEY `user_id` (`user_id`),
                             KEY `room_id` (`room_id`),
                             CONSTRAINT `room_id` FOREIGN KEY (`room_id`) REFERENCES `Rooms` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
                             CONSTRAINT `user_id` FOREIGN KEY (`user_id`) REFERENCES `Users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `User_Room` (`id`, `user_id`, `room_id`) VALUES
                                                         (31,	31,	33),
                                                         (32,	35,	33),
                                                         (35,	35,	35),
                                                         (36,	31,	35),
                                                         (37,	31,	36),
                                                         (38,	31,	39),
                                                         (39,	43,	35),
                                                         (40,	43,	33);

DROP TABLE IF EXISTS `Users`;
CREATE TABLE `Users` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `username` varchar(255) NOT NULL,
                         `password` varchar(255) NOT NULL,
                         `admin` tinyint(1) DEFAULT '0',
                         `email` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `Users` (`id`, `username`, `password`, `admin`, `email`) VALUES
                                                                         (31,	'nass',	'nass',	1,	'waitingemail@noreply.com'),
                                                                         (35,	'amaury',	'Sky@900',	1,	'amauryfranssen@gmail.com'),
                                                                         (43,	'amo',	'Sky@900',	0,	'waiting@noemail.com'),
                                                                         (44,	'khalifa',	'khalifa',	0,	'waiting@email.com');

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
                            `id` int NOT NULL AUTO_INCREMENT,
                            `room_id` int NOT NULL,
                            `user_id` int NOT NULL,
                            `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                            `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
                            `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `Status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `user_id` (`user_id`),
                            KEY `room_id` (`room_id`),
                            CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `Users` (`id`) ON DELETE CASCADE,
                            CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`room_id`) REFERENCES `Rooms` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `messages` (`id`, `room_id`, `user_id`, `username`, `content`, `created_at`, `Status`) VALUES
                                                                                                       (66,	33,	31,	'nass',	'hello',	'2023-12-01 23:01:25',	NULL),
                                                                                                       (67,	33,	31,	'nass',	'hi',	'2023-12-01 23:01:55',	NULL),
                                                                                                       (68,	33,	31,	'nass',	'bonjour',	'2023-12-01 23:05:54',	NULL);

-- 2023-12-10 22:55:39