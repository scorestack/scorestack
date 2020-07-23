import { SavedObjectsType } from '../../../../src/core/server';

export const SavedTemplateObject: SavedObjectsType = {
  name: 'template',
  hidden: false,
  namespaceAgnostic: false,
  mappings: {
    properties: {
      title: {
        type: 'text',
      },
      description: {
        type: 'text',
      },
      protocol: {
        type: 'keyword',
      },
      score_weight: {
        type: 'integer',
      },
      definition: {
        type: 'object',
        dynamic: 'true',
      },
    },
  },
};
