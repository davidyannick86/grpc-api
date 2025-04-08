# SQL / NoSQL

| Feature                | SQL                                                          | NoSQL                                                                        |
| :--------------------- | :----------------------------------------------------------- | :--------------------------------------------------------------------------- |
| **Database**           | Collection of tables                                         | Collection of collections                                                    |
| **Table vs. Collection** | A structured set of data organized into rows and columns     | Similar to a table in SQL. It stores documents, which are JSON-like objects. |
| **Row vs. Document**   | Represents a single record in a table                        | Equivalent of a row in SQL                                                   |
| **Column vs. Field**   | Represents a single attribute or field in a table            | Similar to a column in SQL. It stores data in a key-value pair format        |
| **Schema**             | Defines the structure of the table                           | Schema-less                                                                  |
| **Primary Key vs. _id** | A unique identifier for each row in a table                  | `_id` field serves as the primary key for a document in a collection         |
| **Indexes**            | Used to speed up queries                                     | Used to speed up queries in similar way                                      |

Voici le deuxième tableau transformé en Markdown :

| Feature                                              | SQL                                                              | NoSQL                                                                                              |
| :--------------------------------------------------- | :--------------------------------------------------------------- | :------------------------------------------------------------------------------------------------- |
| **Joins vs. Embedding/Referencing**                   | To combine data from multiple tables                             | Does not support joins in the traditional SQL sense. Embedding is supported.                      |
| **SQL Queries vs. MongoDB Queries**                  | Uses structured query language (SQL)                             | Uses its own query language                                                                        |
| **Transactions**                                     | Allow you to execute multiple operations atomically              | Also supports transactions                                                                         |
| **Aggregation**                                      | Provides aggregate functions like `SUM()`, `COUNT()`, `AVG()`      | Aggregation framework allows you to perform complex data transformations and calculations          |
| **Foreign Keys vs. Manual Referencing: (Table Relations)** | Enforce relationships between tables, ensuring referential integrity | You manually reference related documents using fields, but there's no enforcement of referential integrity |
| **ACID Compliance**                                  | Typically ACID-compliant                                         | Also ACID-compliant                                                                                |
