CREATE DATABASE test;

CREATE TABLE test.Article (
  ArticleName VARCHAR(100) NOT NULL,
  ArticleTag VARCHAR(100) NOT NULL,
  ArticleAuthor VARCHAR(45) NOT NULL,
  AricleContent VARCHAR(500) NOT NULL,
  ArticleReview VARCHAR(500) NOT NULL,
  PRIMARY KEY (ArticleName));

CREATE TABLE test.User (
  Username VARCHAR(100) NOT NULL,
  Password VARCHAR(100) NOT NULL,
  PRIMARY KEY (Username));