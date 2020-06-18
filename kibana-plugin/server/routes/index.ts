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
}
