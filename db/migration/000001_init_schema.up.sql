CREATE TABLE user (
   id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
   username VARCHAR(50) NOT NULL UNIQUE KEY,
   password VARCHAR(255) NOT NULL
);

CREATE TABLE user_details (
   user_id INT(4) NOT NULL,
   first_name VARCHAR(50) NOT NULL,
   last_name VARCHAR(50),
   dob DATETIME,
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (user_id) REFERENCES user(id) 
);

CREATE TABLE role (
   id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
   role VARCHAR(25) NOT NULL
);

CREATE TABLE user_role (
   role_id INT(4) NOT NULL,
   user_id INT(4) NOT NULL,
   FOREIGN KEY (role_id) REFERENCES role(id),
   FOREIGN KEY (user_id) REFERENCES user(id)
);