const mongoHost = process.env.MEDICINE_API_MONGODB_HOST
const mongoPort = process.env.MEDICINE_API_MONGODB_PORT

const mongoUser = process.env.MEDICINE_API_MONGODB_USERNAME
const mongoPassword = process.env.MEDICINE_API_MONGODB_PASSWORD

const database = process.env.MEDICINE_API_MONGODB_DATABASE
const collection = process.env.MEDICINE_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection, "status")) {
        print(`Collections '${collection}' and status already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
if (collections.includes(collection)) {
    db.createCollection(collection)

    // create indexes
    db[collection].createIndex({"id": 1})

    //insert sample data
    let result = db[collection].insertMany([
        {
            "id": "bobulova",
            "name": "Dr.Bobulová",
            "roomNumber": "123"
        }
    ]);

    if (result.writeError) {
        console.error(result)
        print(`Error when writing the data: ${result.errmsg}`)
    }
}
if (collections.includes("status")) {
    db.createCollection("status")

    // create indexes
    db["status"].createIndex({"id": 1})

    //insert sample data
    let result = db["status"].insertMany([
        {
            "id": 1,
            "value": "To_ship",
            "ValidTransitions": [2, 4]
        },
        {
            "id": 2,
            "value": "Shipped",
            "ValidTransitions": [3, 4]
        },
        {
            "id": 3,
            "value": "Delivered",
            "ValidTransitions": []
        },
        {
            "id": 4,
            "value": "Canceled",
            "ValidTransitions": []
        }
    ]);

    if (result.writeError) {
        console.error(result)
        print(`Error when writing the data: ${result.errmsg}`)
    }
}

// exit with success
process.exit(0);