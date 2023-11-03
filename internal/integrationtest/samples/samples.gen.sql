-- Code generated by ddlgen. DO NOT EDIT.
--
-- Date: 2023-11-03T20:38:32+09:00
--

-- source: samples/success_SuccessParen.go:8
-- SuccessParen
--
-- spanddl:   table: CREATE TABLE `SuccessParens`
-- spanddl: options: PRIMARY KEY (`Id`)
CREATE TABLE `SuccessParens` (
    -- ID is a primary key of SuccessParen
    -- ID requires uuid
    `Id`          STRING(36)   NOT NULL,
    -- Description is a description of SuccessParen
    `Description` STRING(1024) NOT NULL
) PRIMARY KEY (`Id`);

-- source: samples/success_SuccessParen.go:21
-- SuccessParenTwo
--
-- spanddl:      table: CREATE TABLE `SuccessParenTwos`
-- spanddl: constraint: FOREIGN KEY (`SuccessParenId`) REFERENCES SuccessParens(`Id`)
-- spanddl:    options: PRIMARY KEY (`Id`)
CREATE TABLE `SuccessParenTwos` (
    -- ID is a primary key of SuccessParenTwo
    -- ID requires uuid
    `Id`             STRING(36)   NOT NULL,
    -- SuccessParenID is a foreign key to SuccessParen.ID
    `SuccessParenId` STRING(36)   NOT NULL,
    -- Description is a description of SuccessParenTwo
    `Description`    STRING(1024) NOT NULL,
    FOREIGN KEY (`SuccessParenId`) REFERENCES SuccessParens(`Id`)
) PRIMARY KEY (`Id`);

-- source: samples/success_SuccessParen.go:34
-- SuccessParenThree
--
-- spanddl: options: PRIMARY KEY (`Id`)
CREATE TABLE SuccessParenThree (
    -- ID is a primary key of SuccessParenThree
    -- ID requires uuid
    `Id`          STRING(36)   NOT NULL,
    -- Description is a description of SuccessParenThree
    `Description` STRING(1024) NOT NULL
) PRIMARY KEY (`Id`);

-- source: samples/success_SuccessParen.go:45
-- spanddl:      table: `SuccessParenFour`
-- spanddl: constraint: FOREIGN KEY (`SuccessParenId`) REFERENCES SuccessParens(`Id`)
-- spanddl:    options: PRIMARY KEY (`Id`)
CREATE TABLE `SuccessParenFour` (
    -- ID is a primary key of SuccessParenFour
    -- ID requires uuid
    `Id`             STRING(36)   NOT NULL,
    -- SuccessParenID is a foreign key to SuccessParen.ID
    `SuccessParenId` STRING(36)   NOT NULL,
    -- Description is a description of SuccessParenFour
    `Description`    STRING(1024) NOT NULL,
    FOREIGN KEY (`SuccessParenId`) REFERENCES SuccessParens(`Id`)
) PRIMARY KEY (`Id`);

-- source: samples/success_SuccessRoot.go:9
-- SuccessRoot
--
-- spanddl:   table: CREATE TABLE success_roots
-- spanddl: options: PRIMARY KEY (`Id`)
CREATE TABLE success_roots (
    -- ID is a primary key of SuccessRoot
    -- ID requires uuid
    `Id`          STRING(36)   NOT NULL,
    -- Description is a description of SuccessRoot
    `Description` STRING(1024) NOT NULL
) PRIMARY KEY (`Id`);

-- source: samples/success_SuccessRoot.go:22
-- SuccessRootTwo
--
-- spanddl:      table: CREATE TABLE success_root_twos
-- spanddl: constraint: FOREIGN KEY (SuccessRootId) REFERENCES success_roots(`Id`)
-- spanddl:    options: PRIMARY KEY (`Id`)
CREATE TABLE success_root_twos (
    -- ID is a primary key of SuccessRootTwo
    -- ID requires uuid
    `Id`            STRING(36)   NOT NULL,
    -- SuccessRootID is a foreign key to SuccessRoot.ID
    `SuccessRootId` STRING(36)   NOT NULL,
    -- Description is a description of SuccessRootTwo
    `Description`   STRING(1024) NOT NULL,
    FOREIGN KEY (SuccessRootId) REFERENCES success_roots(`Id`)
) PRIMARY KEY (`Id`);

-- source: samples/success_SuccessUnformattedFile.go:6
-- SuccessUnformattedFile
-- spanddl:table:success_unformatted_files
CREATE TABLE success_unformatted_files (
    --  ID is a primary key of SuccessUnformattedFile
    `Id` STRING(36) NOT NULL
);
