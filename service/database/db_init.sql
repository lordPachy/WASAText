-- Cleaning up the database
PRAGMA writable_schema = 1;
DELETE FROM sqlite_master;
PRAGMA writable_schema = 0;
VACUUM;
PRAGMA integrity_check;

-- Creating the actual schema
CREATE TABLE users (id TEXT PRIMARY KEY, username TEXT UNIQUE NOT NULL, propic TEXT);
CREATE TABLE privchats (id INTEGER PRIMARY KEY, member1 TEXT NOT NULL, member2 TEXT NOT NULL, FOREIGN KEY (member1) REFERENCES users(username) ON UPDATE CASCADE, FOREIGN KEY (member2) REFERENCES users(username) ON UPDATE CASCADE);
CREATE TABLE privmessages (id INTEGER, messageID INTEGER, PRIMARY KEY (id, messageID), FOREIGN KEY (id) REFERENCES privchats(id), FOREIGN KEY (messageID) REFERENCES messages(id) ON DELETE CASCADE);
CREATE TABLE groupchats (id INTEGER PRIMARY KEY, groupname TEXT NOT NULL, groupphoto TEXT);
CREATE TABLE groupmembers (id INTEGER, member TEXT, PRIMARY KEY (id, member), FOREIGN KEY (member) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE, FOREIGN KEY (id) REFERENCES groupchats(id) ON DELETE CASCADE);
CREATE TABLE groupmessages (id INTEGER, messageID INTEGER, PRIMARY KEY (id, messageID), FOREIGN KEY (id) REFERENCES groupchats(id) ON DELETE CASCADE, FOREIGN KEY (messageID) REFERENCES messages(id) ON DELETE CASCADE);
CREATE TABLE messages (id INTEGER PRIMARY KEY, sender TEXT, created_at TIMESTAMP NOT NULL, content TEXT, photo TEXT, checkmarks INTEGER NOT NULL, replying_to INTEGER, og_sender TEXT, FOREIGN KEY (sender) REFERENCES users(username) ON UPDATE CASCADE, FOREIGN KEY (replying_to) REFERENCES messages(id) ON DELETE SET NULL, FOREIGN KEY (og_sender) REFERENCES users(username) ON UPDATE CASCADE);
CREATE TABLE messagecomments (id INTEGER, commentID INTEGER, PRIMARY KEY(id, commentID), FOREIGN KEY (id) REFERENCES messages(id) ON DELETE CASCADE, FOREIGN KEY (commentID) REFERENCES comments(id) ON DELETE CASCADE);
CREATE TABLE comments (id INTEGER PRIMARY KEY, sender TEXT, reaction TEXT NOT NULL, FOREIGN KEY (sender) REFERENCES users(username) ON UPDATE CASCADE);
CREATE TABLE groupmessageschecks (groupID INTEGER, messageID INTEGER, member TEXT, checkmarks INTEGER NOT NULL, PRIMARY KEY (groupID, messageID, member), FOREIGN KEY (groupID, member) REFERENCES groupmembers(id, member) ON UPDATE CASCADE ON DELETE CASCADE, FOREIGN KEY (groupID, messageID) REFERENCES groupmessages(id, messageID) ON DELETE CASCADE); 

