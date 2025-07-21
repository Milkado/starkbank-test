-- +migrate Up
CREATE TABLE invoice (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
        amount FLOAT,
        tax_id VARCHAR(16),
        due DATE,
        expiration BIGINT,
        fine FLOAT,
        interest FLOAT,
        fee FLOAT,
        status CHAR,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
-- +migrate Down
DROP TABLE invoice;
	
