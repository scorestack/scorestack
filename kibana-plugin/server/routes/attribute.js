export default function (server, dataCluster) {
    server.route({
        path: '/api/scorestack/attribute',
        method: 'GET',
        handler: async (req, h) => {
            let checks = {};

            // Get all attribute indexes
            let attribIndices = await dataCluster.callWithRequest(req, 'indices.get', {
                index: 'attrib_*',
                expand_wildcards: 'open',
            });

            // Get attributes for each check
            for (let attribIndex of Object.keys(attribIndices)) {
                let checkID = attribIndex.split("_").slice(2).join("_");
                let attribDoc = await dataCluster.callWithRequest(req, 'get', {
                    id: 'attributes',
                    index: attribIndex,
                });
                if (checkID in checks === false) {
                    checks[checkID] = {
                        "attributes": {},
                    };
                }
                if (!"attributes" in checks[checkID]) {
                    checks[checkID] = {};
                }
                checks[checkID]["attributes"] = Object.assign(checks[checkID]["attributes"], attribDoc._source);
            }

            // Get names for each check
            for (let checkID of Object.keys(checks)) {
                let checkDoc = await dataCluster.callWithRequest(req, 'get', {
                    id: checkID,
                    index: 'checks',
                    _source_includes: 'name',
                })
                checks[checkID]["name"] = checkDoc._source.name;
            }
            return h.response(checks).code(200);
        }
    })

    server.route({
        path: '/api/scorestack/attribute/{id}/{name}',
        method: 'POST',
        handler: async (req, h) => {
            // Make sure value is in request body
            if ("value" in req.payload === false) {
                return h.response({
                    "statusCode": 400,
                    "error": "Bad Request",
                    "message": 'Request body must contain the "value" attribute',
                }).code(400)

            }

            // Make sure the ID is real
            let attribIndices = await dataCluster.callWithRequest(req, 'indices.get', {
                index: `attrib_*_${req.params["id"]}`,
                expand_wildcards: 'open',
            });

            if (Object.keys(attribIndices).length === 0) {
                return h.response({
                    "statusCode": 404,
                    "error": "Not Found",
                    "message": `Attributes for check ID "${req.params["id"]}" either don't exist or you do not have access to them`,
                }).code(404)
            }

            // Check each attribute index for the attribute we are overwriting
            for (let attribIndex of Object.keys(attribIndices)) {
                let attribDoc = await dataCluster.callWithRequest(req, 'get', {
                    id: 'attributes',
                    index: attribIndex,
                });
                // If the attribute exists in the document, update the document with the new value
                if (req.params["name"] in attribDoc._source) {
                    let newAttrib = {};
                    newAttrib[req.params["name"]] = req.payload["value"];
                    let resp = await dataCluster.callWithRequest(req, 'update', {
                        id: 'attributes',
                        index: attribIndex,
                        body: {
                            "doc": newAttrib,
                        },
                    });
                    return h.response({
                        "statusCode": 200,
                        "message": "Attribute updated",
                    }).code(200)
                }
            }
            // If we fall through to here, the attribute was not found
            return h.response({
                "statusCode": 404,
                "error": "Not Found",
                "message": `Attribute "${req.params["name"]}" for check ID ${req.params["id"]} either doesn't exist or you do not have access to it`,
            }).code(404)
        }
    })
}