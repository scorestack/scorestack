import { SavedObjectsType } from '../../../../src/core/server';

export const SavedTemplateObject: SavedObjectsType = {
  name: 'template',
  hidden: false,
  namespaceAgnostic: false,
  mappings: {
    properties: {
      id: {
        type: 'text',
      },
      title: {
        type: 'text',
      },
      description: {
        type: 'text',
      },
      protocol: {
        type: 'keyword',
      },
    },
  },
};
