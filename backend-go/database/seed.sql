-- USERS TABLE
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('admin', 'worker'))
);

-- PASSWORDS are bcrypt hashes:
-- admin123 => $2a$12$3vTeUuLkuOdZ9k1/X7KnTOjrlTzHH67Y/yyW7r0U9bSH17GhWg.ii
-- workerpass => $2a$12$Xs3yEu.VM97h0jQgfJTXHOw0Hb02alJSilKvDUvqFNDPNDohKYZiq

INSERT INTO users (username, password, name, role) VALUES
  ('admin', '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9', 'Admin User', 'admin'),
  ('worker1', '4ddff7855ff6e876b0c55f88023c2d23ce020906c648228eb771eb720f83c8f7', 'Worker One', 'worker');

-- SHIFTS TABLE
DROP TABLE IF EXISTS shifts;
CREATE TABLE shifts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL,
    start TIME NOT NULL,
    end TIME NOT NULL,
    role TEXT NOT NULL,
    location TEXT
);

INSERT INTO shifts (date, start, end, role, location) VALUES
  ('2025-05-20', '08:00:00', '16:00:00', 'Cashier', 'Jakarta'),
  ('2025-05-21', '08:00:00', '16:00:00', 'Cook', 'Bandung');

-- REQUESTS TABLE
DROP TABLE IF EXISTS requests;
CREATE TABLE requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    shift_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'approved', 'rejected')),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(shift_id) REFERENCES shifts(id)
);
