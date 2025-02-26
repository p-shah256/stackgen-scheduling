// Authentication is handled by MongoDB container environment variables
db = db.getSiblingDB('admin');
db.auth(process.env.MONGO_INITDB_ROOT_USERNAME, process.env.MONGO_INITDB_ROOT_PASSWORD);

db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE);

db.createUser({
  user: process.env.MONGO_INITDB_ROOT_USERNAME,
  pwd: process.env.MONGO_INITDB_ROOT_PASSWORD,
  roles: ["readWrite", "dbAdmin"]
});

// Create collections
db.createCollection('events');
db.createCollection('users');
db.createCollection('availability');

// Create indexes
db.events.createIndex({ "eventId": 1 }, { unique: true });
db.users.createIndex({ "userId": 1 }, { unique: true });
db.availability.createIndex({ "userId": 1, "eventId": 1 });
