CREATE TABLE daily.articles (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(50) NOT NULL,
    content VARCHAR(2000) NOT NULL,
    thumbnail_path VARCHAR(50),
    tags_1 VARCHAR(20),
    tags_2 VARCHAR(20),
    tags_3 VARCHAR(20),
    tags_4 VARCHAR(20),
    tags_5 VARCHAR(20), 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
);
