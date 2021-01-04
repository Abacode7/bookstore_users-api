# bookstore_users-api

Create users table
```
CREATE TABLE `userdb`.`users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(45) NULL,
  `last_name` VARCHAR(45) NULL,
  `email` VARCHAR(45) NOT NULL,
  `date_created` DATETIME NOT NULL,
  `status` VARCHAR(45) NULL,
  `password` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC) VISIBLE);
```