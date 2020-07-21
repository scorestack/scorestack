import { schema } from '@kbn/config-schema';
import { IRouter } from '../../../../src/core/server';

import { PLUGIN_API_BASEURL } from '../../common';

export function defineRoutes(router: IRouter) {
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
      return response.ok({
        body: {
          message: `yeet ${request.query.id}`,
        },
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
      return response.ok({
        body: {
          message: `yeet ${JSON.stringify(request.body)}`,
        },
      });
    }
  );
}
