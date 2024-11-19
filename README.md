# hot-coffee

## Description

hot-coffee — this is a coffee shop management system developed in the Go language using a three-tier architecture. The application provides a RESTful API for managing orders, menu items, and inventory, as well as aggregating data for analytics. The data is stored locally in JSON files.

## Functionality

- **Order Management**: Create, receive, update and delete orders.
- **Menu Item Management**: Create, receive, update and delete menu items.
- **Inventory Management**: Adding, receiving, updating and deleting items in the inventory.
- **Data aggregation**: Data analysis, for example, total sales or popular menu items.
- **Logging**: using the `log/slog` package to log all events and errors.

## Installation

### Requirements

- Go 1.22 and above.

### Compilation

First clone the repository and compile the project:

```bash
git clone <repository_url>
cd ./hot-coffee
go build -o ./hot-coffee
```

## Launch

To start the server, use the following command:

```bash
./hot-coffee --port 8080 --dir data
```

### Command Line Options

- `-port N`: Sets the port on which the server will run.
- `-dir S`: Darectory for stroing data

## Example of use via Postman

### 1. **Create order**

**Method:** `POST`
**URL:** `http://localhost:8080/orders`

#### Steps:

1. Open the Postman.
2. Select the `POST` method.
3. Paste the URL: `http://localhost:8080/orders`.
4. Enter an order in the body of the request:

   ```json
   {
   	"customer_name": "John Doe",
   	"items": [
   		{
   			"product_id": "espresso",
   			"quantity": 2
   		},
   		{
   			"product_id": "croissant",
   			"quantity": 1
   		}
   	]
   }
   ```

5. Click `Send`.

**Response:**

- **Status:** `201 Created`
- **Body:**
  ```json
  {
  	"total_sales": 1500.5
  }
  ```

---

### 2. **Create menu item**

**Method:** `POST`
**URL:** `http://localhost:8080/menu`

#### Steps:

1. Open the Postman.
2. Select the `POST` method.
3. Paste the URL: `http://localhost:8080/menu`.
4. Enter an menu item in the body of the request:

   ```json
   {
   	"product_id": "espresso",
   	"name": "Espresso",
   	"description": "Strong and bold coffee",
   	"price": 2.5,
   	"ingredients": [
   		{
   			"ingredient_id": "espresso_shot",
   			"quantity": 1
   		}
   	]
   }
   ```

5. Click `Send`.

**Response:**

- **Status:** `201 Created`

---

### 3. **Create inventory**

**Method:** `POST`
**URL:** `http://localhost:8080/inventory`

#### Steps:

1. Open the Postman.
2. Select the `POST` method.
3. Paste the URL: `http://localhost:8080/inventory`.
4. Enter an inventory item in the body of the request:

   ```json
   {
   	"ingredient_id": "espresso_shot",
   	"name": "Espresso Shot",
   	"quantity": 500,
   	"unit": "shots"
   }
   ```

5. Click `Send`.

**Response:**

- **Status:** `201 Created`

---

## Logging

All actions and errors are logged using the log/log package. Logs are written to standard output (stdout).

## Errors and processing

In case of errors, the server will return a JSON response with an error message and the corresponding HTTP status.

```json
{
	"error": "not found order with this ID"
}
```

## Architecture

The project uses a three-level architecture:

### 1. Presentation Layer (Handlers)

This layer is responsible for processing HTTP requests and transferring data to the business logic layer. It also generates HTTP responses.

### 2. Business Logic Layer (Services)

The business logic layer contains the basic logic of the application, including data processing and the application of business rules.

### 3. Data access Layer (Repositories)

This layer is responsible for interacting with JSON files for reading and writing data.

## Project structure

```bash
hot-coffee/
├── cmd/
│ └── main.go
├── data/ 						 # Data storage in JSON format
│ ├── inventory.json
│ ├── menu_items.json
│ └── orders.json
├── internal/
| ├── customErrors/
│ │ └── customErrors.go
│ ├── dal/ 			             # Data access Layer (working with JSON)
│ │ ├── inventory_repository.go
│ │ ├── menu_repository.go
│ │ ├── order_repository.go
│ │ └── utils.go
│ ├── flags/
│ │ ├── checkFlags.go
│ │ └── flags.go
│ ├── handlers/ 				 # Presentation Layer
│ │ ├── inventory.go
| | ├── menu.go
│ │ ├── order.go
│ │ ├── reports.go
│ │ └── utils.go
│ ├── models/ 					 # Defining data structures
│ │ ├── inventory_item.go
│ │ ├── menu_item.go
│ │ ├── order.go
│ │ └── utils.go
│ ├── router/
│ │ ├── inventory_router.go
│ │ ├── menu_router.go
│ │ ├── order_router.go
│ │ ├── report_router.go
│ │ └── router.go
│ └── services/ 				 # Business Logic Layer
│   ├── inventory.go
│   ├── menu.go
|   ├── order.go
|   ├── report.go
│   └── utils.go
├── go.mod
├── hot-coffee
└── README.md

```

## Authors

Damir Usetov
