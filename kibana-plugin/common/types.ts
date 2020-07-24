import { Protocol } from './checks/protocol';

export interface ITemplate {
  id: string;
  title: string;
  description: string;
  protocol: Protocol;
}

export interface Template {
  id: string;
  title: string;
  description: string;
  protocol: Protocol;
  score_weight: number;
  definition: {
    [index: string]: any;
  };
}

export interface TemplateRaw {
  title: string;
  description: string;
  protocol: string;
  score_weight: number;
  definition: {
    [index: string]: any;
  };
}
