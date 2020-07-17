import { ITemplate } from '../../../common/types';

export enum ActionType {
  Copy,
  Remove,
  Save,
}

export interface TemplateAction {
  type: ActionType;
  template: ITemplate;
}
