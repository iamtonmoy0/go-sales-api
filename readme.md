# Sales API Project
## Overview

This project is an API designed for managing sales-related operations. It provides endpoints to handle various aspects of sales such as cashiers, categories, discounts, orders, payments, payment types, and products.

### Features

- Cashier Management: Allows the creation, retrieval, updating, and deletion of 
cashiers. Cashiers can be created with a name, email address, phone number, and password. They can also be retrieved by their cashiers. Cashiers have a name cashiers. 
- Category Management: Provides functionalities to manage product categories.
- Discount Management: Enables the management of discounts on products.
- Order Management: Facilitates the creation and tracking of orders.
- Payment Management: Handles different types of payments and payment methods.
Product Management: Manages product details including stock, price, and categories.
Technologies Used
- Golang: The backend of the API is implemented using Go programming language.
- Gorm: Gorm is used as the ORM library to interact with the database.
- JSON Web Tokens (JWT): JWT is used for authentication and authorization purposes.
- MySQL: MySQL database is used to store data related to cashiers, categories, discounts, orders, payments, payment types, and products.
## Setup
Clone the repository.
Install dependencies using
```sh
 go mod tidy
```
Configure the MySQL database connection in the .env file.
Run the application using go 
```sh
run main.go
```.
