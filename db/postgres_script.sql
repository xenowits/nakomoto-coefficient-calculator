-- brew install postgresql

-- brew services start postgresql

-- psql postgres

-- CREATE ROLE xenowits WITH LOGIN PASSWORD ‘password’;

-- ALTER ROLE xenowits CREATEDB;

\c postgres;

CREATE TABLE naka_coefficients (
id SERIAL PRIMARY KEY,
chain_name VARCHAR(50) NOT NULL,
chain_token VARCHAR(50) NOT NULL,
naka_co_prev_val SMALLINT NOT NULL,
naka_co_curr_val SMALLINT NOT NULL
);

\dt;
\d naka_coefficients;

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Cosmos', 'ATOM', -1, 7);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Polygon', 'MATIC', -1, 2);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Binance', 'BNB', -1, 7);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Osmosis', 'OSMO', -1, 4);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Mina', 'MINA', -1, 11);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Solana', 'SOL', -1, 19);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Avalanche', 'AVAX', -1, 24);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Terra', 'LUNA', -1, 7);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Graph', 'GRT', -1, 3);

INSERT INTO naka_coefficients (chain_name, chain_token, naka_co_prev_val, naka_co_curr_val) VALUES ('Thorchain', 'RUNE', -1, 10);

select * from naka_coefficients;

\q