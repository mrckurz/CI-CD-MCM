# Architecture Documentation

## Request flow
**HTTP Request → Router → Handler → Store → Database**

+---------------------------------------------------------------------+
| HTTP REQUEST FLOW - Product Catalog API                             |
| Example: GET /products/42                                           |
+---------------------------------------------------------------------+

  CLIENT
    |  HTTP GET /products/42
    v
+---------------------------------------------------------------------+
| 1. HTTP SERVER          cmd/api/main.go                             |
|    main()                                            (line 14)      |
|    http.ListenAndServe(":"+port, r)                  (line 51)      |
|    PORT from env, default 8080                                      |
+---------------------------------------------------------------------+
    |
    v
+---------------------------------------------------------------------+
| 2. ROUTER               gorilla/mux                                 |
|    r := mux.NewRouter()              cmd/api/main.go:20             |
|    matches "/products/{id:[0-9]+}" with method GET                  |
|    routes registered by PostgresHandler.RegisterRoutes              |
|       internal/handler/postgres_handler.go:24                       |
|       (in-memory variant: internal/handler/handler.go:24)           |
+---------------------------------------------------------------------+
    |  http.ResponseWriter, *http.Request
    v
+---------------------------------------------------------------------+
| 3. HANDLER              PostgresHandler.GetProduct                  |
|    internal/handler/postgres_handler.go:50                          |
|      a) id, _ := strconv.Atoi(mux.Vars(r)["id"])                    |
|      b) p, err := h.Store.GetByID(id)                               |
|      c) err != nil  -> respondError(w, 404, "Product not found")    |
|         err == nil  -> respondJSON(w, 200, p)                       |
|    Other endpoints in same file:                                    |
|      Health, GetProducts, CreateProduct, UpdateProduct,             |
|      DeleteProduct                                                  |
+---------------------------------------------------------------------+
    |  Store.GetByID(id int) (model.Product, error)
    v
+---------------------------------------------------------------------+
| 4. STORE                PostgresStore.GetByID                       |
|    internal/store/postgres.go:67                                    |
|      s.DB.QueryRow(                                                 |
|        "SELECT id, name, price FROM products WHERE id = $1", id     |
|      ).Scan(&p.ID, &p.Name, &p.Price)                               |
|      sql.ErrNoRows -> store.ErrNotFound                             |
|    Store is selected at startup based on DB_HOST env var:           |
|      DB_HOST set    -> PostgresStore (postgres.go)                  |
|      DB_HOST unset  -> MemoryStore   (memory.go)                    |
+---------------------------------------------------------------------+
    |  parameterized SQL via database/sql ($1 placeholder)
    v
+---------------------------------------------------------------------+
| 5. DATABASE             PostgreSQL  (driver: github.com/lib/pq)     |
|    *sql.DB created in NewPostgresStore                              |
|       internal/store/postgres.go:17                                 |
|    schema ensured at startup by EnsureTable()                       |
|       internal/store/postgres.go:33                                 |
|       products(id SERIAL PK, name TEXT NOT NULL,                    |
|                price NUMERIC(10,2) NOT NULL DEFAULT 0)              |
+---------------------------------------------------------------------+
    |
    | rows / err  (response travels back up: DB -> Store -> Handler   |
    |              -> Router -> Client as JSON)                       |
    v

SQL queries used by PostgresStore (internal/store/postgres.go):
  GetAll   line 46   SELECT id, name, price FROM products ORDER BY id
  GetByID  line 69   SELECT id, name, price FROM products WHERE id=$1
  Create   line 80   INSERT INTO products (name, price)
                     VALUES ($1, $2) RETURNING id
  Update   line 88   UPDATE products SET name=$1, price=$2 WHERE id=$3
  Delete   line 105  DELETE FROM products WHERE id=$1

No middleware is registered on the router.

---

## Storage Implementations: MemoryStore vs PostgresStore

* MemoryStore (`internal/store/memory.go`): An in-memory, volatile key-value map protected by a mutex.
  * USAGE: Local development, rapid prototyping, or unit testing (as seen in `handler_test.go`) where data persistence between restarts is not necessary and spinning up a database is too much overhead.
  * TRADE-OFFS: Extremely fast and requires zero external dependencies to run. However, all data is lost upon application restart, it cannot share state across multiple load-balanced API instances (prevents horizontal scaling), and data capacity is limited by the server's available RAM.

* PostgresStore (`internal/store/postgres.go`): A persistent SQL database integration using `database/sql`.
  * USAGE: Production environments, integration testing (e.g., via Docker Compose), and whenever data needs to persist reliably across container/server restarts or scale out horizontally across multiple instances.
  * TRADE-OFFS: Provides durable, ACID-compliant storage that supports large datasets and allows multiple API instances to share the same state. On the downside, it requires managing a separate database service (more operational overhead), introduces network latency for queries, and involves schema management.