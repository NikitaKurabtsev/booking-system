# Resource Booking Management System

## Project Overview

**Idea**: Develop a resource booking management system to help users discover, reserve, and efficiently utilize limited resources. Examples include booking conference rooms in offices, event equipment, sports facilities, or coworking spaces.  
**Hypothesis**: Automating the resource booking process will reduce conflicts from overlapping reservations, lower administrative costs, and improve user satisfaction.

---

## Target Audience
- Offices and companies with large teams.
- Owners of sports facilities or other resources requiring reservations.
- Coworking space administrators.

---

## Use Case Scenario
1. **User Authentication**: A user logs into the system.
2. **Resource Selection**: They select a resource type (room, equipment, facility).
3. **Search & Filter**: Browse available resources via a calendar or filters.
4. **Reservation**: Book the resource for a specific time slot.
5. **Confirmation**: Receive a notification (email/push) confirming the booking.

---

## MVP (Minimal Viable Product)

### Core Features
- **User Registration & Authentication**: Secure sign-up/login functionality.
- **Resource Management**: View available resources and their statuses.
- **Booking System**: Reserve resources via an API interface.
- **Booking Management**: View and cancel personal reservations.
- **Notifications**: Email/push notifications for booking status updates.
- **Feedback Collection**: In-app form for user feedback and interviews.

### Feedback Strategy
- Test MVP with small offices/coworking spaces.
- Collect feedback via embedded forms and user interviews.

---

## Technology Stack

| **Category**       | **Technologies**                                                                 |
|---------------------|---------------------------------------------------------------------------------|
| **Backend**         | Golang                                                                          |
| **Database**        | PostgreSQL (resource/booking storage)                                           |
| **Messaging**       | Kafka (notification processing)                                                 |
| **Caching**         | Redis (fast access to frequently queried data, e.g., resource statuses)         |
| **API**             | REST API (public) / gRPC (internal service communication)                       |
| **Infrastructure**  | Docker (containerization), Kubernetes (orchestration)                           |
| **Monitoring**      | Prometheus, Grafana, Loki (logging and observability)                           |

---

## Project Structure
1. **Handlers (API Layer)**
    - Process incoming requests.
    - Validate input data.

2. **Services (Business Logic)**
    - Check resource availability.
    - Enforce booking rules.
    - Trigger notifications.

3. **Repositories (Data Layer)**
    - Manage database read/write operations.
    - Cache frequent queries via Redis.

---

## API Endpoints

| **Method** | **Endpoint**               | **Description**                              |
|------------|----------------------------|----------------------------------------------|
| `POST`     | `/api/register`            | Register a new user.                         |
| `POST`     | `/api/login`               | Authenticate a user.                         |
| `GET`      | `/api/resources`           | List all available resources.                |
| `POST`     | `/api/bookings`            | Create a new booking.                        |
| `GET`      | `/api/bookings`            | Retrieve a user's active bookings.           |
| `DELETE`   | `/api/bookings/{id}`       | Cancel a booking by ID.                      |

---

## Database Models

### `Resource`
```json
{
  "id": "UUID",
  "name": "string",
  "type": "enum (room, equipment, facility)",
  "status": "enum (available, booked)"
}