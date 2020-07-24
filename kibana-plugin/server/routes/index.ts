import { schema } from '@kbn/config-schema';
import {
  IRouter,
  SavedObjectsServiceSetup,
  RequestHandlerContext,
  SavedObjectsClient,
  SavedObject,
  SavedObjectsFindResponse,
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
      path: `${PLUGIN_API_BASEURL}/template/{id}`,
      validate: {
        params: schema.object({
          id: schema.maybe(schema.string()),
        }),
      },
      options: {
        tags: ['access:template_management-read'],
      },
    },
    async (context: ScoreStackContext, request, response) => {
      const client = context.scorestack.getTemplatesClient();
      const getTemplateById: boolean = request.params.id === undefined;

      let savedObjects: Array<SavedObject<TemplateRaw>>;
      try {
        if (getTemplateById) {
          // If no ID was specified, just get all the templates
          const resp: SavedObjectsFindResponse<TemplateRaw> = await client.find({
            type: 'template',
          });
          savedObjects = resp.saved_objects;
        } else {
          // Only get the template of the specified ID
          const resp: SavedObject<TemplateRaw> = await client.get('template', request.params.id);
          savedObjects = [resp];
        }
      } catch (err) {
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

      // If the client requested a template by ID, don't return it in an array
      let respBody: Template | Template[];
      if (getTemplateById) {
        respBody = templateFromSaved(savedObjects[0]);
      } else {
        respBody = savedObjects.map((obj) => templateFromSaved(obj));
      }

      // Return the template(s)
      return response.ok({
        body: JSON.stringify(respBody),
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

      let res: SavedObject<TemplateSavedObject>;
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
