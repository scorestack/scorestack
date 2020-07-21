import { schema } from '@kbn/config-schema';
import { EuiDataGridBody } from '@elastic/eui/src/components/datagrid/data_grid_body';
import { IRouter, SavedObjectsServiceSetup } from '../../../../src/core/server';

import { PLUGIN_API_BASEURL } from '../../common';

import { SavedTemplateObject } from '../saved_objects';

export function defineRoutes(router: IRouter /* , savedObjects: SavedObjectsServiceSetup*/) {
  router.get(
    {
      path: `${PLUGIN_API_BASEURL}/example`,
      validate: false,
    },
    async (context, request, response) => {
      return response.ok({
        body: {
          time: new Date().toISOString(),
        },
      });
    }
  );

  router.get(
    {
      path: `${PLUGIN_API_BASEURL}/template`,
      validate: {
        query: schema.object({
          id: schema.string(),
        }),
      },
      options: {
        tags: ['access:template_management-read'],
      },
    },
    async (context, request, response) => {
      /*
      const client = savedObjects.getScopedClient(request);
      const template = await client.get('template', request.query.id);
      */
      const client = context.scorestack.getTemplatesClient();

      return response.ok({
        body: await client.get(request.query.id),
      });
    }
  );

  router.post(
    {
      path: `${PLUGIN_API_BASEURL}/template`,
      validate: {
        // TODO: fix this schema object
        body: schema.object({
          id: schema.string(),
          title: schema.string(),
          description: schema.string(),
          protocol: schema.string(),
          score_weight: schema.number(),
          definition: schema.recordOf(schema.string(), schema.any()),
        }),
      },
      options: {
        tags: ['access:template_management-admin'],
      },
    },
    async (context, request, response) => {
      /*
      const client = savedObjects.getScopedClient(request);
      const resp = await client.create('template', { ...request.body });
      */
      const client = context.scorestack.getTemplatesClient();
      const res = await client.create('template', { ...request.body });

      return response.ok({
        body: res,
      });
    }
  );
}
