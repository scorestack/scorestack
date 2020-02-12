export default function (server, dataCluster) {
    server.route({
        path: '/api/scorestack/attributes',
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
                checks[checkID]["attributes"] = Object.assign(attribDoc._source, checks[checkID]["attributes"]);
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
            return checks;
        }
    })
}