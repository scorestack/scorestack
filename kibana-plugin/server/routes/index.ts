import { schema } from '@kbn/config-schema';
import {
  IRouter,
  SavedObjectsServiceSetup,
  RequestHandlerContext,
  SavedObjectsClient,
  SavedObject,
  SavedObjectsFindResponse,
  SavedObjectsFindOptions,
} from '../../../../src/core/server';

import { PLUGIN_API_BASEURL } from '../../common';
import { Template, TemplateRaw } from '../../common/types';
import { Protocol } from '../../common/checks';

import { SavedTemplateObject } from '../saved_objects';

interface ScoreStackContext extends RequestHandlerContext {
  scorestack: {
    getTemplatesClient(): SavedObjectsClient;
  };
}

function templateFromSaved(saved: SavedObject<TemplateRaw>): Template {
  return {
    id: saved.id,
    protocol: Protocol[saved.attributes.protocol],
    ...saved.attributes,
  };
}

function handleClientError(err, response) {
  const payload = err.output.payload;

  // Determine if we can accurately report the error to the client
  if (payload.statusCode === 404) {
    return response.notFound({ body: payload });
  } else {
    // This will dump an error log to Kibana, and I'm not sure if that's ideal
    return response.internalError({
      body: err,
    });
  }
}

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
          page: schema.maybe(schema.number({ min: 1 })),
          perPage: schema.maybe(schema.number({ min: 1 })),
        }),
      },
      options: {
        tags: ['access:template_management-read'],
      },
    },
    async (context: ScoreStackContext, request, response) => {
      const client = context.scorestack.getTemplatesClient();

      // Set any optional values for the object find
      const options: SavedObjectsFindOptions = {
        type: 'template',
      };
      if (request.query.page !== undefined) {
        options.page = request.query.page;
      }
      if (request.query.perPage !== undefined) {
        options.perPage = request.query.perPage;
      }

      // Get the templates from the saved objects API
      let resp: SavedObjectsFindResponse<TemplateRaw>;
      try {
        resp = await client.find(options);
      } catch (err) {
        return handleClientError(err, response);
      }

      return response.ok({
        body: {
          total: resp.total,
          page: resp.page,
          per_page: resp.per_page,
          templates: resp.saved_objects.map((obj) => templateFromSaved(obj)),
        },
      });
    }
  );

  router.get(
    {
      path: `${PLUGIN_API_BASEURL}/template/{id}`,
      validate: {
        params: schema.object({
          id: schema.string(),
        }),
      },
      options: {
        tags: ['access:template_management-read'],
      },
    },
    async (context: ScoreStackContext, request, response) => {
      const client = context.scorestack.getTemplatesClient();

      // Get the template
      let template: Template;
      try {
        const resp: SavedObject<TemplateRaw> = await client.get('template', request.params.id);
        template = templateFromSaved(resp);
      } catch (err) {
        return handleClientError(err, response);
      }

      return response.ok({
        body: template,
      });
    }
  );

  router.post(
    {
      path: `${PLUGIN_API_BASEURL}/template`,
      validate: {
        // TODO: fix this schema object
        body: schema.object({
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
    async (context: ScoreStackContext, request, response) => {
      // Validate the protocol
      const protocol: Protocol = Protocol[request.body.protocol];
      if (protocol === undefined) {
        return response.badRequest({
          body: {
            message: `'${request.body.protocol}' is not a valid protocol`,
          },
        });
      }

      const client = context.scorestack.getTemplatesClient();

      let res: SavedObject<TemplateRaw>;
      try {
        res = await client.create('template', { ...request.body });
      } catch (err) {
        return response.internalError({
          body: err,
        });
      }

      return response.ok({
        body: res,
      });
    }
  );
}
