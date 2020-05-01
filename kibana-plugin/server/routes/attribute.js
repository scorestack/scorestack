export default function (server, dataCluster) {
  server.route({
    path: '/api/scorestack/attribute',
    method: 'GET',
    handler: async (req, h) => {
      // All attributes will be returned in a single object
      const checks = {};

      // Get all attribute indexes
      const attribIndices = await dataCluster.callWithRequest(req, 'indices.get', {
        index: 'attrib_*',
        expand_wildcards: 'open',
      });

      // Get attributes for each check
      for (const attribIndex of Object.keys(attribIndices)) {
        // Check how many documents are in the index
        const countDoc = await dataCluster.callWithRequest(req, 'count', {
          index: attribIndex,
        });

        // Search for all documents in the index
        const searchResults = await dataCluster.callWithRequest(req, 'search', {
          index: attribIndex,
          size: countDoc.count,
        });

        // Add each attribute to the object
        for (const check of searchResults.hits.hits) {
          // Parse the document ID to determine the group
          // TODO: don't rely on parsing the document ID or index ID to determine the group, or ensure that unsafe characters are filtered from group names and check names
          const group = check._id.split('-').slice(-1);

          // Set up the checks object to receive the attributes in the right spot
          if (group in checks === false) {
            checks[group] = {};
          }
          if (check._id in checks[group] === false) {
            checks[group][check._id] = {
              'attributes': {},
            };

            // Add check name
            const checkDoc = await dataCluster.callWithRequest(req, 'get', {
              id: check._id,
              index: 'checks',
              _source_includes: 'name',
            });
            checks[group][check._id].name = checkDoc._source.name;
          }

          // Add attribute contents
          checks[group][check._id].attributes = Object.assign(checks[group][check._id].attributes, check._source);
        }
      }

      return h.response(checks).code(200);
    }
  });

  server.route({
    path: '/api/scorestack/attribute/{id}/{name}',
    method: 'POST',
    handler: async (req, h) => {
      // Make sure value is in request body
      if ('value' in req.payload === false) {
        return h.response({
          'statusCode': 400,
          'error': 'Bad Request',
          'message': 'Request body must contain the "value" attribute',
        }).code(400);
      }

      // Parse the group from the ID
      // TODO: don't rely on parsing the document ID or index ID to determine the group, or ensure that unsafe characters are filtered from group names and check names
      const group = req.params.id.split('-').slice(-1);

      // Make sure the group's index exists
      const attribIndices = await dataCluster.callWithRequest(req, 'indices.get', {
        index: `attrib_*_${group}`,
        expand_wildcards: 'open',
      });

      if (Object.keys(attribIndices).length === 0) {
        return h.response({
          'statusCode': 404,
          'error': 'Not Found',
          'message': `Attributes for group "${group}" either don't exist or you do not have access to them`,
        }).code(404);
      }

      // Check each attribute index for the attribute we are overwriting
      for (const attribIndex of Object.keys(attribIndices)) {
        // Try to get the attribute document for the index
        const attribDoc = await dataCluster.callWithRequest(req, 'get', {
          id: req.params.id,
          index: attribIndex,
        });

        // If the attribute exists in the document, update the document with the new value
        if (req.params.name in attribDoc._source) {
          const newAttrib = {};
          newAttrib[req.params.name] = req.payload.value;
          await dataCluster.callWithRequest(req, 'update', {
            id: req.params.id,
            index: attribIndex,
            body: {
              'doc': newAttrib,
            },
          });
          return h.response({
            'statusCode': 200,
            'message': 'Attribute updated',
          }).code(200);
        }
      }
      // If we fall through to here, the attribute was not found
      return h.response({
        'statusCode': 404,
        'error': 'Not Found',
        'message': `Attribute "${req.params.name}" for check ID ${req.params.id} either doesn't exist or you do not have access to it`,
      }).code(404);
    }
  });
}
