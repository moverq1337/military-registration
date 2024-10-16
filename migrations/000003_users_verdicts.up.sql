CREATE TABLE IF NOT EXISTS users_verdicts
(
    user_id    VARCHAR(255) NOT NULL,
    doctor_id  VARCHAR(255) NOT NULL,
    verdict    BOOLEAN      NOT NULL,
    verdict_date TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, doctor_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (doctor_id) REFERENCES doctors (id)
);
