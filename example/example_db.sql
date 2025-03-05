-- -----------------------------------------------------
-- Schema main_db
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `main_db` DEFAULT CHARACTER SET utf8 ;
USE `main_db` ;

-- -----------------------------------------------------
-- Table `main_db`.`Users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`Users` (
  `UserID` INT NOT NULL AUTO_INCREMENT,
  `Username` VARCHAR(45) NOT NULL,
  `Password(Hash)` VARCHAR(45) NOT NULL,
  `Wallet` INT NOT NULL,
  `Role` INT NOT NULL,
  PRIMARY KEY (`UserID`)
);

-- -----------------------------------------------------
-- Table `main_db`.`ItemTypes`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`ItemTypes` (
  `TypeID` INT NOT NULL AUTO_INCREMENT,
  `ItemName` VARCHAR(45) NOT NULL,
  `ItemDescription` VARCHAR(45) NULL,
  `ImgURL` VARCHAR(45) NULL,
  `ShortDescription` VARCHAR(255) NULL,
  FULLTEXT (ItemName)
  PRIMARY KEY (`TypeID`)
);


-- -----------------------------------------------------
-- Table `main_db`.`Inventory`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`Inventory` (
  `ItemID` INT NOT NULL AUTO_INCREMENT,
  `UserID` INT NOT NULL,
  `TypeID` INT NOT NULL,

  PRIMARY KEY (`ItemID`),
  -- INDEX `UserID_idx` (`UserID`), -- to search for a users inventory fast

  CONSTRAINT `InvUserID` -- a item has a owner
    FOREIGN KEY (`UserID`)
    REFERENCES `main_db`.`Users` (`UserID`),

  CONSTRAINT `InvTypeID` -- a item has a type (which also gives it is name)
    FOREIGN KEY (`TypeID`)
    REFERENCES `main_db`.`ItemTypes` (`TypeID`)
);


-- -----------------------------------------------------
-- Table `main_db`.`Marketplace`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`Marketplace` (
  `OfferID` INT NOT NULL AUTO_INCREMENT,
  `ItemID` INT NOT NULL UNIQUE,
  `Price` INT NOT NULL,
  `CreationDate` DATE NOT NULL,

  PRIMARY KEY (`OfferID`),
  CONSTRAINT `MarketItemID`
    FOREIGN KEY (`ItemID`)
    REFERENCES `main_db`.`Inventory` (`ItemID`),
  FULLTEXT(ItemName
);


-- -----------------------------------------------------
-- Table `main_db`.`Transaction Log`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`TransactionLog` (
  `TransID` INT NOT NULL AUTO_INCREMENT,
  `Price` INT NOT NULL,
  `Date` DATE NOT NULL,
  `ItemID` INT NOT NULL,
  `Buyer` INT  NOT NULL,
  `Seller` INT  NOT NULL,

  PRIMARY KEY (`TransID`),

  CONSTRAINT `TransactionSeller`
    FOREIGN KEY (Seller)
    REFERENCES `main_db`.`Users` (UserID),
  
  CONSTRAINT `TransactionBuyer`
    FOREIGN KEY (Buyer)
    REFERENCES `main_db`.`Users` (UserID),

  CONSTRAINT `TransactionItemID`
    FOREIGN KEY (`ItemID`)
    REFERENCES `main_db`.`Inventory` (`ItemID`)
);

-- -----------------------------------------------------
-- Table `main_db`.`TokenTable`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`TokenTable` (
  `Token` BINARY(16) NOT NULL,
  `UserID` INT NOT NULL,
  `CreatedOn` DATETIME NOT NULL,
  
  PRIMARY KEY (`Token`),

  CONSTRAINT `TokenUser`
    FOREIGN key (`UserID`)
    REFERENCES `main_db`.`Users` (`UserID`)
);


-- -----------------------------------------------------
-- Table `main_db`.`TypeComments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `main_db`.`TypeComments` (
  `CommentID` INT NOT NULL AUTO_INCREMENT,
  `TypeID` INT NOT NULL,
  `UserID` INT NOT NULL,
  `Grade` INT NOT NULL,
  `Comment` VARCHAR(255) NOT NULL,
  `CreatedOn` DATETIME NOT NULL,

  
  PRIMARY KEY (`CommentID`),

  CONSTRAINT `CommentUser`
    FOREIGN key (`UserID`)
    REFERENCES `main_db`.`Users` (`UserID`),
  CONSTRAINT `CommentType`
    FOREIGN key (`TypeID`)
    REFERENCES `main_db`.`ItemTypes` (`TypeID`)
);
