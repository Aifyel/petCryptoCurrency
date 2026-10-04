CREATE TABLE rates (
                       currency   VARCHAR(10) NOT NULL,
                       price      DOUBLE PRECISION NOT NULL CONSTRAINT < 0,
                       fetched_at TIMESTAMPTZ NOT NULL,
                       PRIMARY KEY (currency, fetched_at)
);