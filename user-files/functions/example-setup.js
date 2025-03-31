function setup() {
   // Provision the database with some key-value pairs
   dbSet("key1", "value1");
   dbSet("key2", "value2");

   // Load the JSON file into the database using batch writes
   const result = dbLoadFile("./user-files/functions/db-import.json");
   console.log(result); // Output: File loaded successfully using batch writes

   // Delete entry
   console.log("Going to delete key1 with value: " + dbGet("key1"))
   dbDelete("key1", "value1");
   console.log("Key1 is now " + dbGet("key1"))

   console.log("Database initialized and provisioned");
}