CREATE TABLE users (
    id       SERIAL       PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email    VARCHAR(255) NOT NULL UNIQUE,
    password TEXT         NOT NULL
);

CREATE TABLE resources (
    id     SERIAL       PRIMARY KEY,
    name   VARCHAR(255) NOT NULL UNIQUE,
    type   VARCHAR(50)  NOT NULL UNIQUE,
    status VARCHAR(10)  NOT NULL
        CHECK ( status IN ('available', 'booked'))
);

CREATE TABLE bookings (
    id          SERIAL    PRIMARY KEY,
    resource_id INT       NOT NULL,
    user_id     INT       NOT NULL,
    start_time  TIMESTAMP NOT NULL,
    end_time    TIMESTAMP NOT NULL,

    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id)     REFERENCES users(id)     ON DELETE CASCADE,

    CHECK ( start_time < end_time )
);
