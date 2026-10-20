CREATE TABLE rates (
                       currency   VARCHAR(10) NOT NULL REFERENCES currencies(currency) ON DELETE CASCADE,
                       price      DOUBLE PRECISION NOT NULL CHECK (price >= 0),
                       fetched_at TIMESTAMPTZ NOT NULL,
                       PRIMARY KEY (currency, fetched_at)
);