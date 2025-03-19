start transaction;

CREATE TABLE users (
    id VARCHAR(100) PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    user_name   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE writer (
    id  VARCHAR(100) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    user_name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_Verified BOOLEAN NOT NULL ,
    is_Active BOOLEAN NOT NULL,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE modules (
    id VARCHAR(200) NOT NULL  PRIMARY KEY,
    user_id VARCHAR(200) NOT NULL,
    module_name VARCHAR(500) NOT NULL,
    module_type varchar(255) NOT NULL,
    module_desc varchar(255),
    module_short_name VARCHAR(255),
    module_price VARCHAR(255),
    purchased boolean,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) 

);

CREATE TABLE blogs (
    id VARCHAR(200) NOT NULL  PRIMARY KEY,
    blog_name VARCHAR(500) NOT NULL,
    blog_title VARCHAR(500) NOT NULL,
    blog_content TEXT NOT NULL,
    module_id VARCHAR(200) NOT NULL,
    writer_id VARCHAR(200) NOT NULL,
    writer_name VARCHAR(500),
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_module FOREIGN KEY (module_id) REFERENCES modules(id) ,
    CONSTRAINT fk_writer FOREIGN KEY (writer_id) REFERENCES writer(id) 
);

commit;