export const PLUGIN_ID = 'scorestack';
export const PLUGIN_NAME = 'scorestack';

export interface CheckAttributes {
  [index: string]: {
    [index: string]: {
      name: string;
      attributes: {
        [index: string]: string;
      };
    };
  };
}
