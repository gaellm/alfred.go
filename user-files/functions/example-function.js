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