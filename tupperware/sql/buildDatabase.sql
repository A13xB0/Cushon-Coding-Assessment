-- Create Customers table
CREATE TABLE Customers (
  customer_id INT PRIMARY KEY,
  customer_name VARCHAR(255),
  customer_email VARCHAR(255),
  customer_password VARCHAR(255),
  customer_type VARCHAR(255)
);

-- Create Funds table
CREATE TABLE Funds (
  fund_id INT PRIMARY KEY,
  fund_name VARCHAR(255)
);

-- Create Investments table
CREATE TABLE Investments (
  investment_id INT PRIMARY KEY,
  customer_id INT,
  fund_id INT,
  amount_invested DECIMAL(10, 2),
  date_invested DATE,
  FOREIGN KEY (customer_id) REFERENCES Customers(customer_id),
  FOREIGN KEY (fund_id) REFERENCES Funds(fund_id)
);

-- Microservices User accounts

CREATE USER auth_service IDENTIFIED BY 'password';
GRANT SELECT ON Customers TO auth_service;

CREATE USER customer_service IDENTIFIED BY 'password';
GRANT SELECT (customer_id, customer_name, customer_email, customer_password), UPDATE (email, password) ON Customers TO customer_service;

CREATE USER investment_service IDENTIFIED BY 'password';
GRANT SELECT, INSERT ON Investments TO investment_service;
GRANT SELECT ON Funds TO investment_service;
