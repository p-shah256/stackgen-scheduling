// Authentication is handled by MongoDB container environment variables
db = db.getSiblingDB('scheduling-app');

// Create collections
db.createCollection('events');
db.createCollection('users');
db.createCollection('availability');

// Create indexes
db.events.createIndex({ "eventId": 1 }, { unique: true });
db.users.createIndex({ "userId": 1 }, { unique: true });
db.availability.createIndex({ "userId": 1, "eventId": 1 });
