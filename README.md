# Daily Worker Roster Management System

## 🚀 Project Overview
This system simulates scheduling features used by companies with hourly/daily staff.
It supports worker self-service shift requests, admin shift management, and enforces real-world scheduling constraints.

Features:

    Employee (Worker) Interface

        View assigned shifts

        View available (unassigned) shifts

        Request to work on available shifts

        View status of shift requests (pending, approved, rejected)

    Admin Interface

        Create, edit, and delete shifts

        View and manage all shift requests from workers

        Approve/reject shift requests

        View finalized rosters (approved shift assignments)

    Authentication

        JWT-based login for both workers and admins

    Business Constraints

        No double-booking or overlapping shifts

        Max 1 shift per day, max 5 shifts per week per worker

        Admin can override assignments


## Database Schema
# users
| Field      | Type       | Notes                      |
| ---------- | ---------- | -------------------------- |
| `id`       | INTEGER PK |                            |
| `username` | TEXT       | UNIQUE                     |
| `password` | TEXT       | Hashed (bcrypt)            |
| `name`     | TEXT       | Display name               |
| `role`     | TEXT       | e.g. `"worker"`, `"admin"` |

# shifts
| Field      | Type       | Notes                  |
| ---------- | ---------- | ---------------------- |
| `id`       | INTEGER PK |                        |
| `date`     | DATE       | Stored as `YYYY-MM-DD` |
| `start`    | TIME       | Stored as `HH:MM:SS`   |
| `end`      | TIME       | Stored as `HH:MM:SS`   |
| `role`     | TEXT       | Job role               |
| `location` | TEXT       | Optional               |

# requests
| Field      | Type       | Notes                                   |
| ---------- | ---------- | --------------------------------------- |
| `id`       | INTEGER PK |                                         |
| `user_id`  | INTEGER    | FK to `users.id`                        |
| `shift_id` | INTEGER    | FK to `shifts.id`                       |
| `status`   | TEXT       | `'pending'`, `'approved'`, `'rejected'` |

When you run the dockerfile i added some pre-seed data to be tested


## How to Run
1. Clone the repository

```bash
git clone https://github.com/yourusername/roster-management.git
cd roster-management
```

2. Build and Start All Services

```bash
docker-compose up --build
```

    The backend (Go API) will run on host:8080

    The frontend (Svelte) will run on host:3000

3. Access the App

    Open host:3000 in your browser.
