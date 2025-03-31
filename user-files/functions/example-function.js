//using npm modules
//npm install hello-world-npm  --prefix /tmp/test
//const helloWorld = require('/tmp/test/node_modules/hello-world-npm');
//console.log(helloWorld());

function updateHelpers(helpers) {

    helpers.forEach( (helper) => {

        if (helper.name === "my-city"){

            console.log("Update city from " + helper.value + " to Gotham");
            helper.value = "Gotham";
        }
    });

    return helpers;
}


function alfred(mock, helpers, req, res) {

    console.log(JSON.stringify(mock));
    console.log(JSON.stringify(helpers));
    console.log(JSON.stringify(req));
    console.log(JSON.stringify(res));

    // Retrieve a value from the database
    const value = dbGet("key2");

    //if value in database
    if (value){
        res.body += " " + value +" - edited body with value from database";
    } else {
        res.body += " - edited body";
    }
    
    return res;

}