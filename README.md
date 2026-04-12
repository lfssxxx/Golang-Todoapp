# TodoApp

A simple backend service for managing todo tasks, written in Go.

## 🚀 Tech Stack

* Go
* PostgreSQL
* [Go Migrate](https://github.com/golang-migrate/migrate)
* Docker & Docker Compose
* Makefile
* Git

---


## 🔧 Setup & Run

### 1. Clone the repository

```
git clone <repo-url>
cd todoapp
```

---

### 2. Configure environment variables

Create `.env` from the example:

```
cp .env.example .env
```

---

### 3. Environment control (PostgreSQL)

```
make env-up
```
_or_
```
make env-down
```

---

### 4. Run migrations

```
make migrate-action action=up/down
```

---

### 5. Run the application

```
make todoapp-run
```

---

## 🧹 Cleanup

Remove all containers and database data:

```
make env-cleanup
```

or:

```
docker compose down -v
```

---


## 📌 Notes

* Database data is stored in a Docker volume
* Migrations are executed via a separate container
* Recommended to run using Linux environment

