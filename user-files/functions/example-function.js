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

    res.body += " - edited body";

    return res;

}