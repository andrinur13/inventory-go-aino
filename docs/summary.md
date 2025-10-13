# TWC OTA API Documentation

## Overview

The TWC OTA (Online Travel Agent) API is a RESTful service designed to manage a tourism ticket booking system. The system provides functionalities for managing tourism destinations, ticket bookings, user management, agent management, and transaction processing.

## System Architecture

The system is built using Go (Golang) with the Gin web framework and follows a clean architecture pattern, consisting of the following components:

### Core Components

1. **API Layer**:
   - Handles HTTP requests and responses
   - Implements RESTful endpoints and Websocket connections
   - Uses JWT authentication

2. **Middleware Layer**:
   - Authentication and authorization
   - Request timeouts
   - Permission management

3. **Service Layer**:
   - Business logic implementation
   - Caching service for improved performance

4. **Data Access Layer**:
   - GORM-based database access
   - Repository pattern for data operations
   - Database connection management with retry mechanisms

5. **Configuration System**:
   - Environment-based configuration
   - Support for multiple environments (dev, staging, production)

6. **Logging System**:
   - Structured logging with logrus
   - Integration with Elastic APM for monitoring

### Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Caching**: In-memory cache
- **Documentation**: Swagger
- **Monitoring**: Elastic APM
- **Deployment**: Docker (containerized)

## Database Structure

The database schema is organized around several core entities:

### Core Entities

1. **Users and Agents**: Handles authentication, permissions, and agent management
2. **Master Data**: Manages tourism sites, tickets, tariffs, and pricing information
3. **Booking System**: Processes bookings, transactions, and payment handling
4. **Trip Planning**: Manages multi-day trips with destinations and itineraries
5. **Inventory Management**: Tracks available tickets and QR code management
6. **User Preferences**: Stores user favorites and notifications

For detailed database schema, refer to the [ticketing.sql](../db/ticketing.sql) file.

## API Endpoints

The API is organized around the following core areas:

### Authentication

- `POST /auth/login`: Authenticate a user and obtain JWT token
- `POST /auth/register`: Register a new user

### Ticket Management

- `POST /api/ticket/list`: Get available tickets
- `POST /api/ticket/booking`: Book tickets
- `POST /api/ticket/redeem`: Redeem a ticket
- `POST /api/ticket/checkout`: Checkout and complete a booking
- `GET /api/ticket/cluster`: Get tourism clusters

### Tourism Sites

- `GET /api/site/detail`: Get detailed information about a tourism site
- `GET /api/site/extras`: Get additional information about a tourism site

### Agent Management

- `POST /api/register/agent`: Register a new agent

### Transaction Management

- `GET /api/trx/list/:page/:size`: List transactions with pagination
- `POST /api/trx/create`: Create a new transaction
- `PUT /api/trx/update`: Update transaction details
- `PUT /api/trx/pay`: Update transaction payment status
- `POST /api/trx/info`: Get transaction information by QR code
- `POST /api/trx/number`: Get transaction details by transaction number

### User Favorites

- `POST /api/fav/create`: Create a user favorite
- `POST /api/fav/delete`: Delete a user favorite
- `GET /api/fav/list`: List user favorites
- `POST /api/fav/image`: Upload an image for a favorite

### Application Configuration

- `GET /api/appconfig/detail`: Get application configuration

### WebSocket

- `GET /ws/:name`: WebSocket endpoint for real-time communication

## Authentication and Authorization

The system uses JWT (JSON Web Tokens) for authentication. The token must be provided in the Authorization header of requests to protected endpoints:

```http
Authorization: Bearer <token>
```

### Permission System

The system implements a fine-grained permission system. Each endpoint is protected with specific permission requirements. The middleware checks if the user has the necessary permissions to access the requested resource.

Example permissions include:

- `PERMISSION_MASTER_USER_VIEW`
- `PERMISSION_MASTER_USER_SAVE`

## Error Handling

The API uses standardized error responses with the following structure:

```json
{
  "status": false,
  "message": "Error message",
  "code": "01",
  "data": null
}
```

Common error codes:

- `00`: Success
- `01`: General error
- `99`: Authentication/authorization error

## Deployment

The application can be deployed using Docker and Docker Compose. The included configuration files support easy deployment in various environments.

### Environment Setup

1. Clone the repository
2. Configure the appropriate config file (config.dev.json, config.prod.json, etc.)
3. Start the application:

   ```bash
   docker-compose up -d
   ```

### Configuration Files

Environment-specific configuration files are located in the `config/` directory:

- `config.dev.json`: Development environment
- `config.stg.json`: Staging environment
- `config.prod.json`: Production environment
- `config.local.json`: Local development

## Monitoring and Logging

### Elastic APM

The application integrates with Elastic APM for performance monitoring and tracking:

- Server URL: `https://apm.ainosi.com`
- Service Name: twc-api-ota

### Logging

The application uses structured logging with the following information:

- Message
- Status code
- Success/failure status
- Request/response data

## Development Guidelines

### Code Structure

```plaintext
twc-ota-api/
├── api/                  # API handlers
│   ├── master/           # Master data endpoints
│   ├── public/           # Public endpoints
│   ├── api.router.go     # Route definitions
│   └── websocket.handler.go # WebSocket handlers
├── config/               # Configuration
├── db/                   # Database
│   ├── entities/         # Database entity models
│   └── repositories/     # Data access repositories
├── docs/                 # Documentation
├── logger/               # Logging utilities
├── middleware/           # HTTP middleware
├── requests/             # Request models
├── service/              # Business logic services
├── storage/              # File storage
├── utils/                # Utility functions
├── Dockerfile            # Docker configuration
├── docker-compose.yaml   # Docker Compose configuration
├── go.mod                # Go modules
└── main.go               # Application entry point
```

### Adding New Features

1. Define entities in `db/entities/`
2. Create repository functions in `db/repositories/`
3. Add API handlers in appropriate directories under `api/`
4. Register routes in `api/api.router.go`
5. Add permissions in the route registration

### Testing

Ensure all new features include appropriate testing and that existing functionality is not broken.

## API Response Format

Successful responses follow this structure:

```json
{
  "status": true,
  "message": "Success message",
  "code": "00",
  "data": {
    // Response data
  }
}
```

List responses (with pagination) follow this structure:

```json
{
  "status": true,
  "message": "Success message",
  "code": "00",
  "total_data": 100,
  "current_data": 10,
  "total_pages": 10,
  "page": 1,
  "size": 10,
  "data": [
    // List of items
  ]
}
```

## License and Credits

This application is developed by Ainosi and is proprietary software.

## Contact

For technical support or inquiries, please contact:

- Email: `rino@ainosi.co.id`
