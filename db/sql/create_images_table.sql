CREATE TABLE daily.content_images (
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    article_id INT,
    image_data_path VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_article_id
        FOREIGN KEY (article_id)
        REFERENCES articles (id)
);
