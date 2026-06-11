-- Runs once, on first boot of the postgres container (empty data volume).
-- Each infra service gets its own database inside the single dev instance.
CREATE DATABASE kratos;
CREATE DATABASE hydra;
CREATE DATABASE spicedb;
