CREATE TABLE IF NOT EXISTS users (
    email VARCHAR(255) PRIMARY KEY,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    question_id INT REFERENCES questions(id),
    text TEXT NOT NULL,
    user_email VARCHAR(255) REFERENCES users(email),
    UNIQUE (question_id, user_email)
);

CREATE TABLE IF NOT EXISTS signatures (
    id SERIAL PRIMARY KEY,
    user_email VARCHAR(255) REFERENCES users(email),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sign TEXT NOT NULL,
    UNIQUE (user_email)  
);


INSERT INTO questions (text)
SELECT q
FROM (VALUES
    ('Question 1'),
    ('Question 2'),
    ('Question 3'),
    ('Question 4'),
    ('Question 5'),
    ('Question 6'),
    ('Question 7'),
    ('Question 8'),
    ('Question 9'),
    ('Question 10')
) AS q(q)
WHERE NOT EXISTS (SELECT 1 FROM questions WHERE text = q);
